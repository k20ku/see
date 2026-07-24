package handler

import (
	"log"
	"net/http"

	"github.com/k20ku/see/entity"
	"github.com/k20ku/see/store"
)

type ListItem struct {
	Store *store.ItemStore
}

type item struct {
	Id    entity.ItemId `json:"id"`
	Title string        `json:"title"`
	Url   string        `json:"url"`
}

func (li *ListItem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items := li.Store.All()
	rsp := []item{}
	for _, itm := range items {
		rsp = append(rsp, item{
			Id:    itm.Id,
			Title: itm.Title,
			Url:   itm.Url,
		})
	}
	if err := RespondJSON(ctx, w, http.StatusOK, rsp); err != nil {
		log.Printf("list item failed to respond: %v", err)
	}
}
