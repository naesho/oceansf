package main

import (
	"github.com/naesho/oceansf/define"
	"github.com/naesho/oceansf/model"
)

type MsgHandlerFunc func(req define.Json) (interface{}, error)

var MsgHandler = map[string]MsgHandlerFunc{
	"GetUser": GetUser,
	"GetUserNoWait": GetUserNoWait,
}

func GetUser(req define.Json) (interface{}, error) {
	uid := req["user_id"].(float64)
	return model.LoadUser(int64(uid))
}

func GetUserNoWait(req define.Json) (interface{}, error) {
	uid := req["user_id"].(float64)
	return model.LoadUserNoWait(int64(uid))
}