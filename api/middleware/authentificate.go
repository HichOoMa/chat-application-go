package middlewareFc

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"hichoma.chat.dev/pkg/jwt"
)

func Authentificate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		token := c.Request().Header.Values("token")
		if token == nil {
			c.String(http.StatusUnauthorized, "unauthorized")
		}
		claims, err := jwt.PasreToken(strings.Join(token, ""))
		if err != nil {
			c.Error(err)
		}
		err = claims.StandardClaims.Valid()
		if err != nil {
			c.Error(err)
		}
		isValid := claims.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true)
		if !isValid {
			c.String(http.StatusUnauthorized, "unauthorized")
		}
		c.Request().Header.Set("user", claims.UserID)
		return next(c)
	}
}
