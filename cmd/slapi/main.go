package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"

	"github.com/rs/zerolog/log"
	"serenitylabs.cloud/slapi"
	"serenitylabs.cloud/slapi/pkg/api"
	"serenitylabs.cloud/slapi/pkg/fxutil"
)

func main() {
	opts := make([]fx.Option, 0)

	opts = append(opts,
		fx.Provide(slapi.New),
		fx.Provide(api.GetDefaultConfig),
		fx.Invoke(SetupServer),
		fx.WithLogger(fxutil.NewLogger),
	)

	app := fx.New(
		opts...,
	)
	if app.Err() != nil {
		panic(app.Err())
	}

	app.Run()
}

func SetupServer(lc fx.Lifecycle, shutdown fx.Shutdowner, router *slapi.Core, config *api.Config) *http.Server {
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
				if err != nil && err != http.ErrServerClosed {
					log.Error().Err(err).Send()
					shutdown.Shutdown()
				}

				addr := listener.Addr().(*net.TCPAddr)

				log.Info().Int("port", addr.Port).Str("ip", addr.IP.String()).Msg("listening")

				err = server.Serve(listener)
				if err != nil && err != http.ErrServerClosed {
					log.Error().Err(err).Send()
					shutdown.Shutdown()
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
