package main

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUser_Create(t *testing.T) {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_CONN_STRING")), &gorm.Config{})

	if err != nil {
		t.Errorf("Error was not nil")
		t.Errorf(err.Error())
	}

	db.AutoMigrate(&User{})

	// Clear table
	db.Where("1 = 1").Delete(User{})

	// TODO: Create users array with different values
	user0 := User{Username: "jake", PasswordHash: "123"}
	if createError := user0.Create(db); createError != nil {
		t.Errorf("Could not create user %v %v", user0.Username, createError.Error())
	}

	if createError := user0.Create(db); createError == nil {
		t.Errorf("Created user again and should not have")
	}
}
