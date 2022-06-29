package clients

import (
	"context"
	"fmt"
	"github.com/hosseintrz/torob/rest/conf"
	pb "github.com/hosseintrz/torob/rest/pb/supplier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

type SupplierClient struct {
	Conf   *conf.Config
	Client pb.SupplierClient
}

func NewSupplierClient(conf *conf.Config) *SupplierClient {
	return &SupplierClient{Conf: conf}
}
func (c *SupplierClient) Connect() {
	addr := fmt.Sprintf("%s:%s", c.Conf.Host, c.Conf.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldn't connect to supplier service: %s\n", err)
	}
	c.Client = pb.NewSupplierClient(conn)
}

func (c *SupplierClient) SubmitOffer(storeId, productId, url, description string, price int32) (string, error) {
	req := &pb.OfferReq{
		StoreId:     storeId,
		ProductId:   productId,
		Price:       price,
		Url:         url,
		Description: description,
	}
	res, err := c.Client.SubmitOffer(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.GetResponse(), nil
}

func (c *SupplierClient) AddStore(ownerId, storeName, storeUrl, city string) (string, error) {
	req := &pb.AddStoreReq{
		OwnerId:   ownerId,
		StoreName: storeName,
		StoreUrl:  storeUrl,
		City:      city,
	}
	res, err := c.Client.AddStore(context.Background(), req)
	if err != nil {
		return "", nil
	}
	return res.Response, nil
}

func (c *SupplierClient) GetProductOffers(productId string) ([]*pb.ProdOfferRes, error) {
	offers := make([]*pb.ProdOfferRes, 0)
	req := &pb.ProdOfferReq{
		ProductId: productId,
	}
	stream, err := c.Client.GetProductOffers(context.Background(), req)
	if err != nil {
		return offers, err
	}
	for {
		offer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		//		fmt.Println(offer.StoreName)
		offers = append(offers,offer)
	}
	return offers,nil
}
