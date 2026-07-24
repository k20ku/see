package server

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/k20ku/see/handler"
	"github.com/k20ku/see/store"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status" : "ok"}`))
	})
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
	ai := &handler.AddItem{Store: store.Items, Validate: v}
	mux.Post("/items", ai.ServeHTTP)
	li := &handler.ListItem{Store: store.Items}
	mux.Get("/items", li.ServeHTTP)
	return mux
}
