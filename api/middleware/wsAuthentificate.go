package middlewareFc

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"hichoma.chat.dev/pkg/jwt"
)

func WsAuthentificate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Values("Sec-WebSocket-Protocol")
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
