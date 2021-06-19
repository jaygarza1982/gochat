package main

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUser_Create(t *testing.T) {
	// TODO: Get db, user, and password from env vars
	dsn := "host=localhost user=gochat password=ch4ng3m3p13453 dbname=gochat port=5432 sslmode=disable TimeZone=US/Eastern"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Errorf("Error was not nil")
		t.Errorf(err.Error())
	}

	// TODO: Create users array with different values
	user0 := User{"jake", "123"}
	user0.Create(db)
}
