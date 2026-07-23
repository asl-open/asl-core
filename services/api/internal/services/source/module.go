package source

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/fx"

	sourcerepo "github.com/asl-open/asl-core/services/api/internal/repository/source"
)

var Module = fx.Provide(New)

type Service interface {
	Create(ctx context.Context, in *CreateInput) (sourcerepo.Source, error)
	Get(ctx context.Context, id uuid.UUID) (sourcerepo.Source, error)
	List(ctx context.Context) ([]sourcerepo.Source, error)
}

type Params struct {
	fx.In

	Repo sourcerepo.Repo
}

type service struct {
	repo sourcerepo.Repo
}

func New(params Params) Service {
	return &service{repo: params.Repo}
}
