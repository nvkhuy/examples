package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// AdminSellerPurchaseOrderUploadPo
// @Tags Admin-Seller-PO
// @Summary upload po seller purchase order
// @Description upload po purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/upload_po [put]
func AdminSellerPurchaseOrderUploadPo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminSellerPurchaseOrderUploadPoParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	order, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderUploadPo(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(order)
}

// AdminPurchaseOrderAddDesignComments
// @Tags Admin-Seller-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/po_upload_comments [post]
func AdminSellerPurchaseOrderAddPoUploadComments(c echo.Context) error {
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

// AdminSellerPurchaseOrderPaginatePoUploadComments
// @Tags Admin-Seller-PO
// @Summary Paginate comments
// @Description Paginate comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/po_upload_comments [get]
func AdminSellerPurchaseOrderPaginatePoUploadComments(c echo.Context) error {
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

// AdminSellerPurchaseOrderPoUploadCommentMarkSeen
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
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/po_upload_comments/mark_seen [put]
func AdminSellerPurchaseOrderPoUploadCommentMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPurchaseOrderUploadCommentMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderUploadCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// AdminSellerPurchasePoUploadCommentStatusCount
// @Tags Admin-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/po_upload_comments/status_count [get]
func AdminSellerPurchasePoUploadCommentStatusCount(c echo.Context) error {
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

	result := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderUploadCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerPurchaseOrderMarkDesignApproval
// @Tags Admin-Seller-PO
// @Summary Mark design approval
// @Description Mark design approval
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/mark_design_approval [PUT]
func AdminSellerPurchaseOrderMarkDesignApproval(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPurchaseOrderMarkDesignApprovalParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	results, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderMarkDesignApproval(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(results)
}

// AdminSellerPurchaseOrderUpdateDesign
// @Tags Admin-Seller-PO
// @Summary Update design
// @Description Update design
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/update_design [PUT]
func AdminSellerPurchaseOrderUpdateDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPurchaseOrderUpdateDesignParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	results, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderUpdateDesign(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(results)
}

// AdminSellerPurchaseOrderDesignCommentCreate
// @Tags Admin-Seller-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/design_comments [post]
func AdminSellerPurchaseOrderDesignCommentCreate(c echo.Context) error {
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

	if purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrderShortInfo(orderID); err == nil {
		comment.PurchaseOrder = purchaseOrder
		for _, assigneeID := range purchaseOrder.AssigneeIDs {
			tasks.TrackCustomerIOTask{
				UserID: assigneeID,
				Event:  customerio.EventSellerNewDesignComment,
				Data: comment.GetCustomerIOMetadata(map[string]interface{}{
					"sender": purchaseOrder.SampleMaker,
				}),
			}.Dispatch(c.Request().Context())
		}

	}

	return cc.Success(comment)
}

// AdminSellerPurchaseOrderPaginateDesignComments
// @Tags Admin-Seller-PO
// @Summary Paginate comments
// @Description Paginate comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/design_comments [get]
func AdminSellerPurchaseOrderPaginateDesignComments(c echo.Context) error {
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

// AdminSellerPurchaseOrderDesignCommentMarkSeen
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
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/design_comments/mark_seen [put]
func AdminSellerPurchaseOrderDesignCommentMarkSeen(c echo.Context) error {
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

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPoDesignCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// AdminSellerPurchaseDesignCommentStatusCount
// @Tags Admin-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/design_comments/status_count [get]
func AdminSellerPurchaseDesignCommentStatusCount(c echo.Context) error {
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

	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPoDesignCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerPurchaseOrderUpdateFinalDesign
// @Tags Admin-Seller-PO
// @Summary Update design
// @Description Update design
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/update_final_design [PUT]
func AdminSellerPurchaseOrderUpdateFinalDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPurchaseOrderUpdateFinalDesignParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	results, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderUpdateFinalDesign(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(results)
}

// AdminSellerPurchaseOrderApproveFinalDesign
// @Tags Admin-Seller-PO
// @Summary Update design
// @Description Update design
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/approve_final_design [POST]
func AdminSellerPurchaseOrderApproveFinalDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPurchaseOrderApproveFinalDesignParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	results, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderApproveFinalDesign(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(results)
}

// AdminSellerPurchaseOrderFinalDesignCommentCreate
// @Tags Admin-Seller-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/final_design_comments [post]
func AdminSellerPurchaseOrderFinalDesignCommentCreate(c echo.Context) error {
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

// AdminSellerPurchaseOrderPaginateFinalDesignComments
// @Tags Admin-Seller-PO
// @Summary Paginate comments
// @Description Paginate comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/final_design_comments [get]
func AdminSellerPurchaseOrderPaginateFinalDesignComments(c echo.Context) error {
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

// AdminSellerPurchaseOrderFinalDesignCommentMarkSeen
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
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/final_design_comments/mark_seen [put]
func AdminSellerPurchaseOrderFinalDesignCommentMarkSeen(c echo.Context) error {
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

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerFinalDesignCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// AdminSellerPurchaseFinalDesignCommentStatusCount
// @Tags Admin-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/final_design_comments/status_count [get]
func AdminSellerPurchaseFinalDesignCommentStatusCount(c echo.Context) error {
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

	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerFinalDesignCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerPurchaseOrderRawMaterialSendToBuyer
// @Tags Admin-Seller-PO
// @Summary Update design
// @Description Update design
// @Accept  json
// @Produce  json
// @Param data body models.PurchaseOrder true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/raw_material/send_to_buyer [POST]
func AdminSellerPurchaseOrderRawMaterialSendToBuyer(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPORawMaterialSendToBuyerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	results, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPORawMaterialSendToBuyer(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(results)
}

// AdminSellerPurchaseOrderRawMaterialCommentCreate
// @Tags Admin-Seller-PO
// @Summary Add comments
// @Description Add comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/raw_material_comments [post]
func AdminSellerPurchaseOrderRawMaterialCommentCreate(c echo.Context) error {
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

// AdminSellerPurchaseOrderPaginateRawMaterialComments
// @Tags Admin-Seller-PO
// @Summary Paginate comments
// @Description Paginate comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/raw_material_comments [get]
func AdminSellerPurchaseOrderPaginateRawMaterialComments(c echo.Context) error {
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

// AdminSellerPurchaseOrderRawMaterialCommentMarkSeen
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
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/raw_material_comments/mark_seen [put]
func AdminSellerPurchaseOrderRawMaterialCommentMarkSeen(c echo.Context) error {
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

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPoRawMaterialCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// AdminSellerPurchaseRawMaterialCommentStatusCount
// @Tags Admin-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/raw_material_comments/status_count [get]
func AdminSellerPurchaseRawMaterialCommentStatusCount(c echo.Context) error {
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

	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPoRawMaterialCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// PaginateSellerPurchaseOrderTracking
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
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/logs [get]
func PaginateSellerPurchaseOrderTracking(c echo.Context) error {
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

// SellerPurchaseOrderMarkDelivered
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
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/mark_delivered [put]
func AdminSellerPurchaseOrderMarkDelivered(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPoConfirmDeliveredParams
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	purchaseOrder, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPoConfirmDelivered(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(purchaseOrder)
}

// AdminSellerPurchaseOrderDeliveryFeedback
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
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/delivery_feedback [put]
func AdminSellerPurchaseOrderDeliveryFeedback(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPurchaseOrderDeliveryFeedbackParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderDeliveryFeedback(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Sent")
}

// AdminPurchaseOrderAssignMaker Admin PO assign maker
// @Tags Admin-PO
// @Summary Admin PO assign maker
// @Description Admin PO assign maker
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/assign_maker [post]
func AdminPurchaseOrderAssignMaker(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminAssignMakerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	_, err = repo.NewPurchaseOrderRepo(cc.App.DB).AdminAssignMaker(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Assigned sample room")
}

// AdminSellerApproveRawMaterials
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
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/approve_raw_material [post]
func AdminSellerApproveRawMaterials(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerApproveRawMaterialsParmas
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminSellerApproveRawMaterials(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(purchaseOrder)
}

// AdminPayoutPurchaseOrder Admin PO assign maker
// @Tags Admin-PO
// @Summary Admin PO assign maker
// @Description Admin PO assign maker
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/payout [post]
func AdminSellerPurchaseOrderPayout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPurchaseOrderPayoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderPayout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminPayoutPurchaseOrder Admin PO assign maker
// @Tags Admin-PO
// @Summary Admin PO assign maker
// @Description Admin PO assign maker
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_purchase_orders/{purchase_order_id}/preview_payout [post]
func AdminSellerPurchaseOrderPreviewPayout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerPurchaseOrderPreviewPayoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerPurchaseOrderRepo(cc.App.DB).AdminSellerPurchaseOrderPreviewPayout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
