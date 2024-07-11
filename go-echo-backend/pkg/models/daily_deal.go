package models

// Faq Faq's model
type DailyDeal struct {
	Model

	Status       string `json:"status,omitempty"`
	ProductID    string `json:"product_id,omitempty"`
	DealFrom     *int64 `json:"deal_from,omitempty"`
	DealTo       *int64 `json:"deal_to,omitempty"`
	DiscountType string `json:"discount_type,omitempty"`
	DealValue    int    `json:"deal_value,omitempty"`
}
