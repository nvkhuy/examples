package repo

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FabricCollectionRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewFabricCollectionRepo(db *db.DB) *FabricCollectionRepo {
	return &FabricCollectionRepo{
		db:     db,
		logger: logger.New("repo/fabric_collection"),
	}
}

type CreateFabricCollectionParams struct {
	models.JwtClaimsInfo
	*models.FabricCollection
	FabricIDs []string `json:"fabric_ids" param:"fabric_ids" query:"fabric_ids" form:"fabric_ids"`
}

func (r *FabricCollectionRepo) Create(params CreateFabricCollectionParams) (result *models.FabricCollection, err error) {
	if err = r.db.Model(&params.FabricCollection).Clauses(clause.Returning{}).Create(params.FabricCollection).Error; err != nil {
		return
	}
	if err = r.AddFabric(AddFabricToCollectionParams{
		ID:        params.FabricCollection.ID,
		FabricIDs: params.FabricIDs,
	}); err != nil {
		return
	}
	result = params.FabricCollection
	result.Fabrics, err = r.ListFabricInCollection(result.ID)
	return
}

type UpdateFabricCollectionParams struct {
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	*models.FabricCollection
	FabricIDs []string `json:"fabric_ids" param:"fabric_ids" query:"fabric_ids" form:"fabric_ids"`
}

func (r *FabricCollectionRepo) Update(params UpdateFabricCollectionParams) (result *models.FabricCollection, err error) {
	if params.ID == "" {
		err = errors.New("empty fabric id")
		return
	}
	if params.FabricCollection == nil {
		return
	}
	var updates = models.FabricCollection{
		Name: params.FabricCollection.Name,
		VI:   params.VI,
	}
	params.FabricCollection.ID = params.ID
	if err = r.db.Model(&updates).Clauses(clause.Returning{}).Where("id = ?", params.ID).Updates(&updates).Error; err != nil {
		return
	}
	if err = r.AddFabric(AddFabricToCollectionParams{
		ID:         params.ID,
		FabricIDs:  params.FabricIDs,
		IsOverride: true,
	}); err != nil {
		return
	}
	result = &updates
	result.Fabrics, err = r.ListFabricInCollection(result.ID)
	return
}

type AddFabricToCollectionParams struct {
	models.JwtClaimsInfo
	ID         string   `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	FabricIDs  []string `json:"fabric_ids" param:"fabric_ids" query:"fabric_ids" form:"fabric_ids" validate:"required"`
	IsOverride bool     `json:"is_override"`
}

func (r *FabricCollectionRepo) AddFabric(params AddFabricToCollectionParams) (err error) {
	if params.ID == "" {
		err = errors.New("empty fabric collection id")
		return
	}
	if len(params.FabricIDs) == 0 {
		return
	}
	var relations []models.FabricInCollection
	for _, fabricID := range params.FabricIDs {
		relations = append(relations, models.FabricInCollection{
			FabricID:           fabricID,
			FabricCollectionID: params.ID,
		})
	}
	if params.IsOverride {
		if err = r.db.Unscoped().Where("fabric_collection_id = ?", params.ID).Delete(&models.FabricInCollection{}).Error; err != nil {
			return
		}
	}
	err = r.db.Model(&models.FabricInCollection{}).Create(&relations).Error
	return
}

type RemoveFabricFromCollectionParams struct {
	models.JwtClaimsInfo
	ID        string   `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	FabricIDs []string `json:"fabric_ids" param:"fabric_ids" query:"fabric_ids" form:"fabric_ids" validate:"required"`
}

func (r *FabricCollectionRepo) RemoveFabric(params RemoveFabricFromCollectionParams) (result *models.FabricCollection, err error) {
	if params.ID == "" {
		err = errors.New("empty fabric collection id")
		return
	}

	q := r.db.Unscoped().Where("fabric_collection_id = ?", params.ID)
	if len(params.FabricIDs) > 0 {
		q.Where("fabric_id IN ?", params.FabricIDs)
	}
	err = q.Delete(&models.FabricInCollection{}).Error
	return
}

type DeleteFabricCollectionParams struct {
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
}

func (r *FabricCollectionRepo) Delete(params DeleteFabricCollectionParams) (err error) {
	if params.ID == "" {
		err = errors.New("empty fabric collection id")
		return
	}

	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		if err = r.db.Unscoped().Delete(&models.FabricCollection{}, "id = ?", params.ID).Error; err != nil {
			return
		}
		err = r.db.Unscoped().Delete(&models.FabricInCollection{}, "fabric_collection_id = ?", params.ID).Error
		return
	})

	return
}

type PaginateFabricCollectionParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	IDs         []string `json:"ids" param:"ids" query:"ids" form:"ids"`
	WithFabrics bool     `json:"with_fabrics,omitempty" param:"with_fabrics" query:"with_fabrics" form:"with_fabrics"`
}

func (r *FabricCollectionRepo) Paginate(params PaginateFabricCollectionParams) *query.Pagination {
	if params.Limit == 0 {
		params.Limit = 20
	}
	return query.New(r.db, queryfunc.NewFabricCollectionBuilder(queryfunc.FabricCollectionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			if len(params.IDs) > 0 {
				builder.Where("id IN ?", params.IDs)
			}
			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(fc.name ILIKE @keyword)", sql.Named("keyword", q))
			}
		}).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
}

type DetailsFabricCollectionParams struct {
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
}

func (r *FabricCollectionRepo) Details(params DetailsFabricCollectionParams) (result *models.FabricCollection, err error) {
	var fc models.FabricCollection
	_ = query.New(r.db, queryfunc.NewFabricCollectionBuilder(queryfunc.FabricCollectionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).WhereFunc(func(builder *query.Builder) {
		builder.Where("id = ?", params.ID)
	}).FirstFunc(&fc)

	result = &fc
	return
}

func (r *FabricCollectionRepo) ListFabricCollections(fabricID string) (result []models.FabricCollection, err error) {
	var relations []models.FabricInCollection
	if err = r.db.Model(&models.FabricInCollection{}).Where("fabric_id = ?", fabricID).Find(&relations).Error; err != nil {
		return
	}
	var collectionIDs []string
	for _, v := range relations {
		collectionIDs = append(collectionIDs, v.FabricCollectionID)
	}
	var collections []models.FabricCollection
	if err = r.db.Model(&models.FabricCollection{}).Where("id IN ?", collectionIDs).Find(&collections).Error; err != nil {
		return
	}
	result = collections
	return
}

func (r *FabricCollectionRepo) ListFabricInCollection(collectionID string) (result []models.Fabric, err error) {
	var relations []models.FabricInCollection
	if err = r.db.Model(&models.FabricInCollection{}).Where("fabric_collection_id = ?", collectionID).Find(&relations).Error; err != nil {
		return
	}
	var fabricIDs []string
	for _, v := range relations {
		fabricIDs = append(fabricIDs, v.FabricID)
	}
	var fabrics []models.Fabric
	if err = r.db.Model(&models.Fabric{}).Where("id IN ?", fabricIDs).Find(&fabrics).Error; err != nil {
		return
	}
	result = fabrics
	return
}
