package service

import (
	"errors"
	"studentgit.kata.academy/xp/Library/internal/models"
	"studentgit.kata.academy/xp/Library/internal/repository"
)

type LibraryService struct {
	UserRepo   repository.User
	BookRepo   repository.Book
	AuthorRepo repository.Author
	RentalRepo repository.Rental
}

func (s *LibraryService) RentBook(userID, bookID int) error {
	// Проверка существования пользователя
	_, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Проверка существования книги
	_, err = s.BookRepo.GetBookByID(bookID)
	if err != nil {
		return errors.New("book not found")
	}

	// Если все проверки прошли успешно, арендуем книгу
	return s.RentalRepo.RentBook(bookID, userID)
}

func (s *LibraryService) ReturnBook(bookID, userID int) error {
	err := s.RentalRepo.ReturnBook(bookID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *LibraryService) AddAuthor(name string) {
	s.AuthorRepo.AddAuthor(name)
}

func (s *LibraryService) GetAuthorsWithBooks() ([]models.Authors, error) {
	authors, err := s.AuthorRepo.GetAuthorsWithBooks()
	if err != nil {
		return nil, err
	}
	return authors, err
}

func (s *LibraryService) AddBook(title string, authorId int) error {
	err := s.BookRepo.AddBook(title, authorId)
	if err != nil {
		return err
	}
	return nil
}

func (s *LibraryService) GetAllBooks() ([]models.Book, error) {
	books, err := s.BookRepo.GetAllBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *LibraryService) GetUsers() ([]models.User, error) {
	books, err := s.UserRepo.GetUsersWithRentedBooks()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *LibraryService) GetTopAuthors(limit int) ([]models.Authors, error) {
	return s.AuthorRepo.GetTopAuthors(limit)
}
