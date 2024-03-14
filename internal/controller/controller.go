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

func (h *LibraryHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.sl.GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}

func (h *LibraryHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.sl.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

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
