package tests

import (
	"context"
	"log"
	"testing"

	grpcPort "ads/internal/ports/grpc"
	"ads/internal/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func BenchmarkRestCreateUser(b *testing.B) {
	log.Println("START ^o^ :: BenchMarkRestCreateUser")
	client := getTestClient()

	for i := 0; i <= b.N; i++ {
		_, _ = client.createUser("gopher", "gopher@tin.com")
	}
}

func BenchmarkGRPCCreateUserTest(b *testing.B) {
	client, _ := getClientGRPC()

	for i := 0; i <= b.N; i++ {
		_, _ = client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "gopher"})
	}
}

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

type TestSuite struct {
	suite.Suite
	u	user.User
}

func (suite *TestSuite) SetupTest() {
	suite.u = user.User{
		NickName: "Alex",
		Email: "alex@tin.com",
	}
}

func (suite *TestSuite) TestExample() {
	client := getTestClient()
	u, err := client.createUser("Alex", "alex@tin.com")
	assert.NoError(suite.T(), err)

	suite.Equal(u.Data.UserID, suite.u.UserID)
	suite.Equal(u.Data.NickName, suite.u.NickName)
	suite.Equal(u.Data.Email, suite.u.Email)
	suite.Equal(u.Data.Activate, suite.u.Activate)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

type RESTgRPRCSuite struct {
	suite.Suite
	clientREST *testClient 
	clientGRPC grpcPort.AdServiceClient
}

func (suite *RESTgRPRCSuite) SetupClient() {
	suite.clientREST = getTestClient()
	suite.clientGRPC, _ = getClientGRPC()
}

func (suite *RESTgRPRCSuite) TestRestAndGRPC() {
	cREST := suite.clientREST
	cGRPC := suite.clientGRPC

	uREST, errREST := cREST.createUser("Alex", "Alex@tin.com")
	suite.NoError(errREST)
	uGRPC, errGRPC := cGRPC.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Alex"})
	suite.NoError(errGRPC)
	suite.Equal(uREST.Data.NickName, uGRPC.Name)
}

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
