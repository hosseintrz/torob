package main

import (
	"github.com/hosseintrz/torob/supplier/config"
	"github.com/hosseintrz/torob/supplier/internal/db"
	"github.com/hosseintrz/torob/supplier/internal/services"
	pb "github.com/hosseintrz/torob/supplier/pb/supplier"
	"github.com/hosseintrz/torob/supplier/persistence/mongo"
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
	SupplierRepo := mongo.NewSupplierRepo(database.(*db.MongoDB))

	productService := services.NewSupplierService(SupplierRepo)

	lis, err := net.Listen("tcp", "localhost:9191")
	if err != nil {
		log.Fatal("couldn't start tcp server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSupplierServer(grpcServer, productService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %s", err)
	}

}
