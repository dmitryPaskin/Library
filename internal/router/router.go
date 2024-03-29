package router

import (
	"database/sql"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"studentgit.kata.academy/xp/Library/internal/controller"
	_ "studentgit.kata.academy/xp/Library/internal/docs"
	"studentgit.kata.academy/xp/Library/internal/repository"
	"studentgit.kata.academy/xp/Library/internal/repository/generator"
	"studentgit.kata.academy/xp/Library/internal/service"
)

// @title Library
// @version 1.0
// @description This is implementation library API
// @host localhost:8080
// @BasePath /

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
	generator.GenerateTable(generator.New(r.db))

	serviceRout := service.LibraryService{
		UserRepo:   repository.NewUser(r.db),
		BookRepo:   repository.NewBook(r.db),
		RentalRepo: repository.NewRental(r.db),
		AuthorRepo: repository.NewAuthor(r.db),
	}

	handler := controller.NewController(serviceRout)
	r.r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
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
