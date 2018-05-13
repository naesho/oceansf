package context

import (
	"github.com/ohsean53/oceansf/define"
	"github.com/ohsean53/oceansf/lib"
)

type Session struct {
	UID           int64  `json:"uid"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	RegisterDate  string `json:"register_date"`
	LastLoginDate string `json:"last_login_date"`

	// 국가
	// 언어
	// DID
	// 앱버전
	// 세션 갱신 시간
	// 로그인 날짜
	// 로그인 국가
	// 마켓
	// 등등등 ... etc
}

func GetSessionCacheKey(uid int64) string {
	return define.MemcachePrefix + "session_key:" + lib.Itoa64(uid)
}

func NewSession(uid int64, id, name, email, registerDate, lastLoginDate  string) *Session{
	return &Session{
		UID:           uid,
		Id:            id,
		Name:          name,
		Email:         email,
		RegisterDate:  registerDate,
		LastLoginDate: lastLoginDate,
	}
}