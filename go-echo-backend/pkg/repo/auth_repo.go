package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/oauth"
	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
	"github.com/jinzhu/copier"

	"github.com/engineeringinflow/inflow-backend/pkg/mailer"
	"github.com/thaitanloi365/go-sendgrid"

	"github.com/thaitanloi365/go-utils/values"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepo struct {
	db *db.DB
}

func NewAuthRepo(db *db.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) LoginEmail(form models.LoginEmailForm) (*models.LoginResponse, error) {
	var user models.User
	var err error

	if form.IsAdminPortal {
		var adminPortalRoles = []enums.Role{
			enums.RoleSuperAdmin,
			enums.RoleLeader,
			enums.RoleStaff,
		}
		err = r.db.First(&user, "email = ? AND role IN ?", form.Email, adminPortalRoles).Error
	} else if form.IsSeller {
		err = r.db.First(&user, "email = ? AND role = ?", form.Email, enums.RoleSeller).Error
	} else {
		err = r.db.First(&user, "email = ? AND role = ?", form.Email, enums.RoleClient).Error
	}

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrLoginInvalid
		}

		return nil, eris.Wrap(err, err.Error())
	}

	err = user.ComparePassword(form.Password)
	if err != nil {
		return nil, errs.ErrPasswordIncorrect
	}

	token, err := user.GenerateToken(r.db.Configuration.JWTSecret)
	if err != nil {
		r.db.CustomLogger.ErrorAny(err)
		return nil, eris.Wrap(err, err.Error())
	}

	var lastLogin = r.db.NowFunc().Unix()
	var update = models.User{
		LastLogin:   &lastLogin,
		LoggedOutAt: nil,
		TokenIssuer: user.TokenIssuer,
	}

	err = r.db.Model(&user).Updates(&update).Error
	if err != nil {
		r.db.CustomLogger.ErrorAny(err)
		return nil, eris.Wrap(err, err.Error())
	}

	if user.DeletedAt != nil && user.DeletedAt.Int64 > 0 {
		return nil, errs.ErrAccountInactive
	}

	var response = models.LoginResponse{User: &user, Token: token}

	return &response, nil
}

func (r *AuthRepo) ForgotPassword(form models.ForgotPasswordForm) (*models.ForgotPasswordResponse, error) {
	var user models.User
	var err error

	if form.IsAdminPortal {
		var adminPortalRoles = []enums.Role{
			enums.RoleSuperAdmin,
			enums.RoleLeader,
			enums.RoleStaff,
		}
		err = r.db.Select("ID", "Email", "Role", "TokenIssuer", "TokenResetPassword", "TokenResetPasswordSentAt").
			First(&user, "email = ? AND role IN ?", form.Email, adminPortalRoles).Error
	} else if form.IsSeller {
		err = r.db.Select("ID", "Email", "Role", "TokenIssuer", "TokenResetPassword", "TokenResetPasswordSentAt").
			First(&user, "email = ? AND role = ?", form.Email, enums.RoleSeller).Error
	} else {
		err = r.db.Select("ID", "Email", "Role", "TokenIssuer", "TokenResetPassword", "TokenResetPasswordSentAt").
			First(&user, "email = ? AND role = ?", form.Email, enums.RoleClient).Error
	}
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	if user.TokenResetPasswordSentAt != nil {
		var diff = int(time.Now().Unix() - *user.TokenResetPasswordSentAt)
		var nextIn = int(r.db.Configuration.ResetPasswordResendInterval.Seconds())
		if diff < nextIn {
			var resp = &models.ForgotPasswordResponse{
				Email:         user.Email,
				User:          &user,
				NextInSeconds: nextIn - diff,
			}
			return resp.TransformMessage(), nil
		}

	}

	token, err := user.GenerateResetPasswordTokenAndUpdate(r.db)
	if err != nil {
		return nil, eris.Wrap(err, "GenerateResetPasswordTokenAndUpdate Error")
	}

	redirectURL, err := url.Parse(form.RedirectURL)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	var values = redirectURL.Query()

	values.Add("token", token)
	values.Add("email", user.Email)

	redirectURL.RawQuery = values.Encode()

	var resp = &models.ForgotPasswordResponse{
		Email:       user.Email,
		RedirectURL: redirectURL.String(),
		User:        &user,
	}
	return resp, nil
}

