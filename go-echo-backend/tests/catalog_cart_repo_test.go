package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/engineeringinflow/inflow-backend/services/consumer"
	"github.com/labstack/echo/v4"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestCatalogCartRepo_BuyerInquiries(t *testing.T) {
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

func TestCatalogCartRepo_BuyerCatalogCarts(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/buyer/catalog_carts", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNqdnNvMDVvb2MyYjhmNDVhMW1nIiwiYXVkIjoiY2xpZW50IiwiaXNzIjoiY2p2dDI5bG9vYzJiOGY0NWExbjAiLCJzdWIiOiJjbGllbnQifQ.9yTBQBC7zB1hlm-Upa3jqO-gidiI2vJW_3CIwIyMQOI")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestCatalogCartRepo_BuyerUpdateCatalogCarts(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = []byte(`{"records":[{"items":[{"id":"cmkapp3b2hja012eoef0","created_at":1705553124,"updated_at":1705553124,"product_id":"cmj4k8bb2hj7fi61dqv0","product":{"id":"cmj4k8bb2hj7fi61dqv0","created_at":1705396769,"updated_at":1705458306,"name":"하늘","slug":"haneul","description":"<p>gfrihgihre</p>","category_id":"cm6i4a3b2hj4koef958g","category":null,"price":"0","product_type":"clothing","bullet_points":0,"special_features":"","material":"","gender":"","style":"","country_code":"US","safety_info":"","currency":"USD","trade_unit":"piece","attachments":[{"content_type":"image/jpeg","file_key":"uploads/media/cl6a9lh343v2qdiigju0_products_cmj4fkjb2hj3ktvri4s0.jpeg","metadata":{"name":"2023.jpg","size":16496}}],"fabric_ids":["cm9ntpbb2hj59girefr0","cmj2pdrb2hj2bhd9qst0","cm9okjbb2hj59giregu0"],"source_product_id":"","source":"inflow","is_trending":false},"user_id":"cl4cmjnagm8dsj056um0","variant_id":"cmjjl0jb2hj3ckuh0h2g","variant":{"id":"cmjjl0jb2hj3ckuh0h2g","created_at":1705396769,"updated_at":1705396769,"title":"red,S","product_id":"cmj4k8bb2hj7fi61dqv0","price":"12","min_order":12,"color":"red","size":"S","is_show":false,"images":[{"content_type":"image/jpeg","file_key":"uploads/media/cl6a9lh343v2qdiigju0_products_cmj4k83b2hj7fi61dqug.jpeg","metadata":{"name":"images.jpeg","size":4593}}]},"size":"S","color":"red","fabric_id":"cm9ntpbb2hj59girefr0","fabric":{"id":"cm9ntpbb2hj59girefr0","created_at":1704165093,"updated_at":1704165093,"reference_id":"","fabric_type":"lolo","description":"987","fabric_id":"876","fabric_weight":987,"moq":987,"colors":null,"attachments":[{"content_type":"image/webp","file_key":"uploads/media/cl658pd7qfpgrlqqj93g_fabric_cm9ntp3b2hj59girefqg","metadata":{"name":"631463a45c56303f2d1afbed_Screenshot 2022-09-04 at 09.36.16-min.webp","size":80210}}],"fabric_costings":[{"from":99,"to":9890,"price":"98","processing_time":"098"}],"slug":"lolo"},"attachments":[{"content_type":"image/jpeg","file_key":"uploads/media/cl6a9lh343v2qdiigju0_products_cmj4k83b2hj7fi61dqug.jpeg","metadata":{"name":"images.jpeg","size":4593}}],"unit_price":"12","quantity":121,"total_price":"144"},{"id":"cmkapp3b2hja012eoefg","created_at":1705553124,"updated_at":1705553124,"product_id":"cmj4k8bb2hj7fi61dqv0","product":{"id":"cmj4k8bb2hj7fi61dqv0","created_at":1705396769,"updated_at":1705458306,"name":"하늘","slug":"haneul","description":"<p>gfrihgihre</p>","category_id":"cm6i4a3b2hj4koef958g","category":null,"price":"0","product_type":"clothing","bullet_points":0,"special_features":"","material":"","gender":"","style":"","country_code":"US","safety_info":"","currency":"USD","trade_unit":"piece","attachments":[{"content_type":"image/jpeg","file_key":"uploads/media/cl6a9lh343v2qdiigju0_products_cmj4fkjb2hj3ktvri4s0.jpeg","metadata":{"name":"2023.jpg","size":16496}}],"fabric_ids":["cm9ntpbb2hj59girefr0","cmj2pdrb2hj2bhd9qst0","cm9okjbb2hj59giregu0"],"source_product_id":"","source":"inflow","is_trending":false},"user_id":"cl4cmjnagm8dsj056um0","variant_id":"cmjjl0jb2hj3ckuh0h2g","variant":{"id":"cmjjl0jb2hj3ckuh0h2g","created_at":1705396769,"updated_at":1705396769,"title":"red,S","product_id":"cmj4k8bb2hj7fi61dqv0","price":"12","min_order":12,"color":"red","size":"S","is_show":false,"images":[{"content_type":"image/jpeg","file_key":"uploads/media/cl6a9lh343v2qdiigju0_products_cmj4k83b2hj7fi61dqug.jpeg","metadata":{"name":"images.jpeg","size":4593}}]},"size":"S","color":"red","fabric_id":"cm9okjbb2hj59giregu0","fabric":{"id":"cm9okjbb2hj59giregu0","created_at":1704168013,"updated_at":1704168013,"reference_id":"","fabric_type":"luulong","description":"987","fabric_id":"897","fabric_weight":987,"moq":897,"colors":null,"attachments":[{"content_type":"image/jpeg","file_key":"uploads/media/cl658pd7qfpgrlqqj93g_fabric_cm9okh3b2hj59giregtg.jpeg","metadata":{"name":"_BIS1660.JPG","size":8521468}}],"fabric_costings":[{"from":9,"to":9,"price":"9","processing_time":"0-9"}],"slug":"luulong"},"attachments":[{"content_type":"image/jpeg","file_key":"uploads/media/cl6a9lh343v2qdiigju0_products_cmj4k83b2hj7fi61dqug.jpeg","metadata":{"name":"images.jpeg","size":4593}}],"unit_price":"12","quantity":12,"total_price":"144"}],"cart_id":"cmkapp3b2hja012eoeg0"}]}`)

	var req = httptest.NewRequest(echo.PUT, "/api/v1/buyer/catalog_carts", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNGNtam5hZ204ZHNqMDU2dW0wIiwiYXVkIjoiY2xpZW50IiwiaXNzIjoiY2w0Y21qbmFnbThkc2owNTZ1bWciLCJzdWIiOiJjbGllbnQifQ.YZ_ExUVTluANzLs46fPDfY_a-ESFXXMrjD03jLICEus")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestCatalogCartRepo_BuyerCheckoutCatalogCarts(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()
	consumer.New(app, false)

	var body = []byte(`{"cart_ids":["cmkfmabb2hj8vsnsm3vg"],"payment_type":"card","payment_method_id":"pm_1OEpBYLr6GIPd0Z17nmXIUJV"}`)

	var req = httptest.NewRequest(echo.POST, "/api/v1/buyer/catalog_carts/checkout", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNsNGNtam5hZ204ZHNqMDU2dW0wIiwiYXVkIjoiY2xpZW50IiwiaXNzIjoiY2w0Y21qbmFnbThkc2owNTZ1bWciLCJzdWIiOiJjbGllbnQifQ.YZ_ExUVTluANzLs46fPDfY_a-ESFXXMrjD03jLICEus")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}
