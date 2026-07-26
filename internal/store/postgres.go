package store

import (
	"context"
	"fmt"

	"github.com/k20ku/see/internal/db"
	"github.com/k20ku/see/internal/entity"
)

// Postgres stores items in PostgreSQL through sqlc-generated queries.
type Postgres struct {
	q *db.Queries
}

func NewPostgres(q *db.Queries) *Postgres {
	return &Postgres{q: q}
}

func (p *Postgres) Add(ctx context.Context, i *entity.Item) (*entity.Item, error) {
	row, err := p.q.CreateItem(ctx, db.CreateItemParams{
		Title: i.Title,
		Url:   i.Url,
		Note:  i.Note,
	})
	if err != nil {
		return nil, fmt.Errorf("store add item: %w", err)
	}
	return toEntity(row), nil
}

func (p *Postgres) AddItems(ctx context.Context, items entity.Items) (int64, error) {
	rows := make([]db.CreateItemsParams, 0, len(items))
	for _, itm := range items {
		rows = append(rows, db.CreateItemsParams{
			Title: itm.Title,
			Url:   itm.Url,
			Note:  itm.Note,
		})
	}
	n, err := p.q.CreateItems(ctx, rows)
	if err != nil {
		return -1, fmt.Errorf("store add items: %w", err)
	}
	return n, nil
}

func (p *Postgres) ListItems(ctx context.Context) (entity.Items, error) {
	rows, err := p.q.ListItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("store list items: %w", err)
	}
	items := make(entity.Items, 0, len(rows))
	for _, row := range rows {
		items = append(items, toEntity(row))
	}
	return items, nil
}

// toEntity converts a persistence row into the domain type, keeping pgtype
// out of every other package.
func toEntity(row db.Item) *entity.Item {
	return &entity.Item{
		Id:         entity.ItemId(row.ID),
		Title:      row.Title,
		Url:        row.Url,
		Note:       row.Note,
		CreatedAt:  row.CreatedAt.Time,
		ModifiedAt: row.UpdatedAt.Time,
	}
}
