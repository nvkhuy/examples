package controllers

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// Seller
// @Tags Seller-Inquiry
// @Summary Inquriry list
// @Description Inquriry list
// @Accept  json
// @Produce  json
// @Param date_from query number false "Date from"
// @Param date_to query number false "Date to"
// @Param inquiry_statuses query array false "Inquiry Statuses"
// @Param seller_quotation_filter query string false "Seller Quotation filter"
// @Param statuses query array false "Quotation Status filter"
// @Param page query int false "Page number"
// @Param date_from query int false "Date from"
// @Param date_to query int false "Date to"
// @Param order_reference_id query string false "Order reference"
// @Success 200 {object} query.Pagination{records=[]models.Inquiry}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiries/quotations [get]
func SellerInquiryList(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquirySellerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludeInquiry = true
	params.UserID = claims.GetUserID()

	var result = repo.NewInquirySellerRepo(cc.App.DB).PaginateInquirySellerRequest(params)

	return cc.Success(result)
}

// Seller Inquriry Detail
// @Tags Seller-Inquiry
// @Summary Inquriry Detail
// @Description Inquriry Detail
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiries/quotations/{inquiry_seller_id} [get]
func SellerInquiryDetail(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("inquiry_seller_id")

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewInquirySellerRepo(cc.App.DB).GetInquirySellerRequestByID(id, queryfunc.InquirySellerRequestBuilderOptions{
		IncludeInquiry: true,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: claims.GetRole(),
		},
	})
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerSubmitQuotation CreateFromPayload quotion for inquiry
// @Tags Seller-Inquiry
// @Summary  CreateFromPayload quotion for inquiry
// @Description  CreateFromPayload quotion for inquiry
// @Accept  json
// @Produce  json
// @Param data body models.InquirySellerCreateQuatationParams true "Form"
// @Success 200 {object} models.InquirySeller
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiry_quotations/{inquiry_seller_id}/submit_quotation [post]
func SellerSubmitQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquirySellerCreateQuatationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result, err := repo.NewInquirySellerRepo(cc.App.DB).SellerCreateInquiryQuotation(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateInquiryAuditTask{
		Form: models.InquiryAuditCreateForm{
			InquiryID:   result.InquiryID,
			ActionType:  enums.AuditActionTypeInquirySellerSendQuotation,
			UserID:      claims.GetUserID(),
			Description: "Seller has created new quotation",
		},
	}.Dispatch(c.Request().Context())

	if result.Inquiry != nil && len(result.Inquiry.AssigneeIDs) > 0 {
		for _, userID := range result.Inquiry.AssigneeIDs {
			tasks.TrackCustomerIOTask{
				UserID: userID,
				Event:  customerio.EventSellerSubmitQuotation,
				Data:   result.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())
		}
	}

	return cc.Success(result)
}

