package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

func FixProductAttribute() {
	var cfg = config.New("../deployment/config/prod/secrets.env")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var products []*models.Product
	app.DB.Select("*").Find(&products)

	for index, product := range products {
		var productUpdate = models.Product{
			ProductAttributeMetas: models.ProductAttributeMetas{},
		}

		var productAttributes []*models.ProductAttribute
		app.DB.Where("product_id = ?", product.ID).Select("*").Find(&productAttributes)

		for _, productAttribute := range productAttributes {
			if productAttribute.Name == "Color" {
				for _, itemVal := range productAttribute.Values {
					productUpdate.ProductAttributeMetas = append(productUpdate.ProductAttributeMetas, &models.ProductAttributeMeta{Attribute: enums.ProductAttributeColor, Value: *itemVal})
				}
			}
			if productAttribute.Name == "Size" {
				for _, itemVal := range productAttribute.Values {
					productUpdate.ProductAttributeMetas = append(productUpdate.ProductAttributeMetas, &models.ProductAttributeMeta{Attribute: enums.ProductAttributeSize, Value: *itemVal})
				}
			}
		}

		var err = app.DB.Model(&models.Product{}).Where("id = ?", product.ID).Updates(productUpdate)

		fmt.Printf("Update %d/%d id=%s err=%+v", index, len(products), product.ID, err)
	}

}
