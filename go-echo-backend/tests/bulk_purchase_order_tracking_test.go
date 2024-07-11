package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestBulkPurchaseOrderTrackingRepo_CreateBulkPurchaseOrderTracking(t *testing.T) {
	var app = initApp("local")

	var err = repo.NewBulkPurchaseOrderTrackingRepo(app.DB).CreateBulkPurchaseOrderTrackingTx(app.DB.DB, models.BulkPurchaseOrderTrackingCreateForm{
		PurchaseOrderID: "cl6vlpnbudgn5lhp5ft02",
		ActionType:      enums.BulkPoTrackingActionUpdateProduction,
		UserID:          "123",
	})

	assert.NoError(t, err)

	helper.PrintJSON("DONE")
}
