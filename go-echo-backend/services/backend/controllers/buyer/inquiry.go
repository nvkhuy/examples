package controllers

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// BuyerCreateInquiry
// @Tags Marketplace-Inquiry
// @Summary create inquiry
// @Description create inquiry
// @Accept  json
// @Produce  json
// @Param data body models.InquiryCreateForm true "Form"
// @Success 200 {object} models.Cart
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/create [post]
func BuyerCreateInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.InquiryCreateForm

	err := cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewInquiryBuyerRepo(cc.App.DB).CreateInquiry(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, _ = tasks.CreateCmsNotificationTask{
		Message:          fmt.Sprintf("New inquiry %s was created", result.ReferenceID),
		NotificationType: enums.CmsNotificationTypeNewInquiry,
		Metadata: &models.NotificationMetadata{
			InquiryID: result.ID,
		},
	}.Dispatch(c.Request().Context())

	_, _ = tasks.CreateInquiryAuditTask{
		Form: models.InquiryAuditCreateForm{
			InquiryID:   result.ID,
			ActionType:  enums.AuditActionTypeInquiryCreated,
			UserID:      result.User.ID,
			Description: fmt.Sprintf("%s has been created an inquiry: %s", result.User.Name, result.ReferenceID),
		},
	}.Dispatch(c.Request().Context())

	_, _ = tasks.CreateChatRoomTask{
		UserID:    claims.GetUserID(),
		Role:      claims.GetRole(),
		InquiryID: result.ID,
		BuyerID:   result.UserID,
	}.Dispatch(c.Request().Context())

	if result.User != nil && len(result.User.ContactOwnerIDs) > 0 {
		for _, contactOwnerID := range result.User.ContactOwnerIDs {
			tasks.TrackCustomerIOTask{
				UserID: contactOwnerID,
				Event:  customerio.EventNewInquiry,
				Data:   result.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())

		}
	} else {
		tasks.TrackCustomerIOTask{
			UserID: cc.App.Config.InflowSaleGroupEmail,
			Event:  customerio.EventNewInquiry,
			Data:   result.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	tasks.HubspotSyncInquiryTask{
		InquiryID: result.ID,
		UserID:    claims.GetUserID(),
		IsAdmin:   false,
	}.Dispatch(c.Request().Context())

	_, _ = tasks.UpdateUserProductClassesTask{
		UserID:    claims.GetUserID(),
		InquiryID: result.ID,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}

// BuyerCreateMultipleInquiries
// @Tags Buyer-Inquiry
// @Summary create inquiries
// @Description create inquiries
// @Accept  json
// @Produce  json
// @Param data body []models.InquiryCreateForm true "Form"
// @Success 200 {object} []models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/create_multiple [post]
func BuyerCreateMultipleInquiries(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.CreateMultipleInquiriesRequest

	err := cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewInquiryBuyerRepo(cc.App.DB).CreateMultipleInquiries(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, inquiry := range result {
		_, _ = tasks.CreateCmsNotificationTask{
			Message:          fmt.Sprintf("New inquiry %s was created", inquiry.ReferenceID),
			NotificationType: enums.CmsNotificationTypeNewInquiry,
			Metadata: &models.NotificationMetadata{
				InquiryID: inquiry.ID,
			},
		}.Dispatch(c.Request().Context())

		_, _ = tasks.CreateInquiryAuditTask{
			Form: models.InquiryAuditCreateForm{
				InquiryID:   inquiry.ID,
				ActionType:  enums.AuditActionTypeInquiryCreated,
				UserID:      inquiry.User.ID,
				Description: fmt.Sprintf("%s has been created an inquiry: %s", inquiry.User.Name, inquiry.ReferenceID),
			},
		}.Dispatch(c.Request().Context())

		_, _ = tasks.CreateChatRoomTask{
			UserID:    claims.GetUserID(),
			Role:      claims.GetRole(),
			InquiryID: inquiry.ID,
			BuyerID:   inquiry.UserID,
		}.Dispatch(c.Request().Context())

		if inquiry.User != nil && len(inquiry.User.ContactOwnerIDs) > 0 {
			for _, contactOwnerID := range inquiry.User.ContactOwnerIDs {
				tasks.TrackCustomerIOTask{
					UserID: contactOwnerID,
					Event:  customerio.EventNewInquiry,
					Data:   inquiry.GetCustomerIOMetadata(nil),
				}.Dispatch(c.Request().Context())

			}
		} else {
			tasks.TrackCustomerIOTask{
				UserID: cc.App.Config.InflowSaleGroupEmail,
				Event:  customerio.EventNewInquiry,
				Data:   inquiry.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())
		}

		tasks.HubspotSyncInquiryTask{
			InquiryID: inquiry.ID,
			UserID:    claims.GetUserID(),
			IsAdmin:   false,
		}.Dispatch(c.Request().Context())

		_, _ = tasks.UpdateUserProductClassesTask{
			UserID:    claims.GetUserID(),
			InquiryID: inquiry.ID,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(result)
}

// BuyerUpdateInquiry
// @Tags Marketplace-Inquiry
// @Summary update inquiry
// @Description update inquiry
// @Accept  json
// @Produce  json
// @Param data body models.InquiryUpdateForm true "Form"
// @Success 200 {object} models.Cart
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id} [put]
func BuyerUpdateInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.InquiryUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).UpdateInquiryByID(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, userID := range result.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventBuyerUpdateInquiry,
			Data: map[string]interface{}{
				"before": helper.StructToMap(form),
				"after":  result.GetCustomerIOMetadata(nil),
			},
		}.Dispatch(c.Request().Context())
	}

	tasks.HubspotSyncInquiryTask{
		InquiryID: result.ID,
		UserID:    claims.GetUserID(),
		IsAdmin:   false,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}

// BuyerExtendInquiryEditTimeout
// @Tags Marketplace-Inquiry
// @Summary update inquiry
// @Description update inquiry
// @Accept  json
// @Produce  json
// @Param data body models.InquiryUpdateForm true "Form"
// @Success 200 {object} models.Cart
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/edit_timeout [put]
func BuyerExtendInquiryEditTimeout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.InquiryEditTimeoutUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = repo.NewInquiryRepo(cc.App.DB).UpdateInquiryEditTimeoutByID(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("updated")
}

// BuyerPaginateInquiry
// @Tags Marketplace-Inquiry
// @Summary Inquiry List
// @Description Inquiry List
// @Accept  json
// @Produce  json
// @Param quotation_status query string false "Quotation status"
// @Param page query int false "Page number"
// @Param date_from query int false "Date from"
// @Param date_to query int false "Date to"
// @Param order_reference_id query string false "Order reference"
// @Param statuses query array false "Inquiry Quotation Status filter"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries [get]
func BuyerPaginateInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquiryBuyerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.IncludeCollection = true
	params.JwtClaimsInfo = claims
	params.ExcludeStatuses = append(params.ExcludeStatuses, enums.InquiryStatusClosed)

	var result = repo.NewInquiryBuyerRepo(cc.App.DB).PaginateInquiry(params)

	return cc.Success(result)
}

// BuyerApproveInquiryQuotation
// @Tags Marketplace-Inquiry
// @Summary Approve inquiry quotation
// @Description Approve inquiry quotation
// @Accept  json
// @Produce  json
// @Param data body models.BuyerApproveInquiryQuotationForm true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/approve_quotation [post]
func BuyerApproveInquiryQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerApproveInquiryQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	inquiry, err := repo.NewInquiryBuyerRepo(cc.App.DB).ApproveInquiryQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, assigneeID := range inquiry.AssigneeIDs {
		_, _ = tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventBuyerApproveSkuQuotation,
			Data:   inquiry.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	lastQuotationRequest, _ := repo.NewInquiryAuditRepo(cc.App.DB).GetLastQuotationLog(repo.GetQuotationLogParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		InquiryID:     inquiry.ID,
		ActionType:    enums.AuditActionTypeInquiryAdminSendBuyerQuotation,
	})
	var quotationLogID = ""
	if lastQuotationRequest != nil {
		quotationLogID = lastQuotationRequest.ID
	}

	_, _ = tasks.CreateInquiryAuditTask{
		Form: models.InquiryAuditCreateForm{
			InquiryID:   inquiry.ID,
			ActionType:  enums.AuditActionTypeInquiryBuyerApproveQuotation,
			UserID:      inquiry.User.ID,
			Description: fmt.Sprintf("%s has approved quotation", inquiry.User.Name),
			Metadata: &models.InquiryAuditMetadata{
				After: map[string]interface{}{
					"quotations":          inquiry.AdminQuotations,
					"approve_reject_meta": inquiry.ApproveRejectMeta,
					"quotation_at":        inquiry.QuotationAt,
					"quotation_log_id":    quotationLogID,
				},
			},
		},
	}.Dispatch(c.Request().Context())

	return cc.Success("Approved")
}

// BuyerRejectInquiryQuotation
// @Tags Marketplace-Inquiry
// @Summary Reject inquiry quotation
// @Description Reject inquiry quotation
// @Accept  json
// @Produce  json
// @Param data body models.BuyerRejectInquiryQuotationForm true "Form"
// @Success 200 {string} Rejected
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/reject_quotation [post]
func BuyerRejectInquiryQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerRejectInquiryQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	inquiry, err := repo.NewInquiryBuyerRepo(cc.App.DB).RejectInquiryQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	lastQuotationRequest, _ := repo.NewInquiryAuditRepo(cc.App.DB).GetLastQuotationLog(repo.GetQuotationLogParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		InquiryID:     inquiry.ID,
		ActionType:    enums.AuditActionTypeInquiryAdminSendBuyerQuotation,
	})
	var quotationLogID = ""
	if lastQuotationRequest != nil {
		quotationLogID = lastQuotationRequest.ID
	}

	tasks.CreateInquiryAuditTask{
		Form: models.InquiryAuditCreateForm{
			InquiryID:   inquiry.ID,
			ActionType:  enums.AuditActionTypeInquiryBuyerRejectQuotation,
			UserID:      inquiry.User.ID,
			Description: fmt.Sprintf("%s has rejected quotation", inquiry.User.Name),
			Metadata: &models.InquiryAuditMetadata{
				After: map[string]interface{}{
					"quotations":          inquiry.AdminQuotations,
					"approve_reject_meta": params.ApproveRejectMeta,
					"quotation_at":        inquiry.QuotationAt,
					"quotation_log_id":    quotationLogID,
				},
			},
		},
	}.Dispatch(c.Request().Context())

	for _, assigneeID := range inquiry.AssigneeIDs {
		_, _ = tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventBuyerRejectSkuQuotation,
			Data:   inquiry.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Rejected")
}

