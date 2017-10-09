package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"github.com/ohsaean/oceansf/lib"
)

type Info struct {
	Ip     string
	Port   string
	DbName string
	User   string
	Pass   string
}

var (
	Conn *sql.DB
)

// Initialize instance for using main package
func Init(dbInfo *Info) {
	var err error
	dsn := dbInfo.User + ":" + dbInfo.Pass + "@tcp(" + dbInfo.Ip + ":" + dbInfo.Port + ")/" + dbInfo.DbName + "?charset=utf8"
	Conn, err = sql.Open("mysql", dsn)
	log.Debug("db init")
	lib.CheckError(err)

	err = Conn.Ping()
	lib.CheckError(err)
}
