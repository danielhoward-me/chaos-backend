package sso

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var ssoDevPort int

type Account struct {
	UserId         string `json:"userId"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	Admin          bool   `json:"admin"`
}

func Init(givenSsoDevPort int) {
	ssoDevPort = givenSsoDevPort
}

func Get(authorisation string) (account Account, exists bool, err error) {
	ssoOrigin := "https://sso.danielhoward.me"
	if ssoDevPort != 0 {
		ssoOrigin = fmt.Sprintf("http://local.danielhoward.me:%d", ssoDevPort)
	}
	ssoPath := fmt.Sprintf("%s/api/oauth2/account", ssoOrigin)

	req, err := http.NewRequest("GET", ssoPath, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %s", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authorisation))

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to do request: %s", err)
		return
	}

	if res.StatusCode != 200 {
		exists = false
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body: %s", err)
		return
	}

	if err = json.Unmarshal(data, &account); err != nil {
		err = fmt.Errorf("failed to read body as json: %s", err)
		return
	}

	exists = true
	return
}

func GetUser(id string) (account Account, exists bool, err error)
