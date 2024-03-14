package repository

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"studentgit.kata.academy/xp/Library/internal/models"
)

type Author struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
}

func NewAuthor(db *sql.DB) Author {
	return Author{
		db:         db,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Author) AddAuthor(name string) (int, error) {

	insertQuery := r.sqlBuilder.
		Insert("Authors").
		Columns("Name").
		Values(name).
		Suffix("RETURNING ID").
		RunWith(r.db)

	var authorID int
	err := insertQuery.QueryRow().Scan(&authorID)
	if err != nil {
		return 0, err
	}

	return authorID, nil
}

func (r *Author) GetAuthorsWithBooks() ([]models.Authors, error) {
	query := r.sqlBuilder.
		Select("a.ID AS AuthorID", "a.Name AS AuthorName", "b.ID AS BookID", "b.Title AS BookTitle").
		From("Authors a").
		LeftJoin("Books b ON a.ID = b.AuthorID").
		OrderBy("a.ID, b.ID").
		RunWith(r.db)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []models.Authors
	var currentAuthorID int
	var currentAuthor models.Authors

	for rows.Next() {
		var authorID int
		var authorName string
		var bookID sql.NullInt64
		var bookTitle sql.NullString

		if err := rows.Scan(&authorID, &authorName, &bookID, &bookTitle); err != nil {
			return nil, err
		}

		if authorID != currentAuthorID {
			if currentAuthorID != 0 {
				authors = append(authors, currentAuthor)
			}
			currentAuthorID = authorID
			currentAuthor = models.Authors{
				Id:   authorID,
				Name: authorName,
			}
		}

		if bookID.Valid && bookTitle.Valid {
			currentAuthor.Books = append(currentAuthor.Books, models.Book{
				ID:    int(bookID.Int64),
				Title: bookTitle.String,
			})
		}
	}
	if currentAuthorID != 0 {
		authors = append(authors, currentAuthor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *Author) GetTopAuthors(limit int) ([]models.Authors, error) {
	query := sq.
		Select("a.ID", "a.Name", "COUNT(b.ID) AS BookCount").
		From("Authors a").
		LeftJoin("Books b ON a.ID = b.AuthorID").
		LeftJoin("Rentals r ON b.ID = r.BookID").
		GroupBy("a.ID", "a.Name").
		OrderBy("BookCount DESC").
		Limit(uint64(limit)).
		RunWith(r.db)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topAuthors []models.Authors
	for rows.Next() {
		var author models.Authors
		err = rows.Scan(&author.Id, &author.Name)
		if err != nil {
			return nil, err
		}
		topAuthors = append(topAuthors, author)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topAuthors, nil
}
