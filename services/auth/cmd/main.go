package main

import (
	"github.com/hosseintrz/torob/auth/internal/grpc/client"
	"github.com/hosseintrz/torob/auth/internal/grpc/server"
	pb "github.com/hosseintrz/torob/auth/pb/auth"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	userGrpc := client.InitUserServiceClient("localhost:8181")
	authService := &server.AuthService{
		UserGrpc: userGrpc,
	}
	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("couldn't start auth tcp server: %s", err.Error())
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, authService)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("auth server error : %s", err.Error())
	}
}
