package main

import "gorm.io/gorm"

type User struct {
	Username string
	Password string
}

func (u *User) Create(db *gorm.DB) {
	create(db, u)
}

func create(db *gorm.DB, u *User) {
	db.Create(u)
}

// func (u *User) delete(db *gorm.DB) {
// 	db.Delete(u)
// }
