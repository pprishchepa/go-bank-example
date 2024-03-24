package app

import (
	"io"
	stdlog "log"
	"os"

	"github.com/pprishchepa/go-casino-example/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newLogger(conf config.Config) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(conf.Log.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	writer := io.Writer(os.Stdout)
	if conf.Log.Pretty {
		writer = zerolog.NewConsoleWriter()
	}

	logger := zerolog.New(writer).Level(level)

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	log.Logger = logger.With().Timestamp().Logger()

	return log.Logger, nil
}
