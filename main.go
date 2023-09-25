package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/danielhoward-me/chaos-backend/database"
	"github.com/danielhoward-me/chaos-backend/screenshot/worker"
)

var db *sql.DB
var ChaosDevPort int

func main() {
	flag.IntVar(&ChaosDevPort, "chaos-dev-port", 0, "The custom port to connect to chaos server with for screenshots")
	flag.Parse()
	worker.Init(ChaosDevPort)

	db = database.GetDb()
	os.MkdirAll(os.Getenv("SCREENSHOTPATH"), os.ModePerm)
	createServer()
}
