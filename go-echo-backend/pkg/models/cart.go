package models

type Cart struct {
	Model

	UserID string `json:"user_id"`

	Items []*CartItem `gorm:"-" json:"items,omitempty"`
}

type CartCreateForm struct {
	JwtClaimsInfo

	Items []*CartItemCreateForm `json:"items,omitempty" validate:"required"`
}

type CartItemsCreateForm struct {
	JwtClaimsInfo

	CartID string `param:"cart_id" validate:"required"`

	Items []*CartItemCreateForm `json:"items,omitempty" validate:"required"`
}

type CartItemsRemoveForm struct {
	JwtClaimsInfo

	CartID string `param:"cart_id" validate:"required"`

	ItemIDs []string `json:"item_ids" validate:"required"`
}

type CartUpdateForm struct {
	JwtClaimsInfo

	CartID string                `param:"cart_id" validate:"required"`
	Items  []*CartItemCreateForm `json:"items,omitempty"`
}
