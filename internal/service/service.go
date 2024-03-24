package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pprishchepa/go-bank-example/domain"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/singleflight"
)

//go:generate go run go.uber.org/mock/mockgen -source=service.go -destination=service_mock_test.go -package=service_test

type (
	WalletStoreTx interface {
		domain.WalletStore
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
	WalletStoreTxFactory interface {
		NewTx(ctx context.Context) (WalletStoreTx, error)
	}
	WalletCacheStore interface {
		SaveBalance(ctx context.Context, balance *domain.WalletBalance) error
		GetBalance(ctx context.Context, walletID int) (*domain.WalletBalance, error)
	}
)

type WalletService struct {
	txFactory WalletStoreTxFactory
	cache     WalletCacheStore
	group     singleflight.Group
}

func NewWalletService(txFactory WalletStoreTxFactory, cache WalletCacheStore) *WalletService {
	return &WalletService{
		txFactory: txFactory,
		cache:     cache,
	}
}

func (s *WalletService) GetBalance(ctx context.Context, walletID int) (*domain.WalletBalance, error) {
	balance, err := s.cache.GetBalance(ctx, walletID)
	if err != nil {
		log.Warn().Err(err).Int("walletID", walletID).Msg("could not get balance from cache")
	}
	if balance != nil {
		return balance, nil
	}

	v, err, _ := s.group.Do(strconv.Itoa(walletID), func() (interface{}, error) {
		err := s.runOnceTx(ctx, func(tx WalletStoreTx) error {
			balance, err = domain.NewWalletUseCases(tx).RetrieveBalance(ctx, walletID)
			return err
		})
		if err != nil {
			return nil, err
		}
		if err = s.cache.SaveBalance(ctx, balance); err != nil {
			log.Warn().Err(err).Int("walletID", walletID).Msg("could not save balance to cache")
		}
		return balance, nil
	})
	if err != nil {
		return nil, err
	}

	return v.(*domain.WalletBalance), nil
}

func (s *WalletService) DebitMoney(ctx context.Context, entry domain.DebitEntry) error {
	var balance *domain.WalletBalance

	err := s.runOrRepeatTx(ctx, func(tx WalletStoreTx) error {
		usecase := domain.NewWalletUseCases(tx)
		if err := usecase.DebitMoney(ctx, entry); err != nil {
			return fmt.Errorf("debit money: %w", err)
		}
		var err error
		if balance, err = usecase.RetrieveBalance(ctx, entry.WalletID); err != nil {
			return fmt.Errorf("retrieve balance: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := s.cache.SaveBalance(ctx, balance); err != nil {
		log.Warn().Err(err).Int("walletID", entry.WalletID).Msg("could not update balance in cache")
	}
	return nil
}

func (s *WalletService) CreditMoney(ctx context.Context, entry domain.CreditEntry) error {
	var balance *domain.WalletBalance

	err := s.runOrRepeatTx(ctx, func(tx WalletStoreTx) error {
		usecase := domain.NewWalletUseCases(tx)
		if err := usecase.CreditMoney(ctx, entry); err != nil {
			return fmt.Errorf("credit money: %w", err)
		}
		var err error
		if balance, err = usecase.RetrieveBalance(ctx, entry.WalletID); err != nil {
			return fmt.Errorf("retrieve balance: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := s.cache.SaveBalance(ctx, balance); err != nil {
		log.Warn().Err(err).Int("walletID", entry.WalletID).Msg("could not update balance in cache")
	}
	return nil
}

func (s *WalletService) runOrRepeatTx(ctx context.Context, fn func(tx WalletStoreTx) error) error {
	if err := s.runOnceTx(ctx, fn); err == nil {
		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 200 * time.Millisecond
	expBackoff.MaxInterval = 1 * time.Second
	expBackoff.MaxElapsedTime = 5 * time.Second

	return backoff.Retry(func() error {
		return s.runOnceTx(ctx, fn)
	}, expBackoff)
}

func (s *WalletService) runOnceTx(ctx context.Context, fn func(tx WalletStoreTx) error) error {
	tx, err := s.txFactory.NewTx(ctx)
	if err != nil {
		return fmt.Errorf("new tx: %w", err)
	}

	if fnErr := fn(tx); fnErr != nil {
		if err := tx.Rollback(ctx); err != nil {
			log.Warn().Err(err).Msg("could not rollback tx")
		}
		return fnErr
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
