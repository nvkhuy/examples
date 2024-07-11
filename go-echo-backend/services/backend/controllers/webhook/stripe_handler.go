package webhook

import (
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v74"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StripeHandler struct {
}

type HandlePaymentParams struct {
	PaymentIntent   *stripe.PaymentIntent
	CheckoutSession *stripe.CheckoutSession

	metadata      map[string]string
	latestCharge  *stripe.Charge
	paymentIntent *stripe.PaymentIntent
	actionSource  stripehelper.ActionSource
}

func (hanlder *StripeHandler) HandlePayment(c echo.Context, params *HandlePaymentParams) (err error) {
	var cc = c.(*models.CustomContext)

	if params.PaymentIntent != nil {
		params.paymentIntent = params.PaymentIntent
		params.metadata = params.PaymentIntent.Metadata
	}

	if params.CheckoutSession != nil {
		params.metadata = params.CheckoutSession.Metadata
		params.paymentIntent = params.CheckoutSession.PaymentIntent
		if params.paymentIntent == nil {
			params.paymentIntent, err = stripehelper.GetInstance().GetPaymentIntent(params.CheckoutSession.PaymentIntent.ID)
			if err != nil {
				return err
			}
		}

	}

	as, found := params.metadata["action_source"]
	if !found {
		cc.CustomLogger.Error("Action source in required")
		return nil
	}
	params.actionSource = stripehelper.ActionSource(as)

	if params.paymentIntent == nil {
		cc.CustomLogger.Errorf("PaymentIntent is empty")
		return nil
	}

	if params.paymentIntent.LatestCharge == nil {
		params.paymentIntent, err = stripehelper.GetInstance().GetPaymentIntent(params.paymentIntent.ID)
		if err != nil {
			cc.CustomLogger.Errorf("Get payment intent with latest charges for id=%s err=%+v", params.paymentIntent.ID, err)
			return nil
		}
	}

	if params.paymentIntent.LatestCharge == nil {
		cc.CustomLogger.Errorf("Lastest charges for %s not found", params.paymentIntent.ID)
		return nil

	}

	params.latestCharge, err = stripehelper.GetInstance().GetCharge(params.paymentIntent.LatestCharge.ID)
	if err != nil {
		return err
	}

	if params.latestCharge.ReceiptURL == "" {
		cc.CustomLogger.Errorf("Lastest charges does not have receipt %s", params.paymentIntent.ID)
		return err
	}

	if params.PaymentIntent == nil && params.paymentIntent != nil {
		params.PaymentIntent = params.paymentIntent
	}

	switch params.actionSource {
	case stripehelper.ActionSourceCreatePaymentLink,
		stripehelper.ActionSourceInquiryPayment:
		err = hanlder.HandleInquiryPayment(c, params)

	case stripehelper.ActionSourceBulkPODepositPayment:
		err = hanlder.HandleBulkPODepositPayment(c, params)

	case stripehelper.ActionSourceBulkPOFirstPayment:
		err = hanlder.HandleBulkPOFirstPayment(c, params)

	case stripehelper.ActionSourceBulkPOSecondPayment:
		err = hanlder.HandleBulkPOSecondPayment(c, params)

	case stripehelper.ActionSourceBulkPOFinalPayment:
		err = hanlder.HandleBulkPOFinalPayment(c, params)

	case stripehelper.ActionSourceMultiInquiryPayment:
		err = hanlder.HandleMultiInquiryPayment(c, params)

	case stripehelper.ActionSourceMultiPOPayment:
		err = hanlder.HandleMultiPOPayment(c, params)

	case stripehelper.ActionSourceOrderCartPayment:
		err = hanlder.HandleOrderPayment(c, params)

	default:
		cc.CustomLogger.Errorf("Action source %s is not supported", params.actionSource)
		return eris.Errorf("Action source %s is not supported", params.actionSource)
	}

	return err
}

func (hanlder *StripeHandler) HandleInquiryPayment(c echo.Context, params *HandlePaymentParams) error {
	var cc = c.(*models.CustomContext)

	purchaseOrderReferenceID, found := params.metadata["purchase_order_reference_id"]
	if !found || purchaseOrderReferenceID == "" {
		return eris.New("Purchase order not found")
	}

	var purchaseOrder models.PurchaseOrder
	var err = cc.App.DB.First(&purchaseOrder, "reference_id = ?", purchaseOrderReferenceID).Error
	if err != nil {
		return err
	}

	var inquiry models.Inquiry
	err = cc.App.DB.Select("ID", "ReferenceID").First(&inquiry, "id = ?", purchaseOrder.InquiryID).Error
	if err != nil {
		return err
	}

	cancel, err := cc.App.DB.Locker.AcquireLock(fmt.Sprintf("inquiry_%s", purchaseOrder.InquiryID), time.Second*30)
	if err != nil {
		return err
	}
	defer cancel()

	var updates = models.PurchaseOrder{
		Status:          enums.PurchaseOrderStatusPaid,
		PaymentType:     enums.PaymentTypeCard,
		TransferedAt:    values.Int64(time.Now().Unix()),
		PaymentIntentID: params.paymentIntent.ID,
		MarkAsPaidAt:    values.Int64(time.Now().Unix()),
		ReceiptURL:      params.latestCharge.ReceiptURL,
		ChargeID:        params.latestCharge.ID,
	}

	if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil && params.CheckoutSession.PaymentLink.ID != "" {
		updates.PaymentLinkID = params.CheckoutSession.PaymentLink.ID
	}

	if params.latestCharge.BalanceTransaction != nil {
		updates.TxnID = params.latestCharge.BalanceTransaction.ID
	}

	if len(purchaseOrder.Quotations) > 0 {
		sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeSample
		})
		if sampleQuotation != nil {
			updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
			updates.StartDate = updates.MarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	}

	err = cc.App.DB.Transaction(func(tx *gorm.DB) error {
		// create transaction
		var transaction = models.PaymentTransaction{
			PurchaseOrderID:   purchaseOrder.ID,
			PaidAmount:        purchaseOrder.TotalPrice,
			UserID:            purchaseOrder.UserID,
			TransactionRefID:  purchaseOrder.TransactionRefID,
			TotalAmount:       purchaseOrder.TotalPrice,
			Currency:          purchaseOrder.Currency,
			Status:            enums.PaymentStatusPaid,
			Milestone:         enums.PaymentMilestoneFinalPayment,
			PaymentIntentID:   updates.PaymentIntentID,
			ChargeID:          updates.ChargeID,
			TxnID:             updates.TxnID,
			PaymentPercentage: values.Float64(100),
			PaymentType:       enums.PaymentTypeCard,
			MarkAsPaidAt:      values.Int64(time.Now().Unix()),
			TransactionType:   enums.TransactionTypeCredit,
			PaymentLinkID:     updates.PaymentLinkID,
			Metadata: &models.PaymentTransactionMetadata{
				InquiryID:                inquiry.ID,
				InquiryReferenceID:       inquiry.ReferenceID,
				PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
				PurchaseOrderID:          purchaseOrder.ID,
			},
		}

		if params.latestCharge.BalanceTransaction != nil {
			if stripeCfg, err := stripehelper.GetCurrencyConfig(purchaseOrder.Currency); err == nil {
				transaction.Net = price.NewFromInt(params.latestCharge.BalanceTransaction.Net).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
				transaction.Fee = price.NewFromInt(params.latestCharge.BalanceTransaction.Fee).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
			}
		}

		var sqlResult = tx.Model(&models.PaymentTransaction{}).
			Where("purchase_order_id = ? AND milestone = ? AND transaction_type = ?", purchaseOrder.ID, transaction.Milestone, transaction.TransactionType).
			Updates(&transaction)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			err = tx.Create(&transaction).Error
		}

		updates.PaymentTransactionReferenceID = transaction.ReferenceID
		sqlResult = tx.Model(&models.PurchaseOrder{}).Where("reference_id = ?", purchaseOrderReferenceID).Updates(&updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return eris.New("Purchase order is not found")
		}

		var updateInquiry = models.Inquiry{
			Status:               enums.InquiryStatusFinished,
			BuyerQuotationStatus: enums.InquirySkuStatusApproved,
		}
		return tx.Model(&models.Inquiry{}).Where("id = ?", purchaseOrder.InquiryID).Updates(&updateInquiry).Error
	})
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.CreatePOPaymentInvoiceTask{
		PurchaseOrderID: purchaseOrder.ID,
	}.Dispatch(c.Request().Context())

	return err
}

