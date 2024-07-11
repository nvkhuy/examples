package repo

import (
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingSizeRepo struct {
	db *db.DB
}

func NewSettingSizeRepo(db *db.DB) *SettingSizeRepo {
	return &SettingSizeRepo{
		db: db,
	}
}

type PaginateSettingSizesParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *SettingSizeRepo) PaginateSettingSizes(params PaginateSettingSizesParams) *query.Pagination {
	return query.New(r.db, queryfunc.NewSettingSizeBuilder(queryfunc.SettingSizeBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
}

type GetSettingSizesParams struct {
	models.JwtClaimsInfo

	CountryCodes []enums.CountryCode `json:"country_code" query:"country_code" param:"country_code"`
}

func (r *SettingSizeRepo) GetSettingSizes(params GetSettingSizesParams) ([]*models.SettingSize, error) {
	var records []*models.SettingSize
	var err = query.New(r.db, queryfunc.NewSettingSizeBuilder(queryfunc.SettingSizeBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).WhereFunc(func(builder *query.Builder) {

	}).
		FindFunc(&records)
	return records, err
}

func (r *SettingSizeRepo) CreateSettingSizes(form models.SettingSizeCreateForm) ([]*models.SettingSize, error) {
	var sizes []*models.SettingSize

	for _, v := range form.SizeNames {
		sizes = append(sizes, &models.SettingSize{
			Type: form.Type,
			Name: v,
		})
	}

	var err = r.db.Omit(clause.Associations).Create(&sizes).Error
	if err != nil {
		return nil, err
	}

	return sizes, err
}

func (r *SettingSizeRepo) GetSettingSize(form models.SettingSizeIDForm) (*models.SettingSize, error) {
	var size models.SettingSize

	var err = r.db.First(&size, "id = ?", form.SizeID).Error

	return &size, err
}

func (r *SettingSizeRepo) UpdateSettingSize(form models.SettingSizeUpdateForm) (*models.SettingSize, error) {
	var err = r.db.Model(&models.SettingSize{}).Where("id = ?", form.SizeID).UpdateColumn("Name", form.Name).Error
	if err != nil {
		return nil, err
	}
	return r.GetSettingSize(models.SettingSizeIDForm{
		SizeID: form.SizeID,
	})
}

func (r *SettingSizeRepo) UpdateSettingSizes(form models.SettingSizesUpdateForm) ([]*models.SettingSize, error) {
	var sizes []*models.SettingSize
	var updatedSizeNames []string

	var err = r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now().Unix()
		for _, name := range form.SizeNames {
			var size = models.SettingSize{
				Name: name,
				Type: form.Type,
			}
			var err = tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(func() map[string]interface{} {
					now = now - 1
					return map[string]interface{}{"updated_at": now}
				}()),
				Columns: []clause.Column{
					{Name: "type"},
					{Name: "name"},
				},
			}, clause.Returning{}).
				Create(&size).Error
			if err != nil {
				return err
			}

			updatedSizeNames = append(updatedSizeNames, size.Name)
			sizes = append(sizes, &size)
		}

		if len(updatedSizeNames) > 0 {
			return tx.Delete(&models.SettingSize{}, "type = ? AND name NOT IN ?", form.Type, updatedSizeNames).Error
		}

		return nil
	})

	return sizes, err
}

func (r *SettingSizeRepo) DeleteSettingSize(form models.SettingSizeDeleteForm) error {
	return r.db.Delete(&models.SettingSize{}, "id = ?", form.SizeID).Error
}

func (r *SettingSizeRepo) DeleteSettingSizeType(form models.SettingSizeDeleteTypeForm) error {
	return r.db.Delete(&models.SettingSize{}, "type = ?", form.Type).Error
}

func (r *SettingSizeRepo) UpdateSettingSizeType(form models.SettingSizeUpdateTypeForm) (updates models.SettingSize, err error) {
	err = r.db.Model(&updates).Clauses(clause.Returning{}).
		Where("type = ?", form.Type).UpdateColumn("type", form.NewType).Error
	return
}

func (r *SettingSizeRepo) GetSettingSizeType(form models.GetSettingSizeTypeForm) (*models.SettingSize, error) {
	var result models.SettingSize
	var err = query.New(r.db, queryfunc.NewSettingSizeBuilder(queryfunc.SettingSizeBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: form.GetRole(),
		},
	})).
		Limit(1).
		Where("ss.type = ?", form.Type).
		FirstFunc(&result)

	return &result, err
}
