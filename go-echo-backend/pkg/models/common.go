package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/golang-jwt/jwt"
)

type M map[string]interface{}

type Records struct {
	Records interface{} `json:"records"`
}

type Model struct {
	ID        string     `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`
}

type JwtClaimsInfo struct {
	role     enums.Role     `gorm:"-" json:"-" swaggerignore:"true"`
	userID   string         `gorm:"-" json:"-" swaggerignore:"true"`
	isGhost  bool           `gorm:"-" json:"-" swaggerignore:"true"`
	ghostID  string         `gorm:"-" json:"-" swaggerignore:"true"`
	timezone enums.Timezone `gorm:"-" json:"-" swaggerignore:"true"`
}
type PaginationParams struct {
	Page         int    `json:"page" query:"page" form:"page" validate:"min=0"`
	Limit        int    `json:"limit" query:"limit" form:"limit" validate:"max=100"`
	Keyword      string `json:"keyword" query:"keyword" form:"keyword" validate:"max=200"`
	WithoutCount bool   `json:"-"`
}

type AssetCustomClaims struct {
	FileName string `json:"file_name"`
	jwt.StandardClaims
}

type CheckoutLinkCustomClaims struct {
	LinkID string `json:"link_id"`
	jwt.StandardClaims
}

type RoleConstant struct {
	Value enums.Role `json:"value"`
	Name  string     `json:"name"`
}

type BrandTeamRoleConstant struct {
	Value enums.BrandTeamRole `json:"value"`
	Name  string              `json:"name"`
}
type AccountStatusConstant struct {
	Value enums.AccountStatus `json:"value"`
	Name  string              `json:"name"`
}

type AddressTypeConstant struct {
	Value enums.AddressType `json:"value"`
	Name  string            `json:"name"`
}

type ChatMessageTypeConstant struct {
	Value enums.ChatMessageType `json:"value"`
	Name  string                `json:"name"`
}

type PostStatusConstant struct {
	Value enums.PostStatus `json:"value"`
	Name  string           `json:"name"`
}

type DocumentStatusConstant struct {
	Value enums.DocumentStatus `json:"value"`
	Name  string               `json:"name"`
}

type PaymentStatusConstant struct {
	Value enums.PaymentStatus `json:"value"`
	Name  string              `json:"name"`
}
type RegisterBusinessConstant struct {
	Value   enums.RegisterBusiness `json:"value"`
	Name    string                 `json:"name"`
	IconUrl string                 `json:"icon_url"`
}

type RegisterQuantityConstant struct {
	Value enums.RegisterQuantity `json:"value"`
	Name  string                 `json:"name"`
}

type RegisterAreaConstant struct {
	Value   enums.RegisterArea `json:"value"`
	Name    string             `json:"name"`
	IconUrl string             `json:"icon_url"`
}

type RegisterProductCategory struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type ProductUnitConstant struct {
	Value enums.ProductUnit `json:"value"`
	Name  string            `json:"name"`
}

type ProductTypeConstant struct {
	Value enums.ProductType `json:"value"`
	Name  string            `json:"name"`
}

type InquiryStatusConstant struct {
	Value enums.InquiryStatus `json:"value"`
	Name  string              `json:"name"`
}

type InquiryPriceTypeConstant struct {
	Value enums.InquiryPriceType `json:"value"`
	Name  string                 `json:"name"`
}

type InquiryMOQConstant struct {
	Value enums.InquiryMOQType `json:"value"`
	Name  string               `json:"name"`
}

type CertificationConstant struct {
	Value enums.CertificationType `json:"value"`
	Name  string                  `json:"name"`
}
type SellerQuotationFilterConstant struct {
	Value enums.SellerQuotationFilter `json:"value"`
	Name  string                      `json:"name"`
}
type InquirySkuStatusConstant struct {
	Value enums.InquirySkuStatus `json:"value"`
	Name  string                 `json:"name"`
}

