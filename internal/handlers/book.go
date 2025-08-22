package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/4Noyis/my-library/internal/logger"
	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// BooksHandler routes all book-related requests based on HTTP method
// func BooksHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	switch r.Method {
// 	case http.MethodGet:
// 		if r.URL.Query().Get("id") != "" {
// 			GetOneBookHandler(w, r)
// 		} else {
// 			GetAllBooksHandler(w, r)
// 		}
// 	case http.MethodPost:
// 		CreateBookHandler(w, r)
// 	case http.MethodPatch:
// 		UpdateBookHandler(w, r)
// 	case http.MethodDelete:
// 		DeleteBookHandler(w, r)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// Alternative handlers for specific HTTP methods when using separate routes
func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	if r.URL.Query().Get("id") != "" {
		GetOneBookHandler(w, r)
	} else {
		GetAllBooksHandler(w, r)
	}
}

func PostBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	CreateBookHandler(w, r)
}

func PatchBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	UpdateBookHandler(w, r)
}

func DeleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	DeleteBookHandler(w, r)
}

func GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := services.GetAllBooks()
	if err != nil {
		logger.LogError("GetAllBooks", err, logrus.Fields{
			"handler": "GetAllBooksHandler",
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.LogDebug("Retrieved all books", logrus.Fields{
		"handler": "GetAllBooksHandler",
		"count":   len(books),
	})

	json.NewEncoder(w).Encode(books)
}

func GetOneBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		logger.LogError("GetOneBook", nil, logrus.Fields{
			"handler": "GetOneBookHandler",
			"error":   "id parameter required",
		})
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.LogError("GetOneBook", err, logrus.Fields{
			"handler": "GetOneBookHandler",
			"id_str":  idStr,
		})
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	book, err := services.GetOneBook(id)
	if err != nil {
		logger.LogError("GetOneBook", err, logrus.Fields{
			"handler": "GetOneBookHandler",
			"id":      id,
		})
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	logger.LogDebug("Successfully retrieved book", logrus.Fields{
		"handler": "GetOneBookHandler",
		"id":      id,
		"title":   book.Title,
	})

	json.NewEncoder(w).Encode(book)
}

func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		logger.LogError("DeleteBook", nil, logrus.Fields{
			"handler": "DeleteBookHandler",
			"error":   "id parameter required",
		})
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.LogError("DeleteBook", err, logrus.Fields{
			"handler": "DeleteBookHandler",
			"id_str":  idStr,
		})
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	deletedBook, err := services.DeleteBook(id)
	if err != nil {
		logger.LogError("DeleteBook", err, logrus.Fields{
			"handler": "DeleteBookHandler",
			"id":      id,
		})
		json.NewEncoder(w).Encode(models.Response{
			Status:  "failed",
			Message: "book not found on database",
		})
		return
	}

	logger.LogInfo("Book deleted successfully", logrus.Fields{
		"handler": "DeleteBookHandler",
		"id":      id,
		"title":   deletedBook.Title,
	})

	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "book deleted successfully",
		Book:    &deletedBook,
	})
}

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		logger.LogError("CreateBook", err, logrus.Fields{
			"handler":     "CreateBookHandler",
			"remote_addr": r.RemoteAddr,
		})
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	createdBook, err := services.AddNewBook(newBook)
	if err != nil {
		logger.LogError("CreateBook", err, logrus.Fields{
			"handler": "CreateBookHandler",
			"title":   newBook.Title,
			"isbn":    newBook.ISBN,
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.LogInfo("Book created successfully", logrus.Fields{
		"handler": "CreateBookHandler",
		"id":      createdBook.ID,
		"title":   createdBook.Title,
		"isbn":    createdBook.ISBN,
	})

	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "new book added successfully",
		Book:    &createdBook,
	})
}

func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		logger.LogError("UpdateBook", nil, logrus.Fields{
			"handler": "UpdateBookHandler",
			"error":   "id parameter required",
		})
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.LogError("UpdateBook", err, logrus.Fields{
			"handler": "UpdateBookHandler",
			"id_str":  idStr,
		})
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	var updates models.Book
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		logger.LogError("UpdateBook", err, logrus.Fields{
			"handler":     "UpdateBookHandler",
			"id":          id,
			"remote_addr": r.RemoteAddr,
		})
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	updatedBook, err := services.UpdateBook(id, updates)
	if err != nil {
		logger.LogError("UpdateBook", err, logrus.Fields{
			"handler": "UpdateBookHandler",
			"id":      id,
		})
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	logger.LogInfo("Book updated successfully", logrus.Fields{
		"handler": "UpdateBookHandler",
		"id":      id,
		"title":   updatedBook.Title,
	})

	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "book updated successfully",
		Book:    &updatedBook,
	})
}
