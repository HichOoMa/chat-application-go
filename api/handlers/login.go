package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"hichoma.chat.dev/internal/database"
	"hichoma.chat.dev/internal/models"
	"hichoma.chat.dev/pkg/jwt"
)

func Login(ctx echo.Context) error {
	loginForm := new(models.UserLoginForm)

	// parse request into object
	err := json.NewDecoder(ctx.Request().Body).Decode(&loginForm)
	if err != nil || loginForm.Email == "" || loginForm.Password == "" {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	// verify credentials
	// verify email exist
	user := new(models.User)
	database.FindCollection("users", bson.M{"email": loginForm.Email}, &user)
	if user == new(models.User) {
		return ctx.String(http.StatusBadRequest, "user doesn't exist")
	}
	isAuth, _ := user.PasswordMatch(loginForm.Password)

	if isAuth {
		token, err := jwt.GenerateToken(user.ID.Hex(), user.Email, user.HashedPassword)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, "generate token failed")
		}
		tokenObject := models.TokenResponse{Token: token}
		return ctx.JSON(http.StatusAccepted, tokenObject)
	} else {
		return ctx.JSON(http.StatusUnauthorized, "credential doesn't match")
	}
}
