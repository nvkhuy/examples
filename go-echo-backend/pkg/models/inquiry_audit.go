package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

const InquiryAuditCurrentUserKey = "current_user"
const InquiryAuditCurrentInquiryKey = "current_inquiry"

type InquiryAudit struct {
	Model

	InquiryID           string `gorm:"size:200" json:"inquiry_id"`
	PurchaseOrderID     string `gorm:"size:200" json:"purchase_order_id,omitempty"`
	BulkPurchaseOrderID string `gorm:"size:200" json:"bulk_purchase_order_id,omitempty"`
	UserID              string `gorm:"size:200" json:"user_id"`

	ActionType  enums.AuditActionType `gorm:"size:200" json:"action_type"`
	Description string                `gorm:"size:2000" json:"description,omitempty"`
	Metadata    *InquiryAuditMetadata `json:"metadata,omitempty"`
	FromState   string                `gorm:"size:200" json:"from_state,omitempty"`
	ToState     string                `gorm:"size:200" json:"to_state,omitempty"`

	Notes       string      `gorm:"size:2000" json:"note,omitempty"`
	Attachments Attachments `json:"attachments,omitempty"`
}

type InquiryAuditCreateForm struct {
	InquiryID           string `json:"inquiry_id"`
	UserID              string `json:"user_id"`
	PurchaseOrderID     string `json:"purchase_order_id"`
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id,omitempty"`

	ActionType  enums.AuditActionType `json:"action_type"`
	Description string                `json:"description,omitempty"`
	Metadata    *InquiryAuditMetadata `json:"metadata,omitempty"`
	FromState   string                `json:"from_state,omitempty"`
	ToState     string                `json:"to_state,omitempty"`
}
