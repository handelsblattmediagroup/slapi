package core

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type PrometheusServer struct {
	server *http.Server
}

func NewPrometheusServer(lc fx.Lifecycle, shutdown fx.Shutdowner, c *Config) {
	server := &http.Server{
		Handler: promhttp.Handler(),
		Addr:    c.ListenAddrPrometheus,

		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Minute,
		IdleTimeout:       time.Minute,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				listener, err := net.Listen("tcp", c.ListenAddrPrometheus)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error().Err(err).Send()
					_ = shutdown.Shutdown()
				}

				addr := listener.Addr().(*net.TCPAddr)

				log.Info().Int("port", addr.Port).Str("ip", addr.IP.String()).Msg("listening (prometheus)")

				err = server.Serve(listener)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error().Err(err).Send()
					_ = shutdown.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
