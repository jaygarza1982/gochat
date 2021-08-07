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

func arrayContainsString(array *[]string, value string) bool {
	for _, v := range *array {
		if v == value {
			return true
		}
	}

	return false
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
	user2 := User{}
	faker.FakeData(&user2)
	user3 := User{}
	faker.FakeData(&user3)

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

	// User 0 reads message they just sent
	user0Messages := user0.ReadMessages(db, user1.Username)

	if (*user0Messages)[0].MessageText != userMessage.MessageText {
		t.Errorf("user could not read messages")
	}

	// Other users have a conversation
	// Send a message to user 2 from user 3
	userMessage2 := UserMessage{ReceiverId: user2.Username, MessageText: "Message to user 2"}
	user3.SendMessage(db, &userMessage2, nil)

	// User 2 reads messages
	user2Messages := user2.ReadMessages(db, user3.Username)

	if (*user2Messages)[0].MessageText != userMessage2.MessageText {
		t.Errorf("user 3 could not read messages")
	}

	// Ensure other users cannot see other messages
	user1NewMessages := user1.ReadMessages(db, user3.Username)

	if len(*user1NewMessages) != 0 {
		t.Errorf("possible message leak: user1 has messages from user3")
	}
}

// TODO: Ensure that conversations do not "leak" to other users
// In other words, ensure that conversations are only listed if the user has messages addressed to them
// from a specific user
func TestUser_GetConversations(t *testing.T) {
	db := resetDBConnection(t)

	// Fake a user and password
	user0 := User{Username: "user0"}
	user1 := User{Username: "user1"}
	user2 := User{Username: "user2"}
	user3 := User{Username: "user3"}
	user4 := User{Username: "user4"}
	user5 := User{Username: "user5"}

	if err := user0.Register(db, "12345"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}
	if err := user1.Register(db, "54321"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}
	if err := user2.Register(db, "123user2"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}
	if err := user3.Register(db, "123user3"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}
	if err := user4.Register(db, "123user3"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}
	if err := user5.Register(db, "123user3"); err != nil {
		t.Errorf("got error when login and should not have %v", err.Error())
	}

	// Send message to user 1
	userMessage := UserMessage{ReceiverId: user1.Username, MessageText: "Message to user 1 for convo test"}
	user0.SendMessage(db, &userMessage, nil)

	conversations := user1.GetConversations(db)

	if conversations[0] != user0.Username {
		t.Errorf("got error when getting user0 conversations should have %v", user0.Username)
	}

	// Send another message to user 1 from a different user
	user1NewMessage := UserMessage{ReceiverId: user1.Username, MessageText: "Message from user 2"}
	user2.SendMessage(db, &user1NewMessage, nil)

	// Get new conversations
	conversations = user1.GetConversations(db)

	if !arrayContainsString(&conversations, user0.Username) {
		t.Errorf("got error when getting user1 conversations should have %v", user0.Username)
	}
	if !arrayContainsString(&conversations, user2.Username) {
		t.Errorf("got error when getting user1 conversations should have %v", user2.Username)
	}

	// Send user 3 a message from user 2
	user2MessageFromUser3 := UserMessage{ReceiverId: user3.Username, MessageText: "Message from user 2"}
	user2.SendMessage(db, &user2MessageFromUser3, nil)

	user3Conversations := user3.GetConversations(db)
	if len(user3Conversations) != 1 {
		t.Errorf("user 3 conversations was invalid length, should have %v", 1)
	}
	if user3Conversations[0] != user2.Username {
		t.Errorf("got error when getting user3 conversations should have %v", user2.Username)
	}

	// Send message to user 1
	user4To5Message := UserMessage{ReceiverId: user5.Username}
	user4.SendMessage(db, &user4To5Message, nil)

	// User 4 lists conversations, ensure that we have user 5 in our conversations
	user4Conversations := user4.GetConversations(db)
	if user4Conversations[0] != user5.Username {
		t.Errorf("user 4 did not have user 5 in conversations")
	}
}
