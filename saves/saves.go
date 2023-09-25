package saves

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
)

type Save struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Data       string `json:"data"`
	Screenshot string `json:"screenshot,omitempty"`
}

type RequestSave struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func Get(db *sql.DB, userId string) ([]Save, error) {
	queryString := "IS NULL"
	queryParams := []any{}
	if userId != "" {
		queryString = "= $1"
		queryParams = append(queryParams, userId)
	}

	rows, err := db.Query(fmt.Sprintf("SELECT id, name, data, screenshot FROM saves WHERE user_id %s", queryString), queryParams...)
	if err != nil {
		return nil, fmt.Errorf("there was an error when getting saves from the database: %s", err)
	}
	defer rows.Close()

	saves := []Save{}
	for rows.Next() {
		save := Save{}
		rows.Scan(&save.Id, &save.Name, &save.Data, &save.Screenshot)
		saves = append(saves, save)
	}

	return saves, nil
}

func Delete(db *sql.DB, id int, userId string) (completed bool, err error) {
	res, err := db.Exec("DELETE FROM saves WHERE user_id = $1 AND id = $2", userId, id)
	if err != nil {
		err = fmt.Errorf("there was an error deleting getting saves from the database: %s", err)
		return
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		err = fmt.Errorf("there was an error getting the number of saves deleted from the database: %s", err)
		return
	}

	completed = rowsCount != 0
	return
}

func Create(db *sql.DB, name string, data string, userId string) (save Save, err error) {
	res, err := db.Exec("INSERT INTO saves (name, data, user_id) VALUES ($1, $2, $3)", name, data, userId)
	if err != nil {
		err = fmt.Errorf("there was an error when creating a save: %s", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("there was an error when getting a new save's ID: %s", err)
		return
	}

	save = Save{
		Id:         strconv.FormatInt(id, 10),
		Name:       name,
		Data:       data,
		Screenshot: getScreenshotHash(data),
	}
	return
}

func getScreenshotHash(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
