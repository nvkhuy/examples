package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/lib/pq"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func TestCollectionRepo_CreateCollection(t *testing.T) {
	var app = initApp("local")
	payload := models.CollectionCreateForm{
		Name:        "c1",
		Description: "d1",
		ProductIds:  pq.StringArray{"cl4s2kccr609p5jsrih0", "cl5g3lum2bb1hhprc4d0"},
	}
	helper.PrintJSON(payload)
	resp, err := repo.NewCollectionRepo(app.DB).CreateCollection(payload)
	if err != nil {
		return
	}

	helper.PrintJSON(resp)
}

func TestCollectionRepo_UpdateCollection(t *testing.T) {
	var app = initApp("local")
	payload := models.CollectionUpdateForm{
		CollectionID: "cl65p3vskmp7b1n7pd0g",
		Name:         "c2",
		Description:  "d1",
		ProductIds:   []string{"cl4qjfolurutnq6i2lcg"},
	}
	helper.PrintJSON(payload)
	resp, err := repo.NewCollectionRepo(app.DB).UpdateCollectionByID(payload)
	if err != nil {
		return
	}

	helper.PrintJSON(resp)
}

func TestCollectionRepo_GetCollectionDetailByID(t *testing.T) {
	var app = initApp("local")
	payload := repo.GetCollectionDetailByIDParams{
		CollectionID: "cl66cufskmp7odv1jhmg",
	}
	resp, err := repo.NewCollectionRepo(app.DB).GetCollectionDetailByID(payload)
	if err != nil {
		return
	}

	helper.PrintJSON(resp)
}
