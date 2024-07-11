package routes

import (
	"github.com/engineeringinflow/inflow-backend/pkg/middlewares"
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/auth"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupAuthRoutes(g *echo.Group) {
	var authorizedGroup = g.Group("", router.Middlewares.IsAuthorized())

	authorizedGroup.POST("/zalo/connect", controllers.ZaloConnect)
	authorizedGroup.POST("/zalo/disconnect", controllers.ZaloDisconnect)

	authorizedGroup.POST("/resend_verification_email", controllers.ResendVerificationEmail)
	g.GET("/verify_email", controllers.VerifyEmail, middlewares.IsAuthorizedWithQueryToken(router.App.Config.JWTEmailVerificationSecret))

	g.POST("/login_email", controllers.LoginEmail)

	g.POST("/forgot_password", controllers.ForgotPassword)
	g.POST("/reset_password", controllers.ResetPassword)
	g.POST("/register", controllers.Register)

	g.GET("/oauth/google", controllers.GoogleLogin)
	g.GET("/oauth/google/callback", controllers.GoogleLoginCallback)

	// g.GET("/oauth/shopify/callback", controllers.ShopifyOauthCallback)

	g.POST("/admin/login_email", controllers.AdminLoginEmail)
	g.POST("/admin/forgot_password", controllers.AdminForgotPassword)
	g.POST("/admin/reset_password", controllers.AdminResetPassword)

	g.POST("/seller/register", controllers.SellerRegister)
	g.POST("/seller/login_email", controllers.SellerLoginEmail)
	g.POST("/seller/forgot_password", controllers.SellerForgotPassword)
	g.POST("/seller/reset_password", controllers.SellerResetPassword)
}
