package repo

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/samber/lo"

	"github.com/rotisserie/eris"
	"github.com/stripe/stripe-go/v74"
	stripePrice "github.com/stripe/stripe-go/v74/price"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderCartRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewOrderCartRepo(db *db.DB) *OrderCartRepo {
	return &OrderCartRepo{
		db:     db,
		logger: logger.New("repo/order_cart"),
	}
}

func (repo *OrderCartRepo) GetOrderCart(req *models.GetOrderCartRequest) (*models.GetOrderCartResponse, error) {
	userID := req.GetUserID()
	useRole := req.GetRole()

	var purchaseOrders []*models.PurchaseOrder
	var poCartBuilder = queryfunc.NewOrderCartPurchaseOrderBuilder(queryfunc.OrderCartPurchaseOrderBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: useRole,
		},
		IncludeInquiry: true,
	})
	if err := query.New(repo.db, poCartBuilder).WhereFunc(func(builder *query.Builder) {
		builder.Where("po.status = ?", enums.PurchaseOrderStatusPending)
		builder.Where("po.user_id = ?", userID)
		builder.Where("EXISTS (SELECT 1 FROM order_cart_items oci WHERE oci.purchase_order_id = po.id AND oci.unit_price != 0)")
	}).FindFunc(&purchaseOrders); err != nil {
		return nil, err
	}

	var bulkPurchaseOrders []*models.BulkPurchaseOrder
	var bpoCartBuilder = queryfunc.NewOrderCartBulkBuilder(queryfunc.OrderCartBulkBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: useRole,
		},
		IncludeInquiry: true,
	})
	if err := query.New(repo.db, bpoCartBuilder).WhereFunc(func(builder *query.Builder) {
		builder.Where("bpo.tracking_status IN ?", []enums.BulkPoTrackingStatus{enums.BulkPoTrackingStatusFirstPayment, enums.BulkPoTrackingStatusFinalPayment})
		builder.Where("bpo.user_id = ?", userID)
		builder.Where("EXISTS (SELECT 1 FROM order_cart_items oci WHERE oci.bulk_purchase_order_id = bpo.id AND oci.unit_price != 0)")
	}).FindFunc(&bulkPurchaseOrders); err != nil {
		return nil, err
	}

	return &models.GetOrderCartResponse{
		PurchaseOrders:     purchaseOrders,
		BulkPurchaseOrders: bulkPurchaseOrders,
	}, nil
}

func (repo *OrderCartRepo) OrderCartPreviewCheckout(req *models.OrderCartPreviewCheckoutRequest) (*models.GetOrderCartResponse, error) {
	if len(req.PurchaseOrderCartItemIDs) == 0 && len(req.BulkOrderIDs) == 0 {
		return nil, errs.ErrOrderEmpty
	}
	purchaseOrders, err := repo.GetPurchaseOrderAndOrderItemsByOrderCartItemIDs(req.PurchaseOrderCartItemIDs, req.GetUserID())
	if err != nil {
		return nil, err
	}
	purchaseOrders, err = repo.ProcessPurchaseOrdersPreviewCheckout(purchaseOrders, req.GetUserID(), req.PaymentType)
	if err != nil {
		return nil, err
	}
	bulkOrders, err := repo.GetBulkAndOrderItemsByBulkIDs(req.BulkOrderIDs, req.GetUserID())
	if err != nil {
		return nil, err
	}
	bulkOrders, err = repo.ProcessBulkPurchaseOrdersPreviewCheckout(bulkOrders, req.GetUserID(), req.PaymentType)
	if err != nil {
		return nil, err
	}
	if len(purchaseOrders) > 0 && len(bulkOrders) > 0 {
		if purchaseOrders[0].Currency != bulkOrders[0].Currency {
			return nil, errs.ErrOrderCurrencyMismatch
		}
	}
	if err := repo.db.Transaction(func(tx *gorm.DB) error {
		if len(purchaseOrders) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&purchaseOrders).Error; err != nil {
				return err
			}
		}
		if len(bulkOrders) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&bulkOrders).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &models.GetOrderCartResponse{
		PurchaseOrders:     purchaseOrders,
		BulkPurchaseOrders: bulkOrders,
	}, nil
}

func (repo *OrderCartRepo) ProcessPurchaseOrdersPreviewCheckout(purchaseOrders models.PurchaseOrders, userID string, paymentType enums.PaymentType) (models.PurchaseOrders, error) {
	if len(purchaseOrders) == 0 {
		return purchaseOrders, nil
	}
	currency := purchaseOrders[0].Currency
	for _, po := range purchaseOrders {
		if po.Status == enums.PurchaseOrderStatusPaid {
			return nil, eris.Wrapf(errs.ErrPoIsAlreadyPaid, "purchase_order_id:%s", po.ID)
		}
		if currency != po.Currency {
			return nil, errs.ErrOrderCurrencyMismatch
		}
		if err := ValidateOrderCartItems(po.OrderCartItems); err != nil {
			return nil, err
		}
	}

	inquiryIDs := purchaseOrders.InquiryIDs()
	inquiries := make(models.Inquiries, 0, len(inquiryIDs))
	if err := repo.db.Find(&inquiries, "id IN ? AND user_id = ?", inquiryIDs, userID).Error; err != nil {
		return nil, err
	}
	dbInquiryIDs := inquiries.IDs()
	for _, id := range inquiryIDs {
		if !helper.StringContains(dbInquiryIDs, id) {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "inquiry_id:%s", id)
		}
	}

	mapInquiryIdToInquiry := make(map[string]*models.Inquiry, len(inquiries))
	for _, iq := range inquiries {
		mapInquiryIdToInquiry[iq.ID] = iq
	}

	addressIDs := purchaseOrders.AddressIDs()
	addresses := make([]*models.Address, 0, len(addressIDs))
	if err := query.New(repo.db, queryfunc.NewAddressBuilder(queryfunc.AddressBuilderOptions{})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("a.id IN ?", addressIDs)
		}).
		FindFunc(&addresses); err != nil {
		return nil, err
	}
	mapAddressIdToAddress := make(map[string]*models.Address, len(addresses))
	for _, address := range addresses {
		mapAddressIdToAddress[address.ID] = address
	}

	for _, po := range purchaseOrders {
		inquiry, ok := mapInquiryIdToInquiry[po.InquiryID]
		if !ok && po.InquiryID != "" {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "purchase_order_id:%s", po.ID)
		}
		address, ok := mapAddressIdToAddress[po.ShippingAddressID]
		if !ok && po.ShippingAddressID != "" {
			return nil, eris.Wrapf(errs.ErrAddressNotFound, "purchase_order_id:%s", po.ID)
		}

		var subTotalPrice = price.NewFromFloat(0)
		for _, item := range po.OrderCartItems {
			subTotalPrice = subTotalPrice.Add(item.TotalPrice)
		}

		po.PaymentType = paymentType
		po.SubTotal = subTotalPrice.ToPtr()
		if inquiry != nil {
			po.Currency = inquiry.Currency
			po.TaxPercentage = inquiry.TaxPercentage
			po.ShippingFee = inquiry.ShippingFee
			po.ProductWeight = inquiry.ProductWeight
			po.ShippingAddressID = inquiry.ShippingAddressID
			po.Quotations = inquiry.AdminQuotations

			po.Inquiry = inquiry
		}
		po.ShippingAddress = address

		if err := po.UpdatePrices(); err != nil {
			return nil, err
		}
	}
	return purchaseOrders, nil
}