type LabelDimensionConstant struct {
	Value enums.LabelDimension `json:"value"`
	Name  string               `json:"name"`
}
type LabelMaterialConstant struct {
	Value enums.LabelMaterial `json:"value"`
	Name  string              `json:"name"`
}
type LabelSizeConstant struct {
	Value enums.LabelSize `json:"value"`
	Name  string          `json:"name"`
}
type LabelKindConstant struct {
	Value enums.LabelKind `json:"value"`
	Name  string          `json:"name"`
}
type LabelStatusConstant struct {
	Value enums.LabelStatus `json:"value"`
	Name  string            `json:"name"`
}
type LabelAccessoryConstant struct {
	Value enums.LabelAccessory `json:"value"`
	Name  string               `json:"name"`
}

type InquirySkuRejectReasonConstant struct {
	Value enums.InquirySkuRejectReason `json:"value"`
	Name  string                       `json:"name"`
}

type CardTypeConstant struct {
	Value enums.CardType `json:"value"`
	Name  string         `json:"name"`
}

type CardAttributeConstant struct {
	Value enums.CardAttribute `json:"value"`
}

type FabricRawStatusConstant struct {
	Name  string                `json:"name"`
	Value enums.FabricRawStatus `json:"value"`
}

type FabricBulkProductionConstant struct {
	Name  string                           `json:"name"`
	Value enums.FabricBulkProductionStatus `json:"value"`
}

type QcReportTypeConstant struct {
	Name  string             `json:"name"`
	Value enums.QcReportType `json:"value"`
}

type QcReportResultConstant struct {
	Name  string               `json:"name"`
	Value enums.QcReportResult `json:"value"`
}

type DeliveryTypeConstant struct {
	Name  string             `json:"name"`
	Value enums.DeliveryType `json:"value"`
}

type DeliveryAttributeNameConstant struct {
	Name  string                      `json:"name"`
	Value enums.DeliveryAttributeName `json:"value"`
}

type DeliveryStatusConstant struct {
	Name  string               `json:"name"`
	Value enums.DeliveryStatus `json:"value"`
}

type InflowPaymentMethodConstant struct {
	Name  string                    `json:"name"`
	Value enums.InflowPaymentMethod `json:"value"`
}

type CurrencyConstant struct {
	Name  string         `json:"name"`
	Value enums.Currency `json:"value"`
}

type InquiryTypeConstant struct {
	Name  string            `json:"name"`
	Value enums.InquiryType `json:"value"`
}

type InquiryAttributeConstant struct {
	Value enums.InquiryAttribute `json:"value"`
}

type InquirySizeChartConstant struct {
	Value enums.InquirySizeChart `json:"value"`
}

type ProductAttributeMetaConstant struct {
	Value enums.ProductAttribute `json:"value"`
}

type PoTrackingStatusConstant struct {
	Value enums.PoTrackingStatus `json:"value"`
	Name  string                 `json:"name"`
}

type PoCatalogTrackingStatusConstant struct {
	Value enums.PoCatalogTrackingStatus `json:"value"`
	Name  string                        `json:"name"`
}

type BulkPoTrackingStatusConstant struct {
	Value enums.BulkPoTrackingStatus `json:"value"`
	Name  string                     `json:"name"`
}

type MaterialTypesConstant struct {
	Value enums.BulkPoTrackingStatus `json:"value"`
	Name  string                     `json:"name"`
}

type BulkQCReportTypesConstant struct {
	Value enums.BulkPoTrackingStatus `json:"value"`
	Name  string                     `json:"name"`
}

type PoRawMaterialStatusConstant struct {
	Value enums.PoRawMaterialStatus `json:"value"`
	Name  string                    `json:"name"`
}

type InquiryBuyerStatusConstant struct {
	Name  string                   `json:"name"`
	Value enums.InquiryBuyerStatus `json:"value"`
}

