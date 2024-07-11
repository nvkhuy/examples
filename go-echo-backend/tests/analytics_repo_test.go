package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func TestAnalyticsRepo_PaginatePotentialOverdueInquiries(t *testing.T) {
	var app = initApp()
	repo.NewAnalyticsRepo(app.DB).PaginatePotentialOverdueInquiries(repo.PaginatePotentialOverdueInquiriesParams{
		JwtClaimsInfo:    *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin).SetUserID("cg5anr2llkm6ctpvq8k0"),
		PaginationParams: models.PaginationParams{},
	})
}
