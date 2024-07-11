package price

import (
	"database/sql/driver"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

// Price price https://github.com/shopspring/decimal
type Price decimal.Decimal

func (price *Price) Scan(value interface{}) error {
	if str, ok := value.(string); ok {
		if strings.EqualFold(str, "NaN") || str == "" {
			value = "0.0"
		}

	}

	return (*decimal.Decimal)(price).Scan(value)
}

func (price Price) Value() (driver.Value, error) {
	return (decimal.Decimal)(price).Value()
}

func (price Price) ToPtr() *Price {
	return &price
}

func (price *Price) ToValue() Price {
	return *price
}

func (price Price) MarshalJSON() ([]byte, error) {
	return (decimal.Decimal)(price).MarshalJSON()
}

func (price *Price) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	return (*decimal.Decimal)(price).UnmarshalJSON(b)
}

func (price Price) MarshalCSV() ([]byte, error) {
	return price.MarshalJSON()
}

func (price *Price) UnmarshalCSV(data []byte) error {
	var amount = strings.ReplaceAll(string(data), ",", "")
	return price.UnmarshalJSON([]byte(amount))
}

// ///////
func NewFromPtr(p *Price) Price {
	if p != nil {
		return *p
	}

	return NewFromFloat(0)
}

func NewFromString(value string) Price {
	if value == "" {
		return NewFromFloat(0)
	}

	d, err := decimal.NewFromString(value)
	if err != nil {
		return NewFromFloat(0)
	}

	return Price(d)
}

func NewFromFloat(value float64) Price {
	return Price(decimal.NewFromFloat(value))
}

func NewFromDecimalPtr(value *decimal.Decimal) Price {
	if value == nil {
		return NewFromFloat(0)
	}

	return Price(*value)
}

func NewFromDecimal(value decimal.Decimal) Price {
	return Price(value)
}

func NewFromInt32(value int32) Price {
	return Price(decimal.NewFromInt32(value))
}

func NewFromInt(value int64) Price {
	return Price(decimal.NewFromInt(value))
}

func (price Price) Pow(exponent int64) Price {
	var ex = NewFromInt(exponent).Decimal()
	return Price(price.Decimal().Pow(ex))
}

func (price Price) Abs() Price {
	return Price(price.Decimal().Abs())
}

func (price Price) Decimal() decimal.Decimal {
	return decimal.Decimal(price)
}

func (price Price) Multiple(v Price, skipRound ...bool) Price {
	var v1 = price.Decimal()
	var v2 = v.Decimal()

	if len(skipRound) > 0 && skipRound[0] {

	} else {
		return Price(v1.Mul(v2)).Round()
	}

	return Price(v1.Mul(v2))
}

func (price Price) MultipleFloat64(v float64, skipRound ...bool) Price {
	var v2 = NewFromFloat(v)
	return price.Multiple(v2, skipRound...)

}

func (price Price) MultipleInt(v int64, skipRound ...bool) Price {
	var v2 = NewFromInt(v)
	return price.Multiple(v2, skipRound...)

}

func (price Price) Round(places ...int32) Price {
	var p int32 = 2
	if len(places) > 0 {
		p = places[0]
	}
	return Price(price.Decimal().Round(p))
}

func (price Price) Sub(other Price, skipRound ...bool) Price {
	var v1 = price.Decimal()
	var v2 = other.Decimal()

	if len(skipRound) > 0 && skipRound[0] {

	} else {
		return Price(v1.Sub(v2)).Round()
	}

	return Price(v1.Sub(v2))

}

func (price Price) Add(other Price, skipRound ...bool) Price {
	var v1 = price.Decimal()
	var v2 = other.Decimal()

	if len(skipRound) > 0 && skipRound[0] {

	} else {
		return Price(v1.Add(v2)).Round()
	}

	return Price(v1.Add(v2))

}

func (price Price) AddPtr(other *Price, skipRound ...bool) Price {
	var v = NewFromPtr(other)

	return price.Add(v)

}
func (price Price) SubPtr(other *Price, skipRound ...bool) Price {
	var v = NewFromPtr(other)

	return price.Sub(v)

}
func (price Price) Div(v Price, skipRound ...bool) Price {
	var v1 = price.Decimal()
	var v2 = v.Decimal()

	if len(skipRound) > 0 && skipRound[0] {

	} else {
		return Price(v1.Div(v2)).Round()
	}

	return Price(v1.Div(v2))

}

func (price Price) DivInt(v int64, skipRound ...bool) Price {
	return price.Div(NewFromInt(v), skipRound...)
}

func (price Price) ToFloat64() float64 {
	f, _ := price.Decimal().Float64()
	return f
}

func (price Price) ToInt64() int64 {
	f, _ := price.Decimal().BigFloat().Int64()

	return f
}

func (price Price) Equal(other float64) bool {
	return price.Decimal().Equal(decimal.NewFromFloat(other))
}

func (price Price) LessThan(other float64) bool {
	return price.Decimal().LessThan(decimal.NewFromFloat(other))
}

func (price Price) LessThanOrEqual(other float64) bool {
	return price.Decimal().LessThanOrEqual(decimal.NewFromFloat(other))
}

func (price Price) GreaterThan(other float64) bool {
	return price.Decimal().GreaterThan(decimal.NewFromFloat(other))
}

func (price Price) GreaterThanOrEqual(other float64) bool {
	return price.Decimal().GreaterThanOrEqual(decimal.NewFromFloat(other))
}

func (price Price) InRange(greaterThanOrEqual, lessThanValue float64) bool {
	return price.GreaterThanOrEqual(greaterThanOrEqual) && price.LessThan(lessThanValue)
}

func (price Price) FormatMoney(currency enums.Currency) string {
	var symbol = "$"
	var precision = 2
	switch currency {
	case enums.USD:
		precision = 2
		symbol = "$"
	case enums.VND:
		precision = 0
		symbol = "đ"
	}
	var ac = accounting.Accounting{Symbol: symbol, Precision: precision}

	return ac.FormatMoneyDecimal(price.Decimal())
}
