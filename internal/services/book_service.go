package services

import (
	"github.com/4Noyis/my-library/internal/models"
	"github.com/4Noyis/my-library/internal/repositories"
)

type BookService struct {
	bookRepo *repositories.BookRepository
}

func NewBookService() *BookService {
	return &BookService{
		bookRepo: repositories.NewBookCollection(),
	}
}

func (bs *BookService) GetAllBooks() ([]models.Book, error) {
	return bs.bookRepo.GetAllBooks()
}

func (bs *BookService) GetOneBook(id int) (models.Book, error) {
	return bs.bookRepo.GetOneBook(id)
}

func (bs *BookService) DeleteBook(id int) (models.Book, error) {
	return bs.bookRepo.DeleteBook(id)
}

func (bs *BookService) AddNewBook(book models.Book) (models.Book, error) {
	return bs.bookRepo.AddNewBook(book)
}

func (bs *BookService) UpdateBook(id int, updates models.Book) (models.Book, error) {
	return bs.bookRepo.UpdateBook(id, updates)
}
