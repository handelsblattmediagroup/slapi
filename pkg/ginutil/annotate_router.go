package ginutil

import (
	"github.com/gin-gonic/gin"
	"serenitylabs.cloud/slapi/pkg/api"
)

func AnnotateRouter(vrs *api.VersionedRouterSpec) gin.HandlerFunc {
	version, prefix := vrs.Version, vrs.Prefix
	return func(c *gin.Context) {
		c.Set("slapi.router_version", version)
		c.Set("slapi.router_name", prefix)
	}
}
