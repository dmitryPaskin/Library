package models

type User struct {
	ID          int
	Name        string
	RentedBooks []Book
}

type Book struct {
	ID     int
	Title  string
	Author Authors
	Rented bool
}

type Authors struct {
	Id    int
	Name  string
	Books []Book
}