func (hanlder *StripeHandler) HandleBulkPODepositPayment(c echo.Context, params *HandlePaymentParams) error {
	var cc = c.(*models.CustomContext)
	bulkPurchaseOrderReferenceID, found := params.metadata["bulk_purchase_order_reference_id"]
	if !found || bulkPurchaseOrderReferenceID == "" {
		return eris.New("Bulk PO Reference ID is empty")
	}

	var bulkPurchaseOrder models.BulkPurchaseOrder
	var err = cc.App.DB.First(&bulkPurchaseOrder, "reference_id = ?", bulkPurchaseOrderReferenceID).Error
	if err != nil {
		return err
	}

	var inquiry models.Inquiry
	err = cc.App.DB.Select("ID", "ReferenceID").First(&inquiry, "id = ?", bulkPurchaseOrder.InquiryID).Error
	if err != nil {
		return err
	}

	cancel, err := cc.App.DB.Locker.AcquireLock(fmt.Sprintf("bulk_purchase_order_%s", bulkPurchaseOrder.ID), time.Second*30)
	if err != nil {
		return err
	}
	defer cancel()

	var updates = models.BulkPurchaseOrder{
		DepositTransferedAt:    values.Int64(time.Now().Unix()),
		DepositMarkAsPaidAt:    values.Int64(time.Now().Unix()),
		DepositReceiptURL:      params.latestCharge.ReceiptURL,
		DepositChargeID:        params.latestCharge.ID,
		DepositPaymentIntentID: params.paymentIntent.ID,
		DepositPaidAmount:      bulkPurchaseOrder.DepositAmount,
	}
	if params.latestCharge.BalanceTransaction != nil {
		updates.DepositTxnID = params.latestCharge.BalanceTransaction.ID
	}

	if len(bulkPurchaseOrder.AdminQuotations) > 0 {
		var bulkQuotation = bulkPurchaseOrder.GetBulkQuotation()
		if bulkQuotation != nil {
			updates.LeadTime = int(values.Int64Value(bulkQuotation.LeadTime))
			updates.StartDate = updates.DepositMarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	}

	if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil && params.CheckoutSession.PaymentLink.ID != "" {
		updates.DepositPaymentLinkID = params.CheckoutSession.PaymentLink.ID
	}

	err = cc.App.DB.Transaction(func(tx *gorm.DB) error {
		var transaction = models.PaymentTransaction{
			BulkPurchaseOrderID: bulkPurchaseOrder.ID,
			Currency:            bulkPurchaseOrder.Currency,
			PaidAmount:          bulkPurchaseOrder.DepositAmount,
			PaymentType:         enums.PaymentTypeCard,
			Milestone:           enums.PaymentMilestoneDeposit,
			UserID:              bulkPurchaseOrder.UserID,
			Status:              enums.PaymentStatusPaid,
			BalanceAmount:       bulkPurchaseOrder.DepositAmount,
			PaymentPercentage:   values.Float64(100),
			TotalAmount:         bulkPurchaseOrder.DepositAmount,
			PaymentIntentID:     updates.DepositPaymentIntentID,
			ChargeID:            updates.DepositChargeID,
			TxnID:               updates.DepositTxnID,
			MarkAsPaidAt:        values.Int64(time.Now().Unix()),
			PaymentLinkID:       updates.DepositPaymentLinkID,
			TransactionType:     enums.TransactionTypeCredit,
			Metadata: &models.PaymentTransactionMetadata{
				InquiryID:                    inquiry.ID,
				InquiryReferenceID:           inquiry.ReferenceID,
				BulkPurchaseOrderReferenceID: bulkPurchaseOrder.ReferenceID,
				BulkPurchaseOrderID:          bulkPurchaseOrder.ID,
			},
		}

		if params.latestCharge.BalanceTransaction != nil {
			if stripeCfg, err := stripehelper.GetCurrencyConfig(bulkPurchaseOrder.Currency); err == nil {
				transaction.Net = price.NewFromInt(params.latestCharge.BalanceTransaction.Net).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
				transaction.Fee = price.NewFromInt(params.latestCharge.BalanceTransaction.Fee).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
			}
		}

		var sqlResult = tx.Model(&models.PaymentTransaction{}).
			Where("bulk_purchase_order_id = ? AND milestone = ? AND transaction_type = ?", bulkPurchaseOrder.ID, transaction.Milestone, transaction.TransactionType).
			Updates(&transaction)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			err = tx.Create(&transaction).Error
		}

		updates.DepositPaymentTransactionReferenceID = transaction.ReferenceID
		err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bulkPurchaseOrder.ID).Updates(&updates).Error
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	tasks.CreateBulkPoDepositPaymentInvoiceTask{
		BulkPurchaseOrderID: bulkPurchaseOrder.ID,
	}.Dispatch(c.Request().Context())

	return err
}

func (hanlder *StripeHandler) HandleBulkPOFirstPayment(c echo.Context, params *HandlePaymentParams) error {
	var cc = c.(*models.CustomContext)

	bulkPurchaseOrderReferenceID, found := params.metadata["bulk_purchase_order_reference_id"]
	if !found || bulkPurchaseOrderReferenceID == "" {
		return eris.New("Bulk PO Reference ID is empty")
	}

	var bulkPurchaseOrder models.BulkPurchaseOrder
	var err = cc.App.DB.First(&bulkPurchaseOrder, "reference_id = ?", bulkPurchaseOrderReferenceID).Error
	if err != nil {
		return err
	}

	var inquiry models.Inquiry
	err = cc.App.DB.Select("ID", "ReferenceID").First(&inquiry, "id = ?", bulkPurchaseOrder.InquiryID).Error
	if err != nil {
		return err
	}

	cancel, err := cc.App.DB.Locker.AcquireLock(fmt.Sprintf("bulk_purchase_order_%s", bulkPurchaseOrder.ID), time.Second*30)
	if err != nil {
		return err
	}
	defer cancel()

	var updates = models.BulkPurchaseOrder{
		TrackingStatus:           enums.BulkPoTrackingStatusFirstPaymentConfirmed,
		FirstPaymentTransferedAt: values.Int64(time.Now().Unix()),
		FirstPaymentMarkAsPaidAt: values.Int64(time.Now().Unix()),
		FirstPaymentReceiptURL:   params.latestCharge.ReceiptURL,
		FirstPaymentChargeID:     params.latestCharge.ID,
		FirstPaymentIntentID:     params.paymentIntent.ID,
	}
	if params.latestCharge.BalanceTransaction != nil {
		updates.FirstPaymentTxnID = params.latestCharge.BalanceTransaction.ID
	}

	if len(bulkPurchaseOrder.AdminQuotations) > 0 {
		sampleQuotation, _ := lo.Find(bulkPurchaseOrder.AdminQuotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeSample
		})
		if sampleQuotation != nil {
			updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
			updates.StartDate = updates.FirstPaymentMarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	}

	if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil && params.CheckoutSession.PaymentLink.ID != "" {
		updates.FirstPaymentLinkID = params.CheckoutSession.PaymentLink.ID
	}

	err = cc.App.DB.Transaction(func(tx *gorm.DB) error {
		var transaction = models.PaymentTransaction{
			BulkPurchaseOrderID: bulkPurchaseOrder.ID,
			Currency:            bulkPurchaseOrder.Currency,
			PaidAmount:          bulkPurchaseOrder.FirstPaymentTotal,
			PaymentType:         enums.PaymentTypeCard,
			Milestone:           enums.PaymentMilestoneFirstPayment,
			UserID:              bulkPurchaseOrder.UserID,
			Status:              enums.PaymentStatusPaid,
			BalanceAmount:       bulkPurchaseOrder.FinalPaymentTotal,
			PaymentPercentage:   bulkPurchaseOrder.FirstPaymentPercentage,
			TotalAmount:         bulkPurchaseOrder.TotalPrice,
			PaymentIntentID:     updates.FirstPaymentIntentID,
			ChargeID:            updates.FirstPaymentChargeID,
			TxnID:               updates.FirstPaymentTxnID,
			MarkAsPaidAt:        values.Int64(time.Now().Unix()),
			PaymentLinkID:       updates.FirstPaymentLinkID,
			TransactionType:     enums.TransactionTypeCredit,
			Metadata: &models.PaymentTransactionMetadata{
				InquiryID:                    inquiry.ID,
				InquiryReferenceID:           inquiry.ReferenceID,
				BulkPurchaseOrderReferenceID: bulkPurchaseOrder.ReferenceID,
				BulkPurchaseOrderID:          bulkPurchaseOrder.ID,
			},
		}

		if params.latestCharge.BalanceTransaction != nil {
			if stripeCfg, err := stripehelper.GetCurrencyConfig(bulkPurchaseOrder.Currency); err == nil {
				transaction.Net = price.NewFromInt(params.latestCharge.BalanceTransaction.Net).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
				transaction.Fee = price.NewFromInt(params.latestCharge.BalanceTransaction.Fee).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
			}
		}

		var sqlResult = tx.Model(&models.PaymentTransaction{}).
			Where("bulk_purchase_order_id = ? AND milestone = ? AND transaction_type = ?", bulkPurchaseOrder.ID, transaction.Milestone, transaction.TransactionType).
			Updates(&transaction)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			err = tx.Create(&transaction).Error
		}

		updates.FirstPaymentTransactionReferenceID = transaction.ReferenceID
		var err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bulkPurchaseOrder.ID).Updates(&updates).Error
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	tasks.CreateBulkPoFirstPaymentInvoiceTask{
		BulkPurchaseOrderID: bulkPurchaseOrder.ID,
	}.Dispatch(c.Request().Context())

	return err
}

func (hanlder *StripeHandler) HandleBulkPOSecondPayment(c echo.Context, params *HandlePaymentParams) error {
	var cc = c.(*models.CustomContext)

	bulkPurchaseOrderReferenceID, found := params.metadata["bulk_purchase_order_reference_id"]
	if !found || bulkPurchaseOrderReferenceID == "" {
		return eris.New("Bulk PO Reference ID is empty")
	}

	var bulkPurchaseOrder models.BulkPurchaseOrder
	var err = cc.App.DB.First(&bulkPurchaseOrder, "reference_id = ?", bulkPurchaseOrderReferenceID).Error
	if err != nil {
		return err
	}

	var inquiry models.Inquiry
	err = cc.App.DB.Select("ID", "ReferenceID").First(&inquiry, "id = ?", bulkPurchaseOrder.InquiryID).Error
	if err != nil {
		return err
	}

	cancel, err := cc.App.DB.Locker.AcquireLock(fmt.Sprintf("bulk_purchase_order_%s", bulkPurchaseOrder.ID), time.Second*30)
	if err != nil {
		return err
	}
	defer cancel()

	var updates = models.BulkPurchaseOrder{
		TrackingStatus:            enums.BulkPoTrackingStatusSecondPaymentConfirmed,
		SecondPaymentTransferedAt: values.Int64(time.Now().Unix()),
		SecondPaymentMarkAsPaidAt: values.Int64(time.Now().Unix()),
		SecondPaymentReceiptURL:   params.latestCharge.ReceiptURL,
		SecondPaymentChargeID:     params.latestCharge.ID,
		SecondPaymentIntentID:     params.paymentIntent.ID,
	}
	if params.latestCharge.BalanceTransaction != nil {
		updates.SecondPaymentTxnID = params.latestCharge.BalanceTransaction.ID
	}

	if len(bulkPurchaseOrder.AdminQuotations) > 0 {
		var bulkQuotation = bulkPurchaseOrder.GetBulkQuotation()
		if bulkQuotation != nil {
			updates.LeadTime = int(values.Int64Value(bulkQuotation.LeadTime))
			updates.StartDate = updates.SecondPaymentMarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	}

	if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil && params.CheckoutSession.PaymentLink.ID != "" {
		updates.SecondPaymentLinkID = params.CheckoutSession.PaymentLink.ID
	}

	err = cc.App.DB.Transaction(func(tx *gorm.DB) error {
		var transaction = models.PaymentTransaction{
			BulkPurchaseOrderID: bulkPurchaseOrder.ID,
			Currency:            bulkPurchaseOrder.Currency,
			PaidAmount:          bulkPurchaseOrder.SecondPaymentTotal,
			PaymentType:         enums.PaymentTypeCard,
			Milestone:           enums.PaymentMilestoneSecondPayment,
			UserID:              bulkPurchaseOrder.UserID,
			Status:              enums.PaymentStatusPaid,
			BalanceAmount:       bulkPurchaseOrder.FinalPaymentTotal,
			PaymentPercentage:   values.Float64(bulkPurchaseOrder.SecondPaymentPercentage),
			TotalAmount:         bulkPurchaseOrder.TotalPrice,
			PaymentIntentID:     updates.SecondPaymentIntentID,
			ChargeID:            updates.SecondPaymentChargeID,
			TxnID:               updates.SecondPaymentTxnID,
			MarkAsPaidAt:        values.Int64(time.Now().Unix()),
			PaymentLinkID:       updates.SecondPaymentLinkID,
			TransactionType:     enums.TransactionTypeCredit,
			Metadata: &models.PaymentTransactionMetadata{
				InquiryID:                    inquiry.ID,
				InquiryReferenceID:           inquiry.ReferenceID,
				BulkPurchaseOrderReferenceID: bulkPurchaseOrder.ReferenceID,
				BulkPurchaseOrderID:          bulkPurchaseOrder.ID,
			},
		}

		if params.latestCharge.BalanceTransaction != nil {
			if stripeCfg, err := stripehelper.GetCurrencyConfig(bulkPurchaseOrder.Currency); err == nil {
				transaction.Net = price.NewFromInt(params.latestCharge.BalanceTransaction.Net).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
				transaction.Fee = price.NewFromInt(params.latestCharge.BalanceTransaction.Fee).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
			}
		}

		var sqlResult = tx.Model(&models.PaymentTransaction{}).
			Where("bulk_purchase_order_id = ? AND milestone = ? AND transaction_type = ?", bulkPurchaseOrder.ID, transaction.Milestone, transaction.TransactionType).
			Updates(&transaction)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			err = tx.Create(&transaction).Error
		}

		updates.SecondPaymentTransactionReferenceID = transaction.ReferenceID
		var err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bulkPurchaseOrder.ID).Updates(&updates).Error
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	tasks.CreateBulkPoSecondPaymentInvoiceTask{
		BulkPurchaseOrderID: bulkPurchaseOrder.ID,
	}.Dispatch(c.Request().Context())

	return err
}

func (hanlder *StripeHandler) HandleBulkPOFinalPayment(c echo.Context, params *HandlePaymentParams) error {
	var cc = c.(*models.CustomContext)

	bulkPurchaseOrderReferenceID, found := params.metadata["bulk_purchase_order_reference_id"]
	if !found || bulkPurchaseOrderReferenceID == "" {
		return eris.New("Bulk PO Reference ID is empty")
	}

	var bulkPurchaseOrder models.BulkPurchaseOrder
	var err = cc.App.DB.First(&bulkPurchaseOrder, "reference_id = ?", bulkPurchaseOrderReferenceID).Error
	if err != nil {
		return err
	}

	var inquiry models.Inquiry
	err = cc.App.DB.Select("ID", "ReferenceID").First(&inquiry, "id = ?", bulkPurchaseOrder.InquiryID).Error
	if err != nil {
		return err
	}

	cancel, err := cc.App.DB.Locker.AcquireLock(fmt.Sprintf("bulk_purchase_order_%s", bulkPurchaseOrder.ID), time.Second*30)
	if err != nil {
		return err
	}
	defer cancel()

	var updates = models.BulkPurchaseOrder{
		TrackingStatus:           enums.BulkPoTrackingStatusFinalPaymentConfirmed,
		FinalPaymentTransferedAt: values.Int64(time.Now().Unix()),
		FinalPaymentMarkAsPaidAt: values.Int64(time.Now().Unix()),
		FinalPaymentReceiptURL:   params.latestCharge.ReceiptURL,
		FinalPaymentChargeID:     params.latestCharge.ID,
		FinalPaymentIntentID:     params.paymentIntent.ID,
	}
	if params.latestCharge.BalanceTransaction != nil {
		updates.FinalPaymentTxnID = params.latestCharge.BalanceTransaction.ID
	}

	if len(bulkPurchaseOrder.AdminQuotations) > 0 {
		sampleQuotation, _ := lo.Find(bulkPurchaseOrder.AdminQuotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeSample
		})
		if sampleQuotation != nil {
			updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
			updates.StartDate = updates.FinalPaymentMarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	}

	if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil && params.CheckoutSession.PaymentLink.ID != "" {
		updates.FinalPaymentLinkID = params.CheckoutSession.PaymentLink.ID
	}

	// create transaction
	var transaction = models.PaymentTransaction{
		BulkPurchaseOrderID: bulkPurchaseOrder.ID,
		Currency:            bulkPurchaseOrder.Currency,
		PaymentType:         enums.PaymentTypeCard,
		UserID:              bulkPurchaseOrder.UserID,
		Status:              enums.PaymentStatusPaid,
		PaymentPercentage:   values.Float64(100 - values.Float64Value(bulkPurchaseOrder.FirstPaymentPercentage) - bulkPurchaseOrder.SecondPaymentPercentage),
		TotalAmount:         bulkPurchaseOrder.TotalPrice,
		PaymentIntentID:     updates.FinalPaymentIntentID,
		ChargeID:            updates.FinalPaymentChargeID,
		TxnID:               updates.FinalPaymentTxnID,
		MarkAsPaidAt:        values.Int64(time.Now().Unix()),
		PaidAmount:          bulkPurchaseOrder.FinalPaymentTotal,
		BalanceAmount:       price.NewFromFloat(0).ToPtr(),
		Milestone:           enums.PaymentMilestoneFinalPayment,
		PaymentLinkID:       updates.FinalPaymentLinkID,
		TransactionType:     enums.TransactionTypeCredit,
		Metadata: &models.PaymentTransactionMetadata{
			InquiryID:                    inquiry.ID,
			InquiryReferenceID:           inquiry.ReferenceID,
			BulkPurchaseOrderReferenceID: bulkPurchaseOrder.ReferenceID,
			BulkPurchaseOrderID:          bulkPurchaseOrder.ID,
		},
	}

	if params.latestCharge.BalanceTransaction != nil {
		if stripeCfg, err := stripehelper.GetCurrencyConfig(bulkPurchaseOrder.Currency); err == nil {
			transaction.Net = price.NewFromInt(params.latestCharge.BalanceTransaction.Net).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
			transaction.Fee = price.NewFromInt(params.latestCharge.BalanceTransaction.Fee).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
		}
	}

	err = cc.App.DB.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Model(&models.PaymentTransaction{}).
			Where("bulk_purchase_order_id = ? AND milestone = ? AND transaction_type = ?", bulkPurchaseOrder.ID, transaction.Milestone, transaction.TransactionType).
			Updates(&transaction)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			err = tx.Create(&transaction).Error
		}

		updates.FinalPaymentTransactionReferenceID = transaction.ReferenceID
		sqlResult = tx.Model(&models.BulkPurchaseOrder{}).Where("reference_id = ?", bulkPurchaseOrderReferenceID).Updates(&updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			err = tx.Create(&transaction).Error
		}

		return err
	})
	if err != nil {
		return err
	}

	tasks.CreateBulkPoFinalPaymentInvoiceTask{
		BulkPurchaseOrderID: bulkPurchaseOrder.ID,
		ReCreate:            true,
	}.Dispatch(c.Request().Context())

	return err
}

