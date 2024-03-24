package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pprishchepa/go-bank-example/domain"
	"github.com/pprishchepa/go-bank-example/domain/money"
	"github.com/pprishchepa/go-bank-example/internal/controller/http/v1/model"
	"github.com/rs/zerolog/log"
	_ "github.com/rs/zerolog/log"
)

//go:generate go run go.uber.org/mock/mockgen -source=wallet.go -destination=wallet_mock_test.go -package=v1_test

type WalletService interface {
	GetBalance(ctx context.Context, walletID int) (*domain.WalletBalance, error)
	DebitMoney(ctx context.Context, entry domain.DebitEntry) error
	CreditMoney(ctx context.Context, entry domain.CreditEntry) error
}

type WalletRoutes struct {
	service WalletService
}

func NewWalletRoutes(service WalletService) *WalletRoutes {
	return &WalletRoutes{service: service}
}

func (r WalletRoutes) RegisterRoutes(e *gin.RouterGroup) {
	e.GET("/wallets/:wallet/balance", r.retrieveBalance)
	e.POST("/wallets/:wallet/debit", r.debitMoney)
	e.POST("/wallets/:wallet/credit", r.creditMoney)
}

func (r WalletRoutes) retrieveBalance(c *gin.Context) {
	var req model.GetBalanceRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := r.service.GetBalance(c.Request.Context(), req.WalletID)
	if err != nil {
		log.Err(err).Int("walletId", req.WalletID).Msg("could not get balance")
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": model.BalanceResponse{
		WalletID: balance.WalletID,
		Amount:   balance.Amount.AsInt(),
	}})
}

func (r WalletRoutes) debitMoney(c *gin.Context) {
	var req model.DebitMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.service.DebitMoney(c.Request.Context(), domain.DebitEntry{
		WalletID: req.WalletID,
		Amount:   money.NewFromInt(req.Amount),
	})
	if err != nil {
		log.Err(err).Int("walletId", req.WalletID).Msg("could not debit money")
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.Status(http.StatusOK)
}

func (r WalletRoutes) creditMoney(c *gin.Context) {
	var req model.CreditMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.service.CreditMoney(c.Request.Context(), domain.CreditEntry{
		WalletID: req.WalletID,
		Amount:   money.NewFromInt(req.Amount),
	})
	if err != nil {
		log.Err(err).Int("walletId", req.WalletID).Msg("could not credit money")
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.Status(http.StatusOK)
}
