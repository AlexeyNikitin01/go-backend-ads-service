package tests

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"ads/internal/ads"
	"ads/internal/app"
	"ads/internal/ports/httpgin"
	"ads/internal/tests/mocks"
	"ads/internal/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_CreateUser(t *testing.T) {
	ctx := context.Background()

	repoAd := &mocks.RepositryAd{}
	repoUser := &mocks.RepositoryUser{}
	repoUser.
	On("AddUser", mock.Anything, &user.User{NickName: "Alex", Email: "email@tin.com"}).
	Return(int64(0), nil)

	a := app.NewApp(repoAd, repoUser)

	u, err := a.CreateUser(ctx, "Alex", "email@tin.com")

	log.Println(u, err)
	assert.Nil(t, err)
	assert.Equal(t, u.UserID, int64(0))
}

func TestProductService_CreateAd(t *testing.T) {
	ctx := context.Background()

	repoAd := &mocks.RepositryAd{}
	repoUser := &mocks.RepositoryUser{}
	repoUser.
	On("AddUser", mock.Anything, &user.User{NickName: "Alex", Email: "email@tin.com"}).
	Return(int64(0), nil)

	a := app.NewApp(repoAd, repoUser)

	u, err := a.CreateUser(ctx, "Alex", "email@tin.com")

	log.Println(u, err)
	assert.Nil(t, err)
	assert.Equal(t, u.UserID, int64(0))

	repoAd.
	On("Add", mock.Anything, mock.Anything).
	Return(int64(0), nil)

	ad, err := a.CreateAd(ctx, "title", "text", u.UserID)
	log.Println(ad, err)
	assert.Nil(t, err)
	assert.Equal(t, ad.ID, int64(0))
	assert.Equal(t, ad.Title, "title")
	assert.Equal(t, ad.Text, "text")
	assert.Equal(t, ad.AuthorID, u.UserID)
}

func getTestMockClient(a *mocks.App) *testClient {
	server := httpgin.NewHTTPServer(":18080", a)
	testServer := httptest.NewServer(server.Handler)
	return &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func TestAdService_HandlerCreateUser(t *testing.T) {
	a := &mocks.App{}
	a.On("CreateUser", mock.Anything, "alex", "alex@mai.com").Return(&user.User{
		UserID: 0,
		NickName: "alex",
		Email: "alex@mai.com",
	}, nil)
	client := getTestMockClient(a)
	response, _ := client.createUser("alex", "alex@mai.com")
	fmt.Println("response", response)
	assert.Equal(t, response.Data.NickName, "alex")
	assert.Equal(t, response.Data.UserID, int64(0))
	assert.Equal(t, response.Data.Email, "alex@mai.com")
	assert.False(t, response.Data.Activate, false)
}

func TestAdService_HandlerCreateAd(t *testing.T) {
	a := &mocks.App{}
	a.On("CheckUser", mock.Anything, int64(0)).Return(nil)
	a.On("CreateAd", mock.Anything, "hello", "world", int64(0)).Return(&ads.Ad{
		AuthorID: int64(0),
		Title: "hello",
		Text: "world",
	}, nil)
	client := getTestMockClient(a)

	response, err := client.createAd(0, "hello", "world")

	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
}


func TestAdService_HandlerSearchAd(t *testing.T) {
	a := &mocks.App{}

	var result []*ads.Ad
	result = append(result, &ads.Ad{
		AuthorID: int64(1),
		Title: "hello",
		Text: "world",
	})
	result = append(result, &ads.Ad{
		AuthorID: int64(1),
		Title: "hello",
		Text: "world",
	})

	a.On("SearchAdByName", mock.Anything, mock.Anything).Return(result, nil)

	client := getTestMockClient(a)

	ads, err := client.searchAdByName("h")
	
	assert.NoError(t, err)
	assert.Equal(t, ads.Data[0].ID, int64(0))
}

func TestAdService_HandlerAuthorAds(t *testing.T) {
	a := &mocks.App{}

	var result []*ads.Ad
	result = append(result, &ads.Ad{
		AuthorID: int64(1),
		Title: "hello",
		Text: "world",
	})
	result = append(result, &ads.Ad{
		AuthorID: int64(1),
		Title: "hello",
		Text: "world",
	})

	a.On("ListAdsAuthor", mock.Anything, int64(0)).Return(result, nil)

	client := getTestMockClient(a)

	ads, err := client.listAdsAuthor(int64(0))
	
	assert.NoError(t, err)
	assert.Equal(t, ads.Data[0].ID, int64(0))
}

func TestAdService_HandlerGetAds(t *testing.T) {
	a := &mocks.App{}

	var result []*ads.Ad
	result = append(result, &ads.Ad{
		AuthorID: int64(1),
		Title: "hello",
		Text: "world",
	})
	result = append(result, &ads.Ad{
		AuthorID: int64(1),
		Title: "hello",
		Text: "world",
	})

	a.On("ListAds", mock.Anything).Return(result, nil)

	client := getTestMockClient(a)

	ads, err := client.listAds()
	
	assert.NoError(t, err)
	assert.Equal(t, ads.Data[0].ID, int64(0))
}

func TestAdService_HandlerGetAdsDate(t *testing.T) {
	a := &mocks.App{}

	var result []*ads.Ad
	result = append(result, &ads.Ad{
		AuthorID: int64(1),
		Title: "hello",
		Text: "world",
	})
	result = append(result, &ads.Ad{
		AuthorID: int64(1),
		Title: "hello",
		Text: "world",
	})
	day := int64(time.Now().UTC().Day())
	a.On("ListAdsDate", mock.Anything, day).Return(result, nil)

	client := getTestMockClient(a)

	ads, err := client.listAdsDate(day)
	
	assert.NoError(t, err)
	assert.Equal(t, ads.Data[0].ID, int64(0))
}

func TestAdService_HandlerErrCreateUser(t *testing.T) {
	a := &mocks.App{}
	a.On("CreateUser", mock.Anything, "alex", "alex").Return(nil)
	client := getTestMockClient(a)

	_, err := client.createUser("", "")
	
	assert.Error(t, err)
}
