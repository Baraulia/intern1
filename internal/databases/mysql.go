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
	query := `CREATE TABLE IF NOT EXISTS countries
(
    id serial PRIMARY KEY,
    name varchar(100) not null UNIQUE,
	full_name varchar(150) not null,
	english_name varchar(150) not null,
	alpha_2 varchar(2) not null,
	alpha_3 varchar(3) not null UNIQUE,
	iso int not null UNIQUE,
	location varchar(150) not null,
	location_precise varchar(150) not null,
	url varchar(255) not null
);`
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Start migrations error:%s", err)
		return nil, err
	}
	return db, nil
}
