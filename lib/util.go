package lib

import (
	"fmt"
	"math/rand"
	"time"
)

func RandInt64(min int64, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandInt32(min int32, max int32) int32 {
	return min + rand.Int31n(max-min)
}

func GetNow() time.Time {
	t := time.Now()
	return t
}

func GetUnixTime() int64 {
	t := GetNow().Unix()
	return t
}

func GetMicroTime() int64 {
	t := GetNow().UnixNano()
	return t
}

func GetDateTime() string {
	t := GetNow()
	//t.Format("2016-01-01 23:59:59") timezone 안 없어짐
	y, m, d := t.Date()
	h, i, s := t.Clock()
	ymdhis := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", y, m, d, h, i, s)
	return ymdhis
}

type TimerInfo struct {
	CallBack func()
	Delay    time.Duration
	Info     map[string]interface{}
}
