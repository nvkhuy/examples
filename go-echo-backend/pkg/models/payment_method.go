package models

type UserPaymentMethodCreateForm struct {
	PaymentMethodID string `json:"payment_method_id" validate:"required"`
	IsDefault       bool   `json:"is_default"`
}
