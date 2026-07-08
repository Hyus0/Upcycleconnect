package db

import (
	// "context"
	"database/sql"
	"fmt"
	// "os"
	// "time"

	_ "github.com/go-sql-driver/mysql"
)

// Conn est le pool de connexions partage par l'ensemble de l'API.
// nil tant que NewDB n'a pas reussi a joindre la base.
var Conn *sql.DB

// const driver = "mysql"

// Valeurs par defaut. En production (Docker Compose) elles sont surchargees
// par les variables d'environnement DB_HOST / DB_PORT / DB_USER / DB_PASSWORD / DB_NAME.
const (
	driver   = "mysql"
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = ""
	dbname   = "upcycletest"
)

func NewDB() *sql.DB {
	var sqlInfo = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		user, password, host, port, dbname)
	conn, err := sql.Open(driver, sqlInfo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to database !")
	return conn
}
// getenv retourne la variable d'environnement key ou fallback si vide.
// func getenv(key, fallback string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}
// 	return fallback
// }

// NewDB ouvre le pool MySQL a partir des variables d'environnement, verifie la
// connexion avec un Ping et configure le pool. En cas d'echec, on journalise et
// on retourne nil : les handlers peuvent ainsi detecter une base indisponible
// (Conn == nil) plutot que d'echouer requete par requete.
// func NewDB() *sql.DB {
// 	host := getenv("DB_HOST", defaultHost)
// 	port := getenv("DB_PORT", defaultPort)
// 	user := getenv("DB_USER", defaultUser)
// 	password := getenv("DB_PASSWORD", defaultPassword)
// 	name := getenv("DB_NAME", defaultName)

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
// 		user, password, host, port, name)

// 	conn, err := sql.Open(driver, dsn)
// 	if err != nil {
// 		fmt.Printf("[db] DSN invalide (%s@%s:%s/%s): %v\n", user, host, port, name, err)
// 		return nil
// 	}

// 	// Parametrage du pool de connexions.
// 	conn.SetMaxOpenConns(25)
// 	conn.SetMaxIdleConns(10)
// 	conn.SetConnMaxLifetime(5 * time.Minute)
// 	conn.SetConnMaxIdleTime(2 * time.Minute)

// 	// Verification reelle de la connexion (sql.Open ne se connecte pas).
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	if err := conn.PingContext(ctx); err != nil {
// 		fmt.Printf("[db] base injoignable (%s@%s:%s/%s): %v\n", user, host, port, name, err)
// 		_ = conn.Close()
// 		return nil
// 	}

// 	fmt.Printf("[db] connecte a MySQL %s@%s:%s/%s\n", user, host, port, name)
// 	return conn
// }
