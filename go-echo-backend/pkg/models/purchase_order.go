package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
	"github.com/stripe/stripe-go/v74"
)

type PurchaseOrder struct {
	ID        string     `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	ReferenceID string `gorm:"unique" json:"reference_id"`
	ProductName string `gorm:"size:1000" json:"product_name"`

	ClientReferenceID *string `gorm:"unique;default:null" json:"client_reference_id"`

	Items []*PurchaseOrderItem `gorm:"-" json:"items,omitempty"` // from upload excel

	CheckoutSessionID string             `json:"checkout_session_id"`
	CartItemIDs       pq.StringArray     `gorm:"type:varchar(200)[]" json:"cart_item_ids,omitempty" swaggertype:"array,string"` //legacy
	CartItems         []*InquiryCartItem `gorm:"-" json:"cart_items,omitempty"`                                                 //legacy
	OrderCartItems    []*OrderCartItem   `gorm:"-" json:"order_cart_items,omitempty"`

	User   *User  `gorm:"-" json:"user,omitempty"`
	UserID string `json:"user_id,omitempty"`

	InquiryID string   `gorm:"size:200" json:"inquiry_id,omitempty"`
	Inquiry   *Inquiry `gorm:"-" json:"inquiry,omitempty"`

	Status enums.PurchaseOrderStatus `gorm:"size:50;default:'pending'" json:"status,omitempty"`

	PaymentIntentID string `gorm:"size:100" json:"payment_intent_id,omitempty"`
	ChargeID        string `gorm:"size:100" json:"charge_id,omitempty"`
	TxnID           string `gorm:"size:100" json:"txn_id,omitempty"`

	ReceiptURL string `json:"receipt_url,omitempty"`

	PaymentType enums.PaymentType `gorm:"size:50;not null;default:'bank_transfer'" json:"payment_type,omitempty"`

	TransferedAt   *int64 `json:"transfered_at,omitempty"`
	MarkAsPaidAt   *int64 `json:"mark_as_paid_at,omitempty"`
	MarkAsUnpaidAt *int64 `json:"mark_as_unpaid_at,omitempty"`

	TransactionRefID      string      `json:"transaction_ref_id,omitempty"`
	TransactionAttachment *Attachment `json:"transaction_attachment,omitempty"`

	PayoutTransactionRefID      string      `json:"payout_transaction_ref_id,omitempty"`
	PayoutTransactionAttachment *Attachment `json:"payout_transaction_attachment,omitempty"`
	PayoutTransferedAt          *int64      `json:"payout_transfered_at,omitempty"`
	PayoutMarkAsPaidAt          *int64      `json:"payout_mark_as_paid_at,omitempty"`
	PayoutMarkAsReceivedAt      *int64      `json:"payout_mark_as_received_at,omitempty"`

	Attachments         *Attachments `json:"attachments,omitempty"`
	Document            *Attachments `json:"document,omitempty"`
	Design              *Attachments `json:"design,omitempty"`
	ApproveDesignAt     *int64       `json:"approve_design_at,omitempty"`
	FabricAttachments   *Attachments `json:"fabric_attachments,omitempty"`
	TechpackAttachments *Attachments `json:"techpack_attachments,omitempty"`

	SampleAttachments    *Attachments              `json:"sample_attachments,omitempty"`
	ApproveRejectMeta    *InquiryApproveRejectMeta `json:"approve_reject_meta,omitempty"`
	TrackingStatus       enums.PoTrackingStatus    `gorm:"default:'new'" json:"tracking_status,omitempty"`
	PoRawMaterials       *PoRawMaterialMetas       `json:"po_raw_materials,omitempty"`
	ApproveRawMaterialAt *int64                    `json:"approve_raw_material_at,omitempty"`

	SellerPoRawMaterials *PoRawMaterialMetas `json:"seller_po_raw_materials,omitempty"`

	LogisticInfo *PoLogisticMeta      `json:"logistic_info,omitempty"`
	MakingInfo   *PoMarkingStatusMeta `json:"making_info,omitempty"`
	SubmitInfo   *PoMarkingStatusMeta `json:"submit_info,omitempty"`

	ReceiverConfirmedAt *int64 `json:"receiver_confirmed_at,omitempty"`
	DeliveryStartedAt   *int64 `json:"delivery_started_at,omitempty"`

	AssigneeIDs pq.StringArray `gorm:"type:varchar(200)[]" json:"assignee_ids,omitempty"`
	Assignees   Users          `gorm:"-" json:"assignees,omitempty"`

	PaymentIntentNextAction   *stripe.PaymentIntentNextAction `gorm:"-" json:"payment_intent_next_action,omitempty"`
	PaymentIntentClientSecret string                          `gorm:"-" json:"payment_intent_client_secret,omitempty"`

	PaymentTransaction *PaymentTransaction `gorm:"-" json:"payment_transaction,omitempty"`

	Currency enums.Currency         `gorm:"default:'USD'" json:"currency,omitempty"`
	Feedback *PurchaseOrderFeedback `json:"feedback,omitempty"`

	// Seller or Inflow
	SampleMakerID string `gorm:"default:'inflow'" json:"sample_maker_id,omitempty"`
	SampleMaker   *User  `gorm:"-" json:"sample_maker,omitempty"`

	// Seller info
	SellerPoAttachments       *PoAttachments               `json:"seller_po_attachments,omitempty"`
	SellerTrackingStatus      enums.SellerPoTrackingStatus `gorm:"default:'new'" json:"seller_tracking_status,omitempty"`
	SellerDesign              *SellerPoDesignMeta          `json:"seller_design,omitempty"`
	SellerTechpackAttachments *Attachments                 `json:"seller_techpack_attachments,omitempty"` // Final design got clone from buyer TechpackAttachments
	SellerEstMakingAt         *int64                       `json:"seller_est_making_at,omitempty"`
	SellerEstDeliveryAt       *int64                       `json:"seller_est_delivery_at,omitempty"`
	SellerDesignApprovedAt    *int64                       `json:"seller_design_approved_at,omitempty"`
	SellerSubmitInfo          *PoMarkingStatusMeta         `json:"seller_submit_info,omitempty"`
	SellerLogisticInfo        *PoLogisticMeta              `json:"seller_logistic_info,omitempty"`

	SellerDeliveryStartedAt   *int64 `json:"seller_delivery_started_at,omitempty"`
	SellerDeliveryConfirmedAt *int64 `json:"seller_delivery_confirmed_at,omitempty"`
	SellerDeliveryFeedback    string `json:"seller_delivery_feedback,omitempty"`

	SellerPORejectReason string `json:"seller_po_reject_reason,omitempty"`

	// Extra seller info for admin api
	InquirySeller *InquirySeller `gorm:"-" json:"inquiry_seller,omitempty"`

	InvoiceNumber int      `json:"invoice_number,omitempty"`
	Invoice       *Invoice `gorm:"-" json:"invoice,omitempty"`

	Trackings []*PurchaseOrderTracking `gorm:"-" json:"trackings,omitempty"`

	PaymentTransactionReferenceID string `json:"payment_transaction_reference_id,omitempty"`

	RefundReason string `json:"refund_reason,omitempty"`

	ProductWeight     *float64              `json:"product_weight,omitempty"`
	ShippingAddressID string                `gorm:"size:100" json:"shipping_address_id,omitempty"`
	ShippingAddress   *Address              `gorm:"-" json:"shipping_address,omitempty"`
	Quotations        InquiryQuotationItems `json:"quotations,omitempty"`

	IsCart      *bool `gorm:"default:false" json:"-"`
	FromCatalog bool  `gorm:"default:false" json:"from_catalog"`

	ConfirmedAt *int64 `json:"confirmed_at"`

	PaymentLink   string `json:"payment_link,omitempty"`
	PaymentLinkID string `json:"payment_link_id,omitempty"`

	StartDate      *int64 `json:"start_date,omitempty"`
	LeadTime       int    `json:"lead_time,omitempty"`
	CompletionDate *int64 `json:"completion_date,omitempty"`

	HubspotDealID string `gorm:"size:100"  json:"hubspot_deal_id,omitempty"`

	OrderGroupID string      `gorm:"size:100"  json:"group_id,omitempty"`
	OrderGroup   *OrderGroup `gorm:"-" json:"order_group,omitempty"`

	DesignNote string `gorm:"size:3000" json:"design_note"`

	RoundID     *string           `gorm:"unique;size:100;default:null" json:"round_id,omitempty"`
	RoundStatus enums.RoundStatus `json:"round_status,omitempty"`
	Pricing
	SellerPricing
}

type PurchaseOrders []*PurchaseOrder

type CheckoutForm struct {
	PaymentType enums.PaymentType `json:"payment_type" validate:"required,oneof=card bank_transfer"`
}

type PurchaseOrderAssignPICParam struct {
	JwtClaimsInfo
	PurchaseOrderID string   `json:"purchase_order_id" query:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	AssigneeIDs     []string `json:"assignee_ids" query:"assignee_ids" param:"assignee_ids" validate:"required"`
}

