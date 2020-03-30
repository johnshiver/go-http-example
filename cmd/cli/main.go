package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

const DevUrl = "http://localhost:8080"

type User struct {
	ID       int64
	username string
	password string
	token    string
}

var devUsers = map[string]User{
	"yoda": {
		username: "yoda",
		password: "im-testing-here",
	},
	"luke": {
		username: "luke",
		password: "im-testing-here",
	},
}

func main() {
	rand.Seed(time.Now().Unix())
	var (
		rec, cmd, user string
		limit          int
		start          int64
	)
	flag.StringVar(&user, "user", "yoda", "user to perform action")
	flag.StringVar(&rec, "rec", "luke", "user to receive message")
	flag.StringVar(&cmd, "cmd", "createUsers", "cmd to run")
	flag.IntVar(&limit, "limit", 100, "message limit")
	flag.Int64Var(&start, "start", 0, "message offset")

	flag.Parse()

	switch cmd {
	case "createUsers":
		for _, u := range devUsers {
			createUser(u)
		}
	case "sendMsg":
		sendUser := loginUser(devUsers[user])
		recUser := loginUser(devUsers[rec])
		sendMessage(sendUser, recUser, "image")
		sendMessage(sendUser, recUser, "text")
		sendMessage(sendUser, recUser, "video")
	case "getMsg":
		user := loginUser(devUsers[user])
		recUser := loginUser(devUsers[rec])
		getMessages(user, recUser, start, limit)

	default:
		fmt.Printf("Not a supported action %v\n", cmd)
	}

}
