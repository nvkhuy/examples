package models

import "fmt"

func (tna *TNA) GetLink() string {
	if tna.Inquiry != nil && tna.Inquiry.ID != "" {
		return fmt.Sprintf("/inquiries/%s", tna.Inquiry.ID)
	}
	if tna.PurchaseOrder != nil && tna.PurchaseOrder.ID != "" {
		return fmt.Sprintf("/samples/%s", tna.PurchaseOrder.ID)
	}
	if tna.BulkPurchaseOrder != nil && tna.BulkPurchaseOrder.ID != "" {
		return fmt.Sprintf("/bulks/%s", tna.BulkPurchaseOrder.ID)
	}
	return ""
}

func (tna *TNA) GetCustomerIOMetadata(extras map[string]interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"id":   tna.ID,
		"link": tna.GetLink(),
	}

	if tna.DateFrom > 0 {
		result["date_from"] = tna.DateFrom
	}

	if tna.DateTo > 0 {
		result["date_to"] = tna.DateTo
	}

	if tna.ReferenceID != "" {
		result["reference_id"] = tna.ReferenceID
	}

	if tna.Title != "" {
		result["title"] = tna.Title
	}

	if tna.Comment != "" {
		result["comment"] = tna.Comment
	}

	if tna.OrderType != "" {
		result["order_type"] = tna.OrderType
	}

	if len(tna.AssigneeIDs) > 0 {
		result["assignee_ids"] = tna.AssigneeIDs
	}

	if tna.Inquiry != nil {
		result["inquiry"] = tna.Inquiry.GetCustomerIOMetadata(nil)
	}

	if tna.PurchaseOrder != nil {
		result["purchase_order"] = tna.PurchaseOrder.GetCustomerIOMetadata(nil)
	}

	if tna.BulkPurchaseOrder != nil {
		result["bulk_purchase_order"] = tna.BulkPurchaseOrder.GetCustomerIOMetadata(nil)
	}

	for k, v := range extras {
		result[k] = v
	}

	return result
}
