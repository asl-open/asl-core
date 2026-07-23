package contributor

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/asl-open/asl-core/services/api/internal/http/apierrors"
)

func (r *repo) Get(ctx context.Context, id uuid.UUID) (Contributor, error) {
	row := r.conn.QueryRow(ctx, `
		SELECT id, display_name, roles, bio, handle, created_at, updated_at
		FROM contributors
		WHERE id = $1
	`, id)

	var c Contributor
	err := row.Scan(
		&c.ID,
		&c.DisplayName,
		&c.Roles,
		&c.Bio,
		&c.Handle,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Contributor{}, apierrors.ErrNotFound
		}
		return Contributor{}, fmt.Errorf("failed to get contributor: %w", err)
	}

	return c, nil
}
