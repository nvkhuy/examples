package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// SearchUsers Search users
// @Tags Admin-User
// @Summary Search users
// @Description Search users
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword" default(1)
// @Param limit query int false "Size of page" default(20)
// @Success 200 {object} []models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/search [get]
func SearchUsers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateUsersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	var result = repo.NewUserRepo(cc.App.DB).SearchUsers(params)

	return cc.Success(result)
}

// PaginateUsers Get users
// @Tags Admin-User
// @Summary Get users
// @Description Get users
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword" default(1)
// @Param sorts query []string false "Sort" default("last_login desc") enums("last_login desc","last_login asc")
// @Param page query int false "Page index" default(1)
// @Param limit query int false "Size of page" default(20)
// @Success 200 {object} query.Pagination{records=[]models.User}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users [get]
func PaginateUsers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateUsersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludeContactOwners = true
	var result = repo.NewUserRepo(cc.App.DB).PaginateUsers(params)

	return cc.Success(result)
}

// PaginateUserRecent Get recent users created
// @Tags Admin-User
// @Summary Get recent users created
// @Description Get recent users created
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword" default(1)
// @Param sorts query []string false "Sort" default("last_login desc") enums("last_login desc","last_login asc")
// @Param page query int false "Page index" default(1)
// @Param limit query int false "Size of page" default(20)
// @Success 200 {object} query.Pagination{records=[]models.User}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/recent [get]
func PaginateUserRecent(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateUsersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewUserRepo(cc.App.DB).PaginateRecentUser(params)

	return cc.Success(result)
}

