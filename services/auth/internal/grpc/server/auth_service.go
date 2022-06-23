package server

import (
	"context"
	"fmt"
	"github.com/hosseintrz/torob/auth/conf"
	"github.com/hosseintrz/torob/auth/internal"
	"github.com/hosseintrz/torob/auth/internal/grpc/client"
	"github.com/hosseintrz/torob/auth/internal/utils"
	pb "github.com/hosseintrz/torob/auth/pb/auth"
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
	_, err = a.UserGrpc.AddUser(
		dto.Fullname,
		dto.Email,
		dto.Username,
		dto.Password,
		int32(dto.Role))
	if err != nil {
		return nil, err
	}
	token, err := internal.GetSignedToken()
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

	token, err := internal.GetSignedToken()
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{Response: token}, nil
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
