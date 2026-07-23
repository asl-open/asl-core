package contributor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/database"
)

var Module = fx.Provide(New)

type Role string

const (
	RoleContributor Role = "contributor"
	RoleTranslator  Role = "translator"
	RoleReviewer    Role = "reviewer"
	RoleEditor      Role = "editor"
)

func (r Role) Valid() bool {
	switch r {
	case RoleContributor, RoleTranslator, RoleReviewer, RoleEditor:
		return true
	default:
		return false
	}
}

type Contributor struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Bio         *string
	Handle      *string
	DisplayName string
	Roles       []Role
	ID          uuid.UUID
}

type Repo interface {
	Create(ctx context.Context, c *Contributor) (Contributor, error)
	Get(ctx context.Context, id uuid.UUID) (Contributor, error)
	List(ctx context.Context) ([]Contributor, error)
}

type Params struct {
	fx.In

	Conn database.Conn
}

type repo struct {
	conn database.Conn
}

func New(params Params) Repo {
	return &repo{conn: params.Conn}
}
