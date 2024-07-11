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

// BuyerPaginateCatalogCart
// @Tags Marketplace-Catalog-Cart
// @Summary Preview checkout inquiry
// @Description Preview checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.CatalogCartCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/catalog_carts [get]
func BuyerPaginateCatalogCarts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCatalogCartsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewCatalogCartRepo(cc.App.DB).PaginateCatalogCarts(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerUpdateCatalogCarts
// @Tags Marketplace-Catalog-Cart
// @Summary Preview checkout inquiry
// @Description Preview checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.CatalogCartCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/catalog_carts [put]
func BuyerUpdateCatalogCarts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateCatalogCartsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewCatalogCartRepo(cc.App.DB).UpdateCatalogCarts(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerCreateCatalogCartsOrder
// @Tags Marketplace-Catalog-Cart
// @Summary Preview checkout inquiry
// @Description Preview checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.CatalogCartCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/catalog_carts [put]
func BuyerCreateCatalogCartsOrders(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateCatalogCartOrdersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewCatalogCartRepo(cc.App.DB).CreateCatalogCartOrders(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerMultiCatalogCartCheckout
// @Tags Marketplace-Inquiry
// @Summary Checkout inquiry
// @Description Checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.CatalogCartCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/catalog_carts/checkout [post]
func BuyerMultiCatalogCartCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MultiCatalogCartCheckoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewCatalogCartRepo(cc.App.DB).MultiCatalogCartCheckout(params)
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
					AdminID:   claims.GetUserID(),
					InquiryID: purchaseOrder.InquiryID,
					InquiryReferenceID: func() string {
						if purchaseOrder.Inquiry != nil {
							return purchaseOrder.Inquiry.ReferenceID
						}
						return ""
					}(),
					PurchaseOrderID:          purchaseOrder.ID,
					PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
				},
			}.Dispatch(c.Request().Context())
		}

		if params.PaymentType == enums.PaymentTypeBankTransfer {
			if len(result.Orders) == 1 {
				var purchaseOrder = result.Orders[0]
				var assigneeIDs = purchaseOrder.AssigneeIDs
				if purchaseOrder.Inquiry != nil {
					assigneeIDs = append(assigneeIDs, purchaseOrder.Inquiry.AssigneeIDs...)
				}

				for _, assigneeID := range assigneeIDs {
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

// BuyerMultiCatalogCartCheckout
// @Tags Marketplace-Inquiry
// @Summary Checkout inquiry
// @Description Checkout inquiry
// @Accept  json
// @Produce  json
// @Param data body repo.CatalogCartCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/catalog_carts/checkout_info [get]
func BuyerMultiCatalogCartCheckoutInfo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MultiCatalogCartCheckoutInfoParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewCatalogCartRepo(cc.App.DB).MultiCatalogCartCheckoutInfo(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