type TeamConstant struct {
	Name  string     `json:"name"`
	Value enums.Team `json:"value"`
}

type FeatureConstant struct {
	Name  string            `json:"name"`
	Value enums.FeatureType `json:"value"`
}
type InvoiceTypeConstant struct {
	Name  string            `json:"name"`
	Value enums.InvoiceType `json:"value"`
}

type BrandTypeConstant struct {
	Name  string          `json:"name"`
	Value enums.BrandType `json:"value"`
}

type ShippingMethodConstant struct {
	Value       enums.ShippingMethod `json:"value"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
}

type SupplierTypeConstant struct {
	Value enums.SupplierType `json:"value"`
	Name  string             `json:"name"`
}

type OnboardingConstant struct {
	Value       string `json:"value"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ImageUrl    string `json:"image_url,omitempty"`
}

type InquirySellerStatusConstant struct {
	Value enums.InquirySellerStatus `json:"value"`
	Name  string                    `json:"name"`
}

type CommentTargetTypeConstant struct {
	Value enums.CommentTargetType `json:"value"`
	Name  string                  `json:"name"`
}

type InquirySellerQuotationTypeConstant struct {
	Value enums.InquirySellerQuotationType `json:"value"`
	Name  string                           `json:"name"`
}

type AdsVideoSectionConstant struct {
	Value enums.AdsVideoSection `json:"value"`
	Name  string                `json:"name"`
}

type FabricWeightUnitConstant struct {
	Value enums.FabricWeightUnit `json:"value"`
	Name  string                 `json:"name"`
}
type BulkPOSellerStatusConstant struct {
	Value enums.BulkPurchaseOrderSellerStatus `json:"value"`
	Name  string                              `json:"name"`
}

type BulkPoSellerTrackingStatusConstant struct {
	Value enums.SellerBulkPoTrackingStatus `json:"value"`
	Name  string                           `json:"name"`
}