func (repo *OrderCartRepo) ProcessBulkPurchaseOrdersPreviewCheckout(bulks models.BulkPurchaseOrders, userID string, paymentType enums.PaymentType) ([]*models.BulkPurchaseOrder, error) {
	if len(bulks) == 0 {
		return bulks, nil
	}

	currency := bulks[0].Currency
	for _, bpo := range bulks {
		if bpo.TrackingStatus != enums.BulkPoTrackingStatusFirstPayment && bpo.TrackingStatus != enums.BulkPoTrackingStatusFinalPayment {
			return nil, eris.Wrapf(errs.ErrBulkPoInvalidToCheckout, "bulk_purchase_order_id:%s", bpo.ID)
		}
		if bpo.Currency != currency {
			return nil, errs.ErrOrderCurrencyMismatch
		}
		if err := ValidateOrderCartItems(bpo.OrderCartItems); err != nil {
			return nil, err
		}
	}

	inquiryIDs := bulks.InquiryIDs()
	inquiries := make(models.Inquiries, 0, len(inquiryIDs))
	if err := repo.db.Find(&inquiries, "id IN ?", inquiryIDs).Error; err != nil {
		return nil, err
	}
	dbInquiryIDs := inquiries.IDs()
	for _, id := range inquiryIDs {
		if !helper.StringContains(dbInquiryIDs, id) {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "inquiry_id:%s", id)
		}
	}

	mapInquiryIdToInquiry := make(map[string]*models.Inquiry, len(inquiries))
	for _, iq := range inquiries {
		mapInquiryIdToInquiry[iq.ID] = iq
	}

	addressIDs := bulks.AddressIDs()
	addresses := make([]*models.Address, 0, len(addressIDs))
	if err := query.New(repo.db, queryfunc.NewAddressBuilder(queryfunc.AddressBuilderOptions{})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("a.id IN ?", addressIDs)
		}).
		FindFunc(&addresses); err != nil {
		return nil, err
	}
	mapAddressIdToAddress := make(map[string]*models.Address, len(addresses))
	for _, address := range addresses {
		mapAddressIdToAddress[address.ID] = address
	}

	for _, bpo := range bulks {
		inquiry, ok := mapInquiryIdToInquiry[bpo.InquiryID]
		if !ok && bpo.InquiryID != "" {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "bulk_purchase_order_id:%s", bpo.ID)
		}
		address, ok := mapAddressIdToAddress[bpo.ShippingAddressID]
		if !ok && bpo.ShippingAddressID != "" {
			return nil, eris.Wrapf(errs.ErrAddressNotFound, "bulk_purchase_order_id:%s", bpo.ID)
		}
		var subTotalPrice = price.NewFromFloat(0)
		for _, item := range bpo.OrderCartItems {
			subTotalPrice = subTotalPrice.Add(item.TotalPrice)
		}

		bpo.SubTotal = subTotalPrice.ToPtr()
		if values.Float64Value(bpo.TaxPercentage) == 0 {
			if inquiry != nil {
				bpo.TaxPercentage = inquiry.TaxPercentage
			}
		}

		if bpo.CommercialInvoice != nil {
			bpo.TaxPercentage = bpo.CommercialInvoice.TaxPercentage
		}

		if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment {
			bpo.FirstPaymentType = paymentType
		} else {
			bpo.FinalPaymentType = paymentType
		}

		bpo.Inquiry = inquiry
		bpo.ShippingAddress = address

		if err := bpo.UpdatePrices(); err != nil {
			return nil, err
		}
	}
	return bulks, nil
}

