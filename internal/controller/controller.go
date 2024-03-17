package controller

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"studentgit.kata.academy/xp/Library/internal/service"
)

type LibraryHandler struct {
	sl service.LibraryService
}

func NewController(sl service.LibraryService) LibraryHandler {
	return LibraryHandler{
		sl: sl,
	}
}

// @Summary Rent to book by user ID and book ID
// @Tags library
// @Accept json
// @Produce json
// @Param userID path string true "user ID"
// @Param bookID path string true "book ID"
// @Success 200
// @Router /users/{userID}/books/{bookID}/rent [post]
func (h *LibraryHandler) RentBook(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.sl.RentBook(userID, bookID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Return to book by user ID and book ID
// @Tags library
// @Accept json
// @Produce json
// @Param userID path string true "user ID"
// @Param bookID path string true "book ID"
// @Success 200
// @Router /users/{userID}/books/{bookID}/return [post]
func (h *LibraryHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.sl.ReturnBook(bookID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Add author by name
// @Tags library
// @Accept json
// @Produce json
// @Param author path string true "name author"
// @Success 200
// @Router /authors/{name} [post]
func (h *LibraryHandler) AddAuthor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sl.AddAuthor(req.Name)

	w.WriteHeader(http.StatusOK)
}

// @Summary Get all authours
// @Tags library
// @Produce json
// @Success 200 {object} models.Authors
// @Router /authors [get]
func (h *LibraryHandler) GetAuthorsWithBooks(w http.ResponseWriter, r *http.Request) {
	authors, err := h.sl.GetAuthorsWithBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(authors); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Add book
// @Tags library
// @Accept json
// @Produce json
// @Param book path string true "Title"
// @Param book path int true "Id Author"
// @Success 200
// @Router /book/{name} [post]
func (h *LibraryHandler) AddBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		AuthorID int    `json:"author_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.sl.AddBook(req.Title, req.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get all books
// @Tags library
// @Produce json
// @Success 200 {object} []models.Book
// @Router /books [get]
func (h *LibraryHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.sl.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}

// @Summary Get all users
// @Tags library
// @Produce json
// @Success 200 {object} []models.User
// @Router /users [get]
func (h *LibraryHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.sl.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// @Summary Get toop authors by limit
// @Tags library
// @Accept json
// @Produce json
// @Param author path string true "limit"
// @Success 200 {object} []models.Authors
// @Router /authors/top/{limit} [post]
func (h *LibraryHandler) GetTopAuthors(w http.ResponseWriter, r *http.Request) {
	limitStr := chi.URLParam(r, "limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	authors, err := h.sl.GetTopAuthors(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(authors)
}
