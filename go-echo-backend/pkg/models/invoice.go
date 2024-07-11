package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/datatypes"
)

/*
ALTER SEQUENCE invoices_invoice_number_seq RESTART WITH 1;

SELECT nextval(pg_get_serial_sequence('invoices', 'invoice_number'));
*/
type Invoice struct {
	ID        string     `gorm:"unique" json:"id"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	UserID string `json:"user_id"`

	InvoiceNumber int                 `gorm:"primaryKey;autoIncrement" json:"invoice_number"`
	InvoiceType   enums.InvoiceType   `json:"invoice_type,omitempty"`
	CreatedBy     string              `json:"created_by,omitempty"`
	Document      *Attachment         `json:"document,omitempty"`
	Metadata      InvoiceMetadata     `json:"metadata,omitempty"`
	DueDate       int64               `json:"due_date,omitempty"`
	IssuedDate    int64               `json:"issued_date,omitempty"`
	Currency      enums.Currency      `json:"currency,omitempty"`
	CountryCode   enums.CountryCode   `json:"country_code,omitempty"`
	Status        enums.InvoiceStatus `json:"status" validate:"omitempty,oneof=paid unpaid"`
	Note          string              `json:"note,omitempty"`

	Vendor    *InvoiceParty `json:"vendor,omitempty"`
	Consignee *InvoiceParty `json:"consignee,omitempty"`
	Shipper   *InvoiceParty `json:"shipper,omitempty"`
	Items     InvoiceItems  `json:"items,omitempty"`

	PaymentType          enums.PaymentType   `gorm:"default:'card'" json:"payment_type"`
	PaymentTransactionID string              `json:"payment_transaction_id"`
	PaymentTransaction   *PaymentTransaction `gorm:"-" json:"payment_transaction"`

	InvoicePricing
}

type CreateInvoiceParams struct {
	JwtClaimsInfo

	UserID string `json:"user_id" validate:"required"`

	InvoiceType enums.InvoiceType   `json:"invoice_type,omitempty"`
	CreatedBy   string              `json:"created_by,omitempty"`
	Document    *Attachment         `json:"document,omitempty"`
	Metadata    InvoiceMetadata     `json:"metadata,omitempty"`
	DueDate     int64               `json:"due_date,omitempty"`
	IssuedDate  int64               `json:"issued_date,omitempty"`
	Currency    string              `json:"currency,omitempty"`
	CountryCode string              `json:"country_code,omitempty"`
	Status      enums.InvoiceStatus `json:"status" validate:"omitempty,oneof=paid unpaid"`
	Note        string              `json:"note,omitempty"`

	Vendor    *InvoiceParty `json:"vendor,omitempty"`
	Consignee *InvoiceParty `json:"consignee,omitempty"`
	Shipper   *InvoiceParty `json:"shipper,omitempty"`
	Items     InvoiceItems  `json:"items,omitempty"`

	PaymentType          enums.PaymentType `gorm:"default:'card'" json:"payment_type"`
	PaymentTransactionID string            `json:"payment_transaction_id"`

	InvoicePricing
}

type UpdateInvoiceParams struct {
	JwtClaimsInfo
	InvoiceNumber int `json:"invoice_number,omitempty" param:"invoice_number" query:"invoice_number" form:"invoice_number" validate:"required"`

	InvoiceType enums.InvoiceType   `json:"invoice_type,omitempty"`
	CreatedBy   string              `json:"created_by,omitempty"`
	Document    *Attachment         `json:"document,omitempty"`
	Metadata    datatypes.JSON      `json:"metadata,omitempty"`
	DueDate     int64               `json:"due_date,omitempty"`
	IssuedDate  int64               `json:"issued_date,omitempty"`
	Currency    string              `json:"currency,omitempty"`
	CountryCode string              `json:"country_code,omitempty"`
	Status      enums.InvoiceStatus `json:"status" validate:"omitempty,oneof=paid unpaid"`
	Note        string              `json:"note,omitempty"`
	Tax         float64             `json:"tax,omitempty"`
	Total       float64             `json:"total,omitempty"`

	Vendor    *InvoiceParty `json:"vendor,omitempty"`
	Consignee *InvoiceParty `json:"consignee,omitempty"`
	Shipper   *InvoiceParty `json:"shipper,omitempty"`
	Items     InvoiceItems  `json:"items,omitempty"`
}
