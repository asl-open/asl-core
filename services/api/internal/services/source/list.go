package source

import (
	"context"

	sourcerepo "github.com/asl-open/asl-core/services/api/internal/repository/source"
)

func (s *service) List(ctx context.Context) ([]sourcerepo.Source, error) {
	return s.repo.List(ctx)
}
