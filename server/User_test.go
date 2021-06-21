package main

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
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

	users := [10]User{}

	for i := 0; i < len(users); i++ {
		user := &users[i]
		if err := faker.FakeData(user); err != nil {
			t.Errorf("could not fake data %v", err.Error())
		}

		fakePw := ""
		if err := faker.FakeData(&fakePw); err != nil {
			t.Errorf("Could not set fake pw %v", err.Error())
		}

		if createError := user.Register(db, fakePw); createError != nil {
			t.Errorf("Could not create user %v %v", user.Username, createError.Error())
		}

	}

	// Try to recreate users
	for i := 0; i < len(users); i++ {
		user := &users[i]

		fakePw := ""
		if err := faker.FakeData(&fakePw); err != nil {
			t.Errorf("could not set fake pw %v", err.Error())
		}

		if createError := user.Register(db, fakePw); createError == nil {
			t.Errorf("Created user again and should not have")
		}
	}
}
