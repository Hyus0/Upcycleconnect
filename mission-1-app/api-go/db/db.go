package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func NewDB() *sql.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "upcycletest")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname)

	var conn *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		conn, err = sql.Open("mysql", dsn)
		
		if err == nil {
			if err = conn.Ping(); err == nil {
				fmt.Printf("Connecté à la base de données via : %s\n", host)
				return conn
			}
		}

		fmt.Printf("Tentative %d/10 : En attente de la base de données, nouvelle tentative dans 2s...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	panic(fmt.Sprintf("Impossible de se connecter à la base après 10 tentatives : %v", err))
}