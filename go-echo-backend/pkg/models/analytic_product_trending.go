package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/lib/pq"
)

type AnalyticProductTrending struct {
	ID        string     `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	URL    string                `gorm:"uniqueIndex:idx_product_trending" json:"url,omitempty"`
	Domain enums.EcommerceDomain `gorm:"uniqueIndex:idx_product_trending" json:"domain,omitempty"`
	Size   *string               `gorm:"uniqueIndex:idx_product_trending" json:"size,omitempty"`
	Color  *string               `gorm:"uniqueIndex:idx_product_trending" json:"color,omitempty"`

	CountryCode           enums.CountryCode  `json:"country_code,omitempty"`
	Name                  string             `json:"name,omitempty"`
	Gender                enums.Gender       `json:"gender,omitempty"`
	Description           string             `json:"description,omitempty"`
	Category              string             `json:"category,omitempty"`
	SubCategory           string             `json:"sub_category,omitempty"`
	Images                pq.StringArray     `gorm:"type:varchar(1000)[]" json:"images,omitempty"`
	PrivateImages         pq.StringArray     `gorm:"type:varchar(1000)[]" json:"private_images,omitempty"`
	Price                 price.Price        `json:"price,omitempty"`
	Currency              enums.Currency     `json:"currency,omitempty"`
	NumView               int64              `json:"num_view,omitempty"`
	SoldDescription       string             `json:"sold_description,omitempty"`
	Stock                 int64              `json:"stock,omitempty"`
	Discount              float64            `json:"discount,omitempty"`
	NumberOfRating        int64              `json:"number_of_rating,omitempty"`
	RatingDetails         *RatingDetails     `json:"rating_details,omitempty"`
	Origin                string             `json:"origin,omitempty"`
	Material              string             `json:"material,omitempty"`
	Genuine               string             `json:"genuine,omitempty"`
	Brand                 string             `json:"brand,omitempty"`
	BrandOrigin           string             `json:"brand_origin,omitempty"`
	TradeMark             string             `json:"trade_mark,omitempty"`
	IsWarrantyApplied     string             `json:"is_warranty_applied"`
	FabricType            string             `json:"fabric_type,omitempty"`
	FabricCare            string             `json:"fabric_care,omitempty"`
	BestSellerDescription *string            `json:"best_seller_description"`
	Trending              *string            `json:"trending,omitempty"`
	Tags                  *JsonArrayMetaData `json:"tags,omitempty"`
	CommentOverview       *JsonMetaData      `json:"comment_overview"`
	OverallGrowthRate     float64            `json:"overall_growth_rate"`
	PriceGrowthRate       float64            `json:"price_growth_rate"`
	SoldGrowthRate        float64            `json:"sold_growth_rate"`
	StockGrowthRate       float64            `json:"stock_growth_rate"`
	GrowthRateUpdateAt    int64              `json:"growth_rate_update_at,omitempty"`

	Metadata JsonMetaData `json:"metadata,omitempty"`

	// Fastmoss (Tiktok)
	Sold             int64   `json:"sold,omitempty"`
	TotalSoldCount   float64 `json:"total_sold_count,omitempty"`
	TotalSalesAmount float64 `json:"total_sales_amount,omitempty"`
	SalesAmount      float64 `json:"sales_amount,omitempty"`

	// Trendings
	Trendings          []Trending     `gorm:"-" json:"trendings"`
	ProductTrendingIds pq.StringArray `gorm:"-" json:"product_trending_ids"`
	IsPublish          *bool          `gorm:"default:true" json:"is_publish,omitempty"`
}

func (AnalyticProductTrending) TableName() string {
	return "product_trendings"
}

type AnalyticProductTrendingStats struct {
	Domain      enums.Domain `json:"domain,omitempty"`
	Category    string       `json:"category,omitempty"`
	SubCategory string       `json:"sub_category,omitempty"`
	Count       int          `json:"count,omitempty"`
}

func (AnalyticProductTrendingStats) TableName() string {
	return "product_trendings"
}

type AnalyticProductTrendingTags []*AnalyticProductTrendingTag

type AnalyticProductTrendingTag struct {
	ID                string  `json:"id,omitempty"`
	Key               string  `json:"key,omitempty"`
	Value             string  `json:"value,omitempty"`
	OverallGrowthRate float64 `json:"overall_growth_rate"`
}
