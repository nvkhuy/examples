package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestBulkPurchaseOrderTrackingRepo_PaginateLogsAPI(t *testing.T) {
	var app = initApp("dev")

	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/seller/bulk_purchase_orders/cobon242vss1o35c3c90/logs", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNtZG04dHJiMmhqYzVnMmd2M2RnIiwiYXVkIjoic2VsbGVyIiwiaXNzIjoiY21kbTh0cmIyaGpjNWcyZ3YzZTAiLCJzdWIiOiJzZWxsZXIifQ.8f_iQqXLfiXCpA_p32-cnk8v6mbVE8Wjm4488v6VEdA")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}
