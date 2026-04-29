package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

const (
	driver   = "mysql"
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = ""
	dbname   = "upcycletest"
)

func NewDB() *sql.DB {
	dbHost := getenv("DB_HOST", host)
	dbPort := getenv("DB_PORT", fmt.Sprintf("%d", port))
	dbUser := getenv("DB_USER", user)
	dbPassword := getenv("DB_PASSWORD", password)
	dbName := getenv("DB_NAME", dbname)

	var sqlInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	conn, err := sql.Open(driver, sqlInfo)
	if err != nil {
		panic(err.Error())
	}

	conn.SetConnMaxLifetime(3 * time.Minute)
	conn.SetMaxIdleConns(10)
	conn.SetMaxOpenConns(25)

	if err := conn.Ping(); err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to database!")
	return conn
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
