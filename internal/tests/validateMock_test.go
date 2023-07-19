package tests

import (
	"ads/internal/app"
	"ads/internal/tests/mocks"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAd_EmptyTitle(t *testing.T) {
	a := &mocks.App{}
	a.On("CheckUser", mock.Anything, int64(0)).Return(nil)
	a.On("CreateAd", mock.Anything, "", "world", int64(0)).Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.createAd(int64(0), "", "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongTitle(t *testing.T) {
	a := &mocks.App{}
	title := strings.Repeat("a", 101)
	a.On("CheckUser", mock.Anything, int64(0)).Return(nil)
	a.On("CreateAd", mock.Anything, title, "world", int64(0)).Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.createAd(0, title, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_EmptyText(t *testing.T) {
	a := &mocks.App{}
	a.On("CheckUser", mock.Anything, int64(0)).Return(nil)
	a.On("CreateAd", mock.Anything, "title", "", int64(0)).Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.createAd(0, "title", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongText(t *testing.T) {
	a := &mocks.App{}
	text := strings.Repeat("a", 501)
	a.On("CheckUser", mock.Anything, int64(0)).Return(nil)
	a.On("CreateAd", mock.Anything, "title", text, int64(0)).Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.createAd(0, "title", text)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyTitle(t *testing.T) {
	a := &mocks.App{}
	a.On("UpdateAd", mock.Anything, int64(0), "", "new_world", int64(0)).Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.updateAd(int64(0), int64(0), "", "new_world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongTitle(t *testing.T) {
	a := &mocks.App{}
	title := strings.Repeat("a", 101)
	a.On("UpdateAd", mock.Anything, int64(0), title, "world", int64(0)).Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.updateAd(int64(0), int64(0), title, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyText(t *testing.T) {
	a := &mocks.App{}
	a.On("UpdateAd", mock.Anything, int64(0), "title", "", int64(0)).Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.updateAd(int64(0), int64(0), "title", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongText(t *testing.T) {
	a := &mocks.App{}
	text := strings.Repeat("a", 501)
	a.On("UpdateAd", mock.Anything, int64(0), "title", text, int64(0)).Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.updateAd(int64(0), int64(0), "title", text)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestMockListAdsAuthorErr(t *testing.T) {
	a := &mocks.App{}
	a.On("ListAdsAuthor",  mock.Anything, mock.Anything).Return(
		nil, app.ErrBadRequest,
	)
	client := mockClient(a)

	_, err := client.listAdsAuthor(1)
	assert.ErrorIs(t, err, ErrBadRequest)
}


func TestMockGetAdErr(t *testing.T) {
	a := &mocks.App{}
	a.On("GetAd", mock.Anything, int64(10)).Return(nil, app.ErrBadRequest)

	client := mockClient(a)

	_, err := client.getAd(10)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestListAdsDateErr(t *testing.T) {
	a := &mocks.App{}
	day := int64(time.Now().UTC().Day())
	a.On("ListAdsDate", mock.Anything, day).Return(nil, app.ErrBadRequest)

	client := mockClient(a)

	_, err := client.listAdsDate(day)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestGetUserErrNotFound(t *testing.T) {
	a := &mocks.App{}
	a.On("GetUser", mock.Anything, int64(1)).Return(nil, app.ErrNotFound)
	client := mockClient(a)

	_, err := client.getUser(1)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestErrCreateUser(t *testing.T) {
	a := &mocks.App{}
	a.On("CreateUser", mock.Anything, "", "").Return(nil, app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.createUser("", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestErrDeleteUser(t *testing.T) {
	a := &mocks.App{}
	a.On("DeleteUser", mock.Anything, int64(0)).Return(app.ErrBadRequest)
	client := mockClient(a)

	_, err := client.deleteUser(int64(0))
	assert.Error(t, err)
}

func TestUserUpdateErr(t *testing.T) {
	a := &mocks.App{}
	a.On("UpdateUser", mock.Anything, "NickName", "Email", int64(10), true).Return(nil, app.ErrNotFound)
	client := mockClient(a)

	_, err := client.updateUser("NickName", "Email", int64(10), true)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()
	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(100, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}
