package repo

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FabricRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewFabricRepo(db *db.DB) *FabricRepo {
	return &FabricRepo{
		db:     db,
		logger: logger.New("repo/fabric"),
	}
}

type CreateFabricParams struct {
	models.JwtClaimsInfo
	*models.Fabric
	FabricCollectionIDs []string `json:"fabric_collection_ids" param:"fabric_collection_ids" query:"fabric_collection_ids" form:"fabric_collection_ids"`
}

func (r *FabricRepo) Create(params CreateFabricParams) (result *models.Fabric, err error) {
	if err = r.db.Model(&models.Fabric{}).Clauses(clause.Returning{}).Create(params.Fabric).Error; err != nil {
		return
	}
	result = params.Fabric
	if err = r.AddFabricToCollections(AddFabricToCollectionsParams{ID: params.ID, FabricCollectionIDs: params.FabricCollectionIDs, IsOverride: true}); err != nil {
		return
	}
	if result.FabricCollections, err = NewFabricCollectionRepo(r.db).ListFabricCollections(params.ID); err != nil {
		return
	}
	result.Manufacturers, err = r.ListFabricManufacturers(result.ManufacturerIDs)
	return
}

type UpdateFabricParams struct {
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	*models.Fabric
	FabricCollectionIDs []string `json:"fabric_collection_ids" param:"fabric_collection_ids" query:"fabric_collection_ids" form:"fabric_collection_ids"`
}

func (r *FabricRepo) Update(params UpdateFabricParams) (result *models.Fabric, err error) {
	if params.ID == "" {
		err = errors.New("empty fabric id")
		return
	}
	if params.Fabric == nil {
		params.Fabric = &models.Fabric{}
	}
	params.Fabric.ID = params.ID
	if err = r.db.Model(params.Fabric).Clauses(clause.Returning{}).Where("id = ?", params.ID).Updates(params.Fabric).Error; err != nil {
		return
	}
	result = params.Fabric
	if err = r.AddFabricToCollections(AddFabricToCollectionsParams{ID: params.ID, FabricCollectionIDs: params.FabricCollectionIDs, IsOverride: true}); err != nil {
		return
	}
	if result.FabricCollections, err = NewFabricCollectionRepo(r.db).ListFabricCollections(params.ID); err != nil {
		return
	}
	result.Manufacturers, err = r.ListFabricManufacturers(result.ManufacturerIDs)
	return
}

type PaginateFabricParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	FabricIDs []string `json:"fabric_ids" param:"fabric_ids" query:"fabric_ids" form:"fabric_ids"`
}

func (r *FabricRepo) Paginate(params PaginateFabricParams) *query.Pagination {
	if params.Limit == 0 {
		params.Limit = 20
	}
	return query.New(r.db, queryfunc.NewFabricBuilder(queryfunc.FabricBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			log.Println(params)
			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("fb.fabric_type ILIKE @keyword", sql.Named("keyword", q))
			}
			if len(params.FabricIDs) > 0 {
				builder.Where("fb.id IN ?", params.FabricIDs)
			}
		}).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
}

type DetailsFabricParams struct {
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
}

func (r *FabricRepo) Details(params DetailsFabricParams) (result *models.Fabric, err error) {
	var fb models.Fabric
	err = r.db.Model(&models.Fabric{}).Where("id = ?", params.ID).First(&fb).Error
	if err != nil {
		return
	}
	result = &fb
	if result == nil {
		return
	}
	if result.FabricCollections, err = NewFabricCollectionRepo(r.db).ListFabricCollections(params.ID); err != nil {
		return
	}
	result.Manufacturers, err = r.ListFabricManufacturers(result.ManufacturerIDs)
	return
}

type DeleteFabricParams struct {
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
}

func (r *FabricRepo) Delete(params DeleteFabricParams) (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		if err = r.db.Unscoped().Delete(&models.Fabric{}, "id = ?", params.ID).Error; err != nil {
			return
		}
		err = r.db.Unscoped().Delete(&models.FabricInCollection{}, "fabric_id = ?", params.ID).Error
		return
	})
	return
}

func (r *FabricRepo) ListFabricManufacturers(manufacturerIDs []string) (manufacturers []models.User, err error) {
	var users []models.User
	if err = r.db.Model(&models.User{}).Where("id IN ?", manufacturerIDs).Find(&users).Error; err != nil {
		return
	}
	manufacturers = users
	return
}

type AddFabricToCollectionsParams struct {
	models.JwtClaimsInfo
	ID                  string   `json:"id" param:"id" query:"id" form:"id" validate:"required"` // fabric id
	FabricCollectionIDs []string `json:"fabric_collection_ids" param:"fabric_collection_ids" query:"fabric_collection_ids" form:"fabric_collection_ids"`
	IsOverride          bool     `json:"is_override"`
}

func (r *FabricRepo) AddFabricToCollections(params AddFabricToCollectionsParams) (err error) {
	if params.ID == "" {
		err = errors.New("empty fabric id")
		return
	}
	if len(params.FabricCollectionIDs) == 0 {
		return
	}
	var relations []models.FabricInCollection
	for _, collectionID := range params.FabricCollectionIDs {
		relations = append(relations, models.FabricInCollection{
			FabricID:           params.ID,
			FabricCollectionID: collectionID,
		})
	}
	if params.IsOverride {
		if err = r.db.Unscoped().Where("fabric_id = ?", params.ID).Delete(&models.FabricInCollection{}).Error; err != nil {
			return
		}
	}
	err = r.db.Model(&models.FabricInCollection{}).Create(&relations).Error
	return
}

type PatchFabricReferenceIDParams struct{}

func (r *FabricRepo) PatchFabricReferenceID(params PatchFabricReferenceIDParams) (err error) {
	var fabrics []models.Fabric
	if err = r.db.Model(&models.Fabric{}).Where("reference_id is NULL").Find(&fabrics).Error; err != nil {
		return
	}
	ids := make(map[string]bool)
	uniqueID := func() string {
		id := helper.GenerateFabricReferenceID()
		_, ok := ids[id]
		for ok {
			id = helper.GenerateFabricReferenceID()
			_, ok = ids[id]
		}
		ids[id] = true
		return id
	}

	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		for _, fb := range fabrics {
			updates := models.Fabric{
				ReferenceID: uniqueID(),
			}
			if err = tx.Model(&models.Fabric{}).Where("id = ?", fb.ID).Updates(&updates).Error; err != nil {
				return
			}
		}
		return
	})

	return
}
