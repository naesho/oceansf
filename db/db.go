package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"github.com/ohsean53/oceansf/lib"
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

func Init(dbInfo *Info) {
	var err error
	dsn := dbInfo.User + ":" + dbInfo.Pass + "@tcp(" + dbInfo.Ip + ":" + dbInfo.Port + ")/" + dbInfo.DbName + "?charset=utf8"
	Conn, err = sql.Open("mysql", dsn)
	log.Debug("db init")
	lib.CheckError(err)

	err = Conn.Ping()
	lib.CheckError(err)
}

type DB struct {
	*sql.DB
}

func NewConnection(dbInfo *Info) *DB {
	var err error
	dsn := dbInfo.User + ":" + dbInfo.Pass + "@tcp(" + dbInfo.Ip + ":" + dbInfo.Port + ")/" + dbInfo.DbName + "?charset=utf8"
	Conn, err = sql.Open("mysql", dsn)
	log.Debug("create new mysql connection")
	lib.CheckError(err)
	return &DB{
		Conn,
	}
}

func (db *DB) SomethingWork() {
	log.Debug("extend sql.Db")
}
