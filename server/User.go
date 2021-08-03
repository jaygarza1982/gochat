package main

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           int `gorm:"primaryKey" faker:"-"`
	Username     string
	PasswordHash string
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (u *User) CheckPassword(db *gorm.DB, password string) bool {
	// Obtain user from database
	existing := User{}
	db.First(&existing, "username = ?", u.Username)

	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(password))

	if err != nil {
		fmt.Printf("password is incorrect %v", err.Error())
	}

	return err == nil
}

func (u *User) Register(db *gorm.DB, password string) error {
	if !u.CanCreate(db.First(&User{}, "Username = ?", u.Username).RowsAffected) {
		return errors.New("could not create user username already exists")
	}

	hash, err := hashPassword(password)

	if err != nil {
		fmt.Printf("could not hash password %v", err.Error())
		return err
	}

	u.PasswordHash = hash

	create(db, u)

	return nil
}

// Sends a message to another user
// The message contains who it is to
func (u *User) SendMessage(db *gorm.DB, message *UserMessage, callback func()) {
	// Ensure that the message is from us
	message.SenderId = u.Username

	db.Create(message)

	// Run our optional callback, this could be sending message over a websocket if desired
	if callback != nil {
		callback()
	}
}

func (u *User) ReadMessages(db *gorm.DB, senderId string) *[]UserMessage {
	userMessages := []UserMessage{}

	// Read messages addressed to current user
	// Read messages addressed to other from the current user
	db.Where("(receiver_id = ? AND sender_id = ?) OR (sender_id = ? AND receiver_id = ?)", u.Username, senderId, u.Username, senderId).Find(&userMessages)

	return &userMessages
}

// Return all messages that have receiver as the current user
func (u *User) GetConversations(db *gorm.DB) []string {
	var usernames []string

	// Get distinct chats where receiver is ours
	rows, rowError := db.Raw("SELECT DISTINCT sender_id FROM user_messages WHERE receiver_id = ?", u.Username).Rows()

	if rowError != nil {
		fmt.Printf("Error %v", rowError)
	}

	// For all rows returned, append data to usernames slice
	for rows.Next() {
		var username string
		rows.Scan(&username)

		usernames = append(usernames, username)
	}

	return usernames
}

func (u *User) CanCreate(rows int64) bool {
	return rows == 0
}

func create(db *gorm.DB, u *User) {
	db.Create(u)
}

func (u *User) Delete(db *gorm.DB) {
	fmt.Printf("Deleting %v", u.ID)
	db.Delete(&u, u.ID)
}
