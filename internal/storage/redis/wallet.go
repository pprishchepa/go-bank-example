package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/pprishchepa/go-bank-example/domain"
	"github.com/pprishchepa/go-bank-example/domain/money"
	"github.com/redis/go-redis/v9"
)

type WalletCacheStore struct {
	cache *cache.Cache
}

type cachedWalletBalance struct {
	WalletID int
	Amount   int
}

func NewWalletCacheStore(ring *redis.Ring) *WalletCacheStore {
	return &WalletCacheStore{
		cache: cache.New(&cache.Options{
			Redis:      ring,
			LocalCache: cache.NewTinyLFU(1024, time.Minute),
		}),
	}
}

func (s WalletCacheStore) SaveBalance(ctx context.Context, balance *domain.WalletBalance) error {
	return s.cache.Set(&cache.Item{
		Ctx: ctx,
		Key: s.newKey(balance.WalletID),
		Value: cachedWalletBalance{
			WalletID: balance.WalletID,
			Amount:   balance.Amount.AsInt(),
		},
	})
}

func (s WalletCacheStore) GetBalance(ctx context.Context, walletID int) (*domain.WalletBalance, error) {
	var cachedBalance cachedWalletBalance

	if err := s.cache.Get(ctx, s.newKey(walletID), &cachedBalance); err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			return nil, nil
		}
		return nil, fmt.Errorf("get: %w", err)
	}

	return &domain.WalletBalance{
		WalletID: cachedBalance.WalletID,
		Amount:   money.NewFromInt(cachedBalance.Amount),
	}, nil
}

func (s WalletCacheStore) newKey(walletID int) string {
	return fmt.Sprintf("account:%d:balance", walletID)
}
