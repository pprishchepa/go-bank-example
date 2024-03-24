package v1

import (
	"context"
	"errors"
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
	var reqWallet model.WalletRequest
	if err := c.ShouldBindUri(&reqWallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := r.service.GetBalance(c.Request.Context(), reqWallet.ID)
	if err != nil {
		r.handleError(c, err, reqWallet.ID, "could not get balance")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": model.BalanceResponse{
		WalletID: balance.WalletID,
		Amount:   balance.Amount.AsInt(),
	}})
}

func (r WalletRoutes) debitMoney(c *gin.Context) {
	var reqWallet model.WalletRequest
	if err := c.ShouldBindUri(&reqWallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody model.DebitMoneyRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.service.DebitMoney(c.Request.Context(), domain.DebitEntry{
		WalletID: reqWallet.ID,
		Amount:   money.NewFromInt(reqBody.Amount),
	})
	if err != nil {
		r.handleError(c, err, reqWallet.ID, "could not debit money")
		return
	}

	c.Status(http.StatusOK)
}

func (r WalletRoutes) creditMoney(c *gin.Context) {
	var reqWallet model.WalletRequest
	if err := c.ShouldBindUri(&reqWallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody model.CreditMoneyRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.service.CreditMoney(c.Request.Context(), domain.CreditEntry{
		WalletID: reqWallet.ID,
		Amount:   money.NewFromInt(reqBody.Amount),
	})
	if err != nil {
		r.handleError(c, err, reqWallet.ID, "could not credit money")
		return
	}

	c.Status(http.StatusOK)
}

func (r WalletRoutes) handleError(c *gin.Context, err error, walletID int, msg string) {
	if errors.Is(err, domain.ErrWalletNotFound) {
		log.Debug().Err(err).Int("walletId", walletID).Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet not found"})
		return
	}

	if errors.Is(err, domain.ErrInsufficientFunds) {
		log.Debug().Err(err).Int("walletId", walletID).Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}

	log.Err(err).Int("walletId", walletID).Msg(msg)
	c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
}
