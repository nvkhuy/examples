package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSettingDocRepo_CreateSettingDoc(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewSettingDocRepo(app.DB).Create(repo.SettingDocCreateParams{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetUserID("cl5lei628clrj668m800"),
		Type:          enums.SettingDocNDAType,
		Document: &models.Attachment{
			FileKey: "uploads/media/nda_doc_1.png",
		},
	})

	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestSettingDocRepo_UpdateSettingDoc(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewSettingDocRepo(app.DB).Update(repo.SettingDocUpdateParams{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetUserID("cl5lei628clrj668m800"),
		Type:          enums.SettingDocNDAType,
		Document: &models.Attachment{
			FileKey: "uploads/media/nda_doc_2.png",
		},
	})

	assert.NoError(t, err)
	helper.PrintJSON(result)
}
