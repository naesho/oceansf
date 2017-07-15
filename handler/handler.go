package handler

import "github.com/ohsaean/oceansf/define"

type MsgHandlerFunc func(map[string]interface{}) define.JsonMap

var MsgHandler = map[string]MsgHandlerFunc{
	"ReqLogin": ReqLogin,
}

func ReqLogin(data map[string]interface{}) define.JsonMap {

	ret := define.JsonMap{
		"user_id": data["user_id"].(float64),
		"result":  true,
	}
	return ret
}
