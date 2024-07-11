package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/samber/lo"
)

type PaymentTransactionAlias struct {
	*models.PaymentTransaction
	User *models.User `gorm:"embedded;embeddedPrefix:u__"`
}

type PaymentTransactionBuilderOptions struct {
	QueryBuilderOptions
	IncludeDetails bool
	IncludeInvoice bool
}

func NewPaymentTransactionBuilder(options PaymentTransactionBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*,
	u.id as u__id,
	u.name as u__name,
	u.email as u__email,
	u.avatar as u__avatar

	FROM payment_transactions p
	
	LEFT JOIN users u ON p.user_id = u.id
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM payment_transactions p
	LEFT JOIN users u ON p.user_id = u.id
	`

	return NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithOrderBy("p.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.PaymentTransaction, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var purchaseOrderIDs []string
			var bulkPurchaseOrderIDs []string
			var invoiceNumbers []int

			for rows.Next() {
				var alias PaymentTransactionAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				alias.PaymentTransaction.User = alias.User

				if len(alias.PurchaseOrderIDs) > 0 {
					purchaseOrderIDs = append(purchaseOrderIDs, alias.PurchaseOrderIDs...)
				}

				if len(alias.BulkPurchaseOrderIDs) > 0 {
					bulkPurchaseOrderIDs = append(bulkPurchaseOrderIDs, alias.BulkPurchaseOrderIDs...)
				}

				if !lo.Contains(invoiceNumbers, alias.PaymentTransaction.InvoiceNumber) && alias.PaymentTransaction.InvoiceNumber > 0 {
					invoiceNumbers = append(invoiceNumbers, alias.PaymentTransaction.InvoiceNumber)
				}

				records = append(records, alias.PaymentTransaction)
			}

			var wg sync.WaitGroup

			if options.IncludeInvoice && len(invoiceNumbers) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var invoices []*models.Invoice
					query.New(db, NewInvoiceBuilder(InvoiceBuilderOptions{})).
						Where("iv.invoice_number IN ?", invoiceNumbers).
						FindFunc(&invoices)

					for _, invoice := range invoices {
						for _, record := range records {
							if record.InvoiceNumber == invoice.InvoiceNumber {
								record.Invoice = invoice
							}
						}
					}
				}()
			}

			if options.IncludeDetails {

				if len(purchaseOrderIDs) > 0 {
					wg.Add(1)
					go func() {
						defer wg.Done()

						var purchaseOrders []*models.PurchaseOrder
						query.New(db, NewPurchaseOrderBuilder(PurchaseOrderBuilderOptions{})).
							Where("po.id IN ?", purchaseOrderIDs).
							FindFunc(&purchaseOrders)

						for _, purchaseOrder := range purchaseOrders {
							for _, record := range records {
								if len(record.PurchaseOrderIDs) > 0 {
									if lo.Contains(record.PurchaseOrderIDs, purchaseOrder.ID) {
										record.PurchaseOrders = append(record.PurchaseOrders, purchaseOrder)
									}
								}
							}
						}
					}()
				}

				if len(bulkPurchaseOrderIDs) > 0 {
					wg.Add(1)
					go func() {
						defer wg.Done()

						var bulkPurchaseOrders []*models.BulkPurchaseOrder
						query.New(db, NewBulkPurchaseOrderBuilder(BulkPurchaseOrderBuilderOptions{})).
							Where("bpo.id IN ?", bulkPurchaseOrderIDs).
							FindFunc(&bulkPurchaseOrders)

						for _, bulkPurchaseOrder := range bulkPurchaseOrders {
							for _, record := range records {
								if len(record.BulkPurchaseOrderIDs) > 0 {
									if lo.Contains(record.BulkPurchaseOrderIDs, bulkPurchaseOrder.ID) {
										record.BulkPurchaseOrders = append(record.BulkPurchaseOrders, bulkPurchaseOrder)
									}
								}
							}
						}
					}()
				}

			}
			wg.Wait()

			return records, nil
		})
}

func NewBuyerDashboardPaymentTransactionBuilder(options PaymentTransactionBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*
	FROM payment_transactions p
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1
	FROM payment_transactions p
	`

	return NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithOrderBy("p.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var transactions = make([]*models.PaymentTransaction, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var purchaseOrderIDs []string
			var bulkPurchaseOrderIDs []string

			for rows.Next() {
				var alias PaymentTransactionAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error %v", err)
					continue
				}

				if alias.BulkPurchaseOrderID != "" && !helper.StringContains(bulkPurchaseOrderIDs, alias.BulkPurchaseOrderID) {
					bulkPurchaseOrderIDs = append(bulkPurchaseOrderIDs, alias.BulkPurchaseOrderID)
				}

				for _, id := range alias.PurchaseOrderIDs {
					if id != "" && !helper.StringContains(purchaseOrderIDs, id) {
						purchaseOrderIDs = append(purchaseOrderIDs, id)
					}
				}

				transactions = append(transactions, alias.PaymentTransaction)
			}

			var wg sync.WaitGroup

			if len(purchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var purchaseOrders []*models.PurchaseOrder
					_ = query.New(db, NewPurchaseOrderBuilder(PurchaseOrderBuilderOptions{})).
						Where("po.id IN ?", purchaseOrderIDs).
						FindFunc(&purchaseOrders)

					for _, purchaseOrder := range purchaseOrders {
						for _, record := range transactions {
							if len(record.PurchaseOrderIDs) > 0 {
								if lo.Contains(record.PurchaseOrderIDs, purchaseOrder.ID) {
									record.PurchaseOrders = append(record.PurchaseOrders, purchaseOrder)
								}
								continue
							}
							if record.PurchaseOrderID == purchaseOrder.ID {
								record.PurchaseOrder = purchaseOrder
							}
						}
					}
				}()
			}

			if len(bulkPurchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var bulkPurchaseOrders []*models.BulkPurchaseOrder
					_ = query.New(db, NewBulkPurchaseOrderBuilder(BulkPurchaseOrderBuilderOptions{})).
						Where("bpo.id IN ?", bulkPurchaseOrderIDs).
						FindFunc(&bulkPurchaseOrders)

					for _, bulkPurchaseOrder := range bulkPurchaseOrders {
						for _, record := range transactions {
							if record.BulkPurchaseOrderID == bulkPurchaseOrder.ID {
								record.BulkPurchaseOrder = bulkPurchaseOrder
							}
						}
					}
				}()
			}

			wg.Wait()

			var pendings []*models.BuyerDataAnalyticPendingPayment
			for _, trx := range transactions {
				pending := &models.BuyerDataAnalyticPendingPayment{
					Amount:    trx.TotalAmount.SubPtr(trx.PaidAmount).ToFloat64(),
					Currency:  trx.Currency,
					Milestone: trx.Milestone,
				}
				if trx.BulkPurchaseOrder != nil {
					pending.ID = trx.BulkPurchaseOrderID
					pending.OrderID = trx.BulkPurchaseOrder.ReferenceID
					if trx.BulkPurchaseOrder.Inquiry != nil {
						pending.ProductName = trx.BulkPurchaseOrder.Inquiry.Title
						pending.Attachments = trx.BulkPurchaseOrder.Inquiry.Attachments
						pending.Quantity = *trx.BulkPurchaseOrder.Inquiry.Quantity
					}
				} else if trx.PurchaseOrder != nil {
					pending.ID = trx.PurchaseOrderID
					pending.OrderID = trx.PurchaseOrder.ReferenceID
					if trx.PurchaseOrder.Inquiry != nil {
						pending.ProductName = trx.PurchaseOrder.Inquiry.Title
						pending.Attachments = trx.PurchaseOrder.Inquiry.Attachments
						pending.Quantity = *trx.PurchaseOrder.Inquiry.Quantity
					}
				}
				pendings = append(pendings, pending)
			}

			return pendings, nil
		})
}
