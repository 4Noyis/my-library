package repositories

import (
	"context"
	"log"
	"time"

	"github.com/4Noyis/my-library/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetMongoCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("library").Collection("books")
}

func GetAllBooks(client *mongo.Client) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetMongoCollection(client)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return books, nil
}

func GetOneBook(client *mongo.Client, id int) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetMongoCollection(client)
	var book models.Book
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&book)
	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func DeleteBook(client *mongo.Client, id int) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetMongoCollection(client)
	var book models.Book
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&book)
	if err != nil {
		log.Println("book id not found")
		return models.Book{}, err
	}
	deleted, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		log.Println("can not delete book on collection")
		return models.Book{}, err
	}
	if deleted.DeletedCount == 0 {
		log.Println("no document was deleted")
		return models.Book{}, mongo.ErrNoDocuments
	}

	log.Println("book deleted successfully", deleted.DeletedCount)
	return book, nil
}
