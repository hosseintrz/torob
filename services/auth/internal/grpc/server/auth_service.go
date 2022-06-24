package server

import (
	"context"
	"fmt"
	"github.com/hosseintrz/torob/auth/conf"
	"github.com/hosseintrz/torob/auth/internal/grpc/client"
	"github.com/hosseintrz/torob/auth/internal/jwt"
	"github.com/hosseintrz/torob/auth/internal/utils"
	pb "github.com/hosseintrz/torob/auth/pb/auth"
	pb2 "github.com/hosseintrz/torob/auth/pb/user"
	"net/mail"
	"unicode"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	UserGrpc *client.UserServiceClient
}

type AuthError struct {
	Msg string
}

func NewAuthError(msg string) *AuthError {
	return &AuthError{msg}
}
func (a *AuthError) Error() string {
	return fmt.Sprintf("Auth Error : %s", a.Msg)
}

var (
	ErrShortPassword = NewAuthError("Password length should be greater than 8")
	ErrWrongFormat   = NewAuthError("Password should contain one uppercase,one lowercase and a number")
	ErrWrongEmail    = NewAuthError("Email format is incorrect")
	ErrUserNotFound  = NewAuthError("User not found")
	ErrWrongPassword = NewAuthError("Wrong credentials")
)

func (a *AuthService) Signup(ctx context.Context, dto *pb.SignupRequest) (*pb.AuthResponse, error) {
	err := validate(dto)
	if err != nil {
		return nil, err
	}
	secret, err := conf.GetEnv("SECRET")
	if err != nil {
		return nil, err
	}

	enc, err := utils.Encrypt(secret, dto.Password)
	if err != nil {
		return nil, err
	}
	dto.Password = enc
	userMsg := &pb2.UserMsg{
		Fullname: dto.Fullname,
		Email:    dto.Email,
		Username: dto.Username,
		Password: dto.Password,
		Role:     pb2.UserMsg_Role(dto.Role),
	}
	_, err = a.UserGrpc.AddUser(userMsg)
	if err != nil {
		return nil, err
	}
	token, err := jwt.GetSignedToken(userMsg)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Response: token}, nil
}

func (a *AuthService) Login(ctx context.Context, dto *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, err := a.UserGrpc.GetUser(dto.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	fmt.Println("authservice login : user : ", user)
	secret, err := conf.GetEnv("SECRET")
	if err != nil {
		return nil, err
	}

	enc, err := utils.Encrypt(secret, dto.Password)
	if err != nil {
		return nil, err
	}
	dto.Password = enc
	if user.Password != dto.Password {
		return nil, ErrWrongPassword
	}

	token, err := jwt.GetSignedToken(user)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Response: token}, nil
}

func (a *AuthService) ValidateToken(ctx context.Context, req *pb.ValidationRequest) (*pb.ValidationResponse, error) {
	secret, err := conf.GetEnv("JWT_SECRET")
	if err != nil {
		fmt.Println("error getting secret")
		return nil, err
	}
	user, err := jwt.ValidateToken(req.Token, secret)
	if err != nil {
		return nil, err
	}
	//	fmt.Printf("user is :%v\n", user)
	return &pb.ValidationResponse{
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     pb.ValidationResponse_Role(user.Role),
	}, nil
}

func validate(dto *pb.SignupRequest) error {
	if len(dto.Password) < 8 {
		return ErrShortPassword
	}
	var hasUpper bool
	var hasLower bool
	for _, c := range dto.Password {
		if unicode.IsUpper(c) {
			hasUpper = true
		} else if unicode.IsLower(c) {
			hasLower = true
		}
	}
	if !hasUpper || !hasLower {
		return ErrWrongFormat
	}
	_, err := mail.ParseAddress(dto.Email)
	if err != nil {
		return ErrWrongEmail
	}
	return nil
}
