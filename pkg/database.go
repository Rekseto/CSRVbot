package pkg

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"
	"log"
)

type Database interface {
	InitDatabase() error
	GetDbMap() *gorp.DbMap
}

var dbMap *gorp.DbMap

func InitDatabase(url string) error {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return err
	}

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"}}
	return nil
}

func GetDbMap() *gorp.DbMap {
	return dbMap
}

func CreateTablesIfNotExists() {
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Could not create tables", err)
	}
}
