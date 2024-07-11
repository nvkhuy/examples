package routes

import (
	"github.com/engineeringinflow/inflow-backend/pkg/middlewares"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/common"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupCommonRoutes(g *echo.Group) {
	var authorizedGroup = g.Group("", router.Middlewares.IsAuthorized())
	var authorizedWithAdminRoleGroup = authorizedGroup.Group("", router.Middlewares.CheckRole(enums.RoleSuperAdmin, enums.RoleLeader, enums.RoleStaff))
	var authorizedWithBuyerRoleGroup = authorizedGroup.Group("", router.Middlewares.CheckRole(enums.RoleClient))

	var authorizedRSA = g.Group("")

	g.GET("/attachments/:file_key", controllers.GetAttachment)
	g.GET("/blur_attachments/:file_key", controllers.GetBlurAttachment)
	g.GET("/thumbnail_attachments/:file_key", controllers.GetThumbnailAttachment)
	g.GET("/sitemap", controllers.GetSitemap)
	g.GET("/qrcode", controllers.GetQRCode)
	g.GET("/qrcode/download", controllers.DownloadQRCode)
	g.POST("/qrcode", controllers.CreateQRCodes)

	authorizedRSA.GET("/constants", controllers.GetConstants)
	authorizedRSA.GET("/register_constants", controllers.GetRegisterConstants)
	authorizedRSA.POST("/check_exists", controllers.CheckExists)
	authorizedRSA.GET("/seo/translations", controllers.GetSEOTranslations)
	authorizedRSA.GET("/seo/translations/by_route", controllers.GetSettingSEOByRouteName)
	authorizedRSA.GET("/buyer/seo/translations", controllers.GetBuyerSEOTranslations)
	authorizedRSA.GET("/seller/seo/translations", controllers.GetSellerSEOTranslations)
	authorizedRSA.GET("/admin/seo/translations", controllers.GetAdminSEOTranslations)
	authorizedRSA.GET("/addresses/search", controllers.SearchAddresses)
	authorizedRSA.GET("/settings/bank_infos", controllers.GetSettingBankInfos)

	authorizedRSA.POST("/upload/signatures", controllers.GetS3Signatures)

	authorizedGroup.GET("/download_link/:file_key", controllers.GetAttachmentDownloadLink)

	// S3
	authorizedGroup.POST("/s3_signatures", controllers.GetS3Signatures)
	authorizedGroup.POST("/s3_signature", controllers.GetS3Signature)

	authorizedGroup.GET("/search_addresses", controllers.SearchAddresses)

	authorizedGroup.GET("/settings/taxes", controllers.GetSettingTaxes)
	authorizedGroup.GET("/settings/sizes", controllers.GetSettingSizes)

	authorizedGroup.GET("/push_tokens", controllers.GetPushTokens)
	authorizedGroup.POST("/push_tokens", controllers.CreatePushToken)
	authorizedGroup.DELETE("/push_tokens/:token", controllers.DeletePushToken)

	authorizedGroup.GET("/docs", controllers.GetDocs)

	g.GET("/share/link/:link_id", controllers.GetShareLink)

	g.GET("/share/checkout_link/:link_id", controllers.GetCheckoutInfo, middlewares.IsAuthorizedWithQueryToken(router.App.Config.CheckoutJwtSecret))

	authorizedWithAdminRoleGroup.POST("/admin/share/link/:reference_id", controllers.CreateShareLink)
	authorizedWithBuyerRoleGroup.POST("/buyer/share/link/:reference_id", controllers.CreateShareLink)

	authorizedWithAdminRoleGroup.POST("/admin/share/checkout_link/:reference_id", controllers.CreateCheckoutShareLink)
	authorizedWithBuyerRoleGroup.POST("/buyer/share/checkout_link/:reference_id", controllers.CreateCheckoutShareLink)

	// Zalo
	g.POST("/zns/send", controllers.SendZNS)
}
