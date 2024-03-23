package generator

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v6"
	"log"
	"studentgit.kata.academy/xp/Library/internal/models"
)

type Generator struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
}

func New(db *sql.DB) Generator {
	return Generator{
		db:         db,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (g *Generator) tableIsEmpty(tableName string) (bool, error) {
	var count int
	if err := g.db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count); err != nil {
		return false, err
	}
	return count == 0, nil
}

func (g *Generator) generateAuthors() error {
	insertQuery := g.sqlBuilder.Insert("Authors").Columns("Name")
	for i := 0; i < 100; i++ {
		authorName := gofakeit.Name()
		insertQuery = insertQuery.Values(authorName)
	}

	if _, err := insertQuery.RunWith(g.db).Exec(); err != nil {
		return err
	}
	return nil
}

func (g *Generator) generateBooks() error {
	authors, err := getAuthors(g.db, g.sqlBuilder)
	if err != nil {
		return err
	}

	insertQuery := g.sqlBuilder.Insert("Books").Columns("Title", "AuthorID")
	for _, author := range authors {
		for i := 0; i < 10; i++ {
			bookTitle := gofakeit.Book().Title
			insertQuery = insertQuery.Values(bookTitle, author.Id)
		}
	}
	if _, err := insertQuery.RunWith(g.db).Exec(); err != nil {
		return err
	}
	return nil
}

func getAuthors(db *sql.DB, sb sq.StatementBuilderType) ([]models.Authors, error) {
	rows, err := sb.Select("ID").From("Authors").RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []models.Authors
	for rows.Next() {
		var author models.Authors
		if err = rows.Scan(&author.Id); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return authors, nil
}

func GenerateTable(generator Generator) {
	checkAuthors, err := generator.tableIsEmpty("Authors")
	if err != nil {
		log.Fatal(err)
	}
	if checkAuthors {
		if err = generator.generateAuthors(); err != nil {
			log.Fatal(err)
		}
	}

	checkBooks, err := generator.tableIsEmpty("Books")
	if err != nil {
		log.Fatal(err)
	}

	if checkBooks {
		if err = generator.generateBooks(); err != nil {
			log.Fatal(err)
		}
	}
}
