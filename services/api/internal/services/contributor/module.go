package contributor

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/fx"

	contributorrepo "github.com/asl-open/asl-core/services/api/internal/repository/contributor"
)

var Module = fx.Provide(New)

type Service interface {
	Create(ctx context.Context, in *CreateInput) (contributorrepo.Contributor, error)
	Get(ctx context.Context, id uuid.UUID) (contributorrepo.Contributor, error)
	List(ctx context.Context) ([]contributorrepo.Contributor, error)
}

type Params struct {
	fx.In

	Repo contributorrepo.Repo
}

type service struct {
	repo contributorrepo.Repo
}

func New(params Params) Service {
	return &service{repo: params.Repo}
}
