package domain

import (
	"github.com/pprishchepa/go-bank-example/domain/money"
)

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