func (repo *OrderCartRepo) CheckoutOrders(req *models.OrderCartCheckoutRequest) (*models.OrderCartCheckoutResponse, error) {
	if len(req.PurchaseOrderCartItemIDs) == 0 && len(req.BulkOrderIDs) == 0 {
		return nil, errs.ErrOrderEmpty
	}
	var currency enums.Currency
	purchaseOrders, err := repo.GetPurchaseOrderAndOrderItemsByOrderCartItemIDs(req.PurchaseOrderCartItemIDs, req.GetUserID())
	if err != nil {
		return nil, err
	}
	if len(purchaseOrders) > 0 {
		currency = purchaseOrders[0].Currency
		if err := ValidateCheckoutPurchaseOrders(purchaseOrders, req.PaymentType); err != nil {
			return nil, err
		}
		for _, po := range purchaseOrders {
			if err := ValidateOrderCartItems(po.OrderCartItems); err != nil {
				return nil, err
			}
		}
	}
	bulks, err := repo.GetBulkAndOrderItemsByBulkIDs(req.BulkOrderIDs, req.GetUserID())
	if err != nil {
		return nil, err
	}
	if len(bulks) > 0 {
		currency = bulks[0].Currency
		if err := ValidateCheckoutBulks(bulks, req.PaymentType); err != nil {
			return nil, err
		}
		for _, bpo := range bulks {
			if err := ValidateOrderCartItems(bpo.OrderCartItems); err != nil {
				return nil, err
			}
		}
	}
	if len(purchaseOrders) > 0 && len(bulks) > 0 {
		if purchaseOrders[0].Currency != bulks[0].Currency {
			return nil, errs.ErrOrderCurrencyMismatch
		}
	}

	var purchaseOrdersToUpdate = make([]*models.PurchaseOrder, 0, len(purchaseOrders))
	var bulksToUpdate = make([]*models.BulkPurchaseOrder, 0, len(bulks))
	var inquiryIDsToUpdate = make([]string, 0, len(purchaseOrders))
	var orderCartItemIDsToUpdate []string
	var paymentTransaction models.PaymentTransaction
	var checkoutSessionID = helper.GenerateCheckoutSessionID()
	var totalAmount price.Price
	var transactionRefID = helper.GeneratePaymentTransactionReferenceID()

	if req.PaymentType == enums.PaymentTypeBankTransfer {
		for _, po := range purchaseOrders {
			totalAmount = totalAmount.AddPtr(po.TotalPrice)

			po.Status = enums.PurchaseOrderStatusWaitingConfirm
			po.TransactionRefID = req.TransactionRefID
			po.TransactionAttachment = req.TransactionAttachment
			po.TransferedAt = values.Int64(time.Now().Unix())
			po.CheckoutSessionID = checkoutSessionID
			po.PaymentTransactionReferenceID = transactionRefID

			purchaseOrdersToUpdate = append(purchaseOrdersToUpdate, po)
		}
		for _, bpo := range bulks {
			if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment {
				totalAmount = totalAmount.AddPtr(bpo.FirstPaymentTotal)

				bpo.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirm
				bpo.FirstPaymentTransferedAt = values.Int64(time.Now().Unix())
				bpo.FirstPaymentTransactionRefID = req.TransactionRefID
				bpo.FirstPaymentTransactionAttachment = req.TransactionAttachment
				bpo.FirstPaymentTransactionReferenceID = transactionRefID
				bpo.FirstPaymentCheckoutSessionID = checkoutSessionID
			} else {
				totalAmount = totalAmount.AddPtr(bpo.FinalPaymentTotal)

				bpo.TrackingStatus = enums.BulkPoTrackingStatusFinalPaymentConfirm
				bpo.FinalPaymentTransferedAt = values.Int64(time.Now().Unix())
				bpo.FinalPaymentTransactionRefID = req.TransactionRefID
				bpo.FinalPaymentTransactionAttachment = req.TransactionAttachment
				bpo.FinalPaymentTransactionReferenceID = transactionRefID
				bpo.FinalPaymentCheckoutSessionID = checkoutSessionID
			}
			bulksToUpdate = append(bulksToUpdate, bpo)
		}

		paymentTransaction = models.PaymentTransaction{
			ReferenceID:          transactionRefID,
			PaidAmount:           totalAmount.ToPtr(),
			PaymentType:          req.PaymentType,
			UserID:               req.GetUserID(),
			TransactionRefID:     req.TransactionRefID,
			Status:               enums.PaymentStatusWaitingConfirm,
			TotalAmount:          totalAmount.ToPtr(),
			Currency:             currency,
			CheckoutSessionID:    checkoutSessionID,
			PurchaseOrderIDs:     purchaseOrders.IDs(),
			BulkPurchaseOrderIDs: bulks.IDs(),
			Metadata: &models.PaymentTransactionMetadata{
				PurchaseOrderIDs:     purchaseOrders.IDs(),
				BulkPurchaseOrderIDs: bulks.IDs(),
			},
		}
		if req.TransactionAttachment != nil {
			paymentTransaction.Attachments = &models.Attachments{req.TransactionAttachment}
		}
	}
	if req.PaymentType == enums.PaymentTypeCard {
		var user models.User
		if err = repo.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", req.GetUserID()).Error; err != nil {
			return nil, err
		}
		stripeConfig, err := stripehelper.GetCurrencyConfig(currency)
		if err != nil {
			return nil, err
		}
		for _, po := range purchaseOrders {
			totalAmount = totalAmount.AddPtr(po.TotalPrice)
		}
		for _, bpo := range bulks {
			if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment {
				totalAmount = totalAmount.AddPtr(bpo.FirstPaymentTotal)
			} else {
				totalAmount = totalAmount.AddPtr(bpo.FinalPaymentTotal)
			}
		}
		var stripeParams = stripehelper.CreatePaymentIntentParams{
			Amount:                  totalAmount.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64(),
			Currency:                currency,
			PaymentMethodID:         req.PaymentMethodID,
			CustomerID:              user.StripeCustomerID,
			IsCaptureMethodManually: false,
			Description:             fmt.Sprintf("Charges for checkout session %s", checkoutSessionID),
			PaymentMethodTypes:      []string{string(req.PaymentType)},
			Metadata: map[string]string{
				"purchase_order_cart_item_ids": strings.Join(req.PurchaseOrderCartItemIDs, ","),
				"bulk_order_ids":               strings.Join(req.BulkOrderIDs, ","),
				"user_id":                      req.GetUserID(),
				"checkout_session_id":          checkoutSessionID,
				"action_source":                string(stripehelper.ActionSourceOrderCartPayment),
			},
		}
		pi, err := stripehelper.GetInstance().CreatePaymentIntent(stripeParams)
		if err != nil {
			return nil, err
		}
		if pi.Status != stripe.PaymentIntentStatusSucceeded {
			if pi.NextAction != nil {
				intent, err := stripehelper.GetInstance().ConfirmPaymentIntent(stripehelper.ConfirmPaymentIntentParams{
					PaymentIntentID: pi.ID,
					ReturnURL:       fmt.Sprintf("%s/api/v1/callback/stripe/payment_intents/order_cart/%s/confirm?po_items=%s&bulks=%s", repo.db.Configuration.ServerBaseURL, checkoutSessionID, strings.Join(req.PurchaseOrderCartItemIDs, ","), strings.Join(bulks.IDs(), ",")),
				})
				if err != nil {
					return nil, err
				}

				if intent.Status == stripe.PaymentIntentStatusSucceeded {
					goto PaymentSuccess
				}
				return &models.OrderCartCheckoutResponse{
					PaymentIntentNextAction:   intent.NextAction,
					PaymentIntentClientSecret: intent.ClientSecret,
				}, nil

			} else {
				return nil, eris.Errorf("Payment error with status %s", pi.Status)
			}
		}
	PaymentSuccess:
		for _, po := range purchaseOrders {
			po.PaymentIntentID = pi.ID
			po.Status = enums.PurchaseOrderStatusPaid
			po.TransferedAt = values.Int64(time.Now().Unix())
			po.MarkAsPaidAt = values.Int64(time.Now().Unix())
			po.CheckoutSessionID = checkoutSessionID
			po.PaymentTransactionReferenceID = transactionRefID

			// if len(po.Quotations) > 0 {
			// 	sampleQuotation, _ := lo.Find(po.Quotations, func(item *models.InquiryQuotationItem) bool {
			// 		return item.Type == enums.InquiryTypeSample
			// 	})
			// 	if sampleQuotation != nil {
			// 		updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
			// 		updates.StartDate = updates.MarkAsPaidAt
			// 		updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
			// 	}
			// }

			purchaseOrdersToUpdate = append(purchaseOrdersToUpdate, po)
			if po.InquiryID != "" {
				inquiryIDsToUpdate = append(inquiryIDsToUpdate, po.InquiryID)
			}
			orderCartItemIDsToUpdate = append(orderCartItemIDsToUpdate, models.OrderCartItems(po.OrderCartItems).IDs()...)
		}
		for _, bpo := range bulks {
			if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment {
				bpo.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirmed
				bpo.FirstPaymentTransferedAt = values.Int64(time.Now().Unix())
				bpo.FirstPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
				bpo.FirstPaymentTransactionReferenceID = transactionRefID
				bpo.FirstPaymentIntentID = pi.ID
				bpo.FirstPaymentCheckoutSessionID = checkoutSessionID
			} else {
				bpo.TrackingStatus = enums.BulkPoTrackingStatusFinalPaymentConfirmed
				bpo.FinalPaymentTransferedAt = values.Int64(time.Now().Unix())
				bpo.FinalPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
				bpo.FinalPaymentTransactionReferenceID = transactionRefID
				bpo.FinalPaymentIntentID = pi.ID
				bpo.FinalPaymentCheckoutSessionID = checkoutSessionID
				orderCartItemIDsToUpdate = append(orderCartItemIDsToUpdate, models.OrderCartItems(bpo.OrderCartItems).IDs()...)
			}
			// if len(order.AdminQuotations) > 0 {
			// 	var bulkQuotation = order.GetBulkQuotation()
			// 	if bulkQuotation != nil {
			// 		updates.LeadTime = int(values.Int64Value(bulkQuotation.LeadTime))
			// 		updates.StartDate = updates.FirstPaymentMarkAsPaidAt
			// 		updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
			// 	}
			// }
			bulksToUpdate = append(bulksToUpdate, bpo)
		}
		paymentTransaction = models.PaymentTransaction{
			PaymentIntentID:      pi.ID,
			ReferenceID:          transactionRefID,
			PaidAmount:           totalAmount.ToPtr(),
			PaymentType:          req.PaymentType,
			UserID:               req.GetUserID(),
			Status:               enums.PaymentStatusPaid,
			TotalAmount:          totalAmount.ToPtr(),
			Currency:             currency,
			CheckoutSessionID:    checkoutSessionID,
			PurchaseOrderIDs:     purchaseOrders.IDs(),
			BulkPurchaseOrderIDs: bulks.IDs(),
			TransactionType:      enums.TransactionTypeCredit,
			Metadata: &models.PaymentTransactionMetadata{
				PurchaseOrderIDs:     purchaseOrders.IDs(),
				BulkPurchaseOrderIDs: bulks.IDs(),
			},
		}
	}

	if err := repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&paymentTransaction).Error; err != nil {
			return err
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
		return nil, err
	}

	return &models.OrderCartCheckoutResponse{
		PaymentTransaction: &paymentTransaction,
		CheckoutSessionID:  checkoutSessionID,
		GetOrderCartResponse: models.GetOrderCartResponse{
			PurchaseOrders:     purchaseOrdersToUpdate,
			BulkPurchaseOrders: bulksToUpdate,
		},
	}, nil
}

