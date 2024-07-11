package repo

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/pdf"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvoiceRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewInvoiceRepo(db *db.DB) *InvoiceRepo {
	return &InvoiceRepo{
		db:     db,
		logger: logger.New("repo/Page"),
	}
}

func (r *InvoiceRepo) CreateInvoice(params models.CreateInvoiceParams) (*models.Invoice, error) {
	var invoice models.Invoice
	var err = copier.Copy(&invoice, &params)
	if err != nil {
		return nil, err
	}

	invoice.CreatedBy = params.GetUserID()
	err = r.db.Create(&invoice).Error
	if err != nil {
		return nil, err
	}

	return &invoice, err
}

func (r *InvoiceRepo) UpdateInvoice(params models.UpdateInvoiceParams) (*models.Invoice, error) {
	var find models.Invoice
	_ = r.db.Select("InvoiceNumber", "Status").First(&find, "invoice_number = ?", params.InvoiceNumber).Error
	if find.InvoiceNumber == 0 {
		return nil, errors.New("not found invoice")
	}
	if find.Status == enums.InvoiceStatusPaid {
		return nil, errors.New("cannot update paid invoice")
	}

	var updateInvoice models.Invoice
	var err = copier.Copy(&updateInvoice, &params)
	if err != nil {
		return nil, err
	}
	updateInvoice.CreatedBy = params.GetUserID()

	err = r.db.Model(&models.Invoice{}).Where("invoice_number = ?", params.InvoiceNumber).Updates(&updateInvoice).Error
	if err != nil {
		return nil, err
	}

	return &updateInvoice, err
}

type PaginateInvoicesParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	UserID   string   `json:"user_id" query:"user_id" param:"user_id"`
	Statuses []string `json:"statuses" query:"statuses" param:"statuses"`
}

func (r *InvoiceRepo) PaginateInvoices(params PaginateInvoicesParams) (result *query.Pagination) {
	result = query.New(r.db, queryfunc.NewInvoiceBuilder(queryfunc.InvoiceBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			if params.GetRole().IsAdmin() {
				if params.UserID != "" {
					builder.Where("iv.user_id = ?", params.UserID)
				}
			} else {
				builder.Where("iv.user_id = ?", params.GetUserID())
			}

			if len(params.Statuses) > 0 {
				builder.Where("iv.status IN ?", params.Statuses)
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + " %"
				if strings.HasPrefix(keyword, "PO-") {
					builder.Where("iv.metadata->>'purchase_order_reference_id' = ?", q)
				} else if strings.HasPrefix(keyword, "BPO-") {
					builder.Where("iv.metadata->>'bulk_purchase_order_reference_id' = ?", q)
				} else if strings.HasPrefix(keyword, "IQ-") {
					builder.Where("iv.metadata->>'inquiry_reference_id' = ?", q)
				} else {
					builder.Where("iv.consignee->>'email' ILIKE ? OR iv.consignee->>'name' ILIKE ? OR iv.invoice_number = ?", q, q, keyword)
				}
			}
		}).
		WithoutCount(params.WithoutCount).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return
}

type InvoiceDetailsPrams struct {
	models.JwtClaimsInfo
	InvoiceNumber *int `json:"invoice_number" param:"invoice_number" query:"invoice_number"`
}

func (r *InvoiceRepo) Details(params InvoiceDetailsPrams) (result models.Invoice, err error) {
	err = query.New(r.db, queryfunc.NewInvoiceBuilder(queryfunc.InvoiceBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			if params.InvoiceNumber != nil {
				builder.Where("invoice_number = ?", *params.InvoiceNumber)
			}
		}).
		Limit(1).
		FirstFunc(&result)
	return
}

func (r *InvoiceRepo) IsExitsInvoiceNumber(params InvoiceDetailsPrams) (is bool, err error) {
	var resultID int
	if err = r.db.Model(&models.Invoice{}).Select("invoice_number").
		First(&resultID, "invoice_number = ?", params.InvoiceNumber).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		err = nil
	}
	if resultID == 0 {
		return false, nil
	}
	return true, nil
}

type NextInvoiceNumberParams struct {
	models.JwtClaimsInfo
}

func (r *InvoiceRepo) NextInvoiceNumber() (next int, err error) {
	var resultID int
	if err = r.db.Model(&models.Invoice{}).Select("invoice_number").Order("invoice_number DESC").Limit(1).First(&resultID).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		err = nil
	}
	next = resultID + 1
	return
}

