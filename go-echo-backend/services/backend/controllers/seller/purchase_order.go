package controllers

import (
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// SellerPaginatePurchaseOrders
// @Tags Seller-PO
// @Summary Purchase orders
// @Description Purchase orders
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.PurchaseOrder}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders [get]
func SellerPaginatePurchaseOrders(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.PaginatePurchaseOrdersParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	params.SampleMakerID = params.GetUserID()

	var result = repo.NewPurchaseOrderRepo(cc.App.DB).PaginatePurchaseOrders(params)

	return cc.Success(result)
}

// GetPurchaseOrder
// @Tags Seller-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id} [get]
func SellerGetPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.GetPurchaseOrderParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.IncludeSampleMaker = true
	params.SampleMakerID = params.GetUserID()
	params.IncludeInquirySeller = true
	params.IncludeAssignee = true

	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerApprovePurchaseOrderDesign
// @Tags Seller-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/approve_design [post]
func SellerApprovePurchaseOrderDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerApproveDesignParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerApproveDesign(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerPurchaseOrderAddPoUploadComments
// @Tags Seller-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/po_upload_comments [post]
func SellerPurchaseOrderAddPoUploadComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm
	var orderID = cc.GetPathParamString("purchase_order_id")

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypeSellerPoUpload
	params.TargetID = orderID
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	comment, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(comment)
}

// SellerPurchaseOrderPaginatePoUploadComments
// @Tags Seller-PO
// @Summary Paginate comments
// @Description Paginate comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/po_upload_comments [get]
func SellerPurchaseOrderPaginatePoUploadComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var orderID = cc.GetPathParamString("purchase_order_id")
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
	params.TargetID = orderID
	params.TargetType = enums.CommentTargetTypeSellerPoUpload

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// SellerPurchaseOrderPoUploadCommentMarkSeen
// @Tags Admin-Seller-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/po_upload_comments/mark_seen [put]
func SellerPurchaseOrderPoUploadCommentMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerPurchaseOrderUploadCommentMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPurchaseOrderUploadCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// SellerPurchaseOrderConfirmPo
// @Tags Admin-Seller-PO
// @Summary Confirm po
// @Description Confirm po
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/approve_po [post]
func SellerPurchaseOrderApprovePo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerApprovePoParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerApprovePo(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Confirmed")
}

// SellerPurchaseOrderRejectPo
// @Tags Admin-Seller-PO
// @Summary Confirm po
// @Description Confirm po
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/reject_po [delete]
func SellerPurchaseOrderRejectPo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerRejectPoParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerRejectPo(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Rejected")
}

// SellerPurchasePoUploadCommentStatusCount
// @Tags Seller-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/po_upload_comments/status_count [get]
func SellerPurchasePoUploadCommentStatusCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentStatusCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPurchaseOrderUploadCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerPurchaseOrderDesignCommentCreate
// @Tags Seller-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/design_comments [post]
func SellerPurchaseOrderDesignCommentCreate(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm
	var orderID = cc.GetPathParamString("purchase_order_id")

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypeSellerPoDesign
	params.TargetID = orderID
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	comment, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(comment)
}

// SellerPurchaseOrderPaginateDesignComments
// @Tags Seller-PO
// @Summary Paginate comments
// @Description Paginate comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/design_comments [get]
func SellerPurchaseOrderPaginateDesignComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var orderID = cc.GetPathParamString("purchase_order_id")
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
	params.TargetID = orderID
	params.TargetType = enums.CommentTargetTypeSellerPoDesign

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// SellerPurchaseOrderDesignCommentMarkSeen
// @Tags Admin-Seller-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/design_comments/mark_seen [put]
func SellerPurchaseOrderDesignCommentMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPoDesignCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// SellerPurchaseDesignCommentStatusCount
// @Tags Seller-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/design_comments/status_count [get]
func SellerPurchaseDesignCommentStatusCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentStatusCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPoDesignCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerPurchaseOrderMarkRawMaterial
// @Tags Admin-Seller-PO
// @Summary Confirm po
// @Description Confirm po
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/mark_raw_material [put]
func SellerPurchaseOrderMarkRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerMarkRawMaterialParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerMarkRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Confirmed")
}

// SellerPurchaseOrderFinalDesignCommentCreate
// @Tags Seller-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/final_design_comments [post]
func SellerPurchaseOrderFinalDesignCommentCreate(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm
	var orderID = cc.GetPathParamString("purchase_order_id")

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypeSellerPoFinalDesign
	params.TargetID = orderID
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	comment, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(comment)
}

