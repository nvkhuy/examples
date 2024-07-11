package tests

import (
	"log"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestDataAnalyticRepo_Overview(t *testing.T) {
	var app = initApp("local")
	overview := repo.NewDataAnalyticRepo(app.AnalyticDB).Overview(repo.DataAnalyticPlatformOverviewParam{})
	log.Println(overview)
}

func TestDataAnalyticRepo_SearchProducts(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).SearchProducts(repo.DataAnalyticSearchProductsParam{
		PaginationParams: models.PaginationParams{
			//Keyword: "Necklace",
			Page:  1,
			Limit: 10,
		},
		SortBy:           "growth_rate",
		IsSortDescending: true,
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_RecommendProducts(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).RecommendProducts(repo.DataAnalyticRecommendProductsParam{
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 7,
		},
		//RecommendProductID: "cl6bsk7skmpao012830g",
		ProductID: "cmri85hgjv3t7s5c88eg",
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_TopProducts(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).TopProducts(repo.DataAnalyticTopProductsParam{
		PaginationParams: models.PaginationParams{
			Limit: 10,
		},
		DateFrom:    1702047296,
		DateTo:      1702875504,
		OrderBy:     "sold",
		SubCategory: "Women Coats",
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_TopCategories(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).TopCategories(repo.DataAnalyticTopCategoriesParam{
		PaginationParams: models.PaginationParams{
			Limit:   10,
			Keyword: "Dresses",
		},
		Select:   []string{"sub_category"},
		DateFrom: 1702047296,
		DateTo:   1702875504,
		OrderBy:  "sold",
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_TopMoverProducts(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).TopMoverProducts(repo.DataAnalyticTopMoverProductsParam{
		PaginationParams: models.PaginationParams{
			Limit: 10,
		},
		Category:    "Women Clothing",
		SubCategory: "Women Dresses",
		Domain:      "shein",
		Month:       12,
		Year:        2023,
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_GetProduct(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).GetProduct(repo.DataAnalyticGetProductParam{
		ProductID: "cm1hrp3b2hjacnove90g",
	})

	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_GetProductChart(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).GetProductChart(repo.DataAnalyticGetDAProductChartParam{
		ProductID: "cm2k8oilk3mh15958m8g",
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 7,
		},
		//DateFrom:  1702047296,
		//DateTo:    1703211541,
		PredictNext: 5,
		PredictOn:   enums.PredictOnWeek,
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_GetProductClass(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).GetProductClassGroup(repo.GetAnalyticProductClassGroupParam{
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 30,
		},
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_GetProductTrending(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).GetProductTrendingGroup(repo.GetAnalyticProductTrendingGroupParam{
		PaginationParams: models.PaginationParams{
			Page:    1,
			Limit:   3,
			Keyword: "Fall",
		},
		SubCategories: []string{"Women Jackets"},
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_GetProductGroup(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewDataAnalyticRepo(app.AnalyticDB).ProductGroupURL(repo.DataAnalyticSearchProductsParam{
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 12,
		},
		// Domains: []string{"zara"},
	})

	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_PaginateNewUser(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).PaginateNewUsers(repo.DataAnalyticNewUsersParam{
		DateFrom: 1704870950,
		DateTo:   1706080571,
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_PaginateNewCatalogProduct(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).PaginateNewCatalogProduct(repo.DataAnalyticNewCatalogProductParam{
		DateFrom: 1701870950,
		DateTo:   1706080571,
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_PaginateInquiries(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).PaginateInquiries(repo.DataAnalyticInquiriesParam{
		DateFrom: 1201870950,
		DateTo:   1706080571,
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_PaginatePO(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).PaginatePurchaseOrders(repo.DataAnalyticPurchaseOrdersParam{
		DateFrom: 1201870950,
		DateTo:   1706080571,
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_PaginateBulkPO(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).PaginateBulkPurchaseOrders(repo.DataAnalyticBulkPurchaseOrdersParam{
		DateFrom: 1201870950,
		DateTo:   1706080571,
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_GetOpsBizPerformance(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).GetOpsBizPerformance(repo.DataAnalyticOpsBizPerformance{
		DateFrom: 0,
		DateTo:   1708498009,
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_BuyerDataAnalyticRFQ(t *testing.T) {
	var app = initApp("dev")
	claims := *models.NewJwtClaimsInfo().SetUserID("ci83c8djtqd6mfut3fmg")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).BuyerDataAnalyticRFQ(repo.BuyerDataAnalyticRFQParams{
		JwtClaimsInfo: claims,
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_BuyerDataAnalyticPendingTasks(t *testing.T) {
	var app = initApp("dev")
	claims := *models.NewJwtClaimsInfo().SetUserID("ci83c8djtqd6mfut3fmg")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).BuyerDataAnalyticPendingTasks(repo.BuyerDataAnalyticPendingTasksParams{
		JwtClaimsInfo: claims,
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 10,
		},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_BuyerDataAnalyticPendingPayments(t *testing.T) {
	var app = initApp("local")
	claims := *models.NewJwtClaimsInfo().SetUserID("ck0ovcunr62iue9co3ug")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).BuyerDataAnalyticPendingPaymentsV2(repo.BuyerDataAnalyticPendingPaymentsParams{
		JwtClaimsInfo: claims,
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 8,
		},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_BuyerDataAnalyticTotalStyleProduce(t *testing.T) {
	var app = initApp("dev")
	claims := *models.NewJwtClaimsInfo().SetUserID("cjvso05ooc2b8f45a1mg")
	result, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).BuyerDataAnalyticTotalStyleProduce(repo.BuyerDataAnalyticTotalStyleProduceParams{
		JwtClaimsInfo: claims,
		DateFrom:      0,
		DateTo:        1708501885,
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestDataAnalyticRepo_GetOneAnalyticProduct(t *testing.T) {
	var app = initApp("dev")
	claims := *models.NewJwtClaimsInfo().SetUserID("cjvso05ooc2b8f45a1mg")
	main, others, err := repo.NewDataAnalyticRepo(app.AnalyticDB).WithDB(app.DB).GetBest(repo.GetBestAnalyticProductParams{
		JwtClaimsInfo: claims,
	})
	assert.NoError(t, err)
	helper.PrintJSON(main)
	helper.PrintJSON(others)
}
