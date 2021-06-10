package main

import (
	"encoding/json"
	"fmt"
)

type UserMessage struct {
	SenderId    string
	ReceiverId  string
	MessageText string
}

func (u *UserMessage) ToJSON() string {
	bytes, err := json.Marshal(u)

	if err != nil {
		fmt.Println("Error converting user message to JSON", err)
	}

	return string(bytes)
}

func (u *UserMessage) FromJSON(jsonString string) {
	err := json.Unmarshal([]byte(jsonString), u)

	if err != nil {
		fmt.Println("Could not parse UserMessage JSON", err)
	}
}
