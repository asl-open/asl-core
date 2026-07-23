package source

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/asl-open/asl-core/services/api/internal/http/apierrors"
)

func (r *repo) Get(ctx context.Context, id uuid.UUID) (Source, error) {
	row := r.conn.QueryRow(ctx, `
		SELECT id, title, author, type, edition, language, locator_scheme, created_at, updated_at
		FROM sources
		WHERE id = $1
	`, id)

	var s Source
	err := row.Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return Source{}, apierrors.ErrNotFound
		}
		return Source{}, fmt.Errorf("failed to get source: %w", err)
	}

	return s, nil
}
