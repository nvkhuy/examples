package repo

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"gorm.io/gorm/clause"
)

type TrendingRepo struct {
	db     *db.DB
	adb    *db.DB
	logger *logger.Logger
}

func NewTrendingRepo(db *db.DB) *TrendingRepo {
	return &TrendingRepo{
		db:     db,
		logger: logger.New("repo/trending"),
	}
}

func (r *TrendingRepo) WithADB(adb *db.DB) *TrendingRepo {
	r.adb = adb
	return r
}

type PaginateTrendingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	IDs      []string               `json:"ids" query:"ids" form:"ids" param:"ids"`
	Statuses []enums.TrendingStatus `json:"statuses,omitempty" param:"statuses" query:"statuses" form:"statuses"`
}

func (r *TrendingRepo) PaginateTrendings(params PaginateTrendingParams) *query.Pagination {
	var builder = queryfunc.NewTrendingBuilder(queryfunc.TrendingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{Role: params.GetRole()},
		Adb:                 r.adb,
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.IDs) > 0 {
				builder.Where("id IN ?", params.IDs)
			}
			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(t.name ILIKE @keyword)", sql.Named("keyword", q))
			}
			if len(params.Statuses) > 0 {
				builder.Where("t.status IN ?", params.Statuses)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		OrderBy("t.created_at DESC").
		PagingFunc()

	return result
}

type ListTrendingStatsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	IDs      []string               `json:"ids" query:"ids" form:"ids" param:"ids"`
	Statuses []enums.TrendingStatus `json:"statuses,omitempty" param:"statuses" query:"statuses" form:"statuses"`
}

func (r *TrendingRepo) ListTrendingStats(params ListTrendingStatsParams) (trendings []*models.Trending) {
	r.db.Model(&models.Trending{}).Select("id", "name", "product_trending_ids").Find(&trendings)
	for _, t := range trendings {
		t.Total = len(t.ProductTrendingIDs)
		t.ProductTrendingIDs = []string{}
	}
	return trendings
}

type CreateTrendingsParams struct {
	models.JwtClaimsInfo
	Name               string               `json:"name" param:"name" query:"name" form:"name"`
	Description        string               `json:"description" param:"description" query:"description" form:"description"`
	ProductTrendingIDs pq.StringArray       `json:"product_trending_ids" param:"product_trending_ids" query:"product_trending_ids" form:"product_trending_ids"`
	Status             enums.TrendingStatus `json:"status,omitempty" param:"status" query:"status" form:"status"`
	CoverAttachment    *models.Attachment   `json:"cover_attachment,omitempty" param:"cover_attachment" query:"cover_attachment" form:"cover_attachment"`
	IsAutoCreate       *bool                `json:"is_auto_create,omitempty" param:"is_auto_create" query:"is_auto_create" form:"is_auto_create"`
}

func (r *TrendingRepo) Create(params CreateTrendingsParams) (result *models.Trending, err error) {
	var trending models.Trending
	err = copier.Copy(&trending, &params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}
	if params.Name != "" {
		var _name string
		if _name, err = r.PredictName(params.ProductTrendingIDs); err == nil && _name != "" {
			trending.Name = _name
		}
	}

	if err = r.db.Model(&models.Trending{}).
		Clauses(clause.Returning{}).
		Create(&trending).Error; err != nil {
		return
	}

	result, err = r.Get(GetTrendingParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		ID:            trending.ID,
	})
	return
}

type GetTrendingParams struct {
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
}

func (r *TrendingRepo) Get(params GetTrendingParams) (result *models.Trending, err error) {
	builder := queryfunc.NewTrendingBuilder(queryfunc.TrendingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{Role: params.GetRole()},
		Adb:                 r.adb,
	})
	var trending models.Trending
	err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ID != "" {
				builder.Where("id = ?", params.ID)
			}
		}).
		FirstFunc(&trending)

	result = &trending
	return
}

type UpdateTrendingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ID                 string               `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	Name               string               `json:"name" param:"name" query:"name" form:"name"`
	Description        string               `json:"description" param:"description" query:"description" form:"description"`
	ProductTrendingIDs pq.StringArray       `json:"product_trending_ids" param:"product_trending_ids" query:"product_trending_ids" form:"product_trending_ids"`
	Status             enums.TrendingStatus `json:"status,omitempty" param:"status" query:"status" form:"status"`
	CoverAttachment    *models.Attachment   `json:"cover_attachment,omitempty" param:"cover_attachment" query:"cover_attachment" form:"cover_attachment"`
}

func (r *TrendingRepo) Update(params UpdateTrendingParams) (result *models.Trending, err error) {
	var trending models.Trending
	err = copier.Copy(&trending, &params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	var _name string
	if _name, err = r.PredictName(params.ProductTrendingIDs); err == nil && _name != "" {
		trending.Name = _name
	}

	if err = r.db.Model(&trending).
		Clauses(clause.Returning{}).
		Where("id = ?", params.ID).
		Updates(&trending).Error; err != nil {
		return
	}

	result, err = r.Get(GetTrendingParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		ID:            trending.ID,
	})
	return
}

type DeleteTrendingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
}

func (r *TrendingRepo) Delete(params DeleteTrendingParams) (err error) {
	err = r.db.Unscoped().Delete(&models.Trending{}, "id = ?", params.ID).Error
	return
}

type AddProductToTrendingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ID                 string         `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	ProductTrendingIDs pq.StringArray `json:"product_trending_ids" param:"product_trending_ids" query:"product_trending_ids" form:"product_trending_ids"`
}

