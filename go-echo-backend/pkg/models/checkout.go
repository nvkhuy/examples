package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/price"

type CheckoutItemsCreateForm []*CheckoutItemCreateForm

type CheckoutItemCreateForm struct {
	Name      string `json:"name" validate:"required"`
	ProductID string `json:"product_id" validate:"required"`
	VariantID string `json:"variant_id" validate:"required"`
	ShopID    string `json:"shop_id" validate:"required"`

	Quantity int64 `json:"quantity" validate:"required"`

	UnitAmount price.Price `json:"unit_amount" validate:"required"`
}

type CheckoutCreateURLForm struct {
	Items           []*CheckoutItemCreateForm `json:"items,omitempty" validate:"required"`
	ReferenceID     string                    `json:"reference_id"`
	PaymentMethodID *string                   `json:"payment_method_id" validate:"required"`
	Currency        string                    `json:"currency" validate:"required"`
	SuccessURL      string                    `json:"success_url,omitempty" validate:"required,startswith=http"`
	CanceledURL     string                    `json:"canceled_url,omitempty" validate:"required,startswith=http"`
}
type CheckoutCreateForm struct {
	Items           CheckoutItemsCreateForm `json:"items,omitempty" validate:"required"`
	ReferenceID     string                  `json:"reference_id"`
	PaymentMethodID *string                 `json:"payment_method_id" validate:"required"`
	Currency        string                  `json:"currency" validate:"required"`
}
type ApplyCouponForm struct {
	Items  []*CheckoutItemCreateForm `json:"items,omitempty" validate:"required"`
	Coupon string                    `json:"coupon" validate:"required"`
}

func (items CheckoutItemsCreateForm) GroupByShop() map[string]CheckoutItemsCreateForm {
	var data = map[string]CheckoutItemsCreateForm{}

	for _, item := range items {
		data[item.ShopID] = append(data[item.ShopID], item)
	}

	return data
}
