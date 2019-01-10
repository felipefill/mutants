package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Postgres driver for database/sql
)

var _db *sql.DB

// GetDB gets DB connection, in case of failure it will panic
func GetDB() *sql.DB {
	if _db != nil {
		return _db
	}

	host, name, user, pswd := getDatabaseInfo()
	_db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", user, pswd, host, name))
	if err != nil {
		panic(fmt.Sprintf("Could not connect to database: %s", err.Error()))
	}

	return _db
}

// InjectDatabase uses given database
func InjectDatabase(database *sql.DB) {
	_db = database
}

// MustGetEnvVar tries to retrieve given environment variable, in case of failure it will panic
func MustGetEnvVar(v string) string {
	envVar := os.Getenv(v)

	if envVar == "" {
		panic(fmt.Sprintf("Failed to retrieve %s environment variable", v))
	}

	return envVar
}

func getDatabaseInfo() (host string, name string, user string, pswd string) {
	return MustGetEnvVar("DB_HOST"),
		MustGetEnvVar("DB_NAME"),
		MustGetEnvVar("DB_USER"),
		MustGetEnvVar("DB_PSWD")
}
