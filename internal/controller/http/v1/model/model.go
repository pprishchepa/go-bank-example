package model

type GetBalanceRequest struct {
	WalletID int `uri:"wallet" binding:"required,gt=0"`
}

type DebitMoneyRequest struct {
	WalletID int `uri:"wallet" binding:"required,gt=0"`
	Amount   int `json:"amount" binding:"gt=0"`
}

type CreditMoneyRequest struct {
	WalletID int `uri:"wallet" binding:"required,gt=0"`
	Amount   int `json:"amount" binding:"gt=0"`
}

type BalanceResponse struct {
	WalletID int `json:"walletId"`
	Amount   int `json:"amount"`
}