func (repo *OrderCartRepo) GetOrderCartCheckoutInfo(req *models.OrderCartGetCheckoutInfoRequest) (*models.OrderCartGetCheckoutInfoResponse, error) {
	var paymentTx models.PaymentTransaction
	if err := repo.db.First(&paymentTx, "checkout_session_id = ?", req.CheckoutSessionID).Error; err != nil {
		if repo.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrOrderCheckoutNotFound
		}
		return nil, err
	}

	var purchaseOrders []*models.PurchaseOrder
	var poCartBuilder = queryfunc.NewOrderCartPurchaseOrderBuilder(queryfunc.OrderCartPurchaseOrderBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: req.GetRole(),
		},
		IncludeInquiry:    true,
		IncludeAddress:    true,
		IncludeCollection: true,
	})
	if err := query.New(repo.db, poCartBuilder).WhereFunc(func(builder *query.Builder) {
		builder.Where("po.checkout_session_id = ?", req.CheckoutSessionID)
	}).FindFunc(&purchaseOrders); err != nil {
		return nil, err
	}

	var bulks []*models.BulkPurchaseOrder
	var bpoCartBuilder = queryfunc.NewOrderCartBulkBuilder(queryfunc.OrderCartBulkBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: req.GetRole(),
		},
		IncludeInquiry:    true,
		IncludeAddress:    true,
		IncludeCollection: true,
	})
	if err := query.New(repo.db, bpoCartBuilder).WhereFunc(func(builder *query.Builder) {
		builder.Where("bpo.first_payment_checkout_session_id = @checkoutSessionID OR bpo.final_payment_checkout_session_id = @checkoutSessionID", sql.Named("checkoutSessionID", req.CheckoutSessionID))
	}).FindFunc(&bulks); err != nil {
		return nil, err
	}

	if len(purchaseOrders) == 0 && len(bulks) == 0 {
		return nil, errs.ErrOrderCheckoutNotFound
	}

	return &models.OrderCartGetCheckoutInfoResponse{
		GetOrderCartResponse: models.GetOrderCartResponse{
			PurchaseOrders:     purchaseOrders,
			BulkPurchaseOrders: bulks,
		},
		PaymentTransaction: paymentTx,
	}, nil

}

