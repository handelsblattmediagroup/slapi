package helloworld

import (
	"github.com/acuteaura/slapi/pkg/core"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Provider() (core.VersionedRouterOut, error) {
	var llog = log.With().Str("system", "github").Logger()

	routerGroupFunc := func(group *gin.RouterGroup) {
		group.GET("/", func(c *gin.Context) {
			llog.Info().Msg("got a hello world!")
			c.JSON(200, map[string]any{
				"hello": "world",
			})
		})
	}
	return core.VersionedRouterOut{
		Router: &core.VersionedRouterSpec{
			Version:        0,
			Prefix:         "helloworld",
			RegisterRouter: routerGroupFunc,
		},
	}, nil
}
