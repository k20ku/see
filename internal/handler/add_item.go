package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/k20ku/see/internal/entity"
	"github.com/k20ku/see/internal/respond"
	"github.com/k20ku/see/internal/store"
)

type AddItem struct {
	Store    *store.ItemStore
	Validate *validator.Validate
}

// this implementing http.HandlerFunc
func (ai *AddItem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title"`
		Url   string `json:"url" validate:"required"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&b); err != nil {
		respond.InvalidJSONError(ctx, w)
		return
	}

	// TODO: validator is exposed
	if err := ai.Validate.Struct(b); err != nil {
		respond.ValidationError(ctx, w, err)
		return
	}

	now := time.Now()
	item := &entity.Item{
		Title:      b.Title,
		Url:        b.Url,
		CreatedAt:  now,
		ModifiedAt: now,
	}

	id, err := ai.Store.Add(item)
	if err != nil {
		respond.InternalServerError(ctx, w)
		return
	}

	resp := struct {
		ID entity.ItemId `json:"id"`
	}{ID: id}
	if err := respond.JSON(ctx, w, http.StatusOK, resp); err != nil {
		log.Printf("add item response failed :%v", err)
	}
}
