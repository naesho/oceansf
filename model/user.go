package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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

func NewUser(UID int64) {
	// crate object
}

func Load(uid int64) *User {

	// memcached (cache data)

	// cache fail -> select user data from table
	log.Info("load user data")
}

func (u *User) Save() {
	// update or insert (upsert)
}

func (u *User) Remove() {
	// delete data
}
