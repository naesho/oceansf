package main

import (
	"github.com/labstack/echo"
	"github.com/ohsean53/oceansf/context"
	"io/ioutil"
	"github.com/labstack/gommon/log"
	"encoding/json"
	"github.com/ohsean53/oceansf/cache"
	"github.com/ohsean53/oceansf/db"
	"github.com/ohsean53/oceansf/define"
	"errors"
	"strconv"
)



func SessionCheckMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO 세션 체크는 여기서
		log.Debug("TODO : Session Check ")

		ctx := c.(*context.SessionContext)

		// 로그인 api 인지 검사  해야 함.
		// 로그인 전 단계라면 세션이 아직 생성되지 않은 상태라서 cache 에 get 하는대신 set 을 먼저 해야함.

		api := ctx.ClientRequest["api"].(string)
		var session context.Session
		if isLoginApi(api) == false {
			uid := int64(ctx.ClientRequest["uid"].(float64))

			sessionKey := context.GetSessionCacheKey(uid)
			log.Debug("load session uid : ", uid)
			log.Debug("load session key : ", sessionKey)

			rawData, err := ctx.Cache.Get(sessionKey)
			if err != nil {
				// memcache 에 데이터가 없다.
				log.Info("session expired")
				log.Error(err)
				return err // TODO 세션 만료 error 정의 해야함
			}

			err = json.Unmarshal(rawData, &session)
			if err != nil {
				log.Error("invalid session data in memcached")
				return err // 세션에 잘못된 데이터가 들어감
			}

			if session.UID != uid {
				log.Error("invalid session data, not match UID")
				return errors.New("invalid session data, not match UID") // UID가 일치 하지 않음
			}

			ctx.Session = &session
		}

		if err := next(c); err != nil {
			c.Error(err)
		}

		if ctx.Session.UID > 0 {
			// 세션 갱신은 게이트웨이 나가고 나서 처리..
			uid := ctx.Session.UID
			sessionKey := context.GetSessionCacheKey(uid)
			log.Debug("refresh session key :", sessionKey)

			log.Debug("save session", ctx.Session)

			data, err := json.Marshal(ctx.Session)
			if err != nil {
				return err
			}
			err = ctx.Cache.Set(sessionKey, data, define.SessionExpire)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func CustomContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 매번 connection 생성..
		cc := &context.SessionContext{
			Context: c,
			DB:      db.NewConnection(&config.DB),
			Cache:   cache.NewConnection(config.Memcached.Endpoint),
			Session: &context.Session{},
		}

		var req define.Json
		rawBody, err := ioutil.ReadAll(cc.Request().Body)
		if err != nil {
			log.Error(err)
			return err
		}

		err = json.Unmarshal(rawBody, &req)
		if err != nil {
			log.Error(err)
			return err
		}

		// 클라 요청 파싱 (나중에 별도 미들웨어로?)
		cc.ClientRequest = req

		var api string

		api = cc.ClientRequest["api"].(string)

		lockKey := ""
		if isLoginApi(api) {
			id := cc.ClientRequest["id"].(string)
			lockKey = cache.GetGlobalLockKeyWithId(id)
		} else {
			uid := cc.ClientRequest["uid"].(float64)
			lockKey = cache.GetGlobalLockKey(int64(uid))
		}

		// 중복 처리 방지를 위한 memcached add lock
		lockError := cc.Cache.Lock(lockKey)
		if lockError != nil {
			return lockError
		}
		defer cc.Cache.UnLock(lockKey)

		// 커넥션 정리될때 같이 close
		// 일단 트랜잭션 처리는 해당 핸들러 내에서 담당 하는걸로..
		defer cc.DB.Close()
		err = next(cc)

		// defer 직전에 처리되는 부분
		if err != nil {
			// 에러 상황임
			// cache 반영하지 않는다
			cc.Cache.DiscardAll()
			return err
		}

		// 큐잉 되었던 cache 반영
		cc.Cache.CommitAll()

		return nil
	}
}

func (s *Stats) ProcessMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
