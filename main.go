package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/evalphobia/logrus_fluent"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/naesho/oceansf/cache"
	"github.com/naesho/oceansf/context"
	"github.com/naesho/oceansf/db"
	"github.com/naesho/oceansf/define"
	"github.com/naesho/oceansf/grace"
	"github.com/naesho/oceansf/lib"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

	Session struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		SessionKey string `json:"session_key"`
		// 국가
		// 언어
		// DID
		// 앱버전
		// 세션 갱신 시간
		// 로그인 날짜
		// 로그인 국가
		// 마켓
		// 등등등 ... etc
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

func SessionCheckMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO 세션 체크는 여기서
		log.Debug("TODO : Session Check ")
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

func gateway(c echo.Context) error {

	reqCtx := c.(*context.RequestContext)

	var req define.Json
	rawBody, err := ioutil.ReadAll(reqCtx.Request().Body)
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
	uid := req["user_id"].(float64)

	handlerFunc, ok := MsgHandler[api]
	if ok {
		lockKey := cache.GetGlobalLockKey(int64(uid))

		// 중복 처리 방지를 위한 memcached add lock
		lockError := reqCtx.Cache.Lock(lockKey)
		if lockError != nil {
			return lockError
		}
		defer reqCtx.Cache.UnLock(lockKey)
		ret, err := handlerFunc(req, reqCtx)
		if err != nil {
			// TODO 나중에 appError{ code, msg, error } 와 같은 별도 에러 구조체 정의
			// TODO db rollback, memcache discard
			// panic - recover 됐을 때는 어떻게 하지..
			return c.JSON(http.StatusOK, define.Json{
				"retcode": 1001,
				"retmsg":  err.Error(),
			})
		}

		// 성공
		// TODO memcache cas & delayed 처리?
		// DB 트랜잭션이 성공하고 나서 memcached 에 cas 를 통해 반영해야함 --> 지금은 일단 handlerFunc 에서 처리.
		// memcache 처리 goroutine 띄워서 성공시에 cas 하도록 ? --> channel 을 통해서 완료 시그널 주도록..
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

func CustomContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 매번 connection 생성..
		cc := &context.RequestContext{
			Context: c,
			DB:      db.NewConnection(&config.DB),
			Cache:   cache.NewConnection(config.Memcached.Endpoint),
		}

		// 커넥션 정리될때 같이 close
		defer cc.DB.Close()
		return next(cc)
	}
}

func main() {

	// Setup
	e := echo.New()
	s := NewStats()

	// Middleware before router
	e.Use(CustomContextMiddleware)

	// TODO 프로파일러
	// TODO 서버 상태 체크
	e.Use(SessionCheckMiddleWare)
	// TODO 요청 파라메터 검사
	// TODO 아이피 대역 체크?
	// TODO 버전 체크 (리소스나 클라-서버간 맞춰야 하는 데이터 버전들)
	// TODO 제재 같은거도 체크?

	// Middleware after router
	e.Use(middleware.Logger()) // like nginx access log
	e.Use(middleware.Recover())
	e.Use(s.Process)

	// Router
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello :D")
	})
	e.GET("/stat", s.stat)
	e.POST("/gateway", gateway)
	e.Server.Addr = ":5555"

	// Serve it like a boss
	log.Fatal(grace.Serve(e.Server))
}
