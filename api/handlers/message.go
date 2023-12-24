package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"hichoma.chat.dev/internal/database"
	"hichoma.chat.dev/internal/models"
)

func AddNewMessage(message *models.WsMessage, ctx echo.Context) *models.Message {
	userId := strings.Join(ctx.Request().Header["user"], "")
	date := time.Now().Format("2001-06-06 10:10:10")
	emptyReaction := []models.Reactions{}

	msg := models.Message{
		SenderID:   userId,
		ReceiverID: message.OppositeId,
		Content:    message.Content,
		Date:       date,
		Reactions:  emptyReaction,
	}

	messageId, err := database.CreateCollection("messages", &msg)
	if err != nil {
		return nil
	}

	msg.ID, _ = primitive.ObjectIDFromHex(messageId)

	return &msg
}

func GetUserMessages(ctx echo.Context) error {
	userId := strings.Join(ctx.Request().Header["user"], "")
	filter := bson.A{"$or", bson.A{bson.M{"sender_id": userId}, bson.M{"receiver_id": userId}}}

	result := []models.Message{}
	err := database.FindCollection("messages", filter, &result)
	if err != nil {
		return ctx.JSON(http.StatusAccepted, bson.M{})
	}

	return ctx.JSON(http.StatusOK, result)
}
