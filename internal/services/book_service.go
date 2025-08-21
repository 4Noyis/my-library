package services

import (
	"context"
	"log"
	"time"

	"github.com/4Noyis/my-library/internal/database"
	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/repositories"
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

func AddItemMongoDB(book *models.Book) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := repositories.GetMongoCollection(database.GetClient())
	_, err := collection.InsertOne(ctx, book)
	if err != nil {
		log.Println("inserting item error:", err)
	}
	log.Println("new book added succesfully")

}