func (r *TrendingRepo) AddProductToTrending(params AddProductToTrendingParams) (result *models.Trending, err error) {
	trending, err := r.Get(GetTrendingParams{JwtClaimsInfo: params.JwtClaimsInfo, ID: params.ID})
	if err != nil {
		return
	}
	trending.ProductTrendingIDs = append(trending.ProductTrendingIDs, params.ProductTrendingIDs...)
	trending.ProductTrendingIDs = lo.Uniq(trending.ProductTrendingIDs)

	var updates = models.Trending{
		Name:               trending.Name,
		ProductTrendingIDs: trending.ProductTrendingIDs,
	}

	if err = r.db.Model(&models.Trending{}).Where("id = ?", params.ID).Updates(updates).Error; err != nil {
		return
	}
	result, err = r.Get(GetTrendingParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		ID:            trending.ID,
	})
	return
}

type RemoveProductFromTrendingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ID                 string         `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	ProductTrendingIDs pq.StringArray `json:"product_trending_ids" param:"product_trending_ids" query:"product_trending_ids" form:"product_trending_ids"`
}

func (r *TrendingRepo) RemoveProductFromTrending(params RemoveProductFromTrendingParams) (result *models.Trending, err error) {
	trending, err := r.Get(GetTrendingParams{JwtClaimsInfo: params.JwtClaimsInfo, ID: params.ID})
	if err != nil {
		return
	}
	var productTrendingIds pq.StringArray
	for _, productId := range trending.ProductTrendingIDs {
		exits := lo.Contains(params.ProductTrendingIDs, productId)
		if !exits {
			productTrendingIds = append(productTrendingIds, productId)
		}
	}

	var _name string
	if _name, err = r.PredictName(productTrendingIds); err == nil && _name != "" {
		trending.Name = _name
	}
	var updates = models.Trending{
		Name:               trending.Name,
		ProductTrendingIDs: trending.ProductTrendingIDs,
	}

	if err = r.db.Model(&models.Trending{}).Where("id = ?", params.ID).Updates(updates).Error; err != nil {
		return
	}
	result, err = r.Get(GetTrendingParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		ID:            trending.ID,
	})
	return
}

func (r *TrendingRepo) AutoCreate() (err error) {
	var products []models.AnalyticProductTrending
	err = r.adb.Select("id", "metadata").Where("domain = ?", enums.TagWalk).Find(&products).Error
	if err != nil {
		return
	}
	collections := make(map[string]*models.Trending)

	for _, p := range products {
		if p.Metadata == nil {
			continue
		}
		colName := p.Metadata["collection_name"].(string)
		colImage := p.Metadata["collection_image"].(string)
		if c, ok := collections[colName]; ok {
			c.ProductTrendingIDs = append(c.ProductTrendingIDs, p.ID)
		} else {
			collections[colName] = &models.Trending{
				Name:   colName,
				Status: enums.TrendingStatusDraft,
				CoverAttachment: &models.Attachment{
					FileURL:      colImage,
					FileKey:      colImage,
					ThumbnailURL: colImage,
					ContentType:  "image/jpeg",
				},
				IsAutoCreate: aws.Bool(true),
			}
		}
	}
	for _, v := range collections {
		r.db.Model(&models.Trending{}).Create(v)
	}

	return
}

func (r *TrendingRepo) ReverseOrder() (err error) {
	var trendings []models.Trending
	err = r.db.Select("id", "created_at").Where("is_auto_create = ?", true).Order("created_at").Find(&trendings).Error
	if err != nil {
		return
	}
	n := len(trendings) - 1
	for i, t := range trendings {
		fmt.Println(t.ID)
		r.db.Model(&models.Trending{}).Where("id = ?", t.ID).Update("created_at", trendings[n-i].CreatedAt)
	}

	return
}

func (r *TrendingRepo) AssignCollectionURL() (err error) {
	var trendings []models.Trending
	err = r.db.Select("id", "product_trending_ids").Where("is_auto_create = ?", true).Order("created_at").Find(&trendings).Error
	if err != nil {
		return
	}
	for _, t := range trendings {
		if len(t.ProductTrendingIDs) == 0 {
			continue
		}

		var product models.AnalyticProductTrending
		r.adb.Model(&models.AnalyticProductTrending{}).Where("id = ?", t.ProductTrendingIDs[0]).First(&product)
		if product.Metadata == nil {
			continue
		}

		collectionURL := product.Metadata["collection_url"]
		if collectionURL != "" {
			r.db.Model(&models.Trending{}).Where("id = ?", t.ID).Update("collection_url", collectionURL)
		}
	}
	return
}

func (r *TrendingRepo) PredictName(ids []string) (name string, err error) {
	if len(ids) == 0 {
		return
	}
	var products []models.AnalyticProductTrending
	err = r.adb.Select("id", "name", "sub_category", "category", "domain").Where("id IN ?", pq.StringArray(ids)).Find(&products).Error
	if err != nil {
		return
	}
	for _, p := range products {
		if p.Domain == enums.TagWalk {
			split := strings.Split(p.Name, " ")
			if len(split) > 2 {
				name = fmt.Sprintf("%s %s", split[0], split[1])
				return
			}
		}
	}
	return
}
