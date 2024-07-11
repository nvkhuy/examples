package seeder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/gorm/clause"
)

func (s *Seeder) SeedCategories() *Seeder {
	var categories = []models.Category{
		{
			Name:             "Women",
			ParentCategoryID: aws.String(""),
			Slug:             "women",
			CategoryType:     enums.CategoryOfProduct,
		},
		{
			Name:             "Men",
			ParentCategoryID: aws.String(""),
			Slug:             "men",
			CategoryType:     enums.CategoryOfProduct,
		},
	}
	var err = s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&categories).Error
	if err != nil {
		s.db.CustomLogger.Errorf("Seed categories error=%+v", err)
	}

	return s
}
