package config

type DBInfo struct {
	ip     string
	port   string
	dbName string
	user   string
	pass   string
}

type MemcacheInfo struct {
	ip   string
	port int
}

const (
	DB_DEV = 1 + iota
)

var (
	db_list = []int{
		DB_DEV,
	}
	slave_db_info_list = map[int]*DBInfo{
		DB_DEV: {"127.0.0.1", "3306", "db_oceansf", "admin", "passwdforadmin"},
	}

	master_db_info_list = map[int]*DBInfo{
		DB_DEV: {"127.0.0.1", "3306", "db_oceansf", "admin", "passwdforadmin"},
	}
)

func GetDSN(dbType int, isSlave bool) string {

	var dbInfo *DBInfo
	if isSlave {
		dbInfo = slave_db_info_list[dbType]
	} else {
		dbInfo = master_db_info_list[dbType]
	}

	dsn := dbInfo.user + ":" + dbInfo.pass + "@tcp(" + dbInfo.ip + ":" + dbInfo.port + ")/" + dbInfo.dbName + "?charset=utf8"
	return dsn
}