func (hanlder *StripeHandler) HandleMultiInquiryPayment(c echo.Context, params *HandlePaymentParams) error {
	var cc = c.(*models.CustomContext)

	checkoutSessionID, found := params.metadata["checkout_session_id"]
	if !found || checkoutSessionID == "" {
		return eris.New("Checkout session is not found")
	}

	var orders []*models.PurchaseOrder
	var err = query.New(cc.App.DB, queryfunc.NewPurchaseOrderBuilder(queryfunc.PurchaseOrderBuilderOptions{})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("po.checkout_session_id = ?", checkoutSessionID)
		}).
		FindFunc(&orders)
	if err != nil {
		return err
	}

	err = cc.App.DB.Transaction(func(tx *gorm.DB) error {

		var userID string
		var totalAmount price.Price
		var currency enums.Currency
		var inquiryIDs []string
		var inquiryReferenceIDs []string
		var purchaseOrderIDs []string
		var purchaseOrderReferenceIDs []string
		var paymentTransactionReferenceID = helper.GeneratePaymentTransactionReferenceID()

		for _, purchaseOrder := range orders {
			userID = purchaseOrder.UserID
			totalAmount = totalAmount.AddPtr(purchaseOrder.TotalPrice)
			currency = purchaseOrder.Currency

			purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
			purchaseOrderReferenceIDs = append(purchaseOrderReferenceIDs, purchaseOrder.ReferenceID)

			if purchaseOrder.Inquiry != nil {
				inquiryIDs = append(inquiryIDs, purchaseOrder.Inquiry.ID)
				inquiryReferenceIDs = append(inquiryReferenceIDs, purchaseOrder.Inquiry.ReferenceID)
			}

			cancel, err := cc.App.DB.Locker.AcquireLock(fmt.Sprintf("purchase_order_%s", purchaseOrder.ID), time.Second*30)
			if err != nil {
				return err
			}
			defer cancel()

			var updates = models.PurchaseOrder{
				Status:          enums.PurchaseOrderStatusPaid,
				PaymentType:     enums.PaymentTypeCard,
				TransferedAt:    values.Int64(time.Now().Unix()),
				PaymentIntentID: params.paymentIntent.ID,
				MarkAsPaidAt:    values.Int64(time.Now().Unix()),
				ReceiptURL:      params.latestCharge.ReceiptURL,
				ChargeID:        params.latestCharge.ID,
			}
			if params.latestCharge.BalanceTransaction != nil {
				updates.TxnID = params.latestCharge.BalanceTransaction.ID
			}

			updates.PaymentTransactionReferenceID = paymentTransactionReferenceID

			if len(purchaseOrder.Quotations) > 0 {
				sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
					return item.Type == enums.InquiryTypeSample
				})
				if sampleQuotation != nil {
					updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
					updates.StartDate = updates.MarkAsPaidAt
					updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
				}
			}
			var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
			if sqlResult.Error != nil {
				return sqlResult.Error
			}

			if sqlResult.RowsAffected == 0 {
				return eris.New("Purchase order is not found")
			}

			var tracking = models.PurchaseOrderTracking{
				PurchaseOrderID: purchaseOrder.ID,
				FromStatus:      purchaseOrder.TrackingStatus,
				ToStatus:        updates.TrackingStatus,
				UserID:          purchaseOrder.UserID,
				UserGroup:       enums.PoTrackingUserGroupBuyer,
				ActionType:      enums.PoTrackingActionPaymentReceived,
				Description:     fmt.Sprintf("Payment received for checkout session %s", checkoutSessionID),
				Metadata: &models.PoTrackingMetadata{
					After: map[string]interface{}{
						"checkout_session_id": checkoutSessionID,
					},
				},
			}
			tx.Create(&tracking)

			if purchaseOrder.Inquiry != nil {
				var inquiryAudit = models.InquiryAudit{
					InquiryID:       purchaseOrder.InquiryID,
					PurchaseOrderID: purchaseOrder.ID,
					UserID:          purchaseOrder.UserID,
					ActionType:      enums.AuditActionTypeInquiryAdminRefund,
					Description:     fmt.Sprintf("Payment received for checkout session %s", checkoutSessionID),
					Metadata: &models.InquiryAuditMetadata{
						After: map[string]interface{}{
							"checkout_session_id": checkoutSessionID,
						},
					},
				}
				tx.Create(&inquiryAudit)

			}

		}

		// create transaction
		var transaction = models.PaymentTransaction{
			ReferenceID:       paymentTransactionReferenceID,
			PaidAmount:        totalAmount.ToPtr(),
			UserID:            userID,
			TotalAmount:       totalAmount.ToPtr(),
			Currency:          currency,
			Status:            enums.PaymentStatusPaid,
			Milestone:         enums.PaymentMilestoneFinalPayment,
			PaymentIntentID:   params.PaymentIntent.ID,
			ChargeID:          params.latestCharge.ID,
			PaymentPercentage: values.Float64(100),
			PaymentType:       enums.PaymentTypeCard,
			MarkAsPaidAt:      values.Int64(time.Now().Unix()),
			Metadata: &models.PaymentTransactionMetadata{
				InquiryIDs:                inquiryIDs,
				InquiryReferenceIDs:       inquiryReferenceIDs,
				PurchaseOrderReferenceIDs: purchaseOrderReferenceIDs,
				PurchaseOrderIDs:          purchaseOrderIDs,
			},
			CheckoutSessionID: checkoutSessionID,
		}
		if params.latestCharge.BalanceTransaction != nil {
			if stripeCfg, err := stripehelper.GetCurrencyConfig(currency); err == nil {
				transaction.Net = price.NewFromInt(params.latestCharge.BalanceTransaction.Net).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
				transaction.Fee = price.NewFromInt(params.latestCharge.BalanceTransaction.Fee).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
			}
		}

		if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil && params.CheckoutSession.PaymentLink.ID != "" {
			transaction.PaymentLinkID = params.CheckoutSession.PaymentLink.ID
		}

		if params.latestCharge.BalanceTransaction != nil {
			transaction.TxnID = params.latestCharge.BalanceTransaction.ID
		}

		var sqlResult = tx.Model(&models.PaymentTransaction{}).
			Where("checkout_session_id = ? AND transaction_type = ?", checkoutSessionID, enums.TransactionTypeCredit).
			Updates(&transaction)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			err = tx.Create(&transaction).Error
		}

		if len(inquiryIDs) > 0 {
			var updateInquiry = models.Inquiry{
				Status:               enums.InquiryStatusFinished,
				BuyerQuotationStatus: enums.InquirySkuStatusApproved,
			}

			err = tx.Model(&models.Inquiry{}).Where("id IN ?", inquiryIDs).Updates(&updateInquiry).Error
		}
		return err
	})

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.CreatePOPaymentInvoiceForMultipleItemsTask{
		CheckoutSessionID: checkoutSessionID,
	}.Dispatch(c.Request().Context())

	return err
}

