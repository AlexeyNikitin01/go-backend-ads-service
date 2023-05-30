package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(100, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("Gopher Gopherich", "gopher@go.com") 
	assert.NoError(t, err)

	_, err = client.createUser("Gopher Goshevich", "gopher@go.com") 
	assert.NoError(t, err)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(1, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(1, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}
