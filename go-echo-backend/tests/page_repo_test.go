package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestPage_PageCatalog(t *testing.T) {
	var app = initApp("local")

	resp, err := repo.NewPageRepo(app.DB).PageCatalog("cgmjp7r8f9m7n67fun9g")

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}

func TestPage_GetPageDetailByID(t *testing.T) {
	var app = initApp("local")

	resp, err := repo.NewPageRepo(app.DB).GetPageDetailByID(repo.GetPageDetailByIDParams{
		PageID: "cgmjp7r8f9m7n67fun9g",
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}
