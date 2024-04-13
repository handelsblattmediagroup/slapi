package core

import (
	"github.com/gin-gonic/gin"
)

func AnnotateRouter(vrs *VersionedRouterSpec) gin.HandlerFunc {
	version, prefix := vrs.Version, vrs.Prefix
	return func(c *gin.Context) {
		c.Set(GIN_KV_ROUTER_VERSION_KEY, version)
		c.Set(GIN_KV_ROUTER_NAME_KEY, prefix)
	}
}
