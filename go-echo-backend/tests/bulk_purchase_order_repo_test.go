package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/validation"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/thaitanloi365/go-utils/values"
)

func TestBulkPurchaseOrderRepo_PaginateBulkPurchaseOrderAPI(t *testing.T) {
	var app = initApp("dev")

	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/seller/bulk_purchase_orders?page=1&limit=12", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNtZG04dHJiMmhqYzVnMmd2M2RnIiwiYXVkIjoic2VsbGVyIiwiaXNzIjoiY21kbTh0cmIyaGpjNWcyZ3YzZTAiLCJzdWIiOiJzZWxsZXIifQ.8f_iQqXLfiXCpA_p32-cnk8v6mbVE8Wjm4488v6VEdA")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestBulkPurchaseOrderRepo_PaginateBulkPurchaseOrder(t *testing.T) {
	var app = initApp("dev")

	var resp = repo.NewBulkPurchaseOrderRepo(app.DB).PaginateBulkPurchaseOrder(repo.PaginateBulkPurchaseOrderParams{
		JwtClaimsInfo:    *models.NewJwtClaimsInfo().SetRole(enums.RoleStaff).SetUserID("cg5anr2llkm6ctpvq8k0"),
		IncludeTrackings: true,
		PaginationParams: models.PaginationParams{
			Limit: 20,
		},
	})

	helper.PrintJSON(resp)
}

