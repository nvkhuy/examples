package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFabricCollectionRepo_Paginate(t *testing.T) {
	var app = initApp("local")
	result := repo.NewFabricCollectionRepo(app.DB).Paginate(repo.PaginateFabricCollectionParams{})
	helper.PrintJSON(result)
	return
}

func TestFabricCollectionRepo_Create(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewFabricCollectionRepo(app.DB).Create(repo.CreateFabricCollectionParams{
		FabricCollection: &models.FabricCollection{
			Name: "collection 01",
			Slug: "collection-01",
		},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
	return
}

func TestFabricCollectionRepo_AddFabric(t *testing.T) {
	var app = initApp("local")
	err := repo.NewFabricCollectionRepo(app.DB).AddFabric(repo.AddFabricToCollectionParams{
		ID:        "cm74t3qlk3msun9kgueg",
		FabricIDs: []string{"cm74taalk3msvdhc4hkg", "cm74tf2lk3msvufndnug"},
	})
	assert.NoError(t, err)
	return
}

func TestFabricCollectionRepo_RemoveFabric(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewFabricCollectionRepo(app.DB).RemoveFabric(repo.RemoveFabricFromCollectionParams{
		ID: "cm74t3qlk3msun9kgueg",
		//FabricIDs: []string{"cm74taalk3msvdhc4hkg", "cm74tf2lk3msvufndnug"},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
	return
}
