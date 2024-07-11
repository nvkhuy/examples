package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
)

type Pricing struct {
	SubTotal               *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"sub_total,omitempty"`
	SubTotalAfterDeduction *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"sub_total_after_deduction,omitempty"`
	ShippingFee            *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"shipping_fee,omitempty"`
	TransactionFee         *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"transaction_fee,omitempty"`
	Tax                    *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"tax,omitempty"`
	TotalPrice             *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"total_price,omitempty"`
	TaxPercentage          *float64     `gorm:"type:decimal(20,4);default:0.0" json:"tax_percentage,omitempty"`
}

type SellerPricing struct {
	SellerTotalPrice *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"seller_total_price,omitempty"`
}

type InvoicePricing struct {
	Pricing

	DepositPaidAmount *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"deposit_total,omitempty"`

	FirstPaymentTransactionFee *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"first_payment_transaction_fee,omitempty"`
	FirstPaymentTax            *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"first_payment_tax,omitempty"`
	FirstPaymentSubTotal       *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"first_payment_sub_total,omitempty"`
	FirstPaymentTotal          *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"first_payment_total,omitempty"`
	FirstPaymentPercentage     float64      `gorm:"default:40.0" json:"first_payment_percentage,omitempty"`

	SecondPaymentTransactionFee *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"second_payment_transaction_fee,omitempty"`
	SecondPaymentTax            *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"second_payment_tax,omitempty"`
	SecondPaymentSubTotal       *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"second_payment_sub_total,omitempty"`
	SecondPaymentTotal          *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"second_payment_total,omitempty"`
	SecondPaymentPercentage     float64      `gorm:"default:40.0" json:"second_payment_percentage,omitempty"`

	FinalPaymentTransactionFee *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"final_payment_transaction_fee,omitempty"`
	FinalPaymentTax            *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"final_payment_tax,omitempty"`
	FinalPaymentSubTotal       *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"final_payment_sub_total,omitempty"`
	FinalPaymentTotal          *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"final_payment_total,omitempty"`
	FinalPaymentPercentage     float64      `gorm:"default:40.0" json:"final_payment_percentage,omitempty"`
}

func (p *Pricing) Add(other Pricing) Pricing {
	var newPricing = Pricing{
		SubTotal:               price.NewFromInt(0).ToPtr(),
		SubTotalAfterDeduction: price.NewFromInt(0).ToPtr(),
		ShippingFee:            price.NewFromInt(0).ToPtr(),
		TransactionFee:         price.NewFromInt(0).ToPtr(),
		Tax:                    price.NewFromInt(0).ToPtr(),
		TotalPrice:             price.NewFromInt(0).ToPtr(),
	}

	newPricing.SubTotal = newPricing.SubTotal.AddPtr(p.SubTotal).AddPtr(other.SubTotal).ToPtr()
	newPricing.SubTotalAfterDeduction = newPricing.SubTotalAfterDeduction.AddPtr(p.SubTotalAfterDeduction).AddPtr(other.SubTotalAfterDeduction).ToPtr()
	newPricing.ShippingFee = newPricing.ShippingFee.AddPtr(p.ShippingFee).AddPtr(other.ShippingFee).ToPtr()
	newPricing.TransactionFee = newPricing.TransactionFee.AddPtr(p.TransactionFee).AddPtr(other.TransactionFee).ToPtr()
	newPricing.Tax = newPricing.Tax.AddPtr(p.Tax).AddPtr(other.Tax).ToPtr()
	newPricing.TotalPrice = newPricing.TotalPrice.AddPtr(p.TotalPrice).AddPtr(other.TotalPrice).ToPtr()

	return newPricing
}
