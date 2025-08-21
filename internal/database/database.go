package database

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func GetClient() *mongo.Client {
	return client
}

func ConnectMongoDB() error {

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return errors.New("cannot get uri address")
	}
	opts := options.Client().ApplyURI(uri)

	Client, err := mongo.Connect(opts)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	client = Client

	log.Println("connected to mongodb succesfully")
	return nil
}

func DisconnectMongoDB() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := client.Disconnect(ctx)
		if err != nil {
			log.Println("disconnect error:", err)
		}
		log.Println("Disconnect from mongodb")
	}

}
