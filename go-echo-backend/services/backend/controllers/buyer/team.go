package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// BrandTeamMembers
// @Tags Marketplace-Support
// @Summary Search SupportTicket
// @Description Search SupportTicket
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Success 200 {object} models.SupportTicket
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/teams/members [get]
func BrandTeamMembers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBrandTeamMembersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewUserRepo(cc.App.DB).PaginateBrandTeamMembers(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BrandTeamInvite
// @Tags Marketplace-Support
// @Summary Search SupportTicket
// @Description Search SupportTicket
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Success 200 {object} models.SupportTicket
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/teams/invite [post]
func BrandTeamInvite(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.TeamClientInviteForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewUserRepo(cc.App.DB).TeamClientInvite(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if result.User.ID != "" {
		tasks.HubspotCreateContactTask{
			Data: &hubspot.ContactPropertiesForm{
				Email:          result.User.Email,
				Firstname:      result.User.FirstName,
				Lastname:       result.User.LastName,
				Phone:          result.User.PhoneNumber,
				Lifecyclestage: "lead",
			},
		}.Dispatch(cc.Request().Context())

	}

	if result != nil && result.RedirectURL != "" {
		_, err = tasks.TrackCustomerIOTask{
			UserID: result.User.ID,
			Event:  customerio.EventInviteBrandMember,
			Data: map[string]interface{}{
				"email":      form.Email,
				"link":       result.RedirectURL,
				"invited_by": result.InvitedByUser.GetCustomerIOMetadata(nil),
			},
		}.Dispatch(cc.Request().Context())
	}
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result.User)
}

// UpdateBrandTeamMemberActions
// @Tags Marketplace-Support
// @Summary Search SupportTicket
// @Description Search SupportTicket
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Success 200 {object} models.SupportTicket
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/teams/{member_id}/actions [put]
func UpdateBrandTeamMemberActions(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.UpdateBrandTeamMemberActionsForm
	var err error
	var claims models.JwtClaimsInfo

	claims, err = cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	if err = cc.BindAndValidate(&form); err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewUserRepo(cc.App.DB).UpdateBrandTeamMemberActions(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success("updated!")
}

// DeleteBrandTeamMember
// @Tags Marketplace-Support
// @Summary Search SupportTicket
// @Description Search SupportTicket
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Success 200 {object} models.SupportTicket
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/teams/members/:{member_id} [put]
func DeleteBrandTeamMember(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.DeleteBrandTeamMemberForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	err = repo.NewUserRepo(cc.App.DB).DeleteBrandTeamMember(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}
