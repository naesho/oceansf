package main

import (
	"github.com/BurntSushi/toml"
	"github.com/evalphobia/logrus_fluent"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ohsean53/oceansf/db"
	"github.com/ohsean53/oceansf/grace"
	"github.com/ohsean53/oceansf/lib"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

var config tomlConfig

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

	// Application Configuration
	data, err := ioutil.ReadFile("./config/config.toml")
	lib.CheckError(err)

	_, err = toml.Decode(string(data), &config)
	lib.CheckError(err)

	// Initialize Mysql Client
	//db.Init(&config.DB)

	// Initialize Memcached Client
	//cache.InitMemcache(config.Memcached.Endpoint)

	// Forward to Fluent (log aggregator)
	hook, err := logrus_fluent.NewWithConfig(logrus_fluent.Config{
		Host: config.Fluent.Ip,
		Port: config.Fluent.Port,
	})
	lib.CheckError(err)

	// set custom fire level
	hook.SetLevels([]log.Level{
		log.PanicLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
		log.DebugLevel,
	})

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

	// 미들웨어 처리 순서
	// https://medium.com/veltra-engineering/echo-middleware-in-golang-90e1d301eb27

	e.Use(CustomContextMiddleware)
	// TODO 프로파일러
	// TODO 서버 상태 체크
	e.Use(SessionCheckMiddleWare)
	// TODO 요청 파라메터 검사
	// TODO 아이피 대역 체크?
	// TODO 버전 체크 (리소스나 클라-서버간 맞춰야 하는 데이터 버전들)
	// TODO 제재 같은거도 체크?
	e.Use(s.ProcessMiddleware)
	e.Use(middleware.Logger()) // like nginx access log
	e.Use(middleware.Recover())

	// Router
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello :D")
	})
	e.GET("/stat", s.stat)
	e.POST("/gateway", gateway) // 실제로는 이거 한개 씀
	e.Server.Addr = ":5555"

	// Graceful restart (SIGUSR2 을 받으면 커넥션 끊김 없이 재시작함)
	log.Fatal(grace.Serve(e.Server))
}
