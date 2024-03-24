package domain

import "context"

//go:generate go run go.uber.org/mock/mockgen -source=storage.go -destination=storage_mock_test.go -package=domain_test

type WalletStore interface {
	GetBalance(ctx context.Context, walletID int) (*WalletBalance, error)
	SaveBalance(ctx context.Context, balance *WalletBalance) error
	AddDebitEntry(ctx context.Context, entry DebitEntry) error
	AddCreditEntry(ctx context.Context, entry CreditEntry) error
}
