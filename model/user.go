package model

import (
	"encoding/json"
	"github.com/ohsean53/oceansf/cache"
	"github.com/ohsean53/oceansf/context"
	"github.com/ohsean53/oceansf/define"
	"github.com/ohsean53/oceansf/lib"
	log "github.com/sirupsen/logrus"
)

// database crud object

type User struct {
	UID           int64
	Id            string
	Name          string
	Email         string
	RegisterDate  string
	LastLoginDate string
}

func NewUser(Id string) *User {
	// create object
	return &User{
		UID:           0,
		Id:            Id,
		Name:          "unknown",
		Email:         "unknown@github.com",
		RegisterDate:  "2017-01-01 00:00:00",
		LastLoginDate: "2017-01-01 00:00:00",
	}
}

func getCacheKey(id string) string {
	return define.MemcachePrefix + "user:" + id
}

func (u *User) Load(ctx *context.SessionContext) (err error) {

	//time.Sleep(time.Second * 5)

	dbConn := ctx.DB
	mc := ctx.Cache

	// memcached (cache data)
	key := getCacheKey(u.Id)

	var cacheData []byte
	cacheData, err = mc.Get(key)

	log.Debug("cache get")
	if err == nil {
		if cacheData != nil {
			if err = json.Unmarshal(cacheData, &u); err != nil {
				lib.CheckError(err)
				return err
			}

			log.Debug("cache hit")
			log.Debug(u)
			return nil
		}
	} else {
		log.Error(err)
	}

	// when cache fail -> read db
	// cache fail -> select user data from table
	query := "SELECT SQL_NO_CACHE * " +
		"       FROM USER" +
		"      WHERE id = ?"

	stmt, err := dbConn.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!
	rows, err := stmt.Query(u.Id)
	lib.CheckError(err)

	for rows.Next() {
		err = rows.Scan(&u.UID, &u.Id, &u.Name, &u.Email, &u.RegisterDate, &u.LastLoginDate)
		lib.CheckError(err)
	}
	log.Info("load user data from DB :", u.Id)

	// cache data
	data, err := json.Marshal(u)
	if err != nil {
		lib.CheckError(err)
		return err
	}

	mc.CasDelayed(key, data, cache.EXPIRE)
	lib.CheckError(err)

	return nil
}

func (u *User) Save(ctx *context.SessionContext) error {
	query := "UPDATE USER " +
		"        SET name = ?," +
		"            email = ?," +
		"            register_date = ?," +
		"            last_login_date = ?" +
		"      WHERE id = ?"

	stmt, err := ctx.DB.Prepare(query)
	lib.CheckError(err)

	defer stmt.Close() // danger!
	res, err := stmt.Exec(u.Name, u.Email, u.RegisterDate, u.LastLoginDate, u.Id)
	lib.CheckError(err)

	affectedRow, err := res.RowsAffected()
	if err == nil && affectedRow == 0 && u.UID == 0 {
		query := "INSERT INTO USER (id, name, email, register_date, last_login_date) VALUES (?,?,?,?,?)"
		stmt, err := ctx.DB.Prepare(query)
		lib.CheckError(err)

		res, err = stmt.Exec(u.Id, u.Name, u.Email, u.RegisterDate, u.LastLoginDate)
		lib.CheckError(err)
		if err == nil {
			NewUID, err2 := res.LastInsertId()
			lib.CheckError(err2)
			if err2 == nil {
				u.UID = NewUID
			}
		}

	}

	key := getCacheKey(u.Id)

	// cache data
	data, err := json.Marshal(u)
	if err != nil {
		lib.CheckError(err)
		return err
	}

	ctx.Cache.CasDelayed(key, data, cache.EXPIRE)

	return nil
}

func (u *User) Remove(ctx *context.SessionContext) error {
	query := "DELETE FROM USER WHERE uid = ?"

	stmt, err := ctx.DB.Prepare(query)
	lib.CheckError(err)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(u.UID)
	lib.CheckError(err)

	if err != nil {
		return err
	}

	key := getCacheKey(u.Id)
	err = ctx.Cache.Delete(key)

	if err != nil {
		return err
	}

	return nil
}
