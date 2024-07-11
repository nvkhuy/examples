package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
)

type BulkPurchaseOrderSellerQuotation struct {
	Model
	ID        string `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	BulkPurchaseOrderID string             `gorm:"uniqueIndex:idx_bulk_purchase_order_seller_quotation" json:"bulk_purchase_order_id,omitempty"`
	BulkPurchaseOrder   *BulkPurchaseOrder `gorm:"-" json:"bulk_purchase_order"`

	UserID string `gorm:"uniqueIndex:idx_bulk_purchase_order_seller_quotation" json:"user_id"`
	User   *User  `gorm:"-" json:"user,omitempty"`

	DueDay *int64 `json:"due_day,omitempty"`

	Status enums.BulkPurchaseOrderSellerStatus `gorm:"default:'waiting_for_quotation'" json:"status"`

	DeliveryDate *int64 `json:"delivery_date,omitempty"`

	Currency enums.Currency `json:"currency,omitempty"`

	// Price from Admin
	OrderType   enums.OBOrderType `json:"order_type,omitempty"`
	OfferPrice  *price.Price      `json:"offer_price,omitempty"`
	OfferRemark string            `json:"offer_remark,omitempty"`

	VarianceAmount     *price.Price `json:"variance_amount,omitempty"`
	VariancePercentage *float64     `json:"variance_percentage,omitempty"`

	// Price from seller
	SellerRemark   string       `json:"seller_remark,omitempty"`
	FabricCost     *price.Price `json:"fabric_cost,omitempty"`
	DecorationCost *price.Price `json:"decoration_cost,omitempty"`
	MakingCost     *price.Price `json:"making_cost,omitempty"` // sewing, cut, making, finishing
	OtherCost      *price.Price `json:"other_cost,omitempty"`

	AdminSentAt    *int64                  `json:"admin_sent_at,omitempty"`
	QuotationAt    *int64                  `json:"quotation_at,omitempty"`
	BulkQuotations SellerBulkQuotationMOQs `json:"bulk_quotations,omitempty"`

	ExpectedStartProductionDate *int64 `json:"expected_start_production_date,omitempty"`
	StartProductionDate         *int64 `json:"start_production_date,omitempty"`
	CapacityPerDay              *int64 `json:"capacity_per_day"`

	// Expand response
	UnseenCommentCount *int64 `gorm:"-" json:"unseen_comment_count,omitempty"`

	RejectReason  string       `gorm:"size:2000" json:"reject_reason,omitempty"`           // Seller rejected reasons
	ExpectedPrice *price.Price `gorm:"type:decimal(20,4)" json:"expected_price,omitempty"` // Seller expected
	Note          string       `gorm:"size:2000" json:"note,omitempty"`

	AdminRejectReason string `gorm:"size:2000" json:"admin_reject_reason,omitempty"`

	QuotedPrice *price.Price `json:"quoted_price,omitempty"`
}

type BulkPurchaseOrderSellerQuotations []*BulkPurchaseOrderSellerQuotation
