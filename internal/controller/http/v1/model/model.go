package model

type WalletRequest struct {
	ID int `uri:"wallet" binding:"required,gt=0"`
}

type DebitMoneyRequest struct {
	Amount int `json:"amount" binding:"gt=0"`
}

type CreditMoneyRequest struct {
	Amount int `json:"amount" binding:"gt=0"`
}

type BalanceResponse struct {
	WalletID int `json:"walletId"`
	Amount   int `json:"amount"`
}
