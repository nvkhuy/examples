package stripehelper

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/balancetransaction"
)

type GetAllBalanceTransactionsParams struct {
	PayoutID string `json:"payout_id"`
	Type     string `json:"ytpe"`
}

func (client *StripeClient) GetAllBalanceTransactions(req *GetAllBalanceTransactionsParams) []*stripe.BalanceTransaction {
	var iter = balancetransaction.List(&stripe.BalanceTransactionListParams{
		Payout: stripe.String(req.PayoutID),
		Type:   stripe.String(req.Type),
	})

	return iter.BalanceTransactionList().Data

}

func (client *StripeClient) GetBalanceTransaction(id string) (*stripe.BalanceTransaction, error) {
	item, err := balancetransaction.Get(id, &stripe.BalanceTransactionParams{})

	return item, err

}
