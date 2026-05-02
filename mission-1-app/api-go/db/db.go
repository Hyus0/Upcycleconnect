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

	if err := ensureSchemaCompatibility(conn); err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to database!")
	return conn
}

func ensureSchemaCompatibility(conn *sql.DB) error {
	if err := ensureColumn(conn, "UTILISATEUR", "token", "ALTER TABLE UTILISATEUR ADD COLUMN token VARCHAR(255) NULL"); err != nil {
		return err
	}
	return nil
}

func ensureColumn(conn *sql.DB, tableName, columnName, alterSQL string) error {
	var count int
	query := `
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = ?
	`

	if err := conn.QueryRow(query, tableName, columnName).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err := conn.Exec(alterSQL)
	return err
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
