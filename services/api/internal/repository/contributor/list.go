package contributor

import (
	"context"
	"fmt"
)

func (r *repo) List(ctx context.Context) ([]Contributor, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, display_name, roles, bio, handle, created_at, updated_at
		FROM contributors
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list contributors: %w", err)
	}
	defer rows.Close()

	var contributors []Contributor
	for rows.Next() {
		var c Contributor
		err := rows.Scan(
			&c.ID,
			&c.DisplayName,
			&c.Roles,
			&c.Bio,
			&c.Handle,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contributor: %w", err)
		}
		contributors = append(contributors, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list contributors: %w", err)
	}

	return contributors, nil
}
