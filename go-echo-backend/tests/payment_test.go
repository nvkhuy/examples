package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/stretchr/testify/assert"
)

func TestPayment_Login(t *testing.T) {
	var app = initApp()

	stripehelper.New(app.Config)

	result, err := repo.NewUserRepo(app.DB).SetupIntentForClientSecert("cgqdri7f2jf5imkm9gcg")
	assert.NoError(t, err)

	helper.PrintJSON(result)
}
