package models

func (c *Comment) GetCustomerIOMetadata(extras map[string]interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"id":        c.ID,
		"content":   c.Content,
		"target_id": c.TargetID,
		"type":      c.TargetType,
	}
	if c.Attachments != nil {
		result["attachments"] = c.Attachments.GenerateFileURL()
	}

	if c.PurchaseOrder != nil {
		result["purchase_order"] = c.PurchaseOrder.GetCustomerIOMetadata(nil)
	}

	if c.Inquiry != nil {
		result["inquiry"] = c.Inquiry.GetCustomerIOMetadata(nil)
	}

	for k, v := range extras {
		result[k] = v
	}

	return result

}
