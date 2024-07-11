package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func TestStatsRepo_StatsBuyers(t *testing.T) {
	var app = initApp()

	var resp = repo.NewStatsRepo(app.DB).StatsBuyers(repo.StatsBuyersParams{})

	helper.PrintJSON(resp)
}

func TestStatsRepo_StatsCategories(t *testing.T) {
	var app = initApp()

	var resp = repo.NewStatsRepo(app.DB).StatsCategories(repo.StatsProductsParams{})

	helper.PrintJSON(resp)
}

func TestStatsRepo_StatsProducts(t *testing.T) {
	var app = initApp()

	var resp = repo.NewStatsRepo(app.DB).StatsProducts(repo.StatsProductsParams{})

	helper.PrintJSON(resp)
}

func TestStatsRepo_StatsSuppliers(t *testing.T) {
	var app = initApp()

	var resp = repo.NewStatsRepo(app.DB).StatsSuppliers(repo.StatsSuppliersParams{})

	helper.PrintJSON(resp)
}
