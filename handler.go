package main

import (
	"github.com/labstack/echo"
	"github.com/ohsaean/oceansf/define"
	"github.com/ohsaean/oceansf/model"
)

type MsgHandlerFunc func(echo.Context) *define.JsonMap

var MsgHandler = map[string]MsgHandlerFunc{
	"GetUser": GetUser,
}

func GetUser(c echo.Context) *define.JsonMap {
	ctx := c.(*Context)
	db := ctx.db
	uid := ctx.req["user_id"].(float64)
	return model.Load(db, int64(uid))
}
