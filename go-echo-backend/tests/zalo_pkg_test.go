package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/zalo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZaloPkg_CreateZaloConfig(t *testing.T) {
	var app = initApp("dev")
	err := zalo.SaveZaloConfig(app.DB, zalo.OARefreshTokenResponse{
		AccessToken:  "86ukG17mM04ZRYX54lDP3Y1hBayKwLPA904wMqhwA6GU4H1uKD5fQGCgVYj2ztKDUpH_9IFGQnKI7N869uPZFN1JI7y3kKrKPNrVGWAGKMypMq4jQxGn9mLtNJjtdYih1cGz4MwvB2rRUoChAuvQ7MDRLW4ll5WRTs0q1ncUHK8jLcn4IPuKRXrQ0sX-a7TAIbDcI7FA46WZ73PMR88_HbqV4bTgYnHf6qK6LNYk1KqHLn1INfmROZz6DIT9do1a8mC6U4JR0a4s40S0NTP-4nuyS50e-d9XP1rmMGNx0Nfe9pfx7TrFP702MI4zcdKnHrfu91NERmzvDHStD_8w1ruqJZGiq4KgOGvjR2NLHYbk4oen7OajC2jYAmv7cn0gCNiN1rlY4pOpFYOUKFqV0s0S4ZjwOm2iGN8Kx5i0",
		RefreshToken: "UZJsQw9f0s5vDF1keMe2RdenzcIdS78t277-GjHtNWXB8Sjqm7O-FaKlg4_0TWCuPL-FHTHp62zIHPDTune61Niuj6EhRpCrCKIwGeTH7JKO3OvOeHTv3Ha0yGQeEq1qB3hw3fKGQG4v4U99aKnJ4qX0tcZFRYGES4UX0Cb21Lrm7Uf9rI1YGL1syKdiNLi2CNwXKfTYL0q_EDv_f6Dw0HqXu061BsLWBYcHR_KH31ua1ef1ho04G1CkXqo03Zy49YAdL806C2Gu3wDLZbOK7YXCXM6xTnSUVrs4Q_Ty7XWlQPL6zZKLIdngupZjQ5LnMpMZU_ar43GW3kfOkmORKXCRkZEnG2128aYa0QHlDc43Ny0LY45fUL4MxNZL5XuDMYoHVVm92Z15DebxymaPB-blmZcdTG1W",
		ExpiresIn:    "90000",
	})
	assert.NoError(t, err)
	//helper.PrintJSON(result)
}

func TestZaloPkg_OARefreshToken(t *testing.T) {
	var app = initApp("dev")
	result, err := zalo.OARefreshToken(app.DB, zalo.OARefreshTokenParams{})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}

func TestZaloPkg_SendZNS(t *testing.T) {
	var app = initApp("local")
	result, err := zalo.SendZNS(app.DB, zalo.SendZNSParams{
		Phone:      "84938294687",
		Mode:       "development",
		TemplateID: "311342",
		TemplateData: map[string]interface{}{
			"name":         "Huy",
			"order_code":   "CODE10001",
			"phone_number": "84938294687",
			"start_date":   "20/03/2020",
			"price":        500000,
			"status":       "delivered",
			"date":         "26/01/2024",
		},
		TrackingID: "tracking_id",
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}
