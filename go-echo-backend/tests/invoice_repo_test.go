package tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceRepo_PaginateInvoice(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/admin/invoices?page=1&limit=12&user_id=cjvso05ooc2b8f45a1mg", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6InN1cGVyX2FkbWluIiwiaXNzIjoiY2x1MHNkcmIyaGo5NzFtdmpjYWciLCJzdWIiOiJzdXBlcl9hZG1pbiJ9.f-Co_EV6cAwmaQWS6zRLOfrfDgRT2R5-avvCSbFELlE")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())
}

func TestInvoiceRepo_NextInvoiceNumber(t *testing.T) {
	var app = initApp("dev")
	next, err := repo.NewInvoiceRepo(app.DB).NextInvoiceNumber()

	assert.NoError(t, err)
	log.Println(next)
}

func TestInvoiceRepo_IsExitsInvoiceNumber(t *testing.T) {
	var app = initApp("local")
	next, err := repo.NewInvoiceRepo(app.DB).IsExitsInvoiceNumber(repo.InvoiceDetailsPrams{
		InvoiceNumber: aws.Int(10),
	})

	assert.NoError(t, err)
	log.Println(next)
}

func TestInvoiceRepo_CreateBulkFirstPaymentInvoice(t *testing.T) {
	var app = initApp("prod")

	next, err := repo.NewInvoiceRepo(app.DB).CreateBulkFirstPaymentInvoice(repo.CreateBulkFirstPaymentInvoiceParams{
		BulkPurchaseOrderID: "BPO-MVMW-47823",
	})

	assert.NoError(t, err)
	log.Println(next)
}

func TestInvoiceRepo_CreateBulkFinalPaymentInvoice(t *testing.T) {
	var app = initApp("prod")

	next, err := repo.NewInvoiceRepo(app.DB).CreateBulkFinalPaymentInvoice(repo.CreateBulkFinalPaymentInvoiceParams{
		BulkPurchaseOrderID: "BPO-JGME-52759",
	})

	assert.NoError(t, err)
	log.Println(next)
}

func TestInvoiceRepo_CreateMultiplePurchaseOrderInvoice(t *testing.T) {
	var app = initApp("prod")

	next, err := repo.NewInvoiceRepo(app.DB).CreateMultiplePurchaseOrderInvoice(repo.CreateMultiplePurchaseOrderInvoiceParams{
		CheckoutSessionID: "CK-SAWF-99704",
	})

	assert.NoError(t, err)
	log.Println(next)
}

func TestInvoiceRepo_CreatePurchaseOrderInvoice(t *testing.T) {
	var app = initApp("prod")
	// app.DB.AutoMigrate(&models.Invoice{})

	next, err := repo.NewInvoiceRepo(app.DB).CreatePurchaseOrderInvoice(repo.CreatePurchaseOrderInvoiceParams{
		PurchaseOrderID: "PO-KXDR-68323",
	})

	assert.NoError(t, err)
	log.Println(next)
}

func TestInvoiceRepo_CreateBulkDebitNotes(t *testing.T) {
	var app = initApp("dev")
	// app.DB.AutoMigrate(&models.Invoice{})

	next, err := repo.NewInvoiceRepo(app.DB).CreateBulkDebitNotes(repo.CreateBulkDebitNotesParams{
		BulkPurchaseOrderID: "co143dou3mhqtlakj2h0",
		ReCreate:            true,
	})

	assert.NoError(t, err)
	log.Println(next)
}
func TestInvoiceRepo_CreateBulkCommercialInvoice(t *testing.T) {
	var app = initApp("dev")

	invoice, err := repo.NewInvoiceRepo(app.DB).CreateBulkCommercialInvoice(repo.CreateBulkCommercialInvoiceParams{
		BulkPurchaseOrderID: "co2d4niolqn9ijh9c7rg",
		ReCreate:            true,
	})

	assert.NoError(t, err)

	helper.PrintJSON(invoice.CommercialInvoice)
}
