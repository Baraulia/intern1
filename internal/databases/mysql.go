package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func NewMysqlDB(database *MysqlDB) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		database.Username, database.Password, database.Host, database.Port, database.DBName))
	if err != nil {
		log.Panicf("Database open error:%s", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("DB ping error:%s", err)
		return nil, err
	}
	return db, nil
}
