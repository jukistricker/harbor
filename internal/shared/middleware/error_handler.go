package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"harbor/internal/shared/catalog"
	appErr "harbor/internal/shared/errors"
	"harbor/internal/shared/response"
)

func ErrorHandler(debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var ae *appErr.AppError
		if errors.As(err, &ae) {
			var data any

			if debug && ae.Err != nil {
				data = ae.Err.Error()
			} else {
				data = nil
			}

			response.Send(
				c,
				ae.Status,
				ae.Message,
				data,
			)

			c.Abort()
			return
		}

		// Unknown error -> Internal
		var data any
		if debug {
			data = err.Error()
		} else {
			data = nil
		}

		response.Send(
			c,
			http.StatusInternalServerError,
			catalog.Internal.Message,
			data,
		)

		c.Abort()
	}
}
