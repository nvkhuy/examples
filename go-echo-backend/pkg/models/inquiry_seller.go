package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
)

type InquirySeller struct {
	ID        string `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	InquiryID string   `gorm:"uniqueIndex:idx_inquiry_seller" json:"inquiry_id"`
	Inquiry   *Inquiry `gorm:"-" json:"inquiry"`

	PurchaseOrderID string         `gorm:"uniqueIndex:idx_inquiry_seller" json:"purchase_order_id,omitempty"`
	PurchaseOrder   *PurchaseOrder `gorm:"-" json:"purchase_order"`

	UserID string `gorm:"uniqueIndex:idx_inquiry_seller" json:"user_id"`
	User   *User  `gorm:"-" json:"user,omitempty"`

	DueDay *int64 `json:"due_day,omitempty"`

	Status enums.InquirySellerStatus `gorm:"default:'new'" json:"status"`

	InquirySellerSkus []*InquirySellerSku `gorm:"-" json:"inquiry_seller_skus,omitempty"`

	DeliveryDate *int64 `json:"delivery_date,omitempty"`

	Currency enums.Currency `json:"currency,omitempty"`

	// Price from Admin
	OrderType   enums.OBOrderType `json:"order_type,omitempty"`
	OfferPrice  *price.Price      `json:"offer_price,omitempty"`
	OfferRemark string            `json:"offer_remark,omitempty"`

	VarianceAmount     *price.Price `json:"variance_amount,omitempty"`
	VariancePercentage *float64     `json:"variance_percentage,omitempty"`

	// Price from seller
	FabricCost      *price.Price `json:"fabric_cost,omitempty"`
	DecorationCost  *price.Price `json:"decoration_cost,omitempty"`
	MakingCost      *price.Price `json:"making_cost,omitempty"` // sewing, cut, making, finishing
	OtherCost       *price.Price `json:"other_cost,omitempty"`
	SampleUnitPrice *price.Price `json:"sample_unit_price,omitempty"`
	SampleLeadTime  *int64       `json:"sample_lead_time,omitempty"` // unit: day
	SellerRemark    string       `json:"seller_remark,omitempty"`

	AdminSentAt    *int64                   `json:"admin_sent_at,omitempty"`
	QuotationAt    *int64                   `json:"quotation_at,omitempty"`
	BulkQuotations *SellerBulkQuotationMOQs `json:"bulk_quotations,omitempty"`

	ExpectedStartProductionDate *int64 `json:"expected_start_production_date,omitempty"`
	StartProductionDate         *int64 `json:"start_production_date,omitempty"`
	CapacityPerDay              *int64 `json:"capacity_per_day"`

	// Expand response
	UnseenCommentCount *int64 `gorm:"-" json:"unseen_comment_count,omitempty"`

	RejectReason  string       `gorm:"size:2000" json:"reject_reason,omitempty"`           // Seller rejected reasons
	ExpectedPrice *price.Price `gorm:"type:decimal(20,4)" json:"expected_price,omitempty"` // Seller expected
	Note          string       `gorm:"size:2000" json:"note,omitempty"`

	AdminRejectReason string `gorm:"size:2000" json:"admin_reject_reason,omitempty"`
}

type InquirySellers []*InquirySeller

type InquirySellerCreateQuatationParams struct {
	JwtClaimsInfo
	InquirySellerID string `json:"inquiry_seller_id" param:"inquiry_seller_id" validate:"required"`

	FabricCost      *price.Price `json:"fabric_cost,omitempty"`
	DecorationCost  *price.Price `json:"decoration_cost,omitempty"`
	MakingCost      *price.Price `json:"making_cost,omitempty"` // sewing, cut, making, finishing
	OtherCost       *price.Price `json:"other_cost,omitempty"`
	SampleUnitPrice *price.Price `json:"sample_unit_price,omitempty"`
	SampleLeadTime  *int64       `json:"sample_lead_time,omitempty"` // unit: day
	SellerRemark    string       `json:"seller_remark,omitempty"`

	BulkQuotations *SellerBulkQuotationMOQs `json:"bulk_quotations,omitempty"`

	StartProductionDate *int64 `json:"start_production_date"`
	CapacityPerDay      *int64 `json:"capacity_per_day"`
}
type SellerQuatationPerSkuForm struct {
	InquirySkuID string         `json:"inquiry_sku_id"`
	Price        price.Price    `gorm:"type:decimal(20,4);default:0.0"  json:"price"`
	PriceType    string         `json:"price_type"`
	Currency     enums.Currency `gorm:"default:'USD'" json:"currency"`
}

type GetInquirySellerStatusCountResponse struct {
	TotalNewInquiries                int `json:"total_new_inquiries"`
	TotalSentInquiries               int `json:"total_sent_inquiries"`
	TotalWaitingForApprovedInquiries int `json:"total_waiting_for_approved_inquiries"`
}

type SubmitMultipleInquirySellerQuotationRequest struct {
	JwtClaimsInfo
	Quotations []*InquirySellerCreateQuatationParams `json:"quotations" validate:"required"`
}
