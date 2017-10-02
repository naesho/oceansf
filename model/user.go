package model

import (
	"database/sql"
	"github.com/ohsaean/oceansf/define"
	"github.com/ohsaean/oceansf/lib"
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
		Email:         "nonmae@github.com",
		RegisterDate:  "2017-01-01 00:00:00",
		LastLoginDate: "2017-01-01 00:00:00",
	}
}

func Load(db *sql.DB, uid int64) *define.JsonMap {

	// memcached (cache data)

	// when cache fail -> read db
	user := NewUser(uid)
	// cache fail -> select user data from table
	query := "SELECT SQL_NO_CACHE * FROM USER WHERE user_id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!
	rows, err := stmt.Query(uid)
	lib.CheckError(err)

	for rows.Next() {
		err = rows.Scan(&user.UID, &user.Name, &user.Email, &user.RegisterDate, &user.LastLoginDate)
		lib.CheckError(err)
	}
	log.Info("load user data")

	ret := &define.JsonMap{
		"user":    user,
		"retcode": 100,
	}

	return ret
}

func (u *User) Save() {
	// update or insert (upsert)
}

func (u *User) Remove() {
	// delete data
}
