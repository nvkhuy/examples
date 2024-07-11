package seeder

import "github.com/engineeringinflow/inflow-backend/pkg/db"

// Seeder struct
type Seeder struct {
	db *db.DB
}

// New instance
func New(db *db.DB) *Seeder {
	return &Seeder{
		db: db,
	}
}
