package contributor

import (
	"context"

	"github.com/google/uuid"

	contributorrepo "github.com/asl-open/asl-core/services/api/internal/repository/contributor"
)

func (s *service) Get(ctx context.Context, id uuid.UUID) (contributorrepo.Contributor, error) {
	return s.repo.Get(ctx, id)
}
