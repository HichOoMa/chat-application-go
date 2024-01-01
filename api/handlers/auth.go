package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"hichoma.chat.dev/internal/database"
	"hichoma.chat.dev/internal/models"
	"hichoma.chat.dev/pkg/jwt"
)

type authResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func Register(ctx echo.Context) error {
	userFrom := new(models.UserSignUpForm)

	// parse request into object
	err := json.NewDecoder(ctx.Request().Body).Decode(&userFrom)
	if err != nil || userFrom.Username == "" || userFrom.Email == "" || userFrom.Password == "" {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	// verif if user existance
	userTest := new(models.User)
	database.FindOneCollection("users", bson.M{"email": userFrom.Email}, &userTest)
	if userTest.Email == userFrom.Email {
		return ctx.JSON(http.StatusConflict, "email already exist")
	}

	// create user collection
	newUser := models.User{
		ID:    primitive.NewObjectID(),
		Name:  userFrom.Username,
		Email: userFrom.Email,
	}
	newUser.HashPassword(userFrom.Password)

	// create user in database
	userID, err := database.CreateCollection("users", newUser)
	if err != nil || userID == "" {
		return ctx.JSON(http.StatusConflict, "can't create user")
	}

	// generate token for user
	token, err := jwt.GenerateToken(userID, newUser.Email, newUser.HashedPassword)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "generate token failed")
	}
	// parse token
	tokenObject := authResponse{UserID: userID, Username: newUser.Name, Token: token}
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
	database.FindOneCollection("users", bson.M{"email": loginForm.Email}, &user)
	if user == new(models.User) {
		return ctx.String(http.StatusBadRequest, "user doesn't exist")
	}
	isAuth, _ := user.PasswordMatch(loginForm.Password)

	if isAuth {
		token, err := jwt.GenerateToken(user.ID.Hex(), user.Email, user.HashedPassword)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, "generate token failed")
		}
		tokenObject := authResponse{UserID: user.ID.Hex(), Username: user.Name, Token: token}
		return ctx.JSON(http.StatusAccepted, tokenObject)
	} else {
		return ctx.JSON(http.StatusUnauthorized, "credential doesn't match")
	}
}

func CheckToken(ctx echo.Context) error {
	token := ctx.QueryParams().Get("token")
	if token == "" {
		return echo.ErrUnauthorized
	}

	claims, err := jwt.PasreToken(token)
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

	return ctx.String(http.StatusNoContent, "")
}
