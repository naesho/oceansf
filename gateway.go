package main

import (
	"github.com/labstack/echo"
	"github.com/ohsean53/oceansf/context"
	"github.com/ohsean53/oceansf/define"
	"net/http"
	"github.com/ohsean53/oceansf/apperr"
	"github.com/ohsean53/oceansf/retcode"
)


func isLoginApi(api string) bool {
	switch api {
	case "Login":
		return true
	default:
		return false
	}
}

func gateway(c echo.Context) error {

	ctx := c.(*context.SessionContext)

	var api string

	api = ctx.ClientRequest["api"].(string)
	handlerFunc, ok := MsgHandler[api]

	if ok {
		ret, err := handlerFunc(ctx)
		if err != nil {
			appErr := err.(*apperr.Error)
			// panic - recover 됐을 때는 어떻게 하지..
			return c.JSON(http.StatusOK, define.Json{
				"retcode": appErr.ErrorCode(),
				"retmsg":  err.Error(),
			})
		}

		return c.JSON(http.StatusOK, define.Json{
			"retcode": retcode.Success,
			"data":  ret,
		})
	} else {
		return c.String(http.StatusBadRequest, "handler undefined")
	}
}
