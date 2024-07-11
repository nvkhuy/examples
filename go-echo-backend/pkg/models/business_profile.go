package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
)

// Shop Shop's model
type BusinessProfile struct {
	Model

	OrderTypes                   pq.StringArray     `gorm:"type:varchar(200)[]" json:"order_types" swaggertype:"array,string"`
	ShippingTerms                pq.StringArray     `gorm:"type:varchar(200)[]" json:"shipping_terms,omitempty" swaggertype:"array,string"`
	PoLeadTime                   *int32             `json:"po_lead_time,omitempty"`     // count day
	SampleLeadTime               *int32             `json:"sample_lead_time,omitempty"` // count day
	MonthlyOutput                *float32           `json:"monthly_output,omitempty"`
	MonthlyOutputUnit            enums.OBOutputUnit `json:"monthly_output_unit,omitempty"`
	TotalProductionLine          *int32             `json:"total_production_line,omitempty"`
	TotalWorker                  *int32             `json:"total_worker,omitempty"`
	MCQ                          *int32             `json:"mcq,omitempty"`
	MOQ                          *int32             `json:"moq,omitempty"`
	IncludingDecorationService   *bool              `json:"including_decoration_service,omitempty"`
	IncludingImportExportService *bool              `json:"including_import_export_service,omitempty"`
	CountryCodeOfOrigin          string             `json:"country_code_of_origin,omitempty"`
	MOQTypes                     pq.StringArray     `gorm:"type:varchar(200)[]" json:"moq_types,omitempty" swaggertype:"array,string"`

	ProductCatalogURL         string         `json:"product_catalog_url,omitempty"`
	ProductCatalogAttachments *Attachments   `json:"product_catalog_attachments,omitempty"`
	ProductionLeadTimes       pq.StringArray `gorm:"type:varchar(200)[]" json:"production_lead_times,omitempty"`

	ProductGroups pq.StringArray `gorm:"type:varchar(200)[]" json:"product_groups,omitempty"`

	ProductTypes        pq.StringArray `gorm:"type:varchar(200)[]" json:"product_types,omitempty"`         // options: OBFactoryProductType
	FactoryProductTypes pq.StringArray `gorm:"type:varchar(200)[]" json:"factory_product_types,omitempty"` // TODO: remove. same as product_types

	FabricStock         *bool              `json:"fabric_stock,omitempty"`
	MillFabricTypes     *OBFabricTypeMetas `json:"mill_fabric_types,omitempty"`
	FlatMillFabricTypes pq.StringArray     `gorm:"type:varchar(200)[]" json:"flat_mill_fabric_types,omitempty"`
	ExceptedFabricTypes pq.StringArray     `gorm:"type:varchar(200)[]" json:"excepted_fabric_types,omitempty"` // only for manufacturer seller. Options from: OBFabricType

	DevelopmentServices pq.StringArray `gorm:"type:varchar(200)[]" json:"development_services,omitempty"`
	DecorationServices  pq.StringArray `gorm:"type:varchar(200)[]" json:"decoration_services,omitempty"`

	SewingAccessoryTypes      pq.StringArray `gorm:"type:varchar(200)[]" json:"sewing_accessory_types,omitempty"`
	PackingAccessoryTypes     pq.StringArray `gorm:"type:varchar(200)[]" json:"packing_accessory_types,omitempty"`
	ServiceTypes              pq.StringArray `gorm:"type:varchar(200)[]" json:"service_types,omitempty"`
	InHouseFacilities         pq.StringArray `gorm:"type:varchar(200)[]" json:"in_house_facilities,omitempty"`
	IncludingSecondaryService *bool          `json:"including_secondary_service,omitempty"`
	ProductFocuses            pq.StringArray `gorm:"type:varchar(200)[]" json:"product_focuses,omitempty"`
	DailyOutput               *float32       `json:"daily_output,omitempty"`

	UserID string `gorm:"primaryKey" json:"user_id,omitempty"`
}

type BusinessProfileCreateForm struct {
	JwtClaimsInfo

	OrderTypes                   pq.StringArray     `gorm:"type:varchar(200)[]" json:"order_types" swaggertype:"array,string"`
	ShippingTerms                pq.StringArray     `gorm:"type:varchar(200)[]" json:"shipping_terms" swaggertype:"array,string"`
	PoLeadTime                   *int32             `json:"po_lead_time,omitempty"`     // count day
	SampleLeadTime               *int32             `json:"sample_lead_time,omitempty"` // count day
	MonthlyOutput                *float32           `json:"monthly_output,omitempty"`
	MonthlyOutputUnit            enums.OBOutputUnit `json:"monthly_output_unit,omitempty"`
	TotalProductionLine          *int32             `json:"total_production_line,omitempty"`
	TotalWorker                  *int32             `json:"total_worker,omitempty"`
	MCQ                          *int32             `json:"mcq,omitempty"`
	MOQ                          *int32             `json:"moq,omitempty"`
	IncludingDecorationService   *bool              `json:"including_decoration_service,omitempty"`
	IncludingImportExportService *bool              `json:"including_import_export_service,omitempty"`
	ProductCatalogURL            string             `json:"product_catalog_url,omitempty"`
	ProductCatalogAttachments    *Attachments       `json:"product_catalog_attachments,omitempty"`
	MOQTypes                     pq.StringArray     `gorm:"type:varchar(200)[]" json:"moq_types,omitempty" swaggertype:"array,string"`
	ProductionLeadTimes          pq.StringArray     `gorm:"type:varchar(200)[]" json:"production_lead_times,omitempty"`
	ExceptedFabricTypes          pq.StringArray     `gorm:"type:varchar(200)[]" json:"excepted_fabric_types,omitempty"`
	ProductGroups                pq.StringArray     `gorm:"type:varchar(200)[]" json:"product_groups,omitempty"`
	ProductTypes                 pq.StringArray     `gorm:"type:varchar(200)[]" json:"product_types,omitempty"`
	// FactoryProductTypes          pq.StringArray     `gorm:"type:varchar(200)[]" json:"factory_product_types,omitempty"`
	CountryCodeOfOrigin   string             `json:"country_code_of_origin,omitempty"`
	FabricStock           *bool              `json:"fabric_stock,omitempty"`
	MillFabricTypes       *OBFabricTypeMetas `json:"mill_fabric_types,omitempty"`
	DevelopmentServices   pq.StringArray     `gorm:"type:varchar(200)[]" json:"development_services,omitempty"`
	SewingAccessoryTypes  pq.StringArray     `gorm:"type:varchar(200)[]" json:"sewing_accessory_types,omitempty"`
	PackingAccessoryTypes pq.StringArray     `gorm:"type:varchar(200)[]" json:"packing_accessory_types,omitempty"`

	DecorationServices        pq.StringArray `gorm:"type:varchar(200)[]" json:"decoration_services,omitempty"`
	ServiceTypes              pq.StringArray `gorm:"type:varchar(200)[]" json:"service_types,omitempty"`
	InHouseFacilities         pq.StringArray `json:"in_house_facilities,omitempty"`
	IncludingSecondaryService *bool          `json:"including_secondary_service,omitempty"`
	ProductFocuses            pq.StringArray `json:"product_focuses,omitempty"`
	DailyOutput               *float32       `json:"daily_output,omitempty"`

	User *User `gorm:"-" json:"-"`
}
