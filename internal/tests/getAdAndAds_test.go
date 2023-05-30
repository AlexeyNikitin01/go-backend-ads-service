package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	response, err := client.getAd(0)
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(1))
	assert.False(t, response.Data.Published)
}

func TestGetAdErr(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	_, err = client.getAd(10)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestSearchByName(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	response, err := client.searchAdByName("h")
	assert.NoError(t, err)
	assert.Zero(t, response.Data[0].ID)
	assert.Equal(t, response.Data[0].Title, "hello")
	assert.Equal(t, response.Data[0].Text, "world")
	assert.Equal(t, response.Data[0].AuthorID, int64(1))
	assert.False(t, response.Data[0].Published)
}

func TestListAdsAuthor(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	response, err := client.createAd(1, "best cat", "not for sale")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(1, response.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.listAdsAuthor(1)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestListAdsAuthorErr(t *testing.T) {
	client := getTestClient()

	_, err := client.listAdsAuthor(1)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestListAdsDate(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	response, err := client.createAd(1, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAdsDate(int64(response.Data.CreateDate.Day()))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, response.Data.ID)
	assert.Equal(t, ads.Data[0].Title, response.Data.Title)
	assert.Equal(t, ads.Data[0].Text, response.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, response.Data.AuthorID)
	assert.False(t, ads.Data[0].Published)
}

func TestListAdsDateErr(t *testing.T) {
	client := getTestClient()

	_, err := client.listAdsDate(int64(0))
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	u, err := client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	response, err := client.createAd(u.Data.UserID, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.createAd(u.Data.UserID, "hello", "world")
	assert.NoError(t, err)

	ad, err := client.deleteAd(response.Data.ID, response.Data.AuthorID)

	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, response.Data.ID)
	assert.Equal(t, ad.Data.Title, response.Data.Title)
	assert.Equal(t, ad.Data.Text, response.Data.Text)
	assert.Equal(t, ad.Data.AuthorID, response.Data.AuthorID)
	assert.False(t, ad.Data.Published)
}
