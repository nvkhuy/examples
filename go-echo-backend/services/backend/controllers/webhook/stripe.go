package webhook

import (
	"encoding/json"
	"io"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
)

// ShopifyProductCreate Shopify product create
// @Tags Webhook
// @Summary Shopify product create
// @Description Shopify product create
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword" default(1)
// @Param page query int false "Page index" default(1)
// @Param limit query int false "Size of page" default(20)
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/webhook/stripe [post]
func StripeWebhook(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	data, err := io.ReadAll(cc.Request().Body)
	if err != nil {
		return err
	}
	event, err := webhook.ConstructEvent(data, c.Request().Header.Get("Stripe-Signature"), cc.App.Config.StripeWebhookSecretKey)
	if err != nil {
		return err
	}

	var handler = StripeHandler{}

	switch event.Type {
	case string(stripehelper.PaymentIntentSucceeded):
		var object stripe.PaymentIntent
		err = json.Unmarshal(event.Data.Raw, &object)
		if err != nil {
			cc.CustomLogger.Errorf("Error parsing webhook JSON: %v", err)
			return err
		}

		err = handler.HandlePayment(c, &HandlePaymentParams{
			PaymentIntent: &object,
		})

	case string(stripehelper.CheckoutSessionCompleted):
		var object stripe.CheckoutSession
		err = json.Unmarshal(event.Data.Raw, &object)
		if err != nil {
			cc.CustomLogger.Errorf("Error parsing webhook JSON: %v", err)
			return err
		}

		err = handler.HandlePayment(c, &HandlePaymentParams{
			CheckoutSession: &object,
		})

	}

	if err != nil {
		cc.CustomLogger.Errorf("Handle payment error: %v", err)
		return err
	}

	return cc.Success("Success")
}
