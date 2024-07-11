package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type PaymentRequest struct {
	Model

	Type enums.PaymentRequestType `json:"type,omitempty"`

	Items PaymentRequestItems `json:"items,omitempty"`

	Pricing
}