func (repo *OrderCartRepo) GetPurchaseOrderAndOrderItemsByOrderCartItemIDs(orderItemIDs []string, userID string) (models.PurchaseOrders, error) {
	orderItems := make(models.OrderCartItems, 0, len(orderItemIDs))
	if err := repo.db.Find(&orderItems, "id IN ?", orderItemIDs).Error; err != nil {
		return nil, err
	}
	dbOrderItemIDs := orderItems.IDs()
	for _, id := range orderItemIDs {
		if !helper.StringContains(dbOrderItemIDs, id) {
			return nil, eris.Wrapf(errs.ErrOrderItemNotFound, "order_item_id:%s", id)
		}
	}
	purchaseOrderIDs := make([]string, 0, len(orderItems))
	for _, item := range orderItems {
		if item.PurchaseOrderID == "" || item.BulkPurchaseOrderID != "" {
			return nil, eris.Wrapf(errs.ErrOrderInvalid, "order_item_id:%s", item.ID)
		}
		purchaseOrderIDs = append(purchaseOrderIDs, item.PurchaseOrderID)
	}

	mapPurchaseOrderIdToOrderItems := make(map[string]models.OrderCartItems, len(orderItems))
	for _, item := range orderItems {
		_, ok := mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID]
		if ok {
			mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID] = append(mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID], item)
		} else {
			mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID] = []*models.OrderCartItem{item}
		}
	}

	purchaseOrders := make(models.PurchaseOrders, 0, len(purchaseOrderIDs))
	if len(purchaseOrderIDs) == 0 {
		return purchaseOrders, nil
	}
	if err := repo.db.Find(&purchaseOrders, "id IN ? AND user_id = ?", purchaseOrderIDs, userID).Error; err != nil {
		return nil, err
	}
	dbPurchaseOrderIDs := purchaseOrders.IDs()
	for _, id := range purchaseOrderIDs {
		if !helper.StringContains(dbPurchaseOrderIDs, id) {
			return nil, eris.Wrapf(errs.ErrPONotFound, "purchase_order_id:%s", id)
		}
	}

	for _, po := range purchaseOrders {
		var items, ok = mapPurchaseOrderIdToOrderItems[po.ID]
		if !ok {
			return nil, eris.Wrapf(errs.ErrOrderEmpty, "purchase_order_id:%s", po.ID)
		}
		po.OrderCartItems = items
	}
	return purchaseOrders, nil
}

