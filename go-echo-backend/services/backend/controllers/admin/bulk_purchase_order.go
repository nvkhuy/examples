package controllers

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginateBulkPurchaseOrder
// @Tags Admin-PO
// @Summary PaginateBulkPurchaseOrder
// @Description PaginateBulkPurchaseOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Records{records=[]models.BulkPurchaseOrder}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders [get]
func AdminPaginateBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBulkPurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludeUser = true
	params.IncludeAssignee = true
	params.IncludeTrackings = true
	params.IncludeCollection = true
	var result = repo.NewBulkPurchaseOrderRepo(cc.App.DB).PaginateBulkPurchaseOrder(params)

	return cc.Success(result)
}

// ExportBulkPurchaseOrder
// @Tags Admin-PO
// @Summary PaginateBulkPurchaseOrder
// @Description PaginateBulkPurchaseOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Records{records=[]models.BulkPurchaseOrder}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/export [get]
func AdminExportBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBulkPurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).ExportExcel(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// AdminUpdateBulkPurchaseOrder
// @Tags Admin-PO
// @Summary Get bulk purchase order
// @Description Get bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id} [put]
func AdminUpdateBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.BulkPurchaseOrderUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrder(form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if result.TrackingStatus == enums.BulkPoTrackingStatusWaitingForSubmitOrder {
		tasks.TrackCustomerIOTask{
			UserID: result.UserID,
			Event:  customerio.EventBulkPoSubmitOrder,
			Data:   result.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(result)
}

// AdminApproveBulkPurchaseOrderSubmit
// @Tags Admin-PO
// @Summary Get bulk purchase order
// @Description Get bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/submit [post]
func AdminBulkPurchaseOrderSubmit(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.BulkPurchaseOrderUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).AdminSubmitBulkPurchaseOrder(form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(order)
}

// GetBulkPurchaseOrder
// @Tags Admin-PO
// @Summary Get bulk purchase order
// @Description Get bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id} [get]
func AdminGetBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetBulkPurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.IncludeShippingAddress = true
	params.IncludePaymentTransactions = true
	params.IncludeUser = true
	params.IncludeAssignee = true
	params.IncludeInvoice = true
	params.IncludeItems = true
	params.IncludeSellerQuotation = true
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).GetBulkPurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// AdminSendQuotationToBuyer
// @Tags Admin-PO
// @Summary Get bulk purchase order
// @Description Get bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/send_quotation [post]
func AdminSendQuotationToBuyer(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.SendBulkPurchaseOrderQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).SendQuotationToBuyer(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if params.FirstPaymentPercentage > 0 {
		tasks.TrackCustomerIOTask{
			UserID: order.UserID,
			Event:  customerio.EventBulkPoSubmitQuotation,
			Data:   order.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())

		tasks.CreateUserNotificationTask{
			UserID:           order.UserID,
			Message:          fmt.Sprintf("New quotation was created for %s", order.ReferenceID),
			NotificationType: enums.UserNotificationTypeBulkPoSubmitQuotation,
			Metadata: &models.UserNotificationMetadata{
				AdminID:                      claims.GetUserID(),
				BulkPurchaseOrderID:          order.ID,
				BulkPurchaseOrderReferenceID: order.ReferenceID,
				InquiryID:                    order.Inquiry.ID,
				InquiryReferenceID:           order.Inquiry.ReferenceID,
			},
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Sent")
}

// BulkPurchaseOrderCreateQcReport
// @Tags Admin-PO
// @Summary BulkPurchaseOrderCreateQcReport
// @Description BulkPurchaseOrderCreateQcReport
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/qc_report [post]
func AdminBulkPurchaseOrderCreateQcReport(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminCreateQcReportParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderCreateQcReport(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if params.ApproveQCAt != nil {
		_, _ = tasks.BulkPurchaseQCApproveTask{
			JwtClaimsInfo:       claims,
			BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		}.DispatchAt(time.Unix(*params.ApproveQCAt, 0), asynq.MaxRetry(0))
	}

	tasks.TrackCustomerIOTask{
		UserID: order.UserID,
		Event:  customerio.EventBulkPoCreateQcReport,
		Data:   order.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	tasks.CreateUserNotificationTask{
		UserID:           order.UserID,
		Message:          fmt.Sprintf("New QC report was created for %s", order.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoCreateQcReport,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          order.ID,
			BulkPurchaseOrderReferenceID: order.ReferenceID,
			InquiryID:                    order.Inquiry.ID,
			InquiryReferenceID:           order.Inquiry.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	return cc.Success(order)
}

// PaginateBulkPurchaseOrderTracking
// @Tags Admin-PO
// @Summary PaginateBulkPurchaseOrderTracking
// @Description PaginateBulkPurchaseOrderTracking
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/logs [get]
func PaginateBulkPurchaseOrderTracking(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginateBulkPurchaseOrderTrackingParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.IncludeTrackings = true
	var result = repo.NewBulkPurchaseOrderTrackingRepo(cc.App.DB).PaginateBulkPurchaseOrderTrackings(params)

	return cc.Success(result)
}

// BulkPurchaseOrderUpdateRawMaterial
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/update_raw_material [put]
func AdminBulkPurchaseOrderUpdateRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.BulkPurchaseOrderUpdateRawMaterialParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdateRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if params.ApproveRawMaterialAt != nil {
		_, _ = tasks.BulkPurchaseOrderRawMaterialApproveTask{
			JwtClaimsInfo:       claims,
			BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		}.DispatchAt(time.Unix(*params.ApproveRawMaterialAt, 0), asynq.MaxRetry(0))
	}

	_, _ = tasks.CreateUserNotificationTask{
		UserID:           bulkPO.UserID,
		Message:          fmt.Sprintf("Raw material was updated for %s", bulkPO.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoUpdateRawMaterial,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          bulkPO.ID,
			BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			InquiryID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ID
				}
				return ""
			}(),
			InquiryReferenceID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ReferenceID
				}
				return ""
			}(),
		},
	}.Dispatch(c.Request().Context())

	return cc.Success(bulkPO)
}

// BulkPurchaseOrderUpdateProduction
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/update_production [put]
func AdminBulkPurchaseOrderUpdateProduction(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.BulkPurchaseOrderUpdateProductionParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	purchaseOrder, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdateProduction(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(purchaseOrder)
}

// BulkPurchaseOrderMarkProduction
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/mark_production [put]
func AdminBulkPurchaseOrderMarkProduction(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderUpdateTrackingStatusParams
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.TrackingAction = enums.BulkPoTrackingActionMarkProduction
	params.TrackingStatus = enums.BulkPoTrackingStatusProduction
	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdateTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoMarkProduction,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	tasks.CreateUserNotificationTask{
		UserID:           bulkPO.UserID,
		Message:          fmt.Sprintf("Bulk purchase order %s is on Production", bulkPO.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoMarkProduction,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          bulkPO.ID,
			BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			InquiryID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ID
				}
				return ""
			}(),
			InquiryReferenceID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ReferenceID
				}
				return ""
			}(),
		},
	}.Dispatch(c.Request().Context())

	return cc.Success(bulkPO)
}

