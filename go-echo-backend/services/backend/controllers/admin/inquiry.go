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

// AdminPaginateInquiry
// @Tags Admin-Inquiry
// @Summary Inquiry List
// @Description Inquiry List
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param user_id query string false "UserID"
// @Param date_from query number false "Date from"
// @Param date_to query number false "Date to"
// @Param order_reference_id query string false "Order reference"
// @Param page query int false "Page number"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries [get]
func AdminPaginateInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquiryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.IncludePurchaseOrder = true
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludeAssignee = true
	params.IncludeUser = true
	params.IncludeAuditLog = true
	params.IncludeShippingAddress = true
	params.IncludeCollection = true
	var result = repo.NewInquiryRepo(cc.App.DB).PaginateInquiry(params)
	return cc.Success(result)
}

// AdminInquiryDetail
// @Tags Admin-Inquiry
// @Summary Inquiry Detail
// @Description Inquiry Detail
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id} [get]
func AdminInquiryDetail(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetInquiryByIDParams

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
	params.IncludeShippingAddress = true
	params.IncludePurchaseOrder = true
	params.IncludeAssignee = true
	result, err := repo.NewInquiryRepo(cc.App.DB).GetInquiryByID(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminInquiryUpdate
// @Tags Admin-Inquiry
// @Summary Inquiry update
// @Description Inquiry update
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id} [put]
func AdminInquiryUpdate(c echo.Context) error {
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

	return cc.Success(result)
}

// AdminPaginateInquirySellerRequests Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/requests [get]
func AdminPaginateInquirySellerRequests(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquirySellerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludeUnseenCommentCount = true
	var result = repo.NewInquirySellerRepo(cc.App.DB).PaginateInquirySellerRequest(params)

	return cc.Success(result)
}

// AdminPaginateMatchingSellers Admin paginate matching sellers
// @Tags Admin-Inquiry
// @Summary Admin paginate matching sellers
// @Description Admin paginate matching sellers
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/matching_sellers [get]
func AdminPaginateMatchingSellers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateMatchingSellersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewInquirySellerRepo(cc.App.DB).PaginateMatchingSellers(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSendInquiryToBuyer Admin send quotation to buyer
// @Tags Admin-Inquiry
// @Summary Admin send quotation to buyer
// @Description Admin send quotation to buyer
// @Accept  json
// @Produce  json
// @Param data body models.SendInquiryToBuyerForm true "Form"
// @Success 200 {string} Sent
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/send_to_buyer [post]
func AdminSendInquiryToBuyer(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.SendInquiryToBuyerForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	resp, err := repo.NewInquiryRepo(cc.App.DB).SendInquiryToBuyer(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, _ = tasks.SendInquiryToBuyerTask{
		UserID: resp.Inquiry.User.ID,
		Event:  customerio.EventAdminSentQuotationToBuyer,
		Data:   resp.Inquiry.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	tasks.CreateUserNotificationTask{
		UserID:           resp.Inquiry.UserID,
		Message:          fmt.Sprintf("New quotation was created for %s", resp.Inquiry.ReferenceID),
		NotificationType: enums.UserNotificationTypeInquirySubmitQuotation,
		Metadata: &models.UserNotificationMetadata{
			AdminID:            claims.GetUserID(),
			InquiryID:          resp.Inquiry.ID,
			InquiryReferenceID: resp.Inquiry.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	tasks.CreateInquiryAuditTask{
		Form: models.InquiryAuditCreateForm{
			InquiryID:   resp.Inquiry.ID,
			ActionType:  enums.AuditActionTypeInquiryAdminSendBuyerQuotation,
			UserID:      claims.GetUserID(),
			Description: fmt.Sprintf("%s has sent quotation to buyer %s", resp.Admin.Name, resp.Inquiry.User.Name),
			Metadata: &models.InquiryAuditMetadata{
				After: map[string]interface{}{
					"quotations": form.Quotations,
				},
			},
		},
	}.Dispatch(c.Request().Context())

	return cc.Success("Sent")
}

// AdminSubmitInquiryQuotation Admin submit quotation
// @Tags Admin-Inquiry
// @Summary Admin send quotation to buyer
// @Description Admin send quotation to buyer
// @Accept  json
// @Produce  json
// @Param data body models.SendInquiryToBuyerForm true "Form"
// @Success 200 {string} Sent
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/submit_quotation [post]
func AdminSubmitInquiryQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.SendInquiryToBuyerForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	resp, err := repo.NewInquiryRepo(cc.App.DB).AdminSubmitQuotation(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	// CMS notification
	tasks.CreateCmsNotificationTask{
		Message:          fmt.Sprintf("New quotation was created for %s", resp.Inquiry.ReferenceID),
		NotificationType: enums.CmsNotificationTypeNewInquiryQuotation,
	}.Dispatch(c.Request().Context())

	return cc.Success("Sent")
}

// AdminSubmitMultipleInquiryQuotations Admin submit quotation
// @Tags Admin-Inquiry
// @Summary Admin send quotation to buyer
// @Description Admin send quotation to buyer
// @Accept  json
// @Produce  json
// @Param data body models.SubmitMultipleInquiryQuotationRequest true "Form"
// @Success 200 {string} Sent
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/submit_multiple_quotations [post]
func AdminSubmitMultipleInquiryQuotations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.SubmitMultipleInquiryQuotationRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	form.JwtClaimsInfo = claims

	var admin models.User
	if err := cc.App.DB.Select("Name").First(&admin, "id = ?", claims.GetUserID()).Error; err != nil {
		return err
	}
	resp, err := repo.NewInquiryRepo(cc.App.DB).AdminSubmitMultipleInquiryQuotations(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, inquiry := range resp {
		_, _ = tasks.SendInquiryToBuyerTask{
			UserID: inquiry.User.ID,
			Event:  customerio.EventAdminSentQuotationToBuyer,
			Data:   inquiry.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())

		tasks.CreateUserNotificationTask{
			UserID:           inquiry.UserID,
			Message:          fmt.Sprintf("New quotation was created for %s", inquiry.ReferenceID),
			NotificationType: enums.UserNotificationTypeInquirySubmitQuotation,
			Metadata: &models.UserNotificationMetadata{
				AdminID:            claims.GetUserID(),
				InquiryID:          inquiry.ID,
				InquiryReferenceID: inquiry.ReferenceID,
			},
		}.Dispatch(c.Request().Context())

		tasks.CreateInquiryAuditTask{
			Form: models.InquiryAuditCreateForm{
				InquiryID:   inquiry.ID,
				ActionType:  enums.AuditActionTypeInquiryAdminSendBuyerQuotation,
				UserID:      claims.GetUserID(),
				Description: fmt.Sprintf("%s has sent quotation to buyer %s", admin.Name, inquiry.User.Name),
				Metadata: &models.InquiryAuditMetadata{
					After: map[string]interface{}{
						"quotations": inquiry.AdminQuotations,
					},
				},
			},
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Sent")
}

// AdminInquiryInternalApproveQuotation Internal admin approve quotation
// @Tags Admin-Inquiry
// @Summary Internal admin approve quotation
// @Description Internal admin approve quotation
// @Accept  json
// @Produce  json
// @Param data body models.SendInquiryToBuyerForm true "Form"
// @Success 200 {string} Sent
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/internal_approve_quotation [post]
func AdminInquiryInternalApproveQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.AdminInternalApproveQuotationForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	resp, err := repo.NewInquiryRepo(cc.App.DB).AdminInternalApproveQuotation(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	taskData := map[string]interface{}{
		"product_name":   resp.Inquiry.Title,
		"user_name":      resp.Inquiry.User.Name,
		"email":          resp.Inquiry.User.Email,
		"link_quotation": fmt.Sprintf("%s/login", cc.App.Config.WebAppBaseURL),
	}

	_, _ = tasks.SendInquiryToBuyerTask{
		UserID: resp.Inquiry.User.ID,
		Event:  customerio.EventAdminSentQuotationToBuyer,
		Data:   taskData,
	}.Dispatch(c.Request().Context())

	tasks.CreateUserNotificationTask{
		UserID:           resp.Inquiry.UserID,
		Message:          fmt.Sprintf("New quotation was created for %s", resp.Inquiry.ReferenceID),
		NotificationType: enums.UserNotificationTypeInquirySubmitQuotation,
		Metadata: &models.UserNotificationMetadata{
			AdminID:            claims.GetUserID(),
			InquiryID:          resp.Inquiry.ID,
			InquiryReferenceID: resp.Inquiry.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	return cc.Success("Sent")
}

// AdminInquiryListForCreatingOrder
// @Tags Admin-Inquiry
// @Summary Inquiry list for create order
// @Description Inquiry list for create order
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/for_creating_order [get]
func AdminInquiryListForCreatingOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquiryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.Limit = 1000 // will fetch all inquiry, update this if inquiry count > 1000
	var result = repo.NewInquiryRepo(cc.App.DB).PaginateInquiryForCreateOrder(params)
	return cc.Success(result)
}

// AdminInquiryQuotationHistory
// @Tags Admin-Inquiry
// @Summary AdminInquiryQuotationHistory
// @Description AdminInquiryQuotationHistory
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/quotation_history [get]
func AdminInquiryQuotationHistory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquiryAuditsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var id = cc.GetPathParamString("inquiry_id")
	params.InquiryID = id
	params.JwtClaimsInfo = claims

	var result = repo.NewInquiryRepo(cc.App.DB).InquiryQuotationHistory(params)

	return cc.Success(result)
}

// AdminInquiryLogs
// @Tags Admin-Inquiry
// @Summary Inquiry list for create order
// @Description Inquiry list for create order
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/logs [get]
func AdminInquiryLogs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInquiryAuditsParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewInquiryRepo(cc.App.DB).PaginateInquiryAudits(params)

	return cc.Success(result)
}

// AdminInquiryMarkSeen
// @Tags Admin-Inquiry
// @Summary Inquiry list for create order
// @Description Inquiry list for create order
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/mark_seen [put]
func AdminInquiryMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquiryMarkSeenForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewInquiryRepo(cc.App.DB).AdminMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// Admin Inquiry list for create order
// @Tags Admin-Inquiry
// @Summary Inquiry list for create order
// @Description Inquiry list for create order
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/clone [put]
func AdminCloneInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquiryIDParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	inquiry, err := repo.NewInquiryRepo(cc.App.DB).CloneInquiry(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(inquiry)
}

// AdminInquiryMarkAsPaid
// @Tags Admin-Inquiry
// @Summary Mark as paid
// @Description Mark as paid
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/mark_as_paid [post]
func AdminInquiryMarkAsPaid(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquiryIDParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	// Handle multi inquiry same payment
	var purchaseOrder models.PurchaseOrder
	err = cc.App.DB.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs", "CheckoutSessionID").First(&purchaseOrder, "inquiry_id = ?", params.InquiryID).Error
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.PurchaseOrder = &purchaseOrder
	if purchaseOrder.CheckoutSessionID != "" {
		purchaseOrders, err := repo.NewInquiryRepo(cc.App.DB).MultiInquiryMarkAsPaid(repo.MultiInquiryParams{
			CheckoutSessionID: purchaseOrder.CheckoutSessionID,
			Note:              params.Note,
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		for _, purchaseOrder := range purchaseOrders {
			_, _ = tasks.PurchaseOrderBankTransferConfirmedTask{
				ApprovedByUserID: claims.GetUserID(),
				PurchaseOrderID:  purchaseOrder.ID,
			}.Dispatch(c.Request().Context())
		}

	} else {
		purchaseOrder, err := repo.NewInquiryRepo(cc.App.DB).InquiryMarkAsPaid(params)
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		_, _ = tasks.PurchaseOrderBankTransferConfirmedTask{
			ApprovedByUserID: claims.GetUserID(),
			PurchaseOrderID:  purchaseOrder.ID,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Confirmed")
}

// AdminInquiryMarkAsUnpaid
// @Tags Admin-Inquiry
// @Summary Mark as un-paid
// @Description Mark as un-paid
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/mark_as_unpaid [post]
func AdminInquiryMarkAsUnpaid(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquiryIDParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	// Handle multi inquiry same payment
	var order models.PurchaseOrder
	err = cc.App.DB.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs", "CheckoutSessionID").First(&order, "inquiry_id = ?", params.InquiryID).Error
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if order.CheckoutSessionID != "" {
		purchaseOrders, err := repo.NewInquiryRepo(cc.App.DB).MultiInquiryMarkAsUnpaid(repo.MultiInquiryParams{
			CheckoutSessionID: order.CheckoutSessionID,
			Note:              params.Note,
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		for _, purchaseOrder := range purchaseOrders {
			tasks.PurchaseOrderBankTransferRejectedTask{
				ApprovedByUserID: claims.GetUserID(),
				PurchaseOrderID:  purchaseOrder.ID,
			}.Dispatch(c.Request().Context())
		}

	} else {
		purchaseOrder, err := repo.NewInquiryRepo(cc.App.DB).InquiryMarkAsUnpaid(params)
		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		tasks.PurchaseOrderBankTransferRejectedTask{
			ApprovedByUserID: claims.GetUserID(),
			PurchaseOrderID:  purchaseOrder.ID,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("UnPaid")
}

// AdminInquiryAssignPIC
// @Tags Admin-Inquiry
// @Summary Assign PIC
// @Description Assign PIC
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/assign_pic [put]
func AdminInquiryAssignPIC(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.InquiryAssignPICParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).InquiryAssignPIC(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, userID := range result.AssigneeIDs {
		tasks.AssignInquiryPICTask{
			AssignerID: claims.GetUserID(),
			AssigneeID: userID,
			InquiryID:  params.InquiryID,
		}.Dispatch(c.Request().Context())

	}
	return cc.Success(result)
}

// AdminCreateInquiry Admin send inquiry to seller
// @Tags Admin-Inquiry
// @Summary Admin send inquiry to seller
// @Description Admin send inquiry to seller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries [post]
func AdminCreateInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.InquiryAdminCreateForm

	err := cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).AdminCreateInquiry(form)
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

	tasks.HubspotSyncInquiryTask{
		InquiryID: result.ID,
		UserID:    claims.GetUserID(),
		IsAdmin:   false,
	}.Dispatch(c.Request().Context())

	if form.BuyerId != "" {
		_, _ = tasks.UpdateUserProductClassesTask{
			UserID:          claims.GetUserID(),
			PurchaseOrderID: form.BuyerId,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(result)
}

// AdminArchiveInquiry Admin archive inquiry
// @Tags Admin-Inquiry
// @Summary Admin archive inquiry
// @Description Admin archive inquiry
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/archive [delete]
func AdminArchiveInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminUnarchiveInquiryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewInquiryRepo(cc.App.DB).AdminArchiveInquiry(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}

// AdminUnarchiveInquiry Admin unarchive inquiry
// @Tags Admin-Inquiry
// @Summary Admin unarchive inquiry
// @Description Admin unarchive inquiry
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/unarchive [put]
func AdminUnarchiveInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminUnarchiveInquiryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	err = repo.NewInquiryRepo(cc.App.DB).AdminUnarchiveInquiry(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Unarchived")
}

// AdminDeleteInquiry Admin archive inquiry
// @Tags Admin-Inquiry
// @Summary Admin archive inquiry
// @Description Admin archive inquiry
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/delete [delete]
func AdminDeleteInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminDeleteInquiryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewInquiryRepo(cc.App.DB).AdminDeleteInquiry(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}

// AdminInquiryCreatePaymentLink Admin create payment link
// @Tags Admin-Inquiry
// @Summary Admin create payment link
// @Description Admin create payment link
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/payment_link [post]
func AdminInquiryCreatePaymentLink(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreatePaymentLinkParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	link, err := repo.NewInquiryRepo(cc.App.DB).InquiryCreatePaymentLink(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(link)
}

// AdminCreateBuyerPaymentLink Admin create payment link
// @Tags Admin-Inquiry
// @Summary Admin create payment link
// @Description Admin create payment link
// @Accept  json
// @Produce  json
// @params data body repo.CreateBuyerPaymentLinkRequest true
// @Success 200 {object} repo.CreateBuyerPaymentLinkResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/payment_link [post]
func AdminCreateBuyerPaymentLink(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateBuyerPaymentLinkRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	resp, err := repo.NewInquiryRepo(cc.App.DB).CreateBuyerPaymentLink(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(resp)
}

// AdminInquirySyncSample Admin sync sample data
// @Tags Admin-Inquiry
// @Summary Admin sync sample data
// @Description Admin sync sample data
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/sync_sample [post]
func AdminInquirySyncSample(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquirySyncSampleParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).InquirySyncSample(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSendInquiryToSeller Admin send inquiry to seller
// @Tags Admin-Inquiry
// @Summary Admin send inquiry to seller
// @Description Admin send inquiry to seller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/send_to_seller [post]
func AdminSendInquiryToSeller(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.SendInquiryToSellerParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).SendInquiryToSeller(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if len(result) > 0 {
		inquiry, err := repo.NewInquiryRepo(cc.App.DB).GetInquiryByID(repo.GetInquiryByIDParams{
			InquiryID:     params.InquiryID,
			JwtClaimsInfo: params.JwtClaimsInfo,
		})
		if err == nil {
			for _, record := range result {
				tasks.TrackCustomerIOTask{
					UserID: record.UserID,
					Event:  customerio.EventSellerNewRFQRequest,
					Data: inquiry.GetCustomerIOMetadata(map[string]interface{}{
						"offer_price":  record.OfferPrice,
						"offer_remark": record.OfferRemark,
					}),
				}.Dispatch(cc.Request().Context())

				tasks.CreateChatRoomTask{
					UserID:    claims.GetUserID(),
					Role:      claims.GetRole(),
					InquiryID: inquiry.ID,
					SellerID:  record.UserID,
				}.Dispatch(c.Request().Context())
			}
		}

	}

	return cc.Success(result)
}

// AdminUpdateInquiryCosting Admin update inquiry costing
// @Tags Admin-Inquiry
// @Summary Admin update inquiry costing
// @Description Admin update inquiry costing
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/update_costing [put]
func AdminUpdateInquiryCosting(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var form repo.UpdateInquiryCostingParams
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).UpdateInquiryCosting(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminCloseInquiry
// @Tags Admin-Inquiry
// @Summary Inquiry close
// @Description Inquiry close
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/close [put]
func AdminCloseInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.InquiryCloseForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = repo.NewInquiryRepo(cc.App.DB).CloseInquiry(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Closed")
}

// InquiryAddNote
// @Tags Admin-Inquiry
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/notes [post]
func InquiryAddNote(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypeInquiryInternalNotes
	params.TargetID = cc.GetPathParamString("inquiry_id")
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.NewInquiryNotesTask{
		UserID:         claims.GetUserID(),
		InquiryID:      params.TargetID,
		MentionUserIDs: params.MentionUserIDs,
		Message:        params.Message,
		Attachments:    params.Attachments,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}

// PaginateInquiryNotes
// @Tags Admin-Inquiry
// @Summary Get design comments
// @Description Get design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/notes [get]
func PaginateInquiryNotes(c echo.Context) error {
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
	params.TargetID = cc.GetPathParamString("inquiry_id")
	params.TargetType = enums.CommentTargetTypeInquiryInternalNotes
	params.OrderByQuery = "c.created_at DESC"

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// InquiryNoteMarkSeen
// @Tags Marketplace-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/notes/mark_seen [put]
func InquiryNoteMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquiryNoteMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewInquiryRepo(cc.App.DB).InquiryNoteMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// InquiryNoteStatusCount
// @Tags Marketplace-PO
// @Summary Design comment status count
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/notes/unread_count [get]
func InquiryNoteUnreadCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquiryNoteUnreadCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var result = repo.NewInquiryRepo(cc.App.DB).InquiryNoteUnreadCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminInquiryCommentDelete
// @Tags Admin-Inquiry
// @Summary delete inquiry comment
// @Description delete inquiry comment
// @Accept  json
// @Produce  json
// @Param data body models.ContentCommentCreateForm true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/notes/{comment_id} [delete]
func AdminInquiryCommentDelete(c echo.Context) error {
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

// ExportInquiries
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/export [get]
func ExportInquiries(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.ExportInquiriesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	resp, err := repo.NewInquiryRepo(cc.App.DB).ExportExcel(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(resp)
}

// AdminInquirySellerRequestCreateComment
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
// @Router /api/v1/admin/inquiries/{inquiry_id}/seller_requests/{inquiry_seller_id}/comments [post]
func AdminInquirySellerRequestCreateComment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var requestID = cc.GetPathParamString("inquiry_seller_id")
	var params models.CommentCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypeInquirySellerRequest
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

// AdminInquirySellerRequestPaginateComments
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
// @Router /api/v1/admin/inquiries/{inquiry_id}/seller_requests/{inquiry_seller_id}/comments [get]
func AdminInquirySellerRequestPaginateComments(c echo.Context) error {
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
	params.TargetType = enums.CommentTargetTypeInquirySellerRequest

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
// @Router /api/v1/admin/inquiries/{inquiry_id}/seller_requests/{inquiry_seller_id}/comments/mark_seen [put]
func AdminInquirySellerRequestCommentMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquirySellerQuotationCommentMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewInquirySellerRepo(cc.App.DB).InquiryQuotationCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// AdminInquirySellerStatusCount
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
// @Router /api/v1/admin/inquiries/{inquiry_id}/seller_requests/status_count [get]
func AdminInquirySellerStatusCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquirySellerStatusCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewInquirySellerRepo(cc.App.DB).StatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminInquiryPreviewCheckout
// @Tags Admin-Inquiry
// @Summary Preview checkout inquiry
// @Description Preview checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.InquiryPreviewCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/{inquiry_id}/preview_checkout [post]
func AdminInquiryPreviewCheckout(c echo.Context) error {
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
	result, err := repo.NewInquiryRepo(cc.App.DB).InquiryPreviewCheckout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminInquirySellerAllocationSearchSeller
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
// @Router /api/v1/admin/inquiries/{inquiry_id}/seller_allocations [get]
func AdminInquirySellerAllocationSearchSeller(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InquirySellerAllocationSearchSellerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result := repo.NewInquirySellerRepo(cc.App.DB).InquirySellerAllocationSearchSeller(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminInquiryUpdateAttachments
// @Tags Marketplace-Inquiry
// @Summary CreateFromPayload inquiry collection
// @Description CreateFromPayload inquiry collection
// @Accept  json
// @Produce  json
// @Success 200 {object} models.InquiryCollection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/admin/inquiries/{inquiry_id}/attachments [put]
func AdminInquiryUpdateAttachments(c echo.Context) error {
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

// AdminApproveMultipleInquiryQuotations Internal admin approve quotation
// @Tags Admin-Inquiry
// @Summary Internal admin approve quotation
// @Description Internal admin approve quotation
// @Accept  json
// @Produce  json
// @Param data body models.ApproveMultipleInquiryQuotationsRequest true
// @Success 200 {string} Sent
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/inquiries/approve_multiple_quotations [post]
func AdminApproveMultipleInquiryQuotations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.ApproveMultipleInquiryQuotationsRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	req.JwtClaimsInfo = claims
	inquiries, err := repo.NewInquiryRepo(cc.App.DB).ApproveMultipleInquiryQuotations(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, iq := range inquiries {
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

	return cc.Success("Sent")
}