func (repo *OrderCartRepo) GetBulkAndOrderItemsByBulkIDs(bulkIDs []string, userID string) (models.BulkPurchaseOrders, error) {
	bulks := make(models.BulkPurchaseOrders, 0, len(bulkIDs))
	if len(bulkIDs) == 0 {
		return bulks, nil
	}
	if err := repo.db.Find(&bulks, "id IN ? AND user_id = ?", bulkIDs, userID).Error; err != nil {
		return nil, err
	}
	dbBulkIDs := bulks.IDs()
	for _, id := range bulkIDs {
		if !helper.StringContains(dbBulkIDs, id) {
			return nil, eris.Wrapf(errs.ErrBulkPoNotFound, "bulk_purchase_order_id:%s", id)
		}
	}
	var orderItems models.OrderCartItems
	if err := repo.db.Find(&orderItems, "bulk_purchase_order_id IN ?", bulkIDs).Error; err != nil {
		return nil, err
	}

	mapBulkIdToOrderItems := make(map[string]models.OrderCartItems, len(orderItems))
	for _, item := range orderItems {
		_, ok := mapBulkIdToOrderItems[item.BulkPurchaseOrderID]
		if ok {
			mapBulkIdToOrderItems[item.BulkPurchaseOrderID] = append(mapBulkIdToOrderItems[item.BulkPurchaseOrderID], item)
		} else {
			mapBulkIdToOrderItems[item.BulkPurchaseOrderID] = []*models.OrderCartItem{item}
		}
	}

	for _, bpo := range bulks {
		var items, ok = mapBulkIdToOrderItems[bpo.ID]
		if !ok {
			return nil, eris.Wrapf(errs.ErrOrderEmpty, "bulk_purchase_order_id:%s", bpo.ID)
		}
		bpo.OrderCartItems = items
	}

	return bulks, nil
}
func (repo *OrderCartRepo) GetPurchaseOrderAndOrderItemsByPurchaseOrderIDs(poIDs []string, userID string) (models.PurchaseOrders, error) {
	purchaseOrders := make(models.PurchaseOrders, 0, len(poIDs))
	if len(poIDs) == 0 {
		return purchaseOrders, nil
	}
	if err := repo.db.Find(&purchaseOrders, "id IN ? AND user_id = ?", poIDs, userID).Error; err != nil {
		return nil, err
	}
	dbPurchaseOrderIDs := purchaseOrders.IDs()
	for _, id := range poIDs {
		if !helper.StringContains(dbPurchaseOrderIDs, id) {
			return nil, eris.Wrapf(errs.ErrPONotFound, "purchase_order_id:%s", id)
		}
	}
	var orderItems []*models.OrderCartItem
	if err := repo.db.Find(&orderItems, "purchase_order_id IN ?", poIDs).Error; err != nil {
		return nil, err
	}
	mapPurchaseOrderIdToOrderItems := make(map[string]models.OrderCartItems, len(orderItems))
	for _, item := range orderItems {
		_, ok := mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID]
		if ok {
			mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID] = append(mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID], item)
		} else {
			mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID] = []*models.OrderCartItem{item}
		}
	}
	for _, po := range purchaseOrders {
		var items, ok = mapPurchaseOrderIdToOrderItems[po.ID]
		if !ok {
			return nil, eris.Wrapf(errs.ErrOrderEmpty, "purchase_order_id:%s", po.ID)
		}
		po.OrderCartItems = items
	}

	return purchaseOrders, nil
}
func (repo *OrderCartRepo) GetPurchaseOrderAndOrderItemsByInquiryIDs(inquiryIDs []string, userID string) (models.PurchaseOrders, error) {
	purchaseOrders := make(models.PurchaseOrders, 0, len(inquiryIDs))
	if err := repo.db.Find(&purchaseOrders, "inquiry_id IN ? AND user_id = ?", inquiryIDs, userID).Error; err != nil {
		return nil, err
	}
	dbInquiryIDs := purchaseOrders.InquiryIDs()
	for _, id := range inquiryIDs {
		if !helper.StringContains(dbInquiryIDs, id) {
			return nil, eris.Wrapf(errs.ErrPONotFound, "inquiry_id:%s", id)
		}
	}
	var orderItems []*models.OrderCartItem
	if err := repo.db.Find(&orderItems, "purchase_order_id IN ?", purchaseOrders.IDs()).Error; err != nil {
		return nil, err
	}
	mapPurchaseOrderIdToOrderItems := make(map[string]models.OrderCartItems, len(orderItems))
	for _, item := range orderItems {
		_, ok := mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID]
		if ok {
			mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID] = append(mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID], item)
		} else {
			mapPurchaseOrderIdToOrderItems[item.PurchaseOrderID] = []*models.OrderCartItem{item}
		}
	}
	for _, po := range purchaseOrders {
		var items, ok = mapPurchaseOrderIdToOrderItems[po.ID]
		if !ok {
			return nil, eris.Wrapf(errs.ErrOrderEmpty, "purchase_order_id:%s", po.ID)
		}
		po.OrderCartItems = items
	}

	return purchaseOrders, nil
}

