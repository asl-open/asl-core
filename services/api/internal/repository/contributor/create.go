package contributor

import (
	"context"
	"fmt"
)

func (r *repo) Create(ctx context.Context, c *Contributor) (Contributor, error) {
	row := r.conn.QueryRow(ctx, `
		INSERT INTO contributors (display_name, roles, bio, handle)
		VALUES ($1, $2, $3, $4)
		RETURNING id, display_name, roles, bio, handle, created_at, updated_at
	`, c.DisplayName, c.Roles, c.Bio, c.Handle)

	var created Contributor
	err := row.Scan(
		&created.ID,
		&created.DisplayName,
		&created.Roles,
		&created.Bio,
		&created.Handle,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return Contributor{}, fmt.Errorf("failed to create contributor: %w", err)
	}

	return created, nil
}
