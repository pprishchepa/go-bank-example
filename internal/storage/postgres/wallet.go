package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pprishchepa/go-bank-example/domain"
	"github.com/pprishchepa/go-bank-example/domain/money"
	"github.com/pprishchepa/go-bank-example/internal/entity"
)

// See https://www.postgresql.org/docs/current/errcodes-appendix.html
const errorCodeSerializationFailure = "40001"

type WalletStoreTxFactory struct {
	db *pgxpool.Pool
}

func NewWalletStoreTxFactory(db *pgxpool.Pool) *WalletStoreTxFactory {
	return &WalletStoreTxFactory{db: db}
}

func (f WalletStoreTxFactory) NewTx(ctx context.Context) (*WalletStore, error) {
	tx, err := f.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	return &WalletStore{tx: tx}, nil
}

type WalletStore struct {
	tx pgx.Tx
}

func (s WalletStore) GetBalance(ctx context.Context, walletID int) (*domain.WalletBalance, error) {
	sql := `SELECT amount FROM wallet_balance WHERE wallet_id = $1`

	var amount int

	if err := s.tx.QueryRow(ctx, sql, walletID).Scan(&amount); err != nil {
		return nil, fmt.Errorf("query row: %w", s.recognizeError(err))
	}

	return &domain.WalletBalance{
		WalletID: walletID,
		Amount:   money.NewFromInt(amount),
	}, nil
}

func (s WalletStore) SaveBalance(ctx context.Context, balance *domain.WalletBalance) error {
	sql := `
		INSERT INTO wallet_balance (wallet_id, amount) 
		VALUES ($1, $2)
		ON CONFLICT (wallet_id) DO UPDATE SET amount = $2`

	_, err := s.tx.Exec(ctx, sql, balance.WalletID, balance.Amount.AsInt())
	if err != nil {
		return fmt.Errorf("exec: %w", s.recognizeError(err))
	}

	return nil
}

func (s WalletStore) AddDebitEntry(ctx context.Context, entry domain.DebitEntry) error {
	sql := `
		INSERT INTO wallet_entry (wallet_id, debit_amount) 
		VALUES ($1, $2)`

	_, err := s.tx.Exec(ctx, sql, entry.WalletID, entry.Amount.AsInt())
	if err != nil {
		return fmt.Errorf("exec: %w", s.recognizeError(err))
	}

	return nil
}

func (s WalletStore) AddCreditEntry(ctx context.Context, entry domain.CreditEntry) error {
	sql := `
		INSERT INTO wallet_entry (wallet_id, credit_amount) 
		VALUES ($1, $2)`

	_, err := s.tx.Exec(ctx, sql, entry.WalletID, entry.Amount.AsInt())
	if err != nil {
		return fmt.Errorf("exec: %w", s.recognizeError(err))
	}

	return nil
}

func (s WalletStore) Commit(ctx context.Context) error {
	return s.recognizeError(s.tx.Commit(ctx))
}

func (s WalletStore) Rollback(ctx context.Context) error {
	return s.recognizeError(s.tx.Rollback(ctx))
}

func (s WalletStore) recognizeError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == errorCodeSerializationFailure {
		return entity.ErrTxConflict
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrWalletNotFound
	}

	return err
}
