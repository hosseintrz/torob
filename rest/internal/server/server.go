package server

import (
	"fmt"
	"github.com/hosseintrz/torob/rest/conf"
	"github.com/hosseintrz/torob/rest/internal/gateway/clients"
	"github.com/hosseintrz/torob/rest/internal/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
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

	//s.echo.Pre(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"*"},
	//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
	//		echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin,
	//		echo.HeaderAuthorization, echo.HeaderAccessControlAllowHeaders,
	//		echo.HeaderXRequestedWith,
	//	},
	//	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	//}))
	s.echo.Pre(middleware.CORS())

	type Endpoint struct {
		Method string
		Url    string
	}

	skipper := func(skipList []Endpoint) func(c echo.Context) bool {
		return func(c echo.Context) bool {
			for _, endpoint := range skipList {
				if strings.HasPrefix(c.Request().RequestURI, endpoint.Url) &&
					c.Request().Method == endpoint.Method {
					return true
				}
			}
			return false
		}
	}

	skipList := []Endpoint{
		{Method: http.MethodPost, Url: "/signup"},
		{Method: http.MethodPost, Url: "/login"},
		{Method: http.MethodGet, Url: "/categories"},
		{Method: http.MethodGet, Url: "/products"},
		{Method: http.MethodGet, Url: "/products/:id"},
	}

	s.echo.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return middlewares.JwtValidation(next, s.authGrpc, skipper(skipList))
	})

	//auth service
	s.echo.POST("/signup", s.Signup)
	s.echo.POST("/login", s.Login)
	s.echo.GET("/validateToken", s.ValidateToken)

	//product service
	s.echo.POST("/categories", s.CreateCategory)
	s.echo.GET("/categories", s.GetCategories)
	s.echo.POST("/products", s.CreateProduct)
	s.echo.GET("/products/:id", s.GetProduct)
	s.echo.GET("/products", s.GetProductsByType)

	//supplier service
	s.echo.POST("/stores", s.AddStore)
	s.echo.POST("/offers", s.SubmitOffer)
	s.echo.GET("/offers/:prodId", s.GetProductOffers)
	s.echo.GET("/stores", s.GetOwnerStores)
	s.echo.GET("/stores/:storeId", s.GetStoreInfo)

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
