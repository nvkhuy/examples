package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestSEOTranslationRepo_FetchBuyerSeoTranslation(t *testing.T) {
	var app = initApp("prod")
	resp, err := repo.NewSeoTranslationRepo(app.DB).
		WithSheetAPI(app.SheetAPI).FetchSeoTranslation(&models.FetchSeoTranslationParams{
		Domain: enums.DomainBuyer,
	})
	assert.NoError(t, err)
	helper.PrintJSON(resp)
}

func TestSEOTranslationRepo_GetSEOTranslation(t *testing.T) {
	var app = initApp("prod")
	resp, err := repo.NewSeoTranslationRepo(app.DB).
		WithSheetAPI(app.SheetAPI).GetSEOTranslation(models.GetSEOTranslationForm{
		Domain: enums.DomainWebsite,
	})
	if err != nil {
		return
	}
	helper.PrintJSON(resp)
}
