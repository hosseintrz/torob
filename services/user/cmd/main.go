package main

import (
	"github.com/hosseintrz/torob/user/config"
	db "github.com/hosseintrz/torob/user/internal/db"
	"github.com/hosseintrz/torob/user/internal/persistence"
	"github.com/hosseintrz/torob/user/internal/service"
	pb "github.com/hosseintrz/torob/user/pb/user"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	conf := config.GetConfig()
	database := db.NewDatabase(conf.DbType, conf.DbConn)
	err := database.Connect()
	if err != nil {
		log.Fatalf("unable to connect to db, %v\n", err)
	}
	repo := persistence.NewMongoLayer(database.(*db.MongoDB))

	userService := service.NewUserService(repo)
	lis, err := net.Listen("tcp", "localhost:8181")
	if err != nil {
		log.Fatal("couldn't start tcp server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, userService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %s", err)
	}
}
