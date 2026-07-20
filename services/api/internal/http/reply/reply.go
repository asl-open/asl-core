// Package reply writes a response.Response as an HTTP response. This is
// the one place in the api service allowed to turn a Response into
// bytes on the wire - keeps the JSON encoding/header-setting in one
// spot instead of repeated in every handler.
package reply

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/asl-open/asl-core/pkg/response"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json; charset=utf-8"
)

// JSON writes resp as the response body, using resp.HeaderCode as the
// HTTP status.
func JSON(c *gin.Context, resp response.Response) {
	body, err := json.Marshal(resp)
	if err != nil {
		// resp is always one of our own, JSON-safe types - this is not
		// expected to happen. Fall back to a generic message rather
		// than err.Error(), which could echo whatever made the
		// payload unmarshalable.
		http.Error(c.Writer, "failed to encode response", http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set(headerContentType, contentTypeJSON)
	c.Writer.WriteHeader(resp.HeaderCode)
	_, _ = c.Writer.Write(body)
}