type CreateBulkFinalPaymentInvoiceParams struct {
	models.JwtClaimsInfo

	Bulk                *models.BulkPurchaseOrder `json:"-"`
	BulkPurchaseOrderID string                    `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool                      `json:"re_create"`
	Invoice             *models.Invoice           `json:"-"`
}

func (r *InvoiceRepo) CreateBulkFinalPaymentInvoice(params CreateBulkFinalPaymentInvoiceParams) (*models.BulkPurchaseOrder, error) {
	var bulkPO = params.Bulk
	var err error
	if bulkPO == nil {
		bulkPO, err = NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
			BulkPurchaseOrderID: params.BulkPurchaseOrderID,
			IncludeUser:         true,
			IncludeInvoice:      true,
			IncludeItems:        true,
		})
	}
	if err != nil {
		return nil, err
	}

	var invoice = params.Invoice
	if params.Invoice == nil {
		if bulkPO.FinalPaymentInvoice != nil && !params.ReCreate {
			r.db.CustomLogger.Debugf("Bulk PO %s final payment invoice is already generated", bulkPO.ReferenceID)
			return bulkPO, errs.ErrBulkPoInvoiceAlreadyGenerated
		}

		var invoiceParams = models.CreateInvoiceParams{
			UserID:      bulkPO.UserID,
			InvoiceType: enums.InvoiceTypeBulkPOFinalPayment,
			Vendor:      models.DefaultVendorForOnlinePayment,
			IssuedDate:  bulkPO.CommercialInvoice.IssuedDate,
			DueDate:     bulkPO.CommercialInvoice.DueDate,
			CountryCode: bulkPO.CommercialInvoice.CountryCode,
			Consignee:   bulkPO.CommercialInvoice.Consignee,
			Shipper:     bulkPO.CommercialInvoice.Shipper,
			Status:      enums.InvoiceStatusPaid,
			Currency:    bulkPO.CommercialInvoice.Currency,
			// Note:                 fmt.Sprintf("Final payment, %.f", bulkPO.FirstPaymentPercentage) + "% of total payment" + "\n" + bulkPO.CommercialInvoice.Note,
			PaymentType:          bulkPO.FinalPaymentType,
			PaymentTransactionID: bulkPO.FinalPaymentIntentID,
			InvoicePricing: models.InvoicePricing{
				Pricing: models.Pricing{
					SubTotal:               bulkPO.FinalPaymentSubTotal,
					TransactionFee:         bulkPO.FinalPaymentTransactionFee,
					Tax:                    bulkPO.FinalPaymentTax,
					TotalPrice:             bulkPO.FinalPaymentTotal,
					SubTotalAfterDeduction: bulkPO.SubTotalAfterDeduction,
					TaxPercentage:          bulkPO.CommercialInvoice.TaxPercentage,
					ShippingFee:            bulkPO.CommercialInvoice.ShippingFee,
				},
				FirstPaymentTransactionFee:  bulkPO.FirstPaymentTransactionFee,
				FirstPaymentTax:             bulkPO.FirstPaymentTax,
				FirstPaymentSubTotal:        bulkPO.FirstPaymentSubTotal,
				FirstPaymentTotal:           bulkPO.FirstPaymentTotal,
				FirstPaymentPercentage:      values.Float64Value(bulkPO.FirstPaymentPercentage),
				SecondPaymentTransactionFee: bulkPO.SecondPaymentTransactionFee,
				SecondPaymentTax:            bulkPO.SecondPaymentTax,
				SecondPaymentSubTotal:       bulkPO.SecondPaymentSubTotal,
				SecondPaymentTotal:          bulkPO.SecondPaymentTotal,
				SecondPaymentPercentage:     bulkPO.SecondPaymentPercentage,
				FinalPaymentTransactionFee:  bulkPO.FinalPaymentTransactionFee,
				FinalPaymentTax:             bulkPO.FinalPaymentTax,
				FinalPaymentSubTotal:        bulkPO.FinalPaymentSubTotal,
				FinalPaymentTotal:           bulkPO.FinalPaymentTotal,
			},
			Metadata: models.InvoiceMetadata{
				InvoiceType:                  enums.InvoiceTypeBulkPOFinalPayment,
				InquiryID:                    bulkPO.Inquiry.ID,
				InquiryReferenceID:           bulkPO.Inquiry.ReferenceID,
				BulkPurchaseOrderID:          bulkPO.ID,
				BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
				BulkPurchaseOrderCommercialInvoiceAttachment: bulkPO.CommercialInvoiceAttachment,
				BulkEstimatedProductionLeadTime:              bulkPO.GetQuotationLeadTime(),
			},
		}

		if bulkPO.FinalPaymentType == enums.PaymentTypeBankTransfer && bulkPO.Currency == enums.VND {
			invoiceParams.Vendor = models.DefaultVendorForLocal
		}

		if bulkPO.FinalPaymentType == enums.PaymentTypeBankTransfer {
			invoiceParams.PaymentTransactionID = bulkPO.FinalPaymentTransactionRefID
		}

		for _, item := range bulkPO.CommercialInvoice.Items {
			var invoiceItem = &models.InvoiceItem{
				ItemCode:          item.ID,
				Color:             item.Color,
				Size:              item.Size,
				Description:       fmt.Sprintf("%s-%s", item.Color, item.SizeName),
				TotalQuantity:     item.TotalQuantity,
				UnitPrice:         item.UnitPrice.ToPtr(),
				TotalAmount:       item.TotalAmount.ToPtr(),
				FabricComposition: bulkPO.Inquiry.Composition,
			}
			if item.Size != nil {
				invoiceItem.Description = fmt.Sprintf("%s-%s", item.Color, item.Size.GetSizeDescription())
			}

			invoiceParams.Items = append(invoiceParams.Items, invoiceItem)
		}

		if bulkPO.User != nil {
			invoiceParams.Consignee = &models.InvoiceParty{
				ID:          bulkPO.User.ID,
				Name:        bulkPO.User.Name,
				Email:       bulkPO.User.Email,
				PhoneNumber: bulkPO.User.PhoneNumber,
				CompanyName: bulkPO.User.CompanyName,
			}

			if bulkPO.User.Coordinate == nil && bulkPO.User.CoordinateID != "" {
				var coordinate models.Coordinate
				if err := r.db.First(&coordinate, "id = ?", bulkPO.User.CoordinateID).Error; err == nil {
					invoiceParams.Consignee.Address = coordinate.Display()
				}
			}
		}

		invoice, err = r.CreateInvoice(invoiceParams)
		if err != nil {
			return nil, eris.Wrapf(err, "Generate pdf bulk PO %s", bulkPO.ReferenceID)
		}

	}

	data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
		URL:               fmt.Sprintf("%s/invoices/print/%d", r.db.Configuration.AdminPortalBaseURL, invoice.InvoiceNumber),
		Selector:          "#invoice-ready-to-print",
		Landscape:         true,
		PrintBackground:   true,
		PreferCssPageSize: true,
	})
	if err != nil {
		return nil, eris.Wrapf(err, "Generate pdf bulk PO %s", bulkPO.ReferenceID)
	}

	if len(data) == 0 {
		return nil, eris.Errorf("Generate pdf bulk PO %s empty data", bulkPO.ReferenceID)
	}

	var uploadParams = s3.UploadFileParams{
		Key:         fmt.Sprintf("uploads/%s_%s_%d_final_payment_invoice.pdf", bulkPO.Inquiry.ReferenceID, bulkPO.ReferenceID, invoice.InvoiceNumber),
		Data:        bytes.NewBuffer(data),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(models.ContentTypePDF),
		ACL:         "private",
	}

	_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
	if err != nil {
		return nil, err
	}

	var attachment = models.Attachment{
		ContentType: uploadParams.ContentType,
		FileKey:     uploadParams.Key,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if bulkPO.FinalPaymentInvoice != nil && params.ReCreate {
			err = tx.Delete(&models.Invoice{}, "invoice_number = ?", bulkPO.FinalPaymentInvoice.InvoiceNumber).Error
			if err != nil {
				return err
			}
		}

		err = tx.Model(&models.Invoice{}).Where("invoice_number = ?", invoice.InvoiceNumber).UpdateColumn("Document", &attachment).Error
		if err != nil {
			return err
		}

		bulkPO.FinalPaymentInvoice = invoice
		bulkPO.FinalPaymentInvoice.Document = &attachment

		return tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).UpdateColumn("FinalPaymentInvoiceNumber", invoice.InvoiceNumber).Error
	})
	if err != nil {
		return nil, err
	}

	return bulkPO, err
}

type CreateBulkDepositInvoiceParams struct {
	models.JwtClaimsInfo
	BulkPurchaseOrderID string          `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool            `json:"re_create"`
	Invoice             *models.Invoice `json:"-"`
}

