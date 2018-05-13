package controller

import (
	"github.com/labstack/gommon/log"
	"github.com/ohsean53/oceansf/context"
	"github.com/ohsean53/oceansf/lib"
	"github.com/ohsean53/oceansf/model"
)

type UserController struct {
}

func (UserController) Login(ctx *context.SessionContext, id string, name string, email string) (*model.User, error) {
	user := model.NewUser(id)
	err := user.Load(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if user.UID == 0 {
		// create 유저
		user.Id = id
		user.Name = name
		user.Email = email
		user.RegisterDate  = lib.GetDateTime()
	}

	user.LastLoginDate = lib.GetDateTime()
	err = user.Save(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	
	// session 처리 미들웨어에서 세션을 저장해줌
	ctx.Session.UID = user.UID
	ctx.Session.Id = user.Id
	ctx.Session.Name = user.Name
	ctx.Session.Email = user.Email
	ctx.Session.RegisterDate = user.RegisterDate
	ctx.Session.LastLoginDate = user.LastLoginDate

	return user, nil
}

func (UserController) Remove(ctx *context.SessionContext, id string) error {
	user := model.NewUser(id)
	err := user.Load(ctx)
	if err != nil {
		log.Error(err)
		return err
	}


	err = user.Remove(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	key := context.GetSessionCacheKey(user.UID)
	err = ctx.Cache.Delete(key)

	if err != nil {
		return err
	}

	// 세션이 저장되지 않도록
	ctx.Session.UID = 0

	return nil
}
