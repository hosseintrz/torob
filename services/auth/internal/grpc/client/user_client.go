package client

import (
	"context"
	pb "github.com/hosseintrz/torob/auth/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type UserServiceClient struct {
	Client pb.UserClient
}

func InitUserServiceClient(url string) *UserServiceClient {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldn't connect to user service: %s\n", err)
	}
	client := &UserServiceClient{
		Client: pb.NewUserClient(conn),
	}
	return client
}

func (c *UserServiceClient) AddUser(fullName, email, username, password string, role int32) (string, error) {
	userMsg := &pb.UserMsg{
		Fullname: fullName,
		Email:    email,
		Username: username,
		Password: password,
		Role:     pb.UserMsg_Role(role),
	}
	_, err := c.Client.AddUser(context.Background(), userMsg)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (c *UserServiceClient) GetUser(username string) (*pb.UserMsg, error) {
	req := &pb.UserRequest{
		Username: username,
	}
	user, err := c.Client.GetUser(context.Background(), req)
	return user, err
}

func (c *UserServiceClient) DeleteUser(username string) (string, error) {
	req := &pb.UserRequest{
		Username: username,
	}
	res, err := c.Client.DeleteUser(context.Background(), req)
	return res.Message, err
}

func (c *UserServiceClient) UpdateUser(fullName, email, username, password string, role int32) (string, error) {
	userMsg := &pb.UserMsg{
		Fullname: fullName,
		Email:    email,
		Username: username,
		Password: password,
		Role:     pb.UserMsg_Role(role),
	}
	_, err := c.Client.UpdateUser(context.Background(), userMsg)
	if err != nil {
		return "", err
	}
	return username, nil
}
