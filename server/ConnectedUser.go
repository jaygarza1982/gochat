package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type ConnectedUser struct {
	Id string
}

func (u *ConnectedUser) Read(conn *websocket.Conn, socketsSet *map[*websocket.Conn]ConnectedUser) {
	conn.SetCloseHandler(func(code int, text string) error {
		textString := fmt.Sprintf("Closing client code %d because %s", code, text)

		delete(*socketsSet, conn)

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
		// If it was sent to us, we will not do anything

		// Construct a user message that was sent to us
		newUserMessage := UserMessage{}
		newUserMessage.FromJSON(string(p))

		// TODO: Refactor socket map to be string->socket
		for socket := range *socketsSet {
			// TODO: (Maybe) Don't write to self
			// Might be nice to have inherent confirmation messages sent
			user := (*socketsSet)[socket]
			fmt.Println("User ID", user.Id)

			if user.Id == newUserMessage.ReceiverId {

				socket.WriteJSON(newUserMessage)
			}
		}
	}
}
