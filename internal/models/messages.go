package models

type Message struct {
	SenderID       string      `json:"senderId" bson:"senderId"`
	SenderUsername string      `json:"senderUsername" bson:"senderUsername"`
	Content        string      `json:"content" bson:"content"`
	Date           string      `json:"date" bson:"date"`
	Reactions      []Reactions `json:"reactions" bson:"reactions"`
}

type Reactions struct {
	Emoji  string   `json:"emoji" bson:"emoji"`
	Number int      `json:"number" bson:"number"`
	Users  []string `json:"users" bson:"users"`
}
