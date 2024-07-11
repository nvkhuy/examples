package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/stripe/stripe-go/v74"
)

type OrderCartItems []*OrderCartItem

func (items OrderCartItems) IDs() []string {
	var IDs []string
	for _, item := range items {
		IDs = append(IDs, item.ID)
	}
	return IDs
}

type OrderCartItem struct {
	Model

	PurchaseOrderID     string `json:"purchase_order_id"`
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id"`
	CheckoutSessionID   string `json:"checkout_session_id"`

	Sku       string `json:"sku,omitempty"`
	Style     string `json:"style,omitempty"`
	ColorName string `json:"color_name,omitempty"`
	Size      string `json:"size,omitempty"`

	Qty        int64       `gorm:"default:1" json:"qty"`
	UnitPrice  price.Price `gorm:"type:decimal(20,4);default:0.0" json:"unit_price"`
	TotalPrice price.Price `gorm:"type:decimal(20,4);default:0.0" json:"total_price"`

	WaitingForCheckout *bool `gorm:"default:true" json:"waiting_for_checkout"`

	NoteToSupplier string `json:"note_to_supplier"`
}

type GetOrderCartRequest struct {
	JwtClaimsInfo
	IsFromCatalog bool `json:"is_from_catalog" query:"is_from_catalog"`
}

type GetOrderCartResponse struct {
	PurchaseOrders     []*PurchaseOrder     `json:"purchase_orders"`
	BulkPurchaseOrders []*BulkPurchaseOrder `json:"bulk_purchase_orders"`
}

type OrderCartPreviewCheckoutRequest struct {
	JwtClaimsInfo
	PurchaseOrderCartItemIDs []string          `json:"purchase_order_cart_item_ids"`
	BulkOrderIDs             []string          `json:"bulk_order_ids"`
	PaymentType              enums.PaymentType `json:"payment_type" validate:"required"`
}

type OrderCartCheckoutRequest struct {
	JwtClaimsInfo
	PurchaseOrderCartItemIDs []string `json:"purchase_order_cart_item_ids"`
	BulkOrderIDs             []string `json:"bulk_order_ids"`

	PaymentType     enums.PaymentType `json:"payment_type" validate:"oneof=bank_transfer card"`
	PaymentMethodID string            `json:"payment_method_id" validate:"required_if=PaymentType card"`

	TransactionRefID      string      `json:"transaction_ref_id" validate:"required_if=PaymentType bank_transfer"`
	TransactionAttachment *Attachment `json:"transaction_attachment" validate:"required_if=PaymentType bank_transfer"`
}

type OrderCartCheckoutResponse struct {
	GetOrderCartResponse

	CheckoutSessionID         string                          `json:"checkout_session_id"`
	PaymentTransaction        *PaymentTransaction             `json:"payment_transaction"`
	PaymentIntentNextAction   *stripe.PaymentIntentNextAction `json:"payment_intent_next_action,omitempty"`
	PaymentIntentClientSecret string                          `json:"payment_intent_client_secret,omitempty"`
}

type StripeConfirmOrderCartCheckoutParams struct {
	CheckoutSessionID         string `json:"checkout_session_id" param:"checkout_session_id" query:"checkout_session_id" validate:"required"`
	PaymentIntent             string `json:"payment_intent"  param:"payment_intent" query:"payment_intent" validate:"required"`
	PaymentIntentClientSecret string `json:"payment_intent_client_secret"  param:"payment_intent_client_secret" query:"payment_intent_client_secret"`
	SourceType                string `json:"source_type" param:"source_type" query:"source_type"`
	PurchaseOrderCartItemIDs  string `json:"purchase_order_cart_item_ids" param:"po_items" query:"po_items"`
	BulkOrderIDs              string `json:"bulk_order_ids" param:"bulks" query:"bulks"`
}

type OrderCartGetCheckoutInfoRequest struct {
	JwtClaimsInfo
	CheckoutSessionID string `json:"checkout_session_id" query:"checkout_session_id" validate:"required"`
}

type OrderCartGetCheckoutInfoResponse struct {
	GetOrderCartResponse
	PaymentTransaction PaymentTransaction `json:"payment_transaction"`
}
type CreateBuyerPaymentLinkRequest struct {
	JwtClaimsInfo
	BuyerID          string   `json:"user_id" param:"buyer_id" validate:"required"`
	PurchaseOrderIDs []string `json:"purchase_order_ids"`
	BulkIDs          []string `json:"bulk_purchase_order_ids"`
}
type CreateBuyerPaymentLinkResponse struct {
	PaymentLink string `json:"payment_link"`
}

type GetBuyerOrderCartRequest struct {
	JwtClaimsInfo
	BuyerID    string   `json:"buyer_id" param:"buyer_id" validate:"required"`
	InquiryIDs []string `json:"inquiry_ids"`
	BulkIDs    []string `json:"bulk_purchase_order_ids"`
}
type GetBuyerOrderCartResponse struct {
	GetOrderCartResponse
}
