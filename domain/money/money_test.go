package money_test

import (
	"testing"

	"github.com/pprishchepa/go-bank-example/domain/money"
	"github.com/stretchr/testify/assert"
)

func TestMoney(t *testing.T) {
	m := money.NewFromInt(105023).
		Add(money.NewFromInt(2)).
		Sub(money.NewFromInt(1))

	assert.Equal(t, 105024, m.AsInt())
	assert.True(t, m.IsPositive())
	assert.False(t, m.IsNegative())
}
