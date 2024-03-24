package domain

import (
	"context"
	"fmt"
)

type WalletUseCases struct {
	storage WalletStore
}

func NewWalletUseCases(storage WalletStore) *WalletUseCases {
	return &WalletUseCases{storage: storage}
}

func (c WalletUseCases) RetrieveBalance(ctx context.Context, walletID int) (*WalletBalance, error) {
	return c.storage.GetBalance(ctx, walletID)
}

func (c WalletUseCases) DebitMoney(ctx context.Context, entry DebitEntry) error {
	if !entry.Amount.IsPositive() {
		return fmt.Errorf("invalid amount")
	}

	balance, err := c.storage.GetBalance(ctx, entry.WalletID)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}

	balance.Amount = balance.Amount.Add(entry.Amount)
	if err := c.storage.SaveBalance(ctx, balance); err != nil {
		return fmt.Errorf("save balance: %w", err)
	}

	if err := c.storage.AddDebitEntry(ctx, entry); err != nil {
		return fmt.Errorf("add debit entry: %w", err)
	}

	return nil
}

func (c WalletUseCases) CreditMoney(ctx context.Context, entry CreditEntry) error {
	if !entry.Amount.IsPositive() {
		return fmt.Errorf("invalid amount")
	}

	balance, err := c.storage.GetBalance(ctx, entry.WalletID)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}

	newAmount := balance.Amount.Sub(entry.Amount)
	if newAmount.IsNegative() {
		return ErrInsufficientFunds
	}

	balance.Amount = newAmount
	if err := c.storage.SaveBalance(ctx, balance); err != nil {
		return fmt.Errorf("save balance: %w", err)
	}

	if err := c.storage.AddCreditEntry(ctx, entry); err != nil {
		return fmt.Errorf("add credit entry: %w", err)
	}

	return nil
}
