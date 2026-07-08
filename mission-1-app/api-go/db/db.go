package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

// getEnv lit les variables d'env avec une valeur par défaut
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// NewDB tente de se connecter à la base de données avec une stratégie de retry
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

	// On tente de se connecter jusqu'à 10 fois
	for i := 0; i < 10; i++ {
		conn, err = sql.Open("mysql", dsn)
		
		// Si l'ouverture réussit, on vérifie la connectivité réelle avec Ping
		if err == nil {
			if err = conn.Ping(); err == nil {
				fmt.Printf("Connecté à la base de données via : %s\n", host)
				return conn
			}
		}

		fmt.Printf("Tentative %d/10 : En attente de la base de données, nouvelle tentative dans 2s...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	// Si après 10 fois ça ne fonctionne toujours pas, on arrête tout
	panic(fmt.Sprintf("Impossible de se connecter à la base après 10 tentatives : %v", err))
}