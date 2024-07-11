package stripehelper

import (
	"errors"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
)

// CustomerBillingAddress billing address
type CustomerBillingAddress struct {
	Line1      string
	City       string
	State      string
	Country    string
	PostalCode string
}

// CreateCustomerParams params
type CreateCustomerParams struct {
	PaymentMethodID *string
	Name            string
	Email           string
	BillingAddress  *CustomerBillingAddress
	Metadata        map[string]string
}

// CreateCustomer create customer
func (client *StripeClient) CreateCustomer(params *CreateCustomerParams) (string, error) {
	var desc = fmt.Sprintf("%s's customer.", client.config.GetServerName(config.ServiceBackend))
	var customerParams = &stripe.CustomerParams{
		Name:        stripe.String(params.Name),
		Email:       stripe.String(params.Email),
		Description: stripe.String(desc),
		Params: stripe.Params{
			Metadata: params.Metadata,
		},
	}

	if params.PaymentMethodID != nil && *params.PaymentMethodID != "" {
		customerParams.PaymentMethod = params.PaymentMethodID
		customerParams.InvoiceSettings = &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: params.PaymentMethodID,
		}
	}
	if params.BillingAddress != nil {

		if params.BillingAddress.Line1 != "" {
			customerParams.Address = &stripe.AddressParams{
				Line1:      &params.BillingAddress.Line1,
				City:       &params.BillingAddress.City,
				State:      &params.BillingAddress.State,
				Country:    &params.BillingAddress.Country,
				PostalCode: &params.BillingAddress.PostalCode,
			}
		}

	}

	c, err := customer.New(customerParams)
	if err != nil {
		return "", err
	}

	if c.ID == "" {
		return "", errors.New("customer ID is empty")
	}

	return c.ID, nil
}
