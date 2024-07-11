package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// @Tags Buyer-OrderGroup
// @Summary List order group
// @Description This API allows buyer to list order groups
// @Accept  json
// @Produce  json
// @Param page query int false
// @Param limit query int false
// @Success 200 {object} query.Pagination{Records=[]models.OrderGroup}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/order_groups [get]
func BuyerGetOrderGroupList(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.GetOrderGroupListRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	orderGroups, err := repo.NewOrderGroupRepo(cc.App.DB).GetOrderGroupList(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(orderGroups)
}

// @Tags Buyer-OrderGroup
// @Summary Get order group detail
// @Description This API allows buyer to get order groupbdetail
// @Accept  json
// @Produce  json
// @Param order_group_id path string true
// @Success 200 {object} models.OrderGroup
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/order_groups/:order_group_id [get]
func BuyerGetOrderGroupDetail(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.GetOrderGroupDetailRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	orderGroup, err := repo.NewOrderGroupRepo(cc.App.DB).GetOrderGroupDetail(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(orderGroup)
}

// @Tags Buyer-OrderGroup
// @Summary Create order group
// @Description This API allows buyer to create order group
// @Accept  json
// @Produce  json
// @Param data body models.CreateOrderGroupRequest true
// @Success 200 {object} models.CreateOrderGroupRequest
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/order_groups [post]
func BuyerCreateOrderGroup(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.CreateOrderGroupRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	orderGroups, err := repo.NewOrderGroupRepo(cc.App.DB).CreateOrderGroup(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(orderGroups)
}

// @Tags Buyer-OrderGroup
// @Summary Assign order group
// @Description This API allows buyer to assign order group
// @Accept  json
// @Produce  json
// @Param data body models.AssignOrderGroupRequest true
// @Success 200 {object} models.AssignOrderGroupRequest
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/order_groups/assign [post]
func BuyerAssignOrderGroup(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params models.AssignOrderGroupRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	err = repo.NewOrderGroupRepo(cc.App.DB).AssignOrderGroups(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(true)
}
