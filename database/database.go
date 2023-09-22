package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB
var initialised = false

var PGUSER = os.Getenv("PGUSER")
var PGPASSWORD = os.Getenv("PGPASSWORD")
var PGHOST = os.Getenv("PGHOST")
var PGPORT = os.Getenv("PGPORT")
var PGDATABASE = os.Getenv("PGDATABASE")
var PSSSLMODE = os.Getenv("PSSSLMODE")

var CONNECTION_STRING = fmt.Sprintf(
	"postgres://%s:%s@%s:%s/%s?sslmode=%s",
	PGUSER,
	PGPASSWORD,
	PGHOST,
	PGPORT,
	PGDATABASE,
	PSSSLMODE,
)

func GetDb() *sql.DB {
	if !initialised {
		connect()
	}

	initialised = true
	return db
}

func connect() {
	if initialised {
		return
	}

	fmt.Printf("Connecting to database %s\n", PGDATABASE)

	dbConnection, err := sql.Open("postgres", CONNECTION_STRING)
	if err != nil {
		panic(fmt.Errorf("failed to create database connection: %s", err))
	}

	if err = dbConnection.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %s", err))
	}

	fmt.Printf("Connected to database\n")

	db = dbConnection
}
