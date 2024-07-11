package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/ai"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v74"

	"github.com/rotisserie/eris"
	"github.com/thaitanloi365/go-utils/values"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewUserRepo(db *db.DB) *UserRepo {
	return &UserRepo{
		db:     db,
		logger: logger.New("repo/user"),
	}
}

type PaginateUsersParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	CompanyID string       `json:"company_id" query:"company_id" form:"company_id"`
	Domain    string       `json:"domain" query:"domain" form:"domain"`
	Roles     []enums.Role `json:"roles" query:"roles" form:"roles"`
	UserIDs   []string     `json:"user_ids,omitempty" query:"user_ids" form:"user_ids"`
	Teams     []string     `json:"teams,omitempty" query:"teams" form:"teams"`

	AccountStatuses []enums.AccountStatus `json:"account_statuses" query:"account_statuses" form:"account_statuses"`

	IncludeAssignedInquiryIds bool `json:"include_assigned_inquiry_ids" query:"include_assigned_inquiry_ids" form:"include_assigned_inquiry_ids"`
	IncludeContactOwners      bool `json:"-"`
	IncludeBrandTeam          bool `json:"-"`
}

func (r *UserRepo) PaginateRecentUser(params PaginateUsersParams) *query.Pagination {
	var builder = queryfunc.NewUserBuilder(queryfunc.UserBuilderOptions{
		IncludeContactOwners: params.IncludeContactOwners,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("u.role = ?", enums.RoleClient)

			if len(params.AccountStatuses) > 0 {
				builder.Where("u.account_status IN ?", params.AccountStatuses)
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("( u.email ILIKE @query OR u.name ILIKE @query OR u.first_name ILIKE @query OR u.last_name ILIKE @query )", sql.Named("query", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *UserRepo) PaginateUsers(params PaginateUsersParams) *query.Pagination {
	var builder = queryfunc.NewUserBuilder(queryfunc.UserBuilderOptions{
		IncludeContactOwners:      params.IncludeContactOwners,
		IncludeBrandTeam:          params.IncludeBrandTeam,
		IncludeAssignedInquiryIds: params.IncludeAssignedInquiryIds,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {

			if len(params.Teams) > 0 {
				builder.Where("u.team IN ?", params.Teams)
			}

			if len(params.Roles) > 0 {
				builder.Where("u.role IN ?", params.Roles)
			}

			if len(params.UserIDs) > 0 {
				builder.Where("u.id IN ?", params.UserIDs)
			}

			if len(params.AccountStatuses) > 0 {
				builder.Where("u.account_status IN ?", params.AccountStatuses)
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("( u.email ILIKE @query OR u.name ILIKE @query OR u.first_name ILIKE @query OR u.last_name ILIKE @query )", sql.Named("query", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type SearchUsersParams struct {
	models.PaginationParams

	Roles []string `json:"roles" query:"roles" form:"roles"`

	AccountStatuses []enums.AccountStatus `json:"account_statuses" query:"account_statuses" form:"account_statuses"`

	ForRole enums.Role
}

func (r *UserRepo) SearchUsers(params PaginateUsersParams) []*models.User {
	var builder = queryfunc.NewUserBuilder(queryfunc.UserBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeAddress: true,
	})
	if params.Limit == 0 {
		params.Limit = 20
	}
	var result []*models.User
	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("u.deleted_at IS NULL")

			if len(params.Roles) > 0 {
				builder.Where("u.role IN ?", params.Roles)
			}

			if len(params.AccountStatuses) > 0 {
				builder.Where("u.account_status IN ?", params.AccountStatuses)
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("( u.email ILIKE @query OR u.name ILIKE @query )", sql.Named("query", q))
			}
		}).
		Limit(params.Limit).
		FindFunc(&result)

	return result
}

type GetUserParams struct {
	models.JwtClaimsInfo

	UserID string `param:"user_id" validate:"required"`

	IncludeAddress         bool `json:"-"`
	IsConsistentRead       bool `json:"-"`
	IncludeBusinessProfile bool `json:"-"`
	IncludeContactOwners   bool `json:"-"`
}

func (r *UserRepo) GetUser(params GetUserParams) (*models.User, error) {
	var builder = queryfunc.NewUserBuilder(queryfunc.UserBuilderOptions{
		IncludeAddress:         params.IncludeAddress,
		IncludeBusinessProfile: params.IncludeBusinessProfile,
		IncludeContactOwners:   params.IncludeContactOwners,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IsConsistentRead: params.IsConsistentRead,
	})
	var user models.User
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.GetRole().IsAdmin() {
				builder.Where("u.id = ?", params.UserID)
			} else {
				builder.Where("u.id = ?", params.GetUserID())
			}
		}).
		FirstFunc(&user)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

type GetMeParams struct {
	models.JwtClaimsInfo
	TeamID string `json:"team_id" query:"team_id"`
}

func (r *UserRepo) GetMe(params GetMeParams) (*models.User, error) {
	var builder = queryfunc.NewUserBuilder(queryfunc.UserBuilderOptions{
		IncludeAddress:         true,
		IncludeBusinessProfile: true,
		IncludeBrandTeam:       true,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var user models.User
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("u.id = ?", params.GetUserID())
		}).
		FirstFunc(&user)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	if params.TeamID != "" {
		var brandTeam models.BrandTeam
		err = r.db.First(&brandTeam, "team_id = ? AND user_id = ?", params.TeamID, params.GetUserID()).Error
		if err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrInvalidBrandTeamManager
			}
			return nil, err
		}
		user.BrandTeam = &brandTeam
	} else {
		user.BrandTeam = &models.BrandTeam{
			Role: enums.BrandTeamRoleManager,
			ID:   "",
		}
	}

	return &user, nil
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	var builder = queryfunc.NewUserBuilder(queryfunc.UserBuilderOptions{})
	var err = query.New(r.db, builder).
		Where("u.email = ?", strings.ToLower(email)).
		FirstFunc(&user)
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) AdminUpdateUserByID(form models.AdminUserUpdateForm) (*models.User, error) {
	var update models.User

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	if update.Coordinate != nil {
		err = update.Coordinate.CreateOrUpdate(r.db)
		if err != nil {
			return nil, err
		}
		update.CoordinateID = update.Coordinate.ID
	}

	update.Name = update.GetFullName()
	if form.GetRole().IsAdmin() {
		err = r.db.Omit(clause.Associations).Model(&models.User{}).Where("id = ?", form.UserID).Updates(&update).Error
	} else {
		err = r.db.Omit(clause.Associations).Model(&models.User{}).Where("id = ?", form.GetUserID()).Updates(&update).Error
	}
	if err != nil {
		if e := r.db.CheckUserDuplicateConstraint(err); e != nil {
			return nil, e
		}
		return nil, err
	}

	return r.GetUser(GetUserParams{
		JwtClaimsInfo:  form.JwtClaimsInfo,
		UserID:         form.UserID,
		IncludeAddress: true,
	})
}

func (r *UserRepo) UpdateUserByID(form models.UserUpdateForm) (*models.User, error) {
	var update models.User

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	if update.Coordinate != nil {
		err = update.Coordinate.CreateOrUpdate(r.db)
		if err != nil {
			return nil, err
		}
		update.CoordinateID = update.Coordinate.ID
	}

	update.Name = update.GetFullName()
	if form.GetRole().IsAdmin() {
		err = r.db.Omit(clause.Associations).Model(&models.User{}).Where("id = ?", form.UserID).Updates(&update).Error
	} else {
		err = r.db.Omit(clause.Associations).Model(&models.User{}).Where("id = ?", form.GetUserID()).Updates(&update).Error
	}
	if err != nil {
		if e := r.db.CheckUserDuplicateConstraint(err); e != nil {
			return nil, e
		}
		return nil, err
	}

	return r.GetUser(GetUserParams{
		JwtClaimsInfo:  form.JwtClaimsInfo,
		UserID:         form.UserID,
		IncludeAddress: true,
	})
}

type ArchiveUserParams struct {
	models.JwtClaimsInfo

	UserID string `param:"user_id" validate:"required"`
}

func (r *UserRepo) ArchiveUserByID(params ArchiveUserParams) error {
	var err = r.db.Delete(&models.User{}, "id = ?", params.UserID).Error

	return err
}

type UnarchiveUserParams struct {
	models.JwtClaimsInfo

	UserID string `param:"user_id" validate:"required"`
}

func (r *UserRepo) UnarchiveUserByID(params UnarchiveUserParams) error {
	var err = r.db.Unscoped().Model(&models.User{}).Where("id = ?", params.UserID).UpdateColumn("deleted_at", nil).Error

	return err
}

func (r *UserRepo) UpdatePassword(userID string, form models.UpdateMyPasswordForm) error {
	var user models.User
	var err = r.db.Select("Password").First(&user, "id = ?", userID).Error
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = user.ComparePassword(form.OldPassword)
	if err != nil {
		return errs.ErrOldPasswordIncorrect
	}

	user.Password = form.NewPassword
	err = user.HashPassword()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	err = r.db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("Password", user.Password).Error
	return err
}

func (r *UserRepo) Logout(userID string) error {
	var err = r.db.Transaction(func(tx *gorm.DB) error {
		var updates = map[string]interface{}{
			"token_issuer":  "",
			"logged_out_at": r.db.NowFunc().Unix(),
		}

		var err = tx.Model(&models.User{}).Where("id = ?", userID).Updates(&updates).Error
		if err != nil {
			return err
		}

		var tokens []*models.PushToken
		err = tx.Clauses(clause.Returning{
			Columns: []clause.Column{
				{Name: "token"},
			},
		}).
			Model(&tokens).
			Delete(&models.PushToken{}, "user_id = ?", userID).Error
		if err != nil {
			return err
		}

		for _, token := range tokens {
			err = customerio.GetInstance().Track.DeleteDevice(userID, token.Token)
			r.db.CustomLogger.Debugf("Delete customer.io device user_id=%s token=%s err=%+v", userID, token.Token, err)
		}

		return nil
	})

	return err
}

type RejectUserParams struct {
	models.JwtClaimsInfo

	UserID string `param:"user_id" validate:"required"`
}

func (r *UserRepo) RejectUser(params RejectUserParams) error {
	var updateUser = models.User{
		AccountStatus:          enums.AccountStatusRejected,
		AccountStatusChangedAt: values.Int64(time.Now().Unix()),
	}
	var err = r.db.Select("AccountStatus", "AccountStatusChangedAt").
		Model(&models.User{}).
		Where("id = ?", params.UserID).
		Updates(&updateUser).Error

	return err
}

type DeleteUserParams struct {
	models.JwtClaimsInfo

	UserID string `param:"user_id" validate:"required"`
}

func (r *UserRepo) DeleteUser(params DeleteUserParams) error {
	var err = r.db.Transaction(func(tx *gorm.DB) error {
		var err = tx.Unscoped().Delete(&models.User{}, "id = ?", params.UserID).Error
		if err != nil {
			return err
		}

		var deletedInquiries []*models.Inquiry
		err = tx.Unscoped().Model(&deletedInquiries).
			Clauses(clause.Returning{
				Columns: []clause.Column{{Name: "id"}},
			}).
			Delete(&models.Inquiry{}, "user_id = ?", params.UserID).Error
		if err != nil {
			return err
		}

		var inquiryIDs = lo.Map(deletedInquiries, func(item *models.Inquiry, index int) string {
			return item.ID
		})

		if len(inquiryIDs) > 0 {
			err = tx.Unscoped().Delete(&models.PurchaseOrder{}, "inquiry_id IN ?", inquiryIDs).Error
			if err != nil {
				return err
			}

			err = tx.Unscoped().Delete(&models.InquiryCartItem{}, "inquiry_id IN ?", inquiryIDs).Error
			if err != nil {
				return err
			}

		}

		return err
	})

	return err
}

type ApproveUserParams struct {
	models.JwtClaimsInfo
	UserID string `param:"user_id" validate:"required"`
}

func (r *UserRepo) ApproveUser(params ApproveUserParams) error {
	var updateUser = models.User{
		AccountStatus:          enums.AccountStatusActive,
		AccountStatusChangedAt: values.Int64(time.Now().Unix()),
		AccountVerified:        true,
		AccountVerifiedAt:      values.Int64(time.Now().Unix()),
	}
	var err = r.db.Select("AccountStatus", "AccountStatusChangedAt", "AccountVerified", "AccountVerifiedAt").
		Model(&models.User{}).
		Where("id = ?", params.UserID).
		Updates(&updateUser).Error

	return err
}

func (r *UserRepo) TrackActivity(userID string, form models.UserTrackActivityForm) (map[string]interface{}, error) {
	var updates = map[string]interface{}{
		"last_login": time.Now().Unix(),
	}

	if form.Timezone != "" {
		if _, err := time.LoadLocation(string(form.Timezone)); err == nil {
			updates["timezone"] = form.Timezone
		}
	}

	if form.CountryCode != "" {
		updates["country_code"] = form.CountryCode
	}

	var err = r.db.Model(&models.User{}).Where("id = ?", userID).Updates(&updates).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return updates, nil
}

func (r *UserRepo) CreateClientUser(form models.RegisterForm) (*models.User, error) {
	var user models.User
	var err = copier.Copy(&user, &form)
	if err != nil {
		return nil, err
	}
	user.Role = enums.RoleClient
	user.ID = helper.GenerateXID()

	err = r.db.Create(&user).Error
	if err != nil {
		if e := r.db.CheckUserDuplicateConstraint(err); e != nil {
			return nil, e
		}
		return nil, eris.Wrap(err, err.Error())
	}

	return &user, err
}

func (r *UserRepo) SetupIntentForClientSecert(userID string) (*stripe.SetupIntent, error) {
	var user models.User
	var err = r.db.Select("ID", "Email", "Name", "StripeCustomerID").First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if user.StripeCustomerID == "" {
		var err = user.CreateStripeCustomer(r.db, nil, false)
		if err != nil {
			return nil, eris.Wrap(err, "")
		}
	}

	result, err := stripehelper.GetInstance().SetupIntentForClientSecert(stripehelper.SetupIntentForBankAccountParams{
		StripeCustomerID: user.StripeCustomerID,
	})
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return result, err
}

func (r *UserRepo) GetCustomerIOUser(userID string) (*models.User, error) {
	var user models.User

	var err = query.New(r.db, queryfunc.NewUserCustomerIOBuilder(queryfunc.UserCustomerIOBuilderOptions{})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("u.deleted_at IS NULL and u.id = ?", userID)

		}).
		FirstFunc(&user)

	return &user, err
}

func (r *UserRepo) ChangePassword(form models.ChangePasswordForm) error {
	var update = models.User{
		Password:               form.NewPassword,
		RequiresChangePassword: false,
	}
	update.HashPassword()

	var err = r.db.Model(&models.User{}).Where("id = ?", form.UserID).Updates(&update).Error
	return err
}

func (r *UserRepo) GetShortUserInfo(userID string) (*models.User, error) {
	var user models.User
	var err = r.db.Select("ID", "Name", "Avatar").First(&user, "id = ?", userID).Error

	return &user, err
}

type GetUserPaymentMethodsParams struct {
	UserID string `param:"user_id" validate:"required"`
}

func (r *UserRepo) GetUserPaymentMethods(params GetUserPaymentMethodsParams) ([]*models.UserPaymentMethod, error) {
	return NewPaymentMethodRepo(r.db).GetPaymentMethods(GetPaymentMethodsParams{
		UserID: params.UserID,
	})
}

type InviteUserForm struct {
	models.JwtClaimsInfo

	FirstName   string     `json:"first_name" validate:"required"`
	LastName    string     `json:"last_name" validate:"required"`
	Email       string     `json:"email" validate:"required,email"`
	PhoneNumber string     `json:"phone_number,omitempty" param:"phone_number" query:"phone_number" form:"phone_number" validate:"omitempty,isPhone"`
	Role        enums.Role `json:"role" validate:"required,oneof=staff leader"`
	Team        enums.Team `json:"team" validate:"required,oneof=marketing sales dev designer operator customer_service finance"`
	RedirectURL string     `json:"redirect_url" validate:"required,startswith=http"`

	AccountStatus enums.AccountStatus `json:"-"`
}

type CreateClientForm struct {
	models.JwtClaimsInfo

	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number,omitempty" param:"phone_number" query:"phone_number" form:"phone_number" validate:"omitempty,isPhone"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	RedirectURL string `json:"redirect_url" validate:"required,startswith=http"`
}

func (r *UserRepo) CreateClient(form CreateClientForm) (*models.CreateInvitedUserResponse, error) {
	var user = models.User{
		FirstName:              form.FirstName,
		LastName:               form.LastName,
		Email:                  form.Email,
		Role:                   enums.RoleClient,
		PhoneNumber:            form.PhoneNumber,
		AccountStatus:          enums.AccountStatusActive,
		ContactOwnerIDs:        pq.StringArray([]string{form.GetUserID()}),
		RequiresChangePassword: true,
		Password:               form.Password,
		AccountVerified:        true,
		AccountVerifiedAt:      values.Int64(time.Now().Unix()),
	}
	user.ID = helper.GenerateXID()
	user.Name = user.GetFullName()
	token, err := user.GenerateResetPasswordToken(r.db.Configuration.JWTResetPasswordSecret, r.db.Configuration.JWTResetPasswordExpiry)
	if err != nil {
		return nil, err
	}
	err = r.db.Create(&user).Error
	if err != nil {
		if e := r.db.CheckUserDuplicateConstraint(err); e != nil {
			return nil, e
		}
		return nil, err
	}

	redirectURL, err := url.Parse(form.RedirectURL)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	var values = redirectURL.Query()

	values.Add("token", token)
	values.Add("email", user.Email)

	redirectURL.RawQuery = values.Encode()

	var resp = &models.CreateInvitedUserResponse{
		RedirectURL: redirectURL.String(),
		User:        user,
	}

	return resp, err
}

func (r *UserRepo) CreateInvitedUser(form InviteUserForm) (*models.CreateInvitedUserResponse, error) {
	var user = models.User{
		FirstName:       form.FirstName,
		LastName:        form.LastName,
		Email:           form.Email,
		Role:            form.Role,
		Team:            form.Team,
		PhoneNumber:     form.PhoneNumber,
		AccountStatus:   form.AccountStatus,
		ContactOwnerIDs: pq.StringArray([]string{form.GetUserID()}),
	}
	user.ID = helper.GenerateXID()
	user.Name = user.GetFullName()
	token, err := user.GenerateResetPasswordToken(r.db.Configuration.JWTResetPasswordSecret, r.db.Configuration.JWTResetPasswordExpiry)
	if err != nil {
		return nil, err
	}
	err = r.db.Create(&user).Error
	if err != nil {
		if e := r.db.CheckUserDuplicateConstraint(err); e != nil {
			return nil, e
		}
		return nil, err
	}

	redirectURL, err := url.Parse(form.RedirectURL)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	var values = redirectURL.Query()

	values.Add("token", token)
	values.Add("email", user.Email)

	redirectURL.RawQuery = values.Encode()

	var resp = &models.CreateInvitedUserResponse{
		RedirectURL: redirectURL.String(),
		User:        user,
	}

	return resp, err
}

func (r *UserRepo) CreateTrxInvitedUser(trx *gorm.DB, secret string, expiry time.Duration, form InviteUserForm) (*models.CreateInvitedUserResponse, error) {
	var user = models.User{
		FirstName:     form.FirstName,
		LastName:      form.LastName,
		Email:         form.Email,
		Role:          form.Role,
		Team:          form.Team,
		PhoneNumber:   form.PhoneNumber,
		AccountStatus: form.AccountStatus,
	}
	user.ID = helper.GenerateXID()
	user.Name = user.GetFullName()
	token, err := user.GenerateResetPasswordToken(secret, expiry)
	if err != nil {
		return nil, err
	}
	err = trx.Create(&user).Error
	if err != nil {
		if e := r.db.CheckUserDuplicateConstraint(err); e != nil {
			return nil, e
		}
		return nil, err
	}

	redirectURL, err := url.Parse(form.RedirectURL)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	var values = redirectURL.Query()

	values.Add("token", token)
	values.Add("email", user.Email)

	redirectURL.RawQuery = values.Encode()

	var resp = &models.CreateInvitedUserResponse{
		RedirectURL: redirectURL.String(),
		User:        user,
	}

	return resp, err
}

type PaginateUsersRoleParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *UserRepo) PaginateUsersRole(params PaginateUsersRoleParams) (result []models.UserRoleStat, err error) {
	err = query.New(r.db, queryfunc.NewUserRoleBuilder(queryfunc.UserRoleBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		Where("COALESCE(u.team,'') <> ''").
		FindFunc(&result)
	return
}

func (r *UserRepo) OnboardingSubmit(tx *gorm.DB, form models.BusinessProfileCreateForm) (*models.User, error) {
	var businessProfileUpdate models.BusinessProfile
	var err = copier.Copy(&businessProfileUpdate, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	businessProfileUpdate.UserID = form.GetUserID()

	if businessProfileUpdate.MillFabricTypes != nil {
		var fabricTypes []string
		for _, item := range *businessProfileUpdate.MillFabricTypes {
			fabricTypes = append(fabricTypes, item.FabricValue)
		}
		businessProfileUpdate.FlatMillFabricTypes = fabricTypes
	}

	var businessProfile models.BusinessProfile
	err = tx.Select("ID").First(&businessProfile, "user_id = ?", form.GetUserID()).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			err = tx.Omit(clause.Associations).Model(&models.BusinessProfile{}).Where("user_id = ?", form.GetUserID()).Create(&businessProfileUpdate).Error
			if err != nil {
				return nil, eris.Wrap(err, err.Error())
			}
		}
	} else {
		err = tx.Omit(clause.Associations).Model(&models.BusinessProfile{}).Where("user_id = ?", form.GetUserID()).Updates(&businessProfileUpdate).Error
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	var user = form.User

	if user == nil {
		user, err = r.GetUser(GetUserParams{
			JwtClaimsInfo:    form.JwtClaimsInfo,
			UserID:           form.JwtClaimsInfo.GetUserID(),
			IsConsistentRead: true,
		})
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	if user.OnboardingSubmitAt == nil {
		err = r.db.Omit(clause.Associations).
			Model(&models.User{}).
			Where("id = ? AND onboarding_submit_at IS NULL", form.GetUserID()).
			UpdateColumn("OnboardingSubmitAt", time.Now().Unix()).Error
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	return user, nil
}

type TeamClientInviteForm struct {
	models.JwtClaimsInfo

	FirstName   string                   `json:"first_name" validate:"required"`
	LastName    string                   `json:"last_name" validate:"required"`
	PhoneNumber string                   `json:"phone_number,omitempty" param:"phone_number" query:"phone_number" form:"phone_number" validate:"omitempty,isPhone"`
	Email       string                   `json:"email" validate:"required,email"`
	RedirectURL string                   `json:"redirect_url" validate:"required,startswith=http"`
	Actions     enums.BrandMemberActions `json:"actions,omitempty" param:"actions" query:"actions" form:"actions"`
}

func (r *UserRepo) TeamClientInvite(form TeamClientInviteForm) (resp *models.CreateInvitedUserResponse, err error) {
	var inviteByUser models.User
	err = r.db.Select("ID", "Name", "Role", "Team").First(&inviteByUser, "id = ?", form.GetUserID()).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	var invitedUser models.User
	_ = r.db.Select("ID", "Name", "FirstName", "LastName", "Email", "Role", "Team", "PhoneNumber", "AccountStatus").First(&invitedUser, "email = ?", form.Email)

	var brandTeam models.BrandTeam
	_ = r.db.First(&brandTeam, "user_id = ? and role = ?", inviteByUser.ID, enums.RoleManager).Error

	var countTeamMember int
	_ = r.db.Model(&models.BrandTeam{}).Select("count(1)").Find(&countTeamMember, "team_id = ?", inviteByUser.ID)
	if brandTeam.TeamID == "" && countTeamMember > 0 {
		err = errs.ErrInvalidBrandTeamManager
		return
	}

	if invitedUser.ID == inviteByUser.ID {
		err = errs.ErrCanNotInviteYourself
		return
	}

	if invitedUser.ID != "" {
		if invitedUser.Role != enums.RoleClient {
			return nil, errs.ErrTeamMemberUserNotAbleToJoin
		}

		var countTeam int
		_ = r.db.Model(&models.BrandTeam{}).Select("count(1)").Find(&countTeam, "user_id = ?", invitedUser.ID)
		if countTeam > 0 {
			err = errs.ErrTeamMemberAlreadyInAnotherTeam
			return
		}

		err = r.db.Transaction(func(tx *gorm.DB) (err error) {
			if brandTeam.TeamID == "" {
				brandTeam = models.BrandTeam{
					UserID: inviteByUser.ID,
					TeamID: inviteByUser.ID,
					Role:   enums.BrandTeamRoleManager,
				}
				if err = tx.Create(&brandTeam).Error; err != nil {
					return
				}
			}

			// create if not exits brand team - user invited
			var brandTeamInvited = models.BrandTeam{
				UserID:  invitedUser.ID,
				TeamID:  brandTeam.TeamID,
				Role:    enums.BrandTeamRoleStaff,
				Actions: form.Actions.ToStringSlice(),
			}
			if err = tx.Create(&brandTeamInvited).Error; err != nil {
				err = errs.ErrTeamMemberAlreadyInThisTeam
				return
			}

			return
		})
		resp = &models.CreateInvitedUserResponse{
			InvitedByUser: inviteByUser,
			User:          invitedUser,
		}
		return
	}

	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		// create if not exits brand team - user invite
		if brandTeam.TeamID == "" {
			brandTeam = models.BrandTeam{
				UserID: inviteByUser.ID,
				TeamID: inviteByUser.ID,
				Role:   enums.BrandTeamRoleManager,
			}
			if err = tx.Create(&brandTeam).Error; err != nil {
				return
			}
		}

		// if not exits -> create invited user
		secret := r.db.Configuration.JWTResetPasswordSecret
		expiry := r.db.Configuration.JWTResetPasswordExpiry
		var createForm InviteUserForm
		if err = copier.Copy(&createForm, form); err != nil {
			return
		}
		createForm.Role = enums.RoleClient
		createForm.AccountStatus = enums.AccountStatusPendingReview
		resp, err = r.CreateTrxInvitedUser(tx, secret, expiry, createForm)
		if err != nil {
			return
		}

		// create if not exits brand team - user invited
		var brandTeamInvited = models.BrandTeam{
			UserID:  resp.User.ID,
			TeamID:  brandTeam.TeamID,
			Role:    enums.BrandTeamRoleStaff,
			Actions: form.Actions.ToStringSlice(),
		}
		if err = tx.Create(&brandTeamInvited).Error; err != nil {
			return
		}

		resp.InvitedByUser = inviteByUser
		return
	})

	return
}

func (r *UserRepo) AssignContactOwners(form models.AssignContactOwnersForm) (err error) {
	var user models.User
	err = r.db.Model(&models.User{}).Where("id IN ? AND role IN ?", form.ContactOwnerIDs, []enums.Role{
		enums.RoleLeader,
		enums.RoleStaff,
	}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = eris.Wrap(errs.ErrUserNotFound, "cannot get contact owners")
		return
	}
	if err != nil {
		return
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var ids = pq.StringArray(lo.Uniq(form.ContactOwnerIDs))
		err = tx.Model(&models.User{}).Where("id = ?", form.UserID).UpdateColumn("ContactOwnerIDs", ids).Error
		if err != nil {
			return err
		}

		return tx.Model(&models.Inquiry{}).
			Where("user_id = ? AND (assignee_ids IS NULL OR ARRAY_LENGTH(assignee_ids, 1) == 0)", form.UserID).
			UpdateColumn("AssigneeIDs", ids).Error
	})
	return err
}

type UpdateBrandTeamMemberActionsForm struct {
	models.JwtClaimsInfo
	MemberId string                   `json:"member_id,omitempty" param:"member_id" query:"member_id" form:"member_id"`
	Actions  enums.BrandMemberActions `json:"actions,omitempty" param:"actions" query:"actions" form:"actions"`
}

func (r *UserRepo) UpdateBrandTeamMemberActions(form UpdateBrandTeamMemberActionsForm) (err error) {
	if !r.isTeamManager(form.GetUserID()) {
		err = errs.ErrActionBrandTeamManagerOnly
		return
	}
	var updates models.BrandTeam
	updates.Actions = form.Actions.ToStringSlice()
	err = r.db.Model(&models.BrandTeam{}).Where("user_id = ?", form.MemberId).Updates(&updates).Error
	return
}

type DeleteBrandTeamMemberForm struct {
	models.JwtClaimsInfo
	MemberId string `json:"member_id,omitempty" param:"member_id" query:"member_id" form:"member_id"`
}

func (r *UserRepo) DeleteBrandTeamMember(form DeleteBrandTeamMemberForm) (err error) {
	if !r.isTeamManager(form.GetUserID()) {
		err = errs.ErrActionBrandTeamManagerOnly
		return
	}
	err = r.db.Unscoped().Delete(&models.BrandTeam{}, "user_id = ?", form.MemberId).Error
	return
}

func (r *UserRepo) isTeamManager(userID string) bool {
	var find models.BrandTeam
	_ = r.db.Select("id").First(&find, "user_id = ? and team_id = ?", userID, userID)

	return find.ID != ""
}

func (r *UserRepo) GetLastShippingAddress(claims models.JwtClaimsInfo) (*models.Address, error) {
	var lastInquiryShippingAddressID string
	r.db.Raw(`
	SELECT shipping_address_id 
	FROM inquiries 
	WHERE user_id = ?
	ORDER BY created_at DESC
	`, claims.GetUserID()).Find(&lastInquiryShippingAddressID)
	if lastInquiryShippingAddressID != "" {
		address, err := NewAddressRepo(r.db).GetAddress(GetAddressParams{
			AddressID: lastInquiryShippingAddressID,
		})
		if err != nil {
			return nil, err
		}

		return address, nil
	}

	var user models.User
	var err = r.db.Select("ID", "Name", "PhoneNumber", "Email", "CoordinateID").First(&user, "id = ?", claims.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	var address models.Address
	address.UserID = user.ID
	address.Name = user.Name
	address.PhoneNumber = user.PhoneNumber
	if user.CoordinateID != "" {
		var coordinate models.Coordinate
		if err = r.db.First(&coordinate, "id = ?", user.CoordinateID).Error; err == nil {
			address.Coordinate = &coordinate
			address.CoordinateID = coordinate.ID
		}
	}

	return &address, nil
}

type PaginateBrandTeamMembersParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	TeamID string `json:"team_id" param:"team_id" query:"team_id" form:"team_id"`
}

func (r *UserRepo) PaginateBrandTeamMembers(params PaginateBrandTeamMembersParams) *query.Pagination {
	var memberIDs []string

	if params.TeamID != "" {
		var err = r.db.Model(&models.BrandTeam{}).Select("user_id").Find(&memberIDs, "team_id = ?", params.TeamID).Error
		if err != nil {
			return &query.Pagination{
				Records: []*models.PurchaseOrder{},
			}
		}

		if len(memberIDs) >= 0 && !lo.Contains(memberIDs, params.GetUserID()) {
			return &query.Pagination{
				Records: []*models.PurchaseOrder{},
			}
		}
	} else {
		var user models.User
		var err = r.db.Select("ID", "CreatedAt", "LastLogin", "Name", "Avatar", "Email", "PhoneNumber", "AccountStatus").First(&user, "id = ?", params.GetUserID()).Error
		if err != nil {
			return &query.Pagination{
				Records: []*models.PurchaseOrder{},
			}
		}
		user.BrandTeam = &models.BrandTeam{
			Role: enums.BrandTeamRoleManager,
		}

		return &query.Pagination{
			Records: []*models.User{&user},
		}
	}

	return NewUserRepo(r.db).PaginateUsers(PaginateUsersParams{
		PaginationParams: params.PaginationParams,
		JwtClaimsInfo:    params.JwtClaimsInfo,
		UserIDs:          memberIDs,
		IncludeBrandTeam: true,
	})
}

func (r *UserRepo) GetTeamManagerID(userID string) string {
	var member models.BrandTeam
	// team_id is manager's user_id
	_ = r.db.Select("TeamID").First(&member, "user_id = ?", userID).Error
	if member.TeamID != "" {
		return member.TeamID
	}
	return userID
}

type GetAccessTokenParams struct {
	models.JwtClaimsInfo

	UserID string `json:"user_id" param:"user_id"`
}

type GetAccessTokenResult struct {
	Token string `json:"token,omitempty"`
}

func (r *UserRepo) GetAccessToken(params GetAccessTokenParams) (*GetAccessTokenResult, error) {
	var user models.User
	var err = r.db.First(&user, "id = ?", params.UserID).Error
	if err != nil {
		return nil, err
	}

	user.TokenIssuer = fmt.Sprintf("ghost|%s", params.GetUserID())

	token, err := user.GenerateToken(r.db.Configuration.JWTSecret)
	if err != nil {
		return nil, err
	}

	var result = GetAccessTokenResult{
		Token: token,
	}
	return &result, err
}

func (r *UserRepo) GetTeams(claims models.JwtClaimsInfo) (list []*models.BrandTeam) {
	r.db.Find(&list, "user_id = ?", claims.GetUserID())

	for _, v := range list {
		r.db.Model(&models.User{}).Select("Name").First(&v.TeamName, "id = ?", v.TeamID)
	}

	list = append(list, &models.BrandTeam{
		Role:     enums.BrandTeamRoleManager,
		TeamName: "Personal Account",
	})

	return
}

func (r *UserRepo) CompleteInquiryTutorial(claims models.JwtClaimsInfo) error {
	return r.db.Model(&models.User{}).Where("id = ?", claims.GetUserID()).UpdateColumn("CompletedInquiryTutorialAt", time.Now().Unix()).Error
}

type ConnectZaloParams struct {
	models.JwtClaimsInfo
	ZaloID string `json:"zalo_id" query:"zalo_id" params:"zalo_id" form:"zalo_id" validate:"required"`
}

func (r *UserRepo) ConnectZalo(params ConnectZaloParams) (err error) {
	err = r.db.Model(&models.User{}).
		Where("id = ?", params.GetUserID()).
		Where("COALESCE(zalo_id,'') = ''").
		Update("zalo_id", params.ZaloID).Error
	return
}

type DisconnectZaloParams struct {
	models.JwtClaimsInfo
}

func (r *UserRepo) DisconnectZalo(params DisconnectZaloParams) (err error) {
	err = r.db.Model(&models.User{}).
		Where("id = ?", params.GetUserID()).
		Update("zalo_id", gorm.Expr("NULL")).Error
	return
}

// User - Product Class

type UpdateUserProductClassesParams struct {
	models.JwtClaimsInfo
	UserID              string `json:"user_id"`
	InquiryID           string `json:"inquiry_id"`
	PurchaseOrderID     string `json:"purchase_order_id"`
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id"`
}

func (r *UserRepo) UpdateProductClasses(params UpdateUserProductClassesParams) (err error) {
	if params.UserID == "" {
		err = errors.New("empty param user_id")
		return
	}

	var user models.User
	err = r.db.Model(&models.User{}).Select("id", "product_classes").Where("id = ?", params.UserID).First(&user).Error
	if err != nil {
		return
	}

	validateURL := func(fileURL string) string {
		if fileURL == "" {
			return ""
		}
		ext := filepath.Ext(fileURL)
		if !(ext == ".jpg" || ext == ".jpeg") {
			return ""
		}
		if len(fileURL) > 5 && fileURL[:5] != "https" {
			fileURL = fmt.Sprintf("https://%s/%s", r.db.Configuration.StorageURL, fileURL)
		}
		return fileURL
	}

	var urls []string

	// Inquiry
	if params.InquiryID != "" {
		var rfq models.Inquiry
		_ = r.db.Model(&models.Inquiry{}).Select("id", "attachments").Where("id = ?", params.InquiryID).First(&rfq)
		if rfq.ID != "" && rfq.Attachments != nil {
			for _, att := range *rfq.Attachments {
				imageURL := validateURL(att.FileKey)
				if imageURL != "" {
					urls = append(urls, imageURL)
				}
			}
		}
	}

	// Purchase Order
	if params.PurchaseOrderID != "" {
		var po models.PurchaseOrder
		_ = r.db.Model(&models.PurchaseOrder{}).Select("id", "attachments").Where("id = ?", params.PurchaseOrderID).First(&po)
		if po.ID != "" && po.Attachments != nil {
			for _, att := range *po.Attachments {
				imageURL := validateURL(att.FileKey)
				if imageURL != "" {
					urls = append(urls, imageURL)
				}
			}
		}
	}

	// Bulk Purchase Order
	if params.BulkPurchaseOrderID != "" {
		var bpo models.BulkPurchaseOrder
		_ = r.db.Model(&models.BulkPurchaseOrder{}).Select("id", "attachments").Where("id = ?", params.BulkPurchaseOrderID).First(&bpo)
		if bpo.ID != "" && bpo.Attachments != nil {
			for _, att := range *bpo.Attachments {
				imageURL := validateURL(att.FileKey)
				if imageURL != "" {
					urls = append(urls, imageURL)
				}
			}
		}
	}

	// Classify
	var (
		classifyParams  []ai.ImageClassifyParams
		classifyResults []ai.ImageClassifyResponse
	)
	urls = lo.Uniq(urls)
	for _, fileURL := range urls {
		classifyParams = append(classifyParams, ai.ImageClassifyParams{
			Image:      fileURL,
			Size:       640,
			Confidence: 0.3,
			Overlap:    0,
		})
	}
	if len(classifyParams) == 0 {
		return
	}

	classifyResults, err = ai.ClassifyMultiImage(classifyParams)
	if err != nil {
		return
	}
	var updateProductClasses models.UserProductClasses
	for _, result := range classifyResults {
		for _, predict := range result.Predictions {
			updateProductClasses = append(updateProductClasses, models.UserProductClass{
				Class: predict.Class,
				Conf:  predict.Confidence,
			})
		}
	}
	updateProductClasses = append(updateProductClasses, user.ProductClasses...)
	updateProductClasses = updateProductClasses.Uniq()
	if len(updateProductClasses) == 0 {
		return
	}

	err = r.db.Model(&models.User{}).Where("id = ?", params.UserID).Update("product_classes", updateProductClasses).Error
	return
}
