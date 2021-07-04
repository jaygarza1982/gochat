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
		fmt.Printf("%v", val)
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

func login(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var request LoginRequest

	data := json.NewDecoder(r.Body)
	data.Decode(&request)

	session, _ := store.Get(r, "auth")

	// TODO: Check username and password from DB

	// Set user as authenticated
	session.Values["username"] = request.Username
	session.Save(r, w)
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
		user.Register(db, request.Password)
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
		user.SendMessage(db, &message, func() {
			// Send our message over the socket
			socketsSet[message.ReceiverId].Conn.WriteJSON(message)
		})
	}
}

func setupRoutes(db *gorm.DB) {
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/api/test", testAPI)
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/register", register(db))
	http.HandleFunc("/api/send-message", SendMessage(db))
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
