package source

import (
	"context"

	"github.com/google/uuid"

	sourcerepo "github.com/asl-open/asl-core/services/api/internal/repository/source"
)

func (s *service) Get(ctx context.Context, id uuid.UUID) (sourcerepo.Source, error) {
	return s.repo.Get(ctx, id)
}