type Constants struct {
	Roles             []*RoleConstant             `json:"roles"`
	BrandTeamRoles    []*BrandTeamRoleConstant    `json:"brand_team_roles"`
	AccountStatuses   []*AccountStatusConstant    `json:"account_statuses"`
	BrandTypes        []*BrandTypeConstant        `json:"brand_types"`
	AddressTypes      []*AddressTypeConstant      `json:"address_types"`
	ChatMessageTypes  []*ChatMessageTypeConstant  `json:"chat_message_types"`
	PostStatues       []*PostStatusConstant       `json:"post_statuses"`
	DocumentStatuses  []*DocumentStatusConstant   `json:"document_statuses"`
	PaymentStatuses   []*PaymentStatusConstant    `json:"payment_statuses"`
	DefaultImageTypes []*DefaultImageTypeConstant `json:"default_image_types"`
	FilterRatings     []*FilterRatingConstant     `json:"filter_ratings"`
	FilterMinOrders   []*FilterMinOrderConstant   `json:"filter_min_orders"`

	ProductUnits                 []*ProductUnitConstant                `json:"product_units"`
	ProductTypes                 []*ProductTypeConstant                `json:"product_types"`
	InquiryStatuses              []*InquiryStatusConstant              `json:"inquiry_statuses"`
	InquiryPriceTypes            []*InquiryPriceTypeConstant           `json:"inquiry_price_types"`
	InquiryMOQs                  []*InquiryMOQConstant                 `json:"inquiry_moqs"`
	Certifications               []*CertificationConstant              `json:"certifications"`
	SellerQuotationFilters       []*SellerQuotationFilterConstant      `json:"seller_quotation_filter"`
	InquirySkuStatuses           []*InquirySkuStatusConstant           `json:"inquiry_sku_statuses"`
	LabelDimensions              []*LabelDimensionConstant             `json:"label_dimensions"`
	BarcodeDimensions            []*LabelDimensionConstant             `json:"barcode_dimensions"`
	LabelMaterials               []*LabelMaterialConstant              `json:"label_materials"`
	LabelSizes                   []*LabelSizeConstant                  `json:"label_sizes"`
	LabelKinds                   []*LabelKindConstant                  `json:"label_kinds"`
	LabelStatues                 []*LabelStatusConstant                `json:"label_statuses"`
	LabelAccessories             []*LabelAccessoryConstant             `json:"label_accessories"`
	InquirySkuRejectReasons      []*InquirySkuRejectReasonConstant     `json:"inquiry_sku_reject_reasons"`
	CardTypes                    []*CardTypeConstant                   `json:"card_types"`
	CardAttributes               []*CardAttributeConstant              `json:"card_attributes"`
	FabricRawStatuses            []*FabricRawStatusConstant            `json:"fabric_raw_statuses"`
	FabricBulkProductionStatuses []*FabricBulkProductionConstant       `json:"fabric_bulk_production_statuses"`
	QcReportTypes                []*QcReportTypeConstant               `json:"qc_report_types"`
	QcReportResults              []*QcReportResultConstant             `json:"qc_report_results"`
	DeliveryTypes                []*DeliveryTypeConstant               `json:"delivery_types"`
	DeliveryAttributeNames       []*DeliveryAttributeNameConstant      `json:"delivery_attribute_names"`
	DeliveryStatuses             []*DeliveryStatusConstant             `json:"delivery_statuses"`
	InflowPaymentMethods         []*InflowPaymentMethodConstant        `json:"inflow_payment_methods"`
	Currencies                   []*CurrencyConstant                   `json:"currencies"`
	InquiryTypes                 []*InquiryTypeConstant                `json:"inquiry_types"`
	InquiryAttributes            []*InquiryAttributeConstant           `json:"inquiry_attributes"`
	InquirySizeCharts            []*InquirySizeChartConstant           `json:"inquiry_size_charts"`
	ProductAttributeMetas        []*ProductAttributeMetaConstant       `json:"product_attribute_metas"`
	PoCatalogTrackingStatuses    []*PoCatalogTrackingStatusConstant    `json:"po_catalog_tracking_statusess"`
	PoTrackingStatuses           []*PoTrackingStatusConstant           `json:"po_tracking_statuses"`
	BulkPoTrackingStatuses       []*BulkPoTrackingStatusConstant       `json:"bulk_po_tracking_statuses"`
	PoRawMaterialStatuses        []*PoRawMaterialStatusConstant        `json:"po_raw_material_statuses"`
	InquiryBuyerStatuses         []*InquiryBuyerStatusConstant         `json:"inquiry_buyer_statuses"`
	ShippingMethods              []*ShippingMethodConstant             `json:"shipping_methods"`
	SupplierTypes                []*SupplierTypeConstant               `json:"supplier_types"`
	OBOrderTypes                 []*OnboardingConstant                 `json:"ob_order_types"`
	OBShippingTerms              []*OnboardingConstant                 `json:"ob_shipping_terms"`
	OBProductGroups              []*OnboardingConstant                 `json:"ob_product_groups"`
	OBMOQTypes                   []*OnboardingConstant                 `json:"ob_moq_types"`
	OBLeadTimes                  []*OnboardingConstant                 `json:"ob_lead_times"`
	OBFabricTypes                []*OnboardingConstant                 `json:"ob_fabric_types"`
	OBMillFabricTypes            []*OnboardingConstant                 `json:"ob_mill_fabric_types"`
	OBFactoryProductTypes        []*OnboardingConstant                 `json:"ob_factory_product_types"`
	OBDevelopmentServices        []*OnboardingConstant                 `json:"ob_development_services"`
	OBOutputUnits                []*OnboardingConstant                 `json:"ob_output_units"`
	OBSewingAccessoryTypes       []*OnboardingConstant                 `json:"ob_sewing_accessory_types"`
	OBPackingAccessoryTypes      []*OnboardingConstant                 `json:"ob_packing_accessory_types"`
	OBServiceTypes               []*OnboardingConstant                 `json:"ob_service_types"`
	OBDecorationServices         []*OnboardingConstant                 `json:"ob_decoration_services"`
	OBPaymentTerms               []*OnboardingConstant                 `json:"ob_payment_terms"`
	InquirySellerStatuses        []*InquirySellerStatusConstant        `json:"inquiry_seller_statuses"`
	CommentTargetTypes           []*CommentTargetTypeConstant          `json:"comment_target_types"`
	InquirySellerQuotationTypes  []*InquirySellerQuotationTypeConstant `json:"inquiry_seller_quotation_types"`
	AdsVideoSections             []*AdsVideoSectionConstant            `json:"ads_video_sections"`
	FabricWeightUnits            []*FabricWeightUnitConstant           `json:"fabric_weight_units"`
	MaterialTypes                []*MaterialTypesConstant              `json:"material_types"`
	BulkQCReportTypes            []*BulkQCReportTypesConstant          `json:"bulk_qc_report_types"`
	BulkPOSellerStatuses         []*BulkPOSellerStatusConstant         `json:"bulk_po_seller_statuses"`
	BulkPOSellerTrackingStatuses []*BulkPoSellerTrackingStatusConstant `json:"bulk_po_seller_tracking_statuses"`

	InflowBillingAddresses []*InvoiceParty        `json:"inflow_billing_addresses"`
	InvoiceTypes           []*InvoiceTypeConstant `json:"invoice_types"`
	Teams                  []*TeamConstant        `json:"teams"`
	Features               []*FeatureConstant     `json:"features"`
}

