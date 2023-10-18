package admins

import (
	"github.com/danielhoward-me/chaos-backend/sso"

	"database/sql"
)

func GetAll(db *sql.DB) ([]sso.Account, error)

func IsAdmin(db *sql.DB, id string) (bool, error) {
	row := map[string]bool{}
	if err := db.QueryRow("SELECT EXISTS(id) AS exists FROM admins WHERE id = $1", id).Scan(&row); err != nil {

	}
}
