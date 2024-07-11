package stripehelper

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/payout"
)

type GetAllPayOutsParams struct {
}

func (client *StripeClient) GetAllPayOuts(req *GetAllPayOutsParams) []*stripe.Payout {
	var iter = payout.List(&stripe.PayoutListParams{})

	return iter.PayoutList().Data

}
