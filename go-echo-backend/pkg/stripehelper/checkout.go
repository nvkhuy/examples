package stripehelper

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/charge"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

type CreateCheckoutSessionParams struct {
	SuccessURL         string
	CancelledURL       string
	PaymentMethodTypes []string
	Amount             int64
	Currency           enums.Currency
	StripeCustomerID   string
	ReferenceID        string
}

func (client *StripeClient) CreateCheckoutSession(req *CreateCheckoutSessionParams) (string, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String(req.SuccessURL),
		PaymentMethodTypes: stripe.StringSlice(req.PaymentMethodTypes),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Inflow order"),
					},
					Currency:   stripe.String(string(req.Currency)),
					UnitAmount: stripe.Int64(req.Amount),
				},
			},
		},
		Customer:          stripe.String(req.StripeCustomerID),
		ClientReferenceID: stripe.String(req.ReferenceID),
	}

	if req.CancelledURL != "" {
		params.CancelURL = stripe.String(req.CancelledURL)
	}
	s, err := session.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil

}

type CreateCheckoutProductURLParams struct {
	SuccessURL   string
	CancelledURL string
	ReferenceID  string

	LineItems []*stripe.CheckoutSessionLineItemParams

	Metadata map[string]string

	CustomerEmail string
}

func (client *StripeClient) CreateCheckoutProductURL(req *CreateCheckoutProductURLParams) (string, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL:        stripe.String(req.SuccessURL),
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:         req.LineItems,
		ClientReferenceID: stripe.String(req.ReferenceID),
	}
	if len(req.Metadata) > 0 {
		params.Metadata = make(map[string]string)
		for key, value := range req.Metadata {
			params.Metadata[key] = value
		}
	}

	if req.CustomerEmail != "" {
		params.CustomerEmail = stripe.String(req.CustomerEmail)
	}

	if req.CancelledURL != "" {
		params.CancelURL = stripe.String(req.CancelledURL)
	}
	s, err := session.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil

}

func (client *StripeClient) GetCheckout(id string) (*stripe.CheckoutSession, error) {
	var params = &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")
	return session.Get(id, params)
}

func (client *StripeClient) GetCharge(id string) (*stripe.Charge, error) {
	var params = &stripe.ChargeParams{}
	params.AddExpand("balance_transaction")
	return charge.Get(id, params)
}
