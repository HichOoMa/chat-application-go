package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hichoma.chat.dev/internal/config"
)

var db *mongo.Database

// Initialize database connection
func InitializeDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	optionsClient := options.Client().ApplyURI(config.AppConfig.MongoURI)
	client, err := mongo.Connect(ctx, optionsClient)
	if err != nil {
		log.Fatal("connection failed")
	}
	db = client.Database(config.AppConfig.DatabaseName)
	return err
}

/*
Usage :

	user := User{...}
	CreateCollection("users", &user)
*/
func CreateCollection(collectionName string, model interface{}) (string, error) {
	collection := db.Collection(collectionName)
	result, err := collection.InsertOne(context.Background(), model)
	if err != nil {
		return "", err
	}
	oid, _ := result.InsertedID.(primitive.ObjectID)
	return oid.Hex(), err
}

/*
Usage :

	user := User{}
	filter := bson.M{"name": "test"}
	FindOneCollection("users" ,filter ,&user)
*/
func FindOneCollection(collectionName string, filter interface{}, result interface{}) error {
	collection := db.Collection(collectionName)

	err := collection.FindOne(context.Background(), filter).Decode(result)
	return err
}

/*
Usage :

	filter := bson.M{"name": "test"}
	users, err := FindCollection("users" ,filter)
*/

func FindCollection(collectionName string, filter interface{}) (*mongo.Cursor, error) {
	collection := db.Collection(collectionName)

	result, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*
Usage :

	user := User{}
	filter := bson.M{"name": "test"}
	update := bson.M{"name": "tests"}
	UpdateCollection("users" ,filter ,update)
*/
func UpdateCollection(collectionName string, filter interface{}, update interface{}) error {
	collection := db.Collection(collectionName)

	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

/*
Usage :

	user := User{}
	filter := bson.M{"name": "test"}
	DeleteCollection("users" ,filter)
*/
func DeleteCollection(collectionName string, filter interface{}) error {
	collection := db.Collection(collectionName)

	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}
