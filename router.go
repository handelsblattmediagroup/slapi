package slapi

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"serenitylabs.cloud/slapi/pkg/ginutil"

	"github.com/rs/zerolog/log"
	"serenitylabs.cloud/slapi/pkg/api"

	"github.com/gin-contrib/cors"
	ginlogger "github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Core struct {
	*gin.Engine
	log zerolog.Logger
}

func New(routers api.CoreIn) *Core {
	gin.SetMode(gin.ReleaseMode)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	core := &Core{
		Engine: gin.New(),
		log:    log.With().Logger(),
	}

	logger := ginlogger.WithLogger(func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {

		l := core.log.With().
			Str("id", requestid.Get(c)).
			Str("path", c.Request.URL.Path).
			Dur("latency", latency).
			Str("remote_addr", c.Request.RemoteAddr)

		if value, ok := c.Get("slapi.router_version"); ok {
			l = l.Uint("router_version", value.(uint))
		}

		if value, ok := c.Get("slapi.router_name"); ok {
			l = l.Str("router_name", value.(string))
		}

		return l.Logger()
	})

	core.Use(requestid.New())
	core.Use(ginlogger.SetLogger(logger))
	core.Use(ginutil.ErrorHandler())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, "https://acuteaura.net")
	core.Use(cors.New(corsConfig))

	core.RemoteIPHeaders = []string{"Fly-Client-IP"}

	core.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			Version    string `json:"version"`
			Allocation string `json:"allocation"`
		}{"0.0.2", os.Getenv("FLY_ALLOC_ID")})
	})

	for _, router := range routers.VersionedRouter {
		prefix := fmt.Sprintf("/v%d/%s/", router.Version, router.Prefix)
		core.log.Info().Str("prefix", prefix).Msg("registering router prefix")
		group := core.Group(prefix)
		group.Use(ginutil.AnnotateRouter(router))
		router.RegisterRouter(group)
	}

	return core
}
