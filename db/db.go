package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ohsaean/oceansf/lib"
)

type Info struct {
	Ip     string
	Port   string
	DbName string
	User   string
	Pass   string
}

func NewDB(dbInfo *Info) *sql.DB {
	dsn := dbInfo.User + ":" + dbInfo.Pass + "@tcp(" + dbInfo.Ip + ":" + dbInfo.Port + ")/" + dbInfo.DbName + "?charset=utf8"
	instance, err := sql.Open("mysql", dsn)
	lib.CheckError(err)
	return instance
}
