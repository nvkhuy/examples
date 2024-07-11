package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrendingRepo_Create(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewTrendingRepo(app.DB).WithADB(app.AnalyticDB).
		Create(repo.CreateTrendingsParams{
			Name:               "New Trending 02",
			ProductTrendingIDs: []string{"cpemslqlk3mh3nilacr0"},
		})
	if err != nil {
		return
	}
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestTrendingRepo_Paginate(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewTrendingRepo(app.DB).WithADB(app.AnalyticDB).
		PaginateTrendings(repo.PaginateTrendingParams{
			PaginationParams: models.PaginationParams{
				Page:  0,
				Limit: 12,
			},
		})
	helper.PrintJSON(result)
}

func TestTrendingRepo_Get(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewTrendingRepo(app.DB).WithADB(app.AnalyticDB).
		Get(repo.GetTrendingParams{
			ID: "cpanrkalk3mgbq6cao10",
		})
	if err != nil {
		return
	}
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestTrendingRepo_RemoveProductFromTrending(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewTrendingRepo(app.DB).WithADB(app.AnalyticDB).
		RemoveProductFromTrending(repo.RemoveProductFromTrendingParams{
			ID:                 "cpanrkalk3mgbq6cao10",
			ProductTrendingIDs: []string{"cpa5813b2hj733ak4b00"},
		})
	if err != nil {
		return
	}
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestTrendingRepo_AddProductToTrending(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewTrendingRepo(app.DB).WithADB(app.AnalyticDB).
		AddProductToTrending(repo.AddProductToTrendingParams{
			ID:                 "cpanrkalk3mgbq6cao10",
			ProductTrendingIDs: []string{"cpa5813b2hj733ak4b00"},
		})
	if err != nil {
		return
	}
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestTrendingRepo_AutoCreate(t *testing.T) {
	var app = initApp("dev")
	err := repo.NewTrendingRepo(app.DB).WithADB(app.AnalyticDB).
		AutoCreate()
	if err != nil {
		return
	}
	assert.NoError(t, err)
}

func TestTrendingRepo_ReverseOrder(t *testing.T) {
	var app = initApp("dev")
	err := repo.NewTrendingRepo(app.DB).WithADB(app.AnalyticDB).
		ReverseOrder()
	if err != nil {
		return
	}
	assert.NoError(t, err)
}

func TestTrendingRepo_AssignCollectionURL(t *testing.T) {
	var app = initApp("dev")
	err := repo.NewTrendingRepo(app.DB).WithADB(app.AnalyticDB).
		AssignCollectionURL()
	if err != nil {
		return
	}
	assert.NoError(t, err)
}
