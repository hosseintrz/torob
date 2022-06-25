package main

import (
	"github.com/hosseintrz/torob/product/config"
	"github.com/hosseintrz/torob/product/internal/db"
	"github.com/hosseintrz/torob/product/internal/services"
	pb "github.com/hosseintrz/torob/product/pb/product"
	"github.com/hosseintrz/torob/product/persistence/mongo"
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
	ProductRepo := mongo.NewProductRepo(database.(*db.MongoDB))

	productService := services.NewProductService(ProductRepo)

	lis, err := net.Listen("tcp", "localhost:8282")
	if err != nil {
		log.Fatal("couldn't start tcp server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServer(grpcServer, productService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %s", err)
	}

}
