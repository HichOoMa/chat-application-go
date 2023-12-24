package middlewareFc

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"hichoma.chat.dev/pkg/jwt"
)

func Authentificate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if err := next(ctx); err != nil {
			ctx.Error(err)
		}

		token := ctx.Request().Header.Values("token")
		if token == nil {
			return echo.ErrUnauthorized
		}

		claims, err := jwt.PasreToken(strings.Join(token, ""))
		if err != nil {
			return echo.ErrUnauthorized
		}

		err = claims.StandardClaims.Valid()
		if err != nil {
			return echo.ErrUnauthorized
		}

		isValid := claims.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true)
		if !isValid {
			return echo.ErrUnauthorized
		}

		ctx.Request().Header.Set("user", claims.UserID)
		return next(ctx)
	}
}
