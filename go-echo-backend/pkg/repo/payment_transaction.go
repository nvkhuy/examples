package repo

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentTransactionRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewPaymentTransactionRepo(db *db.DB) *PaymentTransactionRepo {
	return &PaymentTransactionRepo{
		db:     db,
		logger: logger.New("repo/PaymentTransaction"),
	}
}

type PaginatePaymentTransactionsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	Statuses     []enums.PaymentStatus    `json:"statuses" query:"statuses" form:"statuses"`
	UserIDs      []string                 `json:"user_ids" query:"user_ids" form:"user_ids"`
	UserID       string                   `json:"user_id" query:"user_id" form:"user_id"`
	Milestones   []enums.PaymentMilestone `json:"milestones" query:"milestones" form:"milestones"`
	PaidDateFrom int64                    `json:"paid_date_from" query:"paid_date_from" form:"paid_date_from"`
	PaidDateTo   int64                    `json:"paid_date_to" query:"paid_date_to" form:"paid_date_to"`

	IncludeDetails bool `json:"-"`
	IncludeInvoice bool `json:"-"`
}

func (r *PaymentTransactionRepo) PaginatePaymentTransactions(params PaginatePaymentTransactionsParams) *query.Pagination {
	var builder = queryfunc.NewPaymentTransactionBuilder(queryfunc.PaymentTransactionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeDetails: params.IncludeDetails,
		IncludeInvoice: params.IncludeInvoice,
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.GetRole().IsAdmin() {
				if params.UserID != "" {
					builder.Where("p.user_id = ?", params.UserID)
				}
			} else {
				builder.Where("p.user_id = ?", params.GetUserID())
			}

			if len(params.Statuses) > 0 {
				builder.Where("p.status IN ?", params.Statuses)
			}

			if len(params.UserIDs) > 0 {
				builder.Where("p.user_id IN ?", params.UserIDs)
			}

			if len(params.Milestones) > 0 {
				builder.Where("p.milestone IN ?", params.Milestones)
			}

			if params.PaidDateFrom > 0 {
				builder.Where("p.mark_as_paid_at >= ?", params.PaidDateFrom)
			}

			if params.PaidDateTo > 0 {
				builder.Where("p.mark_as_paid_at <= ?", params.PaidDateTo)
			}

			if params.GetRole().IsAdmin() {
				if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
					if strings.HasPrefix(keyword, "pi_") {
						builder.Where("p.payment_intent_id = ?", keyword)
					} else if strings.HasPrefix(keyword, "PAY-") {
						builder.Where("p.reference_id = ?", keyword)
					} else if strings.HasPrefix(keyword, "PO-") {
						builder.Where("p.metadata->>'purchase_order_reference_id' = ?", keyword)
					} else if strings.HasPrefix(keyword, "BPO-") {
						builder.Where("p.metadata->>'bulk_purchase_order_reference_id' = ?", keyword)
					} else if strings.HasPrefix(keyword, "IQ-") {
						builder.Where("p.metadata->>'inquiry_reference_id' = ?", keyword)
					} else {
						builder.Where("(u.name ILIKE @keyword OR u.email ILIKE @keyword)", sql.Named("keyword", "%"+keyword+"%"))
					}
				}
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		WithoutCount(params.WithoutCount).
		PagingFunc()

	return result
}

type GetPaymentTransactionsParams struct {
	models.JwtClaimsInfo
	PaymentTransactionID string `json:"payment_transaction_id" param:"payment_transaction_id" validate:"required"`
	Note                 string `json:"note"`

	IncludeInvoice bool `json:"-"`
	IncludeDetails bool `json:"-"`
}

func (r *PaymentTransactionRepo) GetPaymentTransaction(params GetPaymentTransactionsParams) (result models.PaymentTransaction, err error) {
	var builder = queryfunc.NewPaymentTransactionBuilder(queryfunc.PaymentTransactionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeDetails: params.IncludeDetails,
		IncludeInvoice: params.IncludeInvoice,
	})
	result = models.PaymentTransaction{}
	err = query.New(r.db, builder).
		Limit(1).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.id = ?", params.PaymentTransactionID)
		}).
		FirstFunc(&result)
	return
}

