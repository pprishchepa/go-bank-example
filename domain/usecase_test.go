package domain_test

import (
	"context"
	"testing"

	"github.com/pprishchepa/go-bank-example/domain"
	"github.com/pprishchepa/go-bank-example/domain/money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWalletUseCases_OnlyPositiveDebitAllowed(t *testing.T) {
	store := newFakeWalletStore(25, money.NewFromInt(1000))
	uc := domain.NewWalletUseCases(store)
	err := uc.DebitMoney(context.Background(), domain.DebitEntry{
		WalletID: 25,
		Amount:   money.NewFromInt(-100),
	})
	require.ErrorIs(t, err, domain.ErrInvalidAmount)
}

func TestWalletUseCases_OnlyPositiveCreditAllowed(t *testing.T) {
	store := newFakeWalletStore(25, money.NewFromInt(1000))
	uc := domain.NewWalletUseCases(store)
	err := uc.CreditMoney(context.Background(), domain.CreditEntry{
		WalletID: 25,
		Amount:   money.NewFromInt(-100),
	})
	require.ErrorIs(t, err, domain.ErrInvalidAmount)
}

func TestWalletUseCases_NegativeBalanceIsNotAllowed(t *testing.T) {
	store := newFakeWalletStore(25, money.NewFromInt(1000))
	uc := domain.NewWalletUseCases(store)
	err := uc.CreditMoney(context.Background(), domain.CreditEntry{
		WalletID: 25,
		Amount:   money.NewFromInt(2500),
	})
	require.ErrorIs(t, err, domain.ErrInsufficientFunds)
}

func TestWalletUseCases_RetrieveBalance(t *testing.T) {
	type args struct {
		walletID  int
		mockStore func(svc *MockWalletStore)
	}
	type want struct {
		assertErr     require.ErrorAssertionFunc
		assertBalance func(t require.TestingT, balance *domain.WalletBalance)
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "wallet found",
			args: args{
				walletID: 25,
				mockStore: func(svc *MockWalletStore) {
					svc.EXPECT().GetBalance(gomock.Any(), 25).Return(&domain.WalletBalance{
						WalletID: 25,
						Amount:   money.NewFromInt(251200),
					}, nil)
				},
			},
			want: want{
				assertBalance: func(t require.TestingT, balance *domain.WalletBalance) {
					assert.Equal(t, 25, balance.WalletID)
					assert.Equal(t, 251200, balance.Amount.AsInt())
				},
				assertErr: require.NoError,
			},
		},
		{
			name: "wallet not found",
			args: args{
				walletID: 25,
				mockStore: func(svc *MockWalletStore) {
					svc.EXPECT().GetBalance(gomock.Any(), 25).Return(nil, domain.ErrWalletNotFound)
				},
			},
			want: want{
				assertBalance: func(t require.TestingT, balance *domain.WalletBalance) {
					assert.Nil(t, balance)
				},
				assertErr: func(t require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(t, err, domain.ErrWalletNotFound)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			store := NewMockWalletStore(mockCtrl)
			tt.args.mockStore(store)

			_, err := domain.NewWalletUseCases(store).RetrieveBalance(context.Background(), tt.args.walletID)
			tt.want.assertErr(t, err)
		})
	}
}

type fakeWalletStore struct {
	walletID int
	amount   money.Money
}

func (f *fakeWalletStore) GetBalance(_ context.Context, walletID int) (*domain.WalletBalance, error) {
	if walletID != f.walletID {
		return nil, domain.ErrWalletNotFound
	}
	return &domain.WalletBalance{WalletID: walletID, Amount: f.amount}, nil
}

func (f *fakeWalletStore) SaveBalance(_ context.Context, balance *domain.WalletBalance) error {
	if balance.WalletID != f.walletID {
		return domain.ErrWalletNotFound
	}
	f.amount = balance.Amount
	return nil
}

func (f *fakeWalletStore) AddDebitEntry(_ context.Context, entry domain.DebitEntry) error {
	if entry.WalletID != f.walletID {
		return domain.ErrWalletNotFound
	}
	f.amount = f.amount.Add(entry.Amount)
	return nil
}

func (f *fakeWalletStore) AddCreditEntry(_ context.Context, entry domain.CreditEntry) error {
	if entry.WalletID != f.walletID {
		return domain.ErrWalletNotFound
	}
	f.amount = f.amount.Sub(entry.Amount)
	return nil
}

func newFakeWalletStore(walletID int, amount money.Money) *fakeWalletStore {
	return &fakeWalletStore{
		walletID: walletID,
		amount:   amount,
	}
}
