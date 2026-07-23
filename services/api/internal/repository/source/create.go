package source

import (
	"context"
	"fmt"
)

func (r *repo) Create(ctx context.Context, s *Source) (Source, error) {
	row := r.conn.QueryRow(ctx, `
		INSERT INTO sources (title, author, type, edition, language, locator_scheme)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, author, type, edition, language, locator_scheme, created_at, updated_at
	`, s.Title, s.Author, s.Type, s.Edition, s.Language, s.LocatorScheme)

	var created Source
	err := row.Scan(
		&created.ID,
		&created.Title,
		&created.Author,
		&created.Type,
		&created.Edition,
		&created.Language,
		&created.LocatorScheme,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return Source{}, fmt.Errorf("failed to create source: %w", err)
	}

	return created, nil
}
