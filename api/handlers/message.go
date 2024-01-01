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
		ReceiverID: message.FriendId,
		Content:    message.Content,
		Date:       date,
		Reactions:  emptyReaction,
	}
	if !msg.Validate() {
		return nil
	}
	messageId, err := database.CreateCollection("messages", &msg)
	if err != nil {
		return nil
	}

	msg.ID, _ = primitive.ObjectIDFromHex(messageId)

	return &msg
}

func GetUserChat(ctx echo.Context) error {
	userId := strings.Join(ctx.Request().Header.Values("user"), "")
	chatFriendId := ctx.Param("friend")

	filter := bson.D{{
		Key: "$and",
		Value: bson.A{
			bson.D{{
				Key: "sender_id",
				Value: bson.D{{
					Key: "$in",
					Value: bson.A{
						userId,
						chatFriendId,
					},
				}},
			}},
			bson.D{{
				Key: "receiver_id",
				Value: bson.D{{
					Key: "$in",
					Value: bson.A{
						userId,
						chatFriendId,
					},
				}},
			}},
		},
	}}

	response := []models.Message{}
	result, err := database.FindCollection("messages", filter)
	if err != nil {
		return ctx.JSON(http.StatusNoContent, bson.M{})
	}

	for result.Next(ctx.Request().Context()) {
		var message models.Message
		result.Decode(&message)
		response = append(response, message)
	}

	return ctx.JSON(http.StatusOK, response)
}
