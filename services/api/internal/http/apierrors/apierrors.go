// Package apierrors holds canonical *errorsx.AppError values for the
// api service's handlers to return. Customize per call site with the
// With... builders (e.g. apierrors.ErrNotFound.WithMessage("user 42 not
// found")) - they return a copy, so the shared vars below are never
// mutated.
package apierrors

import (
	"net/http"

	"github.com/asl-open/asl-core/pkg/errorsx"
)

var (
	// ErrBadRequest reports invalid client input, including validation
	// failures.
	ErrBadRequest = errorsx.New().
			WithStatusCode(http.StatusBadRequest).
			WithCode(http.StatusBadRequest).
			WithMessage("Bad Request")

	// ErrNotFound reports a missing resource.
	ErrNotFound = errorsx.New().
			WithStatusCode(http.StatusNotFound).
			WithCode(http.StatusNotFound).
			WithMessage("Not Found")
)
