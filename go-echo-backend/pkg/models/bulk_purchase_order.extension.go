package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/samber/lo"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (record *BulkPurchaseOrder) BeforeSave(tx *gorm.DB) error {
	if len(record.AdditionalItems) > 0 {
		record.AdditionalItems.UpdateTotal()
	}
	return nil
}

func (record *BulkPurchaseOrder) BeforeUpdate(tx *gorm.DB) error {
	if len(record.AdditionalItems) > 0 {
		record.AdditionalItems.UpdateTotal()
	}
	return nil
}

func (record *BulkPurchaseOrder) BeforeCreate(tx *gorm.DB) error {
	if record.ReferenceID == "" {
		var id = helper.GenerateBulkPurchaseOrderReferenceID()
		tx.Statement.SetColumn("ReferenceID", id)
		tx.Statement.AddClauseIfNotExists(clause.OnConflict{
			Columns: []clause.Column{{Name: "reference_id"}},
			DoUpdates: clause.Assignments(func() map[string]interface{} {
				var id = helper.GenerateBulkPurchaseOrderReferenceID()
				return map[string]interface{}{"reference_id": id}
			}()),
		})
	}

	return nil
}

func (record *BulkPurchaseOrder) UpdatePrices() error {
	stripeConfig, err := stripehelper.GetCurrencyConfig(record.Currency)
	if err != nil {
		return err
	}
	if stripeConfig == nil {
		err = errors.New("nil stripe config")
		return err
	}

	// With invoice included
	var finalPaymentTax = record.TaxPercentage
	if record.CommercialInvoice != nil {
		record.ShippingFee = record.CommercialInvoice.ShippingFee
		record.SubTotalAfterDeduction = record.GetInvoiceSubTotal().ToPtr()
		finalPaymentTax = record.CommercialInvoice.TaxPercentage

	} else {
		record.SubTotalAfterDeduction = record.SubTotal
	}

	for _, item := range record.AdditionalItems {
		item.TotalPrice = item.UnitPrice.MultipleInt(item.Qty).ToPtr()
		record.SubTotal = record.SubTotal.AddPtr(item.TotalPrice).ToPtr()
		record.SubTotalAfterDeduction = record.SubTotalAfterDeduction.AddPtr(item.TotalPrice).ToPtr()
	}

	record.FirstPaymentSubTotal = record.GetFirstPaymentAmount().ToPtr()
	if record.FirstPaymentType != enums.PaymentTypeBankTransfer {
		record.FirstPaymentTransactionFee = record.FirstPaymentSubTotal.
			MultipleFloat64(stripeConfig.TransactionFee).
			Add(price.NewFromFloat(stripeConfig.AdditionalFee)).
			ToPtr()
	} else {
		record.FirstPaymentTransactionFee = price.NewFromFloat(0).ToPtr()
	}
	record.FirstPaymentTotal = record.FirstPaymentSubTotal.
		AddPtr(record.FirstPaymentTransactionFee).
		AddPtr(record.FirstPaymentTax).
		ToPtr()

	record.FinalPaymentSubTotal = record.SubTotalAfterDeduction.
		SubPtr(record.DepositPaidAmount).
		SubPtr(record.FirstPaymentSubTotal).
		SubPtr(record.SecondPaymentSubTotal).
		ToPtr()

	record.FinalPaymentDeductionAmount = record.SubTotal.Sub(record.SubTotalAfterDeduction.ToValue()).ToPtr()
	if record.FinalPaymentType != enums.PaymentTypeBankTransfer {
		record.FinalPaymentTransactionFee = record.FinalPaymentSubTotal.
			AddPtr(record.ShippingFee).
			MultipleFloat64(stripeConfig.TransactionFee).
			Add(price.NewFromFloat(stripeConfig.AdditionalFee)).
			ToPtr()
	} else {
		record.FinalPaymentTransactionFee = price.NewFromFloat(0).ToPtr()
	}

	record.FinalPaymentTax = record.SubTotalAfterDeduction.
		AddPtr(record.ShippingFee).
		MultipleFloat64(values.Float64Value(finalPaymentTax)).
		DivInt(100).
		SubPtr(record.FirstPaymentTax).
		ToPtr()
	record.FinalPaymentTotal = record.FinalPaymentTax.
		AddPtr(record.FinalPaymentTransactionFee).
		AddPtr(record.ShippingFee).
		AddPtr(record.FinalPaymentSubTotal).
		SubPtr(record.SampleDeductionAmount).
		ToPtr()

	record.TransactionFee = record.FirstPaymentTransactionFee.Add(record.FinalPaymentTransactionFee.ToValue()).ToPtr()
	record.Tax = price.NewFromFloat(0).AddPtr(record.FirstPaymentTax).Add(record.FinalPaymentTax.ToValue()).ToPtr()
	record.TotalPrice = record.FirstPaymentTotal.Add(record.FinalPaymentTotal.ToValue()).ToPtr()

	if record.DepositPaidAmount != nil {
		record.TotalPrice = record.TotalPrice.Add(record.DepositPaidAmount.ToValue()).ToPtr()
	}

	if record.DepositPaidAmount != nil {
		record.TotalPrice = record.TotalPrice.Add(record.DepositPaidAmount.ToValue()).ToPtr()
	}

	if record.DepositPaidAmount != nil {
		record.TotalPrice = record.TotalPrice.Add(record.DepositPaidAmount.ToValue()).ToPtr()
	}

	if record.DepositPaidAmount != nil {
		record.TotalPrice = record.TotalPrice.Add(record.DepositPaidAmount.ToValue()).ToPtr()
	}

	return nil
}

