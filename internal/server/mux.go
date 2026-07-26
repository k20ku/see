package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/k20ku/see/internal/handler"
	"github.com/k20ku/see/internal/store"
	"github.com/k20ku/see/internal/validate"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status" : "ok"}`))
	})
	store := store.New()
	ai := &handler.AddItem{Store: store, Validate: validate.New()}
	mux.Post("/items", ai.ServeHTTP)
	li := &handler.ListItem{Store: store}
	mux.Get("/items", li.ServeHTTP)
	return mux
}
