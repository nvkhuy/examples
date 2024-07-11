package models

import (
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

func (record *PurchaseOrder) CalcPrices() {
	record.Tax = record.SubTotal.MultipleFloat64(0.1).ToPtr()
	record.TotalPrice = record.SubTotal.AddPtr(record.Tax).AddPtr(record.ShippingFee).ToPtr()
}

func (record *PurchaseOrder) BeforeCreate(tx *gorm.DB) error {
	if record.ReferenceID == "" {
		var id = helper.GeneratePurchaseOrderReferenceID()
		tx.Statement.SetColumn("ReferenceID", id)
		tx.Statement.AddClauseIfNotExists(clause.OnConflict{
			Columns: []clause.Column{{Name: "reference_id"}},
			DoUpdates: clause.Assignments(func() map[string]interface{} {
				var id = helper.GeneratePurchaseOrderReferenceID()
				return map[string]interface{}{"reference_id": id}
			}()),
		})
	}

	return nil
}

func (record *PurchaseOrder) UpdatePrices() error {
	stripeConfig, err := stripehelper.GetCurrencyConfig(record.Currency)
	if err != nil {
		return err
	}

	record.TransactionFee = price.NewFromFloat(0).ToPtr()
	record.SubTotalAfterDeduction = record.SubTotal

	if record.PaymentType != enums.PaymentTypeBankTransfer {
		record.TransactionFee = record.SubTotalAfterDeduction.AddPtr(record.ShippingFee).MultipleFloat64(0.04).Add(price.NewFromFloat(stripeConfig.AdditionalFee)).ToPtr()
	}

	record.Tax = record.SubTotalAfterDeduction.AddPtr(record.ShippingFee).MultipleFloat64(values.Float64Value(record.TaxPercentage)).DivInt(100).ToPtr()

	record.TotalPrice = record.Tax.AddPtr(record.TransactionFee).AddPtr(record.SubTotal).AddPtr(record.ShippingFee).ToPtr()

	return nil
}

func (record *PurchaseOrder) GetCustomerIOMetadata(extras map[string]interface{}) map[string]interface{} {
	var cfg = config.GetInstance()
	var result = map[string]interface{}{
		"brand_po_url":            fmt.Sprintf("%s/samples/%s", cfg.BrandPortalBaseURL, record.ID),
		"admin_po_url":            fmt.Sprintf("%s/samples/%s/customer", cfg.AdminPortalBaseURL, record.ID),
		"admin_po_url_for_seller": fmt.Sprintf("%s/samples/%s/supplier", cfg.AdminPortalBaseURL, record.ID),
		"id":                      record.ID,
		"reference_id":            record.ReferenceID,
	}

	if len(record.CartItems) > 0 {
		result["cart_itmes"] = record.CartItems
	}

	if record.Currency != "" {
		result["currency"] = record.Currency
		result["currency_code"] = record.Currency.GetCustomerIOCode()
	}

	if record.SubTotal != nil {
		result["sub_total"] = *record.SubTotal
	}

	if record.ShippingFee != nil {
		result["shipping_fee"] = *record.ShippingFee
	}

	if record.TransactionFee != nil {
		result["transaction_fee"] = *record.TransactionFee
	}

	if record.Tax != nil {
		result["tax"] = *record.Tax
	}

	if record.TotalPrice != nil {
		result["total_price"] = *record.TotalPrice
	}

	if record.TaxPercentage != nil {
		result["tax_percentage"] = *record.TaxPercentage
	}

	if record.ReceiptURL != "" {
		result["receipt_url"] = record.ReceiptURL
	}

	if record.PaymentType != "" {
		result["payment_type"] = record.PaymentType
	}

	if record.TrackingStatus != "" {
		result["tracking_status"] = record.TrackingStatus
	}

	if record.TransferedAt != nil {
		result["transfered_at"] = *record.TransferedAt
	}

	if record.MarkAsPaidAt != nil {
		result["mark_as_paid_at"] = *record.MarkAsPaidAt
	}

	if record.MarkAsUnpaidAt != nil {
		result["mark_as_unpaid_at"] = *record.MarkAsUnpaidAt
	}

	if record.ApproveRejectMeta != nil {
		result["approve_reject_meta"] = *record.ApproveRejectMeta
	}

	if record.ReceiverConfirmedAt != nil {
		result["receiver_confirmed_at"] = *record.ReceiverConfirmedAt
	}

	if record.DeliveryStartedAt != nil {
		result["delivery_started_at"] = *record.DeliveryStartedAt
	}

	if record.MakingInfo != nil {
		if record.MakingInfo.Attachments != nil {
			record.MakingInfo.Attachments.GenerateFileURL()
		}

		result["making_info"] = record.MakingInfo
	}

	if record.SubmitInfo != nil {
		if record.SubmitInfo.Attachments != nil {
			record.SubmitInfo.Attachments.GenerateFileURL()
		}

		result["submit_info"] = record.SubmitInfo
	}

	if record.LogisticInfo != nil {
		if record.LogisticInfo.Attachments != nil {
			record.LogisticInfo.Attachments.GenerateFileURL()
		}

		result["logistic_info"] = record.LogisticInfo

	}

	if record.TechpackAttachments != nil {
		result["techpack_attachments"] = record.TechpackAttachments.GenerateFileURL()
	}

	if record.SampleAttachments != nil {
		result["sample_attachments"] = record.SampleAttachments.GenerateFileURL()

	}

	if record.TransactionAttachment != nil {
		result["transaction_attachment"] = record.TransactionAttachment.GenerateFileURL()
	}

	if record.PoRawMaterials != nil {
		result["po_raw_materials"] = record.PoRawMaterials.GenerateFileURL()
	}

	if record.Inquiry != nil {
		result["inquiry"] = record.Inquiry.GetCustomerIOMetadata(nil)
	}

	if record.Assignees != nil {
		result["assignees"] = record.Assignees.GetCustomerIOMetadata(nil)
	}

	if record.Invoice != nil {
		result["invoice"] = record.Invoice.GetCustomerIOMetadata()
	}

	if record.User != nil {
		result["user"] = record.User.GetCustomerIOMetadata(nil)
	}

	if record.PaymentTransaction != nil {
		result["payment_transaction"] = record.PaymentTransaction.GetCustomerIOMetadata(nil)
	}

	for k, v := range extras {
		result[k] = v
	}

	return result
}

func (record PurchaseOrder) GenerateRawMaterialRefID(items *PoRawMaterialMetas) error {
	for _, item := range *items {
		if item.ReferenceID == "" {
			var id = helper.GeneratePoRawMaterialReferenceID()
			item.ReferenceID = id
		}
	}

	return nil
}

func (records PurchaseOrders) ToExcel() ([]byte, error) {
	var data = [][]interface{}{
		{"Reference ID", "Inquiry ID", "Buyer", "Product", "Tracking Status", "Assignee", "Sample Room", "Posted Date"},
	}
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
			func() interface{} {
				if record.SampleMaker != nil {
					if record.SampleMaker.CompanyName != "" {
						return record.SampleMaker.CompanyName
					}
					return record.SampleMaker.Name
				}
				return "Inflow Sample Room"
			}(),
			time.Unix(record.CreatedAt, 0).In(helper.DefaultTimezone.GetLocation()).Format(`Mon. Jan 2 2006 3:04 PM MST-0700`),
		})
	}

	return helper.ToExcel(data)

}

func (pos PurchaseOrders) IDs() []string {
	var IDs []string
	for _, po := range pos {
		IDs = append(IDs, po.ID)
	}
	return IDs
}
func (pos PurchaseOrders) InquiryIDs() []string {
	var inquiryIDs []string
	for _, po := range pos {
		if po.InquiryID != "" {
			inquiryIDs = append(inquiryIDs, po.InquiryID)
		}
	}
	return inquiryIDs
}
func (pos PurchaseOrders) AddressIDs() []string {
	var addressIDs []string
	for _, po := range pos {
		if po.ShippingAddressID != "" {
			addressIDs = append(addressIDs, po.ShippingAddressID)
		}
	}
	return addressIDs
}
