package seeder

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/gorm/clause"
)

func (s *Seeder) SeedAccounts() *Seeder {
	var users = []models.User{
		{
			FirstName: "Super Admin",
			LastName:  "Super Admin",
			Email:     "admin@joininflow.io",
			Password:  "1234qwer",
			Role:      enums.RoleSuperAdmin,
		},
	}
	var err = s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error
	if err != nil {
		s.db.CustomLogger.Errorf("Seed accounts error=%+v", err)
	}

	return s
}

func (s *Seeder) SeedLocalAccounts() *Seeder {
	var users = []models.User{
		{
			FirstName: "Super Admin",
			LastName:  "Super Admin",
			Email:     "admin@joininflow.io",
			Password:  "1234qwer",
			Role:      enums.RoleSuperAdmin,
		},
		{
			FirstName:     "Seller",
			LastName:      "Seller",
			Email:         "seller@joininflow.io",
			Password:      "1234qwer",
			Role:          enums.RoleSeller,
			AccountStatus: enums.AccountStatusActive,
		},
		{
			FirstName:     "Buyer",
			LastName:      "Buyer",
			Email:         "buyer@joininflow.io",
			Password:      "1234qwer",
			Role:          enums.RoleClient,
			AccountStatus: enums.AccountStatusActive,
		},
	}
	var err = s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error
	if err != nil {
		s.db.CustomLogger.Errorf("Seed accounts error=%+v", err)
	}

	return s
}
