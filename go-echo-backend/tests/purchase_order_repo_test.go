package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
)

func TestPurchaseOrderRepo_CreatePurchaseOrder(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = []byte(`{"is_paid":false,"quotations":[{"lead_time":10,"price":10,"quantity":100,"type":"bulk","can_delete":false}],"items":[{"unit_price":2,"qty":1,"size":"214235345","color_name":"red","note_to_supplier":"test","id":"8be074a0-754b-496c-b54a-1b9c747a1431"},{"unit_price":5,"qty":1,"size":"12312321","color_name":"green","note_to_supplier":""}],"shipping_address":{"name":"Luu Long","phone_number":"+84 346 374 333","coordinate":{"lat":38.949362,"lng":-76.489805,"address_number":"914","formatted_address":"One, 914 Bay Ridge Rd, Ste 212, Annapolis, MD, 21403, USA","street":"Bay Ridge Rd","level_1":"Annapolis","level_2":"Anne Arundel County","level_3":"Maryland","postal_code":"21403","country_code":"US","label":"One, 914 Bay Ridge Rd, Ste 212, Annapolis, MD, 21403, USA","value":"One, 914 Bay Ridge Rd, Ste 212, Annapolis, MD, 21403, USA"}},"tax_percentage":10,"product_weight":1,"shipping_fee":20,"currency":"USD","size_chart":"New","user_id":"cmdqq4rb2hjdfi4g5a8g","transaction_ref_id":""}`)

	var req = httptest.NewRequest(echo.POST, "/api/v1/admin/purchase_orders", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwiYXVkIjoic3VwZXJfYWRtaW4iLCJpc3MiOiJjbWlmM2JqYjJoamNkcW9xNmlkMCIsInN1YiI6InN1cGVyX2FkbWluIn0.7d8YdxKwyJ5dfMwHrnPdDeP76FMjQEiQf3VQFD_KFnE")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestPurchaseOrderRepo_AdminPaginatePurchaseOrders(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/admin/purchase_orders?page=1&limit=12", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwiYXVkIjoic3VwZXJfYWRtaW4iLCJpc3MiOiJjbWlmM2JqYjJoamNkcW9xNmlkMCIsInN1YiI6InN1cGVyX2FkbWluIn0.7d8YdxKwyJ5dfMwHrnPdDeP76FMjQEiQf3VQFD_KFnE")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// helper.PrintJSONBytes(rec.Body.Bytes())

	ioutil.WriteFile("test.json", rec.Body.Bytes(), 0664)

}

func TestPurchaseOrderRepo_BuyerPaginatePurchaseOrders(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/buyer/purchase_orders?page=1&limit=12", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNGNtam5hZ204ZHNqMDU2dW0wIiwiYXVkIjoiY2xpZW50IiwiaXNzIjoiY2w0Y21qbmFnbThkc2owNTZ1bWciLCJzdWIiOiJjbGllbnQifQ.YZ_ExUVTluANzLs46fPDfY_a-ESFXXMrjD03jLICEus")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// helper.PrintJSONBytes(rec.Body.Bytes())

	ioutil.WriteFile("test.json", rec.Body.Bytes(), 0664)

}

func TestPurchaseOrderRepo_GetPurchaseOrder(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/buyer/purchase_orders/cm1unabb2hj362lpvts0", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNqdnNvMDVvb2MyYjhmNDVhMW1nIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6ImNsaWVudCIsImlzcyI6ImNqdnQyOWxvb2MyYjhmNDVhMW4wIiwic3ViIjoiY2xpZW50In0.IrsTPVsoXJHKrgrC_V0D5ZLEe6B_JmcIpqtmKl0fa2o")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// helper.PrintJSONBytes(rec.Body.Bytes())

	ioutil.WriteFile("test.json", rec.Body.Bytes(), 0664)

}

