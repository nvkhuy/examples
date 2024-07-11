package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

func FixUserFabricType() {
	var cfg = config.New("../deployment/config/dev/secrets.env")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var profiles []*models.BusinessProfile
	app.DB.Select("*").Find(&profiles)

	for index, profile := range profiles {

		var fabricTypes []string

		if profile.MillFabricTypes != nil {
			for _, item := range *profile.MillFabricTypes {
				fabricTypes = append(fabricTypes, item.FabricValue)
			}

			var profileUpdate = models.BusinessProfile{
				FlatMillFabricTypes: fabricTypes,
			}

			var err = app.DB.Model(&models.BusinessProfile{}).Where("id = ?", profile.ID).Updates(profileUpdate)

			fmt.Printf("Update %d/%d id=%s err=%+v", index, len(profiles), profile.ID, err)
		}

	}

}