// SellerPurchaseOrderPaginateFinalDesignComments
// @Tags Seller-PO
// @Summary Paginate comments
// @Description Paginate comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/final_design_comments [get]
func SellerPurchaseOrderPaginateFinalDesignComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var orderID = cc.GetPathParamString("purchase_order_id")
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
	params.TargetID = orderID
	params.TargetType = enums.CommentTargetTypeSellerPoFinalDesign

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// SellerPurchaseOrderFinalDesignCommentMarkSeen
// @Tags Admin-Seller-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/final_design_comments/mark_seen [put]
func SellerPurchaseOrderFinalDesignCommentMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerFinalDesignCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// SellerPurchaseOrderFinalDesignCommentStatusCount
// @Tags Seller-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/final_design_comments/status_count [get]
func SellerPurchaseOrderFinalDesignCommentStatusCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentStatusCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerFinalDesignCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerPurchaseOrderUpdateRawMaterial
// @Tags Admin-Seller-PO
// @Summary Confirm po
// @Description Confirm po
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/update_raw_material [put]
func SellerPurchaseOrderUpdateRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerUpdatePurchaseOrderRawMaterialParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	resutl, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerUpdatePurchaseOrderRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(resutl)
}

// SellerPurchaseOrderRawMaterialCommentCreate
// @Tags Seller-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/raw_material_comments [post]
func SellerPurchaseOrderRawMaterialCommentCreate(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm
	var orderID = cc.GetPathParamString("purchase_order_id")

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypeSellerPoRawMaterial
	params.TargetID = orderID
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	comment, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(comment)
}

// SellerPurchaseOrderPaginateDesignComments
// @Tags Seller-PO
// @Summary Paginate comments
// @Description Paginate comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/raw_material_comments [get]
func SellerPurchaseOrderPaginateRawMaterialComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var orderID = cc.GetPathParamString("purchase_order_id")
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
	params.TargetID = orderID
	params.TargetType = enums.CommentTargetTypeSellerPoRawMaterial

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// SellerPurchaseOrderRawMaterialCommentMarkSeen
// @Tags Admin-Seller-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/raw_material_comments/mark_seen [put]
func SellerPurchaseOrderRawMaterialCommentMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPoRawMaterialCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// SellerPurchaseOrderRawMaterialStatusCount
// @Tags Seller-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/raw_material_comments/status_count [get]
func SellerPurchaseOrderRawMaterialStatusCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentStatusCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPoRawMaterialCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerPurchaseOrderMarkMaking
// @Tags Admin-Seller-PO
// @Summary Confirm po
// @Description Confirm po
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/mark_making [put]
func SellerPurchaseOrderMarkMaking(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerPurchaseOrderMarkMakingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPurchaseOrderMarkMaking(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Confirmed")
}

// SellerPurchaseOrderMarkSubmit
// @Tags Admin-Seller-PO
// @Summary Confirm po
// @Description Confirm po
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/mark_submit [put]
func SellerPurchaseOrderMarkSubmit(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerPurchaseOrderMarkSubmitParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerMarkSubmit(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(purchaseOrder)
}

// SellerPurchaseOrderMarkDelivering
// @Tags Seller-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/mark_delivering [put]
func SellerPurchaseOrderMarkDelivering(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerPurchaseOrderMarkDeliveringParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerMarkDelivering(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	// tasks.CreateUserNotificationTask{
	// 	UserID:           purchaseOrder.UserID,
	// 	Message:          fmt.Sprintf("Order %s is on delivering", purchaseOrder.ReferenceID),
	// 	NotificationType: enums.UserNotificationTypePoMarkDelivering,
	// 	Metadata: &models.UserNotificationMetadata{
	// 		AdminID:                  claims.GetUserID(),
	// 		PurchaseOrderID:          purchaseOrder.ID,
	// 		PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
	// 	},
	// }.Dispatch(c.Request().Context())

	// tasks.TrackCustomerIOTask{
	// 	UserID: purchaseOrder.UserID,
	// 	Event:  customerio.EventPoMarkDelivering,
	// 	Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	// }.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// SellerPaginatePurchaseOrderTracking
// @Tags Seller-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/logs [get]
func SellerPaginatePurchaseOrderTracking(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginatePurchaseOrderTrackingParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.UserGroup = enums.PoTrackingUserGroupSeller
	var result = repo.NewPurchaseOrderTrackingRepo(cc.App.DB).PaginatePurchaseOrderTrackings(params)

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
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/update_raw_material [put]
func PurchaseOrderUpdateRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerPurchaseOrderUpdateRawMaterialParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).SellerPurchaseOrderUpdateRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if params.ApproveRawMaterialAt != nil {
		_, _ = tasks.PurchaseOrderRawMaterialApproveTask{
			JwtClaimsInfo:   claims,
			PurchaseOrderID: params.PurchaseOrderID,
		}.DispatchAt(time.Unix(*params.ApproveRawMaterialAt, 0), asynq.MaxRetry(0))
	}

	return cc.Success(purchaseOrder)
}

// SellerPurchaseOrderReceivePayment
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
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/receive_payment [put]
func SellerPurchaseOrderReceivePayment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerPurchaseOrderReceivePaymentParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPurchaseOrderReceivePayment(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerPurchaseOrderSkipRawMaterial
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
// @Router /api/v1/seller/purchase_orders/{purchase_order_id}/skip_raw_material [put]
func SellerPurchaseOrderSkipRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerPurchaseOrderSkipRawMaterialParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).SellerPurchaseOrderSkipRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}
