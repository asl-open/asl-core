package contributor

import (
	"context"

	contributorrepo "github.com/asl-open/asl-core/services/api/internal/repository/contributor"
)

func (s *service) List(ctx context.Context) ([]contributorrepo.Contributor, error) {
	return s.repo.List(ctx)
}
