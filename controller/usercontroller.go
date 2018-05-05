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

	return user, nil
}

func (UserController) Remove(ctx *context.SessionContext, id string) error {
	// TODO 나중에 session data 에 uid 저장해서 id 없이 uid로 select , delete 가능하게..
	user := model.NewUser(id)
	err := user.Load(ctx)
	if err != nil {
		log.Error(err)
		return err
	}
	return user.Remove(ctx)
}
