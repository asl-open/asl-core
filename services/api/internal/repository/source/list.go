package source

import (
	"context"
	"fmt"
)

func (r *repo) List(ctx context.Context) ([]Source, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, title, author, type, edition, language, locator_scheme, created_at, updated_at
		FROM sources
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var s Source
		err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Author,
			&s.Type,
			&s.Edition,
			&s.Language,
			&s.LocatorScheme,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan source: %w", err)
		}
		sources = append(sources, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}

	return sources, nil
}
