package clients

import (
	"context"
	"fmt"
	"github.com/hosseintrz/torob/rest/conf"
	pb "github.com/hosseintrz/torob/rest/pb/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type AuthClient struct {
	Conf   *conf.Config
	Client pb.AuthClient
}

func NewAuthClient(conf *conf.Config) *AuthClient {
	return &AuthClient{Conf: conf}
}

func (c *AuthClient) Connect() {
	addr := fmt.Sprintf("%s:%s", c.Conf.Host, c.Conf.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldn't connect to auth service: %s\n", err)
	}
	c.Client = pb.NewAuthClient(conn)
}

func (c *AuthClient) Signup(fullName, email, username, password string, role int32) (*pb.AuthResponse, error) {
	req := &pb.SignupRequest{
		Fullname: fullName,
		Email:    email,
		Username: username,
		Password: password,
		Role:     pb.SignupRequest_Role(role),
	}
	res, err := c.Client.Signup(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *AuthClient) Login(username, password string) (*pb.AuthResponse, error) {
	req := &pb.LoginRequest{
		Username: username,
		Password: password,
	}
	res, err := c.Client.Login(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *AuthClient) ValidateToken(token string) (*pb.ValidationResponse, error) {
	req := &pb.ValidationRequest{
		Token: token,
	}
	payload, err := c.Client.ValidateToken(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