type PurchaseOrderPaymentIntentConfirmParams struct {
	InquiryID       string `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`
	PurchaseOrderID string `json:"purchase_order_id"  param:"purchase_order_id" query:"purchase_order_id" validate:"required"`

	PaymentIntent             string `json:"payment_intent"  param:"payment_intent" query:"payment_intent" validate:"required"`
	PaymentIntentClientSecret string `json:"payment_intent_client_secret"  param:"payment_intent_client_secret" query:"payment_intent_client_secret"`
	SourceType                string `json:"source_type" param:"source_type" query:"source_type"`
}

type BulkPurchaseOrderPaymentIntentConfirmParams struct {
	InquiryID                 string `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`
	BulkPurchaseOrderID       string `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	PaymentIntent             string `json:"payment_intent"  param:"payment_intent" query:"payment_intent" validate:"required"`
	PaymentIntentClientSecret string `json:"payment_intent_client_secret"  param:"payment_intent_client_secret" query:"payment_intent_client_secret"`
	SourceType                string `json:"source_type" param:"source_type" query:"source_type"`
}

type StripeConfirmInquiryCartsCheckoutParams struct {
	CheckoutSessionID         string `json:"checkout_session_id" param:"checkout_session_id" query:"checkout_session_id" validate:"required"`
	PaymentIntent             string `json:"payment_intent"  param:"payment_intent" query:"payment_intent" validate:"required"`
	PaymentIntentClientSecret string `json:"payment_intent_client_secret"  param:"payment_intent_client_secret" query:"payment_intent_client_secret"`
	SourceType                string `json:"source_type" param:"source_type" query:"source_type"`
	CartItems                 string `json:"cart_items" param:"cart_items" query:"cart_items"`
}

type PurchaseOrderIDParam struct {
	JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" query:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	Note            string `json:"note"`

	PurchaseOrder *PurchaseOrder `json:"-"`
}
