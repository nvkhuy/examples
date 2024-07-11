package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

const pwd = "Inflow@2023"

func TestAuth_Invite(t *testing.T) {
	var app = initApp("prod")
	customerio.New(app.Config)

	resp, err := repo.NewAuthRepo(app.DB).Register(models.RegisterForm{
		Email:     "guihuytu@gmail.com",
		Password:  pwd,
		FirstName: "Huy",
		LastName:  "Tú",
		BrandRegisterInfo: models.BrandRegisterInfo{
			BrandName: "Annie",
		},
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}
