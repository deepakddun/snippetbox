package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Updated time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

// func will insert a new snippet into the database
func (m *SnippetModel) Insert(ctx context.Context, title string, content string, expires int) (int, error) {

	stmt := `
	INSERT INTO snippets (title, content, created, updated, expires)
	VALUES ($1, $2, NOW(), NOW(), NOW() + $3 * INTERVAL '1 day')
	RETURNING id
`
	var id int

	err := m.DB.QueryRow(ctx, stmt, title, content, expires).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// this will return a specific snippet based on its Ids
func (m *SnippetModel) Get(ctx context.Context, id int) (Snippet, error) {
	stmt := `
	SELECT id , title, content, created, updated, expires
	from snippets where expires > NOW() and id = $1
`
	var s Snippet
	err := m.DB.QueryRow(ctx, stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Updated, &s.Expires)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrorNoRecord
		} else {
			return Snippet{}, err
		}

	}
	return s, nil
}

func (m *SnippetModel) Latest(ctx context.Context) ([]Snippet, error) {
	stmt := `
	SELECT id , title, content, created, updated, expires
	from snippets where expires > NOW() ORDER BY id DESC LIMIT 10
	`

	rows, err := m.DB.Query(ctx, stmt)
	//defer rows.Close()

	if err != nil {
		return nil, err
	}

	snippets, err := pgx.CollectRows(rows, pgx.RowToStructByName[Snippet])

	if err != nil {
		return nil, err
	}

	return snippets, nil

}