func (r *AuthRepo) ResetPassword(form models.ResetPasswordForm) (*models.User, error) {
	var user models.User
	var claims models.JwtClaims

	var err = claims.ValidateToken(r.db.Configuration.JWTResetPasswordSecret, form.TokenResetPassword)
	if err != nil {
		return &user, errs.ErrTokenInvalid
	}

	if form.IsAdminPortal {
		var adminPortalRoles = []enums.Role{
			enums.RoleSuperAdmin,
			enums.RoleLeader,
			enums.RoleStaff,
		}
		err = r.db.Select("TokenResetPassword").First(&user, "id = ? AND role IN ?", claims.ID, adminPortalRoles).Error
	} else if form.IsSeller {
		err = r.db.Select("TokenResetPassword").First(&user, "id = ? AND role = ?", claims.ID, enums.RoleSeller).Error
	} else {
		err = r.db.Select("TokenResetPassword").First(&user, "id = ? AND role = ?", claims.ID, enums.RoleClient).Error

	}
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return &user, errs.ErrUserNotFound
		}
		return &user, err
	}

	if user.TokenResetPassword == nil || *user.TokenResetPassword != form.TokenResetPassword {
		return &user, errs.ErrTokenInvalid
	}

	user.Password = form.NewPassword
	err = user.HashPassword()
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	var tokenResetPassword = ""
	var updates = models.User{
		Password:           user.Password,
		TokenResetPassword: &tokenResetPassword,
		AccountVerified:    true,
		AccountVerifiedAt:  values.Int64(time.Now().Unix()),
	}

	err = r.db.Model(&models.User{}).Where("id = ?", claims.ID).Updates(&updates).Error
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if form.IsAdminPortal {
		user.ID = claims.ID
		user.AccountStatus = enums.AccountStatusActive
		err = user.UpdateAccountStatus(r.db)
		if err != nil {
			return nil, eris.Wrap(err, "UpdateAccountStatus Error")
		}
	}

	return &user, nil
}

type GetOauthUserInfoForm struct {
	State string `json:"state" query:"state" form:"state"`
	Code  string `json:"code" query:"code" form:"code" validate:"required"`
}