type RegisterConstants struct {
	RegisterBusinesses        []*RegisterBusinessConstant `json:"register_businesses"`
	RegisterQuantities        []*RegisterQuantityConstant `json:"register_quantities"`
	RegisterAreas             []*RegisterAreaConstant     `json:"register_areas"`
	RegisterProductCategories []*RegisterProductCategory  `json:"register_product_categories"`
}

type CheckExistsForm struct {
	Email       string `json:"email" validate:"omitempty,required_without=PhoneNumber UserName,email"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,required_without=Email UserName,isPhone"`
	UserName    string `json:"user_name" validate:"omitempty,required_without=Email PhoneNumber"`
}

type CheckExistsResponse struct {
	IsExists bool `json:"is_exists"`
}

type DispatchTaskForm struct {
	Name           string                 `json:"name" validate:"required"`
	Data           map[string]interface{} `json:"data"`
	DelayInSeconds int                    `json:"delay_in_seconds"`
}

type DefaultImageTypeConstant struct {
	Value enums.DefaultImageType `json:"value"`
	URL   string                 `json:"url"`
}

type FilterRatingConstant struct {
	Value enums.FilterRating `json:"value"`
	Name  string             `json:"name"`
}

type FilterMinOrderConstant struct {
	Value enums.FilterMinOrder `json:"value"`
	Name  string               `json:"name"`
}

type BankTransferInfo struct {
	AccountName string `json:"account_name,omitempty"`

	MainOfficeAddress string `json:"main_office_address,omitempty"`

	Uen string `json:"uen,omitempty"`

	AccountNumber string `json:"account_number,omitempty"`

	MultiCurrencyAccountNumber string `json:"multi_currency_account_number,omitempty"`

	BankName    string `json:"bank_name,omitempty"`
	BankCode    string `json:"bank_code,omitempty"`
	Swift       string `json:"swift,omitempty"`
	BankAddress string `json:"bank_address,omitempty"`

	IntermediaryBank  string `json:"intermediary_bank,omitempty"`
	IntermediarySwift string `json:"intermediary_swift,omitempty"`

	Description string `json:"description"`

	Note string `json:"note"`
}

type BankTransferInfos = map[enums.Currency][]BankTransferInfo
