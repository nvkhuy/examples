package repo

import (
	"database/sql"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollectionRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewCollectionRepo(db *db.DB) *CollectionRepo {
	return &CollectionRepo{
		db:     db,
		logger: logger.New("repo/Collection"),
	}
}

type PaginateCollectionParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *CollectionRepo) PaginateCollection(params PaginateCollectionParams) *query.Pagination {
	var builder = queryfunc.NewCollectionBuilder(queryfunc.CollectionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if strings.TrimSpace(params.Keyword) != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("c.name ILIKE @query", sql.Named("query", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *CollectionRepo) SearchCollection(params PaginateCollectionParams) []*models.Collection {
	var builder = queryfunc.NewCollectionBuilder(queryfunc.CollectionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Limit == 0 {
		params.Limit = 20
	}
	var result []*models.Collection
	_ = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
		}).
		Limit(params.Limit).
		FindFunc(&result)

	return result
}

func (r *CollectionRepo) CreateCollection(form models.CollectionCreateForm) (result *models.Collection, err error) {
	if err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		var Collection models.Collection
		err = copier.Copy(&Collection, &form)
		if err != nil {
			return eris.Wrap(err, "")
		}
		Collection.ProductIDs = form.ProductIds
		err = tx.Clauses(clause.Returning{}).Create(&Collection).Error
		if err != nil {
			return eris.Wrap(err, "")
		}
		result = &Collection
		var cps []models.CollectionProduct
		for _, pId := range form.ProductIds {
			cps = append(cps, models.CollectionProduct{
				ProductID:    pId,
				CollectionID: Collection.ID,
			})
		}
		if err = tx.CreateInBatches(cps, len(cps)).Error; err != nil {
			return eris.Wrap(err, "")
		}
		return
	}); err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	result, err = r.GetCollectionByID(result.ID, queryfunc.CollectionBuilderOptions{IsConsistentRead: true})
	return
}

func (r *CollectionRepo) GetCollectionByID(CollectionID string, options queryfunc.CollectionBuilderOptions) (*models.Collection, error) {
	var builder = queryfunc.NewCollectionBuilder(options)
	var Collection models.Collection
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("c.id = ?", CollectionID)
		}).
		FirstFunc(&Collection)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &Collection, nil
}

type GetCollectionDetailByIDParams struct {
	models.JwtClaimsInfo

	CollectionID string `param:"collection_id" validate:"required"`
}

func (r *CollectionRepo) GetCollectionDetailByID(params GetCollectionDetailByIDParams) (*models.Collection, error) {
	var builder = queryfunc.NewCollectionBuilder(queryfunc.CollectionBuilderOptions{})
	var Collection models.Collection
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("c.id = ?", params.CollectionID)
		}).
		FirstFunc(&Collection)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	// Populate products to group of collection
	groups, _ := NewCollectionProductGroupRepo(r.db).GetGroupByCollectionID(params.CollectionID, queryfunc.CollectionProductGroupBuilderOptions{
		IncludeProducts: true,
	})

	Collection.ProductGroups = groups

	return &Collection, nil
}

func (r *CollectionRepo) UpdateCollectionByID(form models.CollectionUpdateForm) (result *models.Collection, err error) {
	if err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		var update models.Collection

		err = copier.Copy(&update, &form)
		if err != nil {
			return
		}
		update.ProductIDs = form.ProductIds

		err = tx.Omit(clause.Associations).Model(&models.Collection{}).Where("id = ?", form.CollectionID).Updates(&update).Error
		if err != nil {
			return
		}

		if err = tx.Unscoped().Delete(&models.CollectionProduct{}, "collection_id = ?", form.CollectionID).Error; err != nil {
			return
		}

		var cps []models.CollectionProduct
		for _, pId := range form.ProductIds {
			cps = append(cps, models.CollectionProduct{
				ProductID:    pId,
				CollectionID: form.CollectionID,
			})
		}
		if err = tx.CreateInBatches(cps, len(cps)).Error; err != nil {
			return
		}

		return
	}); err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	result, err = r.GetCollectionByID(form.CollectionID, queryfunc.CollectionBuilderOptions{IsConsistentRead: true})
	return
}

func (r *CollectionRepo) GetCollectionByIDs(collectionIDs []string, options queryfunc.CollectionBuilderOptions) ([]*models.Collection, error) {
	var builder = queryfunc.NewCollectionBuilder(options)
	var collections []*models.Collection
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("c.id IN (?)", collectionIDs)
		}).
		FindFunc(&collections)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrCategoryNotFound
		}
		return nil, err
	}

	return collections, nil
}

func (r *CollectionRepo) ArchiveCollectionByID(CollectionID string) error {
	var err = r.db.Unscoped().Delete(&models.Collection{}, "id = ?", CollectionID).Error

	return err
}

func (r *CollectionRepo) AddProduct(CollectionID string, productIDs []string) (*models.Collection, error) {
	collection, err := r.GetCollectionByID(CollectionID, queryfunc.CollectionBuilderOptions{})
	if err != nil {
		return nil, err
	}

	for _, productID := range productIDs {
		var collectionProduct models.CollectionProduct
		r.db.Omit(clause.Associations).Where(models.CollectionProduct{ProductID: productID, CollectionID: CollectionID}).FirstOrCreate(&collectionProduct)
	}

	return collection, nil
}

func (r *CollectionRepo) RemoveProduct(collectionID string, productIDs []string) error {
	_, err := r.GetCollectionByID(collectionID, queryfunc.CollectionBuilderOptions{})
	if err != nil {
		return err
	}

	err = r.db.Unscoped().Delete(&models.CollectionProduct{}, "collection_id = ? AND product_id IN (?)", collectionID, productIDs).Error

	return err
}
