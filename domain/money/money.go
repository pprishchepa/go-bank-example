package money

import (
	"encoding/json"
	"math"

	"github.com/shopspring/decimal"
)

const decimalLen = 3

var decimalVal = decimal.NewFromInt(int64(math.Pow(10, float64(decimalLen))))

type Money struct {
	dec decimal.Decimal
}

func (m Money) IsPositive() bool {
	return m.dec.IsPositive()
}

func (m Money) IsNegative() bool {
	return m.dec.IsNegative()
}

func (m Money) Add(v Money) Money {
	return Money{m.dec.Add(v.dec)}
}

func (m Money) Sub(v Money) Money {
	return Money{m.dec.Sub(v.dec)}
}

func (m Money) AsInt() int {
	return int(m.dec.Mul(decimalVal).IntPart())
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.AsInt())
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*m = NewFromInt(v)
	return nil
}

func NewFromDecimal(v decimal.Decimal) Money {
	return Money{v}
}

func NewFromInt(v int) Money {
	return Money{decimal.NewFromInt(int64(v)).Div(decimalVal)}
}
