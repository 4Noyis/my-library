package main

import (
	"log"
	"net/http"
	"time"

	database "github.com/4Noyis/my-library/internal/database"
	"github.com/4Noyis/my-library/internal/handlers"
	"github.com/4Noyis/my-library/internal/models"
)

func main() {

	err := database.ConnectMongoDB()
	if err != nil {
		log.Fatal("failed to connect to mongodb:", err)
	}
	defer database.DisconnectMongoDB()
	http.HandleFunc("/books", handlers.BookHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("server error: ", err)
	}

	// newBook := generateBook()
	// services.AddItemMongoDB(&newBook)

}

func generateBook() models.Book {
	return models.Book{
		ID:          3,
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
