package main

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func resetDBConnection(t *testing.T) *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_CONN_STRING")), &gorm.Config{})

	if err != nil {
		t.Errorf("Error was not nil")
		t.Errorf(err.Error())
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserMessage{})

	// Clear table
	db.Where("1 = 1").Delete(User{})
	db.Where("1 = 1").Delete(UserMessage{})

	return db
}

func TestUser_Create(t *testing.T) {
	db := resetDBConnection(t)

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

func TestUser_CorrectPassword(t *testing.T) {
	db := resetDBConnection(t)

	// Fake a user and password
	user := User{}
	faker.FakeData(&user)

	if err := user.Register(db, "12345"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}

	// Check for correct and incorrect password
	if !user.CheckPassword(db, "12345") {
		t.Errorf("user did not have correct password and should have")
	}

	if user.CheckPassword(db, "123456") {
		t.Errorf("user had correct password and should not have")
	}

	// Check password of nonexistent user
	user1 := User{Username: "jake-not-here"}
	if user1.CheckPassword(db, "123") {
		t.Errorf("nonexistent user was able to login")
	}
}

func TestUser_SendMessage(t *testing.T) {
	db := resetDBConnection(t)

	// Fake a user and password
	user0 := User{}
	faker.FakeData(&user0)
	user1 := User{}
	faker.FakeData(&user1)

	if err := user0.Register(db, "12345"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}
	if err := user1.Register(db, "54321"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}

	// Send message to user 1
	userMessage := UserMessage{ReceiverId: user1.Username, MessageText: "Message to user 1"}
	user0.SendMessage(db, &userMessage, nil)

	// User 1 reads messages
	user1Messages := user1.ReadMessages(db, user0.Username)

	if (*user1Messages)[0].MessageText != userMessage.MessageText {
		t.Errorf("user could not read messages")
	}

	// TODO: Add more users with more messages
	// TODO: Ensure messages do not "leak" when other users request their own messages
}

// TODO: Make test with different users and different conversations
// Ensure that conversations do not "leak" to other users
// In other words, ensure that conversations are only listed if the user has messages addressed to them
// from a specific user
func TestUser_GetConversations(t *testing.T) {
	db := resetDBConnection(t)

	// Fake a user and password
	user0 := User{}
	faker.FakeData(&user0)
	user1 := User{}
	faker.FakeData(&user1)

	if err := user0.Register(db, "12345"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}
	if err := user1.Register(db, "54321"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}

	// Send message to user 1
	userMessage := UserMessage{ReceiverId: user1.Username, MessageText: "Message to user 1 for convo test"}
	user0.SendMessage(db, &userMessage, nil)

	conversations := user1.GetConversations(db)

	if conversations[0] != user0.Username {
		t.Errorf("got error when getting conversations should have %v", user0.Username)
	}
}
