package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// DSN returns the Data Source Name
func postgresqldsn(dbENV string) string {
	// Example: root:@tcp(localhost:3306)/ipaddressservices

	return os.Getenv(dbENV)

}

// Connect to the database -- for all of them
func Connect(dbENV string) {
	var err error

	// Connect to PostgreSQL
	if DB, err = sql.Open("postgres", postgresqldsn(dbENV)); err != nil {
		log.Println("PostGres SQL Driver Error", err)
	}

	// Check if is alive
	if err = DB.Ping(); err != nil {
		log.Println("Database Error", err)
	}
}
