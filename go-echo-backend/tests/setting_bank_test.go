package tests

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"testing"
)

func TestSettingBankRepo_PaginateSettingBank(t *testing.T) {
	var app = initApp("local")
	result := repo.NewSettingBankRepo(app.DB).PaginateSettingBank(repo.PaginateSettingBankParams{
		PaginationParams: models.PaginationParams{
			Limit: 10,
		},
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
	})

	helper.PrintJSON(result)
}

func TestSettingBankRepo_CreateSettingBanks(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewSettingBankRepo(app.DB).CreateSettingBanks(models.SettingBanksForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		CountryCode:   "VN",
		Content:       "vietnam abcdxyz",
		IsDisabled:    aws.Bool(true),
	})
	if err != nil {
		fmt.Println(err)
	}
	helper.PrintJSON(result)
}

func TestSettingBankRepo_UpdateSettingBanks(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewSettingBankRepo(app.DB).UpdateSettingBanks(models.SettingBanksForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		ID:            "ckf2ti7skmp625khtgd0",
		CountryCode:   "VN",
		Content:       "vietnam abcdxyz 2",
		IsDisabled:    aws.Bool(true),
	})
	if err != nil {
		fmt.Println(err)
	}
	helper.PrintJSON(result)
}

func TestSettingTaxRepo_DeleteSettingBank(t *testing.T) {
	var app = initApp("local")
	err := repo.NewSettingBankRepo(app.DB).DeleteSettingBanks(models.DeleteSettingBanksForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		ID:            "ckf2uvfskmp64nec9a1g",
	})
	if err != nil {
		fmt.Println(err)
	}
	helper.PrintJSON("Ok")
}

func TestSettingTaxRepo_DeleteSettingBanksByCountryCode(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewSettingBankRepo(app.DB).DeleteSettingBanksByCountryCode(models.DeleteSettingBanksByCountryCodeForm{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		CountryCode:   "sg",
	})
	if err != nil {
		fmt.Println(err)
	}
	helper.PrintJSON(result)
}
