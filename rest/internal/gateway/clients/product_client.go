package clients

import (
	"context"
	"fmt"
	"github.com/hosseintrz/torob/rest/conf"
	pb "github.com/hosseintrz/torob/rest/pb/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
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
		log.Fatalf("Couldn't connect to product service: %s\n", err)
	}
	c.Client = pb.NewProductClient(conn)
}

func (c *ProductClient) CreateCategory(name, parent, desc string) (string, error) {
	req := pb.CategoryRequest{Name: name, Parent: parent, Desc: desc}
	res, err := c.Client.CreateCategory(context.Background(), &req)
	if err != nil {
		return "", nil
	}
	return res.Message, nil
}

func (c *ProductClient) CreateProduct(name, category, imageUrl string, minPrice int32, fields map[string]string) (string, error) {
	req := &pb.CreateProductReq{
		Name:     name,
		Category: category,
		ImageUrl: imageUrl,
		MinPrice: minPrice,
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

func (c *ProductClient) GetProductsByType(category string) ([]*pb.GetProductsRes, error) {
	products := make([]*pb.GetProductsRes, 0)
	req := &pb.GetProductsReq{
		Category: category,
	}
	stream, err := c.Client.GetProductsByType(context.Background(), req)
	if err != nil {
		return nil, err
	}
	for {
		prod, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		products = append(products, prod)
	}
	return products, nil
}

func (c *ProductClient) GetCategories() ([]*pb.GetCategoriesRes, error) {
	categories := make([]*pb.GetCategoriesRes, 0)
	stream, err := c.Client.GetCategories(context.TODO(), &pb.GetCategoriesReq{})
	if err != nil {
		return nil, err
	}
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
