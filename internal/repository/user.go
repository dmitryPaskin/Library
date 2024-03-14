package repository

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"studentgit.kata.academy/xp/Library/internal/models"
)

type User struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
}

func NewUser(db *sql.DB) User {
	return User{
		db:         db,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *User) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	query := r.sqlBuilder.Select("id", "name").
		From("users").
		Where(sq.Eq{"id": userID})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&user.ID, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *User) GetUsersWithRentedBooks() ([]models.User, error) {
	query := sq.Select("u.ID", "u.Name", "b.ID AS BookID", "b.Title", "a.ID AS AuthorID", "a.Name AS AuthorName").
		From("Users u").
		LeftJoin("Rentals r ON u.ID = r.UserID").
		LeftJoin("Books b ON r.BookID = b.ID").
		LeftJoin("Authors a ON b.AuthorID = a.ID")

	sqlQuery, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Мапим результаты запроса в список пользователей с их арендованными книгами
	usersMap := make(map[int]*models.User)
	for rows.Next() {
		var user models.User
		var book models.Book
		var author models.Authors
		err = rows.Scan(&user.ID, &user.Name, &book.ID, &book.Title, &author.Id, &author.Name)
		if err != nil {
			return nil, err
		}
		// Если пользователь встречается впервые, добавляем его в map
		if _, ok := usersMap[user.ID]; !ok {
			usersMap[user.ID] = &user
		}
		// Добавляем книгу в список арендованных книг пользователя
		usersMap[user.ID].RentedBooks = append(usersMap[user.ID].RentedBooks, models.Book{
			ID:     book.ID,
			Title:  book.Title,
			Author: models.Authors{Id: author.Id, Name: author.Name},
		})
	}

	// Преобразуем map в список пользователей
	var users []models.User
	for _, user := range usersMap {
		users = append(users, *user)
	}

	return users, nil
}
