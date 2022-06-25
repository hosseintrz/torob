package middleware

import (
	"fmt"
	"github.com/hosseintrz/torob/rest/internal/gateway/clients"
	"github.com/labstack/echo/v4"
	"net/http"
)

func JwtValidation(next echo.HandlerFunc, ac *clients.AuthClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("entered middleware")
		if uri := c.Request().RequestURI; uri == "/signup" || uri == "/login" {
			return next(c)
		}
		if _, ok := c.Request().Header["Token"]; !ok {
			return c.String(http.StatusUnauthorized, "token is missing")
		}
		token := c.Request().Header["Token"][0]
		res, err := ac.ValidateToken(token)
		if err != nil {
			return c.String(http.StatusUnauthorized, "token is invalid")
		}
		c.Set("user", res)
		return next(c)
	}
}
