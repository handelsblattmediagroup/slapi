package core

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"net/http"
	"os"
	"time"

	"go.uber.org/fx"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog/log"

	ginlogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type VersionedRouterSpec struct {
	Prefix         string
	Version        uint
	RegisterRouter func(group *gin.RouterGroup)
}

type VersionedRouterOut struct {
	fx.Out

	Router *VersionedRouterSpec `group:"router"`
}

type Core struct {
	*gin.Engine
	log zerolog.Logger
}

type In struct {
	fx.In

	Config          *Config
	VersionedRouter []*VersionedRouterSpec `group:"router"`
	CorsConfig      *cors.Config           `optional:"true"`
	OtelTracer      trace.TracerProvider
}

func New(in In) (*Core, error) {
	_, inFly := os.LookupEnv("FLY_ALLOC_ID")

	logLevel, err := zerolog.ParseLevel(in.Config.LogLevel)
	if err != nil {
		panic("bad log level")
	}
	log.Logger = log.Level(logLevel)

	gin.SetMode(gin.ReleaseMode)
	core := &Core{
		Engine: gin.New(),
		log:    log.With().Logger(),
	}

	logger := ginlogger.WithLogger(func(c *gin.Context, log zerolog.Logger) zerolog.Logger {
		start := time.Now().UTC()
		span := trace.SpanFromContext(c.Request.Context())

		defer c.Next()

		end := time.Now()
		latency := end.Sub(start)

		l := core.log.With().
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("path", c.Request.URL.Path).
			Dur("latency", latency).
			Str("remote_addr", c.Request.RemoteAddr)

		if value, ok := c.Get(GIN_KV_ROUTER_VERSION_KEY); ok {
			l = l.Uint("router_version", value.(uint))
		}

		if value, ok := c.Get(GIN_KV_ROUTER_NAME_KEY); ok {
			l = l.Str("router_name", value.(string))
		}

		return l.Logger()
	})

	core.Use(otelgin.Middleware("slapi"))
	core.Use(ginlogger.SetLogger(logger))
	core.Use(ErrorHandler())

	if inFly {
		core.RemoteIPHeaders = []string{"Fly-Client-IP"}
	} else {
		err := core.SetTrustedProxies(nil)
		if err != nil {
			return nil, err
		}
	}

	core.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			Version string `json:"version"`
		}{"0.0.2"})
	})

	for _, router := range in.VersionedRouter {
		prefix := fmt.Sprintf("/v%d/%s/", router.Version, router.Prefix)
		core.log.Info().Str("prefix", prefix).Msg("registering router prefix")
		group := core.Group(prefix)
		group.Use(AnnotateRouter(router))
		router.RegisterRouter(group)
	}

	return core, nil
}
