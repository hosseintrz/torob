package middlewares

import (
	"github.com/hosseintrz/torob/rest/internal/gateway/clients"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

func JwtValidation(next echo.HandlerFunc, ac *clients.AuthClient, skipper middleware.Skipper) echo.HandlerFunc {
	return func(c echo.Context) error {
		if skipper(c) {
			return next(c)
		}
		//if uri := c.Request().RequestURI; uri == "/signup" || uri == "/login" {
		//		return next(c)
		//	}
		if _, ok := c.Request().Header["Authorization"]; !ok {
			return c.String(http.StatusUnauthorized, "authorization header is missing")
		}
		authHeader := c.Request().Header["Authorization"][0]
		str := strings.Split(authHeader, " ")
		if str[0] != "JWT" {
			return c.String(http.StatusUnauthorized, "authentication method not supported")
		}
		res, err := ac.ValidateToken(str[1])
		if err != nil {
			return c.String(http.StatusUnauthorized, "token is invalid")
		}
		c.Set("user", res)
		return next(c)
	}
}
