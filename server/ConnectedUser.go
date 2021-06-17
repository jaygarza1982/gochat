package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type ConnectedUser struct {
	Id   string
	Conn *websocket.Conn
}

func (u *ConnectedUser) Read(conn *websocket.Conn, socketsSet *map[string]ConnectedUser) {
	conn.SetCloseHandler(func(code int, text string) error {
		textString := fmt.Sprintf("Closing client code %d because %s", code, text)

		delete(*socketsSet, u.Id)

		return errors.New(textString)
	})

	for {
		// Read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Print message
		log.Println(string(p))

		// When we receive a message, we will write it to another user

		// Construct a user message that was sent to us
		newUserMessage := UserMessage{}
		newUserMessage.FromJSON(string(p))

		//Get the user that we are receiving a message from
		if user, ok := (*socketsSet)[newUserMessage.ReceiverId]; ok {
			user.Conn.WriteJSON(newUserMessage)
		}
	}
}
