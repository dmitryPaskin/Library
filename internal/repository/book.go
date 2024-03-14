package repository

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"studentgit.kata.academy/xp/Library/internal/models"
)

type Book struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
}

func NewBook(db *sql.DB) Book {
	return Book{
		db,
		sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Book) GetBookByID(bookID int) (*models.Book, error) {
	var book models.Book
	query := r.sqlBuilder.Select("id", "title", "author_id").
		From("books").
		Where(sq.Eq{"id": bookID})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&book.ID, &book.Title, &book.Author.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

func (r *Book) AddBook(title string, authorID int) error {
	var authorExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM Authors WHERE ID = $1)", authorID).Scan(&authorExists)
	if err != nil {
		return err
	}
	if !authorExists {
		return errors.New("указанного автора нет в списке")
	}

	// Если автор существует, вставляем новую книгу
	_, err = r.db.Exec("INSERT INTO Books (Title, AuthorID) VALUES ($1, $2)", title, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Book) GetAllBooks() ([]models.Book, error) {
	query := r.sqlBuilder.Select("b.ID AS BookID", "b.Title AS BookTitle", "a.ID AS AuthorID", "a.Name AS AuthorName").
		From("Books b").
		Join("Authors a ON b.AuthorID = a.ID")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author.Id, &book.Author.Name)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
