package health

import (
	"context"

	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/database"
)

var Module = fx.Provide(New)

// Checker reports whether the service's dependencies are available. As
// more dependencies are added (cache, external APIs, ...), they get
// checked here too - this is the one place that knows what "ready" means.
type Checker interface {
	Ready(ctx context.Context) error
}

type Params struct {
	fx.In

	DB database.Conn
}

type checker struct {
	db database.Conn
}

func New(p Params) Checker {
	return &checker{db: p.DB}
}

func (c *checker) Ready(ctx context.Context) error {
	return c.db.Ping(ctx)
}
