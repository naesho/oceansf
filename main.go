package main

import (
	"encoding/json"
	"github.com/evalphobia/logrus_fluent"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	JsonMap map[string]interface{}

	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
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

var (
	users = map[int]*user{}
	seq   = 1
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

	var rawJson map[string]interface{}
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
	handler, ok := msgHandler[api]

	if ok {
		ret := handler(rawJson)
		return c.JSON(http.StatusOK, ret)
	} else {
		return c.String(http.StatusInternalServerError, "api_parse_error")
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

	hook, err := logrus_fluent.New("localhost", 24224)
	if err != nil {
		panic(err)
	}

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
	// Setup
	e := echo.New()
	s := NewStats()

	// Middleware
	e.Use(middleware.Logger()) // like Nginx access log
	e.Use(middleware.Recover())
	e.Use(s.Process)

	// Router
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello :D")
	})
	e.GET("/stat", s.stat)
	e.POST("/gateway", gateway)
	e.Server.Addr = ":8080"

	// Serve it like a boss
	log.Fatal(grace.Serve(e.Server))
}
