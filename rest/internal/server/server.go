package server

import (
	"fmt"
	"github.com/hosseintrz/torob/rest/internal/gateway"
	"github.com/hosseintrz/torob/rest/internal/middleware"
	"github.com/labstack/echo/v4"
)

type RestServer struct {
	echo     *echo.Echo
	authGrpc *gateway.AuthClient
	conf     *Config
}

func New(conf *Config, ac *gateway.AuthClient) *RestServer {
	server := &RestServer{
		authGrpc: ac,
		conf:     conf,
	}
	server.echo = echo.New()
	server.echo.HideBanner = true
	server.setupRoutes()
	return server
}

func (s *RestServer) setupRoutes() {
	s.echo.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return middleware.JwtValidation(next, s.authGrpc)
	})

	s.echo.POST("/signup", s.Signup)
	s.echo.POST("/login", s.Login)
	s.echo.GET("/test", func(c echo.Context) error {
		return c.String(200, "something")
	})
	s.echo.GET("/nice", func(c echo.Context) error {
		fmt.Printf("current user is : %s\n", c.Get("user"))
		return c.String(200, "nice path")
	})
}

func (s *RestServer) Serve() chan error {
	addr := fmt.Sprintf("%s:%s", s.conf.Host, s.conf.Port)
	errChan := make(chan error)

	go func() {
		errChan <- s.echo.Start(addr)
	}()
	return errChan
}