func TestBulkPurchaseOrderRepo_BulkPurchaseOrderPreviewCheckoutAndUpdateTax(t *testing.T) {
	var app = initApp("prod")

	var updatesInquiry models.Inquiry
	updatesInquiry.TaxPercentage = values.Float64(8)
	app.DB.Model(&models.Inquiry{}).Where("id = ?", "cnsdhpmpdb85psrjrsmg").Updates(&updatesInquiry)

	var updatesPO models.PurchaseOrder
	updatesPO.TaxPercentage = values.Float64(8)
	app.DB.Model(&models.PurchaseOrder{}).Where("id = ?", "cnsdhpmpdb85psrjrsn0").Updates(&updatesPO)

	var bulkPO models.BulkPurchaseOrder
	var err = app.DB.Select("ID", "CommercialInvoice").First(&bulkPO, "id = ?", "co5p9bti7kqtgl3b3ih0").Error
	assert.NoError(t, err)

	bulkPO.CommercialInvoice.TaxPercentage = values.Float64(8)

	var updatesBulk models.BulkPurchaseOrder
	updatesBulk.TaxPercentage = values.Float64(8)
	updatesBulk.CommercialInvoice = bulkPO.CommercialInvoice
	app.DB.Model(&models.BulkPurchaseOrder{}).Where("id = ?", "co5p9bti7kqtgl3b3ih0").Updates(&updatesBulk)

	resp, err := repo.NewBulkPurchaseOrderRepo(app.DB).BulkPurchaseOrderPreviewCheckout(repo.BulkPurchaseOrderPreviewCheckoutParams{
		BulkPurchaseOrderID: "co5p9bti7kqtgl3b3ih0",
		PaymentType:         enums.PaymentTypeCard,
		JwtClaimsInfo:       *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin).SetUserID(""),
		Milestone:           enums.PaymentMilestoneFinalPayment,
		UpdatePricing:       true,
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}
func TestBulkPurchaseOrderRepo_BulkPurchaseOrderPreviewCheckout(t *testing.T) {
	var app = initApp("prod")

	var updatesInquiry models.Inquiry
	updatesInquiry.TaxPercentage = values.Float64(8)
	app.DB.Model(&models.Inquiry{}).Where("id = ?", "cnsdhpmpdb85psrjrsmg").Updates(&updatesInquiry)

	var updatesPO models.PurchaseOrder
	updatesPO.TaxPercentage = values.Float64(8)
	app.DB.Model(&models.PurchaseOrder{}).Where("id = ?", "cnsdhpmpdb85psrjrsn0").Updates(&updatesPO)

	var updatesBulk models.BulkPurchaseOrder
	updatesBulk.TaxPercentage = values.Float64(8)
	app.DB.Model(&models.BulkPurchaseOrder{}).Where("id = ?", "co5p9bti7kqtgl3b3ih0").Updates(&updatesBulk)

	resp, err := repo.NewBulkPurchaseOrderRepo(app.DB).BulkPurchaseOrderPreviewCheckout(repo.BulkPurchaseOrderPreviewCheckoutParams{
		BulkPurchaseOrderID: "co5p9bti7kqtgl3b3ih0",
		PaymentType:         enums.PaymentTypeCard,
		JwtClaimsInfo:       *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin).SetUserID(""),
		Milestone:           enums.PaymentMilestoneFinalPayment,
		UpdatePricing:       true,
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}

func TestBulkPurchaseOrderRepo_SendQuotationToBuyer(t *testing.T) {
	var app = initApp("prod")
	var rawData = []byte(`{"admin_quotations":[{"lead_time":30,"price":6.8,"quantity":1,"type":"bulk","can_delete":false}],"first_payment_percentage":50,"quotation_note":"please help to arrange the 50%payment to proceed the bulk order and we will send the final shipping price once have the final weight of shipment and get the quote for arrange the shipment and you can pay balance amount of bulk and shipping cost before the shipment ready for shipping. ","shipping_fee":"0"}`)

	var params models.SendBulkPurchaseOrderQuotationParams
	var err = json.Unmarshal(rawData, &params)
	assert.NoError(t, err)

	params.BulkPurchaseOrderID = "cldd3jbeh9ihn4f92b7g"
	params.SetUserID("ckq2n5tc0gstp98hesp0")
	params.SetRole(enums.RoleLeader)

	resp, err := repo.NewBulkPurchaseOrderRepo(app.DB).SendQuotationToBuyer(params)
	assert.NoError(t, err)

	helper.PrintJSON(resp)
}

func TestBulkPurchaseOrderRepoUpdateBulkPurchaseOrder(t *testing.T) {
	var app = initApp("prod")
	var rawData = []byte(`{"packing_note":"","shipping_method":"cif","shipping_note":"","additional_requirements":"","packing_attachments":[],"shipping_attachments":[],"attachments":[{"content_type":"image/jpeg","file_key":"uploads/media/anonymous_rfq_attachments_clddl4jeh9ihn4f92b8g.jpeg","metadata":{"name":"IMG_1141.jpeg","size":96768}}],"items":[
		{
			"color_name": "black",
			"size": "S",
			"qty": 50
		},
		{
			"color_name": "black",
			"size": "M",
			"qty": 50
		},
		{
			"color_name": "black",
			"size": "L",
			"qty": 50
		}
	]}`)

	var params models.BulkPurchaseOrderUpdateForm
	var err = json.Unmarshal(rawData, &params)
	assert.NoError(t, err)

	params.BulkPurchaseOrderID = "cldd3jbeh9ihn4f92b7g"
	params.SetUserID("ckq2n5tc0gstp98hesp0")
	params.SetRole(enums.RoleLeader)

	resp, err := repo.NewBulkPurchaseOrderRepo(app.DB).UpdateBulkPurchaseOrder(params)
	assert.NoError(t, err)

	helper.PrintJSON(resp)
}

func TestBulkPurchaseOrderRepo_ExportExcel(t *testing.T) {
	var app = initApp("local")
	_, _ = repo.NewBulkPurchaseOrderRepo(app.DB).ExportExcel(repo.PaginateBulkPurchaseOrderParams{})
}

func TestBulkPurchaseOrderRepo_GetInvoiceSubTotal(t *testing.T) {
	var app = initApp("dev")
	order, err := repo.NewBulkPurchaseOrderRepo(app.DB).GetBulkPurchaseOrder(repo.GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID:        "cml2ocbb2hj2gpgsbmr0",
		IncludeShippingAddress:     true,
		IncludePaymentTransactions: true,
		IncludeUser:                true,
		IncludeAssignee:            true,
		IncludeInvoice:             true,
		IncludeItems:               true,
	})
	assert.NoError(t, err)

	helper.PrintJSON(order)
}

func TestBulkPurchaseOrderRepo_BulkPurchaseOrderConfirmQCReport(t *testing.T) {
	// var app = initApp("dev")
	e := echo.New()

	var rawData = []byte(`
	{
		"shipping_fee": "12",
		"shipping_address": {
		  "id": "0d97e78655d29ffae8d68c6cd737a274",
		  "created_at": 1700554247,
		  "updated_at": 1700554247,
		  "user_id": "",
		  "name": "Thuy Nguyen",
		  "phone_number": "+84 947 420 900",
		  "address_type": "primary",
		  "coordinate_id": "9ebbf96ba1917390446e78b650317a51",
		  "coordinate": {
			"id": "9ebbf96ba1917390446e78b650317a51",
			"address_number": "232",
			"formatted_address": "Nod nod, 232 Đường Võ Văn Tần, P.5, Q.3, Hồ Chí Minh, Vietnam",
			"street": "Đường Võ Văn Tần",
			"level_1": "Hồ Chí Minh",
			"level_2": "Q.3",
			"level_3": "P.5",
			"postal_code": "",
			"country_code": "VN",
			"timezone_name": "Asia/Ho_Chi_Minh",
			"timezone_offset": 25200
		  }
		},
		"order_closing_doc": {
		  "order_closing_date": 1701836641,
		  "purchase_order_number": "BPO-ZYIL-69174",
		  "purchase_order_date": 1701836503,
		  "order_items": [
			{
			  "id": "5870bd91-4bee-48db-bc95-b66f37811cbb",
			  "color": "red",
			  "size": {
				"XS": {
				  "po_quantity": 10,
				  "actual": 8
				},
				"S": {
				  "po_quantity": 12,
				  "actual": 12
				}
			  },
			  "unit_price": {
				"po_quantity": 11111,
				"actual": 11111
			  },
			  "total_quantity": {
				"po_quantity": 22,
				"actual": 20
			  },
			  "total_amount": {
				"po_quantity": 244442,
				"actual": 222220
			  }
			},
			{
			  "id": "1",
			  "color": "blue",
			  "size": {
				"XS": {
				  "po_quantity": 222,
				  "actual": 222
				},
				"S": {
				  "po_quantity": 11,
				  "actual": 11
				}
			  },
			  "unit_price": {
				"po_quantity": 11111,
				"actual": 11111
			  },
			  "total_quantity": {
				"po_quantity": 233,
				"actual": 233
			  },
			  "total_amount": {
				"po_quantity": 2588863,
				"actual": 2588863
			  }
			},
			{
			  "id": "47787c3b-6153-4255-9136-7ba58e90edc5",
			  "color": "green",
			  "size": {
				"XS": {
				  "po_quantity": 21,
				  "actual": 8
				},
				"S": {
				  "po_quantity": null,
				  "actual": null
				}
			  },
			  "unit_price": {
				"po_quantity": 11111,
				"actual": 11111
			  },
			  "total_quantity": {
				"po_quantity": 21,
				"actual": 8
			  },
			  "total_amount": {
				"po_quantity": 233331,
				"actual": 88888
			  }
			}
		  ]
		},
		"commercial_invoice": {
		  "invoice_number": 22,
		  "vendor": {
			"name": "INFLOW COMPANY LIMITED",
			"address": "48 HUYNH MAN DAT, WARD 19, BINH THANH DISTRICT, HO CHI MINH CITY, VIETNAM updated",
			"phone_number": "+84 (876) 543 2198",
			"contact_name": "Khanh Le",
			"email": "khanhle@joininflow.io"
		  },
		  "consignee": {
			"name": "Thuy Nguyen",
			"address": "Nod nod, 232 Đường Võ Văn Tần, P.5, Q.3, Hồ Chí Minh, Vietnam updated",
			"phone_number": "+84 947 420 900"
		  },
		  "currency": "VND",
		  "issued_date": 1701836.784,
		  "shipping_fee": 12,
		  "tax_percentage": 6,
		  "items": [
			{
			  "id": "0",
			  "color": "red",
			  "size": {
				"XS": 8,
				"S": 12
			  },
			  "unit_price": 11111,
			  "total_quantity": 20,
			  "total_amount": 222220
			},
			{
			  "id": "1",
			  "color": "blue",
			  "size": {
				"XS": 222,
				"S": 11
			  },
			  "unit_price": 11111,
			  "total_quantity": 233,
			  "total_amount": 2588863
			},
			{
			  "id": "2",
			  "color": "green",
			  "size": {
				"XS": 8,
				"S": null
			  },
			  "unit_price": 11111,
			  "total_quantity": 8,
			  "total_amount": 88888
			}
		  ],
		  "country_code": "VN",
		  "status": "paid",
		  "invoice_type": "bulk_po_final_payment"
		}
	  }
	`)

	req := httptest.NewRequest("PUT", "/api/v1/admin/bulk_purchase_orders/clmpg8rb2hj3luoqdj6g/confirm_qc_report", bytes.NewBuffer(rawData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var params repo.BulkPurchaseOrderConfirmQcReportParams
	var err = e.Binder.Bind(&params, c)
	assert.NoError(t, err)
	helper.PrintJSON(params)

	params.BulkPurchaseOrderID = "clmpg8rb2hj3luoqdj6g"
	err = validation.RegisterValidation().Validate(&params)
	assert.NoError(t, err)

	// return
	// order, err := repo.NewBulkPurchaseOrderRepo(app.DB).BulkPurchaseOrderConfirmQCReport(params)
	// assert.NoError(t, err)

	// helper.PrintJSON(order)
}

func TestBulkPurchaseOrderRepo_AdminResetBulkPurchaseOrder(t *testing.T) {
	var app = initApp("prod")
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(app.DB).ResetBulkPurchaseOrder(repo.ResetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: "clr8aubb2hjeq366f5a0",
	})
	assert.NoError(t, err)
	helper.PrintJSON(bulkPO)
}

func TestBulkPurchaseOrderRepo_UploadExcel(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = new(bytes.Buffer)
	var writer = multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "Order Summary_Denim March_Inflow.xlsx")
	assert.NoError(t, err)

	content, err := ioutil.ReadFile("./Order Summary_Denim March_Inflow.xlsx")
	assert.NoError(t, err)

	_, err = io.Copy(part, bytes.NewBuffer(content))
	assert.NoError(t, err)

	writer.Close()

	var req = httptest.NewRequest(echo.POST, "/api/v1/admin/users/cjvso05ooc2b8f45a1mg/upload_bulks", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Timzeon", "Asia/Saigon")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwidHoiOiJBc2lhL1NhaWdvbiIsImNpZCI6IiIsImN0eXBlIjoiYnV5ZXIiLCJhdWQiOiJzdXBlcl9hZG1pbiIsImlzcyI6ImNsNHI4am9sdXJ1dG5xNmkybGZnIiwic3ViIjoic3VwZXJfYWRtaW4ifQ.S_OQiJmaQ3xc_Hmq6GaTGuq34YNr_vW6-CuohlYktdo")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestBulkPurchaseOrderRepo_UploadBOM(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	data, err := repo.NewBulkPurchaseOrderRepo(app.DB).UploadBOM(repo.UploadBOMParams{
		BulkPurchaseOrderID: "cmahqm3b2hj5g2d9rv6g",
		FileKey:             "uploads/media/cg5anr2llkm6ctpvq8k0_bom_cmf0tjbb2hjdff5iqab0.xlsx",
	})

	assert.NoError(t, err)

	helper.PrintJSON(data)
}
