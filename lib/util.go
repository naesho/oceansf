package lib

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
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

func GetDateTime() string {
	t := GetNow()
	layoutYmdHis := "2006-01-02 15:04:05"
	return t.Format(layoutYmdHis)
}

func GetDateYmd() string {
	t := GetNow()
	layout := "20060102"
	return t.Format(layout)
}

func Atoi64(s string) int64 {
	integer, err := strconv.Atoi(s)
	CheckError(err)
	return int64(integer)
}

func Atoi32(s string) int32 {
	integer, err := strconv.Atoi(s)
	CheckError(err)
	return int32(integer)
}

func Atoi(s string) int {
	integer, err := strconv.Atoi(s)
	CheckError(err)
	return integer
}

func Itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Itoa32(i int32) string {
	return strconv.Itoa(int(i))
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}

func GetUnixTime() int64 {
	t := GetNow().Unix()
	return t
}

func GetMicroTime() int64 {
	t := GetNow().UnixNano()
	return t
}

type TimerInfo struct {
	CallBack func()
	Delay    time.Duration
	Info     map[string]interface{}
}

func CheckError(err error) {
	if err != nil {
		log.Error(err)
	}
}
