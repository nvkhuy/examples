package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/lib/pq"
)

type PaymentTransactions []*PaymentTransaction

type PaymentTransaction struct {
	Model

	ReferenceID string            `gorm:"unique" json:"reference_id,omitempty"`
	PaymentType enums.PaymentType `gorm:"size:50;not null;default:'bank_transfer'" json:"payment_type,omitempty"`

	PurchaseOrder   *PurchaseOrder `gorm:"-" json:"purchase_order,omitempty"`
	PurchaseOrderID string         `json:"purchase_order_id,omitempty"`

	BulkPurchaseOrder   *BulkPurchaseOrder `gorm:"-" json:"bulk_purchase_order,omitempty"`
	BulkPurchaseOrderID string             `json:"bulk_purchase_order_id,omitempty"`

	CheckoutSessionID    string               `json:"checkout_session_id,omitempty"`
	PurchaseOrderIDs     pq.StringArray       `gorm:"type:varchar(100)[]" json:"purchase_order_ids,omitempty"`
	PurchaseOrders       []*PurchaseOrder     `gorm:"-" json:"purchase_orders,omitempty"`
	BulkPurchaseOrderIDs pq.StringArray       `gorm:"type:varchar(100)[]" json:"bulk_purchase_order_ids,omitempty"`
	BulkPurchaseOrders   []*BulkPurchaseOrder `gorm:"-" json:"bulk_purchase_orders,omitempty"`

	Currency enums.Currency `gorm:"default:'USD'" json:"currency,omitempty"`

	PaidAmount        *price.Price `json:"paid_amount,omitempty"`
	TotalAmount       *price.Price `json:"total_amount,omitempty"`
	PaymentPercentage *float64     `gorm:"type:decimal(20,4);default:100.0" json:"payment_percentage,omitempty"`
	BalanceAmount     *price.Price `json:"balance_amount,omitempty"`

	Attachments       *Attachments           `json:"attachments,omitempty"`
	Remark            string                 `json:"remark,omitempty"`
	PaymentCompleted  *bool                  `gorm:"-" json:"payment_completed,omitempty"`
	Milestone         enums.PaymentMilestone `json:"milestone,omitempty"` // First payment, final payment
	UserID            string                 `json:"user_id,omitempty"`
	BankTransferInfos BankTransferInfos      `gorm:"-" json:"bank_transfer_infos,omitempty"`
	TransactionRefID  string                 `json:"transaction_ref_id,omitempty"`
	User              *User                  `gorm:"-" json:"user,omitempty"`
	Note              string                 `json:"note,omitempty"`
	Status            enums.PaymentStatus    `gorm:"size:50;default:'pending'" json:"status,omitempty"`

	// Stripe
	PaymentIntentID string `gorm:"size:100" json:"payment_intent_id,omitempty"`
	ChargeID        string `gorm:"size:100" json:"charge_id,omitempty"`
	ReceiptURL      string `json:"receipt_url,omitempty"`
	TxnID           string `gorm:"size:100" json:"txn_id,omitempty"`

	MarkAsPaidAt   *int64 `json:"mark_as_paid_at,omitempty"`
	MarkAsUnpaidAt *int64 `json:"mark_as_unpaid_at,omitempty"`

	PaymentLinkID string `json:"payment_link_id"`

	TransactionType enums.TransactionType `gorm:"default:'credit'" json:"transaction_type"`

	Metadata *PaymentTransactionMetadata `json:"metadata,omitempty"`

	InvoiceNumber int      `json:"invoice_number"`
	Invoice       *Invoice `gorm:"-" json:"invoice,omitempty"`

	RefundReason string `json:"refund_reason,omitempty"`

	Fee *price.Price `gorm:"type:decimal(20,4);default:0.0"`
	Net *price.Price `gorm:"type:decimal(20,4);default:0.0"`
}
