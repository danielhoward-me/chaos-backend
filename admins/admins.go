package admins

import (
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