func (r *InvoiceRepo) CreateBulkDepositInvoice(params CreateBulkDepositInvoiceParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
		IncludeItems:        true,
		IncludeInvoice:      true,
	})
	if err != nil {
		return nil, err
	}

	var invoice = params.Invoice

	if invoice == nil {
		if bulkPO.FirstPaymentInvoice != nil && !params.ReCreate {
			r.db.CustomLogger.Debugf("Bulk PO %s first payment invoice is already generated", bulkPO.ReferenceID)
			return bulkPO, errs.ErrBulkPoInvoiceAlreadyGenerated
		}

		var invoiceParams = models.CreateInvoiceParams{
			UserID:               bulkPO.UserID,
			InvoiceType:          enums.InvoiceTypeBulkPODepositPayment,
			Status:               enums.InvoiceStatusPaid,
			Currency:             string(bulkPO.Currency),
			Vendor:               models.DefaultVendorForOnlinePayment,
			PaymentType:          bulkPO.FirstPaymentType,
			PaymentTransactionID: bulkPO.FirstPaymentIntentID,
			InvoicePricing: models.InvoicePricing{
				Pricing: models.Pricing{
					SubTotal:               bulkPO.SubTotal,
					TransactionFee:         bulkPO.TransactionFee,
					Tax:                    bulkPO.Tax,
					TotalPrice:             bulkPO.TotalPrice,
					TaxPercentage:          bulkPO.TaxPercentage,
					SubTotalAfterDeduction: bulkPO.FirstPaymentSubTotal,
				},
				DepositPaidAmount:           bulkPO.DepositPaidAmount,
				FirstPaymentTransactionFee:  bulkPO.FirstPaymentTransactionFee,
				FirstPaymentTax:             bulkPO.FirstPaymentTax,
				FirstPaymentSubTotal:        bulkPO.FirstPaymentSubTotal,
				FirstPaymentTotal:           bulkPO.FirstPaymentTotal,
				SecondPaymentTransactionFee: bulkPO.SecondPaymentTransactionFee,
				SecondPaymentTax:            bulkPO.SecondPaymentTax,
				SecondPaymentSubTotal:       bulkPO.SecondPaymentSubTotal,
				SecondPaymentTotal:          bulkPO.SecondPaymentTotal,
				FirstPaymentPercentage:      values.Float64Value(bulkPO.FirstPaymentPercentage),
				FinalPaymentTransactionFee:  bulkPO.FinalPaymentTransactionFee,
				FinalPaymentTax:             bulkPO.FinalPaymentTax,
				FinalPaymentSubTotal:        bulkPO.FinalPaymentSubTotal,
				FinalPaymentTotal:           bulkPO.FinalPaymentTotal,
			},
			Metadata: models.InvoiceMetadata{
				InvoiceType:                     enums.InvoiceTypeBulkPODepositPayment,
				InquiryID:                       bulkPO.Inquiry.ID,
				InquiryReferenceID:              bulkPO.Inquiry.ReferenceID,
				BulkPurchaseOrderID:             bulkPO.ID,
				BulkPurchaseOrderReferenceID:    bulkPO.ReferenceID,
				BulkEstimatedProductionLeadTime: bulkPO.GetQuotationLeadTime(),
			},
		}

		if bulkPO.FirstPaymentType == enums.PaymentTypeBankTransfer && bulkPO.Currency == enums.VND {
			invoiceParams.Vendor = models.DefaultVendorForLocal
		}

		if bulkPO.FirstPaymentType == enums.PaymentTypeBankTransfer {
			invoiceParams.PaymentTransactionID = bulkPO.FirstPaymentTransactionRefID
		}

		if bulkPO.FirstPaymentReceivedAt != nil {
			invoiceParams.IssuedDate = *bulkPO.FirstPaymentReceivedAt
		}

		if bulkPO.FirstPaymentMarkAsPaidAt != nil {
			invoiceParams.IssuedDate = *bulkPO.FirstPaymentMarkAsPaidAt
		}

		if bulkPO.User != nil {
			invoiceParams.Consignee = &models.InvoiceParty{
				ID:          bulkPO.User.ID,
				Name:        bulkPO.User.Name,
				Email:       bulkPO.User.Email,
				PhoneNumber: bulkPO.User.PhoneNumber,
				CompanyName: bulkPO.User.CompanyName,
			}
			if bulkPO.User.Coordinate == nil && bulkPO.User.CoordinateID != "" {
				var coordinate models.Coordinate
				if err := r.db.First(&coordinate, "id = ?", bulkPO.User.CoordinateID).Error; err == nil {
					invoiceParams.Consignee.Address = coordinate.Display()
				}
			}
		}

		invoiceParams.Items = append(invoiceParams.Items, &models.InvoiceItem{
			ItemCode:      bulkPO.ReferenceID,
			Description:   bulkPO.DepositNote,
			TotalQuantity: 1,
			UnitPrice:     bulkPO.DepositPaidAmount,
			TotalAmount:   bulkPO.DepositPaidAmount,
		})

		invoice, err = r.CreateInvoice(invoiceParams)
		if err != nil {
			return nil, eris.Wrapf(err, "Generate pdf bulk PO %s", bulkPO.ReferenceID)
		}
	}

	data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
		URL:               fmt.Sprintf("%s/invoices/print/%d", r.db.Configuration.AdminPortalBaseURL, invoice.InvoiceNumber),
		Selector:          "#invoice-ready-to-print",
		Landscape:         true,
		PrintBackground:   true,
		PreferCssPageSize: true,
	})
	if err != nil {
		return nil, eris.Wrapf(err, "Generate pdf bulk PO %s", bulkPO.ReferenceID)
	}

	if len(data) == 0 {
		return nil, eris.Errorf("Generate pdf bulk PO %s empty data", bulkPO.ReferenceID)
	}

	var uploadParams = s3.UploadFileParams{
		Key:         fmt.Sprintf("uploads/%s_%s_%d_deposit_invoice.pdf", bulkPO.Inquiry.ReferenceID, bulkPO.ReferenceID, invoice.InvoiceNumber),
		Data:        bytes.NewBuffer(data),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(models.ContentTypePDF),
		ACL:         "private",
	}

	_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
	if err != nil {
		return nil, err
	}

	var attachment = models.Attachment{
		ContentType: uploadParams.ContentType,
		FileKey:     uploadParams.Key,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if bulkPO.DepositInvoice != nil && params.ReCreate {
			err = tx.Delete(&models.Invoice{}, "invoice_number = ?", bulkPO.DepositInvoice.InvoiceNumber).Error
			if err != nil {
				return err
			}
		}
		err = tx.Model(&models.Invoice{}).Where("invoice_number = ?", invoice.InvoiceNumber).UpdateColumn("Document", &attachment).Error
		if err != nil {
			return err
		}

		bulkPO.DepositInvoice = invoice
		bulkPO.DepositInvoice.Document = &attachment

		return tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).UpdateColumn("DepositInvoiceNumber", invoice.InvoiceNumber).Error
	})
	if err != nil {
		return nil, err
	}

	return bulkPO, err
}

type CreateBulkFirstPaymentInvoiceParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string          `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool            `json:"re_create"`
	Invoice             *models.Invoice `json:"-"`
}

func (r *InvoiceRepo) CreateBulkFirstPaymentInvoice(params CreateBulkFirstPaymentInvoiceParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
		IncludeItems:        true,
		IncludeInvoice:      true,
	})
	if err != nil {
		return nil, err
	}

	var invoice = params.Invoice

	if invoice == nil {
		if bulkPO.FirstPaymentInvoice != nil && bulkPO.FirstPaymentInvoice.Document != nil && !params.ReCreate {
			r.db.CustomLogger.Debugf("Bulk PO %s first payment invoice is already generated", bulkPO.ReferenceID)
			return bulkPO, errs.ErrBulkPoInvoiceAlreadyGenerated
		}

		var invoiceParams = models.CreateInvoiceParams{
			UserID:      bulkPO.UserID,
			InvoiceType: enums.InvoiceTypeBulkPOFirstPayment,
			Status:      enums.InvoiceStatusPaid,
			Currency:    string(bulkPO.Currency),
			Vendor:      models.DefaultVendorForOnlinePayment,
			// Note:                 fmt.Sprintf("First payment, %.f", bulkPO.FirstPaymentPercentage) + "% of total payment",
			PaymentType:          bulkPO.FirstPaymentType,
			PaymentTransactionID: bulkPO.FirstPaymentIntentID,
			InvoicePricing: models.InvoicePricing{
				Pricing: models.Pricing{
					SubTotal:               bulkPO.SubTotal,
					TransactionFee:         bulkPO.TransactionFee,
					Tax:                    bulkPO.Tax,
					TotalPrice:             bulkPO.TotalPrice,
					TaxPercentage:          bulkPO.TaxPercentage,
					SubTotalAfterDeduction: bulkPO.FirstPaymentSubTotal,
				},
				FirstPaymentTransactionFee: bulkPO.FirstPaymentTransactionFee,
				FirstPaymentTax:            bulkPO.FirstPaymentTax,
				FirstPaymentSubTotal:       bulkPO.FirstPaymentSubTotal,
				FirstPaymentTotal:          bulkPO.FirstPaymentTotal,
				FirstPaymentPercentage:     values.Float64Value(bulkPO.FirstPaymentPercentage),
				FinalPaymentTransactionFee: bulkPO.FinalPaymentTransactionFee,
				FinalPaymentTax:            bulkPO.FinalPaymentTax,
				FinalPaymentSubTotal:       bulkPO.FinalPaymentSubTotal,
				FinalPaymentTotal:          bulkPO.FinalPaymentTotal,
			},
			Metadata: models.InvoiceMetadata{
				InvoiceType:                     enums.InvoiceTypeBulkPOFirstPayment,
				InquiryID:                       bulkPO.Inquiry.ID,
				InquiryReferenceID:              bulkPO.Inquiry.ReferenceID,
				BulkPurchaseOrderID:             bulkPO.ID,
				BulkPurchaseOrderReferenceID:    bulkPO.ReferenceID,
				BulkEstimatedProductionLeadTime: bulkPO.GetQuotationLeadTime(),
			},
		}

		if bulkPO.FirstPaymentType == enums.PaymentTypeBankTransfer && bulkPO.Currency == enums.VND {
			invoiceParams.Vendor = models.DefaultVendorForLocal
		}

		if bulkPO.FirstPaymentType == enums.PaymentTypeBankTransfer {
			invoiceParams.PaymentTransactionID = bulkPO.FirstPaymentTransactionRefID
		}

		if bulkPO.FirstPaymentReceivedAt != nil {
			invoiceParams.IssuedDate = *bulkPO.FirstPaymentReceivedAt
		}

		if bulkPO.FirstPaymentMarkAsPaidAt != nil {
			invoiceParams.IssuedDate = *bulkPO.FirstPaymentMarkAsPaidAt
		}

		if bulkPO.User != nil {
			invoiceParams.Consignee = &models.InvoiceParty{
				ID:          bulkPO.User.ID,
				Name:        bulkPO.User.Name,
				Email:       bulkPO.User.Email,
				PhoneNumber: bulkPO.User.PhoneNumber,
				CompanyName: bulkPO.User.CompanyName,
			}
			if bulkPO.User.Coordinate == nil && bulkPO.User.CoordinateID != "" {
				var coordinate models.Coordinate
				if err := r.db.First(&coordinate, "id = ?", bulkPO.User.CoordinateID).Error; err == nil {
					invoiceParams.Consignee.Address = coordinate.Display()
				}
			}
		}

		for _, item := range bulkPO.Items {
			invoiceParams.Items = append(invoiceParams.Items, &models.InvoiceItem{
				ItemCode:      item.ID,
				Color:         item.ColorName,
				SizeName:      item.Size,
				Description:   fmt.Sprintf("%s-%s", item.ColorName, item.Size),
				TotalQuantity: item.Qty,
				UnitPrice:     item.UnitPrice,
				TotalAmount:   item.TotalPrice,
			})
		}

		invoice, err = r.CreateInvoice(invoiceParams)
		if err != nil {
			return nil, eris.Wrapf(err, "Generate pdf bulk PO %s", bulkPO.ReferenceID)
		}
	}

	data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
		URL:               fmt.Sprintf("%s/invoices/print/%d", r.db.Configuration.AdminPortalBaseURL, invoice.InvoiceNumber),
		Selector:          "#invoice-ready-to-print",
		Landscape:         true,
		PrintBackground:   true,
		PreferCssPageSize: true,
	})
	if err != nil {
		return nil, eris.Wrapf(err, "Generate pdf bulk PO %s", bulkPO.ReferenceID)
	}

	if len(data) == 0 {
		return nil, eris.Errorf("Generate pdf bulk PO %s empty data", bulkPO.ReferenceID)
	}

	var uploadParams = s3.UploadFileParams{
		Key:         fmt.Sprintf("uploads/%s_%s_%d_first_payment_invoice.pdf", bulkPO.Inquiry.ReferenceID, bulkPO.ReferenceID, invoice.InvoiceNumber),
		Data:        bytes.NewBuffer(data),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(models.ContentTypePDF),
		ACL:         "private",
	}

	_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
	if err != nil {
		return nil, err
	}

	var attachment = models.Attachment{
		ContentType: uploadParams.ContentType,
		FileKey:     uploadParams.Key,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if bulkPO.FirstPaymentInvoice != nil && params.ReCreate {
			err = tx.Delete(&models.Invoice{}, "invoice_number = ?", bulkPO.FirstPaymentInvoice.InvoiceNumber).Error
			if err != nil {
				return err
			}
		}
		err = tx.Model(&models.Invoice{}).Where("invoice_number = ?", invoice.InvoiceNumber).UpdateColumn("Document", &attachment).Error
		if err != nil {
			return err
		}

		bulkPO.FirstPaymentInvoice = invoice
		bulkPO.FirstPaymentInvoice.Document = &attachment

		return tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).UpdateColumn("FirstPaymentInvoiceNumber", invoice.InvoiceNumber).Error
	})
	if err != nil {
		return nil, err
	}

	return bulkPO, err
}

type CreateBulkSecondPaymentInvoiceParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string          `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool            `json:"re_create"`
	Invoice             *models.Invoice `json:"-"`
}

