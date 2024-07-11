package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func TestSettingSEORepo_PaginateSettingSEO(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewSettingSEORepo(app.DB).PaginateSettingSEOLanguageGroup(repo.PaginateSettingSEOParams{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
	})

	helper.PrintJSON(result)
}