func (repo *OrderCartRepo) CreateBuyerPaymentLink(req *models.CreateBuyerPaymentLinkRequest) (*models.CreateBuyerPaymentLinkResponse, error) {
	if len(req.PurchaseOrderIDs) == 0 && len(req.BulkIDs) == 0 {
		return nil, errs.ErrOrderEmpty
	}
	var user models.User
	if err := repo.db.First(&user, "id = ? ", req.BuyerID).Error; err != nil {
		if repo.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	var currency enums.Currency
	var paymentType = enums.PaymentTypeCard

	purchaseOrders, err := repo.GetPurchaseOrderAndOrderItemsByPurchaseOrderIDs(req.PurchaseOrderIDs, req.BuyerID)
	if err != nil {
		return nil, err
	}
	if len(purchaseOrders) > 0 {
		currency = purchaseOrders[0].Currency
		if err := ValidateCheckoutPurchaseOrders(purchaseOrders, paymentType); err != nil {
			return nil, err
		}
		for _, po := range purchaseOrders {
			if err := ValidateOrderCartItems(po.OrderCartItems); err != nil {
				return nil, err
			}
		}
	}

	bulks, err := repo.GetBulkAndOrderItemsByBulkIDs(req.BulkIDs, req.BuyerID)
	if err != nil {
		return nil, err
	}
	if len(bulks) > 0 {
		currency = bulks[0].Currency
		if err := ValidateCheckoutBulks(bulks, paymentType); err != nil {
			return nil, err
		}
		for _, bpo := range bulks {
			if err := ValidateOrderCartItems(bpo.OrderCartItems); err != nil {
				return nil, err
			}
		}

	}
	if len(purchaseOrders) > 0 && len(bulks) > 0 {
		if purchaseOrders[0].Currency != bulks[0].Currency {
			return nil, errs.ErrOrderCurrencyMismatch
		}
	}

	stripeConfig, err := stripehelper.GetCurrencyConfig(currency)
	if err != nil {
		return nil, err
	}
	var lineItems []*stripe.PaymentLinkLineItemParams
	var totalTax = price.NewFromFloat(0)
	var totalTransactionFee = price.NewFromFloat(0)
	var totalShippingFee = price.NewFromFloat(0)

	if len(purchaseOrders) > 0 {
		for _, po := range purchaseOrders {
			for _, item := range po.OrderCartItems {
				priceItem, err := stripePrice.New(&stripe.PriceParams{
					Currency:   stripe.String(string(currency)),
					UnitAmount: stripe.Int64(item.UnitPrice.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
					ProductData: &stripe.PriceProductDataParams{
						Name: stripe.String(fmt.Sprintf("%s - %s - %s (Sample)", po.ProductName, item.Size, item.ColorName)),
					},
				})
				if err != nil {
					return nil, err
				}
				lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
					Price:    &priceItem.ID,
					Quantity: stripe.Int64(int64(item.Qty)),
				})
			}
			if po.Tax.GreaterThan(0) {
				totalTax = totalTax.AddPtr(po.Tax)
			}
			if po.ShippingFee.GreaterThan(0) {
				totalShippingFee = totalShippingFee.AddPtr(po.ShippingFee)
			}
			if po.TransactionFee.GreaterThan(0) {
				totalTransactionFee = totalTransactionFee.AddPtr(po.TransactionFee)
			}
		}
	}
	if len(bulks) > 0 {
		for _, bpo := range bulks {
			var itemPricePct float64
			if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment {
				itemPricePct = *bpo.FirstPaymentPercentage

				totalTransactionFee = totalTransactionFee.AddPtr(bpo.FirstPaymentTransactionFee)
			} else {
				itemPricePct = 100 - *bpo.FirstPaymentPercentage

				totalTransactionFee = totalTransactionFee.AddPtr(bpo.FinalPaymentTransactionFee)
				totalTax = totalTax.AddPtr(bpo.FinalPaymentTax)
				totalShippingFee = totalShippingFee.AddPtr(bpo.ShippingFee)
			}
			var itemPricePctString = price.NewFromFloat(itemPricePct).Decimal().String()
			for _, item := range bpo.OrderCartItems {
				var itemPrice = item.UnitPrice.MultipleFloat64(itemPricePct).DivInt(100)
				priceItem, err := stripePrice.New(&stripe.PriceParams{
					Currency:   stripe.String(string(currency)),
					UnitAmount: stripe.Int64(itemPrice.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
					ProductData: &stripe.PriceProductDataParams{
						Name: stripe.String(fmt.Sprintf("%s - %s - %s (%s%%)", bpo.ProductName, item.Size, item.ColorName, itemPricePctString)),
					},
				})
				if err != nil {
					return nil, err
				}
				lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
					Price:    &priceItem.ID,
					Quantity: stripe.Int64(int64(item.Qty)),
				})
			}
		}
	}
	if totalShippingFee.GreaterThan(0) {
		priceItemShippingFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalShippingFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Shipping Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemShippingFee.ID,
		})
	}
	if totalTransactionFee.GreaterThan(0) {
		priceItemTxFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalTransactionFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Transaction Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTxFee.ID,
		})
	}
	if totalTax.GreaterThan(0) {
		priceItemTax, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalTax.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Tax"),
			},
		})
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTax.ID,
		})
	}

	var checkoutSessionID = helper.GenerateCheckoutSessionID()
	var poItemsIDs []string
	for _, po := range purchaseOrders {
		itemIDs := lo.Map(po.OrderCartItems, func(item *models.OrderCartItem, idx int) string {
			return item.ID
		})
		poItemsIDs = append(poItemsIDs, itemIDs...)
	}
	var stripeParams = stripehelper.CreatePaymentLinkParams{
		Currency: currency,
		Metadata: map[string]string{
			"purchase_order_cart_item_ids": strings.Join(poItemsIDs, ","),
			"bulk_order_ids":               strings.Join(req.BulkIDs, ","),
			"user_id":                      req.BuyerID,
			"checkout_session_id":          checkoutSessionID,
			"action_source":                string(stripehelper.ActionSourceOrderCartPayment),
		},
		RedirectURL: fmt.Sprintf("%s/order-checkout?checkout_session_id=%s", repo.db.Configuration.BrandPortalBaseURL, checkoutSessionID),
		LineItems:   lineItems,
	}

	pl, err := stripehelper.GetInstance().CreatePaymentLink(stripeParams)
	if err != nil {
		return nil, err
	}
	var samplePaymentLink = helper.AddURLQuery(pl.URL,
		map[string]string{
			"prefilled_email": user.Email,
		},
	)
	return &models.CreateBuyerPaymentLinkResponse{
		PaymentLink: samplePaymentLink,
	}, nil
}

func (repo *OrderCartRepo) GetBuyerOrderCartPreview(req *models.GetBuyerOrderCartRequest) (*models.GetBuyerOrderCartResponse, error) {
	if len(req.InquiryIDs) == 0 && len(req.BulkIDs) == 0 {
		return nil, errs.ErrOrderEmpty
	}
	purchaseOrders, err := repo.GetPurchaseOrderAndOrderItemsByInquiryIDs(req.InquiryIDs, req.BuyerID)
	if err != nil {
		return nil, err
	}
	purchaseOrders, err = repo.ProcessPurchaseOrdersPreviewCheckout(purchaseOrders, req.BuyerID, enums.PaymentTypeCard)
	if err != nil {
		return nil, err
	}
	bulks, err := repo.GetBulkAndOrderItemsByBulkIDs(req.BulkIDs, req.BuyerID)
	if err != nil {
		return nil, err
	}
	bulks, err = repo.ProcessBulkPurchaseOrdersPreviewCheckout(bulks, req.BuyerID, enums.PaymentTypeCard)
	if err != nil {
		return nil, err
	}
	if len(purchaseOrders) > 0 && len(bulks) > 0 {
		if purchaseOrders[0].Currency != bulks[0].Currency {
			return nil, errs.ErrOrderCurrencyMismatch
		}
	}

	var bulkPendingPOs models.PurchaseOrders
	if len(bulks) > 0 {
		var bulkPendingPoIDs []string
		if err := repo.db.Model(&models.PurchaseOrder{}).Select("ID").Find(&bulkPendingPoIDs, "id IN ? AND status = ?", bulks.PurchaseOrderIDs(), enums.PurchaseOrderStatusPending).Error; err != nil {
			return nil, err
		}
		bulkPendingPOs, err = repo.GetPurchaseOrderAndOrderItemsByPurchaseOrderIDs(bulkPendingPoIDs, req.BuyerID)
		if err != nil {
			return nil, err
		}
		bulkPendingPOs, err = repo.ProcessPurchaseOrdersPreviewCheckout(bulkPendingPOs, req.BuyerID, enums.PaymentTypeCard)
		if err != nil {
			return nil, err
		}
		if len(bulkPendingPOs) > 0 {
			if bulkPendingPOs[0].Currency != bulks[0].Currency {
				return nil, errs.ErrOrderCurrencyMismatch
			}

			var mapPoIDtoPurchaseOrder = make(map[string]*models.PurchaseOrder, len(bulkPendingPOs))
			for _, po := range bulkPendingPOs {
				mapPoIDtoPurchaseOrder[po.ID] = po
			}
			for _, bpo := range bulks {
				po, ok := mapPoIDtoPurchaseOrder[bpo.PurchaseOrderID]
				if ok {
					bpo.PurchaseOrder = po
				}
			}

		}
	}

	if err := repo.db.Transaction(func(tx *gorm.DB) error {
		if len(purchaseOrders) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&purchaseOrders).Error; err != nil {
				return err
			}
		}
		if len(bulks) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&bulks).Error; err != nil {
				return err
			}
		}
		if len(bulkPendingPOs) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&bulkPendingPOs).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &models.GetBuyerOrderCartResponse{
		GetOrderCartResponse: models.GetOrderCartResponse{
			PurchaseOrders:     purchaseOrders,
			BulkPurchaseOrders: bulks,
		},
	}, nil
}

