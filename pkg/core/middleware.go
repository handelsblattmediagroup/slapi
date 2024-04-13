package core

import "github.com/gin-gonic/gin"

type Middleware interface {
	Middleware(c *gin.Context)
}
