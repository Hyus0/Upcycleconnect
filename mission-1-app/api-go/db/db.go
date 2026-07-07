package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

const driver = "mysql"

const (
	defaultHost     = "127.0.0.1"
	defaultPort     = "3306"
	defaultUser     = "root"
	defaultPassword = "password"
	defaultName     = "upcycletest"
)

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func NewDB() *sql.DB {
	host := getenv("DB_HOST", defaultHost)
	port := getenv("DB_PORT", defaultPort)
	user := getenv("DB_USER", defaultUser)
	password := getenv("DB_PASSWORD", defaultPassword)
	name := getenv("DB_NAME", defaultName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, password, host, port, name)

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		fmt.Printf("[db] invalid DSN (%s@%s:%s/%s): %v\n", user, host, port, name, err)
		return nil
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(5 * time.Minute)
	conn.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		fmt.Printf("[db] unavailable database (%s@%s:%s/%s): %v\n", user, host, port, name, err)
		_ = conn.Close()
		return nil
	}

	fmt.Printf("[db] connected to MySQL %s@%s:%s/%s\n", user, host, port, name)
	return conn
}
