package tests

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzTestCreateUser_Fuzz(f *testing.F) {
	testcases := []string{
		"alex",
		"gopher",
	}

	for _, tc := range testcases {
		f.Add(tc)
	}
	
	client := getTestClient()

	f.Fuzz(func(t *testing.T, s string) {
		log.Println(s)
		u, err := client.createUser(s, "gopher@tin.com")
		assert.NoError(t, err)
		got := u.Data.NickName
		expect := s

		if got != expect {
			t.Errorf("For (%v) Expect: %v, but got: %v", s, expect, got)
		}
	})
}