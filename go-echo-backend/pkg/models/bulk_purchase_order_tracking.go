package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

// Faq Faq's model
type BulkPurchaseOrderTracking struct {
	Model

	Attachments      *Attachments               `json:"attachments,omitempty"`
	Notes            string                     `json:"notes,omitempty"`
	ActionType       enums.BulkPoTrackingAction `json:"action_type"`
	Description      string                     `json:"description,omitempty"`
	FromStatus       enums.BulkPoTrackingStatus `json:"from_status,omitempty"`
	ToStatus         enums.BulkPoTrackingStatus `json:"to_status,omitempty"`
	ParentTrackingID string                     `json:"parent_tracking_id,omitempty"`
	PurchaseOrderID  string                     `json:"purchase_order_id,omitempty"`
	UserID           string                     `json:"user_id,omitempty"`
	CreatedByUserID  string                     `json:"created_by_user_id,omitempty"`
	ReportStatus     enums.QcReportResult       `json:"report_status,omitempty"`

	UserGroup enums.PoTrackingUserGroup `gorm:"default:'buyer'" json:"user_group,omitempty"`

	Metadata *PoTrackingMetadata `json:"metadata,omitempty"`
}

type BulkPurchaseOrderTrackingCreateForm struct {
	Attachments      *Attachments               `json:"attachments,omitempty"`
	ActionType       enums.BulkPoTrackingAction `json:"action_type"`
	Description      string                     `json:"description,omitempty"`
	FromStatus       enums.BulkPoTrackingStatus `json:"from_status,omitempty"`
	ToStatus         enums.BulkPoTrackingStatus `json:"to_status,omitempty"`
	ParentTrackingID string                     `json:"parent_tracking_id,omitempty"`
	UserID           string                     `json:"user_id,omitempty"`
	CreatedByUserID  string                     `json:"created_by_user_id,omitempty"`
	PurchaseOrderID  string                     `json:"purchase_order_id,omitempty"`
	ReportStatus     enums.QcReportResult       `json:"report_status,omitempty"`
	Metadata         *PoTrackingMetadata        `json:"metadata,omitempty"`
	UserGroup        enums.PoTrackingUserGroup  `gorm:"default:'buyer'" json:"user_group,omitempty"`
}