// BulkPurchaseOrderMarkPps
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/mark_pps [put]
func AdminBulkPurchaseOrderMarkRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderUpdateTrackingStatusParams
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.TrackingAction = enums.BulkPoTrackingActionMarkRawMaterial
	params.TrackingStatus = enums.BulkPoTrackingStatusRawMaterial
	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdateTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoMarkRawMaterial,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	tasks.CreateUserNotificationTask{
		UserID:           bulkPO.UserID,
		Message:          fmt.Sprintf("Bulk purchase order %s is on Raw Materials", bulkPO.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoMarkRawMaterial,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          bulkPO.ID,
			BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			InquiryID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ID
				}
				return ""
			}(),
			InquiryReferenceID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ReferenceID
				}
				return ""
			}(),
		},
	}.Dispatch(c.Request().Context())

	return cc.Success(bulkPO)
}

// BulkPurchaseOrderMarkPps
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/mark_pps [put]
func AdminBulkPurchaseOrderMarkPps(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderUpdateTrackingStatusParams
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.TrackingAction = enums.BulkPoTrackingActionMarkPps
	params.TrackingStatus = enums.BulkPoTrackingStatusPps
	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdateTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoMarkPps,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	tasks.CreateUserNotificationTask{
		UserID:           bulkPO.UserID,
		Message:          fmt.Sprintf("Bulk purchase order %s is on PPS", bulkPO.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoMarkPps,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          bulkPO.ID,
			BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			InquiryID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ID
				}
				return ""
			}(),
			InquiryReferenceID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ReferenceID
				}
				return ""
			}(),
		},
	}.Dispatch(c.Request().Context())

	return cc.Success(bulkPO)
}

