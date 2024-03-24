package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pprishchepa/go-bank-example/internal/config"
	"github.com/pprishchepa/go-bank-example/internal/pkg/pgxmigrator"
	"github.com/pprishchepa/go-bank-example/migrations"
	"go.uber.org/fx"
)

func newPostgresClient(lc fx.Lifecycle, conf config.Config) (*pgxpool.Pool, error) {
	connString := strings.Join([]string{
		"user=" + conf.Postgres.User,
		"password=" + conf.Postgres.Password,
		"dbname=" + conf.Postgres.Database,
		"host=" + conf.Postgres.Host,
		"port=" + fmt.Sprintf("%d", conf.Postgres.Port),
		"sslmode=" + conf.Postgres.SSLMode,
		"connect_timeout=" + fmt.Sprintf("%d", conf.Postgres.ConnTimeout),
		"pool_max_conns=" + fmt.Sprintf("%d", conf.Postgres.MaxConn),
	}, " ")

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("init pgxpool: %w", err)
	}

	lc.Append(fx.StopHook(func() {
		db.Close()
	}))

	return db, nil
}

func migrate(db *pgxpool.Pool) error {
	return pgxmigrator.NewMigrator().Up(db, migrations.FS)
}
