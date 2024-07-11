package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/lib/pq"
)

// Product model
type Product struct {
	Model

	Name             string            `gorm:"size:1000" json:"name,omitempty"`
	Slug             string            `gorm:"index:idx_slug,unique" json:"slug"`
	ShortDescription string            `gorm:"size:5000" json:"short_description,omitempty"`
	Description      string            `gorm:"size:5000" json:"description,omitempty"`
	Vi               *ProductContent   `json:"vi,omitempty"`
	QRCode           string            `gorm:"size:200" json:"qr_code,omitempty"`
	CategoryID       string            `gorm:"size:100" json:"category_id,omitempty"`
	Category         *Category         `gorm:"-"  json:"category"`
	Sku              string            `gorm:"size:100" json:"sku,omitempty"`
	Price            price.Price       `gorm:"type:decimal(20,4);default:0.0" json:"price"`
	SoldQuantity     int               `json:"sold_quantity,omitempty"`
	ProductType      enums.ProductType `json:"product_type,omitempty" gorm:"default:clothing;size:100"` // clothing, fabric, graphic....

	BulletPoints    int               `json:"bullet_points"`
	SpecialFeatures string            `gorm:"size:500" json:"special_features"`
	Material        string            `gorm:"size:500" json:"material"`
	Gender          string            `gorm:"size:500" json:"gender"`
	Style           string            `gorm:"size:500" json:"style"`
	CountryCode     enums.CountryCode `gorm:"size:100" json:"country_code"`
	SafetyInfo      string            `gorm:"size:500" json:"safety_info"`
	Currency        enums.Currency    `gorm:"size:100" json:"currency"`

	ReadyToShip bool `gorm:"default:false" json:"ready_to_ship,omitempty"`
	DailyDeal   bool `gorm:"default:false" json:"daily_deal,omitempty"`

	RatingCount int     `json:"rating_count,omitempty"`
	RatingStar  float32 `json:"rating_star,omitempty"`

	TradeUnit enums.ProductUnit `gorm:"default:piece" json:"trade_unit,omitempty"` // piece, pairs, boxes
	MinOrder  int               `gorm:"default:0" json:"min_order,omitempty"`

	Attachments *Attachments   `json:"attachments,omitempty"`
	FabricIDs   pq.StringArray `gorm:"type:varchar(200)[]" json:"fabric_ids,omitempty"`

	SourceProductID string       `gorm:"default:null;unique;size:100" json:"source_product_id"`
	Source          enums.Source `gorm:"default:'inflow'" json:"source"`

	ProductAttributeMetas ProductAttributeMetas `json:"product_attribute_metas,omitempty"`

	Variants       []*Variant      `gorm:"-" json:"variants,omitempty"`
	Fabrics        []*Fabric       `gorm:"-" json:"fabrics,omitempty"`
	IsTrending     *bool           `gorm:"default:false" json:"is_trending,omitempty"`
	ProductClasses []*ProductClass `gorm:"-" json:"product_classes,omitempty"`
}

type Products []*Product

// ProductResponse Product's model
type ProductResponse struct {
	ID               string       `json:"id,omitempty"`
	Name             string       `json:"name,omitempty"`
	Description      string       `json:"description,omitempty"`
	ShortDescription string       `json:"short_description,omitempty"`
	Sku              string       `json:"sku"`
	SoldQuantity     int          `json:"sold_quantity"`
	Price            price.Price  `json:"price"`
	Attachments      *Attachments `json:"attachments"`
	ReadyToShip      bool         `json:"ready_to_ship"`
	ProductType      string       `json:"product_type"`
}