// SellerSubmitQuotation CreateFromPayload quotion for inquiry
// @Tags Seller-Inquiry
// @Summary  CreateFromPayload quotion for inquiry
// @Description  CreateFromPayload quotion for inquiry
// @Accept  json
// @Produce  json
// @Param data body models.InquirySellerCreateQuatationParams true "Form"
// @Success 200 {object} models.InquirySeller
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiry_quotations/{inquiry_seller_id}/approve_offer [post]
func SellerApproveOffer(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquirySellerApproveOfferParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewInquirySellerRepo(cc.App.DB).InquirySellerApproveOffer(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success("Approved")
}

// SellerRejectOffer Seller rejects offer
// @Tags Seller-Inquiry
// @Summary Seller rejects offer
// @Description Seller rejects offer
// @Accept  json
// @Produce  json
// @Param data body models.InquirySellerCreateQuatationParams true "Form"
// @Success 200 {object} models.InquirySeller
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiry_quotations/{inquiry_seller_id}/reject_offer [delete]
func SellerRejectOffer(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquirySellerRejectOfferParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewInquirySellerRepo(cc.App.DB).InquirySellerRejectOffer(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success("Rejected")
}

// Seller
// @Tags Seller-Inquiry
// @Summary Inquriry list
// @Description Inquriry list
// @Accept  json
// @Produce  json
// @Param date_from query number false "Date from"
// @Param date_to query number false "Date to"
// @Param inquiry_statuses query array false "Inquiry Statuses"
// @Param seller_quotation_filter query string false "Seller Quotation filter"
// @Param page query int false "Page number"
// @Param date_from query int false "Date from"
// @Param date_to query int false "Date to"
// @Success 200 {object} query.Pagination{records=[]models.Inquiry}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiries/{id}/audits [get]
func SellerInquiryAuditList(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquirySellerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludeInquiry = true

	var result = repo.NewInquirySellerRepo(cc.App.DB).PaginateInquirySellerRequest(params)

	return cc.Success(result)
}

// Seller
// @Tags Seller-Inquiry
// @Summary Inquriry list
// @Description Inquriry list
// @Accept  json
// @Produce  json
// @Param date_from query number false "Date from"
// @Param date_to query number false "Date to"
// @Param inquiry_statuses query array false "Inquiry Statuses"
// @Param seller_quotation_filter query string false "Seller Quotation filter"
// @Param page query int false "Page number"
// @Param date_from query int false "Date from"
// @Param date_to query int false "Date to"
// @Success 200 {object} query.Pagination{records=[]models.Inquiry}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiries/{id}/latest_audit [get]
func SellerInquiryAuditLatest(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquirySellerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludeInquiry = true

	var result = repo.NewInquirySellerRepo(cc.App.DB).PaginateInquirySellerRequest(params)

	return cc.Success(result)
}

// SellerInquiryCreateComment CreateFromPayload inquiry comment
// @Tags Seller-Inquiry
// @Summary  CreateFromPayload quotion for inquiry
// @Description  CreateFromPayload quotion for inquiry
// @Accept  json
// @Produce  json
// @Param data body models.InquirySellerCreateQuatationParams true "Form"
// @Success 200 {object} models.InquirySeller
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiry_quotations/{inquiry_seller_id}/comments [post]
func SellerInquiryCreateComment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm
	var id = cc.GetPathParamString("inquiry_seller_id")

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	sellerRequest, err := repo.NewInquirySellerRepo(cc.App.DB).
		GetInquirySellerRequestByID(id, queryfunc.InquirySellerRequestBuilderOptions{IncludeInquiry: true})
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.TargetType = enums.CommentTargetTypeInquirySellerDesign
	params.TargetID = sellerRequest.InquiryID
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.CreateCmsNotificationTask{
		Message:          fmt.Sprintf("Seller comment on inquiry design %s", sellerRequest.Inquiry.ReferenceID),
		NotificationType: enums.CmsNotificationTypeInquirySellerDesignComment,
		Metadata: &models.NotificationMetadata{
			CommentID:          result.ID,
			InquiryID:          sellerRequest.Inquiry.ID,
			InquiryReferenceID: sellerRequest.Inquiry.ReferenceID,
			InquirySellerID:    sellerRequest.ID,
		},
	}.Dispatch(c.Request().Context())

	for _, userID := range sellerRequest.Inquiry.AssigneeIDs {
		if sender, err := repo.NewUserRepo(cc.App.DB).GetShortUserInfo(params.GetUserID()); err == nil {
			tasks.TrackCustomerIOTask{
				UserID: userID,
				Event:  customerio.EventSellerCommentOnSellerRequest,
				Data: result.GetCustomerIOMetadata(map[string]interface{}{
					"sender":  sender,
					"request": sellerRequest.GetCustomerIOMetadata(nil),
				}),
			}.Dispatch(cc.Request().Context())
		}
	}

	return cc.Success(result)
}

// SellerInquiryCommentList
// @Tags Seller-Inquiry
// @Summary Inquiry seller list
// @Description Inquiry seller list
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiry_quotations/{inquiry_seller_id}/comments [get]
func SellerInquiryCommentList(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquirySellerCommentsParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	var results = repo.NewCommentRepo(cc.App.DB).PaginateInquirySelerComments(params)

	return cc.Success(results)
}

// SellerInquiryCommentMarkSeen
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
// @Router /api/v1/seller/inquiry_quotations/{inquiry_seller_id}/comments/mark_seen [put]
func SellerInquiryCommentMarkSeen(c echo.Context) error {
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

// SellerInquiryCommentStatusCount
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
// @Router /api/v1/seller/inquiry_quotations/{inquiry_seller_id}/status_count [get]
func SellerInquiryCommentStatusCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquirySellerDesignCommenStatusCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewInquirySellerRepo(cc.App.DB).InquiryDesignCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerSubmitMultipleQuotations
// @Tags Seller-Inquiry
// @Summary  Seller submit multiple quotations for inquiry
// @Description  Seller submit multiple quotations for inquiry
// @Accept  json
// @Produce  json
// @Param data body models.SubmitMultipleInquirySellerQuotationRequest true "Form"
// @Success 200 {object} models.InquirySellers
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/inquiry_quotations/submit_multiple_quotations [post]
func SellerSubmitMultipleQuotations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.SubmitMultipleInquirySellerQuotationRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	result, err := repo.NewInquirySellerRepo(cc.App.DB).SubmitMultipleInquirySellerQuotations(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, iqSeller := range result {
		tasks.CreateInquiryAuditTask{
			Form: models.InquiryAuditCreateForm{
				InquiryID:   iqSeller.InquiryID,
				ActionType:  enums.AuditActionTypeInquirySellerSendQuotation,
				UserID:      claims.GetUserID(),
				Description: "Seller has created new quotation",
			},
		}.Dispatch(c.Request().Context())

		if iqSeller.Inquiry != nil && len(iqSeller.Inquiry.AssigneeIDs) > 0 {
			for _, userID := range iqSeller.Inquiry.AssigneeIDs {
				tasks.TrackCustomerIOTask{
					UserID: userID,
					Event:  customerio.EventSellerSubmitQuotation,
					Data:   iqSeller.GetCustomerIOMetadata(nil),
				}.Dispatch(c.Request().Context())
			}
		}
	}

	return cc.Success(result)
}
