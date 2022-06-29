package server

import (
	"fmt"
	"github.com/hosseintrz/torob/rest/conf"
	"github.com/hosseintrz/torob/rest/internal/gateway/clients"
	"github.com/hosseintrz/torob/rest/internal/middleware"
	"github.com/labstack/echo/v4"
)

type RestServer struct {
	echo         *echo.Echo
	authGrpc     *clients.AuthClient
	productGrpc  *clients.ProductClient
	supplierGrpc *clients.SupplierClient
	conf         *conf.Config
}

func New(conf *conf.Config, authConf *conf.Config, prodConf *conf.Config, suppConf *conf.Config) *RestServer {
	server := &RestServer{
		conf:         conf,
		authGrpc:     clients.NewAuthClient(authConf),
		productGrpc:  clients.NewProductClient(prodConf),
		supplierGrpc: clients.NewSupplierClient(suppConf),
	}

	server.authGrpc.Connect()
	server.productGrpc.Connect()
	server.supplierGrpc.Connect()

	server.echo = echo.New()
	server.echo.HideBanner = true
	server.setupRoutes()
	return server
}

func (s *RestServer) setupRoutes() {
	s.echo.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return middleware.JwtValidation(next, s.authGrpc)
	})

	//auth service
	s.echo.POST("/signup", s.Signup)
	s.echo.POST("/login", s.Login)

	//product service
	s.echo.POST("/categories", s.CreateCategory)
	s.echo.POST("/products", s.CreateProduct)
	s.echo.GET("/products/:id", s.GetProduct)

	//supplier service
	s.echo.POST("/stores", s.AddStore)
	s.echo.POST("/offers", s.SubmitOffer)
	s.echo.GET("/offers/:prodId", s.GetProductOffers)

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
