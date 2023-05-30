package tests

import (
	"context"
	"log"
	"testing"

	grpcPort "ads/internal/ports/grpc"
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
