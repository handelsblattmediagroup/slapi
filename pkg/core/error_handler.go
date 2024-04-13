package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		c.JSON(http.StatusInternalServerError, struct {
			Err string `json:"err"`
		}{Err: err.Error()})
	}
}
