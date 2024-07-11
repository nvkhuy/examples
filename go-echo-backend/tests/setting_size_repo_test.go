package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestSettingSizeRepo_PaginateSettingSize(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewSettingSizeRepo(app.DB).PaginateSettingSizes(repo.PaginateSettingSizesParams{
		PaginationParams: models.PaginationParams{
			Limit: 10,
		},
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
	})

	helper.PrintJSON(result)
}
func TestSettingTaxRepo_CreateSettingSize(t *testing.T) {
	var app = initApp("dev")
	app.DB.AutoMigrate(&models.SettingSize{})

	result, err := repo.NewSettingSizeRepo(app.DB).CreateSettingSizes(models.SettingSizeCreateForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		Type:          "letters",
		SizeNames: []string{
			"small (s)",
		},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestSettingTaxRepo_UpdateSettingSize(t *testing.T) {
	var app = initApp("dev")
	app.DB.AutoMigrate(&models.SettingSize{})

	result, err := repo.NewSettingSizeRepo(app.DB).UpdateSettingSizes(models.SettingSizesUpdateForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		Type:          "test",
		SizeNames:     []string{"S", "M"},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestSettingTaxRepo_DeleteSettingSize(t *testing.T) {
	var app = initApp("dev")
	err := repo.NewSettingSizeRepo(app.DB).DeleteSettingSize(models.SettingSizeDeleteForm{
		SettingSizeIDForm: models.SettingSizeIDForm{
			SizeID:        "ckftrtberad5cd7mtpo0",
			JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		},
	})
	assert.NoError(t, err)
}

func TestSettingTaxRepo_DeleteSettingSizesByType(t *testing.T) {
	var app = initApp("dev")
	var err = repo.NewSettingSizeRepo(app.DB).DeleteSettingSizeType(models.SettingSizeDeleteTypeForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		Type:          "asdfasdf",
	})
	assert.NoError(t, err)
}

func TestSettingTaxRepo_UpdateSettingSizeType(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewSettingSizeRepo(app.DB).UpdateSettingSizeType(models.SettingSizeUpdateTypeForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		Type:          "UK02",
		NewType:       "UK01",
	})
	if err != nil {
		return
	}
	helper.PrintJSON(result)
}

func TestSettingSizeRepo_GetSettingSizeType(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewSettingSizeRepo(app.DB).GetSettingSizeType(models.GetSettingSizeTypeForm{
		Type: "US",
	})
	if err != nil {
		return
	}
	helper.PrintJSON(result)
}