func (hanlder *StripeHandler) HandleMultiPOPayment(c echo.Context, params *HandlePaymentParams) error {
	var cc = c.(*models.CustomContext)

	checkoutSessionID, found := params.metadata["checkout_session_id"]
	if !found || checkoutSessionID == "" {
		return eris.New("Checkout session is not found")
	}

	var orders []*models.PurchaseOrder
	var err = query.New(cc.App.DB, queryfunc.NewPurchaseOrderBuilder(queryfunc.PurchaseOrderBuilderOptions{})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("po.checkout_session_id = ?", checkoutSessionID)
		}).
		FindFunc(&orders)
	if err != nil {
		return err
	}

	err = cc.App.DB.Transaction(func(tx *gorm.DB) error {

		var userID string
		var totalAmount price.Price
		var currency enums.Currency
		var inquiryIDs []string
		var inquiryReferenceIDs []string
		var purchaseOrderIDs []string
		var purchaseOrderReferenceIDs []string
		var paymentTransactionReferenceID = helper.GeneratePaymentTransactionReferenceID()

		for _, purchaseOrder := range orders {
			userID = purchaseOrder.UserID
			totalAmount = totalAmount.AddPtr(purchaseOrder.TotalPrice)
			currency = purchaseOrder.Currency

			purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
			purchaseOrderReferenceIDs = append(purchaseOrderReferenceIDs, purchaseOrder.ReferenceID)

			if purchaseOrder.Inquiry != nil {
				inquiryIDs = append(inquiryIDs, purchaseOrder.Inquiry.ID)
				inquiryReferenceIDs = append(inquiryReferenceIDs, purchaseOrder.Inquiry.ReferenceID)
			}

			cancel, err := cc.App.DB.Locker.AcquireLock(fmt.Sprintf("purchase_order_%s", purchaseOrder.ID), time.Second*30)
			if err != nil {
				return err
			}
			defer cancel()

			var updates = models.PurchaseOrder{
				Status:          enums.PurchaseOrderStatusPaid,
				PaymentType:     enums.PaymentTypeCard,
				TransferedAt:    values.Int64(time.Now().Unix()),
				PaymentIntentID: params.paymentIntent.ID,
				MarkAsPaidAt:    values.Int64(time.Now().Unix()),
				ReceiptURL:      params.latestCharge.ReceiptURL,
				ChargeID:        params.latestCharge.ID,
			}
			if params.latestCharge.BalanceTransaction != nil {
				updates.TxnID = params.latestCharge.BalanceTransaction.ID
			}

			if len(purchaseOrder.Quotations) > 0 {
				sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
					return item.Type == enums.InquiryTypeSample
				})
				if sampleQuotation != nil {
					updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
					updates.StartDate = updates.MarkAsPaidAt
					updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
				}
			}

			updates.PaymentTransactionReferenceID = paymentTransactionReferenceID
			var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
			if sqlResult.Error != nil {
				return sqlResult.Error
			}

			if sqlResult.RowsAffected == 0 {
				return eris.New("Purchase order is not found")
			}

			var tracking = models.PurchaseOrderTracking{
				PurchaseOrderID: purchaseOrder.ID,
				FromStatus:      purchaseOrder.TrackingStatus,
				ToStatus:        updates.TrackingStatus,
				UserID:          purchaseOrder.UserID,
				UserGroup:       enums.PoTrackingUserGroupBuyer,
				ActionType:      enums.PoTrackingActionPaymentReceived,
				Description:     fmt.Sprintf("Payment received for checkout session %s", checkoutSessionID),
				Metadata: &models.PoTrackingMetadata{
					After: map[string]interface{}{
						"checkout_session_id": checkoutSessionID,
					},
				},
			}
			tx.Create(&tracking)

			if purchaseOrder.Inquiry != nil {
				var inquiryAudit = models.InquiryAudit{
					InquiryID:       purchaseOrder.InquiryID,
					PurchaseOrderID: purchaseOrder.ID,
					UserID:          purchaseOrder.UserID,
					ActionType:      enums.AuditActionTypeInquiryAdminRefund,
					Description:     fmt.Sprintf("Payment received for checkout session %s", checkoutSessionID),
					Metadata: &models.InquiryAuditMetadata{
						After: map[string]interface{}{
							"checkout_session_id": checkoutSessionID,
						},
					},
				}
				tx.Create(&inquiryAudit)

			}

		}

		// create transaction
		var transaction = models.PaymentTransaction{
			ReferenceID:       paymentTransactionReferenceID,
			PaidAmount:        totalAmount.ToPtr(),
			UserID:            userID,
			TotalAmount:       totalAmount.ToPtr(),
			Currency:          currency,
			Status:            enums.PaymentStatusPaid,
			Milestone:         enums.PaymentMilestoneFinalPayment,
			PaymentIntentID:   params.paymentIntent.ID,
			ChargeID:          params.latestCharge.ID,
			PaymentPercentage: values.Float64(100),
			PaymentType:       enums.PaymentTypeCard,
			MarkAsPaidAt:      values.Int64(time.Now().Unix()),
			Metadata: &models.PaymentTransactionMetadata{
				InquiryIDs:                inquiryIDs,
				InquiryReferenceIDs:       inquiryReferenceIDs,
				PurchaseOrderReferenceIDs: purchaseOrderReferenceIDs,
				PurchaseOrderIDs:          purchaseOrderIDs,
			},
			CheckoutSessionID: checkoutSessionID,
		}
		if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil && params.CheckoutSession.PaymentLink.ID != "" {
			transaction.PaymentLinkID = params.CheckoutSession.PaymentLink.ID
		}

		if params.latestCharge.BalanceTransaction != nil {
			transaction.TxnID = params.latestCharge.BalanceTransaction.ID
		}

		if params.latestCharge.BalanceTransaction != nil {
			if stripeCfg, err := stripehelper.GetCurrencyConfig(currency); err == nil {
				transaction.Net = price.NewFromInt(params.latestCharge.BalanceTransaction.Net).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
				transaction.Fee = price.NewFromInt(params.latestCharge.BalanceTransaction.Fee).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
			}
		}

		var sqlResult = tx.Model(&models.PaymentTransaction{}).
			Where("checkout_session_id = ? AND transaction_type = ?", checkoutSessionID, enums.TransactionTypeCredit).
			Updates(&transaction)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			err = tx.Create(&transaction).Error
		}

		if len(inquiryIDs) > 0 {
			var updateInquiry = models.Inquiry{
				Status:               enums.InquiryStatusFinished,
				BuyerQuotationStatus: enums.InquirySkuStatusApproved,
			}

			err = tx.Model(&models.Inquiry{}).Where("id IN ?", inquiryIDs).Updates(&updateInquiry).Error
		}

		return err

	})

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.CreatePOPaymentInvoiceForMultipleItemsTask{
		CheckoutSessionID: checkoutSessionID,
	}.Dispatch(c.Request().Context())

	return err
}

