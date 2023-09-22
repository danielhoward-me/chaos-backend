package main

import (
	"database/sql"

	"github.com/danielhoward-me/chaos-backend/database"
)

var db *sql.DB

func main() {
	db = database.GetDb()
	createServer()
}