// BuyerInquiryQuotationHistory
// @Tags Marketplace-Inquiry
// @Summary Approve inquiry quotation
// @Description Approve inquiry quotation
// @Accept  json
// @Produce  json
// @Param data body models.BuyerApproveInquiryQuotationForm true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/quotation_history [get]
func BuyerInquiryQuotationHistory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquiryAuditsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var id = cc.GetPathParamString("inquiry_id")
	params.InquiryID = id
	params.JwtClaimsInfo = claims

	var result = repo.NewInquiryRepo(cc.App.DB).InquiryQuotationHistory(params)

	return cc.Success(result)
}

// BuyerInquiryLogs
// @Tags Inquiry
// @Summary Inquiry history log
// @Description Inquiry history log
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/logs [get]
func BuyerInquiryLogs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquiryAuditsParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewInquiryRepo(cc.App.DB).PaginateInquiryAudits(params)

	return cc.Success(result)
}

// BuyerCloneInquiry
// @Tags Marketplace-Inquiry
// @Summary Clone an inquiry
// @Description Clone an inquiry
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/clone [post]
func BuyerCloneInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquiryIDParam
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	inquiry, err := repo.NewInquiryRepo(cc.App.DB).CloneInquiryAndQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(inquiry)
}

// BuyerGetInquiry
// @Tags Marketplace-Inquiry
// @Summary Get an inquiry
// @Description Get an inquiry
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id} [get]
func BuyerGetInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetInquiryByIDParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludePurchaseOrder = true
	params.IncludeShippingAddress = true
	params.IncludeCollection = true
	params.InquiryBuilderOptions.IncludeAuditLog = true
	params.UserID = claims.GetUserID()
	result, err := repo.NewInquiryRepo(cc.App.DB).GetInquiryByID(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// BuyerInquiryCollections
// @Tags Marketplace-Inquiry
// @Summary Get inquiry collections
// @Description Get inquiry collections
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.InquiryCollection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/collections [get]
func BuyerInquiryCollections(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	reuslt := repo.NewInquiryRepo(cc.App.DB).PaginateInquiryCollections(repo.PaginateInquiryCollectionParams{
		UserID: claims.ID,
	})

	return cc.Success(reuslt)
}

// BuyerInquiryCollectionCreate
// @Tags Marketplace-Inquiry
// @Summary CreateFromPayload inquiry collection
// @Description CreateFromPayload inquiry collection
// @Accept  json
// @Produce  json
// @Success 200 {object} models.InquiryCollection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/collections [post]
func BuyerInquiryCollectionCreate(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquiryCollectionUpdateForm

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.UserID = claims.ID

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	inquiryCollection, err := repo.NewInquiryRepo(cc.App.DB).CreateInquiryCollection(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(inquiryCollection)
}

// BuyerInquiryUpdateCartItems
// @Tags Marketplace-Inquiry
// @Summary CreateFromPayload inquiry collection
// @Description CreateFromPayload inquiry collection
// @Accept  json
// @Produce  json
// @Success 200 {object} models.InquiryCollection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/update_cart_items [put]
func BuyerInquiryUpdateCartItems(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquiryCartItemsUpdateForm

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	result, err := repo.NewInquiryRepo(cc.App.DB).UpdateCartItems(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerInquiryUpdateAttachments
// @Tags Marketplace-Inquiry
// @Summary CreateFromPayload inquiry collection
// @Description CreateFromPayload inquiry collection
// @Accept  json
// @Produce  json
// @Success 200 {object} models.InquiryCollection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/attachments [put]
func BuyerInquiryUpdateAttachments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateAttachmentsParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	result, err := repo.NewInquiryBuyerRepo(cc.App.DB).UpdateAttachments(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerInquiryCartItems
// @Tags Marketplace-Inquiry
// @Summary CreateFromPayload inquiry collection
// @Description CreateFromPayload inquiry collection
// @Accept  json
// @Produce  json
// @Success 200 {object} models.InquiryCollection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/cart_items [get]
func BuyerInquiryCartItems(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.GetInquiryCartItemsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).GetCartItems(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerInquiryCart
// @Tags Marketplace-Inquiry
// @Summary Buyer inquiry cart
// @Description Buyer inquiry cart
// @Accept  json
// @Produce  json
// @Success 200 {object} models.InquiryCollection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/cart [get]
func BuyerInquiryCart(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCartsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewInquiryRepo(cc.App.DB).PaginateCarts(params)

	return cc.Success(result)
}

// BuyerInquiryPreviewCheckout
// @Tags Marketplace-Inquiry
// @Summary Preview checkout inquiry
// @Description Preview checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.InquiryCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/preview_checkout [post]
func BuyerInquiryPreviewCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquiryPreviewCheckoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.UserID = claims.GetUserID()
	params.UpdatePricing = true
	result, err := repo.NewInquiryRepo(cc.App.DB).InquiryPreviewCheckout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerInquiryCheckout
// @Tags Marketplace-Inquiry
// @Summary Checkout inquiry
// @Description Checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.InquiryCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/checkout [post]
func BuyerInquiryCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquiryCheckoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	purchaseOrder, err := repo.NewInquiryRepo(cc.App.DB).InquiryCheckout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	_, _ = tasks.CreateChatRoomTask{
		UserID:          claims.GetUserID(),
		Role:            claims.GetRole(),
		PurchaseOrderID: purchaseOrder.ID,
		BuyerID:         purchaseOrder.UserID,
	}.Dispatch(c.Request().Context())

	tasks.CreateInquiryAuditTask{
		Form: models.InquiryAuditCreateForm{
			InquiryID:       purchaseOrder.InquiryID,
			ActionType:      enums.AuditActionTypeInquirySamplePoCreated,
			UserID:          purchaseOrder.UserID,
			Description:     fmt.Sprintf("New sample PO %s has been created for inquiry", purchaseOrder.ReferenceID),
			PurchaseOrderID: purchaseOrder.ID,
		},
	}.Dispatch(c.Request().Context())

	tasks.CreateUserNotificationTask{
		UserID:           purchaseOrder.UserID,
		Message:          fmt.Sprintf("New sample PO %s has been created for inquiry", purchaseOrder.ReferenceID),
		NotificationType: enums.UserNotificationTypePoCreated,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                  claims.GetUserID(),
			InquiryID:                purchaseOrder.InquiryID,
			InquiryReferenceID:       purchaseOrder.Inquiry.ReferenceID,
			PurchaseOrderID:          purchaseOrder.ID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	if params.PaymentType == enums.PaymentTypeBankTransfer {
		for _, assigneeID := range purchaseOrder.Inquiry.AssigneeIDs {
			tasks.TrackCustomerIOTask{
				Event:  customerio.EventPoWaitingConfirmBankTransfer,
				UserID: assigneeID,
				Data:   purchaseOrder.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())
		}
	}

	return cc.Success(purchaseOrder)
}

// BuyerInquiryRemoveItems
// @Tags Marketplace-Inquiry
// @Summary Checkout inquiry
// @Description Checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.InquiryCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/{inquiry_id}/remove_items [delete]
func BuyerInquiryRemoveItems(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.InquiryRemoveItemsForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = repo.NewInquiryRepo(cc.App.DB).InquiryRemoveItems(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Removed")
}

// BuyerCloseInquiry
// @Tags Admin-Inquiry
// @Summary Inquiry close
// @Description Inquiry close
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/inquiries/{inquiry_id}/close [put]
func BuyerCloseInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.BuyerInquiryCloseForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = repo.NewInquiryBuyerRepo(cc.App.DB).CloseInquiry(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Closed")
}

// BuyerCancelInquiry
// @Tags Admin-Inquiry
// @Summary Inquiry cancel
// @Description Inquiry cancel
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/inquiries/{inquiry_id}/cancel [delete]
func BuyerCancelInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.BuyerInquiryCancelForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = repo.NewInquiryBuyerRepo(cc.App.DB).CancelInquiry(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Canceled")
}

// BuyerUpdateInquiryLogs
// @Tags Admin-Inquiry
// @Summary Inquiry cancel
// @Description Inquiry cancel
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/inquiries/{inquiry_id}/logs [put]
func BuyerUpdateInquiryLogs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.UpdateInquiryLogsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).UpdateInquiryLogs(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerCreateInquiryLogs
// @Tags Admin-Inquiry
// @Summary Inquiry cancel
// @Description Inquiry cancel
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/inquiries/{inquiry_id}/logs [delete]
func BuyerDeleteInquiryLogs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.DeleteInquiryLogsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).DeleteInquiryLogs(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerConfirmInquiry
// @Tags Admin-Inquiry
// @Summary Inquiry confirm
// @Description Inquiry confirm
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/inquiries/{inquiry_id}/confirm [post]
func BuyerConfirmInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerConfirmInquiryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).BuyerConfirmInquiry(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerApproveMultipleInquiryQuotations
// @Tags Marketplace-Inquiry
// @Summary Approve inquiry quotation
// @Description Approve inquiry quotation
// @Accept  json
// @Produce  json
// @Param data body models.ApproveMultipleInquiryQuotationsRequest true
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/approve_multiple_quotations [post]
func BuyerApproveMultipleInquiryQuotations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.ApproveMultipleInquiryQuotationsRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, "")
	}

	req.JwtClaimsInfo = claims
	inquiries, err := repo.NewInquiryRepo(cc.App.DB).ApproveMultipleInquiryQuotations(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, iq := range inquiries {
		for _, assigneeID := range iq.AssigneeIDs {
			_, _ = tasks.TrackCustomerIOTask{
				UserID: assigneeID,
				Event:  customerio.EventBuyerApproveSkuQuotation,
				Data:   iq.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())
		}
		_, _ = tasks.CreateInquiryAuditTask{
			Form: models.InquiryAuditCreateForm{
				InquiryID:   iq.ID,
				ActionType:  enums.AuditActionTypeInquiryBuyerApproveQuotation,
				UserID:      iq.User.ID,
				Description: fmt.Sprintf("%s has approved quotation", iq.User.Name),
				Metadata: &models.InquiryAuditMetadata{
					After: map[string]interface{}{
						"quotations":          iq.AdminQuotations,
						"approve_reject_meta": iq.ApproveRejectMeta,
						"quotation_at":        iq.QuotationAt,
					},
				},
			},
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Approved")
}

// BuyerRejectMultipleInquiryQuotations
// @Tags Marketplace-Inquiry
// @Summary Reject inquiry quotation
// @Description Reject inquiry quotation
// @Accept  json
// @Produce  json
// @Param data body models.BuyerRejectInquiryQuotationForm true "Form"
// @Success 200 {string} Rejected
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/inquiries/reject_multiple_quotations [post]
func BuyerRejectMultipleInquiryQuotations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.RejectMultipleInquiryQuotationsRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, "")
	}

	req.JwtClaimsInfo = claims
	inquiries, err := repo.NewInquiryBuyerRepo(cc.App.DB).RejectMultipleInquiryQuotations(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, iq := range inquiries {
		tasks.CreateInquiryAuditTask{
			Form: models.InquiryAuditCreateForm{
				InquiryID:   iq.ID,
				ActionType:  enums.AuditActionTypeInquiryBuyerRejectQuotation,
				UserID:      iq.UserID,
				Description: fmt.Sprintf("%s has rejected quotation", iq.User.Name),
				Metadata: &models.InquiryAuditMetadata{
					After: map[string]interface{}{
						"quotations":          iq.AdminQuotations,
						"approve_reject_meta": iq.ApproveRejectMeta,
						"quotation_at":        iq.QuotationAt,
					},
				},
			},
		}.Dispatch(c.Request().Context())

		for _, assigneeID := range iq.AssigneeIDs {
			_, _ = tasks.TrackCustomerIOTask{
				UserID: assigneeID,
				Event:  customerio.EventBuyerRejectSkuQuotation,
				Data:   iq.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())
		}
	}

	return cc.Success("Rejected")
}
