package tests

import (
	"bytes"
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

func TestBulkPurchaseOrderRepo_FirstPayout(t *testing.T) {
	var app = initApp("dev")

	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = bytes.NewBuffer([]byte(`{"payout_percentage": 0}`))
	var req = httptest.NewRequest(echo.POST, "/api/v1/admin/seller_bulk_purchase_orders/cobon242vss1o35c3c90/first_payout", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNmE5bGgzNDN2MnFkaWlnanUwIiwiYXVkIjoic3RhZmYiLCJpc3MiOiJjb2ZwZXVubDNnYjB1Mmw2N2k2ZyIsInN1YiI6InN0YWZmOmRldiJ9.1swoL-c3-8q1MbWetv9BMymY3IRWHLKNsTSYm_IVYeQ")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestBulkPurchaseOrderRepo_SellerAllocations(t *testing.T) {
	var app = initApp("dev")

	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/admin/bulk_purchase_orders/coj0j15ik4rmpklmb0sg/seller_allocations", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNmE5bGgzNDN2MnFkaWlnanUwIiwiYXVkIjoic3RhZmYiLCJpc3MiOiJjb2ZwZXVubDNnYjB1Mmw2N2k2ZyIsInN1YiI6InN0YWZmOmRldiJ9.1swoL-c3-8q1MbWetv9BMymY3IRWHLKNsTSYm_IVYeQ")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestBulkPurchaseOrderRepo_RejectSellerQuotation(t *testing.T) {
	var app = initApp("dev")

	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = bytes.NewBuffer([]byte(`{"reject_reason":"ssssss"}`))
	var req = httptest.NewRequest(echo.DELETE, "/api/v1/admin/seller_bulk_purchase_orders/col268pu0vkk57pttipg/reject", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNmE5bGgzNDN2MnFkaWlnanUwIiwiYXVkIjoic3RhZmYiLCJpc3MiOiJjb2ZwZXVubDNnYjB1Mmw2N2k2ZyIsInN1YiI6InN0YWZmOmRldiJ9.1swoL-c3-8q1MbWetv9BMymY3IRWHLKNsTSYm_IVYeQ")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}
func TestBulkPurchaseOrderRepo_ApproveSellerQuotation(t *testing.T) {
	var app = initApp("dev")

	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.POST, "/api/v1/admin/seller_bulk_purchase_orders/col268pu0vkk57pttipg/approve", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNmE5bGgzNDN2MnFkaWlnanUwIiwiYXVkIjoic3RhZmYiLCJpc3MiOiJjb2ZwZXVubDNnYjB1Mmw2N2k2ZyIsInN1YiI6InN0YWZmOmRldiJ9.1swoL-c3-8q1MbWetv9BMymY3IRWHLKNsTSYm_IVYeQ")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestBulkPurchaseOrderRepo_SellerQuotation(t *testing.T) {
	var app = initApp("dev")

	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/admin/bulk_purchase_orders/coj2t9f2te3lr2bqmf8g/seller_quotations?statuses=new", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNjU4cGQ3cWZwZ3JscXFqOTNnIiwiYXVkIjoibGVhZGVyIiwiaXNzIjoiY29mamNjNTRsMzcycHNoMzVpazAiLCJzdWIiOiJsZWFkZXI6ZGV2In0.AYMkP-1o_kqI--cb1BUlSAelEL11pKHo9kL0V5bO1M4")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestSellerBulkPurchaseOrderRepo_PaginateBulkPurchaseOrderMatchingSellers(t *testing.T) {
	var app = initApp("dev")
	params := repo.PaginateBulkPurchaseOrderMatchingSellersParams{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSeller),
		SellerID:      "",
		ProductGroups: []string{"clothing"},
		ProductTypes:  []string{"blouse"},
		FabricTypes:   []string{"sequin"},
	}
	var result = repo.NewSellerBulkPurchaseOrderRepo(app.DB).PaginateBulkPurchaseOrderMatchingSellers(params)
	helper.PrintJSON(result)
}
