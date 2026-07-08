package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

// Fonction utilitaire pour lire les variables d'env avec une valeur par défaut
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func NewDB() *sql.DB {
	// Lecture des configurations depuis les variables d'environnement
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "upcycletest")

	// Construction du DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	// Vérification de la connexion
	if err := conn.Ping(); err != nil {
		panic(fmt.Sprintf("Impossible de se connecter à la base : %v", err))
	}

	var conn *sql.DB
		var err error
		for i := 0; i < 10; i++ {
			conn, err = sql.Open("mysql", dsn)
			if err == nil && conn.Ping() == nil {
				fmt.Println("Connecté à la base de données !")
				return conn
			}
			fmt.Println("En attente de la base de données, nouvelle tentative dans 2s...")
			time.Sleep(2 * time.Second)
		}
		panic("Impossible de se connecter à la base après 10 tentatives")
	fmt.Println("Connecté à la base de données via :", host)
	return conn
}