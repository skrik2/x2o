package middleware

import (
	"net/http"

	"github.com/skrik2/x2o/app/objects"

	"github.com/gin-gonic/gin"
)

// AbortWithError aborts the request with a JSON error response and adds the error to gin context for access logging.
func AbortWithError(c *gin.Context, status int, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(status, objects.ErrorResponse{
		Error: objects.Error{
			Type:    http.StatusText(status),
			Message: err.Error(),
		},
	})
}
