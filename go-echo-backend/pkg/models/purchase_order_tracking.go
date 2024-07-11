package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

// Faq Faq's model
type PurchaseOrderTracking struct {
	Model

	Attachments *Attachments `json:"attachments,omitempty"`
	Notes       string       `json:"notes,omitempty"`

	ActionType       enums.PoTrackingAction    `json:"action_type"`
	Description      string                    `json:"description,omitempty"`
	FromStatus       enums.PoTrackingStatus    `json:"from_status,omitempty"`
	ToStatus         enums.PoTrackingStatus    `json:"to_status,omitempty"`
	ParentTrackingID string                    `json:"parent_tracking_id,omitempty"`
	PurchaseOrderID  string                    `json:"purchase_order_id,omitempty"`
	UserID           string                    `json:"user_id,omitempty"`
	CreatedByUserID  string                    `json:"created_by_user_id,omitempty"`
	UserGroup        enums.PoTrackingUserGroup `gorm:"default:'buyer'" json:"user_group,omitempty"`

	Metadata *PoTrackingMetadata `json:"metadata,omitempty"`
}

type PurchaseOrderTrackingCreateForm struct {
	Attachments      *Attachments              `json:"attachments,omitempty"`
	ActionType       enums.PoTrackingAction    `json:"action_type"`
	Description      string                    `json:"description,omitempty"`
	FromStatus       enums.PoTrackingStatus    `json:"from_status,omitempty"`
	ToStatus         enums.PoTrackingStatus    `json:"to_status,omitempty"`
	ParentTrackingID string                    `json:"parent_tracking_id,omitempty"`
	UserID           string                    `json:"user_id,omitempty"`
	CreatedByUserID  string                    `json:"created_by_user_id,omitempty"`
	PurchaseOrderID  string                    `json:"purchase_order_id,omitempty"`
	UserGroup        enums.PoTrackingUserGroup `gorm:"default:'buyer'" json:"user_group,omitempty"`

	Metadata *PoTrackingMetadata `json:"metadata,omitempty"`
}
