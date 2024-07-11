package tests

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"log"
	"testing"
)

func TestRWDFabricPriceRepo_FetchPrice(t *testing.T) {
	var app = initApp("local")
	url, err := repo.NewRWDFabricPriceRepo(app.DB).WithSheetAPI(app.SheetAPI).FetchPrice(&models.FetchRWDFabricPriceParams{
		MaterialType: enums.WovenMaterial,
		From:         "A3",
		To:           "L",
	})
	if err != nil {
		return
	}
	log.Println(url)
}

func TestRWDFabricPriceRepo_Paginate(t *testing.T) {
	var app = initApp("local")
	result := repo.NewRWDFabricPriceRepo(app.DB).Paginate(&models.PaginateRWDFabricPriceParams{
		FabricType:  aws.String(string(enums.KnitMaterial)),
		Material:    aws.String("Knit Polo Interlock filagen FL4885-305081"),
		Composition: aws.String("POLY/ POLY BLEND"),
		Weight:      aws.Float64(140),
		CutWidth:    aws.Float64(170),
	})
	helper.PrintJSON(result)
}

func TestRWDFabricPriceRepo_Vine(t *testing.T) {
	var app = initApp("local")
	result := repo.NewRWDFabricPriceRepo(app.DB).Vine(&models.PaginateRWDFabricPriceParams{
		FabricType: aws.String(string(enums.KnitMaterial)),
		Material:   aws.String("Knit Polo Interlock filagen FL4885-305081"),
	})
	helper.PrintJSON(result)
}