func (r *InvoiceRepo) CreateBulkSecondPaymentInvoice(params CreateBulkSecondPaymentInvoiceParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
		IncludeItems:        true,
		IncludeInvoice:      true,
	})
	if err != nil {
		return nil, err
	}

	var invoice = params.Invoice

	if invoice == nil {
		if bulkPO.SecondPaymentInvoice != nil && !params.ReCreate {
			r.db.CustomLogger.Debugf("Bulk PO %s second payment invoice is already generated", bulkPO.ReferenceID)
			return bulkPO, errs.ErrBulkPoInvoiceAlreadyGenerated
		}

		var invoiceParams = models.CreateInvoiceParams{
			UserID:      bulkPO.UserID,
			InvoiceType: enums.InvoiceTypeBulkPOSecondPayment,
			Status:      enums.InvoiceStatusPaid,
			Currency:    string(bulkPO.Currency),
			Vendor:      models.DefaultVendorForOnlinePayment,
			// Note:                 fmt.Sprintf("Second payment, %.f", bulkPO.SecondPaymentPercentage) + "% of total payment",
			PaymentType:          bulkPO.SecondPaymentType,
			PaymentTransactionID: bulkPO.SecondPaymentIntentID,
			InvoicePricing: models.InvoicePricing{
				Pricing: models.Pricing{
					SubTotal:               bulkPO.SubTotal,
					TransactionFee:         bulkPO.TransactionFee,
					Tax:                    bulkPO.Tax,
					TotalPrice:             bulkPO.TotalPrice,
					TaxPercentage:          bulkPO.TaxPercentage,
					SubTotalAfterDeduction: bulkPO.SecondPaymentSubTotal,
				},
				SecondPaymentTransactionFee: bulkPO.SecondPaymentTransactionFee,
				SecondPaymentTax:            bulkPO.SecondPaymentTax,
				SecondPaymentSubTotal:       bulkPO.SecondPaymentSubTotal,
				SecondPaymentTotal:          bulkPO.SecondPaymentTotal,
				SecondPaymentPercentage:     bulkPO.SecondPaymentPercentage,
				FinalPaymentTransactionFee:  bulkPO.FinalPaymentTransactionFee,
				FinalPaymentTax:             bulkPO.FinalPaymentTax,
				FinalPaymentSubTotal:        bulkPO.FinalPaymentSubTotal,
				FinalPaymentTotal:           bulkPO.FinalPaymentTotal,
			},
			Metadata: models.InvoiceMetadata{
				InvoiceType:                     enums.InvoiceTypeBulkPOSecondPayment,
				InquiryID:                       bulkPO.Inquiry.ID,
				InquiryReferenceID:              bulkPO.Inquiry.ReferenceID,
				BulkPurchaseOrderID:             bulkPO.ID,
				BulkPurchaseOrderReferenceID:    bulkPO.ReferenceID,
				BulkEstimatedProductionLeadTime: bulkPO.GetQuotationLeadTime(),
			},
		}

		if bulkPO.SecondPaymentType == enums.PaymentTypeBankTransfer && bulkPO.Currency == enums.VND {
			invoiceParams.Vendor = models.DefaultVendorForLocal
		}

		if bulkPO.SecondPaymentType == enums.PaymentTypeBankTransfer {
			invoiceParams.PaymentTransactionID = bulkPO.SecondPaymentTransactionRefID
		}

		if bulkPO.SecondPaymentReceivedAt != nil {
			invoiceParams.IssuedDate = *bulkPO.SecondPaymentReceivedAt
		}

		if bulkPO.SecondPaymentMarkAsPaidAt != nil {
			invoiceParams.IssuedDate = *bulkPO.SecondPaymentMarkAsPaidAt
		}

		if bulkPO.User != nil {
			invoiceParams.Consignee = &models.InvoiceParty{
				ID:          bulkPO.User.ID,
				Name:        bulkPO.User.Name,
				Email:       bulkPO.User.Email,
				PhoneNumber: bulkPO.User.PhoneNumber,
				CompanyName: bulkPO.User.CompanyName,
			}
			if bulkPO.User.Coordinate == nil && bulkPO.User.CoordinateID != "" {
				var coordinate models.Coordinate
				if err := r.db.First(&coordinate, "id = ?", bulkPO.User.CoordinateID).Error; err == nil {
					invoiceParams.Consignee.Address = coordinate.Display()
				}
			}
		}

		for _, item := range bulkPO.Items {
			invoiceParams.Items = append(invoiceParams.Items, &models.InvoiceItem{
				ItemCode:      item.ID,
				Color:         item.ColorName,
				SizeName:      item.Size,
				Description:   fmt.Sprintf("%s-%s", item.ColorName, item.Size),
				TotalQuantity: item.Qty,
				UnitPrice:     item.UnitPrice,
				TotalAmount:   item.TotalPrice,
			})
		}

		invoice, err = r.CreateInvoice(invoiceParams)
		if err != nil {
			return nil, eris.Wrapf(err, "Generate pdf bulk PO %s", bulkPO.ReferenceID)
		}
	}

	data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
		URL:               fmt.Sprintf("%s/invoices/print/%d", r.db.Configuration.AdminPortalBaseURL, invoice.InvoiceNumber),
		Selector:          "#invoice-ready-to-print",
		Landscape:         true,
		PrintBackground:   true,
		PreferCssPageSize: true,
	})
	if err != nil {
		return nil, eris.Wrapf(err, "Generate pdf bulk PO %s", bulkPO.ReferenceID)
	}

	if len(data) == 0 {
		return nil, eris.Errorf("Generate pdf bulk PO %s empty data", bulkPO.ReferenceID)
	}

	var uploadParams = s3.UploadFileParams{
		Key:         fmt.Sprintf("uploads/%s_%s_%d_second_payment_invoice.pdf", bulkPO.Inquiry.ReferenceID, bulkPO.ReferenceID, invoice.InvoiceNumber),
		Data:        bytes.NewBuffer(data),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(models.ContentTypePDF),
		ACL:         "private",
	}

	_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
	if err != nil {
		return nil, err
	}

	var attachment = models.Attachment{
		ContentType: uploadParams.ContentType,
		FileKey:     uploadParams.Key,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if bulkPO.SecondPaymentInvoice != nil && params.ReCreate {
			err = tx.Delete(&models.Invoice{}, "invoice_number = ?", bulkPO.SecondPaymentInvoice.InvoiceNumber).Error
			if err != nil {
				return err
			}
		}
		err = tx.Model(&models.Invoice{}).Where("invoice_number = ?", invoice.InvoiceNumber).UpdateColumn("Document", &attachment).Error
		if err != nil {
			return err
		}
		bulkPO.SecondPaymentInvoice = invoice
		bulkPO.SecondPaymentInvoice.Document = &attachment
		return tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).UpdateColumn("SecondPaymentInvoiceNumber", invoice.InvoiceNumber).Error
	})
	if err != nil {
		return nil, err
	}

	return bulkPO, err
}

