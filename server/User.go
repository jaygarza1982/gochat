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

func checkPassword(password string) bool {

	return false
	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// return err == nil
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