func (hanlder *StripeHandler) HandleOrderPayment(c echo.Context, params *HandlePaymentParams) error {
	var cc = c.(*models.CustomContext)

	checkoutSessionID, found := params.metadata["checkout_session_id"]
	if !found || checkoutSessionID == "" {
		return eris.New("Checkout session is not found")
	}
	poCartItemIDsString, found := params.metadata["purchase_order_cart_item_ids"]
	if !found {
		poCartItemIDsString = ""
	}
	bulkIDsString, found := params.metadata["bulk_order_ids"]
	if !found {
		bulkIDsString = ""
	}
	userID, found := params.metadata["user_id"]
	if !found || userID == "" {
		return eris.New("user_id is not found")
	}

	if poCartItemIDsString == "" && bulkIDsString == "" {
		return eris.New("Order is not empty")
	}
	var poCartItemIDs, bulkIDs []string
	if poCartItemIDsString != "" {
		poCartItemIDs = strings.Split(poCartItemIDsString, ",")
	}
	if bulkIDsString != "" {
		bulkIDs = strings.Split(bulkIDsString, ",")
	}

	orderCartRepo := repo.NewOrderCartRepo(cc.App.DB)

	purchaseOrders, err := orderCartRepo.GetPurchaseOrderAndOrderItemsByOrderCartItemIDs(poCartItemIDs, userID)
	if err != nil {
		return err
	}
	bulks, err := orderCartRepo.GetBulkAndOrderItemsByBulkIDs(bulkIDs, userID)
	if err != nil {
		return err
	}

	purchaseOrdersToUpdate := make([]*models.PurchaseOrder, 0, len(purchaseOrders))
	inquiryIDsToUpdate := make([]string, 0, len(purchaseOrders))
	bulksToUpdate := make([]*models.BulkPurchaseOrder, 0, len(bulks))
	var orderCartItemIDsToUpdate []string
	var paymentTransaction models.PaymentTransaction
	var paymentTransactionReferenceID = helper.GeneratePaymentTransactionReferenceID()
	var currency = enums.Currency(strings.ToUpper(string(params.paymentIntent.Currency)))
	stripeCfg, err := stripehelper.GetCurrencyConfig(currency)
	if err != nil {
		return err
	}

	for _, po := range purchaseOrders {
		po.Status = enums.PurchaseOrderStatusPaid
		po.ReceiptURL = params.latestCharge.ReceiptURL
		po.ChargeID = params.latestCharge.ID
		po.TransferedAt = values.Int64(time.Now().Unix())
		po.MarkAsPaidAt = values.Int64(time.Now().Unix())
		po.PaymentIntentID = params.paymentIntent.ID
		po.PaymentTransactionReferenceID = paymentTransactionReferenceID
		po.CheckoutSessionID = checkoutSessionID
		if params.latestCharge.BalanceTransaction != nil {
			po.TxnID = params.latestCharge.BalanceTransaction.ID
		}
		if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil {
			po.PaymentLinkID = params.CheckoutSession.PaymentLink.ID
			po.PaymentLink = params.CheckoutSession.PaymentLink.URL
		}
		purchaseOrdersToUpdate = append(purchaseOrdersToUpdate, po)
		if po.InquiryID != "" {
			inquiryIDsToUpdate = append(inquiryIDsToUpdate, po.InquiryID)
		}
		orderCartItemIDsToUpdate = append(orderCartItemIDsToUpdate, models.OrderCartItems(po.OrderCartItems).IDs()...)
	}
	for _, bpo := range bulks {
		if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment || bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPaymentConfirmed {
			bpo.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirmed
			bpo.FirstPaymentTransferedAt = values.Int64(time.Now().Unix())
			bpo.FirstPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
			bpo.FirstPaymentTransactionReferenceID = paymentTransactionReferenceID
			bpo.FirstPaymentIntentID = params.paymentIntent.ID
			bpo.FirstPaymentReceiptURL = params.latestCharge.ReceiptURL
			bpo.FirstPaymentChargeID = params.latestCharge.ID
			bpo.FirstPaymentCheckoutSessionID = checkoutSessionID
			if params.latestCharge.BalanceTransaction != nil {
				bpo.FirstPaymentTxnID = params.latestCharge.BalanceTransaction.ID
			}
			if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil {
				bpo.FirstPaymentLinkID = params.CheckoutSession.PaymentLink.ID
				bpo.FirstPaymentLink = params.CheckoutSession.PaymentLink.URL
			}
		} else {
			bpo.TrackingStatus = enums.BulkPoTrackingStatusFinalPaymentConfirmed
			bpo.FinalPaymentTransferedAt = values.Int64(time.Now().Unix())
			bpo.FinalPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
			bpo.FinalPaymentTransactionReferenceID = paymentTransactionReferenceID
			bpo.FinalPaymentIntentID = params.paymentIntent.ID
			bpo.FinalPaymentReceiptURL = params.latestCharge.ReceiptURL
			bpo.FinalPaymentChargeID = params.latestCharge.ID
			bpo.FinalPaymentCheckoutSessionID = checkoutSessionID
			if params.latestCharge.BalanceTransaction != nil {
				bpo.FinalPaymentTxnID = params.latestCharge.BalanceTransaction.ID
			}
			if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil {
				bpo.FinalPaymentLinkID = params.CheckoutSession.PaymentLink.ID
				bpo.FinalPaymentLink = params.CheckoutSession.PaymentLink.URL

			}
			orderCartItemIDsToUpdate = append(orderCartItemIDsToUpdate, models.OrderCartItems(bpo.OrderCartItems).IDs()...)
		}
		bulksToUpdate = append(bulksToUpdate, bpo)
	}
	paymentTransaction = models.PaymentTransaction{
		PaymentIntentID:      params.paymentIntent.ID,
		ReferenceID:          paymentTransactionReferenceID,
		PaidAmount:           price.NewFromInt(params.paymentIntent.Amount).DivInt(stripeCfg.SmallestUnitFactor).ToPtr(),
		PaymentType:          enums.PaymentTypeCard,
		UserID:               userID,
		Status:               enums.PaymentStatusPaid,
		TotalAmount:          price.NewFromInt(params.paymentIntent.Amount).DivInt(stripeCfg.SmallestUnitFactor).ToPtr(),
		Currency:             currency,
		ChargeID:             params.latestCharge.ID,
		ReceiptURL:           params.latestCharge.ReceiptURL,
		CheckoutSessionID:    checkoutSessionID,
		PurchaseOrderIDs:     purchaseOrders.IDs(),
		BulkPurchaseOrderIDs: bulks.IDs(),
		TransactionType:      enums.TransactionTypeCredit,
		MarkAsPaidAt:         values.Int64(time.Now().Unix()),
		Metadata: &models.PaymentTransactionMetadata{
			PurchaseOrderIDs:     purchaseOrders.IDs(),
			BulkPurchaseOrderIDs: bulks.IDs(),
		},
	}
	if params.latestCharge.BalanceTransaction != nil {
		paymentTransaction.Net = price.NewFromInt(params.latestCharge.BalanceTransaction.Net).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()
		paymentTransaction.Fee = price.NewFromInt(params.latestCharge.BalanceTransaction.Fee).DivInt(stripeCfg.SmallestUnitFactor).ToPtr()

	}

	if params.CheckoutSession != nil && params.CheckoutSession.PaymentLink != nil && params.CheckoutSession.PaymentLink.ID != "" {
		paymentTransaction.PaymentLinkID = params.CheckoutSession.PaymentLink.ID
	}

	if params.latestCharge.BalanceTransaction != nil {
		paymentTransaction.TxnID = params.latestCharge.BalanceTransaction.ID
	}

	if err := cc.App.DB.Transaction(func(tx *gorm.DB) error {
		paymentUpdateResult := tx.Model(&paymentTransaction).
			Clauses(clause.Returning{}).
			Where("checkout_session_id = ? AND transaction_type = ?", checkoutSessionID, enums.TransactionTypeCredit).
			Updates(&paymentTransaction)
		if paymentUpdateResult.Error != nil {
			return paymentUpdateResult.Error
		}
		if paymentUpdateResult.RowsAffected == 0 {
			if err := tx.Create(&paymentTransaction).Error; err != nil {
				return err
			}
		}

		if len(purchaseOrdersToUpdate) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&purchaseOrdersToUpdate).Error; err != nil {
				return err
			}
		}
		if len(bulksToUpdate) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&bulksToUpdate).Error; err != nil {
				return err
			}
		}
		if len(inquiryIDsToUpdate) > 0 {
			if err := tx.Model(&models.Inquiry{}).Where("id IN ?", inquiryIDsToUpdate).
				UpdateColumn("Status", enums.InquiryStatusFinished).Error; err != nil {
				return err
			}
		}
		if len(orderCartItemIDsToUpdate) > 0 {
			if err := tx.Model(&models.OrderCartItem{}).Where("id IN ?", orderCartItemIDsToUpdate).
				Updates(&models.OrderCartItem{CheckoutSessionID: checkoutSessionID, WaitingForCheckout: values.Bool(false)}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.CreatePaymentInvoiceTask{
		PaymentTransactionID: paymentTransaction.ID,
	}.Dispatch(c.Request().Context())

	return nil
}