type CreatePurchaseOrderInvoiceParams struct {
	PurchaseOrderID string          `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	ReCreate        bool            `json:"re_create"`
	Invoice         *models.Invoice `json:"-"`
}

func (r *InvoiceRepo) CreatePurchaseOrderInvoice(params CreatePurchaseOrderInvoiceParams) (*models.PurchaseOrder, error) {
	purchaseOrder, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		IncludeInvoice:  true,
		IncludeUsers:    true,
	})
	if err != nil {
		return nil, err
	}

	var invoice = params.Invoice

	if invoice == nil {
		if purchaseOrder.Invoice != nil && !params.ReCreate {
			r.db.CustomLogger.Debugf("PO %s fpayment invoice is already generated", purchaseOrder.ReferenceID)
			return purchaseOrder, errs.ErrPOInvoiceAlreadyGenerated
		}

		var invoiceParams = models.CreateInvoiceParams{
			UserID:      purchaseOrder.UserID,
			InvoiceType: enums.InvoiceTypeInquiry,
			Status:      enums.InvoiceStatusPaid,
			Currency:    string(purchaseOrder.Currency),
			Vendor:      models.DefaultVendorForOnlinePayment,
			InvoicePricing: models.InvoicePricing{
				Pricing: purchaseOrder.Pricing,
			},
			IssuedDate:           time.Now().Unix(),
			Note:                 fmt.Sprintf("Charges for %s/%s", purchaseOrder.ReferenceID, purchaseOrder.Inquiry.ReferenceID),
			PaymentType:          purchaseOrder.PaymentType,
			PaymentTransactionID: purchaseOrder.PaymentIntentID,
			Metadata: models.InvoiceMetadata{
				InvoiceType:              enums.InvoiceTypeInquiry,
				InquiryID:                purchaseOrder.Inquiry.ID,
				InquiryReferenceID:       purchaseOrder.Inquiry.ReferenceID,
				PurchaseOrderID:          purchaseOrder.ID,
				PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
			},
		}

		if purchaseOrder.PaymentType == enums.PaymentTypeBankTransfer && purchaseOrder.Currency == enums.VND {
			invoiceParams.Vendor = models.DefaultVendorForLocal
		}

		if purchaseOrder.PaymentType == enums.PaymentTypeBankTransfer {
			invoiceParams.PaymentTransactionID = purchaseOrder.TransactionRefID
		}

		if purchaseOrder.TransferedAt != nil {
			invoiceParams.IssuedDate = *purchaseOrder.TransferedAt
		}

		if purchaseOrder.MarkAsPaidAt != nil {
			invoiceParams.IssuedDate = *purchaseOrder.MarkAsPaidAt
		}

		if purchaseOrder.User != nil {
			invoiceParams.Consignee = &models.InvoiceParty{
				ID:          purchaseOrder.User.ID,
				Name:        purchaseOrder.User.Name,
				Email:       purchaseOrder.User.Email,
				PhoneNumber: purchaseOrder.User.PhoneNumber,
				CompanyName: purchaseOrder.User.CompanyName,
			}

			if purchaseOrder.User.Coordinate == nil && purchaseOrder.User.CoordinateID != "" {
				var coordinate models.Coordinate
				if err := r.db.First(&coordinate, "id = ?", purchaseOrder.User.CoordinateID).Error; err == nil {
					invoiceParams.Consignee.Address = coordinate.Display()
				}
			}
		}

		for _, item := range purchaseOrder.CartItems {
			var parts = []string{
				purchaseOrder.Inquiry.SkuNote,
				item.ColorName,
				item.Size,
			}
			var nonEmptyParts = lo.Filter(parts, func(item string, index int) bool {
				return item != ""
			})

			var cartItem = &models.InvoiceItem{
				ItemCode:      item.ID,
				Color:         item.ColorName,
				SizeName:      item.Size,
				Description:   strings.Join(nonEmptyParts, "/"),
				TotalQuantity: item.Qty,
				UnitPrice:     item.UnitPrice.ToPtr(),
				TotalAmount:   item.TotalPrice.ToPtr(),
			}

			invoiceParams.Items = append(invoiceParams.Items, cartItem)
		}

		invoice, err = r.CreateInvoice(invoiceParams)
		if err != nil {
			return nil, eris.Wrapf(err, "Generate pdf purchase order %s failed", purchaseOrder.ReferenceID)
		}
	}

	data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
		URL:               fmt.Sprintf("%s/invoices/print/%d", r.db.Configuration.AdminPortalBaseURL, invoice.InvoiceNumber),
		Selector:          "#invoice-ready-to-print",
		Landscape:         true,
		PrintBackground:   true,
		PreferCssPageSize: true,
	})
	if err != nil {
		return nil, eris.Wrapf(err, "Generate pdf purchase order %s failed", purchaseOrder.ReferenceID)
	}

	if len(data) == 0 {
		return nil, eris.Errorf("Generate pdf purchase order %s failed empty data", purchaseOrder.ReferenceID)
	}
	var uploadParams = s3.UploadFileParams{
		Key:         fmt.Sprintf("uploads/%s_%s_%d_invoice.pdf", purchaseOrder.Inquiry.ReferenceID, purchaseOrder.ReferenceID, invoice.InvoiceNumber),
		Data:        bytes.NewBuffer(data),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(models.ContentTypePDF),
		ACL:         "private",
	}

	_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
	if err != nil {
		return nil, err
	}

	var updatedInvoice = models.Invoice{
		InvoiceNumber: invoice.InvoiceNumber,
		Document: &models.Attachment{
			ContentType: uploadParams.ContentType,
			FileKey:     uploadParams.Key,
		},
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if purchaseOrder.Invoice != nil && params.ReCreate {
			err = tx.Delete(&models.Invoice{}, "invoice_number = ?", purchaseOrder.Invoice.InvoiceNumber).Error
			if err != nil {
				return err
			}
		}
		err = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).UpdateColumn("InvoiceNumber", invoice.InvoiceNumber).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Invoice{}).Where("invoice_number = ?", invoice.InvoiceNumber).Updates(&updatedInvoice).Error
		if err != nil {
			return err
		}

		purchaseOrder.Invoice = invoice
		purchaseOrder.Invoice.Document = updatedInvoice.Document

		return err
	})
	if err != nil {
		return nil, err
	}

	return purchaseOrder, err
}

type CreateMultiplePurchaseOrderInvoiceParams struct {
	CheckoutSessionID string          `json:"checkout_session_id"`
	ReCreate          bool            `json:"re_create"`
	Invoice           *models.Invoice `json:"-"`
}

type CreateMultiplePurchaseOrderInvoiceResponse struct {
	PurchaseOrders     []*models.PurchaseOrder    `json:"purchase_orders"`
	PaymentTransaction *models.PaymentTransaction `json:"payment_transaction"`
}

func (r *InvoiceRepo) CreateMultiplePurchaseOrderInvoice(params CreateMultiplePurchaseOrderInvoiceParams) (*CreateMultiplePurchaseOrderInvoiceResponse, error) {
	var purchaseOrders []*models.PurchaseOrder
	var err = query.New(r.db, queryfunc.NewPurchaseOrderBuilder(queryfunc.PurchaseOrderBuilderOptions{
		IncludeCartItems: true,
		IncludeInvoice:   true,
		IncludeUsers:     true,
	})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("po.checkout_session_id = ?", params.CheckoutSessionID)
		}).
		FindFunc(&purchaseOrders)

	if len(purchaseOrders) == 0 {
		return nil, errs.ErrPONInvalid
	}
	var resp = CreateMultiplePurchaseOrderInvoiceResponse{
		PurchaseOrders: purchaseOrders,
	}

	var paymentTransaction models.PaymentTransaction
	err = query.New(r.db, queryfunc.NewPaymentTransactionBuilder(queryfunc.PaymentTransactionBuilderOptions{
		IncludeInvoice: true,
	})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.checkout_session_id = ?", params.CheckoutSessionID)
		}).
		FirstFunc(&paymentTransaction)
	if err != nil {
		return nil, err
	}

	var invoice = params.Invoice

	if invoice == nil {
		if paymentTransaction.Invoice != nil && !params.ReCreate {
			r.db.CustomLogger.Debugf("Payment transaction %s invoice is already generated", paymentTransaction.ReferenceID)
			return &resp, errs.ErrPaymentTransactionInvoiceAlreadyGenerated
		}

		var invoiceParams = models.CreateInvoiceParams{
			UserID:               paymentTransaction.UserID,
			InvoiceType:          enums.InvoiceTypeMultipleInquiry,
			Status:               enums.InvoiceStatusPaid,
			Currency:             string(paymentTransaction.Currency),
			Vendor:               models.DefaultVendorForOnlinePayment,
			IssuedDate:           paymentTransaction.CreatedAt,
			PaymentType:          paymentTransaction.PaymentType,
			PaymentTransactionID: paymentTransaction.PaymentIntentID,
			Metadata: models.InvoiceMetadata{
				InvoiceType:       enums.InvoiceTypeMultipleInquiry,
				CheckoutSessionID: paymentTransaction.CheckoutSessionID,
				InquiryIDs: lo.Map(purchaseOrders, func(item *models.PurchaseOrder, index int) string {
					return item.Inquiry.ID
				}),
				InquiryReferenceIDs: lo.Map(purchaseOrders, func(item *models.PurchaseOrder, index int) string {
					return item.Inquiry.ReferenceID
				}),
				PurchaseOrderIDs: lo.Map(purchaseOrders, func(item *models.PurchaseOrder, index int) string {
					return item.ID
				}),
				PurchaseOrderReferenceIDs: lo.Map(purchaseOrders, func(item *models.PurchaseOrder, index int) string {
					return item.ReferenceID
				}),
			},
		}

		invoiceParams.Note = fmt.Sprintf("Charges for %s", strings.Join(invoiceParams.Metadata.PurchaseOrderReferenceIDs, "/"))

		for _, purchaseOrder := range purchaseOrders {
			invoiceParams.InvoicePricing.Pricing = invoiceParams.InvoicePricing.Pricing.Add(purchaseOrder.Pricing)

			if invoiceParams.Consignee == nil {
				invoiceParams.Consignee = &models.InvoiceParty{
					ID:          purchaseOrder.User.ID,
					Name:        purchaseOrder.User.Name,
					Email:       purchaseOrder.User.Email,
					PhoneNumber: purchaseOrder.User.PhoneNumber,
					CompanyName: purchaseOrder.User.CompanyName,
				}

				if purchaseOrder.User.Coordinate == nil && purchaseOrder.User.CoordinateID != "" {
					var coordinate models.Coordinate
					if err := r.db.First(&coordinate, "id = ?", purchaseOrder.User.CoordinateID).Error; err == nil {
						invoiceParams.Consignee.Address = coordinate.Display()
					}
				}
			}

			for _, item := range purchaseOrder.CartItems {
				var parts = []string{
					purchaseOrder.ReferenceID,
					purchaseOrder.Inquiry.SkuNote,
					item.ColorName,
					item.Size,
				}
				var nonEmptyParts = lo.Filter(parts, func(item string, index int) bool {
					return item != ""
				})

				var cartItem = &models.InvoiceItem{
					ItemCode:      item.ID,
					Color:         item.ColorName,
					SizeName:      item.Size,
					Description:   strings.Join(nonEmptyParts, "/"),
					TotalQuantity: item.Qty,
					UnitPrice:     item.UnitPrice.ToPtr(),
					TotalAmount:   item.TotalPrice.ToPtr(),
				}

				invoiceParams.Items = append(invoiceParams.Items, cartItem)
			}

		}

		if paymentTransaction.PaymentType == enums.PaymentTypeBankTransfer && paymentTransaction.Currency == enums.VND {
			invoiceParams.Vendor = models.DefaultVendorForLocal
		}

		if paymentTransaction.PaymentType == enums.PaymentTypeBankTransfer {
			invoiceParams.PaymentTransactionID = paymentTransaction.TransactionRefID
		}

		invoice, err = r.CreateInvoice(invoiceParams)
		if err != nil {
			return nil, eris.Wrapf(err, "Generate pdf payment transaction %s failed", paymentTransaction.ReferenceID)
		}

	}

	data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
		URL:               fmt.Sprintf("%s/invoices/print/%d", r.db.Configuration.AdminPortalBaseURL, invoice.InvoiceNumber),
		Selector:          "#invoice-ready-to-print",
		Landscape:         true,
		PrintBackground:   true,
		PreferCssPageSize: true,
	})
	if err != nil {
		return nil, eris.Wrapf(err, "Generate pdf payment transaction %s, invoice %d failed", paymentTransaction.ReferenceID, invoice.InvoiceNumber)
	}

	if len(data) == 0 {
		return nil, eris.Errorf("Generate pdf purchase order %s failed empty data", paymentTransaction.ReferenceID)
	}
	var uploadParams = s3.UploadFileParams{
		Key:         fmt.Sprintf("uploads/multiple_rfqs_%s_%d_invoice.pdf", paymentTransaction.ReferenceID, invoice.InvoiceNumber),
		Data:        bytes.NewBuffer(data),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(models.ContentTypePDF),
		ACL:         "private",
	}

	_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
	if err != nil {
		return nil, err
	}

	var updatedInvoice = models.Invoice{
		InvoiceNumber: invoice.InvoiceNumber,
		Document: &models.Attachment{
			ContentType: uploadParams.ContentType,
			FileKey:     uploadParams.Key,
		},
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if paymentTransaction.Invoice != nil && params.ReCreate {
			err = tx.Delete(&models.Invoice{}, "invoice_number = ?", paymentTransaction.Invoice.InvoiceNumber).Error
			if err != nil {
				return err
			}
		}
		err = tx.Model(&models.PaymentTransaction{}).Where("id = ?", paymentTransaction.ID).UpdateColumn("InvoiceNumber", invoice.InvoiceNumber).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Invoice{}).Where("invoice_number = ?", invoice.InvoiceNumber).Updates(&updatedInvoice).Error
		if err != nil {
			return err
		}

		paymentTransaction.Invoice = invoice
		paymentTransaction.Invoice.Document = updatedInvoice.Document
		resp.PaymentTransaction = &paymentTransaction

		return err
	})
	if err != nil {
		return nil, err
	}

	return &resp, err
}

type GetInvoiceAttachmentParams struct {
	models.JwtClaimsInfo

	InvoiceNumber int  `json:"invoice_number" param:"invoice_number" validate:"required"`
	ReCreate      bool `json:"re_create" query:"re_create"`
}

func (r *InvoiceRepo) GetInvoiceAttachment(params GetInvoiceAttachmentParams) (*models.Attachment, error) {
	var invoice models.Invoice
	var err = r.db.First(&invoice, "invoice_number = ?", params.InvoiceNumber).Error
	if err != nil {
		return nil, err
	}

	if invoice.Document != nil && !params.ReCreate {
		return invoice.Document, nil
	}

	switch invoice.InvoiceType {
	case enums.InvoiceTypeBulkPODepositPayment:
		_, err = r.CreateBulkDepositInvoice(CreateBulkDepositInvoiceParams{
			BulkPurchaseOrderID: invoice.Metadata.BulkPurchaseOrderID,
			ReCreate:            true,
			Invoice:             &invoice,
		})

	case enums.InvoiceTypeBulkPOFirstPayment:
		_, err = r.CreateBulkFirstPaymentInvoice(CreateBulkFirstPaymentInvoiceParams{
			BulkPurchaseOrderID: invoice.Metadata.BulkPurchaseOrderID,
			ReCreate:            true,
			Invoice:             &invoice,
		})

	case enums.InvoiceTypeBulkPOSecondPayment:
		_, err = r.CreateBulkSecondPaymentInvoice(CreateBulkSecondPaymentInvoiceParams{
			BulkPurchaseOrderID: invoice.Metadata.BulkPurchaseOrderID,
			ReCreate:            true,
			Invoice:             &invoice,
		})

	case enums.InvoiceTypeBulkPOFinalPayment:
		_, err = r.CreateBulkFinalPaymentInvoice(CreateBulkFinalPaymentInvoiceParams{
			BulkPurchaseOrderID: invoice.Metadata.BulkPurchaseOrderID,
			ReCreate:            true,
			Invoice:             &invoice,
		})

	case enums.InvoiceTypeInquiry:
		_, err = r.CreatePurchaseOrderInvoice(CreatePurchaseOrderInvoiceParams{
			PurchaseOrderID: invoice.Metadata.PurchaseOrderID,
			ReCreate:        true,
			Invoice:         &invoice,
		})

	case enums.InvoiceTypeMultipleInquiry:
		_, err = r.CreateMultiplePurchaseOrderInvoice(CreateMultiplePurchaseOrderInvoiceParams{
			CheckoutSessionID: invoice.Metadata.CheckoutSessionID,
			ReCreate:          true,
			Invoice:           &invoice,
		})
	}

	return invoice.Document, err
}

type CreateBulkDebitNotesParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool   `json:"re_create"`
}

func (r *InvoiceRepo) CreateBulkDebitNotes(params CreateBulkDebitNotesParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
		IncludeInvoice:      true,
		IncludeItems:        true,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	var requireUpdate = false

	if bulkPO.DebitNoteAttachment == nil && bulkPO.CommercialInvoice != nil || (params.ReCreate && bulkPO.CommercialInvoice != nil) {
		data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
			URL:               fmt.Sprintf("%s/bulks/%s/commercial-invoice", r.db.Configuration.AdminPortalBaseURL, params.BulkPurchaseOrderID),
			Selector:          "#commercial-invoice-ready-to-print",
			Landscape:         true,
			PrintBackground:   true,
			PreferCssPageSize: true,
		})
		if err != nil {
			return nil, err
		}

		var uploadParams = s3.UploadFileParams{
			Key:         fmt.Sprintf("uploads/%s_%s_debit_note.pdf", bulkPO.Inquiry.ReferenceID, bulkPO.ReferenceID),
			Data:        bytes.NewBuffer(data),
			Bucket:      r.db.Configuration.AWSS3StorageBucket,
			ContentType: "application/pdf",
			ACL:         "private",
		}

		_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
		if err != nil {
			return nil, err
		}

		var attachment = models.Attachment{
			ContentType: uploadParams.ContentType,
			FileKey:     uploadParams.Key,
		}
		updates.DebitNoteAttachment = &attachment
		bulkPO.DebitNoteAttachment = &attachment
		requireUpdate = true
	}

	if requireUpdate {
		err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
		if err != nil {
			return nil, err
		}
	}

	return bulkPO, err
}

type CreateBulkCommercialInvoiceParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool   `json:"re_create"`
}

func (r *InvoiceRepo) CreateBulkCommercialInvoice(params CreateBulkCommercialInvoiceParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
		IncludeInvoice:      true,
		IncludeItems:        true,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	var requireUpdate = false

	if bulkPO.CommercialInvoiceAttachment == nil && bulkPO.CommercialInvoice != nil || (params.ReCreate && bulkPO.CommercialInvoice != nil) {
		data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
			URL:               fmt.Sprintf("%s/bulks/%s/commercial-invoice", r.db.Configuration.AdminPortalBaseURL, params.BulkPurchaseOrderID),
			Selector:          "#commercial-invoice-ready-to-print",
			Landscape:         true,
			PrintBackground:   true,
			PreferCssPageSize: true,
		})
		if err != nil {
			return nil, err
		}

		var uploadParams = s3.UploadFileParams{
			Key:         fmt.Sprintf("uploads/%s_%s_commercial_invoice.pdf", bulkPO.Inquiry.ReferenceID, bulkPO.ReferenceID),
			Data:        bytes.NewBuffer(data),
			Bucket:      r.db.Configuration.AWSS3StorageBucket,
			ContentType: "application/pdf",
			ACL:         "private",
		}

		_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
		if err != nil {
			return nil, err
		}

		var attachment = models.Attachment{
			ContentType: uploadParams.ContentType,
			FileKey:     uploadParams.Key,
		}

		bulkPO.CommercialInvoice.Status = enums.InvoiceStatusPaid
		bulkPO.CommercialInvoiceAttachment = &attachment

		updates.CommercialInvoiceAttachment = bulkPO.CommercialInvoiceAttachment
		updates.CommercialInvoice = bulkPO.CommercialInvoice
		requireUpdate = true
	}

	if requireUpdate {
		err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
		if err != nil {
			return nil, err
		}
	}

	return bulkPO, err
}

type CreateBulkPurchaseOrderInvoiceParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool   `json:"re_create"`
}

func (r *InvoiceRepo) CreateBulkPurchaseOrderInvoice(params CreateBulkPurchaseOrderInvoiceParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
		IncludeInvoice:      true,
		IncludeItems:        true,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	var requireUpdate = false

	if bulkPO.CommercialInvoiceAttachment == nil && bulkPO.CommercialInvoice != nil || (params.ReCreate && bulkPO.CommercialInvoice != nil) {
		data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
			URL:               fmt.Sprintf("%s/bulks/%s/commercial-invoice", r.db.Configuration.AdminPortalBaseURL, params.BulkPurchaseOrderID),
			Selector:          "#commercial-invoice-ready-to-print",
			Landscape:         true,
			PrintBackground:   true,
			PreferCssPageSize: true,
		})
		if err != nil {
			return nil, err
		}

		var uploadParams = s3.UploadFileParams{
			Key:         fmt.Sprintf("uploads/%s_%s_commercial_invoice.pdf", bulkPO.Inquiry.ReferenceID, bulkPO.ReferenceID),
			Data:        bytes.NewBuffer(data),
			Bucket:      r.db.Configuration.AWSS3StorageBucket,
			ContentType: "application/pdf",
			ACL:         "private",
		}

		_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
		if err != nil {
			return nil, err
		}

		var attachment = models.Attachment{
			ContentType: uploadParams.ContentType,
			FileKey:     uploadParams.Key,
		}
		updates.CommercialInvoiceAttachment = &attachment
		bulkPO.CommercialInvoiceAttachment = &attachment
		requireUpdate = true
	}

	if requireUpdate {
		err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
		if err != nil {
			return nil, err
		}
	}

	return bulkPO, err
}

type CreateOrderInvoiceRequest struct {
	models.JwtClaimsInfo
	PaymentTransactionID string `json:"payment_transaction_id" validate:"required"`
}

func (r *InvoiceRepo) CreatePaymentInvoice(req *CreateOrderInvoiceRequest) (*models.Invoice, error) {
	var payment models.PaymentTransaction
	if err := r.db.First(&payment, "id = ?", req.PaymentTransactionID).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrPaymentTransactionNotFound
		}
		return nil, err
	}
	if payment.Status != enums.PaymentStatusPaid {
		return nil, errs.ErrPaymentTransactionNotPaid
	}
	var existingInvoice models.Invoice
	if err := r.db.Select("ID").First(&existingInvoice, "payment_transaction_id = ?", req.PaymentTransactionID).Error; err != nil {
		if !r.db.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	var user models.User
	if err := r.db.First(&user, "id = ?", payment.UserID).Error; err != nil {
		return nil, err
	}
	var addressCoordinate models.Coordinate
	if user.CoordinateID != "" {
		if err := r.db.First(&addressCoordinate, "id = ?", user.CoordinateID).Error; err != nil {
			return nil, err
		}
	}

	var invoice = models.Invoice{
		ID:                   existingInvoice.ID,
		PaymentTransactionID: payment.ID,
		UserID:               payment.UserID,
		Status:               enums.InvoiceStatusPaid,
		Currency:             payment.Currency,
		Vendor:               models.DefaultVendorForOnlinePayment,
		Consignee: &models.InvoiceParty{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			CompanyName: user.CompanyName,
			Address:     addressCoordinate.Display(),
		},
		IssuedDate:  *payment.MarkAsPaidAt,
		PaymentType: payment.PaymentType,
		InvoicePricing: models.InvoicePricing{
			Pricing: models.Pricing{
				TotalPrice: payment.TotalAmount,
			},
		},
		Metadata: models.InvoiceMetadata{
			PurchaseOrderIDs:     payment.PurchaseOrderIDs,
			BulkPurchaseOrderIDs: payment.BulkPurchaseOrderIDs,
		},
	}
	if payment.PaymentType == enums.PaymentTypeBankTransfer && payment.Currency == enums.VND {
		invoice.Vendor = models.DefaultVendorForLocal
	}

	var orderCartRepo = NewOrderCartRepo(r.db)
	var subTotal = price.NewFromFloat(0)
	var tax = price.NewFromFloat(0)
	var shippingFee = price.NewFromFloat(0)
	var transactionFee = price.NewFromFloat(0)
	var orderReferenceIDs []string
	var orderCartItems models.OrderCartItems
	if len(payment.PurchaseOrderIDs) > 0 {
		var purchaseOrders, err = orderCartRepo.GetPurchaseOrderAndOrderItemsByPurchaseOrderIDs(payment.PurchaseOrderIDs, payment.UserID)
		if err != nil {
			return nil, err
		}
		for _, po := range purchaseOrders {
			orderReferenceIDs = append(orderReferenceIDs, po.ReferenceID)
			subTotal = subTotal.AddPtr(po.SubTotal)
			tax = tax.AddPtr(po.Tax)
			shippingFee = shippingFee.AddPtr(po.ShippingFee)
			transactionFee = transactionFee.AddPtr(po.TransactionFee)
			orderCartItems = append(orderCartItems, po.OrderCartItems...)
		}
	}
	if len(payment.BulkPurchaseOrderIDs) > 0 {
		var bulks, err = orderCartRepo.GetBulkAndOrderItemsByBulkIDs(payment.BulkPurchaseOrderIDs, payment.UserID)
		if err != nil {
			return nil, err
		}
		for _, bulk := range bulks {
			orderReferenceIDs = append(orderReferenceIDs, bulk.ReferenceID)

			if bulk.FirstPaymentTransactionReferenceID == payment.ReferenceID {
				subTotal = subTotal.AddPtr(bulk.FirstPaymentSubTotal)
				tax = tax.AddPtr(bulk.FirstPaymentTax)
				transactionFee = transactionFee.AddPtr(bulk.FirstPaymentTransactionFee)

			} else if bulk.FinalPaymentTransactionReferenceID == payment.ReferenceID {
				subTotal = subTotal.AddPtr(bulk.FinalPaymentSubTotal)
				tax = tax.AddPtr(bulk.FinalPaymentTax)
				transactionFee = transactionFee.AddPtr(bulk.FinalPaymentTransactionFee)
				shippingFee = shippingFee.AddPtr(bulk.ShippingFee)
			}
			orderCartItems = append(orderCartItems, bulk.OrderCartItems...)
		}
	}
	invoice.SubTotal = &subTotal
	invoice.Tax = &tax
	invoice.ShippingFee = &shippingFee
	invoice.TransactionFee = &transactionFee
	invoice.Note = fmt.Sprintf("Charges for %s", strings.Join(orderReferenceIDs, ", "))
	for _, item := range orderCartItems {
		var parts = []string{
			item.Sku,
			item.ColorName,
			item.Size,
		}
		var nonEmptyParts = lo.Filter(parts, func(item string, index int) bool {
			return item != ""
		})
		invoice.Items = append(invoice.Items, &models.InvoiceItem{
			ItemCode:      item.ID,
			Color:         item.ColorName,
			SizeName:      item.Size,
			Description:   strings.Join(nonEmptyParts, "/"),
			TotalQuantity: item.Qty,
			UnitPrice:     item.UnitPrice.ToPtr(),
			TotalAmount:   item.TotalPrice.ToPtr(),
		})
	}

	if err := r.db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, UpdateAll: true}).
		Create(&invoice).Error; err != nil {
		return nil, err
	}
	data, err := pdf.New(r.db.Configuration).GetPDF(pdf.GetPDFParams{
		URL:               fmt.Sprintf("%s/invoices/print/%d", r.db.Configuration.AdminPortalBaseURL, invoice.InvoiceNumber),
		Selector:          "#invoice-ready-to-print",
		Landscape:         true,
		PrintBackground:   true,
		PreferCssPageSize: true,
	})
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errs.ErrInvoiceGeneratePdfError
	}
	var uploadParams = s3.UploadFileParams{
		Key:         fmt.Sprintf("uploads/invoice_%s_%d_invoice.pdf", payment.ID, invoice.InvoiceNumber),
		Data:        bytes.NewBuffer(data),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(models.ContentTypePDF),
		ACL:         "private",
	}

	_, err = s3.New(r.db.Configuration).UploadFile(uploadParams)
	if err != nil {
		return nil, err
	}

	var document = &models.Attachment{
		ContentType: uploadParams.ContentType,
		FileKey:     uploadParams.Key,
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&models.PaymentTransaction{}).Where("id = ?", payment.ID).UpdateColumn("InvoiceNumber", invoice.InvoiceNumber).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Invoice{}).Where("id = ?", invoice.ID).UpdateColumn("Document", document).Error; err != nil {
			return err
		}

		invoice.Document = document
		return nil
	}); err != nil {
		return nil, err
	}

	return &invoice, nil

}
