package database

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/4Noyis/my-library/internal/logger"
	"github.com/sirupsen/logrus"
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

	logger.LogDebug("Attempting to connect to MongoDB", logrus.Fields{
		"operation": "ConnectMongoDB",
		"uri_set":   uri != "",
	})

	opts := options.Client().ApplyURI(uri)

	Client, err := mongo.Connect(opts)
	if err != nil {
		logger.LogError("ConnectMongoDB", err, logrus.Fields{
			"operation": "mongo.Connect",
		})
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Client.Ping(ctx, nil)
	if err != nil {
		logger.LogError("ConnectMongoDB", err, logrus.Fields{
			"operation": "ping",
		})
		return err
	}

	client = Client

	logger.LogInfo("Connected to MongoDB successfully", logrus.Fields{
		"operation": "ConnectMongoDB",
		"database":  "library",
	})
	return nil
}

func DisconnectMongoDB() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := client.Disconnect(ctx)
		if err != nil {
			logger.LogError("DisconnectMongoDB", err, logrus.Fields{
				"operation": "disconnect",
			})
		} else {
			logger.LogInfo("Disconnected from MongoDB", logrus.Fields{
				"operation": "DisconnectMongoDB",
			})
		}
	}
}

func Collection(collectionName string) *mongo.Collection {
	return client.Database("library").Collection(collectionName)
}
