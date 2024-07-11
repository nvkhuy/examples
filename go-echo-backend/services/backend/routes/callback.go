package routes

import (
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/callback"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupCallbackRoutes(g *echo.Group) {
	g.GET("/stripe/payment_intents/inquiries/:inquiry_id/purchase_orders/:purchase_order_id/confirm", controllers.StripeConfirmPurchaseOrder)
	g.GET("/stripe/payment_intents/inquiries/:inquiry_id/bulk_purchase_orders/:bulk_purchase_order_id/confirm", controllers.StripeConfirmBulkPurchaseOrder)
	g.GET("/stripe/payment_intents/inquiry_carts/:checkout_session_id/confirm", controllers.StripeConfirmInquiryCartsCheckout)
	g.GET("/stripe/payment_intents/order_cart/:checkout_session_id/confirm", controllers.StripeConfirmOrderCartsCheckout)
	g.GET("/zalo/permission", controllers.CallbackZaloOA)
}
