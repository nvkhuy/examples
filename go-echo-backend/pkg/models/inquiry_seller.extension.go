package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/price"

func (record *InquirySeller) GetCustomerIOMetadata(extras map[string]interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"id":                record.ID,
		"purchase_order_id": record.PurchaseOrderID,
		"user_id":           record.UserID,
	}

	if record.DueDay != nil {
		result["due_day"] = record.DueDay
	}

	if record.Status != "" {
		result["status"] = record.Status

	}

	if record.DeliveryDate != nil {
		result["delivery_date"] = *record.DeliveryDate
	}

	if record.OfferPrice != nil {
		result["offer_price"] = *record.OfferPrice
	}

	if record.OrderType != "" {
		result["order_type"] = record.OrderType
	}

	if record.OfferRemark != "" {
		result["offer_remark"] = record.OfferRemark
	}

	if record.VarianceAmount != nil {
		result["variance_amount"] = *record.VarianceAmount

	}
	if record.VariancePercentage != nil {
		result["variance_percentage"] = *record.VariancePercentage

	}
	if record.FabricCost != nil {
		result["fabric_cost"] = *record.FabricCost

	}
	if record.DecorationCost != nil {
		result["decoration_cost"] = *record.DecorationCost

	}
	if record.MakingCost != nil {
		result["making_cost"] = *record.MakingCost

	}
	if record.OtherCost != nil {
		result["other_cost"] = *record.OtherCost
	}

	if record.SellerRemark != "" {
		result["seller_remark"] = record.SellerRemark
	}

	if record.BulkQuotations != nil {
		result["bulk_quotations"] = *record.BulkQuotations
	}

	if record.SampleUnitPrice != nil {
		result["sample_unit_price"] = *record.SampleUnitPrice
	}

	if record.SampleLeadTime != nil {
		result["sample_lead_time"] = *record.SampleLeadTime
	}

	if record.AdminSentAt != nil {
		result["admin_sent_at"] = *record.AdminSentAt
	}

	if record.QuotationAt != nil {
		result["quotation_at"] = *record.QuotationAt
	}

	if record.ExpectedStartProductionDate != nil {
		result["expected_start_production_date"] = *record.ExpectedStartProductionDate
	}

	if record.StartProductionDate != nil {
		result["start_production_date"] = *record.StartProductionDate
	}

	if record.CapacityPerDay != nil {
		result["capacity_per_day"] = *record.CapacityPerDay
	}

	if record.Inquiry != nil {
		result["inquiry"] = record.Inquiry.GetCustomerIOMetadata(nil)
	}
	for k, v := range extras {
		result[k] = v
	}

	return result
}

func (record *InquirySeller) GetSampleUnitPrice() (p price.Price) {
	p = price.NewFromFloat(0).AddPtr(record.DecorationCost).AddPtr(record.FabricCost).AddPtr(record.MakingCost).AddPtr(record.OtherCost)

	return
}

func (iqs InquirySellers) IDs() []string {
	var iqIDs = make([]string, 0, len(iqs))
	for _, iq := range iqs {
		iqIDs = append(iqIDs, iq.ID)
	}
	return iqIDs
}

func (iqs InquirySellers) InquiryIDs() []string {
	inquiryIDs := make([]string, 0, len(iqs))
	for _, iq := range iqs {
		inquiryIDs = append(inquiryIDs, iq.InquiryID)
	}
	return inquiryIDs
}
