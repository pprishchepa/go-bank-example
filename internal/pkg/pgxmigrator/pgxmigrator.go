package pgxmigrator

import (
	"embed"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

type Migrator struct{}

func NewMigrator() Migrator {
	return Migrator{}
}

func (m Migrator) Up(db *pgxpool.Pool, fs embed.FS) error {
	hostname := net.JoinHostPort(db.Config().ConnConfig.Host, strconv.Itoa(int(db.Config().ConnConfig.Port)))

	mg, err := m.newInstance(db, fs)
	if err != nil {
		return fmt.Errorf("init migrator: %s: %w", hostname, err)
	}
	defer func() {
		if _, err := mg.Close(); err != nil {
			log.Err(err).Str("dbhost", hostname).Msg("could not close migrator")
		}
	}()

	log.Info().Str("dbhost", hostname).Msg("migrating...")

	if err := mg.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
			log.Info().Str("dbhost", hostname).Msg(err.Error())
			return nil
		}
		log.Err(err).Str("dbhost", hostname).Msg("could not migrate")
		return fmt.Errorf("migrate: %s: %w", hostname, err)
	}

	log.Info().Str("dbhost", hostname).Msg("migration succeeded")
	return nil
}

func (m Migrator) newInstance(db *pgxpool.Pool, fs embed.FS) (*migrate.Migrate, error) {
	source, err := iofs.New(fs, ".")
	if err != nil {
		return nil, fmt.Errorf("create source: %w", err)
	}

	stdDB := stdlib.OpenDBFromPool(db)

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 2 * time.Second
	expBackoff.MaxInterval = 5 * time.Second
	expBackoff.MaxElapsedTime = 5 * time.Minute

	var driver database.Driver
	err = backoff.Retry(func() error {
		var err error
		driver, err = postgres.WithInstance(stdDB, &postgres.Config{})
		if err != nil {
			log.Warn().Err(err).Msg("database connection issue, retrying...")
			return err
		}
		return nil
	}, expBackoff)
	if err != nil {
		return nil, fmt.Errorf("create driver: %w", err)
	}

	return migrate.NewWithInstance("iofs", source, db.Config().ConnConfig.Database, driver)
}
