package main

import (
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

var count int = 0

//We will use this as a set
var socketsSet = map[string]ConnectedUser{}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)

	// TODO: assign a username or userId when user logs in
	// TODO:
	// session, _ := store.Get(r, "auth")
	// Check if username is in our session
	// if val, ok := session.Values["username"]; ok {

	// }

	count++
	clientId := fmt.Sprintf("%d", count)

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
	log.Println("TEST ENDPOINT HIT!")
}

func login(w http.ResponseWriter, r *http.Request) {
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	session, _ := store.Get(r, "auth")

	// TODO: Check these from a database
	username := r.FormValue("username")
	// password := r.FormValue("password")

	// Set user as authenticated
	session.Values["username"] = username
	session.Save(r, w)
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/api/test", testAPI)
	http.HandleFunc("/api/login", login)
}

func main() {
	// Database setup
	dsn := os.Getenv("DB_CONN_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Could not open DB! %v", err.Error())
	}

	db.AutoMigrate(&User{})

	// Start HTTP server
	fmt.Println("Server started...")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
