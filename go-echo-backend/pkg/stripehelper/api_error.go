package stripehelper

import "github.com/rotisserie/eris"

var (
	ErrNoDefaultPaymentMethod = eris.New("No have default payment method")
)
