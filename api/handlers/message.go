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
	userId := strings.Join(ctx.Request().Header.Values("user"), "")
	date := time.Now().Format(time.RFC3339)
	emptyReaction := []models.Reactions{}

	msg := models.Message{
		ID:         primitive.NewObjectID(),
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
	userId := strings.Join(ctx.Request().Header.Values("user"), "")
	filter := bson.A{"$or", bson.A{bson.M{"sender_id": userId}, bson.M{"receiver_id": userId}}}

	result := []models.Message{}
	err := database.FindCollection("messages", filter, &result)
	if err != nil {
		return ctx.JSON(http.StatusAccepted, bson.M{})
	}

	return ctx.JSON(http.StatusOK, result)
}
