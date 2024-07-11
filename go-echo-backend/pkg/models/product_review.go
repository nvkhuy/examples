package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

// "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

// ProductReview's model
type ProductReview struct {
	Model

	RatingStar  int          `json:"rating_star"`
	Status      string       `json:"status"`
	Body        string       `json:"body"`
	ProductID   string       `json:"product_id"`
	OrderID     string       `json:"order_id"`
	UserID      string       `json:"user_id"`
	ShopID      string       `json:"shop_id"`
	User        *User        `gorm:"-" json:"user"`
	Attachments *Attachments `gorm:"-" json:"attachments,omitempty"`
}

type ProductReviewCreateForm struct {
	RatingStar int        `json:"rating_star"`
	Body       string     `json:"body"`
	ProductID  string     `json:"product_id"`
	OrderID    string     `json:"order_id"`
	UserID     string     `json:"user_id"`
	ForRole    enums.Role `json:"-"`
}
