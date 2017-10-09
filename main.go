package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/evalphobia/logrus_fluent"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ohsaean/oceansf/cache"
	"github.com/ohsaean/oceansf/db"
	"github.com/ohsaean/oceansf/define"
	"github.com/ohsaean/oceansf/grace"
	"github.com/ohsaean/oceansf/lib"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type (
	tomlConfig struct {
		DB        db.Info       `toml:"database"`
		Fluent    FluentInfo    `toml:"fluent"`
		Memcached MemcachedInfo `toml:"memcached"`
	}

	FluentInfo struct {
		Ip   string
		Port int
	}

	MemcachedInfo struct {
		Endpoint string
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

	var req define.Json
	rawBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
		return err
	}

	err = json.Unmarshal(rawBody, &req)
	if err != nil {
		log.Error(err)
		return err
	}

	var api string
	api = req["api"].(string)
	handlerFunc, ok := MsgHandler[api]

	if ok {
		ret, err := handlerFunc(req)
		if err != nil {
			// TODO 나중에 appError{ code, msg, error } 와 같은 별도 에러 구조체 정의
			return c.JSON(http.StatusOK, define.Json{
				"retcode": 1001,
				"retmsg":  err.Error(),
			})
		}
		return c.JSON(http.StatusOK, ret)
	} else {
		return c.String(http.StatusNotFound, "api_parsing_error")
	}
}

func (s *Stats) stat(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}

func init() {

	var config tomlConfig

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	// Application Configuration
	data, err := ioutil.ReadFile("./config/config.toml")
	lib.CheckError(err)

	_, err = toml.Decode(string(data), &config)
	lib.CheckError(err)

	log.Debug("endpoint of memcached : " + config.Memcached.Endpoint)
	log.Debug("endpoint of db : " + config.DB.Ip)

	// Initialize Mysql Client
	db.Init(&config.DB)

	// Initialize ElastiCache (Memcached Cluster) Client
	cache.InitMemcache(config.Memcached.Endpoint)

	// Forward to Fluent (log aggregator)
	hook, err := logrus_fluent.New(config.Fluent.Ip, config.Fluent.Port)
	lib.CheckError(err)

	// set custom fire level
	//hook.SetLevels([]log.Level{
	//	log.PanicLevel,
	//	log.ErrorLevel,
	//})

	// Set static tag
	hook.SetTag("td.log.server")

	// Ignore field
	hook.AddIgnore("context")

	// Filter func
	hook.AddFilter("error", logrus_fluent.FilterError)

	log.AddHook(hook)
}

func main() {

	// Setup
	e := echo.New()
	s := NewStats()

	// Middleware
	e.Use(middleware.Logger()) // like nginx access log
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
