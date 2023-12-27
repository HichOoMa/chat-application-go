package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"hichoma.chat.dev/internal/database"
	"hichoma.chat.dev/internal/models"
)

type friendReq struct {
	ID string `json:"_id"`
}

func AddUserFriend(ctx echo.Context) error {
	userId := strings.Join(ctx.Request().Header.Values("user"), "")

	bodyDTO := friendReq{}
	ctx.Bind(&bodyDTO)
	friendId, err := primitive.ObjectIDFromHex(bodyDTO.ID)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	friend := models.Friend{}
	err = database.FindOneCollection("users", bson.M{"_id": friendId}, &friend)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "friend id doesn't exist")
	}
	if userId == friend.ID.Hex() {
		return ctx.String(http.StatusBadRequest, "cannot be friend with yourself")
	}

	userID, _ := primitive.ObjectIDFromHex(userId)
	err = database.UpdateCollection(
		"users",
		bson.M{"_id": userID},
		bson.M{"$push": bson.M{"friends": friend}},
	)
	if err != nil {
		return ctx.String(http.StatusConflict, "can't add friend")
	}

	return ctx.String(http.StatusNoContent, "")
}

func GetUserFriendList(ctx echo.Context) error {
	userId := strings.Join(ctx.Request().Header.Values("user"), "")
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	user := models.User{}
	filter := bson.M{"_id": userIdObj}
	err = database.FindOneCollection("users", filter, &user)
	if err != nil {
		return ctx.String(http.StatusConflict, "user collection not found")
	}
	return ctx.JSON(http.StatusOK, user.Friends)
}
