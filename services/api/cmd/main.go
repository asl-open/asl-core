package main

import (
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/services/api/internal"
)

func main() {
	fx.New(
		internal.Module,
		fx.RecoverFromPanics(),
	).Run()
}