// BulkPurchaseOrderUpdatePps
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/update_pps [put]
func AdminBulkPurchaseOrderUpdatePps(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.BulkPurchaseOrderUpdatePpsParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdatePps(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateUserNotificationTask{
		UserID:           bulkPO.UserID,
		Message:          fmt.Sprintf("Pre-Production info was updated for %s", bulkPO.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoUpdatePps,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          bulkPO.ID,
			BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			InquiryID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ID
				}
				return ""
			}(),
			InquiryReferenceID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ReferenceID
				}
				return ""
			}(),
		},
	}.Dispatch(c.Request().Context())

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoUpdatePps,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())
	return cc.Success(bulkPO)
}

// BulkPurchaseOrderMarkQc
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/mark_qc [put]
func AdminBulkPurchaseOrderMarkQc(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderUpdateTrackingStatusParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.TrackingAction = enums.BulkPoTrackingActionMarkQc
	params.TrackingStatus = enums.BulkPoTrackingStatusQc
	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdateTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoMarkQc,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	tasks.CreateUserNotificationTask{
		UserID:           bulkPO.UserID,
		Message:          fmt.Sprintf("Bulk purchase order %s is on QC", bulkPO.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoMarkQc,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          bulkPO.ID,
			BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			InquiryID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ID
				}
				return ""
			}(),
			InquiryReferenceID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ReferenceID
				}
				return ""
			}(),
		},
	}.Dispatch(c.Request().Context())

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoMarkQc,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success(bulkPO)
}

// BulkPurchaseOrderConfirmQcReport
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/confirm_qc_report [put]
func AdminBulkPurchaseOrderConfirmQcReport(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.BulkPurchaseOrderConfirmQcReportParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderConfirmQCReport(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoConfirmQCReport,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	tasks.CreateBulkPoAttachmentPDFsTask{
		BulkPurchaseOrderID: bulkPO.ID,
	}.Dispatch(c.Request().Context())

	return cc.Success(bulkPO)
}

// AdminBulkPurchaseOrderMarkDelivering
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/mark_delivering [put]
func AdminBulkPurchaseOrderMarkDelivering(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderMarkDeliveringParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).MarkDelivering(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateUserNotificationTask{
		UserID:           bulkPO.UserID,
		Message:          fmt.Sprintf("Bulk purchase order %s is on delivering", bulkPO.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoMarkDelivering,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          bulkPO.ID,
			BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			InquiryID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ID
				}
				return ""
			}(),
			InquiryReferenceID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ReferenceID
				}
				return ""
			}(),
		},
	}.Dispatch(c.Request().Context())

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoMarkDelivering,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success(bulkPO)
}