func TestPurchaseOrderRepo_PurchaseOrderRepoTests(t *testing.T) {
	var app = initApp("dev")

	var resp = repo.NewPurchaseOrderRepo(app.DB).PaginatePurchaseOrders(repo.PaginatePurchaseOrdersParams{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetUserID("cl4cmjnagm8dsj056um0"),
		PaginationParams: models.PaginationParams{
			Page:    1,
			Limit:   10,
			Keyword: "fr",
		},
		//TeamID:           "cjvso05ooc2b8f45a1mg",
		// Statuses: []enums.PurchaseOrderStatus{
		// 	enums.PurchaseOrderStatusPaid,
		// },
		// InquiryStatuses: []enums.InquiryBuyerStatus{
		// 	enums.InquiryBuyerStatusApproved,
		// },
	})

	helper.PrintJSON(resp)
}

func TestPurchaseOrderRepo_BuyerGivePurchaseOrderFeedback(t *testing.T) {
	var app = initApp("local")

	var err = repo.NewPurchaseOrderRepo(app.DB).BuyerGivePurchaseOrderFeedback(repo.PurchaseOrderFeedbackParams{
		PurchaseOrderID: "cl6u5bcm8dkhnnn6qca0",
		PurchaseOrderFeedback: models.PurchaseOrderFeedback{
			Visual: &models.PurchaseOrderFeedbackDetails{
				MeetExpectation: true,
			},
			Form: &models.PurchaseOrderFeedbackDetails{
				MeetExpectation: true,
			},
			Measurements: &models.PurchaseOrderFeedbackDetails{
				NeedsImprovement: "need to improve 2",
			},
			Workmanship: &models.PurchaseOrderFeedbackDetails{
				NeedsImprovement: "need to improve 3",
			},
			FabricTrimQuality: &models.PurchaseOrderFeedbackDetails{
				NeedsImprovement: "need to improve 4",
			},
			Printing: &models.PurchaseOrderFeedbackDetails{
				NeedsImprovement: "need to improve 5",
			},
			Embroidery: &models.PurchaseOrderFeedbackDetails{
				NeedsImprovement: "need to improve 6",
			},
			WashingDyeing: &models.PurchaseOrderFeedbackDetails{
				MeetExpectation: true,
			},
		},
	})

	assert.NoError(t, err)
}

func TestPurchaseOrderRepo_ExportExcel(t *testing.T) {
	var app = initApp("local")
	_, _ = repo.NewPurchaseOrderRepo(app.DB).ExportExcel(repo.PaginatePurchaseOrdersParams{})
}

func TestPurchaseOrderRepo_SellerPoUploadCommentStatusCount(t *testing.T) {
	var app = initApp("dev")

	var params repo.PoCommentStatusCountParams

	params.PurchaseOrderID = "clit108qier7erkd6ai0"
	params.JwtClaimsInfo = *models.NewJwtClaimsInfo().SetUserID("cl712kdbt6ca83b9kk3g")
	result := repo.NewSellerPurchaseOrderRepo(app.DB).SellerPurchaseOrderUploadCommentStatusCount(params)
	helper.PrintJSON(result)
}

func TestPurchaseOrderRepo_AdminSellerPurchaseOrderPayout(t *testing.T) {
	var app = initApp("dev")

	var params repo.AdminSellerPurchaseOrderPayoutParams

	params.PurchaseOrderID = "cmodnkjb2hje0fl1cnl0"
	params.JwtClaimsInfo = *models.NewJwtClaimsInfo().SetUserID("cl658pd7qfpgrlqqj93g")
	params.TransactionRefID = "34534fg"
	params.TransactionAttachment = &models.Attachment{
		ContentType: "image/jpeg",
		FileKey:     "uploads/media/cl658pd7qfpgrlqqj93g_bank_transfer_cmp404bb2hj8l38ilhj0.jpeg",
		Metadata: map[string]interface{}{
			"name": "_BIS1654.JPG",
			"size": 8594470,
		},
	}
	result, err := repo.NewSellerPurchaseOrderRepo(app.DB).AdminSellerPurchaseOrderPayout(params)

	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestPurchaseOrderRepo_BuyerReject(t *testing.T) {
	var app = initApp("local")

	var params repo.BuyerRejectDesignParams

	params.PurchaseOrderID = "co0fgius2p8inl39q1u0"
	params.JwtClaimsInfo = *models.NewJwtClaimsInfo().SetUserID("cl658pd7qfpgrlqqj93g")
	result, err := repo.NewPurchaseOrderRepo(app.DB).BuyerRejectDesign(params)

	assert.NoError(t, err)
	helper.PrintJSON(result)
}
