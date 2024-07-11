package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

func FixHubspoContacts() {
	var cfg = config.New("../deployment/config/prod/env.json")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	app.DB.AutoMigrate(&models.User{})
	var users []*models.User
	app.DB.Select("ID", "Email").Find(&users, "role IN ? AND COALESCE(hubspot_contact_id,'') = ''", []enums.Role{enums.RoleClient})

	var client = hubspot.New(app.Config)
	for index, user := range users {
		contacts, err := client.SearchContactsByEmail([]string{user.Email})
		if err != nil {
			fmt.Printf("Update err %d/%d id=%s err=%+v", index, len(users), user.ID, err)
			continue
		}

		if len(contacts.Results) == 0 {
			fmt.Printf("Update err %d/%d id=%s empty results", index, len(users), user.ID)
			continue
		}

		err = app.DB.Model(&models.User{}).Where("id = ?", user.ID).UpdateColumn("HubspotContactID", contacts.Results[0].ID).Error

		fmt.Printf("Update %d/%d id=%s err=%+v", index, len(users), user.ID, err)
	}

}