func (order BulkPurchaseOrder) GetQuotedPrice() (price price.Price) {
	if order.AdminQuotations != nil && len(order.AdminQuotations) > 0 {
		firstQuote := order.AdminQuotations[0]
		return firstQuote.Price
	}

	return
}

func (order BulkPurchaseOrder) GetBulkQuotation() *InquiryQuotationItem {
	if order.AdminQuotations != nil && len(order.AdminQuotations) > 0 {
		return order.AdminQuotations[0]
	}

	return nil
}

func (order BulkPurchaseOrder) GetFinalPaymentAmount() price.Price {
	return (order.TotalPrice.Sub(price.NewFromPtr(order.FirstPaymentTotal))).Abs()
}

func (order BulkPurchaseOrder) GetQuotationLeadTime() *int64 {
	if order.AdminQuotations != nil && len(order.AdminQuotations) > 0 {
		leadTime := order.AdminQuotations[0]
		return leadTime.LeadTime
	}

	return values.Int64(0)
}

func (order BulkPurchaseOrder) GetFirstPaymentAmount() price.Price {
	var percentage = values.Float64Value(order.FirstPaymentPercentage)
	if percentage >= 1 && percentage <= 100 {
		return order.SubTotal.MultipleFloat64(percentage).DivInt(100)
	}

	return order.SubTotal.MultipleFloat64(0)
}

func (order BulkPurchaseOrder) GetBalanceAmountAfterFirstPayment() price.Price {
	return order.TotalPrice.Sub(price.NewFromPtr(order.FirstPaymentTotal))
}

func (order BulkPurchaseOrder) GetInvoiceSubTotal() price.Price {
	var total = price.NewFromFloat(0)
	if order.CommercialInvoice != nil && order.CommercialInvoice.SubTotal != nil && order.CommercialInvoice.SubTotal.GreaterThan(0) {
		return order.CommercialInvoice.SubTotal.ToValue()
	}

	for _, item := range order.CommercialInvoice.Items {
		total = total.Add(item.TotalAmount)
	}

	return total
}

func (order BulkPurchaseOrder) GenerateRawMaterialRefID(items *PoRawMaterialMetas) error {
	for _, item := range *items {
		if item.ReferenceID == "" {
			var id = helper.GeneratePoRawMaterialReferenceID()
			item.ReferenceID = id
		}
	}

	return nil
}

