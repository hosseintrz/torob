package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserDto struct {
	FullName string `json:"fullname"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int32  `json:"role"`
}
type CredentialDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s RestServer) Signup(c echo.Context) error {
	fmt.Println("in signup")
	u := new(UserDto)
	if err := c.Bind(u); err != nil {
		return err
	}
	res, err := s.authGrpc.Signup(u.FullName, u.Email, u.UserName, u.Password, u.Role)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (s RestServer) Login(c echo.Context) error {
	dto := new(CredentialDto)
	if err := c.Bind(dto); err != nil {
		return err
	}
	res, err := s.authGrpc.Login(dto.Username, dto.Password)
	if err != nil {
		return c.String(404, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
