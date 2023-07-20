package tests

import (
	"context"
	"net"
	"testing"
	"time"

	"ads/internal/adapters/adrepo"
	"ads/internal/adapters/userrepo"
	"ads/internal/app"
	grpcPort "ads/internal/ports/grpc"
	"ads/internal/tests/mocks"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

func getClientGRPC() (grpcPort.AdServiceClient, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lis := bufconn.Listen(1024 * 1024)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcPort.UnaryServerInterceptorPanicMethod),
		grpc.ChainUnaryInterceptor(grpcPort.UnaryServerInterceptorLogMethod),
	)
	defer srv.Stop()

	logrus.SetFormatter(new(logrus.JSONFormatter))

	db := &mocks.RepositoryDbUser{}

	a := app.NewApp(adrepo.New(), userrepo.New(), db)

	svc := grpcPort.NewService(a)
	grpcPort.RegisterAdServiceServer(srv, svc)

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, _ := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpcPort.UnaryClientInterceptor),
	)
	defer conn.Close()

	client := grpcPort.NewAdServiceClient(conn)

	return client, ctx
}

func Client(t *testing.T) (grpcPort.AdServiceClient, context.Context) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcPort.UnaryServerInterceptorPanicMethod),
		grpc.ChainUnaryInterceptor(grpcPort.UnaryServerInterceptorLogMethod),
	)
	t.Cleanup(func() {
		srv.Stop()
	})

	logrus.SetFormatter(new(logrus.JSONFormatter))

	db := &mocks.RepositoryDbUser{}

	a := app.NewApp(adrepo.New(), userrepo.New(), db)

	svc := grpcPort.NewService(a)
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpcPort.UnaryClientInterceptor),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	return grpcPort.NewAdServiceClient(conn), ctx
}

func TestGRRPCCreateUser(t *testing.T) {
	client, ctx := Client(t)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
}

func TestGRRPCCreateUserErr(t *testing.T) {
	client, ctx := Client(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: ""})
	assert.Error(t, err)
}

func TestGRRPCCreateAd(t *testing.T) {
	client, ctx := Client(t)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: u.Id})
	assert.NoError(t, err, "client.CreateAd")

	assert.Equal(t, "world", res.Text)
	assert.Equal(t, "hello", res.Title)
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, int64(0), res.AuthorId)
	assert.Equal(t, false, res.Published)
}

func TestGRRPCCreateAdErr(t *testing.T) {
	client, ctx := Client(t)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "", Text: "", UserId: u.Id})
	assert.Error(t, err)
}

func TestGRRPCChangeAdStatus(t *testing.T) {
	client, ctx := Client(t)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")
	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: u.Id})
	assert.NoError(t, err, "client.CreateAd")
	res, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: ad.AuthorId, AdId: ad.Id, Published: true})
	assert.NoError(t, err, "client.ChangeStatusAd")

	assert.Equal(t, ad.Text, res.Text)
	assert.Equal(t, ad.Title, res.Title)
	assert.Equal(t, ad.Id, res.Id)
	assert.Equal(t, ad.AuthorId, res.AuthorId)
	assert.Equal(t, true, res.Published)
}

func TestGRRPCUpdateAd(t *testing.T) {
	client, ctx := Client(t)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")
	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: u.Id})
	assert.NoError(t, err, "client.CreateAd")
	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: ad.AuthorId, AdId: ad.Id, Title: "", Text: ""})
	assert.Error(t, err)
}

func TestGRRPCListAds(t *testing.T) {
	client, ctx := Client(t)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: u.Id})
	assert.NoError(t, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: ad.AuthorId, AdId: ad.Id, Published: true})
	assert.NoError(t, err, "client.ChangeStatusAd")

	ad, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: u.Id})
	assert.NoError(t, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: ad.AuthorId, AdId: ad.Id, Published: true})
	assert.NoError(t, err, "client.ChangeStatusAd")

	ads, err := client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err, "client.ChangeStatusAd")

	assert.Len(t, ads.GetList(), 2)
}

func TestGRRPCGetUser(t *testing.T) {
	client, ctx := Client(t)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")

	u, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, u.Name, "alex")
	assert.Equal(t, u.Id, int64(0))
}

func TestGRRPCGetUserErr(t *testing.T) {
	client, ctx := Client(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 11})
	assert.Error(t, err)
}

func TestGRRPCDeleteUser(t *testing.T) {
	client, ctx := Client(t)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.DeleteUser")
}

func TestGRRPCDeleteUserErr(t *testing.T) {
	client, ctx := Client(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: 11})
	assert.Error(t, err)
}

func TestGRRPCDeleteAd(t *testing.T) {
	client, ctx := Client(t)
	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "alex"})
	assert.NoError(t, err, "client.CreateUser")

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: u.Id})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AuthorId: u.Id, AdId: ad.Id})
	assert.NoError(t, err, "client.GetUser")
}
