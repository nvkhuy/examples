package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/lib/pq"
)

type Inquiries []*Inquiry

// Faq Faq's model
type Inquiry struct {
	Model

	ReferenceID  string              `gorm:"unique" json:"reference_id"`
	Title        string              `gorm:"size:1000" json:"title,omitempty" validate:"required"`
	Requirement  string              `gorm:"size:1000" json:"requirement,omitempty"`
	SkuNote      string              `gorm:"size:1000" json:"sku_note,omitempty"`
	UserID       string              `gorm:"size:100" json:"user_id,omitempty" validate:"required"`
	ExpiredDate  *int64              `json:"expired_date,omitempty"`
	DeliveryDate *int64              `json:"delivery_date,omitempty"` // date when admin send quotation
	Status       enums.InquiryStatus `gorm:"default:'new'" json:"status,omitempty"`

	Attachments         *Attachments `json:"attachments,omitempty"`
	Document            *Attachments `json:"document,omitempty"`
	Design              *Attachments `json:"design,omitempty"`
	FabricAttachments   *Attachments `json:"fabric_attachments,omitempty"`
	TechpackAttachments *Attachments `json:"techpack_attachments,omitempty"`

	Currency   enums.Currency `gorm:"default:'USD'" json:"currency,omitempty"`
	CategoryID string         `gorm:"size:100" json:"category_id,omitempty"`

	Category *Category `json:"category,omitempty" gorm:"-"`
	User     *User     `gorm:"-" json:"user,omitempty"`

	BuyerQuotationStatus enums.InquirySkuStatus `gorm:"default:'new;size:100'" json:"buyer_quotation_status,omitempty"`

	NewSeenAt    *int64 `json:"new_seen_at,omitempty"`
	UpdateSeenAt *int64 `json:"update_seen_at,omitempty"`

	Quantity         *int64                 `json:"quantity,omitempty"`
	PriceType        enums.InquiryPriceType `json:"price_type,omitempty"`
	ExpectedPrice    *price.Price           `gorm:"type:decimal(20,4);default:0.0" json:"expected_price,omitempty"`
	PriceDestination string                 `gorm:"size:1000" json:"price_desination,omitempty"`

	// Quotation
	InternalQuotationApprovedAt *int64 `json:"internal_quotation_approved_at,omitempty"`
	InternalQuotationApprovedBy string `json:"internal_quotation_approved_by,omitempty"` // Admin/CEO id
	InternalQuotationCreatedBy  string `json:"internal_quotation_created_by,omitempty"`  // Admin id
	InternalQuotationCreatedAt  *int64 `json:"internal_quotation_created_at,omitempty"`

	QuotationApprovedAt    *int64                        `json:"quotation_approved_at,omitempty"`
	QuotationAt            *int64                        `json:"quotation_at,omitempty"`
	AdminQuotations        InquiryQuotationItems         `json:"admin_quotations,omitempty"`
	ApproveRejectMeta      *InquiryApproveRejectMeta     `json:"approve_reject_meta,omitempty"`
	ApproveRejectMetaItems InquiryApproveRejectMetaItems `json:"approve_reject_meta_items,omitempty"`
	CollectionID           string                        `json:"collection_id,omitempty"`
	Collection             *InquiryCollection            `gorm:"-" json:"collection,omitempty"`

	ShippingAddressID string       `gorm:"size:100" json:"shipping_address_id,omitempty"`
	ShippingAddress   *Address     `gorm:"-" json:"shipping_address,omitempty"`
	ShippingFee       *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"shipping_fee,omitempty"`
	TaxPercentage     *float64     `gorm:"type:decimal(20,4);default:0.0" json:"tax_percentage,omitempty"`

	PurchaseOrder *PurchaseOrder     `gorm:"-" json:"purchase_order,omitempty"`
	CartItems     []*InquiryCartItem `gorm:"-" json:"cart_items,omitempty"` //legacy

	SizeList         string   `gorm:"size:1000" json:"size_list,omitempty"`
	SizeChart        string   `gorm:"size:1000" json:"size_chart,omitempty"`
	Composition      string   `gorm:"size:1000" json:"composition,omitempty"`
	StyleNo          string   `gorm:"size:1000" json:"style_no,omitempty"`
	FabricName       string   `gorm:"size:1000" json:"fabric_name,omitempty"`
	FabricWeight     *float64 `json:"fabric_weight,omitempty"`
	FabricWeightUnit string   `gorm:"default:'gsm';size:1000" json:"fabric_weight_unit,omitempty"`
	ColorList        string   `gorm:"size:1000" json:"color_list,omitempty"`
	ProductWeight    *float64 `json:"product_weight,omitempty"`

	QuotedPrice *price.Price `gorm:"-" json:"quoted_price,omitempty"`

	AuditLogs []*InquiryAudit `gorm:"-" json:"audit_logs,omitempty"`

	AssigneeIDs pq.StringArray `gorm:"type:varchar(200)[]" json:"assignee_ids,omitempty"`
	Assignees   Users          `gorm:"-" json:"assignees,omitempty"`
	CreatedBy   string         `gorm:"size:100" json:"created_by,omitempty"`

	SamplePaymentLinkID string `gorm:"size:100" json:"sample_payment_link_id,omitempty"`
	SamplePaymentLink   string `gorm:"size:1000" json:"sample_payment_link,omitempty"`

	CloseReason  *InquiryApproveRejectMeta `json:"close_reason,omitempty"`
	CancelReason *InquiryApproveRejectMeta `json:"cancel_reason,omitempty"`

	// est costing
	CMCost     *price.Price `json:"cm_cost,omitempty"`
	FabricCost *price.Price `json:"fabric_cost,omitempty"`
	BasicCost  *price.Price `json:"basic_cost,omitempty"`

	EstTACost    *price.Price `json:"est_ta_cost,omitempty"`
	EstOtherCost *price.Price `json:"est_other_cost,omitempty"`
	EstTotalCost *price.Price `json:"est_total_cost,omitempty"`

	ProductID    string               `gorm:"size:100" json:"product_id,omitempty"`
	Requirements *InquiryRequirements `json:"requirements,omitempty"`

	HubspotDealID         string `gorm:"size:100"  json:"hubspot_deal_id,omitempty"`
	EditTimeout           *int64 `json:"edit_timeout,omitempty"`
	IsEditTimeoutExtended *bool  `gorm:"type:bool;default:false" json:"is_edit_timeout_extended,omitempty"`

	OrderGroupID string      `json:"order_group_id,omitempty"`
	OrderGroup   *OrderGroup `gorm:"-" json:"order_group,omitempty"`
}

type InquiryAdminCreateForm struct {
	JwtClaimsInfo

	AssigneeId string `json:"assignee_id"`
	BuyerId    string `json:"buyer_id" validate:"required"`
	CreatedBy  string `json:"created_by"`

	Title       string `json:"title,omitempty" validate:"required"`
	Requirement string `json:"requirement,omitempty"`
	SkuNote     string `json:"sku_note,omitempty"`

	Attachments         *Attachments `json:"attachments,omitempty"`
	Design              *Attachments `json:"design,omitempty"`
	FabricAttachments   *Attachments `json:"fabric_attachments,omitempty"`
	TechpackAttachments *Attachments `json:"techpack_attachments,omitempty"`

	Status           enums.InquiryStatus    `gorm:"default:'new'" json:"status"`
	PriceType        enums.InquiryPriceType `json:"price_type,omitempty"`
	Quantity         *int64                 `json:"quantity,omitempty"`
	ExpectedPrice    price.Price            `json:"expected_price,omitempty"`
	PriceDestination string                 `json:"price_desination,omitempty"`
	CollectionID     string                 `json:"collection_id,omitempty"`

	SizeList         string   `json:"size_list,omitempty"`
	SizeChart        string   `json:"size_chart,omitempty"`
	Composition      string   `json:"composition,omitempty"`
	StyleNo          string   `json:"style_no,omitempty"`
	FabricName       string   `json:"fabric_name,omitempty"`
	FabricWeight     *float64 `json:"fabric_weight,omitempty"`
	FabricWeightUnit string   `gorm:"default:'gsm'" json:"fabric_weight_unit,omitempty"`
	ColorList        string   `json:"color_list,omitempty"`
	ProductWeight    *float64 `json:"product_weight,omitempty"`

	Currency        enums.Currency `gorm:"default:'USD'" json:"currency" validate:"required"`
	ShippingAddress *Address       `json:"shipping_address"`
	CountryCode     string         `json:"country_code"`
}

type InquiryCreateForm struct {
	JwtClaimsInfo

	Title   string `json:"title,omitempty" validate:"required"`
	SkuNote string `json:"sku_note,omitempty"`

	Attachments         *Attachments `json:"attachments,omitempty" validate:"required"`
	FabricAttachments   *Attachments `json:"fabric_attachments,omitempty"`
	TechpackAttachments *Attachments `json:"techpack_attachments,omitempty"`

	Quantity      int64       `json:"quantity,omitempty" validate:"required"`
	ExpectedPrice price.Price `json:"expected_price,omitempty"`

	SizeList         string   `json:"size_list,omitempty" validate:"required"`
	SizeChart        string   `json:"size_chart,omitempty"`
	Composition      string   `json:"composition,omitempty"`
	StyleNo          string   `json:"style_no,omitempty"`
	FabricName       string   `json:"fabric_name,omitempty"`
	FabricWeight     *float64 `json:"fabric_weight,omitempty"`
	FabricWeightUnit string   `gorm:"default:'gsm'" json:"fabric_weight_unit,omitempty"`
	ColorList        string   `json:"color_list,omitempty"`
	ProductWeight    *float64 `json:"product_weight,omitempty"`

	Currency        enums.Currency `gorm:"default:'USD'" json:"currency" validate:"required"`
	ShippingAddress *Address       `json:"shipping_address"`

	// Requirement string `json:"requirement,omitempty"`
	// Design              *Attachments `json:"design,omitempty"`
	// Status           enums.InquiryStatus    `gorm:"default:'new'" json:"status"`
	// PriceType        enums.InquiryPriceType `json:"price_type,omitempty"`
	// ShippingFee price.Price `json:"shipping_fee"`
	// CountryCode string      `json:"country_code"`
	// PriceDestination string                 `json:"price_desination,omitempty"`

	// seller est costing
	// CMCost     *price.Price `json:"cm_cost"`
	// FabricCost *price.Price `json:"fabric_cost"`
	// BasicCost  *price.Price `json:"basic_cost"`
	// EstTACost    *price.Price `json:"est_ta_cost"` // Trims && Accesories
	// EstOtherCost *price.Price `json:"est_other_cost"`
	// EstTotalCost *price.Price `json:"est_total_cost"`
	// ProductID    string               `json:"product_id"`
	// Requirements *InquiryRequirements `json:"requirements"`
	// EditTimeout  int64                `json:"edit_timeout"`
	OrderGroupID string `json:"order_group_id,omitempty"`
}
type CreateMultipleInquiriesRequest struct {
	JwtClaimsInfo
	Inquiries []*InquiryCreateForm `json:"inquiries" validate:"required"`
}

type InquiryUserCreateInfo struct {
	Name        string `json:"name,omitempty" validate:"required"`
	BrandName   string `json:"brand_name,omitempty" validate:"required"`
	Email       string `gorm:"unique;type:citext;default:null" json:"email,omitempty" validate:"required"`
	PhoneNumber string `json:"phone_number,omitempty" param:"phone_number" query:"phone_number" form:"phone_number" validate:"isPhone"`
}

type InquiryUpdateForm struct {
	JwtClaimsInfo

	InquiryID string `param:"inquiry_id" validate:"required"`

	Title               string              `json:"title,omitempty" validate:"required"`
	Requirement         string              `json:"requirement,omitempty"`
	SkuNote             string              `json:"sku_note,omitempty"`
	ExpiredDate         *int64              `json:"expired_date,omitempty"`
	Attachments         *Attachments        `json:"attachments,omitempty"`
	Document            *Attachments        `json:"document,omitempty"`
	FabricAttachments   *Attachments        `json:"fabric_attachments,omitempty"`
	TechpackAttachments *Attachments        `json:"techpack_attachments,omitempty"`
	Design              *Attachments        `json:"design,omitempty"`
	CategoryID          string              `json:"category_id,omitempty"`
	CollectionID        string              `json:"collection_id,omitempty"`
	Status              enums.InquiryStatus `gorm:"default:'new'" json:"status"`
	ExpectedPrice       price.Price         `json:"expected_price,omitempty"`

	SizeList         string   `json:"size_list,omitempty"`
	SizeChart        string   `json:"size_chart,omitempty"`
	Composition      string   `json:"composition,omitempty"`
	StyleNo          string   `json:"style_no,omitempty"`
	FabricName       string   `json:"fabric_name,omitempty"`
	FabricWeight     *float64 `json:"fabric_weight,omitempty"`
	FabricWeightUnit string   `gorm:"default:'gsm'" json:"fabric_weight_unit,omitempty"`

	ColorList       string      `json:"color_list,omitempty"`
	ProductWeight   *float64    `json:"product_weight,omitempty"`
	ShippingAddress *Address    `json:"shipping_address"`
	ShippingFee     price.Price `json:"shipping_fee"`

	Currency     enums.Currency `json:"currency"`
	Quantity     *int64         `json:"quantity,omitempty"`
	OrderGroupID string         `json:"order_group_id,omitempty"`
}

type InquiryEditTimeoutUpdateForm struct {
	JwtClaimsInfo

	InquiryID   string `param:"inquiry_id" validate:"required"`
	EditTimeout int64  `json:"edit_timeout" validate:"required"`
}

type SellerRequestQuotationInfo struct {
	OfferPrice                  *price.Price      `json:"offer_price" param:"offer_price" validate:"required"`
	OfferRemark                 string            `json:"offer_remark" param:"offer_remark"`
	SellerID                    string            `json:"seller_id" param:"seller_id"`
	VarianceAmount              *price.Price      `json:"variance_amount,omitempty"`
	VariancePercentage          *float64          `json:"variance_percentage,omitempty"`
	OrderType                   enums.OBOrderType `json:"order_type,omitempty"`
	ExpectedStartProductionDate *int64            `json:"expected_start_production_date"`
}

type InquiryMarkSeenForm struct {
	JwtClaimsInfo

	InquiryID  string `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`
	ActionType string `json:"action_type" param:"action_type" query:"action_type" validate:"required"`
}

type InquiryCart struct {
	Items []*Inquiry `json:"items" query:"items" param:"items" validate:"required"`
}

type InquiryIDParam struct {
	JwtClaimsInfo

	InquiryID string `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`
	Note      string `json:"note"`

	PurchaseOrderID string         `json:"purchase_order_id"`
	PurchaseOrder   *PurchaseOrder `json:"-"`
}

type InquiryRemoveItemsForm struct {
	JwtClaimsInfo

	InquiryID string `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`

	ItemIDs []string `json:"item_ids" validate:"required"`
}

type InquiryAssignPICParam struct {
	JwtClaimsInfo
	InquiryID   string   `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`
	AssigneeIDs []string `json:"assignee_ids" query:"assignee_ids" param:"assignee_ids" validate:"required"`
}

type InquiryCloseForm struct {
	JwtClaimsInfo

	InquiryID   string                    `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`
	CloseReason *InquiryApproveRejectMeta `json:"close_reason,omitempty"`
}

type SubmitMultipleInquiryQuotationRequest struct {
	JwtClaimsInfo
	Quotations []*SendInquiryToBuyerForm `json:"quotations" validate:"required"`
}

type InquiryMultiplePreviewCheckoutParams struct {
	JwtClaimsInfo
	InquiryIDs []string `param:"inquiry_ids" validate:"required"`
}

type ApproveMultipleInquiryQuotationsRequest struct {
	JwtClaimsInfo
	InquiryIDs []string `json:"inquiry_ids" validate:"required"`
}
