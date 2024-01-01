package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `json:"_id"         bson:"_id"`
	SenderID   string             `json:"sender_id"   bson:"sender_id"`
	ReceiverID string             `json:"receiver_id" bson:"receiver_id"`
	Content    string             `json:"content"     bson:"content"`
	Date       string             `json:"date"        bson:"date"`
	Reactions  []Reactions        `json:"reactions"   bson:"reactions"`
}

type Reactions struct {
	Emoji  string   `json:"emoji"  bson:"emoji"`
	Number int      `json:"number" bson:"number"`
	Users  []string `json:"users"  bson:"users"`
}

type WsMessage struct {
	FriendId string `json:"friend_id" bson:"friend_id"`
	Content  string `json:"content"   bson:"content"`
}

func (msg *Message) Validate() bool {
	if msg.ID == primitive.NilObjectID || msg.ReceiverID == "" || msg.SenderID == "" || msg.Content == "" || msg.Date == "" {
		return false
	}
	return true
}
