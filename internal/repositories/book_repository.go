package repositories

import (
	"context"
	"log"
	"time"

	"github.com/4Noyis/my-library/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func AddNewBook(client *mongo.Client, book models.Book) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetMongoCollection(client)

	// Find the highest existing ID
	var lastBook models.Book
	opts := options.FindOne().SetSort(bson.D{{"id", -1}})

	err := collection.FindOne(ctx, bson.D{}, opts).Decode(&lastBook)
	nextID := 1
	if err == nil {
		nextID = lastBook.ID + 1
	} else if err != mongo.ErrNoDocuments {
		log.Println("error finding last book:", err)
		return book, err
	}

	// Set the auto-incremented ID and timestamps
	book.ID = nextID
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	ss, err := collection.InsertOne(ctx, book)
	if err != nil {
		log.Println("book can not added")
		return book, err
	}
	log.Println("book added successfully new book id:", ss.InsertedID)
	return book, nil
}

func UpdateBook(client *mongo.Client, id int, updates models.Book) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := GetMongoCollection(client)
	
	// First check if book exists
	var existingBook models.Book
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&existingBook)
	if err != nil {
		log.Println("book id not found")
		return models.Book{}, err
	}

	// Build update document dynamically based on provided fields
	updateDoc := bson.M{}
	
	if updates.ISBN != "" {
		updateDoc["isbn"] = updates.ISBN
	}
	if updates.Title != "" {
		updateDoc["title"] = updates.Title
	}
	if updates.Author != "" {
		updateDoc["author"] = updates.Author
	}
	if updates.Publisher != "" {
		updateDoc["publisher"] = updates.Publisher
	}
	if !updates.PublishedAt.IsZero() {
		updateDoc["published_at"] = updates.PublishedAt
	}
	if updates.Genre != "" {
		updateDoc["genre"] = updates.Genre
	}
	if updates.Language != "" {
		updateDoc["language"] = updates.Language
	}
	if updates.Pages != 0 {
		updateDoc["pages"] = updates.Pages
	}
	if updates.Description != "" {
		updateDoc["description"] = updates.Description
	}
	if updates.CoverURL != "" {
		updateDoc["coverURL"] = updates.CoverURL
	}
	if updates.Location != "" {
		updateDoc["location"] = updates.Location
	}
	
	// Always update the updated_at field
	updateDoc["updated_at"] = time.Now()

	// Perform the update
	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": updateDoc})
	if err != nil {
		log.Println("failed to update book:", err)
		return models.Book{}, err
	}

	// Return the updated book
	var updatedBook models.Book
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedBook)
	if err != nil {
		log.Println("failed to retrieve updated book:", err)
		return models.Book{}, err
	}

	log.Println("book updated successfully with id:", id)
	return updatedBook, nil
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

	log.Println("book deleted successfully deletedCount: ", deleted.DeletedCount)
	return book, nil
}
