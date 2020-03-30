package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var (
	textMessage = map[string]interface{}{
		"type": "text",
		"text": "hey there",
	}
	imageMessage = map[string]interface{}{
		"type":   "image",
		"url":    "https://www.good-images.com/",
		"height": 1,
		"width":  5,
	}
	videoMessage = map[string]interface{}{
		"type":   "video",
		"url":    "https://www.good-videos.com/",
		"source": "youtube",
	}
)

type msgPayload struct {
	Sender    int64                  `json:"sender"`
	Recipient int64                  `json:"recipient"`
	Content   map[string]interface{} `json:"content"`
}

func sendMessage(sender, receiver User, messageType string) {
	var (
		resp    loginResp
		payload msgPayload
		content map[string]interface{}
	)
	client := resty.New()

	switch messageType {
	case "text":
		content = textMessage
	case "image":
		content = imageMessage
	case "video":
		content = videoMessage
	}

	payload = msgPayload{
		Sender:    sender.ID,
		Recipient: receiver.ID,
		Content:   content,
	}

	reqJSON, _ := json.Marshal(&payload)
	rJSON := bytes.NewBuffer(reqJSON)

	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+sender.token).
		SetBody(rJSON).
		SetResult(&resp).
		Post(DevUrl + "/messages")

	if err != nil {
		panic(err)
	}

	fmt.Println(response.String())
}

func getMessages(u, rec User, start int64, limit int) {
	client := resty.New()
	qParams := map[string]string{
		"recipient": fmt.Sprintf("%d", rec.ID),
		"start":     fmt.Sprintf("%d", start),
		"limit":     fmt.Sprintf("%d", limit),
	}

	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+u.token).
		SetQueryParams(qParams).
		Get(DevUrl + "/messages")
	if err != nil {
		panic(err)
	}
	fmt.Println(response.String())
}
