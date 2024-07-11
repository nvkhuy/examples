package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/ai"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"
	"log"
	"path/filepath"
)

type ProductClassRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewProductClassRepo(db *db.DB) *ProductClassRepo {
	return &ProductClassRepo{
		db:     db,
		logger: logger.New("repo/ProductClass"),
	}
}

type ProductClassCreate struct {
	models.JwtClaimsInfo
	ProductID string  `json:"product_id,omitempty" validate:"required"`
	Class     string  `json:"class,omitempty" validate:"required"`
	Conf      float64 `json:"conf,omitempty"`
}

func (r *ProductClassRepo) Create(params ProductClassCreate) (err error) {
	var create models.ProductClass
	if err = copier.Copy(&create, &params); err != nil {
		return err
	}
	err = r.db.Model(&models.ProductClass{}).Create(&create).Error
	return
}

type ProductClassUpdate struct {
	models.JwtClaimsInfo
	ProductID string  `json:"product_id,omitempty" validate:"required"`
	Class     string  `json:"class,omitempty" validate:"required"`
	Conf      float64 `json:"conf,omitempty"`
}

func (r *ProductClassRepo) Update(params ProductClassUpdate) (err error) {
	var updates models.ProductClass
	if err = copier.Copy(&updates, &params); err != nil {
		return err
	}
	err = r.db.Model(&models.ProductClass{}).
		Where("product_id = ?", params.ProductID).
		Where("class = ?", params.Class).
		Updates(&updates).Error
	return
}

type ProductClassUpsert struct {
	models.JwtClaimsInfo
	ProductID string  `json:"product_id,omitempty" validate:"required"`
	Class     string  `json:"class,omitempty" validate:"required"`
	Conf      float64 `json:"conf,omitempty"`
}

func (r *ProductClassRepo) Upsert(params ProductClassUpsert) (err error) {
	var upserts models.ProductClass
	if err = copier.Copy(&upserts, &params); err != nil {
		return err
	}
	r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}, {Name: "class"}},
		DoUpdates: clause.AssignmentColumns([]string{"conf"}),
	}).Create(&upserts)
	return
}

type ProductClassBatchUpsert struct {
	models.JwtClaimsInfo
	ProductClasses []models.ProductClass `json:"product_classes,omitempty"`
}

func (r *ProductClassRepo) BatchUpsert(params ProductClassBatchUpsert) (err error) {
	r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}, {Name: "class"}},
		DoUpdates: clause.AssignmentColumns([]string{"conf"}),
	}).CreateInBatches(&params.ProductClasses, 100)
	return
}

type ProductClassList struct {
	models.JwtClaimsInfo
	models.PaginationParams
	ProductID string   `json:"product_id,omitempty" param:"product_id" query:"product_id"`
	Classes   []string `json:"classes,omitempty" param:"classes" query:"classes"`
}

func (r *ProductClassRepo) List(params ProductClassList) (classes []*models.ProductClass, err error) {
	var builder = queryfunc.NewProductClassBuilder(queryfunc.ProductClassBuilderOptions{})
	err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ProductID != "" {
				builder.Where("product_id = ?", params.ProductID)
			}
			if len(params.Classes) > 0 {
				builder.Where("class IN ?", params.Classes)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		Find(&classes)
	return
}

func (r *ProductClassRepo) ClassifyProduct() (err error) {
	var products []models.Product
	err = r.db.Model(&models.Product{}).
		Select("id", "attachments").
		Where("attachments is not null").
		Find(&products).Error
	if err != nil {
		return
	}
	for _, p := range products {
		if p.Attachments == nil {
			continue
		}
		var (
			classifyParams  []ai.ImageClassifyParams
			classifyResults []ai.ImageClassifyResponse
		)
		for _, att := range *p.Attachments {
			if att.FileURL == "" {
				continue
			}
			ext := filepath.Ext(att.FileURL)
			if !(ext == ".jpg" || ext == ".jpeg") {
				continue
			}
			classifyParams = append(classifyParams, ai.ImageClassifyParams{
				Image:      att.FileURL,
				Size:       640,
				Confidence: 0.5,
				Overlap:    0,
			})
		}
		if len(classifyParams) == 0 {
			continue
		}

		classifyResults, err = ai.ClassifyMultiImage(classifyParams)
		if err != nil {
			return
		}
		classes := make(map[string]float64)
		for _, result := range classifyResults {
			for _, pred := range result.Predictions {
				classes[pred.Class] = pred.Confidence
			}
		}
		upsertParams := ProductClassBatchUpsert{}
		for class, conf := range classes {
			upsertParams.ProductClasses = append(upsertParams.ProductClasses, models.ProductClass{
				ProductID: p.ID,
				Class:     class,
				Conf:      conf,
			})
		}
		_ = r.BatchUpsert(upsertParams)
		log.Print("Insert ProductID: ", p.ID)
	}
	return
}
