// Package response defines the JSON envelope HTTP responses use, and
// maps application errors (pkg/errorsx) to it.
package response

import (
	"net/http"

	"github.com/asl-open/asl-core/pkg/errorsx"
)

// Response is the JSON envelope every HTTP response uses.
type Response struct {
	Payload    any    `json:"payload,omitempty"`
	Message    string `json:"message" example:"Not Found"`
	RequestID  string `json:"request_id,omitempty"`
	Code       int    `json:"code" example:"404"`
	HeaderCode int    `json:"-"`
}

// Internal is the generic response for an error that isn't a known
// *errorsx.AppError - its real content is never sent to the client.
func Internal() Response {
	return Response{
		Message:    "Internal Server Error",
		Code:       http.StatusInternalServerError,
		HeaderCode: http.StatusInternalServerError,
	}
}

// Err maps err to a Response. Any error that isn't (or doesn't wrap) an
// *errorsx.AppError is treated as unexpected and reported via Internal.
func Err(err error) Response {
	appErr, ok := errorsx.As(err)
	if !ok {
		return Internal()
	}

	return Response{
		Message:    appErr.Message(),
		Code:       appErr.Code(),
		HeaderCode: appErr.StatusCode(),
	}
}

// WithRequestID returns a copy of r with RequestID set.
func (r Response) WithRequestID(id string) Response {
	r.RequestID = id
	return r
}
