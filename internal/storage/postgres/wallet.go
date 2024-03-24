package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pprishchepa/go-bank-example/domain"
	"github.com/shopspring/decimal"
)

type WalletStoreTxFactory struct {
	db *pgxpool.Pool
}

func NewWalletStoreTxFactory(db *pgxpool.Pool) *WalletStoreTxFactory {
	return &WalletStoreTxFactory{db: db}
}

func (f WalletStoreTxFactory) NewTx(ctx context.Context) (*WalletTx, error) {
	tx, err := f.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	return &WalletTx{tx: tx}, nil
}

type WalletTx struct {
	tx pgx.Tx
}

func (a WalletTx) GetBalance(ctx context.Context, walletID int) (*domain.WalletBalance, error) {
	sql := `SELECT amount FROM wallet_balance WHERE wallet_id = $1`

	var amount int

	if err := a.tx.QueryRow(ctx, sql, walletID).Scan(&amount); err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &domain.WalletBalance{
		WalletID: walletID,
		Amount:   intToDec(amount),
	}, nil
}

func (a WalletTx) SaveBalance(ctx context.Context, balance *domain.WalletBalance) error {
	sql := `
		INSERT INTO wallet_balance (wallet_id, amount) 
		VALUES ($1, $2)
		ON CONFLICT (wallet_id) DO UPDATE SET amount = $2`

	_, err := a.tx.Exec(ctx, sql, balance.WalletID, balance.Amount)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (a WalletTx) AddDebitEntry(ctx context.Context, entry domain.DebitEntry) error {
	sql := `
		INSERT INTO wallet_entry (wallet_id, debit_amount) 
		VALUES ($1, $2)`

	_, err := a.tx.Exec(ctx, sql, entry.WalletID, entry.Amount)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (a WalletTx) AddCreditEntry(ctx context.Context, entry domain.CreditEntry) error {
	sql := `
		INSERT INTO wallet_entry (wallet_id, credit_amount) 
		VALUES ($1, $2)`

	_, err := a.tx.Exec(ctx, sql, entry.WalletID, entry.Amount)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (a WalletTx) Commit(ctx context.Context) error {
	return a.tx.Commit(ctx)
}

func (a WalletTx) Rollback(ctx context.Context) error {
	return a.tx.Rollback(ctx)
}

var precision = decimal.NewFromInt(100)

func intToDec(v int) decimal.Decimal {
	return decimal.NewFromInt(int64(v)).Div(precision)
}

//func dec2Int(d decimal.Decimal) int {
//	return int(d.IntPart()*precision) + int(d.Exponent())
//}
