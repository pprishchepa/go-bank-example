package app

import (
	"net/http"

	"github.com/pprishchepa/go-casino-example/internal/config"
	httpctrl "github.com/pprishchepa/go-casino-example/internal/controller/http"
	httpv1 "github.com/pprishchepa/go-casino-example/internal/controller/http/v1"
	"github.com/pprishchepa/go-casino-example/internal/pkg/fxlog"
	"github.com/pprishchepa/go-casino-example/internal/service"
	"github.com/pprishchepa/go-casino-example/internal/storage/postgres"
	"github.com/pprishchepa/go-casino-example/internal/storage/redis"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.NewConfig,
			newLogger,
			newPostgresClient,
			newRedisClient,
			redis.NewWalletCacheStore,
			postgres.NewWalletStoreTxFactory,
			newWalletStoreTxFactory,
			service.NewWalletService,
			httpv1.NewWalletRoutes,
			httpctrl.NewRouter,
			newHTTPServer,
			func(v *service.WalletService) httpv1.WalletService { return v },
			func(v *redis.WalletCacheStore) service.WalletCacheStore { return v },
		),
		fx.WithLogger(func(logger zerolog.Logger) fxevent.Logger {
			return fxlog.NewZerologAdapter(logger.With().Str("logger", "fx").Logger())
		}),
		fx.Invoke(automaxprocs),
		fx.Invoke(migrate),
		fx.Invoke(func(*http.Server) {}),
	)
}

func automaxprocs() error {
	_, err := maxprocs.Set(maxprocs.Logger(func(s string, i ...interface{}) {
		log.Info().Str("logger", "automaxprocs").Msgf(s, i...)
	}))
	return err
}
