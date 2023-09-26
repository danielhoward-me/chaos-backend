package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/danielhoward-me/chaos-backend/database"
	"github.com/danielhoward-me/chaos-backend/screenshot/worker"
	"github.com/danielhoward-me/chaos-backend/sso"
)

var db *sql.DB
var ChaosDevPort int
var SsoDevPort int

func main() {
	flag.IntVar(&ChaosDevPort, "chaos-dev-port", 0, "The custom port to connect to chaos with for screenshots")
	flag.IntVar(&SsoDevPort, "sso-dev-port", 0, "The custom port to connect to sso")
	flag.Parse()
	worker.Init(ChaosDevPort)
	sso.Init(SsoDevPort)

	db = database.GetDb()
	os.MkdirAll(os.Getenv("SCREENSHOTPATH"), os.ModePerm)
	createServer()
}
