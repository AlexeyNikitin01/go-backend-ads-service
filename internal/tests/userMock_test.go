package tests

import (
	"ads/internal/tests/mocks"
	"fmt"
	"ads/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	a := &mocks.App{}
	a.On("CreateUser", mock.Anything, "alex", "alex@mai.com").Return(&user.User{
		UserID: 0,
		NickName: "alex",
		Email: "alex@mai.com",
	}, nil)
	client := mockClient(a)
	response, _ := client.createUser("alex", "alex@mai.com")
	fmt.Println("response", response)
	assert.Equal(t, response.Data.NickName, "alex")
	assert.Equal(t, response.Data.UserID, int64(0))
	assert.Equal(t, response.Data.Email, "alex@mai.com")
	assert.False(t, response.Data.Activate, false)
}

func TestUserUpdate(t *testing.T) {
	a := &mocks.App{}
	a.On("UpdateUser", mock.Anything, "alex", "alex@mai.com", int64(0), true).Return(&user.User{
		UserID: 0,
		NickName: "alex",
		Email: "alex@mai.com",
		Activate: true,
	}, nil)
	client := mockClient(a)

	response, err := client.updateUser("alex", "alex@mai.com", int64(0), true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Activate)
}

func TestDeleteUser(t *testing.T) {
	a := &mocks.App{}
	a.On("CheckUser", mock.Anything, int64(0)).Return(nil)
	a.On("DeleteUser", mock.Anything, int64(0)).Return(nil)
	client := mockClient(a)

	response, err := client.deleteUser(0)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.UserID, int64(0))
}

func TestGetUser(t *testing.T) {
	a := &mocks.App{}
	a.On("GetUser", mock.Anything, int64(1)).Return(&user.User{
		UserID: 1,
		NickName: "alex",
		Email: "alex@mai.com",
		Activate: false,
	}, nil)
	client := mockClient(a)

	response, err := client.getUser(1)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.NickName, "alex")
	assert.Equal(t, response.Data.UserID, int64(1))
	assert.Equal(t, response.Data.Email, "alex@mai.com")
	assert.False(t, response.Data.Activate)
}
