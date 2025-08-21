package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/services"
)

func BookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		if r.URL.Query().Get("id") != "" {
			GetOneBook(w, r)
		} else {
			getAllBooks(w)
		}
	} else if r.Method == "POST" {

	} else if r.Method == "PATCH" {

	} else if r.Method == "DELETE" {
		DeleteBook(w, r)
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
func GetOneBook(w http.ResponseWriter, r *http.Request) error {
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

func DeleteBook(w http.ResponseWriter, r *http.Request) error {
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
		json.NewEncoder(w).Encode(models.DeleteResponse{
			Status:  "failed",
			Message: "book not found on database",
		})
		return errors.New("book not found")
	}

	json.NewEncoder(w).Encode(models.DeleteResponse{
		Status:  "success",
		Message: "book deleted successfully",
		Book:    deletedBook,
	})
	return nil

}

func NewBook(w http.ResponseWriter, r *http.Request) {

}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

}
