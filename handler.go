package main

import (
	"github.com/ohsean53/oceansf/context"
	"github.com/ohsean53/oceansf/controller"
	"github.com/ohsean53/oceansf/define"
)

type MsgHandlerFunc func(ctx *context.SessionContext) (interface{}, error)

var MsgHandler = map[string]MsgHandlerFunc{
	"Login":      Login,
	"RemoveUser": RemoveUser,
}

func Login(ctx *context.SessionContext) (interface{}, error) {
	id := ctx.ClientRequest["id"].(string)
	name := ctx.ClientRequest["name"].(string)
	email := ctx.ClientRequest["email"].(string)
	uc := controller.UserController{}
	return uc.Login(ctx, id, name, email)
}

func RemoveUser(ctx *context.SessionContext) (interface{}, error) {
	id := ctx.Session.Id
	uc := controller.UserController{}
	err := uc.Remove(ctx, id)
	return define.Json{}, err
}
