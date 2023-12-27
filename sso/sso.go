package sso

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func makeSsoRequest(path string, authorisation string, method string, body string) (res *http.Response, err error) {
	ssoOrigin := "https://sso.danielhoward.me"
	if ssoDevPort != 0 {
		ssoOrigin = fmt.Sprintf("http://local.danielhoward.me:%d", ssoDevPort)
	}
	ssoPath := fmt.Sprintf("%s/%s", ssoOrigin, path)

	req, err := http.NewRequest(method, ssoPath, strings.NewReader(body))
	if err != nil {
		err = fmt.Errorf("failed to create request: %s", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authorisation))

	client := http.DefaultClient
	res, err = client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to do request: %s", err)
		return
	}

	return
}

func readBody[T any](res *http.Response) (body T, err error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body: %s", err)
		return
	}

	if err = json.Unmarshal(data, &body); err != nil {
		err = fmt.Errorf("failed to read body as json: %s", err)
		return
	}

	return
}

func Get(authorisation string) (account Account, exists bool, err error) {
	res, err := makeSsoRequest("/api/oauth2/account", authorisation, "GET", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode != 200 {
		exists = false
		return
	}

	account, err = readBody[Account](res)
	if err != nil {
		fmt.Println(err)
		return
	}

	exists = true
	return
}

type GetUserResponse struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profilePicture"`
	Successful     bool   `json:"successful"`
}

func GetUser(id string, username string, authorisation string) (account GetUserResponse, exists bool, err error) {
	if id == "" && username == "" {
		err = fmt.Errorf("id or username must be given")
		return
	}

	reqBody := map[string]string{
		"id":       id,
		"username": username,
	}
	reqBodyString, err := json.Marshal(reqBody)
	if err != nil {
		return
	}

	res, err := makeSsoRequest("/api/get-user", authorisation, "POST", string(reqBodyString))
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := readBody[GetUserResponse](res)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !body.Successful {
		exists = false
		return
	}

	exists = true
	account = body

	return
}
