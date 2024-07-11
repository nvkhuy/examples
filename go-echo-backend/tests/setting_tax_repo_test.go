package tests

import (
	"fmt"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func TestSettingTaxRepo_PaginateSettingTax(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewSettingTaxRepo(app.DB).PaginateSettingTaxes(repo.PaginateSettingTaxParams{
		PaginationParams: models.PaginationParams{
			Limit: 10,
		},
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
	})

	helper.PrintJSON(result)
}

func TestSettingTaxRepo_CreateSettingTax(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewSettingTaxRepo(app.DB).CreateSettingTax(models.CreateSettingTaxForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		CountryCode:   "US",
		TaxPercentage: 47,
		DateAffected:  1696303500,
	})
	if err != nil {
		fmt.Println(err)
	}
	helper.PrintJSON(result)
}

func TestSettingTaxRepo_UpdateSettingTax(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewSettingTaxRepo(app.DB).UpdateSettingTax(models.UpdateSettingTaxForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		TaxID:         "ckdq36vskmp2jbvfmgfg",
		CountryCode:   "VN",
		TaxPercentage: 100,
		DateAffected:  1696305968,
	})
	if err != nil {
		fmt.Println(err)
	}
	helper.PrintJSON(result)
}

func TestSettingTaxRepo_GetAffectedSettingTax(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewSettingTaxRepo(app.DB).GetAffectedSettingTax(models.GetAffectedSettingTaxForm{
		CurrencyCode: enums.VND,
	})
	if err != nil {
		fmt.Println(err)
	}
	helper.PrintJSON(result)
}
