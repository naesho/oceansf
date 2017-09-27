package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/evalphobia/logrus_fluent"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ohsaean/oceansf/db"
	"github.com/ohsaean/oceansf/define"
	"github.com/ohsaean/oceansf/grace"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type (
	Context struct {
		echo.Context
		db  *sql.DB
		req define.JsonMap
	}

	Session struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		SessionKey string `json:"session_key"`
	}

	Stats struct {
		Uptime       time.Time      `json:"uptime"`
		RequestCount uint64         `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		mutex        sync.RWMutex
	}
)

func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: map[string]int{},
	}
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++
		return nil
	}
}

func gateway(c echo.Context) error {

	ctx := c.(*Context)

	var rawJson define.JsonMap
	rawBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
		return err
	}

	err = json.Unmarshal(rawBody, &rawJson)
	if err != nil {
		log.Error(err)
		return err
	}

	var api string
	api = rawJson["api"].(string)
	handlerFunc, ok := MsgHandler[api]
	ctx.req = rawJson

	if ok {
		ret := handlerFunc(c)
		return c.JSON(http.StatusOK, ret)
	} else {
		return c.String(http.StatusNotFound, "api_parsing_error")
	}
}

func CheckError(err error) {
	if err != nil {
		log.Error(err)
	}
}
func (s *Stats) stat(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	// TODO fluent instance 만들기
	hook, err := logrus_fluent.New("localhost", 24224)
	CheckError(err)

	// set custom fire level
	//hook.SetLevels([]log.Level{
	//	log.PanicLevel,
	//	log.ErrorLevel,
	//})

	// set static tag
	hook.SetTag("log.server")

	// ignore field
	hook.AddIgnore("context")

	// filter func
	hook.AddFilter("error", logrus_fluent.FilterError)

	log.AddHook(hook)
}

func main() {
	data, err := ioutil.ReadFile("./config/config.toml")
	CheckError(err)
	fmt.Print(string(data))

	var dbInfo db.Info
	if _, err := toml.Decode(string(data), &dbInfo); err != nil {
		// handle error
	}

	log.Debug(dbInfo)
	dbConn := db.NewDB(&dbInfo)

	// Setup
	e := echo.New()
	s := NewStats()

	// Middleware

	// Create a middleware to extend default context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &Context{c, dbConn, nil}
			return h(ctx)
		}
	})
	e.Use(middleware.Logger()) // like nginx access log
	e.Use(middleware.Recover())
	e.Use(s.Process)

	// Router
	e.GET("/", func(c echo.Context) error {
		cc := c.(*Context)
		return cc.String(http.StatusOK, "Hello :D")
	})
	e.GET("/stat", s.stat)
	e.POST("/gateway", gateway)
	e.Server.Addr = ":8080"

	// Serve it like a boss
	log.Fatal(grace.Serve(e.Server))
}
