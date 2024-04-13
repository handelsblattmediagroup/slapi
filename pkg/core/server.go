package core

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"net"
	"net/http"
	"time"
)

func SetupServer(lc fx.Lifecycle, shutdown fx.Shutdowner, router *Core, config *Config) *http.Server {
	server := &http.Server{
		Handler: router,
		Addr:    config.ListenAddr,

		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Minute,
		IdleTimeout:       time.Minute,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				listener, err := net.Listen("tcp", config.ListenAddr)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error().Err(err).Send()
					_ = shutdown.Shutdown()
				}

				addr := listener.Addr().(*net.TCPAddr)

				log.Info().Int("port", addr.Port).Str("ip", addr.IP.String()).Msg("listening")

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

	return server
}
