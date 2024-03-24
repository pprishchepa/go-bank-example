package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pprishchepa/go-casino-example/internal/service"
	"github.com/pprishchepa/go-casino-example/internal/storage/postgres"
)

type walletStoreTxFactory struct {
	factory *postgres.WalletStoreTxFactory
}

func (f walletStoreTxFactory) NewTx(ctx context.Context) (service.WalletStoreTx, error) {
	v, err := f.factory.NewTx(ctx)
	return v, err
}

func newWalletStoreTxFactory(db *pgxpool.Pool) service.WalletStoreTxFactory {
	return &walletStoreTxFactory{factory: postgres.NewWalletStoreTxFactory(db)}
}
