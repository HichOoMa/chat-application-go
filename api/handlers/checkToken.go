package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"hichoma.chat.dev/pkg/jwt"
)

type checkTokenForm struct {
}

func CheckToken(ctx echo.Context) error {
	token := ctx.Request().Header.Values("token")
	if token == nil {
		ctx.String(http.StatusUnauthorized, "unauthorized")
	}

	claims, err := jwt.PasreToken(strings.Join(token, ""))
	if err != nil {
		ctx.String(http.StatusUnauthorized, "unauthorized")
	}

	err = claims.StandardClaims.Valid()
	if err != nil {
		ctx.String(http.StatusUnauthorized, "unauthorized")
	}

	isValid := claims.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true)
	if !isValid {
		ctx.String(http.StatusUnauthorized, "unauthorized")
	}

	return ctx.String(http.StatusNoContent, "")
}
