package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFabricRepo_Create(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewFabricRepo(app.DB).Create(repo.CreateFabricParams{
		Fabric: &models.Fabric{
			FabricType:      "Denim 02",
			FabricID:        "denim-02",
			FabricWeight:    12,
			MOQ:             10,
			Colors:          []string{"black", "white"},
			ManufacturerIDs: []string{"cl1i06sdc1giiu6epaeg", "cl12kl0q98bfgt0h37vg"},
			FabricCostings: &models.FabricCostings{
				{
					From:           100,
					To:             200,
					Price:          price.NewFromFloat(4.99),
					ProcessingTime: "3-40 days",
				},
			},
		},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestFabricRepo_Update(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewFabricRepo(app.DB).Update(repo.UpdateFabricParams{
		ID: "cm5t602lk3mhskj6kumg",
		Fabric: &models.Fabric{
			FabricType:      "Denim 01",
			FabricID:        "denim-01",
			FabricWeight:    12,
			MOQ:             10,
			Colors:          []string{"black", "white", "orange"},
			ManufacturerIDs: []string{"cl1i06sdc1giiu6epaeg", "cl12kl0q98bfgt0h37vg"},
			FabricCostings: &models.FabricCostings{
				{
					From:           100,
					To:             200,
					Price:          price.NewFromFloat(4.99),
					ProcessingTime: "3-40 days",
				},
			},
		},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestFabricRepo_Paginate(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewFabricRepo(app.DB).Paginate(repo.PaginateFabricParams{})
	helper.PrintJSON(result)
	return
}

func TestFabricRepo_Details(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewFabricRepo(app.DB).Details(repo.DetailsFabricParams{
		ID: "cm781qilk3mg03dpte50",
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
	return
}

func TestFabricRepo_PatchReferenceID(t *testing.T) {
	var app = initApp("dev")
	err := repo.NewFabricRepo(app.DB).PatchFabricReferenceID(repo.PatchFabricReferenceIDParams{})
	assert.NoError(t, err)
	return
}
