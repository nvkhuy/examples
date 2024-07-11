package tests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/stretchr/testify/assert"
)

func TestInquiryRepo_BuyerInquiries(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/buyer/inquiries?page=1&limit=12&keyword=jay", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNtYjFtMHJiMmhqZjU3anJhZGpnIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6ImNsaWVudCIsImlzcyI6ImNtYjFtMHJiMmhqZjU3anJhZGswIiwic3ViIjoiY2xpZW50In0.KLStPVNn7wx86X4iJfdhYbnkUx-MOppWFo23pZ0jkUM")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestInquiryRepo_BuyerGetInquiry(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/buyer/inquiries/coij966mj3ba22ur9an0", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNqdnNvMDVvb2MyYjhmNDVhMW1nIiwiYXVkIjoiY2xpZW50IiwiaXNzIjoiY2p2dDI5bG9vYzJiOGY0NWExbjAiLCJzdWIiOiJjbGllbnQifQ.9yTBQBC7zB1hlm-Upa3jqO-gidiI2vJW_3CIwIyMQOI")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}
func TestInquiryRepo_AdminDeleteInquiry(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.DELETE, "/api/v1/admin/inquiries/cl4s2gccr609p5jsrigg/delete", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwidHoiOiJBc2lhL1NhaWdvbiIsImNpZCI6IiIsImN0eXBlIjoiYnV5ZXIiLCJhdWQiOiJzdXBlcl9hZG1pbiIsImlzcyI6ImNsNHI4am9sdXJ1dG5xNmkybGZnIiwic3ViIjoic3VwZXJfYWRtaW4ifQ.S_OQiJmaQ3xc_Hmq6GaTGuq34YNr_vW6-CuohlYktdo")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestInquiryRepo_BuyerUpdateCartItems(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var data = []byte(`
	{
		"items": [
			{
				"id": "cmacm4rb2hjcrg25opmg",
				"created_at": 1704250131,
				"updated_at": 1704250131,
				"deleted_at": null,
				"inquiry_id": "cl8vnkfpc8cqejad34m0",
				"checkout_session_id": "",
				"sku": "",
				"style": "",
				"color": "",
				"color_name": "Grey",
				"size": "Plus 18",
				"qty": 22,
				"unit_price": "10.95",
				"total_price": "240.9",
				"waiting_for_checkout": false,
				"note_to_supplier": ""
			}
		]
	}
	`)
	var req = httptest.NewRequest(echo.POST, "/api/v1/buyer/inquiries/cl8vnkfpc8cqejad34m0/update_cart_items", bytes.NewBuffer(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNrMG92Y3VucjYyaXVlOWNvM3VnIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6ImNsaWVudCIsImlzcyI6Imdob3N0fGNnNWFucjJsbGttNmN0cHZxOGswIiwic3ViIjoiY2xpZW50In0.epJvaYF4HmluDEmD-wIHLThE25Ca2YD42zP84WhFdDA")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestInquiryRepo_BuyerGetCartItems(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/buyer/inquiry_carts/cart", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNrMG92Y3VucjYyaXVlOWNvM3VnIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6ImNsaWVudCIsImlzcyI6Imdob3N0fGNnNWFucjJsbGttNmN0cHZxOGswIiwic3ViIjoiY2xpZW50In0.epJvaYF4HmluDEmD-wIHLThE25Ca2YD42zP84WhFdDA")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}
func TestInquiryRepo_ExportInquiries(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/admin/inquiries/export?statuses=new&buyer_quotation_statuses=new", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6InN1cGVyX2FkbWluIiwiaXNzIjoiY2x1MHNkcmIyaGo5NzFtdmpjYWciLCJzdWIiOiJzdXBlcl9hZG1pbiJ9.f-Co_EV6cAwmaQWS6zRLOfrfDgRT2R5-avvCSbFELlE")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestInquiryRepo_GetInquiryRemindAdmin(t *testing.T) {
	var app = initApp("prod")
	info := models.NewJwtClaimsInfo()
	info.SetRole(enums.RoleSuperAdmin)
	var result = repo.NewInquiryRepo(app.DB).GetInquiryRemindAdmin()

	helper.PrintJSON(result)
}

func TestInquiryRepo_PaginateInquiry(t *testing.T) {
	var app = initApp("dev")
	info := models.NewJwtClaimsInfo()
	info.SetRole(enums.RoleSuperAdmin)
	var result = repo.NewInquiryRepo(app.DB).PaginateInquiry(repo.PaginateInquiryParams{
		JwtClaimsInfo: *info,
		AssigneeID:    "cl1ledua2jjv1foe37og",
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 12,
		},
	})

	helper.PrintJSON(result)
}

func TestInquiryRepo_InquiryPreviewCheckout(t *testing.T) {
	var app = initApp("dev")

	result, err := repo.NewInquiryRepo(app.DB).InquiryPreviewCheckout(repo.InquiryPreviewCheckoutParams{
		InquiryID:     "clrdhm3b2hjc6991ui2g",
		PaymentType:   "card",
		UserID:        "cjvso05ooc2b8f45a1mg",
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetUserID("cjvso05ooc2b8f45a1mg"),
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestInquiryRepo_GetInquiryByID(t *testing.T) {
	var app = initApp()
	result, err := repo.NewInquiryRepo(app.DB).GetInquiryByID(repo.GetInquiryByIDParams{
		InquiryID: "ck20k3dfpdpa89aef7d0",
		InquiryBuilderOptions: queryfunc.InquiryBuilderOptions{
			IncludePurchaseOrder:   true,
			IncludeShippingAddress: true,
		},
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestInquiryRepo_InquiryRemoveItems(t *testing.T) {
	var app = initApp()
	var err = repo.NewInquiryRepo(app.DB).InquiryRemoveItems(models.InquiryRemoveItemsForm{
		InquiryID: "ckcg4gk74n32532dijf0",
		ItemIDs:   []string{"ckd6asuodjknobt06tcg"},
	})
	assert.NoError(t, err)

}

func TestInquiryRepo_InquiryAssignPIC(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewInquiryRepo(app.DB).InquiryAssignPIC(models.InquiryAssignPICParam{
		InquiryID:   "ciqdf24tq0u8gav5nudg",
		AssigneeIDs: []string{"cg5anr2llkm6ctpvq8k0"},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestInquiryRepo_InquiryMarkAsPaid(t *testing.T) {
	var app = initApp("local")
	_, err := repo.NewInquiryRepo(app.DB).InquiryMarkAsPaid(models.InquiryIDParam{
		InquiryID: "ckrols77k645b59nn0a0",
	})
	assert.NoError(t, err)
}

func TestInquiryRepo_AdminArchiveInquiry(t *testing.T) {
	var app = initApp("local")
	err := repo.NewInquiryRepo(app.DB).AdminArchiveInquiry(repo.AdminUnarchiveInquiryParams{
		InquiryID: "cl6u40sm8dkhnnn6qc90",
	})
	log.Println("done")
	assert.NoError(t, err)
}

func TestInquiryRepo_PaginateInquiryAudits(t *testing.T) {
	var app = initApp("dev")
	var result = repo.NewInquiryRepo(app.DB).PaginateInquiryAudits(repo.PaginateInquiryAuditsParams{
		InquiryID: "cldir6jh335juevc868g",
	})

	helper.PrintJSON(result)
}

func TestInquiryRepo_ExportExcel(t *testing.T) {
	var app = initApp("dev")
	resp, err := repo.NewInquiryRepo(app.DB).ExportExcel(repo.ExportInquiriesParams{})
	if err != nil {
		return
	}
	helper.PrintJSON(resp)
}

func TestInquiryRepo_PaginateCarts(t *testing.T) {
	var app = initApp("dev")
	var resp = repo.NewInquiryRepo(app.DB).PaginateCarts(repo.PaginateCartsParams{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetUserID(""),
	})

	helper.PrintJSON(resp)
}
