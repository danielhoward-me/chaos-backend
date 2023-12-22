package admins

import (
	"github.com/danielhoward-me/chaos-backend/sso"

	"database/sql"
	"fmt"
)

// func GetAll(db *sql.DB) ([]sso.Account, error)

func IsAdmin(db *sql.DB, id string) (exists bool, err error) {
	if err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM admins WHERE id = $1);", id).Scan(&exists); err != nil {
		err = fmt.Errorf("failed to test if user is admin: %s", err)
		return
	}

	return
}

func IsLocalAdmin(db *sql.DB, account sso.Account) (bool, error) {
	if account.Admin {
		return true, nil
	}

	return IsAdmin(db, account.UserId)
}
