package tests

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPaymentTransaction_PaymentTransactions(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/admin/payment_transactions?page=1&limit=12", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwiYXVkIjoic3VwZXJfYWRtaW4iLCJpc3MiOiJjbWlmM2JqYjJoamNkcW9xNmlkMCIsInN1YiI6InN1cGVyX2FkbWluIn0.7d8YdxKwyJ5dfMwHrnPdDeP76FMjQEiQf3VQFD_KFnE")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestPaymentTransaction_PaymentTransactionList(t *testing.T) {
	var app = initApp()
	app.DB.DB.AutoMigrate(&models.PaymentTransaction{})

	resp := repo.NewPaymentTransactionRepo(app.DB).
		PaginatePaymentTransactions(repo.PaginatePaymentTransactionsParams{
			JwtClaimsInfo:  *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
			IncludeDetails: true,
			PaginationParams: models.PaginationParams{
				Limit:   10,
				Keyword: "IQ-IZPX-94989",
			},
		})

	helper.PrintJSON(resp)
}

func TestPaymentTransaction_GetPaymentTransaction(t *testing.T) {
	var app = initApp("dev")

	resp, err := repo.NewPaymentTransactionRepo(app.DB).
		GetPaymentTransaction(repo.GetPaymentTransactionsParams{
			PaymentTransactionID: "cm0r6srb2hj5e4cnpd60",
			JwtClaimsInfo:        *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
			IncludeInvoice:       true,
			IncludeDetails:       true,
		})
	assert.NoError(t, err)
	helper.PrintJSON(resp)
}

func TestPaymentTransaction_ApprovePaymentTransaction(t *testing.T) {
	var app = initApp("dev")

	_, err := repo.NewPaymentTransactionRepo(app.DB).
		ApprovePaymentTransactions(repo.GetPaymentTransactionsParams{
			PaymentTransactionID: "clobkljb2hja8pk087eg",
			JwtClaimsInfo:        *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		})
	assert.NoError(t, err)
}

func TestPaymentTransaction_RejectPaymentTransactions(t *testing.T) {
	var app = initApp("dev")

	result, err := repo.NewPaymentTransactionRepo(app.DB).
		RejectPaymentTransactions(repo.GetPaymentTransactionsParams{
			PaymentTransactionID: "cm0r6srb2hj5e4cnpd60",
			Note:                 "no",
			JwtClaimsInfo:        *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		})

	helper.PrintJSON(result)
	assert.NoError(t, err)
}

func TestPaymentTransaction_ToExcel(t *testing.T) {
	var app = initApp("dev")

	var result = repo.NewPaymentTransactionRepo(app.DB).PaginatePaymentTransactions(repo.PaginatePaymentTransactionsParams{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleLeader),
		PaginationParams: models.PaginationParams{
			Limit: 100,
		},
		IncludeDetails: true,
	})

	var paymentTransactions models.PaymentTransactions = result.Records.([]*models.PaymentTransaction)

	data, err := paymentTransactions.ToExcel()
	assert.NoError(t, err)

	ioutil.WriteFile("test.xlsx", data, 0664)
}
