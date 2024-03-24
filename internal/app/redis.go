package app

import (
	"context"
	"fmt"

	"github.com/pprishchepa/go-bank-example/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func newRedisClient(lc fx.Lifecycle, conf config.Config) (*redis.Ring, error) {
	redis.SetLogger(&redisLoggerAdapter{})

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			conf.Redis.Host: fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		},
	})

	lc.Append(fx.StopHook(func() {
		if err := ring.Close(); err != nil {
			log.Warn().Msg(err.Error())
		}
	}))

	return ring, nil
}

type redisLoggerAdapter struct{}

func (r *redisLoggerAdapter) Printf(_ context.Context, format string, v ...interface{}) {
	log.Warn().Msg(fmt.Sprintf(format, v...))
}
