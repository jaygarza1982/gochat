package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID           int `gorm:"primaryKey"`
	Username     string
	PasswordHash string
	password     string
}

func (u *User) CanCreate(rows int64) bool {
	return rows == 0
}

func (u *User) Create(db *gorm.DB) error {
	existing := User{Username: u.Username}

	if !u.CanCreate(db.First(&existing).RowsAffected) {
		return errors.New("could not create user username already exists")
	}

	create(db, u)

	return nil
}

func create(db *gorm.DB, u *User) {
	db.Create(u)
}

func (u *User) Delete(db *gorm.DB) {
	fmt.Printf("Deleting %v", u.ID)
	db.Delete(&u, u.ID)
}
