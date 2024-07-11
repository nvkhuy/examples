package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
)

type AddressRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewAddressRepo(db *db.DB) *AddressRepo {
	return &AddressRepo{
		db:     db,
		logger: logger.New("repo/Address"),
	}
}

type GetAddressParams struct {
	AddressID string
}

func (r *AddressRepo) GetAddress(params GetAddressParams) (*models.Address, error) {
	var builder = queryfunc.NewAddressBuilder(queryfunc.AddressBuilderOptions{})

	var result models.Address
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("a.id = ?", params.AddressID)
		}).
		Limit(1).
		FirstFunc(&result)

	return &result, err
}
