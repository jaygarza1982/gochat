package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var count int = 0

//We will use this as a set
var socketsSet = map[*websocket.Conn]ConnectedUser{}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)

	// TODO: assign a username or userId when user logs in
	count++
	clientId := fmt.Sprintf("%d", count)

	// Add our new user to the sockets map
	newUser := ConnectedUser{clientId}
	socketsSet[ws] = newUser

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Client %s connected\n", clientId)

	// Listen on this socket forever until it closes
	newUser.Read(ws, &socketsSet)
}

func testAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
	log.Println("TEST ENDPOINT HIT!")
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/api/test", testAPI)
}

func main() {
	fmt.Println("Server started...")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
