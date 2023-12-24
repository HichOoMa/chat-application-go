package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"

	"hichoma.chat.dev/internal/database"
	"hichoma.chat.dev/internal/models"
	"hichoma.chat.dev/pkg/jwt"
)

func Register(ctx echo.Context) error {
	userFrom := new(models.UserSignUpForm)

	// parse request into object
	err := json.NewDecoder(ctx.Request().Body).Decode(&userFrom)
	if err != nil || userFrom.Username == "" || userFrom.Email == "" || userFrom.Password == "" {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	// verif if user existance
	userTest := new(models.User)
	database.FindCollection("users", bson.M{"email": userFrom.Email}, &userTest)
	if userTest.Email == userFrom.Email {
		return ctx.JSON(http.StatusConflict, "email already exist")
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
		return ctx.JSON(http.StatusConflict, "can't create user")
	}

	// generate token for user
	token, err := jwt.GenerateToken(userID, newUser.Email, newUser.HashedPassword)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "generate token failed")
	}
	// parse token
	tokenObject := models.TokenResponse{Token: token}
	return ctx.JSON(http.StatusAccepted, tokenObject)
}

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
