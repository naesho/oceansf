package db

import (
	"database/sql"
	"github.com/ohsaean/gogpd/lib"
	"github.com/ohsaean/oceansf/config"
)

func GetInstance(dbType int, isSlave bool) *sql.DB {
	instance, err := sql.Open("mysql", config.GetDSN(dbType, isSlave))
	lib.CheckError(err)
	return instance
}
