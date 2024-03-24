package domain

import (
	"errors"
)

var ErrWalletNotFound = errors.New("wallet not found")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrInvalidAmount = errors.New("invalid amount")
