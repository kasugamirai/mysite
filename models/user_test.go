package models_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"xy.com/mysite/models"
)

var testDB *gorm.DB

func setup() {
	var err error
	testDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect to test database")
	}
	err = testDB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Print(err)
	}
}

func TestUserModel(t *testing.T) {
	setup()

	t.Run("CreateUser", testCreateUser)
	t.Run("GetUserByID", testGetUserByID)
	t.Run("GetUserByUsername", testGetUserByUsername)
	t.Run("GetUserByEmail", testGetUserByEmail)
	t.Run("UpdateUser", testUpdateUser)
	t.Run("DeleteUser", testDeleteUser)
	t.Run("AuthenticateUser", testAuthenticateUser)
}

func testCreateUser(t *testing.T) {
	user := &models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password",
	}
	err := models.CreateUser(testDB, user)
	assert.NoError(t, err)
	assert.NotEqual(t, user.ID, 0)
}

func testGetUserByID(t *testing.T) {
	user, err := models.GetUserByID(testDB, 1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.ID, uint(1))
}

func testGetUserByUsername(t *testing.T) {
	user, err := models.GetUserByUsername(testDB, "testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Username, "testuser")
}

func testGetUserByEmail(t *testing.T) {
	user, err := models.GetUserByEmail(testDB, "testuser@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Email, "testuser@example.com")
}

func testUpdateUser(t *testing.T) {
	user, _ := models.GetUserByID(testDB, 1)
	user.Email = "updated@example.com"
	err := models.UpdateUser(testDB, user)
	assert.NoError(t, err)

	updatedUser, _ := models.GetUserByID(testDB, 1)
	assert.Equal(t, updatedUser.Email, "updated@example.com")
}

func testDeleteUser(t *testing.T) {
	err := models.DeleteUser(testDB, 1)
	assert.NoError(t, err)

	user, err := models.GetUserByID(testDB, 1)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func testAuthenticateUser(t *testing.T) {
	user := &models.User{
		Email:    "authuser@example.com",
		Password: "password",
	}
	err := models.CreateUser(testDB, user)
	if err != nil {
		fmt.Println(err)
	}

	authUser, err := models.AuthenticateUser(testDB, "authuser@example.com", "password")
	assert.NoError(t, err)
	assert.NotNil(t, authUser)
	assert.Equal(t, authUser.Email, "authuser@example.com")
}
