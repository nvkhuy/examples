package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSysNotification_Create(t *testing.T) {
	var app = initApp("local")

	claims := models.NewJwtClaimsInfo().SetUserID("cl9hs5onmcbsbpjm11mg")
	resp, err := repo.NewSysNotificationRepo(app.DB).Create(repo.CreateSysNotificationsParams{
		JwtClaimsInfo: *claims,
		SysNotification: models.SysNotification{
			Name:    "tell all",
			Message: "i'm here",
		},
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}

func TestSysNotification_MarkSeen(t *testing.T) {
	var app = initApp("local")

	claims := models.NewJwtClaimsInfo().SetUserID("cl9hs5onmcbsbpjm11mg")
	err := repo.NewSysNotificationRepo(app.DB).MarkSeen(repo.MarkSeenSysNotificationParams{
		JwtClaimsInfo:  *claims,
		NotificationID: "cm6fpfilk3mj1cdlh1jg",
	})

	assert.NoError(t, err)
}

func TestSysNotification_MarkSeenAll(t *testing.T) {
	var app = initApp("local")

	claims := models.NewJwtClaimsInfo().SetUserID("cl9hs5onmcbsbpjm11mg")
	err := repo.NewSysNotificationRepo(app.DB).MarkSeenAll(repo.MarkSeenAllSysNotificationParams{
		JwtClaimsInfo: *claims,
	})

	assert.NoError(t, err)
}

func TestSysNotification_Paginate(t *testing.T) {
	var app = initApp("local")

	claims := models.NewJwtClaimsInfo().SetUserID("cl9hs5onmcbsbpjm11mg")
	result := repo.NewSysNotificationRepo(app.DB).Paginate(repo.PaginateSysNotificationsParams{
		JwtClaimsInfo: *claims,
	})

	helper.PrintJSON(result)
}
