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
	return repositories.GetAllBooks()
}

func GetOneBook(id int) (models.Book, error) {
	return repositories.GetOneBook(id)
}

func DeleteBook(id int) (models.Book, error) {
	return repositories.DeleteBook(id)
}

func AddNewBook(book models.Book) (models.Book, error) {
	return repositories.AddNewBook(book)
}

func UpdateBook(id int, updates models.Book) (models.Book, error) {
	return repositories.UpdateBook(id, updates)
}

func AddItemMongoDB(book *models.Book) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Collection("books")
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
