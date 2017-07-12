package main

type MsgHandlerFunc func(map[string]interface{}) JsonMap

var msgHandler = map[string]MsgHandlerFunc{
	"ReqLogin": ReqLogin,
}

func ReqLogin(data map[string]interface{}) JsonMap {

	ret := JsonMap{
		"user_id": data["user_id"].(float64),
		"result":  true,
	}
	return ret
}
