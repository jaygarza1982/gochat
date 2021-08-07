package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//We will use this as a set
var socketsSet = map[string]ConnectedUser{}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)

	session, _ := store.Get(r, "auth")
	clientId := ""

	// Check if username is in our session
	if val, ok := session.Values["username"]; ok {
		clientId = fmt.Sprintf("%v", val)
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("401 - User session was not found. Please login first."))
	}

	// Add our new user to the sockets map
	newUser := ConnectedUser{clientId, ws}
	socketsSet[clientId] = newUser

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Client %s connected\n", clientId)

	// Listen on this socket forever until it closes
	newUser.Read(ws, &socketsSet)
}

func testAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func login(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var request LoginRequest

		data := json.NewDecoder(r.Body)
		data.Decode(&request)

		user := User{Username: request.Username}
		if valid := user.CheckPassword(db, request.Password); valid {
			// Set user as authenticated
			session, _ := store.Get(r, "auth")
			session.Values["username"] = request.Username
			session.Save(r, w)

			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("401 - User credentials are incorrect"))
		}
	}
}

func register(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type RegisterRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var request RegisterRequest

		data := json.NewDecoder(r.Body)
		data.Decode(&request)

		user := User{Username: request.Username}
		if err := user.Register(db, request.Password); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Error in registration."))
		}
	}
}

func ListConversations(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth")
		username := ""

		// Check if username is in our session
		if val, ok := session.Values["username"]; ok {
			username = fmt.Sprintf("%v", val)
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("401 - User session was not found. Please login first."))
		}

		user := User{Username: username}
		conversations := user.GetConversations(db)

		if bytes, err := json.Marshal(conversations); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Error in registration."))
		} else {
			w.Write(bytes)
		}
	}
}

func ListMessages(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth")
		username := ""

		// Check if username is in our session
		if val, ok := session.Values["username"]; ok {
			username = fmt.Sprintf("%v", val)
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("401 - User session was not found. Please login first."))
		}

		// Get messages request
		type MessageRequest struct {
			Username string `json:"username"`
		}

		var request MessageRequest

		data := json.NewDecoder(r.Body)
		data.Decode(&request)

		// Read current users messages from the requested user
		user := User{Username: username}
		messages := user.ReadMessages(db, request.Username)

		// Send messages data
		if bytes, err := json.Marshal(messages); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Error in registration."))
		} else {
			w.Write(bytes)
		}
	}
}

func SendMessage(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth")
		username := ""

		// Check if username is in our session
		if val, ok := session.Values["username"]; ok {
			fmt.Printf("Username %v is sending a message\n", val)
			username = fmt.Sprintf("%v", val)
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("401 - User session was not found. Please login first."))
		}

		var message UserMessage

		data := json.NewDecoder(r.Body)
		data.Decode(&message)

		user := User{Username: username}
		user.SendMessage(db, &message, func(sentMessage *UserMessage) {
			// Ensure the socket exists
			if _, ok := socketsSet[message.ReceiverId]; !ok {
				fmt.Printf("Receiver %v not found within sockets\n", message.ReceiverId)

				return
			}

			fmt.Printf("Sending message to %v\n", message.ReceiverId)

			// Send our message ID over the socket
			// The receiver will later query this to find the message text
			// This is because I have not yet found a way to proxy WSS
			socketsSet[message.ReceiverId].Conn.WriteJSON(sentMessage.ID)
		})
	}
}

func setupRoutes(db *gorm.DB) {
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/api/test", testAPI)
	http.HandleFunc("/api/login", login(db))
	http.HandleFunc("/api/register", register(db))
	http.HandleFunc("/api/send-message", SendMessage(db))
	// TODO: Will list messages from a specified user and list messages to a specified user where username is user logged in
	// Example: User A -> message0. User B -> message1. If user A is logged in, it will return message0 and message1
	// Since this is all part of the conversation between the two users
	http.HandleFunc("/api/list-messages", ListMessages(db))
	http.HandleFunc("/api/conversations", ListConversations(db))
}

func main() {
	// Database setup
	dsn := os.Getenv("DB_CONN_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Could not open DB! %v", err.Error())
	}

	db.AutoMigrate(&User{}, &UserMessage{})

	// Start HTTP server
	fmt.Println("Server started...")
	setupRoutes(db)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
