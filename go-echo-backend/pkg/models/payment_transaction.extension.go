package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/runner"
	"github.com/samber/lo"
	"github.com/thaitanloi365/go-utils/values"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (item *PaymentTransaction) BeforeCreate(tx *gorm.DB) error {
	if item.ReferenceID == "" {
		var id = helper.GeneratePaymentTransactionReferenceID()
		tx.Statement.SetColumn("ReferenceID", id)
		tx.Statement.AddClauseIfNotExists(clause.OnConflict{
			Columns: []clause.Column{{Name: "reference_id"}},
			DoUpdates: clause.Assignments(func() map[string]interface{} {
				var id = helper.GeneratePaymentTransactionReferenceID()
				return map[string]interface{}{"reference_id": id}
			}()),
		})
	}

	return nil
}

func (item PaymentTransaction) GetCustomerIOMetadata(extras map[string]interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"id":           item.ID,
		"reference_id": item.ReferenceID,
	}

	if item.PaymentType != "" {
		result["payment_type"] = item.PaymentType

	}
	if item.Currency != "" {
		result["currency"] = item.Currency
		result["currency_code"] = item.Currency.GetCustomerIOCode()

	}
	if item.PaidAmount != nil {
		result["paid_amount"] = *item.PaidAmount

	}
	if item.TotalAmount != nil {
		result["total_amount"] = *item.TotalAmount

	}
	if item.PaymentPercentage != nil {
		result["payment_percentage"] = *item.PaymentPercentage

	}
	if item.BalanceAmount != nil {
		result["balance_amount"] = *item.BalanceAmount

	}
	if item.Remark != "" {
		result["remark"] = item.Remark

	}

	if item.Milestone != "" {
		result["milestone"] = item.Milestone
	}

	if item.TransactionRefID != "" {
		result["transaction_ref_id"] = item.TransactionRefID
	}

	if item.Note != "" {
		result["note"] = item.Note
	}

	if item.PaymentLinkID != "" {
		result["payment_link_id"] = item.PaymentLinkID
	}

	if item.TransactionType != "" {
		result["transaction_type"] = item.TransactionType
	}

	if item.Metadata != nil {
		result["metadata"] = *item.Metadata
	}

	if item.Attachments != nil {
		result["attachments"] = item.Attachments.GenerateFileURL()
	}
	for k, v := range extras {
		result[k] = v
	}

	return result
}

