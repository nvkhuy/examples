package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
)

type SendInquiryToBuyerForm struct {
	JwtClaimsInfo

	InquiryID     string                `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`
	Quotations    InquiryQuotationItems `json:"quotations" param:"quotations" query:"quotations"`
	ProductWeight *float64              `json:"product_weight" param:"product_weight" query:"product_weight"`
	ShippingFee   *price.Price          `json:"shipping_fee"`
	TaxPercentage *float64              `json:"tax_percentage" validate:"min=0,max=100"`
}

type ApproveInquiryQuotationItem struct {
	Quantity *int64 `json:"quantity"`
}

type AdminInternalApproveQuotationForm struct {
	JwtClaimsInfo

	InquiryID string `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`
}

type RejectMultipleInquiryQuotationsRequest struct {
	JwtClaimsInfo
	InquiryIDs []string `json:"inquiry_ids" validate:"required"`
	Reason     string   `json:"reason" validate:"required"`
}
