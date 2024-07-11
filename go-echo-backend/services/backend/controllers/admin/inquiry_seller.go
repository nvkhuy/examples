package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// AdminInquirySellerCreateComment
// @Tags Admin-Inquiry
// @Summary Add seller request comments
// @Description Add seller request comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiry_seller/{inquiry_seller_id}/design_comments [post]
func AdminInquirySellerCreateDesignComment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var requestID = cc.GetPathParamString("inquiry_seller_id")
	var params models.CommentCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypeInquirySellerDesign
	params.TargetID = requestID
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	comment, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if request, err := repo.NewInquirySellerRepo(cc.App.DB).GetInquirySellerRequestByID(requestID, queryfunc.InquirySellerRequestBuilderOptions{}); err == nil {
		if sender, err := repo.NewUserRepo(cc.App.DB).GetShortUserInfo(params.GetUserID()); err == nil {
			tasks.TrackCustomerIOTask{
				UserID: request.UserID,
				Event:  customerio.EventAdminDesignCommentOnSellerRequest,
				Data: comment.GetCustomerIOMetadata(map[string]interface{}{
					"sender":  sender,
					"request": request.GetCustomerIOMetadata(nil),
				}),
			}.Dispatch(cc.Request().Context())
		}
	}

	return cc.Success(comment)
}

// AdminInquirySellerPaginateDesignComments
// @Tags Admin-Inquiry
// @Summary Paginate design comments
// @Description Paginate design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiry_sellers/{inquiry_seller_id}/design_comments [get]
func AdminInquirySellerPaginateDesignComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var requestID = cc.GetPathParamString("inquiry_seller_id")
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
	params.TargetID = requestID
	params.TargetType = enums.CommentTargetTypeInquirySellerDesign

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// AdminPurchaseOrderDesignCommentMarkSeen
// @Tags Admin-Inquiry
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiry_sellers/{inquiry_seller_id}/design_comments/mark_seen [put]
func AdminInquirySellerPaginateDesignCommentsMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquirySellerDesignCommentMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewInquirySellerRepo(cc.App.DB).InquiryDesignCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// AdminInquirySellerApproveQuotation
// @Tags Admin-Inquiry
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiry_sellers/{inquiry_seller_id}/approve [post]
func AdminInquirySellerApproveQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminInquirySellerApproveQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	iqSeller, err := repo.NewInquirySellerRepo(cc.App.DB).AdminInquirySellerApproveQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.AdminApproveSellerQuotationTask{
		AdminID:         claims.GetUserID(),
		InquirySellerID: params.InquirySellerID,
	}.Dispatch(c.Request().Context())

	tasks.CreateChatRoomTask{
		UserID:          claims.GetUserID(),
		Role:            claims.GetRole(),
		PurchaseOrderID: iqSeller.PurchaseOrderID,
		SellerID:        iqSeller.UserID,
	}.Dispatch(c.Request().Context())

	return cc.Success("Approved")
}

// AdminInquirySellerRejectQuotation
// @Tags Admin-Inquiry
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiry_sellers/{inquiry_seller_id}/reject [delete]
func AdminInquirySellerRejectQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminInquirySellerRejectQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	err = repo.NewInquirySellerRepo(cc.App.DB).AdminInquirySellerRejectQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.AdminRejectSellerQuotationTask{
		AdminID:         claims.GetUserID(),
		InquirySellerID: params.InquirySellerID,
	}.Dispatch(c.Request().Context())

	return cc.Success("Rejected")
}