func (records PaymentTransactions) ToExcel() ([]byte, error) {
	var cfg = config.GetInstance()

	var data = make([][]interface{}, len(records)+1)
	data[0] = []interface{}{"Product Name", "Reference ID", "User", "PO or Bulk ID", "Payment Type", "Milestone", "Payment Details", "Paid Amount", "Total Amount", "Paid Date"}

	var runner = runner.New(10)
	defer runner.Release()

	for index, record := range records {
		index := index + 1
		record := record

		runner.Submit(func() {
			data[index] = []interface{}{
				func(xlsx *excelize.File, sheetName, cell string) error {
					var productName = ""
					if record.PurchaseOrder != nil && record.PurchaseOrder.Inquiry != nil && record.PurchaseOrder.Inquiry.Title != "" {
						productName = record.PurchaseOrder.Inquiry.Title
						var display = fmt.Sprintf("%s/inquiries/%s/customer", cfg.AdminPortalBaseURL, record.PurchaseOrder.Inquiry.ID)

						xlsx.SetCellStr(sheetName, cell, productName)
						return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
							Display: &display,
							Tooltip: &productName,
						})
					}

					if record.BulkPurchaseOrder != nil && record.BulkPurchaseOrder.ProductName != "" {
						productName = record.BulkPurchaseOrder.ProductName
						var display = fmt.Sprintf("%s/bulks/%s/customer", cfg.AdminPortalBaseURL, record.BulkPurchaseOrder.ID)

						xlsx.SetCellStr(sheetName, cell, productName)
						return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
							Display: &display,
							Tooltip: &productName,
						})
					}

					return xlsx.SetCellStr(sheetName, cell, productName)
				},
				func(xlsx *excelize.File, sheetName, cell string) error {
					var display = fmt.Sprintf("%s/payments/%s/overview", cfg.AdminPortalBaseURL, record.ID)

					xlsx.SetCellStr(sheetName, cell, record.ReferenceID)
					return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
						Display: &display,
						Tooltip: &record.ReferenceID,
					})
				},
				func(xlsx *excelize.File, sheetName, cell string) error {
					var name = ""
					if record.User != nil {
						name = record.User.Name
						var display = fmt.Sprintf("%s/users/%s/overview", cfg.AdminPortalBaseURL, record.User.ID)

						xlsx.SetCellStr(sheetName, cell, name)
						return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
							Display: &display,
							Tooltip: &name,
						})
					}

					return xlsx.SetCellStr(sheetName, cell, name)
				},
				func(xlsx *excelize.File, sheetName, cell string) error {
					if record.PurchaseOrder != nil {
						var display = fmt.Sprintf("%s/samples/%s/customer", cfg.AdminPortalBaseURL, record.PurchaseOrder.ID)

						xlsx.SetCellStr(sheetName, cell, record.PurchaseOrder.ReferenceID)
						return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
							Display: &display,
							Tooltip: &record.PurchaseOrder.ReferenceID,
						})
					}

					if record.BulkPurchaseOrder != nil {
						var display = fmt.Sprintf("%s/bulks/%s/customer", cfg.AdminPortalBaseURL, record.BulkPurchaseOrder.ID)

						xlsx.SetCellStr(sheetName, cell, record.BulkPurchaseOrder.ReferenceID)
						return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
							Display: &display,
							Tooltip: &record.BulkPurchaseOrder.ReferenceID,
						})
					}

					if len(record.PurchaseOrders) > 0 {
						var ids = lo.Map(record.PurchaseOrders, func(item *PurchaseOrder, index int) string {
							return item.ReferenceID
						})

						if len(record.PurchaseOrders) == 1 {
							var display = fmt.Sprintf("%s/samples/%s/customer", cfg.AdminPortalBaseURL, record.PurchaseOrders[0].ID)

							xlsx.SetCellStr(sheetName, cell, record.PurchaseOrders[0].ReferenceID)
							return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
								Display: &display,
								Tooltip: &record.PurchaseOrders[0].ReferenceID,
							})
						}

						var idsStr = strings.Join(ids, "\n")
						var display = fmt.Sprintf("%s/payments/%s/overview", cfg.AdminPortalBaseURL, record.ID)
						xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
							Display: &display,
							Tooltip: &idsStr,
						})

						return xlsx.SetCellStr(sheetName, cell, idsStr)
					}

					return xlsx.SetCellStr(sheetName, cell, "")
				},
				func(xlsx *excelize.File, sheetName, cell string) error {
					var displayName = record.PaymentType.DisplayName()
					if record.PaymentType == enums.PaymentTypeCard && record.PaymentIntentID != "" {
						var display = fmt.Sprintf("%s/payments/%s", cfg.StripeDashboardURL, record.PaymentIntentID)

						xlsx.SetCellStr(sheetName, cell, displayName)
						return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
							Display: &display,
							Tooltip: &displayName,
						})

					}

					if record.PaymentType == enums.PaymentTypeBankTransfer && record.Attachments != nil && len(*record.Attachments) > 0 {
						var display = (*record.Attachments)[0].GenerateFileURL().FileURL
						xlsx.SetCellStr(sheetName, cell, displayName)
						return xlsx.SetCellHyperLink(sheetName, cell, display, "External", excelize.HyperlinkOpts{
							Display: &display,
							Tooltip: &displayName,
						})

					}

					return xlsx.SetCellValue(sheetName, cell, displayName)
				},
				func(xlsx *excelize.File, sheetName, cell string) error {
					var displayName = record.Milestone.DisplayName()
					return xlsx.SetCellValue(sheetName, cell, displayName)
				},
				func(xlsx *excelize.File, sheetName, cell string) error {
					var paymentDetails []string

					if record.PurchaseOrder != nil {
						paymentDetails = append(paymentDetails, fmt.Sprintf("Sample SubTotal: %s", record.PurchaseOrder.SubTotal.FormatMoney(record.Currency)))
						paymentDetails = append(paymentDetails, fmt.Sprintf("Sample Transaction Fee: %s", record.PurchaseOrder.TransactionFee.FormatMoney(record.Currency)))
						paymentDetails = append(paymentDetails, fmt.Sprintf("Sample Shipping Fee: %s", record.PurchaseOrder.ShippingFee.FormatMoney(record.Currency)))
						paymentDetails = append(paymentDetails, fmt.Sprintf("Sample Tax (%s): %s", fmt.Sprintf("%0.f", values.Float64Value(record.PurchaseOrder.TaxPercentage))+"%", record.PurchaseOrder.Tax.FormatMoney(record.Currency)))
						paymentDetails = append(paymentDetails, "")
						paymentDetails = append(paymentDetails, fmt.Sprintf("Sample Total: %s", record.PurchaseOrder.TotalPrice.FormatMoney(record.Currency)))
					} else if len(record.PurchaseOrders) > 0 {
						if len(record.PurchaseOrders) == 1 {
							paymentDetails = append(paymentDetails, fmt.Sprintf("Sample SubTotal: %s", record.PurchaseOrders[0].SubTotal.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, fmt.Sprintf("Sample Transaction Fee: %s", record.PurchaseOrders[0].TransactionFee.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, fmt.Sprintf("Sample Shipping Fee: %s", record.PurchaseOrders[0].ShippingFee.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, fmt.Sprintf("Sample Tax (%s): %s", fmt.Sprintf("%0.f", values.Float64Value(record.PurchaseOrders[0].TaxPercentage))+"%", record.PurchaseOrders[0].Tax.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, "")
							paymentDetails = append(paymentDetails, fmt.Sprintf("Sample Total: %s", record.PurchaseOrders[0].TotalPrice.FormatMoney(record.Currency)))
						} else {
							for _, record := range record.PurchaseOrders {
								paymentDetails = append(paymentDetails, fmt.Sprintf("Sample %s SubTotal: %s", record.ReferenceID, record.SubTotal.FormatMoney(record.Currency)))
								paymentDetails = append(paymentDetails, fmt.Sprintf("Sample %s Transaction Fee: %s", record.ReferenceID, record.TransactionFee.FormatMoney(record.Currency)))
								paymentDetails = append(paymentDetails, fmt.Sprintf("Sample %s Shipping Fee: %s", record.ReferenceID, record.ShippingFee.FormatMoney(record.Currency)))
								paymentDetails = append(paymentDetails, fmt.Sprintf("Sample %s Tax (%s): %s", record.ReferenceID, fmt.Sprintf("%0.f", values.Float64Value(record.TaxPercentage))+"%", record.Tax.FormatMoney(record.Currency)))
								paymentDetails = append(paymentDetails, "")
								paymentDetails = append(paymentDetails, fmt.Sprintf("Sample %s Total: %s", record.ReferenceID, record.TotalPrice.FormatMoney(record.Currency)))
							}
						}

					} else if record.BulkPurchaseOrder != nil {
						if record.Milestone == enums.PaymentMilestoneDeposit || record.Milestone == enums.PaymentMilestoneFinalPayment {
							if record.BulkPurchaseOrder.DepositPaidAmount.GreaterThan(0) {
								paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk deposit amount: %s", record.BulkPurchaseOrder.DepositPaidAmount.FormatMoney(record.Currency)))
								paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk deposit note: %s", record.BulkPurchaseOrder.DepositNote))
							}
						}

						if record.Milestone == enums.PaymentMilestoneFirstPayment || record.Milestone == enums.PaymentMilestoneFinalPayment {
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk 1st Percentage: %s", fmt.Sprintf("%0.f", values.Float64Value(record.BulkPurchaseOrder.FirstPaymentPercentage))+"%"))
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk 1st SubTotal: %s", record.BulkPurchaseOrder.FirstPaymentSubTotal.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk 1st Transaction Fee: %s", record.BulkPurchaseOrder.FirstPaymentTransactionFee.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, "")
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk 1st Total: %s", record.BulkPurchaseOrder.FirstPaymentTotal.FormatMoney(record.Currency)))
						}

						if record.Milestone == enums.PaymentMilestoneFinalPayment {
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk Final Percentage: %s", fmt.Sprintf("%0.f", 100-values.Float64Value(record.BulkPurchaseOrder.FirstPaymentPercentage))+"%"))
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk Final SubTotal: %s", record.BulkPurchaseOrder.FinalPaymentSubTotal.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk Final Transaction Fee: %s", record.BulkPurchaseOrder.FinalPaymentTransactionFee.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk Final Tax (%s): %s", fmt.Sprintf("%0.f", values.Float64Value(record.BulkPurchaseOrder.TaxPercentage))+"%", record.BulkPurchaseOrder.FinalPaymentTax.FormatMoney(record.Currency)))
							paymentDetails = append(paymentDetails, "")
							paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk Final Total: %s", record.BulkPurchaseOrder.FinalPaymentTotal.FormatMoney(record.Currency)))
						}

						paymentDetails = append(paymentDetails, "\n")
						paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk SubTotal: %s", record.BulkPurchaseOrder.SubTotal.FormatMoney(record.Currency)))
						paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk Tax (%s): %s", fmt.Sprintf("%0.f", values.Float64Value(record.BulkPurchaseOrder.TaxPercentage))+"%", record.BulkPurchaseOrder.Tax.FormatMoney(record.Currency)))
						paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk Shipping Fee: %s", record.BulkPurchaseOrder.ShippingFee.FormatMoney(record.Currency)))
						paymentDetails = append(paymentDetails, fmt.Sprintf("Bulk Total: %s", record.BulkPurchaseOrder.TotalPrice.FormatMoney(record.Currency)))
					}

					return xlsx.SetCellStr(sheetName, cell, strings.Join(paymentDetails, "\n"))
				},
				func() interface{} {
					if record.PaidAmount != nil {
						return record.PaidAmount.FormatMoney(record.Currency)
					}
					return ""
				}(),
				func() interface{} {
					if record.TotalAmount != nil {
						return record.TotalAmount.FormatMoney(record.Currency)
					}
					return ""
				}(),
				func() interface{} {
					var createdAtTime = time.Unix(record.CreatedAt, 0)
					if record.MarkAsPaidAt != nil {
						createdAtTime = time.Unix(*record.MarkAsPaidAt, 0)
					}
					return createdAtTime.In(helper.DefaultTimezone.GetLocation()).Format(`Mon. Jan 2 2006 3:04 PM MST-0700`)
				}(),
			}
		})

	}

	runner.Wait()

	return helper.ToExcel(data)
}
