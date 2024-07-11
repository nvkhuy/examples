package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v74"
)

// StripeConfirmPurchaseOrder Stripe confirm purchase order
// @Tags Webhook
// @Summary Stripe confirm purchase order
// @Description Stripe confirm purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id param string true "Purchase Order ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/callback/stripe/payment_intents/inquiries/{inquiry_id}/purchase_orders/{purchase_order_id}/confirm [get]
func StripeConfirmPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PurchaseOrderPaymentIntentConfirmParams

	var redirectURL = fmt.Sprintf("%s/inquiries", cc.App.Config.BrandPortalBaseURL)

	var err = cc.BindAndValidate(&params)
	if err != nil {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": err.Error(),
		})
		return cc.Redirect(http.StatusPermanentRedirect, url)
	}
	redirectURL = fmt.Sprintf("%s/inquiries/%s", cc.App.Config.BrandPortalBaseURL, params.InquiryID)

	intent, err := stripehelper.GetInstance().GetPaymentIntent(params.PaymentIntent)
	if err != nil {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": err.Error(),
		})
		return cc.Redirect(http.StatusPermanentRedirect, url)
	}

	if intent.Status != stripe.PaymentIntentStatusSucceeded {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": "Payment failed, please try again",
		})

		return cc.Redirect(http.StatusPermanentRedirect, url)
	}

	return cc.Redirect(http.StatusPermanentRedirect, redirectURL)
}

// StripeConfirmPurchaseOrder Stripe confirm purchase order
// @Tags Webhook
// @Summary Stripe confirm purchase order
// @Description Stripe confirm purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id param string true "Purchase Order ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/callback/stripe/payment_intents/inquiries/{inquiry_id}/bulk_purchase_orders/{bulk_purchase_order_id}/confirm [post]
func StripeConfirmBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.BulkPurchaseOrderPaymentIntentConfirmParams

	var redirectURL = fmt.Sprintf("%s/bulks", cc.App.Config.BrandPortalBaseURL)

	var err = cc.BindAndValidate(&params)
	if err != nil {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": err.Error(),
		})
		return cc.Redirect(http.StatusPermanentRedirect, url)
	}
	redirectURL = fmt.Sprintf("%s/bulks/%s", cc.App.Config.BrandPortalBaseURL, params.InquiryID)

	intent, err := stripehelper.GetInstance().GetPaymentIntent(params.PaymentIntent)
	if err != nil {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": err.Error(),
		})
		return cc.Redirect(http.StatusPermanentRedirect, url)
	}

	if intent.Status != stripe.PaymentIntentStatusSucceeded {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": "Payment failed, please try again",
		})

		return cc.Redirect(http.StatusPermanentRedirect, url)
	}

	return cc.Redirect(http.StatusPermanentRedirect, redirectURL)
}

// StripeConfirmInquiryCartsCheckout Stripe confirm checkout session
// @Tags Webhook
// @Summary Stripe confirm purchase order
// @Description Stripe confirm purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id param string true "Purchase Order ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/callback/stripe/inquiry_carts/checkout/{checkout_session_id}/confirm [post]
func StripeConfirmInquiryCartsCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.StripeConfirmInquiryCartsCheckoutParams

	var redirectURL = fmt.Sprintf("%s/inquiry-checkout", cc.App.Config.BrandPortalBaseURL)

	var err = cc.BindAndValidate(&params)
	if err != nil {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": err.Error(),
			"cart_items":    params.CartItems,
		})
		return cc.Redirect(http.StatusPermanentRedirect, url)
	}

	intent, err := stripehelper.GetInstance().GetPaymentIntent(params.PaymentIntent)
	if err != nil {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": err.Error(),
			"cart_items":    params.CartItems,
		})
		return cc.Redirect(http.StatusPermanentRedirect, url)
	}

	if intent.Status != stripe.PaymentIntentStatusSucceeded {
		var url = helper.AddURLQuery(redirectURL, map[string]string{
			"error_message": "Payment failed, please try again",
			"cart_items":    params.CartItems,
		})

		return cc.Redirect(http.StatusPermanentRedirect, url)
	}

	var url = helper.AddURLQuery(redirectURL, map[string]string{
		"checkout_session_id": params.CheckoutSessionID,
	})
	return cc.Redirect(http.StatusPermanentRedirect, url)
}

// StripeConfirmOrderCartCheckout Stripe confirm checkout session
// @Tags Webhook
// @Summary Stripe confirm purchase order
// @Description Stripe confirm purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id param string true "Purchase Order ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/callback/stripe/order_cart/checkout/{checkout_session_id}/confirm [get]
func StripeConfirmOrderCartsCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.StripeConfirmOrderCartCheckoutParams

	var redirectURL, err = url.Parse(fmt.Sprintf("%s/order-checkout", cc.App.Config.BrandPortalBaseURL))

	if err != nil {
		return err
	}
	urlParams := redirectURL.Query()

	if err := cc.BindAndValidate(&params); err != nil {
		urlParams.Add("error_message", err.Error())
		redirectURL.RawQuery = urlParams.Encode()
		return cc.Redirect(http.StatusPermanentRedirect, redirectURL.String())
	}

	var poItemIDs, bulkIDs = []string{}, []string{}
	if params.PurchaseOrderCartItemIDs != "" {
		poItemIDs = strings.Split(params.PurchaseOrderCartItemIDs, ",")
	}
	if params.BulkOrderIDs != "" {
		bulkIDs = strings.Split(params.BulkOrderIDs, ",")
	}

	poItemIDsString, err := json.Marshal(poItemIDs)
	if err != nil {
		poItemIDsString = []byte("")
	}
	bulkIDsString, err := json.Marshal(bulkIDs)
	if err != nil {
		bulkIDsString = []byte("")
	}

	intent, err := stripehelper.GetInstance().GetPaymentIntent(params.PaymentIntent)
	if err != nil {
		urlParams.Add("error_message", err.Error())
		urlParams.Add("po_items", string(poItemIDsString))
		urlParams.Add("bulks", string(bulkIDsString))
		redirectURL.RawQuery = urlParams.Encode()
		return cc.Redirect(http.StatusPermanentRedirect, redirectURL.String())
	}

	if intent.Status != stripe.PaymentIntentStatusSucceeded {
		urlParams.Add("error_message", "Payment failed, please try again")
		urlParams.Add("po_items", string(poItemIDsString))
		urlParams.Add("bulks", string(bulkIDsString))
		redirectURL.RawQuery = urlParams.Encode()

		return cc.Redirect(http.StatusPermanentRedirect, redirectURL.String())
	}

	urlParams.Add("checkout_session_id", params.CheckoutSessionID)
	redirectURL.RawQuery = urlParams.Encode()
	return cc.Redirect(http.StatusPermanentRedirect, redirectURL.String())
}
