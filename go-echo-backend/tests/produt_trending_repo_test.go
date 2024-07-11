package tests

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductTrendingRepo_Paginate(t *testing.T) {
	var app = initApp("local")
	result := repo.NewProductTrending(app.AnalyticDB).WithDB(app.DB).
		PaginateProductTrendings(repo.PaginateProductTrendingParams{
			ProductIDs: []string{"cpa5813b2hj733ak4b00"},
		})
	helper.PrintJSON(result)
}

func TestProductTrendingRepo_ListProductTrendingDomain(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewProductTrending(app.AnalyticDB).WithDB(app.DB).ListProductTrendingDomain(repo.ListProductTrendingDomainParams{})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestProductTrendingRepo_ListProductTrendingCategory(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewProductTrending(app.AnalyticDB).WithDB(app.DB).ListProductTrendingCategory(repo.ListProductTrendingCategoryParams{})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestProductTrendingRepo_PaginateGroup(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewProductTrending(app.AnalyticDB).WithDB(app.DB).
		PaginateProductTrendingGroup(repo.PaginateProductTrendingGroupParams{})
	helper.PrintJSON(result)
}

func TestProductTrendingRepo_Create(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewProductTrending(app.AnalyticDB).WithDB(app.DB).
		CreateProductTrending(repo.CreateProductTrendingParams{
			Products: []*models.AnalyticProductTrending{
				{
					Name:   "Casablanca Fall/Winter 2022 — Look 69",
					Domain: enums.Inflow,
					Size:   aws.String("X"),
					Color:  aws.String("White"),
				},
				{
					Name:   "Casablanca Fall/Winter 2022 — Look 69",
					Domain: enums.Inflow,
					Size:   aws.String("XL"),
					Color:  aws.String("White"),
				},
			},
		})
	if err != nil {
		return
	}
	helper.PrintJSON(result)
}