func (order *BulkPurchaseOrder) GetCustomerIOMetadata(extras map[string]interface{}) map[string]interface{} {
	var cfg = config.GetInstance()
	var result = map[string]interface{}{
		"brand_bulk_po_url": fmt.Sprintf("%s/bulks/%s", cfg.BrandPortalBaseURL, order.ID),
		"admin_bulk_po_url": fmt.Sprintf("%s/bulks/%s/overview", cfg.AdminPortalBaseURL, order.ID),
		"id":                order.ID,
		"reference_id":      order.ReferenceID,
	}

	if order.Status != "" {
		result["status"] = order.Status
	}

	if order.Currency != "" {
		result["currency"] = order.Currency
		result["currency_customerio_code"] = order.Currency.GetCustomerIOCode()
	}

	if order.TransactionFee != nil {
		result["transaction_fee"] = *order.TransactionFee
	}

	if order.Tax != nil {
		result["tax"] = *order.Tax
	}

	if order.TotalPrice != nil {
		result["total_price"] = *order.TotalPrice
	}

	if order.SubTotal != nil {
		result["sub_total"] = *order.SubTotal
	}

	if order.TaxPercentage != nil {
		result["tax_percentage"] = *order.TaxPercentage
	}

	if order.ShippingFee != nil {
		result["shipping_fee"] = *order.ShippingFee
	}

	if order.FirstPaymentType != "" {
		result["first_payment_payment_type"] = order.FirstPaymentType
	}

	if order.FirstPaymentTax != nil {
		result["first_payment_tax"] = *order.FirstPaymentTax
	}

	if order.FirstPaymentTransactionFee != nil {
		result["first_payment_transaction_fee"] = *order.FirstPaymentTransactionFee
	}

	if order.FirstPaymentReceivedAt != nil {
		result["first_payment_received_at"] = *order.FirstPaymentReceivedAt
	}

	if order.FirstPaymentSubTotal != nil {
		result["first_payment_sub_total"] = *order.FirstPaymentSubTotal
	}

	if order.FirstPaymentTotal != nil {
		result["first_payment_total"] = *order.FirstPaymentTotal
	}

	if order.FirstPaymentPercentage != nil {
		result["first_payment_percentage"] = order.FirstPaymentPercentage
	}

	if order.FirstPaymentTransferedAt != nil {
		result["first_payment_transfered_at"] = *order.FirstPaymentTransferedAt
	}

	if order.FirstPaymentMarkAsPaidAt != nil {
		result["first_payment_mark_as_paid_at"] = *order.FirstPaymentMarkAsPaidAt
	}

	if order.FirstPaymentMarkAsUnpaidAt != nil {
		result["first_payment_mark_as_unpaid_at"] = *order.FirstPaymentMarkAsUnpaidAt
	}

	if order.FirstPaymentIntentID != "" {
		result["first_payment_intent_id"] = order.FirstPaymentIntentID
	}

	if order.FirstPaymentChargeID != "" {
		result["first_payment_charge_id"] = order.FirstPaymentChargeID
	}

	if order.FinalPaymentTax != nil {
		result["final_payment_tax"] = *order.FinalPaymentTax
	}

	if order.FinalPaymentTransactionFee != nil {
		result["final_payment_transaction_fee"] = *order.FinalPaymentTransactionFee
	}

	if order.FinalPaymentIntentID != "" {
		result["final_payment_intent_id"] = order.FinalPaymentIntentID
	}

	if order.FinalPaymentChargeID != "" {
		result["final_payment_charge_id"] = order.FinalPaymentChargeID
	}

	if order.FinalPaymentReceivedAt != nil {
		result["final_payment_received_at"] = *order.FinalPaymentReceivedAt
	}

	if order.FinalPaymentTotal != nil {
		result["final_payment_total"] = order.FinalPaymentTotal
	}

	if order.FinalPaymentTransferedAt != nil {
		result["final_payment_transfered_at"] = *order.FinalPaymentTransferedAt

	}
	if order.FinalPaymentMarkAsPaidAt != nil {
		result["final_payment_mark_as_paid_at"] = *order.FinalPaymentMarkAsPaidAt

	}
	if order.FinalPaymentMarkAsUnpaidAt != nil {
		result["final_payment_mark_as_unpaid_at"] = *order.FinalPaymentMarkAsUnpaidAt
	}

	if order.QuotationLeadTime != nil {
		result["quotation_lead_time"] = *order.QuotationLeadTime
	}

	if order.ReceiverConfirmedAt != nil {
		result["receiver_confirmed_at"] = *order.ReceiverConfirmedAt
	}

	if order.QuotationAt != nil {
		result["quotation_at"] = *order.QuotationAt
	}

	if order.DeliveryStartedAt != nil {
		result["delivery_started_at"] = *order.DeliveryStartedAt
	}

	if order.SubmittedAt != nil {
		result["submitted_at"] = *order.SubmittedAt
	}

	if order.DeliveredAt != nil {
		result["delivered_at"] = *order.DeliveredAt
	}

	if order.FinalPaymentType != "" {
		result["final_payment_payment_type"] = order.FinalPaymentType
	}

	if order.TrackingStatus != "" {
		result["tracking_status"] = order.TrackingStatus
	}

	if order.QuotationNote != "" {
		result["quotation_note"] = order.QuotationNote
	}

	if order.PoQcReports != nil {
		result["po_qc_reports"] = order.PoQcReports.GenerateFileURL()
	}

	if order.PpsInfo != nil {
		result["pps_info"] = order.PpsInfo.GenerateFileURL()
	}

	if order.LogisticInfo != nil {
		result["logistic_info"] = order.LogisticInfo.GenerateFileURL()
	}

	if order.ProductionInfo != nil {
		result["production_info"] = order.ProductionInfo.GenerateFileURL()
	}

	if order.Inquiry != nil {
		result["inquiry"] = order.Inquiry.GetCustomerIOMetadata(nil)
	}

	if order.PackingAttachments != nil && len(*order.PackingAttachments) > 0 {
		result["packing_attachments"] = order.PackingAttachments.GenerateFileURL()

	}

	if order.Attachments != nil && len(*order.Attachments) > 0 {
		result["requirements"] = order.Attachments.GenerateFileURL()

	}

	if order.ShippingAttachments != nil && len(*order.ShippingAttachments) > 0 {
		result["shipping_attachments"] = order.ShippingAttachments.GenerateFileURL()
	}

	if order.CommercialInvoiceAttachment != nil {
		result["commercial_invoice_attachment"] = order.CommercialInvoiceAttachment.GenerateFileURL()
	}

	if order.PoRawMaterials != nil && len(*order.PoRawMaterials) > 0 {
		result["po_raw_materials"] = order.PoRawMaterials.GenerateFileURL()
	}

	if order.Assignees != nil {
		result["assignees"] = order.Assignees.GetCustomerIOMetadata(nil)
	}

	for k, v := range extras {
		result[k] = v
	}

	return result
}