type GetOauthUserInfoResponse struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	Gender        string `json:"gender"`
	GivenName     string `json:"given_name"`
	Hd            string `json:"hd"`
	ID            string `json:"id"`
	Link          string `json:"link"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (r *AuthRepo) GetOauthUserInfo(form GetOauthUserInfoForm) (*GetOauthUserInfoResponse, error) {
	token, err := oauth.GetInstance().Exchange(context.Background(), form.Code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	var resp GetOauthUserInfoResponse

	err = json.NewDecoder(response.Body).Decode(&resp)

	return &resp, err
}

func (r *AuthRepo) Register(form models.RegisterForm) (*models.LoginResponse, error) {
	var user models.User
	var err = copier.Copy(&user, &form)
	if err != nil {
		return nil, err
	}

	if form.IsSeller {
		user.Role = enums.RoleSeller
	} else {
		user.Role = enums.RoleClient
	}

	if user.Coordinate != nil {
		err = user.Coordinate.CreateOrUpdate(r.db)
		if err != nil {
			return nil, err
		}
		user.CoordinateID = user.Coordinate.ID
	}

	user.Name = user.GetFullName()
	user.ID = helper.GenerateXID()
	token, err := user.GenerateToken(r.db.Configuration.JWTSecret)
	if err != nil {
		return nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		user.OnboardingSubmitAt = values.Int64(r.db.NowFunc().Unix())
		err = tx.Create(&user).Error
		if err != nil {
			if e := r.db.CheckUserDuplicateConstraint(err); e != nil {
				return e
			}
			return eris.Wrap(err, err.Error())
		}

		if form.BusinessProfile != nil {
			form.BusinessProfile.User = &user
			_, err = NewUserRepo(r.db).OnboardingSubmit(tx, *form.BusinessProfile)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	var response = models.LoginResponse{User: &user, Token: token}

	return &response, err
}

type GoogleLoginCallbackResponse struct {
	RedirectURL string
	IsNew       bool         `json:"-"`
	User        *models.User `json:"-"`
}

func (r *AuthRepo) GoogleLoginCallback(form GetOauthUserInfoForm) (*GoogleLoginCallbackResponse, error) {
	var loginURL = fmt.Sprintf("%s/login", r.db.Configuration.BrandPortalBaseURL)

	userInfo, err := r.GetOauthUserInfo(form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var user models.User

	switch form.State {
	case "admin":
		var adminPortalRoles = []enums.Role{
			enums.RoleSuperAdmin,
			enums.RoleLeader,
			enums.RoleStaff,
		}
		err = r.db.First(&user, "email = ? AND role IN ?", userInfo.Email, adminPortalRoles).Error
		loginURL = fmt.Sprintf("%s/login", r.db.Configuration.AdminPortalBaseURL)
	case "seller":
		err = r.db.First(&user, "email = ? AND role = ?", userInfo.Email, enums.RoleSeller).Error
		loginURL = fmt.Sprintf("%s/login", r.db.Configuration.SellerPortalBaseURL)
	case "client":
		err = r.db.First(&user, "email = ? AND role = ?", userInfo.Email, enums.RoleClient).Error

	default:
		return nil, errs.New(422, "Invalid request")
	}

	redirectURL, err1 := url.Parse(loginURL)
	if err1 != nil {
		return nil, err1
	}

	var resp = GoogleLoginCallbackResponse{
		RedirectURL: loginURL,
	}

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			if form.State == "admin" {
				var params = redirectURL.Query()
				params.Add("error_message", errs.ErrUserNotFound.Message)
				redirectURL.RawQuery = params.Encode()

				resp.IsNew = true
				resp.User = &user
				resp.RedirectURL = redirectURL.String()
				return &resp, nil
			}
			user = models.User{
				Email:             userInfo.Email,
				FirstName:         userInfo.GivenName,
				LastName:          userInfo.FamilyName,
				Name:              userInfo.Name,
				SocialProvider:    "google",
				Role:              enums.RoleClient,
				AccountVerified:   true,
				AccountVerifiedAt: values.Int64(time.Now().Unix()),
				AccountStatus:     enums.AccountStatusActive,
			}
			user.ID = helper.GenerateXID()
			token, err := user.GenerateToken(r.db.Configuration.JWTSecret)
			if err != nil {
				return nil, eris.Wrap(err, err.Error())
			}

			err = r.db.Clauses(clause.OnConflict{
				DoNothing: true,
			}).
				Create(&user).Error
			if err != nil {
				return nil, eris.Wrap(err, err.Error())
			}

			var params = redirectURL.Query()
			params.Add("token", token)
			params.Add("email", userInfo.Email)
			redirectURL.RawQuery = params.Encode()

			resp.IsNew = true
			resp.User = &user
			resp.RedirectURL = redirectURL.String()
			return &resp, nil
		}
		return nil, err
	} else {
		var updates = models.User{
			AccountVerified:   true,
			AccountVerifiedAt: values.Int64(time.Now().Unix()),
			AccountStatus:     enums.AccountStatusActive,
		}

		err = r.db.Model(&models.User{}).Where("email = ?", userInfo.Email).Updates(&updates).Error
		if err != nil {
			return nil, err
		}
	}

	token, err := user.GenerateToken(r.db.Configuration.JWTSecret)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var params = redirectURL.Query()
	params.Add("token", token)
	params.Add("email", userInfo.Email)
	redirectURL.RawQuery = params.Encode()

	resp.IsNew = false
	resp.User = &user
	resp.RedirectURL = redirectURL.String()

	return &resp, nil
}

func (r *AuthRepo) VerifyEmail(form models.VerifyEmailForm) error {
	var user models.User

	var err = r.db.Select("ID", "VerificationToken", "AccountVerified").First(&user, "id = ?", form.GetUserID()).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return errs.ErrUserNotFound
		}
		return err
	}

	if user.VerificationToken == nil || *user.VerificationToken != form.Token {
		return errs.ErrTokenInvalid
	}

	if user.AccountVerified {
		return errs.ErrAccountAlreadyVerified
	}

	var verificationToken = ""
	var updates = models.User{
		VerificationToken: &verificationToken,
		AccountVerified:   true,
		AccountVerifiedAt: values.Int64(time.Now().Unix()),
		AccountStatus:     enums.AccountStatusActive,
	}

	err = r.db.Model(&models.User{}).Where("id = ?", form.GetUserID()).Updates(&updates).Error
	if err != nil {
		return nil
	}

	return nil
}

func (r *AuthRepo) SendVerificationEmail(email string, token string) error {
	var webURL = fmt.Sprintf("%s/verify-email", r.db.Configuration.BrandPortalBaseURL)
	redirectURL, err := url.Parse(webURL)
	if err != nil {
		return err
	}

	var values = redirectURL.Query()
	values.Add("token", token)
	values.Add("email", email)
	redirectURL.RawQuery = values.Encode()

	var mailerInstance = mailer.GetInstance()
	err = mailerInstance.Send(sendgrid.SendMailParams{
		Email:      email,
		TemplateID: string(mailer.TemplateIDConfirmMail),
		Data: map[string]interface{}{
			"link": redirectURL.String(),
		},
	})

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return nil
}

func (r *AuthRepo) SendWelcomeBoardEmail(email string) error {
	var loginURL = fmt.Sprintf("%s/login", r.db.Configuration.WebAppBaseURL)

	var mailerInstance = mailer.GetInstance()
	var err = mailerInstance.Send(sendgrid.SendMailParams{
		Email:      email,
		TemplateID: string(mailer.TemplateIDWelcomToBoard),
		Data: map[string]interface{}{
			"login_url": loginURL,
		},
	})

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return nil
}

type ResendVerificationEmailParams struct {
	UserID string       `json:"-"`
	User   *models.User `json:"-"`
}

func (r *AuthRepo) ResendVerificationEmail(params ResendVerificationEmailParams) (*models.ResendVerificationEmailResponse, error) {
	var user = params.User

	if user == nil {
		user = new(models.User)
		var err = r.db.First(user, "id = ?", params.UserID).Error
		if err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrUserNotFound
			}

			return nil, eris.Wrap(err, err.Error())
		}
	}

	if user.AccountVerified {
		return nil, errs.ErrAccountAlreadyVerified
	}

	if user.VerificationSentAt != nil {
		var diff = int(time.Now().Unix() - *user.VerificationSentAt)
		var nextIn = int(r.db.Configuration.ResetPasswordResendInterval.Seconds())
		if diff < nextIn {
			var resp = &models.ResendVerificationEmailResponse{
				NextInSeconds: nextIn - diff,
			}
			return resp.TransformMessage(), nil
		}

	}

	token, err := user.GenerateVerificationToken(r.db)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	var webURL = fmt.Sprintf("%s/api/v1/verify_email", r.db.Configuration.ServerBaseURL)
	var redirectURL = helper.AddURLQuery(webURL, map[string]string{
		"token":        token,
		"redirect_url": r.db.Configuration.BrandPortalBaseURL,
	})

	err = customerio.GetInstance().Track.Track(user.ID, string(customerio.EventEmailVerification), map[string]interface{}{
		"email": user.Email,
		"name":  user.Name,
		"link":  redirectURL,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return &models.ResendVerificationEmailResponse{}, nil
}

// func (r *AuthRepo) AuthShopify(params shopify.AuthorizeShopParams) (*shopify.AuthorizeShopResponse, error) {
// 	resp, err := shopify.GetInstance().AuthorizeShop(params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var user models.User
// 	err = r.adb.First(&user, "email = ?", resp.Shop.Email).Error
// 	if err != nil {
// 		if r.adb.IsRecordNotFoundError(err) {
// 			user.ID = helper.GenerateXID()
// 			user.Name = resp.Shop.ShopOwner
// 			user.Email = resp.Shop.Email
// 			user.Role = enums.RoleClient
// 			user.CustomerType = enums.CustomerTypeSeller
// 			user.Password = "1234qwerR!"
// 			err = r.adb.CreateFromPayload(&user).Error
// 			if err != nil {
// 				return nil, err
// 			}
// 			goto Next
// 		}

// 		return nil, err
// 	}

// Next:

// 	var shop = models.Shop{
// 		Name:          resp.Shop.Name,
// 		UserID:        user.ID,
// 		AddressLevel1: resp.Shop.Address1,
// 		AddressLevel2: resp.Shop.Address2,
// 		Status:        enums.AccountStatusPendingReview,
// 	}

// 	var existingShop models.Shop
// 	var sqlResult = r.adb.
// 		Model(&existingShop).
// 		Clauses(clause.Returning{
// 			Columns: []clause.Column{
// 				{Name: "id"},
// 			},
// 		}).
// 		Where("user_id = ?", shop.UserID).
// 		Updates(&shop)
// 	if sqlResult.RowsAffected == 0 {
// 		shop.ID = helper.GenerateXID()
// 		err = r.adb.CreateFromPayload(&shop).Error
// 		if err != nil {
// 			return nil, err
// 		}

// 		if sqlResult.Error != nil {
// 			return nil, sqlResult.Error
// 		}
// 	} else {
// 		shop.ID = existingShop.ID
// 	}

// 	var shopChannel = models.ShopChannel{
// 		UserID:                  user.ID,
// 		ShopName:                resp.Shop.Domain,
// 		Token:                   resp.Token,
// 		SourcePrimaryLocationID: fmt.Sprintf("%d", resp.Shop.PrimaryLocationId),
// 		Channel:                 enums.ShopChannelShopify,
// 		ShopID:                  shop.ID,
// 	}
// 	shopChannel.ID = helper.GenerateXID()

// 	sqlResult = r.adb.Model(&models.ShopChannel{}).Where("user_id = ? AND shop_name = ?", shopChannel.UserID, shopChannel.ShopName).Updates(&shopChannel)
// 	if sqlResult.RowsAffected == 0 {
// 		err = r.adb.CreateFromPayload(&shopChannel).Error
// 		if err != nil {
// 			return nil, err
// 		}

// 		if sqlResult.Error != nil {
// 			return nil, sqlResult.Error
// 		}
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	if locations, err := shopify.GetInstance().NewClient(resp.Shop.Domain, resp.Token).Location.List(nil); err == nil {
// 		var activeLocations = lo.Filter(locations, func(item goshopify.Location, index int) bool {
// 			return item.Active
// 		})

// 		var locations []*models.ShopChannelLocation
// 		for _, activeLocation := range activeLocations {
// 			locations = append(locations, &models.ShopChannelLocation{
// 				ID:            fmt.Sprintf("%d", activeLocation.ID),
// 				CreatedAt:     activeLocation.CreatedAt.Unix(),
// 				UpdatedAt:     activeLocation.UpdatedAt.Unix(),
// 				ShopChannelID: shopChannel.ID,
// 				UserID:        user.ID,
// 				Name:          activeLocation.Name,
// 			})
// 		}

// 		if len(locations) > 0 {
// 			r.adb.Clauses(clause.OnConflict{UpdateAll: true}).CreateFromPayload(&locations)
// 		}
// 	}

// 	return resp, err
// }

func (r *AuthRepo) AuthorizeUrl(shopName string) string {
	return shopify.GetInstance().GetAuthURL(shopName, nil)

}
