package main

import (
	"database/sql"
	"os"

	"github.com/danielhoward-me/chaos-backend/database"
)

var db *sql.DB

func main() {
	db = database.GetDb()
	os.MkdirAll(os.Getenv("SCREENSHOTPATH"), os.ModePerm)
	createServer()
}
