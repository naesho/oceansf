package main

import (
	"github.com/ohsaean/oceansf/define"
	"github.com/ohsaean/oceansf/model"
)

type MsgHandlerFunc func(req define.Json) (interface{}, error)

var MsgHandler = map[string]MsgHandlerFunc{
	"GetUser": GetUser,
}

func GetUser(req define.Json) (interface{}, error) {
	uid := req["user_id"].(float64)
	return model.LoadUser(int64(uid))
}
