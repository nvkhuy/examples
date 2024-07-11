package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type InvoiceAlias struct {
	*models.Invoice
}

type InvoiceBuilderOptions struct {
	QueryBuilderOptions
}

func NewInvoiceBuilder(options InvoiceBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ iv.*
	FROM invoices iv
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM invoices iv
	`
	var orderBy = "iv.updated_at DESC"
	return NewBuilder(rawSQL, countSQL).
		WithOrderBy(orderBy).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Invoice, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var paymentTransactionIDs []string

			for rows.Next() {
				var alias InvoiceAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				paymentTransactionIDs = append(paymentTransactionIDs, alias.PaymentTransactionID)
				records = append(records, alias.Invoice)
			}
			var wg = new(sync.WaitGroup)

			if len(paymentTransactionIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var payments []*models.PaymentTransaction
					if err := db.Find(&payments, "id IN ?", paymentTransactionIDs).Error; err != nil {
						return
					}
					for _, payment := range payments {
						for _, record := range records {
							if record.PaymentTransactionID == payment.ID {
								record.PaymentTransaction = payment
							}
						}
					}
				}()
			}

			wg.Wait()
			return &records, nil
		})
}
