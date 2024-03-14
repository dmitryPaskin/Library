package repository

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type Rental struct {
	db         *sql.DB
	sqlBuilder sq.StatementBuilderType
}

func NewRental(db *sql.DB) Rental {
	return Rental{
		db,
		sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Rental) ReturnBook(bookID, userID int) error {
	// Проверяем, арендована ли книга пользователем
	isRented, err := r.IsBookRentedByUser(bookID, userID)
	if err != nil {
		return err
	}
	if !isRented {
		return errors.New("book is not rented by user")
	}

	// Обновляем запись об аренде, устанавливая дату возврата
	updateQuery := r.sqlBuilder.Update("Rentals").
		Set("ReturnDate", time.Now()).
		Where(sq.Eq{"BookID": bookID, "UserID": userID})

	sqlQuery, args, err := updateQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rental) IsBookRentedByUser(bookID, userID int) (bool, error) {
	var count int
	query := r.sqlBuilder.Select("COUNT(*)").
		From("Rentals").
		Where(sq.Eq{"BookID": bookID, "UserID": userID, "ReturnDate": nil})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Rental) RentBook(bookID, userID int) error {
	isRented, err := r.IsBookRented(bookID)
	if err != nil {
		return err
	}
	if isRented {
		return errors.New("book already rented")
	}

	// Создаем запись об аренде книги пользователем
	insertQuery := r.sqlBuilder.Insert("Rentals").
		Columns("UserID", "BookID", "RentDate").
		Values(userID, bookID, time.Now())

	sqlQuery, args, err := insertQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rental) IsBookRented(bookID int) (bool, error) {
	var rented bool
	query := r.sqlBuilder.Select("COUNT(*) > 0").
		From("Rentals").
		Where(sq.Eq{"BookID": bookID, "ReturnDate": nil})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	err = r.db.QueryRow(sqlQuery, args...).Scan(&rented)
	if err != nil {
		return false, err
	}

	return rented, nil
}
