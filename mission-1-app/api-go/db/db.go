package db

import (
	"database/sql"
	"fmt"
	 _ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

const (
	driver   = "mysql"
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = ""
	dbname   = "partielb_go"
)

func NewDB() *sql.DB {
		var sqlInfo = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		 			user, password,host, port, dbname)
	conn, err := sql.Open(driver, sqlInfo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to database !")
	return conn
}
