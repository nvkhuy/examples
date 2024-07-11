package stripehelper

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/rotisserie/eris"
)

type Config struct {
	Currency           enums.Currency
	SmallestUnitFactor int64
	SmallestAmount     float64
	AdditionalFee      float64
	CalcFees           func(amount price.Price) price.Price
	TransactionFee     float64
}

var CurrencyConfig = map[enums.Currency]*Config{
	enums.SGD: {
		Currency:           enums.SGD,
		SmallestUnitFactor: 100, // cent -> dollar,
		SmallestAmount:     0.5, //  0.5$
		AdditionalFee:      0.3,
		TransactionFee:     0.04, // 4%
		CalcFees: func(amount price.Price) price.Price {
			return price.NewFromFloat(0)
		},
	},
	enums.USD: {
		Currency:           enums.USD,
		SmallestUnitFactor: 100, // cent -> dollar,
		SmallestAmount:     0.5, //  0.5$
		AdditionalFee:      0.3,
		TransactionFee:     0.04, // 4%
		CalcFees: func(amount price.Price) price.Price {
			return price.NewFromFloat(0)
		},
	},
	enums.VND: {
		Currency:           enums.VND,
		SmallestUnitFactor: 1,    // VND
		SmallestAmount:     1200, //  50 cents = 11.6 ~ 12 VND
		AdditionalFee:      8000,
		TransactionFee:     0.04, // 4%
		CalcFees: func(amount price.Price) price.Price {
			return price.NewFromFloat(0)

		},
	},
}

func GetCurrencyConfig(currency enums.Currency) (*Config, error) {
	c, ok := CurrencyConfig[currency]
	if !ok {
		return nil, eris.Wrap(errs.ErrCountryNotSupported.WithMessage(fmt.Sprintf("%v is not supported", currency)), "")
	}

	return c, nil
}