func (bulkPO *BulkPurchaseOrder) IsFirstPaymentPaid() bool {
	return bulkPO.FirstPaymentIntentID != "" || (bulkPO.FirstPaymentType == enums.PaymentTypeBankTransfer && bulkPO.FirstPaymentTransactionRefID != "")
}

func (bulkPO *BulkPurchaseOrder) IsSecondPaymentPaid() bool {
	return bulkPO.SecondPaymentIntentID != "" || (bulkPO.SecondPaymentType == enums.PaymentTypeBankTransfer && bulkPO.SecondPaymentTransactionRefID != "")
}

func (bulkPO *BulkPurchaseOrder) IsFinalPaymentPaid() bool {
	return bulkPO.FinalPaymentIntentID != "" || (bulkPO.FinalPaymentType == enums.PaymentTypeBankTransfer && bulkPO.FinalPaymentTransactionRefID != "")
}

func (records BulkPurchaseOrders) ToExcel() ([]byte, error) {
	var data = [][]interface{}{
		{"Reference ID", "Inquiry ID", "Buyer", "Product", "Tracking Status", "Assignee", "Posted Date"},
	}
	var sb strings.Builder
	sb.WriteString("id,user,product,tracking status,assignee,created date\n")
	for _, record := range records {
		data = append(data, []interface{}{
			record.ReferenceID,
			func() string {
				if record.Inquiry != nil {
					return record.Inquiry.ReferenceID
				}
				return record.InquiryID
			}(),
			func() string {
				if record.User != nil {
					return record.User.Name
				}
				return ""
			}(),
			func() string {
				if record.Inquiry != nil {
					return record.Inquiry.Title
				}
				return ""
			}(),
			record.TrackingStatus.DisplayName(),
			func() interface{} {
				if len(record.Assignees) > 0 {
					var names = lo.Map(record.Assignees, func(item *User, index int) string {
						return item.Name
					})

					return strings.Join(names, ",")
				}
				return nil
			}(),
			time.Unix(record.CreatedAt, 0).In(helper.DefaultTimezone.GetLocation()).Format(`Mon. Jan 2 2006 3:04 PM MST-0700`),
		})
	}

	return helper.ToExcel(data)
}

func (bpos BulkPurchaseOrders) IDs() []string {
	var IDs []string
	for _, bpo := range bpos {
		IDs = append(IDs, bpo.ID)
	}
	return IDs
}
func (bpos BulkPurchaseOrders) InquiryIDs() []string {
	var inquiryIDs []string
	for _, bpo := range bpos {
		if bpo.InquiryID != "" {
			inquiryIDs = append(inquiryIDs, bpo.InquiryID)
		}
	}
	return inquiryIDs
}
func (bpos BulkPurchaseOrders) PurchaseOrderIDs() []string {
	var purchaseOrderIDs []string
	for _, bpo := range bpos {
		if bpo.PurchaseOrderID != "" {
			purchaseOrderIDs = append(purchaseOrderIDs, bpo.PurchaseOrderID)
		}
	}
	return purchaseOrderIDs
}

func (bpos BulkPurchaseOrders) AddressIDs() []string {
	var addressIDs []string
	for _, bpo := range bpos {
		if bpo.ShippingAddressID != "" {
			addressIDs = append(addressIDs, bpo.ShippingAddressID)
		}
	}
	return addressIDs
}
