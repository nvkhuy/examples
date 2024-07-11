package stripehelper

import (
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/rotisserie/eris"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/refund"
	"github.com/stripe/stripe-go/v74/setupintent"
)

// CreatePaymentIntentParams params
type CreatePaymentIntentParams struct {
	Amount                  int64
	Currency                enums.Currency
	Description             string
	PaymentMethodID         string
	CustomerID              string
	IsCaptureMethodManually bool
	Metadata                map[string]string
	PaymentMethodTypes      []string

	Shipping *stripe.ShippingDetailsParams
}

// CreatePaymentIntent payment intent
func (client *StripeClient) CreatePaymentIntent(params CreatePaymentIntentParams) (*stripe.PaymentIntent, error) {
	var intentParams = &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(params.Amount),
		Currency: stripe.String(string(params.Currency.DefaultIfInvalid())),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Confirm:       stripe.Bool(true),
		PaymentMethod: stripe.String(params.PaymentMethodID),
		Description:   stripe.String(params.Description),
		Customer:      stripe.String(params.CustomerID),
		Shipping:      params.Shipping,
	}

	if len(params.PaymentMethodTypes) > 0 {
		intentParams.PaymentMethodTypes = stripe.StringSlice(params.PaymentMethodTypes)
	}

	for key, value := range params.Metadata {
		intentParams.AddMetadata(key, value)
	}

	if params.IsCaptureMethodManually {
		intentParams.CaptureMethod = stripe.String(string(stripe.PaymentIntentCaptureMethodManual))
	}

	pi, err := paymentintent.New(intentParams)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			if stripeErr.Code == stripe.ErrorCodeBalanceInsufficient {
				return nil, eris.Wrap(errs.ErrPreAuthCardFailed, "")
			}
		}

		return nil, eris.Wrap(err, "")
	}

	if pi.ID == "" {
		return nil, eris.Wrap(errs.ErrStripePaymentIntentIDInvalid, "")
	}

	return pi, nil
}

// CapturePaymentIntentParams params
type CapturePaymentIntentParams struct {
	PaymentIntentID string
	Amount          int64
}

// CapturePaymentIntent capture payment intent
func (client *StripeClient) CapturePaymentIntent(params CapturePaymentIntentParams) (*stripe.PaymentIntent, error) {
	var value = params.Amount

	pi, err := paymentintent.Get(params.PaymentIntentID, nil)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if value > pi.Amount || value == 0 {
		value = pi.Amount
		params.Amount = pi.Amount
	}

	var captureParams = &stripe.PaymentIntentCaptureParams{
		AmountToCapture: stripe.Int64(value),
	}

	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded, stripe.PaymentIntentStatusCanceled:
		return pi, nil
	}

	pi, err = paymentintent.Capture(pi.ID, captureParams)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return pi, nil
}

type CancelOrRefundPaymentIntentParams struct {
	PaymentIntentID string
	Amount          int64
	OrderID         string
}

func (client *StripeClient) CancelOrRefundPaymentIntent(params CancelOrRefundPaymentIntentParams) (*stripe.PaymentIntent, *stripe.Refund, error) {
	pi, err := paymentintent.Get(params.PaymentIntentID, nil)
	if err != nil {
		return nil, nil, err
	}

	if pi.Status == stripe.PaymentIntentStatusSucceeded {
		refundParams := &stripe.RefundParams{
			Amount:        stripe.Int64(pi.Amount),
			PaymentIntent: stripe.String(params.PaymentIntentID),
		}

		if params.Amount > 0 && params.Amount < pi.Amount {
			refundParams.Amount = stripe.Int64(params.Amount)
		}
		if params.OrderID != "" {
			refundParams.AddMetadata("order_id", params.OrderID)
		}

		ref, err := refund.New(refundParams)
		if err != nil {
			return nil, nil, eris.Wrap(err, "")
		}
		return pi, ref, nil
	}

	if pi.Status == stripe.PaymentIntentStatusCanceled {
		return pi, nil, errs.ErrPaymentRefunded
	}

	intent, err := paymentintent.Cancel(params.PaymentIntentID, nil)
	if err != nil {
		return nil, nil, eris.Wrap(err, "")
	}

	return intent, nil, nil

}

type RefundPaymentIntentParams struct {
	PaymentIntentID string
	Amount          int64
	Metadata        map[string]string
}

func (client *StripeClient) RefundPaymentIntent(params RefundPaymentIntentParams) (*stripe.Refund, error) {
	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(params.PaymentIntentID),
	}
	if params.Amount > 0 {
		refundParams.Amount = stripe.Int64(params.Amount)
	}

	for k, v := range params.Metadata {
		refundParams.AddMetadata(k, v)
	}

	ref, err := refund.New(refundParams)

	return ref, err

}

type CancelPaymentIntentParams struct {
	PaymentIntentID string
}

func (client *StripeClient) CancelPaymentIntent(params CancelPaymentIntentParams) (*stripe.PaymentIntent, error) {
	return paymentintent.Cancel(params.PaymentIntentID, nil)

}
func (client *StripeClient) GetRefundOfOrder(paymentIntentID string, orderID string) (*stripe.Refund, error) {
	var params = &stripe.RefundListParams{
		PaymentIntent: stripe.String(paymentIntentID),
	}

	i := refund.List(params)
	var foundRefundedID string
	for i.Next() {
		r := i.Refund()
		if id, ok := r.Metadata["order_id"]; ok {
			if id == orderID {
				foundRefundedID = r.ID
				break
			}
		}
	}

	if foundRefundedID == "" {
		return nil, eris.New("Refund not found")
	}

	return refund.Get(foundRefundedID, nil)

}

func (client *StripeClient) GetPaymentIntent(piID string) (*stripe.PaymentIntent, error) {
	var params = &stripe.PaymentIntentParams{}
	params.AddExpand("invoice")
	params.AddExpand("latest_charge.balance_transaction")
	return paymentintent.Get(piID, params)
}

type SetupIntentForBankAccountParams struct {
	StripeCustomerID string `json:"stripe_customer_id" query:"stripe_customer_id" form:"stripe_customer_id"`
}

func (client *StripeClient) SetupIntentForClientSecert(form SetupIntentForBankAccountParams) (*stripe.SetupIntent, error) {
	params := &stripe.SetupIntentParams{
		Customer: stripe.String(form.StripeCustomerID),
	}
	result, err := setupintent.New(params)

	if err != nil {
		return nil, err
	}

	return result, nil
}

type ConfirmPaymentIntentParams struct {
	PaymentIntentID string `json:"payment_intent_id"`
	ReturnURL       string `json:"return_url"`
}

func (client *StripeClient) ConfirmPaymentIntent(params ConfirmPaymentIntentParams) (*stripe.PaymentIntent, error) {
	var confirmParams = &stripe.PaymentIntentConfirmParams{
		ReturnURL: stripe.String(params.ReturnURL),
	}
	result, err := paymentintent.Confirm(params.PaymentIntentID, confirmParams)
	return result, err
}
