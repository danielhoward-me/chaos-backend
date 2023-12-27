package admins

import (
	"github.com/danielhoward-me/chaos-backend/sso"

	"database/sql"
	"fmt"
)

func GetAll(db *sql.DB) ([]sso.Account, error) {
	rows, err := db.Query("SELECT id FROM admins")
	if err != nil {
		return nil, fmt.Errorf("there was an error when fetching the admins: %s", err)
	}
	defer rows.Close()

	admins := []sso.Account{}
	for rows.Next() {
		admin := sso.Account{}
		rows.Scan(&admin.UserId)
		admins = append(admins, admin)
	}

	return admins, nil
}

func Remove(db *sql.DB, id string) error {
	res, err := db.Exec("DELETE FROM admins WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("there was an error removing an admin from the database: %s", err)
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("there was an error getting the number of admins removed from the database: %s", err)
	}

	if rowsCount != 1 {
		return fmt.Errorf("no rows were affected when removing an admin from the database: %s", err)
	}

	return nil
}

func New(db *sql.DB, id string) (err error) {
	_, err = db.Exec("INSERT INTO admins(id) VALUES ($1)", id)
	return
}

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
