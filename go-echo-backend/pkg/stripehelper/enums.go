package stripehelper

// Account Type
type AccountType string

var (
	AccountTypeStandard AccountType = "standard"
	AccountTypeExpress  AccountType = "express"
	AccountTypeCustom   AccountType = "custom"
)

func (a AccountType) String() string {
	return string(a)
}

// Payment Type
type PaymentMethodType string

var (
	PaymentMethodTypeCard          PaymentMethodType = "card"
	PaymentMethodTypeUsBankAccount PaymentMethodType = "us_bank_account"
)

func (a PaymentMethodType) String() string {
	return string(a)
}

// PayoutEvent
type PayoutEvent string

var (
	PaymentIntentProcessing     PayoutEvent = "payment_intent.processing"
	PaymentIntentRequriesAction PayoutEvent = "payment_intent.requires_action"
	PaymentIntentFailed         PayoutEvent = "payment_intent.payment_failed"
	PaymentIntentSucceeded      PayoutEvent = "payment_intent.succeeded"
	CheckoutSessionCompleted    PayoutEvent = "checkout.session.completed"
)
