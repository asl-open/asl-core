// Package errorsx provides AppError, an application error that carries
// the information needed to report it consistently over HTTP: a stable
// code, an HTTP status, and a message safe to show a client.
package errorsx

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is an application error. Builder methods (With...) return a
// new value rather than mutating the receiver, so a canonical error
// (e.g. a shared package-level var) can be customized per call site
// without one caller's change leaking into another's.
type AppError struct {
	message    string
	code       int
	statusCode int
	internal   bool
}

// New returns a baseline AppError defaulting to 200 OK; use the With...
// methods to turn it into an actual error.
func New() *AppError {
	return &AppError{
		code:       http.StatusOK,
		message:    http.StatusText(http.StatusOK),
		statusCode: http.StatusOK,
	}
}

func (e *AppError) Code() int       { return e.code }
func (e *AppError) Message() string { return e.message }
func (e *AppError) StatusCode() int { return e.statusCode }

// Internal reports whether this error represents an unexpected failure
// (as opposed to an expected, well-formed application error like
// validation or not-found) - callers can use it to decide whether the
// error is worth logging at a higher severity.
func (e *AppError) Internal() bool { return e.internal }

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s (internal=%t)", e.code, e.message, e.internal)
}

func (e *AppError) clone() *AppError {
	c := *e
	return &c
}

// WithCode returns a copy of e with its stable error code set.
func (e *AppError) WithCode(code int) *AppError {
	c := e.clone()
	c.code = code
	return c
}

// WithMessage returns a copy of e with its client-facing message set.
func (e *AppError) WithMessage(message string) *AppError {
	c := e.clone()
	c.message = message
	return c
}

// WithStatusCode returns a copy of e with its HTTP status set.
func (e *AppError) WithStatusCode(statusCode int) *AppError {
	c := e.clone()
	c.statusCode = statusCode
	return c
}

// WithInternal returns a copy of e marked as an unexpected/internal
// error.
func (e *AppError) WithInternal() *AppError {
	c := e.clone()
	c.internal = true
	return c
}

// As reports whether err is, or wraps, an *AppError.
func As(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	if appErr == nil {
		return nil, false
	}
	return appErr, ok
}
