package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func adminLogin(db *db.DB) (*models.LoginResponse, error) {
	// seed account
	cred := models.LoginEmailForm{
		Email:         "admin@joininflow.io",
		Password:      "1234qwer",
		IsAdminPortal: true,
	}

	loginResp, err := repo.NewAuthRepo(db).LoginEmail(cred)
	if err != nil {
		return nil, err
	}
	return loginResp, nil
}
