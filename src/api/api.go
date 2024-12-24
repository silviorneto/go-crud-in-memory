package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func NewHandler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(jsonMiddleware)

	repo := NewInMemoryStore()

	r.Route("/api", func(r chi.Router) {
		r.Get("/users", FindAll(repo))
		r.Get("/users/{id}", FindById(repo))
		r.Post("/users", Insert(repo))
		r.Delete("/users/{id}", Delete(repo))
		r.Put("/users/{id}", Update(repo))
	})

	return r
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func sendJSON(w http.ResponseWriter, response Response, statusCode int) {
	data, err := json.Marshal(response)
	if err != nil {
		slog.Error("Failed to marshal json data")
		sendJSON(w, Response{Error: "Something went wrong"}, http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}

func FindAll(repo *InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := repo.GetUsers()

		sendJSON(w, Response{Data: users}, http.StatusOK)
	}
}

func FindById(repo *InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			slog.Error("Failed to parse id")
			sendJSON(w, Response{Error: "Invalid id"}, http.StatusBadRequest)
			return
		}

		user, err := repo.GetUserById(uuid)
		if err != nil && errors.Is(err, ErrorNotFound) {
			sendJSON(w, Response{Error: err.Error()}, http.StatusNotFound)
			return
		}

		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}

func Insert(repo *InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.MaxBytesReader(w, r.Body, 1000)

		var body User
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			slog.Error("Failed to decode body")
			sendJSON(w, Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		repo.CreateUser(body)

		sendJSON(w, Response{}, http.StatusCreated)
	}
}

func Delete(repo *InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			slog.Error("Failed to parse id")
			sendJSON(w, Response{Error: "Invalid id"}, http.StatusBadRequest)
			return
		}

		if err := repo.DeleteUser(uuid); err != nil && errors.Is(err, ErrorNotFound) {
			sendJSON(w, Response{Error: err.Error()}, http.StatusNotFound)
		}

		sendJSON(w, Response{}, http.StatusOK)
	}
}

func Update(repo *InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.MaxBytesReader(w, r.Body, 1000)

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			slog.Error("Failed to parse id")
			sendJSON(w, Response{Error: "Invalid id"}, http.StatusBadRequest)
			return
		}

		var body User
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			slog.Error("Failed to decode body")
			sendJSON(w, Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		repo.UpdateUser(body, uuid)
	}
}
