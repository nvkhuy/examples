package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
)

type TNA struct { // Time And Action
	Model
	ReferenceID  string            `gorm:"not null" json:"reference_id"`
	Title        string            `json:"title"`
	SubTitle     string            `json:"sub_title"`
	Comment      string            `json:"comment"`
	DateFrom     int64             `json:"date_from"`
	DateTo       int64             `json:"date_to"`
	OrderType    enums.InquiryType `json:"order_type"`
	AssigneeIDs  pq.StringArray    `gorm:"type:varchar(200)[]" json:"assignee_ids"`
	Dependencies pq.StringArray    `gorm:"type:varchar(200)[]" json:"dependencies"`
	Assignees    []User            `gorm:"-" json:"assignees,omitempty"`
	UserID       string            `gorm:"-" json:"user_id,omitempty"`

	Inquiry           *Inquiry           `gorm:"-" json:"inquiry,omitempty"`
	PurchaseOrder     *PurchaseOrder     `gorm:"-" json:"purchase_order,omitempty"`
	BulkPurchaseOrder *BulkPurchaseOrder `gorm:"-" json:"bulk_purchase_order,omitempty"`
}
