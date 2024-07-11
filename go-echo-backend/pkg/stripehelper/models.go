package stripehelper

import "github.com/stripe/stripe-go/v74"

// PaymentMethodWithDefault default
type PaymentMethodWithDefault struct {
	*stripe.PaymentMethod
	IsDefault bool `json:"is_default"`
}
