package store

import (
	"errors"
	"sync"

	"github.com/k20ku/see/entity"
)

var (
	Items       = &ItemStore{Items: map[entity.ItemId]*entity.Item{}}
	ErrNotFound = errors.New("not found")
)

type ItemStore struct {
	mu     sync.RWMutex
	LastId entity.ItemId
	Items  map[entity.ItemId]*entity.Item
}

func (is *ItemStore) Add(i *entity.Item) (entity.ItemId, error) {
	is.mu.Lock()
	defer is.mu.Unlock()
	is.LastId++
	i.Id = is.LastId
	is.Items[i.Id] = i
	return i.Id, nil
}

// Returns all items, which are sorted by Id
func (is *ItemStore) All() entity.Items {
	items := make([]*entity.Item, len(is.Items))
	for i, item := range items {
		items[i-1] = item
	}
	return items
}