// AdminBulkPurchaseOrderMarkDelivering
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/mark_delivered [put]
func AdminBulkPurchaseOrderMarkDelivered(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderMarkDeliveredParams
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).AdminMarkDelivered(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// AdminBulkPurchaseOrderMarkFirstPayment
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/mark_first_payment [put]
func AdminBulkPurchaseOrderMarkFirstPayment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderMarkFirstPaymentParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderMarkFirstPayment(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.BulkPurchaseOrderBankTransferConfirmedTask{
		ApprovedByUserID:    claims.GetUserID(),
		BulkPurchaseOrderID: bulkPO.ID,
		Milestone:           enums.PaymentMilestoneFirstPayment,
	}.Dispatch(cc.Request().Context())

	return cc.Success(bulkPO)
}

// AdminBulkPurchaseOrderMarkFinalPayment
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/mark_final_payment [put]
func AdminBulkPurchaseOrderMarkFinalPayment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderMarkFinalPaymentParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderMarkFinalPayment(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.BulkPurchaseOrderBankTransferConfirmedTask{
		ApprovedByUserID:    claims.GetUserID(),
		BulkPurchaseOrderID: bulkPO.ID,
		Milestone:           enums.PaymentMilestoneFinalPayment,
	}.Dispatch(cc.Request().Context())

	return cc.Success(bulkPO)
}

// AdminBulkPurchaseOrderMarkAsPaid
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/mark_as_paid [put]
func AdminBulkPurchaseOrderMarkAsPaid(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.BulkPurchaseOrderMarkAsPaidParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	_, err = repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderMarkAsPaid(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success("Marked")
}

// AdminBulkPurchaseApproveQc
// @Tags Admin-PO
// @Summary approve qc
// @Description approve qc
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/approve_qc [put]
func AdminBulkPurchaseApproveQc(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderUpdateTrackingStatusParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.TrackingAction = enums.BulkPoTrackingActionBuyerApproveQc
	params.TrackingStatus = enums.BulkPoTrackingStatusSubmit
	params.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdateTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(order)
}

// AdminBulkPurchaseOrderAssignPIC
// @Tags Admin-Inquiry
// @Summary Assign PIC
// @Description Assign PIC
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/assign_pic [put]
func AdminBulkPurchaseOrderAssignPIC(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.BulkPurchaseOrderAssignPICParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderAssignPIC(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, userID := range result.AssigneeIDs {
		tasks.AssignBulkPurchaseOrderPICTask{
			AssignerID:          claims.GetUserID(),
			AssigneeID:          userID,
			BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		}.Dispatch(c.Request().Context())
	}
	return cc.Success(result)
}

// AdminBulkPurchaseOrderGetSamplePO
// @Tags Admin-PO
// @Summary BulkPurchaseOrderGetSamplePO
// @Description BulkPurchaseOrderGetSamplePO
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/sample_po [get]
func AdminBulkPurchaseOrderGetSamplePO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderGetSamplePOParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderGetSamplePO(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminBulkPurchaseBuyerApproveRawMaterial
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/approve_raw_material [post]
func AdminBulkPurchaseApproveRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseBuyerApproveRawMaterialParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderBuyerApproveRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(order)
}

// AdminCreateBulkPurchaseInvoice
// @Tags Admin-PO
// @Summary create purchase order invoice
// @Description create purchase order invoice
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/invoice [post]
func AdminCreateBulkPurchaseInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateBulkPurchaseInvoiceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).CreateBulkPurchaseInvoice(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(order)
}

// AdminBulkPurchaseOrderPreviewCheckout
// @Tags Admin-PO
// @Summary Preview checkout bulk purchase order
// @Description Preview checkout bulk purchase order
// @Accept  json
// @Produce  json
// @Param data body repo.BulkPurchaseOrderPreviewCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/preview_checkout [post]
func AdminBulkPurchaseOrderPreviewCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderPreviewCheckoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.UpdatePricing = true
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderPreviewCheckout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminResetBulkPurchaseOrder
// @Tags Admin-PO
// @Summary Reset bulk purchase order
// @Description Reset bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/reset [put]
func AdminResetBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.ResetBulkPurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).ResetBulkPurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// AdminBulkPurchaseOrderAddNote
// @Tags Admin-BulkPurchaseOrder
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/notes [post]
func AdminBulkPurchaseOrderAddNote(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypeBulkPOInternalNotes
	params.TargetID = cc.GetPathParamString("bulk_purchase_order_id")
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.NewBulkPONotesTask{
		UserID:              claims.GetUserID(),
		BulkPurchaseOrderID: params.TargetID,
		MentionUserIDs:      params.MentionUserIDs,
		Message:             params.Message,
		Attachments:         params.Attachments,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}

// PaginateBulkPurchaseOrderNotes
// @Tags Admin-BulkPurchaseOrder
// @Summary Get design comments
// @Description Get design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/notes [get]
func PaginateBulkPurchaseOrderNotes(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCommentsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.TargetID = cc.GetPathParamString("bulk_purchase_order_id")
	params.TargetType = enums.CommentTargetTypePOInternalNotes
	params.OrderByQuery = "c.created_at DESC"

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// PurchaseOrderNoteMarkSeen
// @Tags Marketplace-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/notes/mark_seen [put]
func AdminBulkPurchaseOrderNoteMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims
	params.TargetID = cc.GetPathParamString("bulk_purchase_order_id")
	params.TargetType = enums.CommentTargetTypeBulkPOInternalNotes
	err = repo.NewCommentRepo(cc.App.DB).MarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// BulkPurchaseOrderNoteUnreadCount
// @Tags Marketplace-PO
// @Summary Design comment status count
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/notes/unread_count [get]
func AdminBulkPurchaseOrderNoteUnreadCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetUnreadCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims
	params.TargetID = cc.GetPathParamString("bulk_purchase_order_id")
	params.TargetType = enums.CommentTargetTypeBulkPOInternalNotes

	var result = repo.NewCommentRepo(cc.App.DB).GetUnreadCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminBulkPurchaseOrderCommentDelete
// @Tags Admin-BulkPurchaseOrder
// @Summary delete inquiry comment
// @Description delete inquiry comment
// @Accept  json
// @Produce  json
// @Param data body models.ContentCommentCreateForm true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/notes/{comment_id} [delete]
func AdminBulkPurchaseOrderCommentDelete(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteCommentParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewCommentRepo(cc.App.DB).DeleteComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}

// AdminBulkPurchaseOrderCommentDelete
// @Tags Admin-BulkPurchaseOrder
// @Summary delete inquiry comment
// @Description delete inquiry comment
// @Accept  json
// @Produce  json
// @Param data body models.ContentCommentCreateForm true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/deposit [post]
func AdminBulkPurchaseOrderDeposit(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateDepositParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).CreateDeposit(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminBulkPurchaseOrderStageCommentsCreate
// @Tags Admin-BulkPurchaseOrder
// @Summary create stage comment
// @Description create stage comment
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.BulkPurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{purchase_order_id}/stage_comments [post]
func AdminBulkPurchaseOrderStageCommentsCreate(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.BulkPurchaseOrderStageCommentsParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).StageCommentsCreate(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventBulkPoNewComment,
		Data: purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
			"comment": params.Comment,
			"comment_attachments": func() models.Attachments {
				return params.Attachments.GenerateFileURL()
			}(),
		}),
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// AdminBulkPurchaseOrderUploadBOM
// @Tags Admin-BulkPurchaseOrder
// @Summary Upload BOM
// @Description Upload BOM
// @Accept  json
// @Produce  json
// @Param data body models.ContentCommentCreateForm true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/upload_bom [post]
func AdminBulkPurchaseOrderUploadBOM(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UploadBOMParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	boms, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).UploadBOM(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(boms)
}

// AdminBulkPurchaseOrderPaginateBOM
// @Tags Admin-BulkPurchaseOrder
// @Summary Paginate BOM
// @Description Paginate BOM
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/{bulk_purchase_order_id}/upload_bom [get]
func AdminBulkPurchaseOrderPaginateBOM(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBOMParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewBulkPurchaseOrderRepo(cc.App.DB).PaginateBOM(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminCreateMultipleBulks
// @Tags Admin-Bulk
// @Summary Admin Create multiple bulks
// @Description Admin Create multiple bulks
// @Accept  json
// @Produce  json
// @Param data body []models.CreateMultipleBulkPurchaseOrdersRequest true
// @Success 200 {object} []models.BulkPurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/admin/bulk_purchase_order/create_multiple [post]
func AdminCreateMultipleBulks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.CreateMultipleBulkPurchaseOrdersRequest

	err := cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	bulks, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).CreateMultipleBulkPurchaseOrders(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, bulk := range bulks {
		_, _ = tasks.CreateChatRoomTask{
			UserID:              claims.GetUserID(),
			Role:                claims.GetRole(),
			BulkPurchaseOrderID: bulk.ID,
			BuyerID:             bulk.UserID,
		}.Dispatch(c.Request().Context())
		if bulk.PurchaseOrder.ID != "" {
			_, _ = tasks.CreateChatRoomTask{
				UserID:          claims.GetUserID(),
				Role:            claims.GetRole(),
				PurchaseOrderID: bulk.PurchaseOrder.ID,
				BuyerID:         bulk.UserID,
			}.Dispatch(c.Request().Context())
		}

		_, _ = tasks.UpdateUserProductClassesTask{
			UserID:          claims.GetUserID(),
			PurchaseOrderID: bulk.ID,
		}.Dispatch(c.Request().Context())

		_, _ = tasks.UpdateUserProductClassesTask{
			UserID:          claims.GetUserID(),
			PurchaseOrderID: bulk.ID,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(bulks)
}

// AdminSubmitMultipleBulkQuotations
// @Tags Admin-PO
// @Summary Admin submit multiple bulk quotations
// @Description Admin submit multiple bulk quotations
// @Accept  json
// @Produce  json
// @Param data body models.SubmitMultipleBulkQuotationsRequest true
// @Success 200 {object} models.BulkPurchaseOrders
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/submit_multiple_quotations [post]
func AdminSubmitMultipleBulkQuotations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.SubmitMultipleBulkQuotationsRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulks, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).SubmitMultipleBulkQuotations(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, bulk := range bulks {
		if *bulk.FirstPaymentPercentage > 0 {
			tasks.TrackCustomerIOTask{
				UserID: bulk.UserID,
				Event:  customerio.EventBulkPoSubmitQuotation,
				Data:   bulk.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())

			tasks.CreateUserNotificationTask{
				UserID:           bulk.UserID,
				Message:          fmt.Sprintf("New quotation was created for %s", bulk.ReferenceID),
				NotificationType: enums.UserNotificationTypeBulkPoSubmitQuotation,
				Metadata: &models.UserNotificationMetadata{
					AdminID:                      claims.GetUserID(),
					BulkPurchaseOrderID:          bulk.ID,
					BulkPurchaseOrderReferenceID: bulk.ReferenceID,
				},
			}.Dispatch(c.Request().Context())
		}
	}

	return cc.Success("Sent")
}

// UploadBulkPurchaseOrder
// @Tags Admin-User
// @Summary delete inquiry comment
// @Description delete inquiry comment
// @Accept  json
// @Produce  json
// @Param data body models.UploadBulksRequest true "Form"
// @Success 200 {object} models.UploadBulksResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/bulk_purchase_orders/upload_file [post]
func AdminUploadBulkPurchaseOrders(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.UploadBulksRequest

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return err
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	results, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).UploadBulks(&params)
	if err != nil {
		return err
	}

	return cc.Success(results)
}
