package tests

import (
	"net/http/httptest"
	"time"

	"testing"

	"ads/internal/ads"
	"ads/internal/ports/httpgin"
	"ads/internal/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func mockClient(a *mocks.App) *testClient{
	server := httpgin.NewHTTPServer(":18080", a)
	testServer := httptest.NewServer(server.Handler)
	return &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func TestMockCreateAd(t *testing.T) {
	a := &mocks.App{}
	a.On("CheckUser", mock.Anything, int64(0)).Return(nil)
	a.On("CreateAd", mock.Anything, "title", "text", int64(0)).Return(&ads.Ad{
		ID: 0,
		AuthorID: 0,
		Title: "title",
		Text: "text",
	}, nil)
	client := mockClient(a)

	response, err := client.createAd(0, "title", "text")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "title")
	assert.Equal(t, response.Data.Text, "text")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)

}

func TestMockChangeAdStatus(t *testing.T) {
	a := &mocks.App{}
	a.On("ChangeAdStatus", mock.Anything, int64(0), true, int64(0)).Return(
		&ads.Ad{
			ID: 0,
			AuthorID: 0,
			Title: "title",
			Text: "text",
			Published: true,
		}, nil,
	)
	client := mockClient(a)

	response, err := client.changeAdStatus(int64(0), int64(0), true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)
}

func TestMockUpdateAd(t *testing.T) {
	a := &mocks.App{}
	a.On("UpdateAd", mock.Anything, int64(0), "привет", "мир", int64(0)).Return(
		&ads.Ad{
			ID: 0,
			AuthorID: 0,
			Title: "привет",
			Text: "мир",
			Published: true,
		}, nil,
	)
	client := mockClient(a)

	response, err := client.updateAd(int64(0), int64(0), "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestMockAdService_HandlerGetAdsDate(t *testing.T) {
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

	client := mockClient(a)

	ads, err := client.listAdsDate(day)
	
	assert.NoError(t, err)
	assert.Equal(t, ads.Data[0].ID, int64(0))
}