// GetUser Get user
// @Tags Admin-User
// @Summary Get user
// @Description Get user
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID of user"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id} [get]
func GetUser(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetUserParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.IncludeBusinessProfile = true
	params.IncludeAddress = true
	params.IncludeContactOwners = true
	u, err := repo.NewUserRepo(cc.App.DB).GetUser(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(u)
}

// UpdateUser Update user
// @Tags Admin-User
// @Summary Update user
// @Description Update user
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Param data body models.UserUpdateForm true "Form"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id} [put]
func UpdateUser(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.AdminUserUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	u, err := repo.NewUserRepo(cc.App.DB).AdminUpdateUserByID(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.SyncCustomerIOUserTask{
		UserID: u.ID,
	}.Dispatch(cc.Request().Context())

	return cc.Success(u)
}

// ArchiveUser Archive user
// @Tags Admin-User
// @Summary Archive user
// @Description Archive user
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/archive [delete]
func ArchiveUser(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.ArchiveUserParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewUserRepo(cc.App.DB).ArchiveUserByID(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}

// UnarchiveUser Unarchive user
// @Tags Admin-User
// @Summary Unarchive user
// @Description Unarchive user
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/archive [delete]
func UnarchiveUser(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UnarchiveUserParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewUserRepo(cc.App.DB).UnarchiveUserByID(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Unarchived")
}

// ApproveUser Approve user
// @Tags Admin-User
// @Summary Approve user
// @Description Approve user
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/approve [post]
func ApproveUser(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.ApproveUserParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewUserRepo(cc.App.DB).ApproveUser(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.ApproveUserTask{
		UserID: params.UserID,
	}.Dispatch(c.Request().Context())

	return cc.Success("Approved")
}

// InviteUser Invite user
// @Tags Admin-User
// @Summary Invite user
// @Description Invite user
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/invite [post]
func InviteUser(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.InviteUserForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	form.AccountStatus = enums.AccountStatusActive
	result, err := repo.NewUserRepo(cc.App.DB).CreateInvitedUser(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, err = tasks.TrackCustomerIOTask{
		UserID: result.User.ID,
		Event:  customerio.EventAdminInviteNewUser,
		Data: map[string]interface{}{
			"email":        form.Email,
			"first_name":   form.FirstName,
			"last_name":    form.LastName,
			"role":         form.Role,
			"team":         form.Team,
			"role_display": form.Role.DisplayName(),
			"team_display": form.Team.DisplayName(),
			"link":         result.RedirectURL,
		},
	}.Dispatch(cc.Request().Context())
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result.User)
}

// CreateClient CreateFromPayload client
// @Tags Admin-User
// @Summary CreateFromPayload client
// @Description CreateFromPayload client
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/client [post]
func CreateClient(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.CreateClientForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result, err := repo.NewUserRepo(cc.App.DB).CreateClient(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if result.User.Role == enums.RoleClient && cc.App.Config.IsProd() {
		tasks.HubspotCreateContactTask{
			Data: &hubspot.ContactPropertiesForm{
				Email:     result.User.Email,
				Firstname: result.User.FirstName,
				Lastname:  result.User.LastName,
				Phone:     result.User.PhoneNumber,
			},
		}.Dispatch(cc.Request().Context())
	}

	tasks.SyncCustomerIOUserTask{
		UserID: result.User.ID,
	}.Dispatch(cc.Request().Context())

	_, err = tasks.TrackCustomerIOTask{
		UserID: result.User.ID,
		Event:  customerio.EventAdminInviteClient,
		Data: map[string]interface{}{
			"email":      form.Email,
			"first_name": form.FirstName,
			"last_name":  form.LastName,
			"link":       result.RedirectURL,
		},
	}.Dispatch(cc.Request().Context())
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result.User)
}

// RejectUser Reject user
// @Tags Admin-User
// @Summary Reject user
// @Description Reject user
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/reject [delete]
func RejectUser(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.RejectUserParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewUserRepo(cc.App.DB).RejectUser(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.TrackCustomerIOTask{
		UserID: params.UserID,
		Event:  customerio.EventNotifyUserRejected,
	}.Dispatch(c.Request().Context())

	return cc.Success("Rejected")
}

// ChangePassword change password
// @Tags Admin-User
// @Summary Change password
// @Description Change password
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/change_password [put]
func ChangePassword(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.ChangePasswordForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = repo.NewUserRepo(cc.App.DB).ChangePassword(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Changed")
}

// DeleteUser Delete user
// @Tags Admin-User
// @Summary Delete user
// @Description Delete user
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/delete [delete]
func DeleteUser(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteUserParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewUserRepo(cc.App.DB).DeleteUser(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.SyncCustomerIOUserTask{
		UserID: params.UserID,
	}.Dispatch(cc.Request().Context())

	return cc.Success("Deleted")
}

// GetUserPaymentMethods Get user payment methods
// @Tags Admin-User
// @Summary Get user payment methods
// @Description Get user payment methods
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Security ApiKeyAuth
// @Router /api/v1/admin/users/{user_id}/payment_methods [get]
func GetUserPaymentMethods(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	list, err := repo.NewPaymentMethodRepo(cc.App.DB).GetPaymentMethods(repo.GetPaymentMethodsParams{
		UserID: cc.GetPathParamString("user_id"),
	})
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(list)
}

// GetUserBanks Get user banks
// @Tags Admin-User
// @Summary Get user payment methods
// @Description Get user payment methods
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Security ApiKeyAuth
// @Router /api/v1/admin/users/{user_id}/banks [get]
func GetUserBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	list, err := repo.NewUserBankRepo(cc.App.DB).GetUserBankInfos(repo.GetUserBankInfosParams{
		UserID: cc.GetPathParamString("user_id"),
	})
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(list)
}

// UpdateUser Update user
// @Tags Admin-User
// @Summary Update user
// @Description Update user
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Param data body models.UserUpdateForm true "Form"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/assign_owners [put]
func AssignContactOwners(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.AssignContactOwnersForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = repo.NewUserRepo(cc.App.DB).AssignContactOwners(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Assigned")
}

// GetAccessToken Get access token
// @Tags Admin-User
// @Summary Get access token
// @Description Get access token
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/access_token [post]
func GetAccessToken(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetAccessTokenParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewUserRepo(cc.App.DB).GetAccessToken(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UploadBulkPurchaseOrder
// @Tags Admin-User
// @Summary delete inquiry comment
// @Description delete inquiry comment
// @Accept  json
// @Produce  json
// @Param data body models.ContentCommentCreateForm true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/upload_bulks [post]
func UploadBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UploadExcelParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return err
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	results, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).UploadExcel(params)
	if err != nil {
		return err
	}

	return cc.Success(results)
}

// UploadBulkPurchaseOrder
// @Tags Admin-User
// @Summary delete inquiry comment
// @Description delete inquiry comment
// @Accept  json
// @Produce  json
// @Param data body models.ContentCommentCreateForm true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/users/{user_id}/activities [get]
func GetActivities(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.UserIDParam

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return err
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := cc.App.CustomerIOClient.GetActivities(params.UserID, params.GetActivitiesParams)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
