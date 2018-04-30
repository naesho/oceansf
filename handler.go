package main

import (
	"github.com/naesho/oceansf/define"
	"github.com/naesho/oceansf/model"
	"github.com/naesho/oceansf/context"
)

type MsgHandlerFunc func(req define.Json, reqCtx *context.RequestContext) (interface{}, error)

var MsgHandler = map[string]MsgHandlerFunc{
	"GetUser": GetUser,
	"GetUserNoWait": GetUserNoWait,
}

func GetUser(req define.Json, reqCtx *context.RequestContext) (interface{}, error) {
	uid := req["user_id"].(float64)
	return model.LoadUser(int64(uid), reqCtx)
}

func GetUserNoWait(req define.Json, reqCtx *context.RequestContext) (interface{}, error) {
	uid := req["user_id"].(float64)
	return model.LoadUserNoWait(int64(uid), reqCtx)
}