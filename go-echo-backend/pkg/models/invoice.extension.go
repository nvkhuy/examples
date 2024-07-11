package models

func (m *Invoice) GetCustomerIOMetadata() map[string]interface{} {
	var result = map[string]interface{}{
		"invoice_number": m.InvoiceNumber,
	}

	if m.DueDate > 0 {
		result["due_date"] = m.DueDate
	}

	if m.IssuedDate > 0 {
		result["issued_date"] = m.IssuedDate
	}

	if m.Currency != "" {
		result["currency"] = m.Currency
		result["currency_code"] = m.Currency.GetCustomerIOCode()
	}

	if m.CountryCode != "" {
		result["country_code"] = m.CountryCode
	}

	if m.Status != "" {
		result["status"] = m.Status
	}

	if m.Note != "" {
		result["note"] = m.Note
	}

	if m.Tax != nil {
		result["tax"] = *m.Tax
	}

	if m.TotalPrice != nil {
		result["total_price"] = *m.TotalPrice
	}

	if m.SubTotal != nil {
		result["sub_total"] = *m.SubTotal
	}

	if m.TaxPercentage != nil {
		result["tax_percentage"] = *m.TaxPercentage
	}

	if m.TransactionFee != nil {
		result["transaction_fee"] = *m.TransactionFee
	}

	if m.ShippingFee != nil {
		result["shipping_fee"] = *m.ShippingFee
	}

	if m.Document != nil {
		m.Document = m.Document.GenerateFileURL()
		result["document"] = m.Document
	}

	return result
}
