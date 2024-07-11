package main

import (
	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/runner"
)

func FixBlurhash() {
	var cfg = config.New("../deployment/config/prod/env.json")
	var logger = logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var products []*models.Product
	app.DB.Find(&products, "attachments is not null")

	var runner = runner.New(5)

	for index, product := range products {
		index := index
		product := product

		runner.Submit(func() {
			if product.Attachments != nil && len(*product.Attachments) > 0 {
				var attachments models.Attachments
				for _, attachment := range *product.Attachments {
					attachment.GetBlurhash()
					attachments = append(attachments, attachment)

				}
				var err = app.DB.Model(&models.Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{"attachments": attachments}).Error
				logger.Debugf("Update %d/%d product attachments err=%+v", index, len(products), err)
			}
		})

	}

	runner.Wait()

}
