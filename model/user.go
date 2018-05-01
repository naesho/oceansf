package model

import (
	"encoding/json"
	"github.com/naesho/oceansf/cache"
	"github.com/naesho/oceansf/context"
	"github.com/naesho/oceansf/define"
	"github.com/naesho/oceansf/lib"
	log "github.com/sirupsen/logrus"
)

// database crud object

type User struct {
	UID           int64
	Name          string
	Email         string
	RegisterDate  string
	LastLoginDate string
}

func NewUser(UID int64) *User {
	// crate object
	return &User{
		UID:           UID,
		Name:          "unknown",
		Email:         "unknown@github.com",
		RegisterDate:  "2017-01-01 00:00:00",
		LastLoginDate: "2017-01-01 00:00:00",
	}
}

func getCacheKey(uid int64) string {
	return define.MemcachePrefix + "user:" + lib.Itoa64(uid)
}

func Login(uid int64, reqCtx *context.RequestContext) (u *User, err error) {
	u, err = LoadUser(uid, reqCtx)

	u.LastLoginDate = lib.GetDateTime()

	// cache data
	data, err := json.Marshal(u)
	if err != nil {
		lib.CheckError(err)
		return u, err
	}

	query :=
		"UPDATE USER " +
		"   SET last_login_date = ?" +
		"   WHERE uid = ?"

	stmt, err := reqCtx.DB.Prepare(query)
	lib.CheckError(err)
	defer stmt.Close() // danger!
	_, err = stmt.Exec(u.LastLoginDate, u.UID)
	lib.CheckError(err)

	key := getCacheKey(uid)
	reqCtx.Cache.CasDelayed(key, data, cache.EXPIRE)

	return u, nil
}

func LoadUser(uid int64, reqCtx *context.RequestContext) (u *User, err error) {

	//time.Sleep(time.Second * 5)

	dbConn := reqCtx.DB
	mc := reqCtx.Cache

	// memcached (cache data)
	key := getCacheKey(uid)

	var cacheData []byte
	cacheData, err = mc.Get(key)

	log.Debug("cache get")
	if err == nil {
		if cacheData != nil {
			if err = json.Unmarshal(cacheData, &u); err != nil {
				lib.CheckError(err)
				return nil, err
			}

			log.Debug("cache hit")
			log.Debug(u)
			return u, nil
		}
	} else {
		log.Error(err)
	}

	// when cache fail -> read db
	u = NewUser(uid)
	// cache fail -> select user data from table
	query := "SELECT SQL_NO_CACHE * FROM USER WHERE uid = ?"

	stmt, err := dbConn.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!
	rows, err := stmt.Query(uid)
	lib.CheckError(err)

	for rows.Next() {
		err = rows.Scan(&u.UID, &u.Name, &u.Email, &u.RegisterDate, &u.LastLoginDate)
		lib.CheckError(err)
	}
	log.Info("load user data from DB")

	// cache data
	data, err := json.Marshal(u)
	if err != nil {
		lib.CheckError(err)
		return nil, err
	}

	mc.CasDelayed(key, data, cache.EXPIRE)
	lib.CheckError(err)

	return u, nil
}

func (u *User) Save() {

}

func (u *User) Remove() {

}
