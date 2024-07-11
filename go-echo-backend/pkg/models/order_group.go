package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderGroup struct {
	Model
	Name   string `gorm:"size:1000" json:"name,omitempty"`
	UserID string `gorm:"index;size:100" json:"user_id"`
	User   *User  `gorm:"-" json:"user,omitempty"`

	Inquiries []*Inquiry           `gorm:"-" json:"inquiries,omitempty"`
	Samples   []*PurchaseOrder     `gorm:"-" json:"samples,omitempty"`
	Bulks     []*BulkPurchaseOrder `gorm:"-" json:"bulks,omitempty"`
}

type OrderGroupAlias struct {
	OrderGroup
	InquiriesJson datatypes.JSON `gorm:"column:inquiries_json"`
	SamplesJson   datatypes.JSON `gorm:"column:samples_json"`
	BulksJson     datatypes.JSON `gorm:"column:bulks_json"`
}

type OrderGroups []*OrderGroup

func (og OrderGroups) IDs() []string {
	var orderGroupIDs []string
	for _, o := range og {
		orderGroupIDs = append(orderGroupIDs, o.ID)
	}
	return orderGroupIDs
}

func (d *OrderGroup) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = helper.GenerateOrderGroupID()
	}
	return
}

type CreateOrderGroupRequest struct {
	JwtClaimsInfo
	Name   string `json:"name" validate:"required"`
	UserID string `json:"user_id"`
}

type GetOrderGroupListRequest struct {
	JwtClaimsInfo
	PaginationParams
	OrderGroupStatus enums.OrderGroupStatus `json:"order_group_status" query:"order_group_status"`
	UserID           string                 `json:"user_id" query:"user_id"`
}
type GetOrderGroupDetailRequest struct {
	JwtClaimsInfo
	OrderGroupID string `json:"order_group_id" param:"order_group_id" validate:"required"`
}

type AssignOrderGroupRequest struct {
	JwtClaimsInfo
	OrderType    enums.OrderGroupType `json:"order_type" validate:"required"`
	OrderIDs     []string             `json:"order_ids" validate:"required"`
	OrderGroupID string               `json:"order_group_id" validate:"required"`
}

type CollectionPreviewCheckoutRequest struct {
	JwtClaimsInfo
	OrderGroupID string `json:"order_group_id" validate:"required"`
}
