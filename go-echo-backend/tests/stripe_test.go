package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/stretchr/testify/assert"
)

func TestStripe_GetAllPayOuts(t *testing.T) {
	var config = initConfig()

	var result = stripehelper.New(config).GetAllPayOuts(&stripehelper.GetAllPayOutsParams{})

	helper.PrintJSON(result)

}

func TestStripe_GetAllBalanceTransactions(t *testing.T) {
	var config = initConfig()

	var result = stripehelper.New(config).GetAllBalanceTransactions(&stripehelper.GetAllBalanceTransactionsParams{
		PayoutID: "po_1OpHtWLr6GIPd0Z1IQAFizCN",
		Type:     "charge",
	})

	helper.PrintJSON(result)
}

func TestStripe_GetBalanceTransaction(t *testing.T) {
	var config = initConfig()

	result, err := stripehelper.New(config).GetBalanceTransaction("txn_3OmpX2Lr6GIPd0Z11p99m95w")
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestStripe_GetCharge(t *testing.T) {
	var config = initConfig()

	result, err := stripehelper.New(config).GetCharge("ch_3Npno2Lr6GIPd0Z11O8hvRBJ")
	assert.NoError(t, err)

	helper.PrintJSON(result)

}

func TestStripe_ConfirmPaymentIntent(t *testing.T) {
	var config = initConfig()

	result, err := stripehelper.New(config).ConfirmPaymentIntent(stripehelper.ConfirmPaymentIntentParams{
		PaymentIntentID: "pi_3O9PqHLr6GIPd0Z10EOeROlj",
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)

}
func TestStripe_GetPaymentIntent(t *testing.T) {
	var config = initConfig()

	result, err := stripehelper.New(config).GetPaymentIntent("pi_3Oz9n0Lr6GIPd0Z11B6fyjiY")
	assert.NoError(t, err)

	helper.PrintJSON(result)

}
