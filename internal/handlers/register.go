package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"hichoma.chat.dev/internal/database"
	"hichoma.chat.dev/internal/models"
	"hichoma.chat.dev/pkg/jwt"
)

type tokenResponse struct {
	Token string `json:"token"`
}

func Register(c echo.Context) error {
	userFrom := new(models.UserSignUpForm)

	// parse request into struct
	err := json.NewDecoder(c.Request().Body).Decode(&userFrom)
	if err != nil || userFrom.Username == "" || userFrom.Email == "" || userFrom.Password == "" {
		return c.String(http.StatusBadRequest, "bad request from client")
	}

	// verif if user existance
	userTest := new(models.User)
	database.FindCollection("users", bson.M{"email": userFrom.Email}, &userTest)
	if userTest.Email == userFrom.Email {
		return c.JSON(http.StatusConflict, "email already exist")
	}

	// create user collection
	newUser := models.User{
		Name:  userFrom.Username,
		Email: userFrom.Email,
	}
	newUser.HashPassword(userFrom.Password)

	// create user in database
	userID, err := database.CreateCollection("users", newUser)
	if err != nil || userID == "" {
		fmt.Println(err)
		return c.JSON(http.StatusConflict, "can't create user")
	}

	// generate token for user
	token, err := jwt.GenerateToken(userID, newUser.Email, newUser.HashedPassword)
	if err != nil {
		return c.String(http.StatusInternalServerError, "generate token failed")
	}
	// parse token
	tokenObject := tokenResponse{Token: token}
	return c.JSON(http.StatusAccepted, tokenObject)
}
