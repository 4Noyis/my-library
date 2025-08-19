package pkg

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/4Noyis/my-library/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func ConnectMongoDB() error {

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return errors.New("Cannot get uri address")
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

func GetMongoCollection() *mongo.Collection {
	return client.Database("library").Collection("books")
}
func genereteBook() models.Book {
	return models.Book{
		ID:          1,
		ISBN:        "978-0-452-28423-4",
		Title:       "1984",
		Author:      "George Orwell",
		Publisher:   "Plume",
		PublishedAt: time.Date(1949, 6, 8, 0, 0, 0, 0, time.UTC),
		Genre:       "Dystopian Fiction",
		Language:    "English",
		Pages:       328,
		Description: "A dystopian social science fiction novel about totalitarian control and surveillance.",
		CoverURL:    "https://example.com/covers/1984.jpg",
		Location:    "C-3-08",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func AddItemMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetMongoCollection()
	_, err := collection.InsertOne(ctx, genereteBook())
	if err != nil {
		log.Println("inserting item error:", err)
	}

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
