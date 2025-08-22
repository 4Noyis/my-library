package repositories

import (
	"context"
	"time"

	"github.com/4Noyis/my-library/internal/database"
	"github.com/4Noyis/my-library/internal/logger"
	"github.com/4Noyis/my-library/internal/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BookRepository struct {
	collection string
}

func NewBookCollection() *BookRepository {
	return &BookRepository{
		collection: "books",
	}
}

// func GetMongoCollection(client *mongo.Client) *mongo.Collection {
// 	return client.Database("library").Collection("books")
// }

func GetAllBooks() ([]models.Book, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Collection("books")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		logger.LogDatabaseOperation("find_all", "books", nil, time.Since(start).Milliseconds(), err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	if err = cursor.All(ctx, &books); err != nil {
		logger.LogDatabaseOperation("find_all", "books", nil, time.Since(start).Milliseconds(), err)
		return nil, err
	}

	logger.LogDatabaseOperation("find_all", "books", nil, time.Since(start).Milliseconds(), nil)
	logger.LogDebug("Retrieved all books from database", logrus.Fields{
		"count":    len(books),
		"duration": time.Since(start).Milliseconds(),
	})

	return books, nil
}

func GetOneBook(id int) (models.Book, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Collection("books")
	var book models.Book
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&book)

	logger.LogDatabaseOperation("find_one", "books", id, time.Since(start).Milliseconds(), err)

	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func AddNewBook(book models.Book) (models.Book, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Collection("books")
	// Find the highest existing ID
	var lastBook models.Book
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})

	err := collection.FindOne(ctx, bson.D{}, opts).Decode(&lastBook)
	nextID := 1
	if err == nil {
		nextID = lastBook.ID + 1
	} else if err != mongo.ErrNoDocuments {
		logger.LogError("AddNewBook", err, logrus.Fields{
			"operation": "find_last_id",
		})
		return book, err
	}

	// Set the auto-incremented ID and timestamps
	book.ID = nextID
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	ss, err := collection.InsertOne(ctx, book)

	logger.LogDatabaseOperation("insert", "books", book.ID, time.Since(start).Milliseconds(), err)

	if err != nil {
		logger.LogError("AddNewBook", err, logrus.Fields{
			"operation": "insert_one",
			"book_id":   book.ID,
		})
		return book, err
	}

	logger.LogInfo("Book added successfully", logrus.Fields{
		"book_id":     book.ID,
		"inserted_id": ss.InsertedID,
		"title":       book.Title,
		"duration_ms": time.Since(start).Milliseconds(),
	})

	return book, nil
}

func UpdateBook(id int, updates models.Book) (models.Book, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Collection("books")

	// First check if book exists
	var existingBook models.Book
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&existingBook)
	if err != nil {
		logger.LogDatabaseOperation("find_for_update", "books", id, time.Since(start).Milliseconds(), err)
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
		logger.LogDatabaseOperation("update", "books", id, time.Since(start).Milliseconds(), err)
		return models.Book{}, err
	}

	// Return the updated book
	var updatedBook models.Book
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedBook)
	if err != nil {
		logger.LogError("UpdateBook", err, logrus.Fields{
			"operation": "retrieve_updated",
			"book_id":   id,
		})
		return models.Book{}, err
	}

	logger.LogDatabaseOperation("update", "books", id, time.Since(start).Milliseconds(), nil)
	logger.LogInfo("Book updated successfully", logrus.Fields{
		"book_id":     id,
		"title":       updatedBook.Title,
		"duration_ms": time.Since(start).Milliseconds(),
	})

	return updatedBook, nil
}

func DeleteBook(id int) (models.Book, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Collection("books")
	var book models.Book
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&book)
	if err != nil {
		logger.LogDatabaseOperation("find_for_delete", "books", id, time.Since(start).Milliseconds(), err)
		return models.Book{}, err
	}

	deleted, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		logger.LogDatabaseOperation("delete", "books", id, time.Since(start).Milliseconds(), err)
		return models.Book{}, err
	}

	if deleted.DeletedCount == 0 {
		logger.LogError("DeleteBook", mongo.ErrNoDocuments, logrus.Fields{
			"operation":     "delete_one",
			"book_id":       id,
			"deleted_count": deleted.DeletedCount,
		})
		return models.Book{}, mongo.ErrNoDocuments
	}

	logger.LogDatabaseOperation("delete", "books", id, time.Since(start).Milliseconds(), nil)
	logger.LogInfo("Book deleted successfully", logrus.Fields{
		"book_id":       id,
		"title":         book.Title,
		"deleted_count": deleted.DeletedCount,
		"duration_ms":   time.Since(start).Milliseconds(),
	})

	return book, nil
}