func ValidateCheckoutPurchaseOrders(purchaseOrders []*models.PurchaseOrder, paymentType enums.PaymentType) error {
	currency := purchaseOrders[0].Currency
	for _, po := range purchaseOrders {
		switch {
		case po.Status == enums.PurchaseOrderStatusPaid:
			return eris.Wrapf(errs.ErrPoIsAlreadyPaid, "purchase_order_id:%s", po.ID)
		case po.PaymentType != paymentType:
			return eris.Wrapf(errs.ErrOrderInvalid, "purchase_order_id:%s", po.ID)
		case po.TotalPrice.LessThanOrEqual(0):
			return eris.Wrapf(errs.ErrOrderInvalid, "purchase_order_id:%s", po.ID)
		case po.SubTotal.LessThanOrEqual(0):
			return eris.Wrapf(errs.ErrOrderInvalid, "purchase_order_id:%s", po.ID)
		case po.PaymentType != enums.PaymentTypeCard && po.PaymentType != enums.PaymentTypeBankTransfer:
			return eris.Wrapf(errs.ErrOrderInvalid, "purchase_order_id:%s", po.ID)
		case po.Currency == "" || po.Currency != currency:
			return eris.Wrapf(errs.ErrOrderCurrencyMismatch, "purchase_order_id:%s", po.ID)
		}
	}
	return nil
}

func ValidateCheckoutBulks(bulks []*models.BulkPurchaseOrder, paymentType enums.PaymentType) error {
	currency := bulks[0].Currency
	for _, bpo := range bulks {
		switch {
		case bpo.TrackingStatus != enums.BulkPoTrackingStatusFirstPayment && bpo.TrackingStatus != enums.BulkPoTrackingStatusFinalPayment:
			return eris.Wrapf(errs.ErrOrderInvalid, "bulk_purchase_order_id:%s", bpo.ID)
		case bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment && bpo.FirstPaymentType != paymentType:
			return eris.Wrapf(errs.ErrOrderInvalid, "bulk_purchase_order_id:%s", bpo.ID)
		case bpo.TrackingStatus == enums.BulkPoTrackingStatusFinalPayment && bpo.FinalPaymentType != paymentType:
			return eris.Wrapf(errs.ErrOrderInvalid, "bulk_purchase_order_id:%s", bpo.ID)
		case bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment && bpo.FirstPaymentTotal.LessThanOrEqual(0):
			return eris.Wrapf(errs.ErrOrderInvalid, "bulk_purchase_order_id:%s", bpo.ID)
		case bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPayment && bpo.FirstPaymentSubTotal.LessThanOrEqual(0):
			return eris.Wrapf(errs.ErrOrderInvalid, "bulk_purchase_order_id:%s", bpo.ID)
		case bpo.FirstPaymentType != enums.PaymentTypeCard && bpo.FirstPaymentType != enums.PaymentTypeBankTransfer:
			return eris.Wrapf(errs.ErrOrderInvalid, "purchase_order_id:%s", bpo.ID)
		case bpo.FinalPaymentTotal.LessThanOrEqual(0):
			return eris.Wrapf(errs.ErrOrderInvalid, "bulk_purchase_order_id:%s", bpo.ID)
		case bpo.FinalPaymentSubTotal.LessThanOrEqual(0):
			return eris.Wrapf(errs.ErrOrderInvalid, "bulk_purchase_order_id:%s", bpo.ID)
		case bpo.FinalPaymentType != enums.PaymentTypeCard && bpo.FinalPaymentType != enums.PaymentTypeBankTransfer:
			return eris.Wrapf(errs.ErrOrderInvalid, "purchase_order_id:%s", bpo.ID)
		case bpo.Currency == "" || bpo.Currency != currency:
			return eris.Wrapf(errs.ErrOrderCurrencyMismatch, "bulk_purchase_order_id:%s", bpo.ID)
		}
	}
	return nil
}

func ValidateOrderCartItems(orderItems []*models.OrderCartItem) error {
	for _, item := range orderItems {
		switch {
		case item.UnitPrice.LessThanOrEqual(0):
			return eris.Wrapf(errs.ErrOrderItemInvalid, "order_cart_item_id:%s, purchase_order_id:%s", item.ID, item.PurchaseOrderID)
		case item.TotalPrice.LessThanOrEqual(0):
			return eris.Wrapf(errs.ErrOrderItemInvalid, "order_cart_item_id:%s, purchase_order_id:%s", item.ID, item.PurchaseOrderID)
		case item.CheckoutSessionID != "":
			return eris.Wrapf(errs.ErrOrderItemIsAlreadyPaid, "order_cart_item_id:%s, purchase_order_id:%s", item.ID, item.PurchaseOrderID)
		}
	}
	return nil
}
