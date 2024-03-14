package router

import (
	"database/sql"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"studentgit.kata.academy/xp/Library/internal/controller"
	"studentgit.kata.academy/xp/Library/internal/repository"
	"studentgit.kata.academy/xp/Library/internal/service"
)

type Router struct {
	db *sql.DB
	r  *chi.Mux
}

func NewRouter(db *sql.DB) Router {
	r := chi.NewRouter()
	return Router{
		db: db,
		r:  r,
	}
}
func (r *Router) StartRouter() {
	serviceRout := service.LibraryService{
		UserRepo:   repository.NewUser(r.db),
		BookRepo:   repository.NewBook(r.db),
		RentalRepo: repository.NewRental(r.db),
		AuthorRepo: repository.NewAuthor(r.db),
	}

	handler := controller.NewController(serviceRout)

	r.r.Post("/users/{userID}/books/{bookID}/rent", handler.RentBook)
	r.r.Post("/users/{userID}/books/{bookID}/return", handler.ReturnBook)
	r.r.Post("/authors/{name}", handler.AddAuthor)
	r.r.Get("/authors", handler.GetAuthorsWithBooks)
	r.r.Get("/books", handler.GetAllBooks)
	r.r.Get("/users", handler.GetUsers)
	r.r.Get("/authors/top/{limit}", handler.GetTopAuthors)

	log.Printf("Start server...")
	http.ListenAndServe(":8080", r.r)
}
