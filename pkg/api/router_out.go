package api

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type VersionedRouterSpec struct {
	Prefix         string
	Version        uint
	RegisterRouter func(group *gin.RouterGroup)
}

type CoreIn struct {
	fx.In

	VersionedRouter []*VersionedRouterSpec `group:"router"`
	OtelTracer      trace.TracerProvider
}

type VersionedRouterOut struct {
	fx.Out

	Router *VersionedRouterSpec `group:"router"`
}
