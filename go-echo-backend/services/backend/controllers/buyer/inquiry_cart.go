package controllers

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
)

// BuyerMultiInquiryPreviewCheckout
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
// @Router /api/v1/inquiry_carts/preview_checkout [post]
func BuyerMultiInquiryPreviewCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MultiInquiryPreviewCheckoutParams

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
	result, err := repo.NewInquiryRepo(cc.App.DB).MultiInquiryPreviewCheckout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerMultiInquiryCheckout
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
// @Router /api/v1/inquiry_carts/checkout [post]
func BuyerMultiInquiryCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MultiInquiryCheckoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).MultiInquiryCheckout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if len(result.Orders) > 0 {
		for _, purchaseOrder := range result.Orders {
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
		}

		if params.PaymentType == enums.PaymentTypeBankTransfer {
			if len(result.Orders) == 1 {
				var purchaseOrder = result.Orders[0]
				for _, assigneeID := range purchaseOrder.Inquiry.AssigneeIDs {
					tasks.TrackCustomerIOTask{
						Event:  customerio.EventPoWaitingConfirmBankTransfer,
						UserID: assigneeID,
						Data:   purchaseOrder.GetCustomerIOMetadata(nil),
					}.Dispatch(c.Request().Context())
				}
			} else {
				var assigneeIDs []string
				var metadata []map[string]interface{}
				for _, order := range result.Orders {
					assigneeIDs = append(assigneeIDs, order.AssigneeIDs...)
					metadata = append(metadata, order.GetCustomerIOMetadata(nil))
				}

				for _, assigneeID := range lo.Uniq(assigneeIDs) {
					tasks.TrackCustomerIOTask{
						Event:  customerio.EventPoMultipleItemsWaitingConfirmBankTransfer,
						UserID: assigneeID,
						Data: result.PaymentTransaction.GetCustomerIOMetadata(map[string]interface{}{
							"purchase_orders": metadata,
						}),
					}.Dispatch(c.Request().Context())
				}
			}
		}

	}

	return cc.Success(result)
}

// BuyerMultiInquiryCheckout
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
// @Router /api/v1/inquiry_carts/checkout_info [get]
func BuyerMultiInquiryCheckoutInfo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MultiInquiryCheckoutInfoParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInquiryRepo(cc.App.DB).MultiInquiryCheckoutInfo(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
