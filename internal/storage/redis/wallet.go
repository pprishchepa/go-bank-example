package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/pprishchepa/go-bank-example/domain"
	"github.com/redis/go-redis/v9"
)

type WalletCacheStore struct {
	cache *cache.Cache
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
		Ctx:   ctx,
		Key:   fmt.Sprintf("account:%d:balance", balance.WalletID),
		Value: balance,
	})
}

func (s WalletCacheStore) GetBalance(ctx context.Context, walletID int) (*domain.WalletBalance, error) {
	var balance domain.WalletBalance

	if err := s.cache.Get(ctx, s.newKey(walletID), &balance); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &balance, nil
}

func (s WalletCacheStore) newKey(walletID int) string {
	return fmt.Sprintf("account:%d:balance", walletID)
}
