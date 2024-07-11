package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
)

type CollectionProductGroupRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewCollectionProductGroupRepo(db *db.DB) *CollectionProductGroupRepo {
	return &CollectionProductGroupRepo{
		db:     db,
		logger: logger.New("repo/CollectionProductGroup"),
	}
}

type PaginateCollectionProductGroupsParams struct {
	models.PaginationParams

	ForRole enums.Role
}

type CollectionProductGroupParams struct {
	ForRole enums.Role
}

func (r *CollectionProductGroupRepo) GetGroupByCollectionID(collectionID string, options queryfunc.CollectionProductGroupBuilderOptions) ([]*models.CollectionProductGroup, error) {
	var builder = queryfunc.NewCollectionProductGroupBuilder(options)
	var groups []*models.CollectionProductGroup
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("g.collection_id = ?", collectionID)
		}).
		FindFunc(&groups)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return groups, nil
}
