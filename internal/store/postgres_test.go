package store

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/k20ku/see/internal/db"
	"github.com/k20ku/see/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestPostgres_ListItems(t *testing.T) {
	// TODO: review test method to connect DB
	t.SkipNow()
	ctx := context.Background()

	port := 5432
	if _, defined := os.LookupEnv("CI"); defined {
		port = 5432
	}

	conn, err := pgx.Connect(
		ctx,
		fmt.Sprintf("postgres://see:seedbpass@localhost:%d/see?sslmode=disable", port),
	)
	require.NoError(t, err, "TestPostgres list items: connection failed")

	t.Cleanup(func() {
		_ = conn.Close(ctx)
	})
	tx, err := conn.Begin(ctx)
	require.NoError(t, err,
		"TestPostgres list items: beginning transaction failed",
	)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
	})
	q := db.New(conn)
	qtx := q.WithTx(tx)
	p := NewPostgres(qtx)
	wants := []*entity.Item{
		{Title: "foo", Url: "https://example.com/foo", Note: "FOO"},
		{Title: "bar", Url: "https://example.com/bar", Note: "BAR"},
		{Title: "baz", Url: "https://example.com/baz", Note: "BAZ"},
	}
	n, err := p.AddItems(ctx, wants)
	require.NoError(t, err, "TestPostgres list items: add items failed")

	require.Equal(t, int64(len(wants)), n, "TestPostgres list item: numbers added equals returned?")
	gots, err := p.ListItems(ctx)
	require.NoError(t, err, "TestPostgres: list items failed")
	for _, got := range gots {
		ok := slices.ContainsFunc[entity.Items](
			wants,
			func(want *entity.Item) bool {
				return got.Title == want.Title &&
					got.Url == want.Url &&
					got.Note == want.Note
			},
		)
		require.Truef(t, ok, "got is %+v", got)
		fmt.Printf("\ngot=%#v\n", got)
	}
}
