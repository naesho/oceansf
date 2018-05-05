package main

import (
	"github.com/ohsean53/oceansf/context"
	"github.com/ohsean53/oceansf/controller"
	"github.com/ohsean53/oceansf/define"
)

type MsgHandlerFunc func(req define.Json, ctx *context.SessionContext) (interface{}, error)

var MsgHandler = map[string]MsgHandlerFunc{
	"Login":      Login,
	"RemoveUser": RemoveUser,
}

func Login(req define.Json, ctx *context.SessionContext) (interface{}, error) {
	id := req["id"].(string)
	name := req["name"].(string)
	email := req["email"].(string)
	uc := controller.UserController{}
	return uc.Login(ctx, id, name, email)
}

func RemoveUser(req define.Json, ctx *context.SessionContext) (interface{}, error) {
	id := req["id"].(string)
	uc := controller.UserController{}
	err := uc.Remove(ctx, id)
	return define.Json{
		"retcode": 100,
	}, err
}
