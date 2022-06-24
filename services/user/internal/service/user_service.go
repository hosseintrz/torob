package service

import (
	"context"
	"fmt"
	"github.com/hosseintrz/torob/user/internal/model"
	"github.com/hosseintrz/torob/user/internal/persistence"
	pb "github.com/hosseintrz/torob/user/pb/user"

	"github.com/hosseintrz/torob/user/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type UserService struct {
	pb.UnimplementedUserServer
	Repository persistence.UserRepository
}

func NewUserService(repo persistence.UserRepository) *UserService {
	return &UserService{Repository: repo}
}

func (u *UserService) AddUser(ctx context.Context, userMsg *pb.UserMsg) (*pb.UserResponse, error) {
	_, err := u.Repository.GetUser(userMsg.GetUsername())
	if err == nil {
		return nil, errors.ErrUserExists
	}
	id, err := u.Repository.AddUser(model.NewUser(
		"",
		userMsg.Fullname,
		userMsg.Username,
		userMsg.Email,
		userMsg.Password,
		time.Now(),
		int32(userMsg.Role),
	))
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Message: fmt.Sprintf("user with id %d added successfully", id),
	}, nil
}

func (u *UserService) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserMsg, error) {
	user, err := u.Repository.GetUser(req.GetUsername())
	if err != nil {
		return nil, err
	}
	fmt.Println("userrole: ", user.Role)
	fmt.Println("role : ", pb.UserMsg_Role(user.Role))
	return &pb.UserMsg{
		Fullname:    user.FullName,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		CreatedDate: timestamppb.New(user.CreatedDate),
		Role:        pb.UserMsg_Role(user.Role),
	}, nil
}

func (u *UserService) DeleteUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	_, err := u.Repository.DeleteUser(req.GetUsername())
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Message: fmt.Sprintf("user %s deleted successfully", req.GetUsername()),
	}, nil
}

func (u *UserService) UpdateUser(ctx context.Context, userMsg *pb.UserMsg) (*pb.UserResponse, error) {
	user := model.NewUser(
		"",
		userMsg.Fullname,
		userMsg.Username,
		userMsg.Email,
		userMsg.Password,
		time.Now(),
		int32(userMsg.Role))

	_, err := u.Repository.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Message: fmt.Sprintf("user %s updated successfully", userMsg.Username),
	}, nil
}
