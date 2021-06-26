package main

import (
	"errors"
	"fmt"
	"strconv"

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
	db.First(&existing, "ID = ?", u.ID)

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
func (u *User) SendMessage(db *gorm.DB, message *UserMessage) {
	// Ensure that the message is from us
	message.SenderId = strconv.Itoa(u.ID)

	db.Create(message)
}

func (u *User) ReadMessages(db *gorm.DB, senderId string) *[]UserMessage {
	userMessages := []UserMessage{}
	db.Where("receiver_id = ? AND sender_id = ?", strconv.Itoa(u.ID), senderId).Find(&userMessages)

	return &userMessages
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
