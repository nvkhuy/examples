package payos

import "time"

type APIResponse[K any] struct {
	Code      string `json:"code,omitempty"`
	Desc      string `json:"desc,omitempty"`
	Data      K      `json:"data,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type PaymentLinkItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type CreatePaymentLinkRequest struct {
	OrderCode    int                `json:"orderCode"`
	Amount       int                `json:"amount"`
	Description  string             `json:"description"`
	BuyerName    string             `json:"buyerName"`
	BuyerEmail   string             `json:"buyerEmail"`
	BuyerPhone   string             `json:"buyerPhone"`
	BuyerAddress string             `json:"buyerAddress"`
	Items        []*PaymentLinkItem `json:"items"`
	CancelURL    string             `json:"cancelUrl"`
	ReturnURL    string             `json:"returnUrl"`
	ExpiredAt    int                `json:"expiredAt"`
	Signature    string             `json:"signature"`
}

type CreatePaymentLinkResponse struct {
	Bin           string `json:"bin"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	Amount        int    `json:"amount"`
	Description   string `json:"description"`
	OrderCode     int    `json:"orderCode"`
	PaymentLinkID string `json:"paymentLinkId"`
	Status        string `json:"status"`
	CheckoutURL   string `json:"checkoutUrl"`
	QrCode        string `json:"qrCode"`
}

type PaymentLinkTransaction struct {
	Reference              string    `json:"reference"`
	Amount                 int       `json:"amount"`
	AccountNumber          string    `json:"accountNumber"`
	Description            string    `json:"description"`
	TransactionDateTime    time.Time `json:"transactionDateTime"`
	VirtualAccountName     string    `json:"virtualAccountName"`
	VirtualAccountNumber   any       `json:"virtualAccountNumber"`
	CounterAccountBankID   any       `json:"counterAccountBankId"`
	CounterAccountBankName any       `json:"counterAccountBankName"`
	CounterAccountName     any       `json:"counterAccountName"`
	CounterAccountNumber   any       `json:"counterAccountNumber"`
}
type GetPaymentLinkResponse struct {
	ID                 string                    `json:"id"`
	OrderCode          int                       `json:"orderCode"`
	Amount             int                       `json:"amount"`
	AmountPaid         int                       `json:"amountPaid"`
	AmountRemaining    int                       `json:"amountRemaining"`
	Status             string                    `json:"status"`
	CreatedAt          time.Time                 `json:"createdAt"`
	Transactions       []*PaymentLinkTransaction `json:"transactions"`
	CancellationReason any                       `json:"cancellationReason"`
	CanceledAt         any                       `json:"canceledAt"`
}
