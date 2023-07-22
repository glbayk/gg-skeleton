package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func MongoClientConnect() *mongo.Client {
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	c, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	// defer func() {
	// 	if err = MongoClient.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	var result bson.M
	if err := c.Database(os.Getenv("MONGO_INITDB_DATABASE")).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return c
}

type LogMessageType struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Message   string             `bson:"message,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
}

func Create(message LogMessageType) error {
	coll := MongoClient.Database(os.Getenv("MONGO_INITDB_DATABASE")).Collection("logs")
	_, err := coll.InsertOne(context.TODO(), message)
	return err
}

func Read(id string) error {
	var res LogMessageType
	coll := MongoClient.Database(os.Getenv("MONGO_INITDB_DATABASE")).Collection("logs")
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&res)
	return err
}
