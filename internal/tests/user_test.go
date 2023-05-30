package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("Alex", "alex@mai.com")
	assert.NoError(t, err)
	fmt.Println(response)
	assert.Equal(t, response.Data.NickName, "Alex")
	assert.Equal(t, response.Data.UserID, int64(0))
	assert.Equal(t, response.Data.Email, "alex@mai.com")
	assert.False(t, response.Data.Activate, false)
}

func TestUserUpdate(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("hello", "world")
	assert.NoError(t, err)

	response, err = client.updateUser(response.Data.NickName, response.Data.Email, response.Data.UserID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Activate)

	response, err = client.updateUser(response.Data.NickName, response.Data.Email, response.Data.UserID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Activate)

	response, err = client.updateUser(response.Data.NickName, response.Data.Email, response.Data.UserID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Activate)
}

func TestDeleteUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("hello", "world")
	assert.NoError(t, err)

	_, err = client.createUser("hello", "world")
	assert.NoError(t, err)

	response, err := client.deleteUser(1)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.UserID, int64(1))
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("hello", "world")
	assert.NoError(t, err)

	_, err = client.createUser("Alex", "alex@mai.com")
	assert.NoError(t, err)

	response, err := client.getUser(1)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.NickName, "Alex")
	assert.Equal(t, response.Data.UserID, int64(1))
	assert.Equal(t, response.Data.Email, "alex@mai.com")
	assert.False(t, response.Data.Activate, false)
}

func TestGetUserErrNotFound(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("hello", "world")
	assert.NoError(t, err)

	_, err = client.getUser(1)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestErrCreateUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestErrDeleteUser(t *testing.T) {
	client := getTestClient()
	_, err := client.deleteUser(int64(1))
	assert.Error(t, err)
}

func TestUserUpdateErr(t *testing.T) {
	client := getTestClient()
	_, err := client.updateUser("NickName", "Email", 10, true)
	assert.ErrorIs(t, err, ErrNotFound)
}
