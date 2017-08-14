package model

import (
	"github.com/ohsaean/oceansf/config"
	"github.com/ohsaean/oceansf/db"
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

func Load(uid int64) *User {

	// memcached (cache data)

	// when cache fail -> read db
	user := NewUser(uid)
	// cache fail -> select user data from table
	DB := db.GetInstance(config.DB_DEV, false)
	query := "SELECT SQL_NOCACHE * FROM user WHERE uid = ?"

	DB.Prepare(query)

	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!
	rows, err := stmt.Query(uid)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&user.UID, &user.Name, &user.Email, &user.RegisterDate, &user.LastLoginDate)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Info("load user data")
	return user
}

func (u *User) Save() {
	// update or insert (upsert)
}

func (u *User) Remove() {
	// delete data
}
