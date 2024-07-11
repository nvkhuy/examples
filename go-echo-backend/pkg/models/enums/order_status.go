package enums

type OrderStatus string

var (
	// OrderStatusDesign         OrderStatus = "design"
	// OrderStatusSampling       OrderStatus = "sampling"
	// OrderStatusBulkProduction OrderStatus = "bulk_production"
	// OrderStatusQc             OrderStatus = "qc"
	// OrderStatusShipping       OrderStatus = "shipping"
	OrderStatusDraft          OrderStatus = "draft"
	OrderStatusWaitingPayment OrderStatus = "waiting_payment"
	OrderStatusConfirmed      OrderStatus = "confirmed"
	OrderStatusPaid           OrderStatus = "paid"
	OrderStatusDelivered      OrderStatus = "delivered"
	OrderStatusProducing      OrderStatus = "producing"
	OrderStatusCancelled      OrderStatus = "cancelled"
)

func (p OrderStatus) String() string {
	return string(p)
}

func (p OrderStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case OrderStatusWaitingPayment:
		name = "Waiting Payment"
	case OrderStatusConfirmed:
		name = "Confirmed"
	case OrderStatusPaid:
		name = "Paid"
	case OrderStatusDelivered:
		name = "Delivered"
	}

	return name
}
