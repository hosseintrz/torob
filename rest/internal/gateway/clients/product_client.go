package clients

import (
	"context"
	"fmt"
	"github.com/hosseintrz/torob/rest/conf"
	pb "github.com/hosseintrz/torob/rest/pb/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ProductClient struct {
	Conf   *conf.Config
	Client pb.ProductClient
}

func NewProductClient(conf *conf.Config) *ProductClient {
	return &ProductClient{Conf: conf}
}
func (c *ProductClient) Connect() {
	addr := fmt.Sprintf("%s:%s", c.Conf.Host, c.Conf.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldn't connect to category service: %s\n", err)
	}
	c.Client = pb.NewProductClient(conn)
}

func (c *ProductClient) CreateCategory(name string, parent string) (string, error) {
	req := pb.CategoryRequest{Name: name, Parent: parent}
	res, err := c.Client.CreateCategory(context.Background(), &req)
	if err != nil {
		return "", nil
	}
	return res.Message, nil
}

func (c *ProductClient) CreateProduct(name string, category string, fields map[string]string) (string, error) {
	req := &pb.CreateProductReq{
		Name:     name,
		Category: category,
		Fields:   fields,
	}
	res, err := c.Client.CreateProduct(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.Message, nil
}

func (c *ProductClient) GetProduct(id string) (*pb.GetProductRes, error) {
	req := &pb.GetProductReq{
		Id: id,
	}
	res, err := c.Client.GetProduct(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return res, nil

}
