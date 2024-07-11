package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func TestInquirySellerRepo_ApproveQuotation(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = bytes.NewBuffer([]byte(`{"transaction_ref_id":"sss","transaction_attachment":{"content_type":"image/jpeg","file_key":"uploads/media/cl6a9lh343v2qdiigju0_bank_transfer_couopjjdat4b8h9dhsa0.jpeg","path":"2023.jpg","metadata":{"name":"2023.jpg","size":16496}}}`))
	var req = httptest.NewRequest(echo.POST, "/api/v1/admin/inquiry_sellers/cp0pckqhrhpaiakbcqc0/approve", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNmE5bGgzNDN2MnFkaWlnanUwIiwiYXVkIjoic3RhZmYiLCJpc3MiOiJjb2ZwZXVubDNnYjB1Mmw2N2k2ZyIsInN1YiI6InN0YWZmOmRldiJ9.1swoL-c3-8q1MbWetv9BMymY3IRWHLKNsTSYm_IVYeQ")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestInquirySellerRepo_PayoutToSeller(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = bytes.NewBuffer([]byte(`{"transaction_ref_id":"sss","transaction_attachment":{"content_type":"image/jpeg","file_key":"uploads/media/cl6a9lh343v2qdiigju0_bank_transfer_couopjjdat4b8h9dhsa0.jpeg","path":"2023.jpg","metadata":{"name":"2023.jpg","size":16496}}}`))
	var req = httptest.NewRequest(echo.POST, "/api/v1/admin/purchase_orders/coitn22hr8g5bhlc4p8g/payout_to_seller", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNmE5bGgzNDN2MnFkaWlnanUwIiwiYXVkIjoic3RhZmYiLCJpc3MiOiJjb2ZwZXVubDNnYjB1Mmw2N2k2ZyIsInN1YiI6InN0YWZmOmRldiJ9.1swoL-c3-8q1MbWetv9BMymY3IRWHLKNsTSYm_IVYeQ")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestInquirySellerRepo_StatusCount(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/admin/inquiries/cojruuuqps4octv673c0/seller_requests/status_count", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNmE5bGgzNDN2MnFkaWlnanUwIiwiYXVkIjoic3RhZmYiLCJpc3MiOiJjb2ZwZXVubDNnYjB1Mmw2N2k2ZyIsInN1YiI6InN0YWZmOmRldiJ9.1swoL-c3-8q1MbWetv9BMymY3IRWHLKNsTSYm_IVYeQ")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestInquirySellerRepo_Seller(t *testing.T) {
	var app = initApp("dev")
	resp := repo.NewInquirySellerRepo(app.DB).InquirySellerAllocationSearchSeller(repo.InquirySellerAllocationSearchSellerParams{
		InquiryID:        "cle8h4tnhvrn09h2fds0",
		PaginationParams: models.PaginationParams{
			// Keyword: "do",
		},
	})

	helper.PrintJSON(resp)
}

func TestInquirySellerRepo_PaginateMatchingSellers(t *testing.T) {
	var app = initApp("dev")

	var params repo.PaginateMatchingSellersParams
	params.ProductGroups = []string{enums.OBProductGroupClothing.String()}
	params.ProductTypes = []string{enums.OBFactoryProductTypeBlouse.String(), enums.OBFactoryProductTypeBodysuit.String(), enums.OBFactoryProductTypeCamisole.String()}
	params.FabricTypes = []string{enums.OBFabricTypeSatin.String()}

	var result = repo.NewInquirySellerRepo(app.DB).PaginateMatchingSellers(params)

	helper.PrintJSON(result)
}

func TestInquirySellerRepo_GetInquirySellerRequestByID(t *testing.T) {
	var app = initApp("local")

	claims := *models.NewJwtClaimsInfo().SetUserID("cir3f55ocd6fbjvrsv10")
	id := "clbj3piebq59evpn5a6g"

	result, err := repo.NewInquirySellerRepo(app.DB).GetInquirySellerRequestByID(id, queryfunc.InquirySellerRequestBuilderOptions{
		IncludeInquiry: true,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: claims.GetRole(),
		},
	})

	assert.NoError(t, err)
	helper.PrintJSON(result)
}