func (r *PaymentTransactionRepo) ApprovePaymentTransactions(params GetPaymentTransactionsParams) (*models.PaymentTransaction, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("transaction_%s", params.PaymentTransactionID), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var transaction models.PaymentTransaction
	err = r.db.First(&transaction, "id = ?", params.PaymentTransactionID).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrPaymentTransactionNotFound
		}
		return nil, err
	}

	if transaction.Status != enums.PaymentStatusWaitingConfirm {
		return nil, errs.ErrPaymentTransactionIsInvalid
	}

	var purchaseOrders = make(models.PurchaseOrders, 0, len(transaction.PurchaseOrderIDs))
	if len(transaction.PurchaseOrderIDs) > 0 {
		if err := r.db.Find(&purchaseOrders, "id IN ?", []string(transaction.PurchaseOrderIDs)).Error; err != nil {
			return nil, err
		}
		var inquiries = make(models.Inquiries, 0, len(purchaseOrders.InquiryIDs()))
		if err := r.db.Find(&inquiries, "id IN ?", purchaseOrders.InquiryIDs()).Error; err != nil {
			return nil, err
		}
		var mapInquiryIDToInquiry = make(map[string]*models.Inquiry, len(inquiries))
		for _, iq := range inquiries {
			mapInquiryIDToInquiry[iq.ID] = iq
		}

		for _, po := range purchaseOrders {
			var quotations = po.Quotations
			if len(quotations) == 0 {
				var inquiry, ok = mapInquiryIDToInquiry[po.InquiryID]
				if ok {
					quotations = inquiry.AdminQuotations
				}
			}
			sampleQuotation, found := lo.Find(quotations, func(item *models.InquiryQuotationItem) bool {
				return item.Type == enums.InquiryTypeSample
			})
			if found {
				po.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
				po.StartDate = values.Int64(time.Now().Unix())
				po.CompletionDate = values.Int64(time.Unix(*po.StartDate, 0).AddDate(0, 0, po.LeadTime).Unix())
			}

			po.Status = enums.PurchaseOrderStatusPaid
			po.MarkAsPaidAt = values.Int64(time.Now().Unix())

		}

	}
	var bulks = make(models.BulkPurchaseOrders, 0, len(transaction.BulkPurchaseOrderIDs))
	if len(transaction.BulkPurchaseOrderIDs) > 0 {
		if err := r.db.Find(&bulks, "id IN ?", []string(transaction.BulkPurchaseOrderIDs)).Error; err != nil {
			return nil, err
		}
		for _, bpo := range bulks {
			if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPaymentConfirm {
				var quotations = bpo.AdminQuotations
				bulkQuotation, found := lo.Find(quotations, func(item *models.InquiryQuotationItem) bool {
					return item.Type == enums.InquiryTypeBulk
				})
				if found {
					bpo.LeadTime = int(values.Int64Value(bulkQuotation.LeadTime))
					bpo.StartDate = values.Int64(time.Now().Unix())
					bpo.CompletionDate = values.Int64(time.Unix(*bpo.StartDate, 0).AddDate(0, 0, bpo.LeadTime).Unix())
				}
				bpo.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirmed
				bpo.FirstPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
			}
			if bpo.TrackingStatus == enums.BulkPoTrackingStatusFinalPaymentConfirm {
				bpo.TrackingStatus = enums.BulkPoTrackingStatusFinalPaymentConfirmed
				bpo.FinalPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
			}

		}
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Returning{}).Model(&transaction).Where("id = ?", transaction.ID).
			Updates(&models.PaymentTransaction{Status: enums.PaymentStatusPaid, MarkAsPaidAt: values.Int64(time.Now().Unix())}).Error; err != nil {
			return err
		}

		if len(purchaseOrders) > 0 {
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&purchaseOrders).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Inquiry{}).Where("id IN ?", purchaseOrders.InquiryIDs()).
				UpdateColumn("Status", enums.InquiryStatusFinished).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.OrderCartItem{}).Where("purchase_order_id IN ?", purchaseOrders.IDs()).
				Updates(&models.OrderCartItem{CheckoutSessionID: transaction.CheckoutSessionID, WaitingForCheckout: values.Bool(false)}).Error; err != nil {
				return err
			}
		}
		if len(bulks) > 0 {
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&bulks).Error; err != nil {
				return err
			}
			var finalPaymentBulkIDs []string
			for _, bpo := range bulks {
				if bpo.TrackingStatus == enums.BulkPoTrackingStatusFinalPaymentConfirmed {
					finalPaymentBulkIDs = append(finalPaymentBulkIDs, bpo.ID)
				}
			}
			if err := tx.Model(&models.OrderCartItem{}).Where("bulk_purchase_order_id IN ?", finalPaymentBulkIDs).
				Updates(&models.OrderCartItem{CheckoutSessionID: transaction.CheckoutSessionID, WaitingForCheckout: values.Bool(false)}).Error; err != nil {
				return err
			}

		}
		return nil
	}); err != nil {
		return nil, err
	}
	transaction.PurchaseOrders = purchaseOrders
	transaction.BulkPurchaseOrders = bulks

	return &transaction, nil

}

