package server

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strings"
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

func (s RestServer) CreateCategory(c echo.Context) error {
	dto := new(struct {
		Name   string
		Parent string
	})
	if err := c.Bind(dto); err != nil {
		c.Error(err)
	}
	res, err := s.productGrpc.CreateCategory(dto.Name, dto.Parent)
	if err != nil {
		c.Error(err)
	}
	return c.JSON(200, res)
}

func (s RestServer) CreateProduct(c echo.Context) error {
	body := c.Request().Body
	dto := new(struct {
		Name     string            `json:"name"`
		Category string            `json:"category"`
		Fields   map[string]string `json:"fields"`
	})
	err := json.NewDecoder(body).Decode(&dto)
	fmt.Println(dto.Fields)
	if err != nil {
		return c.String(http.StatusBadRequest, "wrong data format")
	}
	res, err := s.productGrpc.CreateProduct(dto.Name, dto.Category, dto.Fields)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(200, res)
}

func (s RestServer) GetProduct(c echo.Context) error {
	url := c.Request().URL
	paths := strings.Split(url.String(), "/")
	id := paths[len(paths)-1]
	res, err := s.productGrpc.GetProduct(id)
	if err != nil {
		return c.String(404, err.Error())
	}
	fmt.Println(res)
	fields := make(map[string]string)
	err = mapstructure.Decode(res.Fields, &fields)
	fields["id"] = res.GetId()
	fields["name"] = res.GetName()
	fields["category"] = res.GetCategory()
	if err != nil {
		return c.String(404, err.Error())
	}
	return c.JSONPretty(200, fields, " ")
}
