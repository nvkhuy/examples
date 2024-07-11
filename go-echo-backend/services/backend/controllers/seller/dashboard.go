package controllers

import (
	"strconv"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
)

// SellerDashboardRevenue Get revenue
// @Tags Seller-Dashboard
// @Summary Get revenue
// @Description Get revenue
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param item_id path string true "Item ID"
// @Success 200 {object} models.Order
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/dashboard/revenue [get]
func SellerDashboardRevenue(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var result models.SellerDashboardRevenueResponse
	for i := 1; i <= 12; i++ {
		result.MonthlyTrend = append(result.MonthlyTrend, models.MonthlyRevenue{
			Month:   strconv.Itoa(i),
			Revenue: price.NewFromFloat(100),
		})
	}

	for i := 1; i <= 4; i++ {
		result.WeeklyTrend = append(result.WeeklyTrend, models.WeekLyRevenue{
			Week:    strconv.Itoa(i),
			Revenue: price.NewFromFloat(100),
		})
	}

	// result.CurrentMonth = strconv.Itoa(6)
	// result.CurrentMonthName = "June"
	result.MonthlyIncreaseRate = 12
	result.WeeklyIncreaseRate = -11
	result.MonthlyRevenue = price.NewFromFloat(2500)
	result.WeeklyRevenue = price.NewFromFloat(200)

	return cc.Success(result)
}

// SellerDashboard Get general info
// @Tags Seller-Dashboard
// @Summary Get general info
// @Description Get general info
// @Accept  json
// @Produce  json
// @Success 200 {object} models.SellerDashboardResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/dashboard [get]
func SellerDashboard(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return err
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return err
	}

	params.JwtClaimsInfo = claims
	params.Limit = 2

	var products = repo.NewProductRepo(cc.App.DB).PaginateProducts(params)

	return cc.Success(models.SellerDashboardResponse{
		RecommendedProducts: products.Records,
		News:                []string{},
		Messages:            []string{},
	})

}
