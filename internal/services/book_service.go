package services

import (
	"context"
	"time"

	"github.com/4Noyis/my-library/internal/database"
	"github.com/4Noyis/my-library/internal/logger"
	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/repositories"
	"github.com/sirupsen/logrus"
)

func GetAllBooks() ([]models.Book, error) {
	return repositories.GetAllBooks(database.GetClient())
}

func GetOneBook(id int) (models.Book, error) {
	return repositories.GetOneBook(database.GetClient(), id)
}

func DeleteBook(id int) (models.Book, error) {
	return repositories.DeleteBook(database.GetClient(), id)
}

func AddNewBook(book models.Book) (models.Book, error) {
	return repositories.AddNewBook(database.GetClient(), book)
}

func UpdateBook(id int, updates models.Book) (models.Book, error) {
	return repositories.UpdateBook(database.GetClient(), id, updates)
}

func AddItemMongoDB(book *models.Book) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := repositories.GetMongoCollection(database.GetClient())
	_, err := collection.InsertOne(ctx, book)
	if err != nil {
		logger.LogError("AddItemMongoDB", err, logrus.Fields{
			"operation": "insert_one",
			"title":     book.Title,
		})
		return
	}

	logger.LogInfo("Book added successfully via AddItemMongoDB", logrus.Fields{
		"operation": "AddItemMongoDB",
		"title":     book.Title,
	})
}
