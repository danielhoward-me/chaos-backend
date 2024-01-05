package saves

import (
	screenshotUtils "github.com/danielhoward-me/chaos-backend/screenshot/utils"

	"database/sql"
	"fmt"
	"strconv"
)

type Save struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Data       string `json:"data"`
	Screenshot string `json:"screenshot,omitempty"`
	UserId     string `json:"userId"`
}

type RequestSave struct {
	Name     string `json:"name"`
	Data     string `json:"data"`
	IsPreset bool   `json:"isPreset"`
}

func GetUsers(db *sql.DB, userId string) ([]Save, error) {
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
		save := Save{UserId: userId}
		rows.Scan(&save.Id, &save.Name, &save.Data, &save.Screenshot)
		saves = append(saves, save)
	}

	return saves, nil
}

func Exists(db *sql.DB, id int) (bool, error) {
	var exists bool = true
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM saves WHERE id = $1)", id).Scan(&exists); err != nil {
		return false, fmt.Errorf("there was an error when checking if a save exists: %s", err)
	}

	return exists, nil
}

func Get(db *sql.DB, id int) (Save, error) {
	rows, err := db.Query("SELECT name, data, screenshot, user_id FROM saves WHERE id = $1", id)
	if err != nil {
		return Save{}, fmt.Errorf("there was an error when getting a save from the database: %s", err)
	}
	defer rows.Close()

	save := Save{Id: strconv.FormatInt(int64(id), 10)}
	rows.Scan(&save.Name, &save.Data, &save.Screenshot, &save.UserId)

	return save, nil
}

func Delete(db *sql.DB, id int, userId string, isAdmin bool) (completed bool, err error) {
	userQuery := "user_id = $1"
	if isAdmin {
		userQuery = fmt.Sprintf("(%s OR user_id IS NULL)", userQuery)
	}

	res, err := db.Exec(fmt.Sprintf("DELETE FROM saves WHERE %s AND id = $2", userQuery), userId, id)
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
	var id string

	var userIdArg any = userId
	if userId == "0" {
		userIdArg = nil
	}

	row := db.QueryRow("INSERT INTO saves (name, data, user_id) VALUES ($1, $2, $3) RETURNING id", name, data, userIdArg)
	if err = row.Scan(&id); err != nil {
		err = fmt.Errorf("there was an error when creating a save: %s", err)
		return
	}

	save = Save{
		Id:         id,
		Name:       name,
		Data:       data,
		Screenshot: screenshotUtils.Hash(data),
	}
	return
}

func ChangeName(db *sql.DB, id string, name string) error {
	_, err := db.Exec("UPDATE saves SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return fmt.Errorf("there was an error when changing save name: %s", err)
	}
	return nil
}