func (r *PaymentTransactionRepo) RejectPaymentTransactions(params GetPaymentTransactionsParams) (*models.PaymentTransaction, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("transaction_%s", params.PaymentTransactionID), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var transaction models.PaymentTransaction
	err = r.db.First(&transaction, "id = ?", params.PaymentTransactionID).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrPaymentTransactionNotFound
		}
		return nil, err
	}

	if transaction.Status != enums.PaymentStatusWaitingConfirm {
		return nil, errs.ErrPaymentTransactionIsInvalid
	}

	var purchaseOrders = make(models.PurchaseOrders, 0, len(transaction.PurchaseOrderIDs))
	if len(transaction.PurchaseOrderIDs) > 0 {
		if err := r.db.Find(&purchaseOrders, "id IN ?", []string(transaction.PurchaseOrderIDs)).Error; err != nil {
			return nil, err
		}
		for _, po := range purchaseOrders {
			po.Status = enums.PurchaseOrderStatusUnpaid
			po.MarkAsUnpaidAt = values.Int64(time.Now().Unix())
		}
	}

	var bulks = make(models.BulkPurchaseOrders, 0, len(transaction.BulkPurchaseOrderIDs))
	if len(transaction.BulkPurchaseOrderIDs) > 0 {
		if err := r.db.Find(&bulks, "id IN ?", []string(transaction.BulkPurchaseOrderIDs)).Error; err != nil {
			return nil, err
		}
		for _, bpo := range bulks {
			if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPaymentConfirm {
				bpo.FirstPaymentMarkAsUnpaidAt = values.Int64(time.Now().Unix())
			}
			if bpo.TrackingStatus == enums.BulkPoTrackingStatusFinalPaymentConfirm {
				bpo.FinalPaymentMarkAsUnpaidAt = values.Int64(time.Now().Unix())
			}
		}
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Returning{}).Model(&transaction).Where("id = ?", transaction.ID).
			Updates(&models.PaymentTransaction{Status: enums.PaymentStatusUnpaid, MarkAsUnpaidAt: values.Int64(time.Now().Unix())}).Error; err != nil {
			return err
		}
		if len(purchaseOrders) > 0 {
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&purchaseOrders).Error; err != nil {
				return err
			}
		}
		if len(bulks) > 0 {
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&bulks).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &transaction, nil
}

type GetPaymentTransactionAttachmentsParams struct {
	models.JwtClaimsInfo

	PaymentTransactionID string `json:"payment_transaction_id" param:"payment_transaction_id" query:"payment_transaction_id" validate:"required"`
}

type GetPaymentTransactionAttachmentsResult struct {
	DownloadURLs []string `json:"download_urls"`
}

func (r *PaymentTransactionRepo) GetPaymentTransactionAttachments(params GetPaymentTransactionAttachmentsParams) (*GetPaymentTransactionAttachmentsResult, error) {
	trans, err := r.GetPaymentTransaction(GetPaymentTransactionsParams{
		JwtClaimsInfo:        params.JwtClaimsInfo,
		PaymentTransactionID: params.PaymentTransactionID,
	})
	if err != nil {
		return nil, err
	}

	var invoice models.Invoice
	err = r.db.Select("Document").First(&invoice, "payment_transaction_reference_id = ?", trans.ReferenceID).Error
	if err != nil {
		return nil, err
	}

	if trans.PaymentType == enums.PaymentTypeCard {
		if trans.PaymentIntentID == "" {
			return nil, eris.Errorf("Payment transaction %s has invalid payment intent", trans.ReferenceID)
		}

		pi, err := stripehelper.GetInstance().GetPaymentIntent(trans.PaymentIntentID)
		if err != nil {
			return nil, err
		}

		var resp = GetPaymentTransactionAttachmentsResult{
			DownloadURLs: []string{
				pi.LatestCharge.ReceiptURL,
				NewCommonRepo(r.db).GetDownloadLink(GetDownloadLinkParams{
					FileKey: invoice.Document.FileKey,
				}),
			},
		}

		return &resp, nil
	}

	var resp = GetPaymentTransactionAttachmentsResult{
		DownloadURLs: []string{
			NewCommonRepo(r.db).GetDownloadLink(GetDownloadLinkParams{
				FileKey: invoice.Document.FileKey,
			}),
		},
	}

	return &resp, nil
}

func (r *PaymentTransactionRepo) ExportExcel(params PaginatePaymentTransactionsParams) (*models.Attachment, error) {
	params.WithoutCount = true
	params.IncludeDetails = true
	params.Limit = 100
	params.Page = 1
	var results []*models.PaymentTransaction
Loop:
	var result = r.PaginatePaymentTransactions(params)
	if result == nil || result.Records == nil {
		return nil, errors.New("empty response")
	}

	trans, ok := result.Records.([]*models.PaymentTransaction)
	if !ok {
		return nil, errors.New("empty response")
	}

	results = append(results, trans...)

	if len(trans) == params.Limit {
		params.Page++
		goto Loop
	}

	fileContent, err := models.PaymentTransactions(results).ToExcel()
	if err != nil {
		return nil, err
	}

	var contentType = models.ContentTypeXLSX
	url := fmt.Sprintf("uploads/inquiries/export/export_payment_user_%s%s", params.GetUserID(), contentType.GetExtension())
	_, err = s3.New(r.db.Configuration).UploadFile(s3.UploadFileParams{
		Data:         bytes.NewReader(fileContent),
		Bucket:       r.db.Configuration.AWSS3StorageBucket,
		ContentType:  string(contentType),
		ACL:          "private",
		Key:          url,
		CacheControl: values.String("public, must-revalidate, proxy-revalidate, max-age=0"),
	})
	if err != nil {
		return nil, err
	}

	var resp = models.Attachment{
		FileKey:     url,
		ContentType: string(contentType),
	}
	return &resp, err
}
