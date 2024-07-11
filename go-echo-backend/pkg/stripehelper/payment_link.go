package stripehelper

import (
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/rotisserie/eris"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentlink"
)

// CreatePaymentLinkParams params
type CreatePaymentLinkParams struct {
	Currency    enums.Currency
	Metadata    map[string]string
	RedirectURL string
	LineItems   []*stripe.PaymentLinkLineItemParams
}

// CreatePaymentIntent payment intent
func (client *StripeClient) CreatePaymentLink(params CreatePaymentLinkParams) (*stripe.PaymentLink, error) {
	var linkParams = &stripe.PaymentLinkParams{
		Currency: stripe.String(string(params.Currency)),
		AfterCompletion: &stripe.PaymentLinkAfterCompletionParams{
			Type: stripe.String(string(stripe.PaymentLinkAfterCompletionTypeRedirect)),
			Redirect: &stripe.PaymentLinkAfterCompletionRedirectParams{
				URL: stripe.String(params.RedirectURL),
			},
		},
	}

	for key, value := range params.Metadata {
		linkParams.AddMetadata(key, value)
	}
	linkParams.Metadata = params.Metadata

	linkParams.LineItems = params.LineItems

	pl, err := paymentlink.New(linkParams)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			if stripeErr.Code == stripe.ErrorCodeBalanceInsufficient {
				return nil, eris.Wrap(errs.ErrPreAuthCardFailed, "")
			}
		}

		return nil, eris.Wrap(err, "")
	}

	if pl.ID == "" {
		return nil, eris.Wrap(errs.ErrStripePaymentLinkIDInvalid, "")
	}

	return pl, nil
}
