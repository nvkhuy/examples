package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type AnalyticProduct struct {
	Model

	URL    string                `gorm:"index:idx_product_url_domain_size_color,unique" json:"url,omitempty"`
	Domain enums.EcommerceDomain `gorm:"index:idx_product_url_domain_size_color,unique" json:"domain,omitempty"`
	Size   string                `gorm:"index:idx_product_url_domain_size_color,unique;default:null" json:"size,omitempty"`
	Color  string                `gorm:"index:idx_product_url_domain_size_color,unique;default:null" json:"color,omitempty"`

	CountryCode           enums.CountryCode  `json:"country_code,omitempty"`
	Name                  string             `json:"name,omitempty"`
	Gender                enums.Gender       `json:"gender,omitempty"`
	Description           string             `json:"description,omitempty"`
	Category              string             `json:"category,omitempty"`
	SubCategory           string             `json:"sub_category,omitempty"`
	Images                pq.StringArray     `gorm:"type:varchar(1000)[]" json:"images,omitempty"`
	PrivateImages         pq.StringArray     `gorm:"type:varchar(1000)[]" json:"private_images,omitempty"`
	Price                 float64            `json:"price,omitempty"`
	Currency              enums.Currency     `json:"currency,omitempty"`
	NumView               int64              `json:"num_view,omitempty"`
	Sold                  int64              `json:"sold,omitempty"`
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
	BestSellerDescription string             `json:"best_seller_description"`
	Trending              string             `json:"trending,omitempty"`
	Tags                  *JsonArrayMetaData `json:"tags,omitempty"`
	CommentOverview       *JsonMetaData      `json:"comment_overview"`
	Metadata              *JsonMetaData      `json:"metadata,omitempty"`

	Score                          float64         `json:"score,omitempty"`
	AverageSalesIncreasePercentage float64         `json:"average_sales_increase_percentage,omitempty"`
	ProductClasses                 []*ProductClass `gorm:"-" json:"product_classes,omitempty"`
	AnalyticProductGrowthRate
}

type AnalyticProductGrowthRate struct {
	OverallGrowthRate float64 `gorm:"" json:"overall_growth_rate"`
	PriceGrowthRate   float64 `gorm:"" json:"price_growth_rate"`
	SoldGrowthRate    float64 `gorm:"" json:"sold_growth_rate"`
	StockGrowthRate   float64 `gorm:"" json:"stock_growth_rate"`
}

func (AnalyticProduct) TableName() string {
	return "products"
}

type AnalyticProductChanges struct {
	Model
	ProductId     string                `gorm:"index:idx_product_changes_product_id" json:"product_id,omitempty"`
	URL           string                `gorm:"index:idx_product_changes_url_scrape_date,unique" json:"url,omitempty"`
	ScrapeDate    string                `gorm:"index:idx_product_changes_url_scrape_date,unique"  json:"scrape_date,omitempty"`
	CountryCode   enums.CountryCode     `json:"country_code,omitempty"`
	Domain        enums.EcommerceDomain `json:"domain,omitempty"`
	Trending      string                `json:"trending,omitempty"`
	Price         float64               `json:"price"`
	Sold          float64               `json:"sold"`
	Stock         float64               `json:"stock"`
	Discount      float64               `json:"discount"`
	RatingAverage float64               `json:"rating_average"`
	ReviewsCount  int                   `json:"reviews_count"`
	Rev           float64               `json:"rev,omitempty"`
	IsPrediction  bool                  `gorm:"-" json:"is_prediction"`
}

func (AnalyticProductChanges) TableName() string {
	return "product_changes"
}

type DataAnalyticPlatformOverviews []*DataAnalyticPlatformOverview

type DataAnalyticPlatformOverview struct {
	Domain      enums.EcommerceDomain `json:"domain,omitempty"`
	CountryCode enums.CountryCode     `json:"country_code,omitempty"`
	Total       int                   `json:"total,omitempty"`
}

type DataAnalyticCategory struct {
	Rank        int                         `json:"rank,omitempty"`
	Category    string                      `json:"category,omitempty"`
	SubCategory string                      `json:"sub_category,omitempty"`
	Sold        int64                       `json:"sold,omitempty"`
	Rev         float64                     `json:"rev,omitempty"`
	Domains     datatypes.JSONSlice[string] `json:"domains,omitempty"`
}

type AnalyticProductClass struct {
	ProductID string    `gorm:"primaryKey" json:"product_id,omitempty"`
	Class     string    `gorm:"primaryKey" json:"class,omitempty"`
	Conf      float64   `json:"conf,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`
}

func (AnalyticProductClass) TableName() string {
	return "product_classes"
}
