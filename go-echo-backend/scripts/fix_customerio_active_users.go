package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func FixCustomerioActiveUsers() {
	var cfg = config.New("../deployment/config/prod/secrets.env")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var customerioClient = customerio.New(cfg)
	var users []*models.User
	app.DB.Select("ID").Find(&users)

	for index, user := range users {
		data, err := repo.NewUserRepo(app.DB).GetCustomerIOUser(user.ID)
		if err != nil {
			continue
		}

		err = customerioClient.Track.Identify(user.ID, data.GetCustomerIOMetadata(nil))

		fmt.Printf("Update %d/%d id=%s err=%+v", index, len(users), user.ID, err)
	}

}
