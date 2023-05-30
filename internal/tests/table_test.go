package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableTestCreateAd(t *testing.T) {
	client := getTestClient()

	u, err := client.createUser("nick", "@")
	assert.NoError(t, err)

	type AdField struct {
		UserId int64
		Title string
		Text string
	}

	type Test struct {
		Name   string
		In     AdField
		Expect AdField
	}

	tests := [...]Test{
		{"Correct field",
		AdField{u.Data.UserID, "title", "text"},
		AdField{u.Data.UserID, "title", "text"},
		},
		{"Correct field 2",
			AdField{u.Data.UserID, "hello", "hello_text"},
			AdField{u.Data.UserID, "hello", "hello_text"},
		},
		{"Correct field 3",
			AdField{u.Data.UserID, "Gopher", "textGopher"},
			AdField{u.Data.UserID, "Gopher", "textGopher"},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			ad, err := client.createAd(test.In.UserId, test.In.Title, test.In.Text)
			assert.NoError(t, err)
			if ad.Data.AuthorID != test.Expect.UserId ||
			ad.Data.Title != test.Expect.Title ||
			ad.Data.Text != test.Expect.Text {
			t.Fatalf(`test %q: expect %v got %v`, test.Name, test.Expect, ad)
		}
		})
	}
}

func TestTableCreateUser(t *testing.T) {
	client := getTestClient()

	type UserInField struct {
		NickName string
		Email string
	}

	type UserExpectField struct {
		UserId int64
		NickName string
		Email string
		Activate bool
	}

	type Test struct {
		Name string
		In UserInField
		Expect UserExpectField
	}

	tests := []Test{
		{"Field create user",
		UserInField{"Gopher", "Gopher@tin.com"},
		UserExpectField{int64(0),"Gopher", "Gopher@tin.com", false}},
		{"Field create user",
		UserInField{"Alexey", "golanger@tin.com"},
		UserExpectField{int64(1), "Alexey", "golanger@tin.com", false}},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {
			user, err := client.createUser(test.In.NickName, test.In.Email)
			assert.NoError(t, err)
			if user.Data.UserID != test.Expect.UserId ||
			user.Data.NickName != test.Expect.NickName ||
			user.Data.Email != test.Expect.Email ||
			user.Data.Activate != test.Expect.Activate {
			t.Fatalf(`test %q: expect %v got %v`, test.Name, test.Expect, user)
		}
		})
	}
}
