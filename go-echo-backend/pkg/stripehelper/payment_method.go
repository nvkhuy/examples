package stripehelper

import (
	"errors"

	"github.com/rotisserie/eris"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentmethod"
)

// AddPaymentMethod add pm
func (client *StripeClient) AddPaymentMethod(customerID, paymentMethodID string) (*stripe.PaymentMethod, error) {
	var params = &stripe.PaymentMethodAttachParams{
		Customer: &customerID,
	}

	pm, err := paymentmethod.Attach(paymentMethodID, params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return pm, nil
}

// HasPaymentMethod check payment method by customer id
func (client *StripeClient) HasPaymentMethod(customerID string) bool {
	if customerID == "" {
		return false
	}

	var params = &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String("card"),
	}
	params.Limit = stripe.Int64(100)
	var iter = paymentmethod.List(params)

	var hasPaymentMethod = false
	for iter.Next() {
		hasPaymentMethod = true
		break

	}

	return hasPaymentMethod
}

func (client *StripeClient) GetPaymentMethod(pmID string) (*stripe.PaymentMethod, error) {
	if pmID == "" {
		return nil, errors.New("Payment method not found")
	}

	pm, err := paymentmethod.Get(pmID, nil)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return pm, nil
}

// DetachPaymentMethod get payment method by id
func (client *StripeClient) DetachPaymentMethod(pmID string) (*stripe.PaymentMethod, error) {
	if pmID == "" {
		return nil, errors.New("Payment method not found")
	}

	pm, err := paymentmethod.Detach(pmID, nil)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return pm, nil
}

// UpdatePaymentMethodParams params
type UpdatePaymentMethodParams struct {
	PaymentMethodID string
	CustomerID      string
	IsDefault       *bool
	BillingDetails  *stripe.PaymentMethodBillingDetailsParams
}

func (client *StripeClient) UpdatePaymentMethod(params UpdatePaymentMethodParams) (*stripe.PaymentMethod, error) {
	if params.PaymentMethodID == "" {
		return nil, errors.New("Payment method not found")
	}
	var p = &stripe.PaymentMethodParams{
		BillingDetails: params.BillingDetails,
	}

	pm, err := paymentmethod.Update(params.PaymentMethodID, p)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if params.IsDefault != nil && *params.IsDefault {
		_, err = customer.Update(
			params.CustomerID,
			&stripe.CustomerParams{
				InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
					DefaultPaymentMethod: stripe.String(params.PaymentMethodID),
				},
			},
		)
	}
	return pm, nil
}

func (client *StripeClient) GetPaymentMethods(customerID string) []*PaymentMethodWithDefault {
	var list = []*PaymentMethodWithDefault{}

	if customerID == "" {
		return list
	}

	var params = &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String("card"),
	}
	params.AddExpand("data.customer")
	params.Limit = stripe.Int64(100)

	var iter = paymentmethod.List(params)
	for iter.Next() {

		var p = iter.PaymentMethod()
		var pm = &PaymentMethodWithDefault{
			PaymentMethod: p,
			IsDefault:     false,
		}
		if p.Customer != nil && p.Customer.InvoiceSettings != nil && p.Customer.InvoiceSettings.DefaultPaymentMethod != nil {
			pm.IsDefault = p.Customer.InvoiceSettings.DefaultPaymentMethod.ID == p.ID
		}

		list = append(list, pm)
	}

	return list
}

func (client *StripeClient) SetAsDefaultPaymentMethod(customerID, paymentMethodID string) (*stripe.PaymentMethod, error) {
	var params = &stripe.PaymentMethodAttachParams{
		Customer: &customerID,
	}

	pm, err := paymentmethod.Attach(paymentMethodID, params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	_, err = customer.Update(
		customerID,
		&stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(paymentMethodID),
			},
		},
	)

	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return pm, nil
}

// GetDefaultPaymentMethod get default payment method
func (client *StripeClient) GetDefaultPaymentMethod(customerID string) (*stripe.PaymentMethod, error) {
	cus, err := customer.Get(customerID, nil)
	if err != nil {
		return nil, err
	}

	if cus.InvoiceSettings == nil {
		return nil, ErrNoDefaultPaymentMethod
	}

	if cus.InvoiceSettings.DefaultPaymentMethod == nil {
		return nil, ErrNoDefaultPaymentMethod
	}

	pm, err := paymentmethod.Get(cus.InvoiceSettings.DefaultPaymentMethod.ID, nil)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return pm, nil
}
