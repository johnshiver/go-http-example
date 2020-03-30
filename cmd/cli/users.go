package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResp struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

func loginUser(u User) User {
	client := resty.New()
	var lRes loginResp

	r, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(&loginPayload{Username: u.username, Password: u.password}).
		Post(DevUrl + "/login")

	if err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(r.Body())
	err = json.NewDecoder(buff).Decode(&lRes)
	if err != nil {
		panic(err)
	}
	u.token = lRes.Token
	u.ID = lRes.ID
	return u
}

type createUserResp struct {
	ID int64
}

func createUser(u User) User {
	var resp createUserResp
	client := resty.New()

	r, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(&loginPayload{Username: u.username, Password: u.password}).
		SetResult(&resp).
		Post(DevUrl + "/users")

	if err != nil {
		panic(err)
	}
	fmt.Println(r.String())
	u.ID = resp.ID
	return u
}
