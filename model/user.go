package model

import (
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/naesho/oceansf/cache"
	"github.com/naesho/oceansf/db"
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

func LoadUser(uid int64) (u *User, err error) {

	dbConn := db.Conn
	mc := cache.Mcache

	// memcached (cache data)
	key := getCacheKey(uid)

	var item *memcache.Item
	item, err = mc.Get(key)
	if err != memcache.ErrCacheMiss {
		if err = json.Unmarshal(item.Value, &u); err != nil {
			lib.CheckError(err)
			return nil, err
		}

		log.Debug("cache hit")
		log.Debug(string(item.Value))
		return u, nil
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
	log.Info("load user data")

	// cache data
	data, err := json.Marshal(u)
	if err != nil {
		lib.CheckError(err)
		return nil, err
	}

	item = &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: 10,
	}
	err = mc.Set(item)
	lib.CheckError(err)

	return u, nil
}

func (u *User) Save() {

}

func (u *User) Remove() {

}
