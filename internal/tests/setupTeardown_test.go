package tests

import (
	"context"
	grpcPort "ads/internal/ports/grpc"
	"ads/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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
