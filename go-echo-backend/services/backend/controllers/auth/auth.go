package controllers

import (
	"net/http"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/oauth"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"

	"github.com/rotisserie/eris"
)

// Register
// @Tags Auth
// @Summary Register
// @Description Register
// @Accept  json
// @Produce  json
// @Param data body models.RegisterForm true "Form"
// @Success 200 {object} models.LoginResponse
// @Failure 404 {object} errs.Error
// @Router /api/v1/register [post]
func Register(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var form models.RegisterForm
	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	response, err := repo.NewAuthRepo(cc.App.DB).Register(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, err = repo.NewAuthRepo(cc.App.DB).ResendVerificationEmail(repo.ResendVerificationEmailParams{
		UserID: response.User.ID,
		User:   response.User,
	})
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.OnboardUserTask{
		UserID: response.User.ID,
	}.DispatchIn(c.Request().Context(), time.Second*3)

	return cc.Success(response)

}

// LoginEmail Login by email
// @Tags Auth
// @Summary Login by email
// @Description Login by email
// @Accept  json
// @Produce  json
// @Param data body models.LoginEmailForm true "Login Form"
// @Success 200 {object} models.LoginResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/login_email [post]
func LoginEmail(c echo.Context) error {
	var form models.LoginEmailForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	response, err := repo.NewAuthRepo(cc.App.DB).LoginEmail(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if response.User.Role == enums.RoleClient && !response.User.AccountVerified {
		_, _ = repo.NewAuthRepo(cc.App.DB).ResendVerificationEmail(repo.ResendVerificationEmailParams{
			User:   response.User,
			UserID: response.User.ID,
		})
	}

	return cc.Success(response)

}

// ForgotPassword Forgot password
// @Tags Auth
// @Summary Forgot password
// @Description Forgot password
// @Accept  json
// @Produce  json
// @Param data body models.ForgotPasswordForm true "Login Form"
// @Success 200 {object} models.ForgotPasswordResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/forgot_password [post]
func ForgotPassword(c echo.Context) error {
	var form models.ForgotPasswordForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	response, err := repo.NewAuthRepo(cc.App.DB).ForgotPassword(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, err = tasks.TrackCustomerIOTask{
		UserID: response.User.ID,
		Event:  customerio.EventResetPassword,
		Data: map[string]interface{}{
			"email": form.Email,
			"link":  response.RedirectURL,
		},
	}.Dispatch(cc.Request().Context())
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(response)

}

// ResetPassword Reset password
// @Tags Auth
// @Summary Reset password
// @Description Reset password
// @Accept  json
// @Produce  json
// @Param data body models.ResetPasswordForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/reset_password [post]
func ResetPassword(c echo.Context) error {
	var form models.ResetPasswordForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, err = repo.NewAuthRepo(cc.App.DB).ResetPassword(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Password is changed")

}

// GoogleLogin Google login
// @Tags Auth
// @Summary Google login
// @Description Google login
// @Accept  json
// @Produce  json
// @Param data body models.ResetPasswordForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/oauth/google [get]
func GoogleLogin(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.GoogleOauthForm

	var err = cc.BindAndValidate(&form)
	if err != nil {
		return err
	}

	var url = oauth.GetInstance().AuthCodeURL(form.Platform.String(), oauth2.ApprovalForce)

	return cc.Redirect(http.StatusTemporaryRedirect, url)

}

// GoogleLoginCallback Google login callback
// @Tags Auth
// @Summary Google callback
// @Description Google callback
// @Accept  json
// @Produce  json
// @Param data body models.ResetPasswordForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/oauth/google/callback [get]
func GoogleLoginCallback(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var form repo.GetOauthUserInfoForm
	var err = cc.BindAndValidate(&form)
	if err != nil {
		return err
	}

	resp, err := repo.NewAuthRepo(cc.App.DB).GoogleLoginCallback(form)
	if err != nil {
		return err
	}

	return cc.Redirect(http.StatusPermanentRedirect, resp.RedirectURL)

}

// VerifyEmail VerifyEmail
// @Tags Auth
// @Summary Verify Email
// @Description Verify Email
// @Accept  json
// @Produce  json
// @Param data body models.VerifyEmailForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/verify_email [get]
func VerifyEmail(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var form models.VerifyEmailForm
	var err = cc.BindAndValidate(&form)
	if err != nil {
		var link = helper.AddURLQuery(form.RedirectURL, map[string]string{"error_message": err.Error()})
		return cc.Redirect(http.StatusPermanentRedirect, link)
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		var link = helper.AddURLQuery(form.RedirectURL, map[string]string{"error_message": err.Error()})
		return cc.Redirect(http.StatusPermanentRedirect, link)
	}

	form.JwtClaimsInfo = claims
	err = repo.NewAuthRepo(cc.App.DB).VerifyEmail(form)
	if err != nil {
		if eris.Is(err, errs.ErrAccountAlreadyVerified) {
			return cc.Redirect(http.StatusPermanentRedirect, form.RedirectURL)
		}
		var link = helper.AddURLQuery(form.RedirectURL, map[string]string{"error_message": err.Error()})
		return cc.Redirect(http.StatusPermanentRedirect, link)
	}

	tasks.ApproveUserTask{
		UserID: claims.GetUserID(),
	}.Dispatch(c.Request().Context())

	return cc.Redirect(http.StatusPermanentRedirect, form.RedirectURL)
}

// ResendVerificationEmail ResendVerificationEmail
// @Tags Auth
// @Summary Resend Verification Email
// @Description Resend Verification Email
// @Accept  json
// @Produce  json
// @Param data body models.ResendVerificationEmailForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/resend_verification_email [post]
func ResendVerificationEmail(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return err
	}

	result, err := repo.NewAuthRepo(cc.App.DB).ResendVerificationEmail(repo.ResendVerificationEmailParams{
		UserID: claims.GetUserID(),
	})
	if err != nil {
		return err
	}

	return cc.Success(result)
}

// // ShopifyOauthCallback Shopify oauth callback
// // @Tags Auth
// // @Summary Shopify oauth callback
// // @Description Shopify oauth callback
// // @Accept  json
// // @Produce  json
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/oauth/shopify/callback [get]
// func ShopifyOauthCallback(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)

// 	var params shopify.AuthorizeShopParams

// 	var err = cc.BindAndValidate(&params)
// 	if err != nil {
// 		return err
// 	}

// 	resp, err := repo.NewAuthRepo(cc.App.DB).AuthShopify(params)
// 	if err != nil {
// 		return err
// 	}

// 	tasks.ShopifyCreateWebhooksTask{
// 		ClientInfo: &shopify.ClientInfo{
// 			ShopName: resp.Shop.Domain,
// 			Token:    resp.Token,
// 		},
// 	}.Dispatch(cc.Request().Context())

// 	tasks.ShopifySyncChannelProductsTask{
// 		LocationID: fmt.Sprintf("%d", resp.Shop.PrimaryLocationId),
// 		ClientInfo: &shopify.ClientInfo{
// 			ShopName: resp.Shop.Domain,
// 			Token:    resp.Token,
// 		},
// 	}.Dispatch(cc.Request().Context())

// 	return cc.Success(resp)
// }

// AdminLoginEmail Login by email
// @Tags Auth
// @Summary Login by email
// @Description Login by email
// @Accept  json
// @Produce  json
// @Param data body models.LoginEmailForm true "Login Form"
// @Success 200 {object} models.LoginResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/login_email [post]
func AdminLoginEmail(c echo.Context) error {
	var form models.LoginEmailForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.IsAdminPortal = true
	response, err := repo.NewAuthRepo(cc.App.DB).LoginEmail(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(response)

}

// admin/resources

// AdminForgotPassword Forgot password
// @Tags Auth
// @Summary Forgot password
// @Description Forgot password
// @Accept  json
// @Produce  json
// @Param data body models.ForgotPasswordForm true "Login Form"
// @Success 200 {object} models.ForgotPasswordResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/forgot_password [post]
func AdminForgotPassword(c echo.Context) error {
	var form models.ForgotPasswordForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.IsAdminPortal = true
	response, err := repo.NewAuthRepo(cc.App.DB).ForgotPassword(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, err = tasks.TrackCustomerIOTask{
		UserID: response.User.ID,
		Event:  customerio.EventResetPassword,
		Data: map[string]interface{}{
			"email": form.Email,
			"link":  response.RedirectURL,
		},
	}.Dispatch(cc.Request().Context())
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(response)

}

// AdminResetPassword Reset password
// @Tags Auth
// @Summary Reset password
// @Description Reset password
// @Accept  json
// @Produce  json
// @Param data body models.ResetPasswordForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/reset_password [post]
func AdminResetPassword(c echo.Context) error {
	var form models.ResetPasswordForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.IsAdminPortal = true
	_, err = repo.NewAuthRepo(cc.App.DB).ResetPassword(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Password is changed")

}

// seller/resources

// SellerRegister
// @Tags Auth
// @Summary Register
// @Description Register
// @Accept  json
// @Produce  json
// @Param data body models.RegisterForm true "Form"
// @Success 200 {object} models.LoginResponse
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/register [post]
func SellerRegister(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var form models.RegisterForm
	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.IsSeller = true
	response, err := repo.NewAuthRepo(cc.App.DB).Register(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, _ = tasks.OnboardSellerTask{
		UserID: response.User.ID,
	}.DispatchIn(c.Request().Context(), time.Second*3)

	return cc.Success(response)

}

// SellerLoginEmail Login by email
// @Tags Auth
// @Summary Login by email
// @Description Login by email
// @Accept  json
// @Produce  json
// @Param data body models.LoginEmailForm true "Login Form"
// @Success 200 {object} models.LoginResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/login_email [post]
func SellerLoginEmail(c echo.Context) error {
	var form models.LoginEmailForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.IsSeller = true
	response, err := repo.NewAuthRepo(cc.App.DB).LoginEmail(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(response)

}

// AdminForgotPassword Forgot password
// @Tags Auth
// @Summary Forgot password
// @Description Forgot password
// @Accept  json
// @Produce  json
// @Param data body models.ForgotPasswordForm true "Login Form"
// @Success 200 {object} models.ForgotPasswordResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/forgot_password [post]
func SellerForgotPassword(c echo.Context) error {
	var form models.ForgotPasswordForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.IsSeller = true
	response, err := repo.NewAuthRepo(cc.App.DB).ForgotPassword(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, err = tasks.TrackCustomerIOTask{
		UserID: response.User.ID,
		Event:  customerio.EventResetPassword,
		Data: map[string]interface{}{
			"email": form.Email,
			"link":  response.RedirectURL,
		},
	}.Dispatch(cc.Request().Context())
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(response)

}

// SellerResetPassword Reset password
// @Tags Auth
// @Summary Reset password
// @Description Reset password
// @Accept  json
// @Produce  json
// @Param data body models.ResetPasswordForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/reset_password [post]
func SellerResetPassword(c echo.Context) error {
	var form models.ResetPasswordForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.IsSeller = true
	_, err = repo.NewAuthRepo(cc.App.DB).ResetPassword(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Password is changed")

}

// ZaloConnect Zalo connect
// @Tags Auth
// @Summary Reset password
// @Description Reset password
// @Accept  json
// @Produce  json
// @Param data body models.ResetPasswordForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/zalo/connect [post]
func ZaloConnect(c echo.Context) error {
	var params repo.ConnectZaloParams
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return err
	}
	params.JwtClaimsInfo = claims

	if err = cc.BindAndValidate(&params); err != nil {
		return err
	}

	err = repo.NewUserRepo(cc.App.DB).ConnectZalo(params)

	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success("connected")

}

// ZaloDisconnect Zalo connect
// @Tags Auth
// @Summary Reset password
// @Description Reset password
// @Accept  json
// @Produce  json
// @Param data body models.ResetPasswordForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/zalo/disconnect [post]
func ZaloDisconnect(c echo.Context) error {
	var params repo.DisconnectZaloParams
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return err
	}
	params.JwtClaimsInfo = claims

	if err = cc.BindAndValidate(&params); err != nil {
		return err
	}

	err = repo.NewUserRepo(cc.App.DB).DisconnectZalo(params)

	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success("disconnected")

}