type ProductCreateForm struct {
	JwtClaimsInfo

	Name               string                         `json:"name,omitempty" validate:"required"`
	Description        string                         `json:"description,omitempty"`
	ShortDescription   string                         `json:"short_description,omitempty"`
	Vi                 *ProductContent                `json:"vi,omitempty"`
	ShopID             string                         `json:"shop_id"`
	CategoryID         string                         `json:"category_id,omitempty"`
	Sku                string                         `json:"sku"`
	Price              price.Price                    `json:"price"`
	ProductType        enums.ProductType              `json:"product_type"` // clothing, fabric, graphic....
	BulletPoints       int                            `json:"bullet_points"`
	SpecialFeatures    string                         `json:"special_features"`
	Material           string                         `json:"material"`
	Gender             string                         `json:"gender"`
	Style              string                         `json:"style"`
	CountryCode        enums.CountryCode              `json:"country_code"`
	SafetyInfo         string                         `json:"safety_info"`
	Currency           enums.Currency                 `json:"currency"`
	Variants           []*VariantAttributeUpdateForm  `json:"variants,omitempty"`
	QuantityPriceTiers []*QuantityPriceTierUpdateForm `json:"quantity_price_tiers,omitempty"`
	Attachments        *Attachments                   `json:"attachments,omitempty"`
	FabricIDs          pq.StringArray                 `gorm:"type:varchar(200)[]" json:"fabric_ids,omitempty"`

	ReadyToShip bool `json:"ready_to_ship"`
	DailyDeal   bool `json:"daily_deal"`

	RatingCount int     `json:"rating_count,omitempty"`
	RatingStar  float32 `json:"rating_star,omitempty"`

	TradeUnit enums.ProductUnit `json:"trade_unit,omitempty"` // piece, pairs, boxes
	MinOrder  int               `json:"min_order,omitempty"`

	ForRole enums.Role `json:"-"`
}

type ProductUpdateForm struct {
	JwtClaimsInfo
	ProductID string `param:"product_id" validate:"required"`

	Name               string                         `json:"name,omitempty" validate:"required"`
	Description        string                         `json:"description,omitempty"`
	ShortDescription   string                         `json:"short_description,omitempty"`
	Vi                 *ProductContent                `json:"vi,omitempty"`
	ShopID             string                         `json:"shop_id,omitempty"`
	CategoryID         string                         `json:"category_id,omitempty"`
	Sku                string                         `json:"sku,omitempty"`
	Price              price.Price                    `json:"price"`
	ProductType        enums.ProductType              `json:"product_type,omitempty"` // clothing, fabric, graphic....
	BulletPoints       int                            `json:"bullet_points,omitempty"`
	SpecialFeatures    string                         `json:"special_features,omitempty"`
	Material           string                         `json:"material,omitempty"`
	Gender             string                         `json:"gender,omitempty"`
	Style              string                         `json:"style,omitempty"`
	CountryCode        enums.CountryCode              `json:"country_code,omitempty"`
	SafetyInfo         string                         `json:"safety_info,omitempty"`
	Currency           enums.Currency                 `json:"currency,omitempty"`
	Variants           []*VariantAttributeUpdateForm  `json:"variants,omitempty,omitempty"`
	QuantityPriceTiers []*QuantityPriceTierUpdateForm `json:"quantity_price_tiers,omitempty"`
	Attachments        *Attachments                   `json:"attachments,omitempty"`
	FabricIDs          pq.StringArray                 `gorm:"type:varchar(200)[]" json:"fabric_ids,omitempty"`

	ReadyToShip bool `json:"ready_to_ship,omitempty"`
	DailyDeal   bool `json:"daily_deal,omitempty"`

	RatingCount int     `json:"rating_count,omitempty"`
	RatingStar  float32 `json:"rating_star,omitempty"`

	TradeUnit enums.ProductUnit `json:"trade_unit,omitempty"` // piece, pairs, boxes
	MinOrder  int               `json:"min_order,omitempty"`
}

// ProductSearchResponse Product's model
type ProductSearchResponse struct {
	// Model

	Products []*Product `json:"products,omitempty"`
}

// ProductDetailWithVariant Product's model
type ProductDetailWithVariant struct {
	// Model

	*Product
	Variants           []*Variant           `json:"variants,omitempty"`
	Options            []*ProductAttribute  `json:"options,omitempty"`
	QuantityPriceTiers []*QuantityPriceTier `json:"quantity_price_tiers,omitempty"`
}

type GenerateProductSlugParams struct {
	JwtClaimsInfo
}

type GetProductQRCodeParams struct {
	JwtClaimsInfo
	Logo      string
	Bucket    string
	ProductId string `json:"product_id" query:"product_id" form:"product_id" validate:"required"`
	URL       string `json:"url" query:"url" form:"url"`
	Override  bool   `json:"override" query:"override" form:"override"`
}

type ExportProductsResponse struct {
	FileKey string `json:"file_key"`
}
