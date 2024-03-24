package domain

import (
	"errors"

	"github.com/pprishchepa/go-bank-example/domain/money"
)

var ErrWalletNotFound = errors.New("wallet not found")
var ErrInsufficientFunds = errors.New("insufficient funds")

type Wallet struct {
	ID int
}

type WalletBalance struct {
	WalletID int
	Amount   money.Money
}

type DebitEntry struct {
	WalletID int
	Amount   money.Money
}

type CreditEntry struct {
	WalletID int
	Amount   money.Money
}
