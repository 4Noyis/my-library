package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/services"
)

func BookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		if r.URL.Query().Get("id") != "" {
			getOneBook(w, r)
		} else {
			getAllBooks(w)
		}
	} else if r.Method == "POST" {
		newBook(w, r)
	} else if r.Method == "PATCH" {
		updateBook(w, r)
	} else if r.Method == "DELETE" {
		deleteBook(w, r)
	}

}

func getAllBooks(w http.ResponseWriter) error {
	books, err := services.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return errors.New("server error")
	}

	json.NewEncoder(w).Encode(books)
	return nil
}
func getOneBook(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return errors.New("id parameter required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return errors.New("invalid id format")
	}

	book, err := services.GetOneBook(id)
	if err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return errors.New("book not found")
	}

	json.NewEncoder(w).Encode(book)
	return nil
}

func deleteBook(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return errors.New("id parameter required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return errors.New("invalid id format")
	}
	deletedBook, err := services.DeleteBook(id)
	if err != nil {
		json.NewEncoder(w).Encode(models.Response{
			Status:  "failed",
			Message: "book not found on database",
		})
		return errors.New("book not found")
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "book deleted successfully",
		Book:    &deletedBook,
	})
	return nil

}

func newBook(w http.ResponseWriter, r *http.Request) error {
	var newBook models.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		log.Printf("POST failed: invalid JSON from %s", r.RemoteAddr)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return errors.New("invalid JSON request")
	}
	createdBook, err := services.AddNewBook(newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return errors.New("failed to create book")
	}
	json.NewEncoder(w).Encode(models.Response{
		Status:  "status",
		Message: "new book added successfully",
		Book:    &createdBook,
	})
	return nil
}

func updateBook(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return errors.New("id parameter required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return errors.New("invalid id format")
	}

	var updates models.Book
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Printf("PATCH failed: invalid JSON from %s", r.RemoteAddr)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return errors.New("invalid JSON request")
	}

	updatedBook, err := services.UpdateBook(id, updates)
	if err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return errors.New("book not found")
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "book updated successfully",
		Book:    &updatedBook,
	})
	return nil
}
