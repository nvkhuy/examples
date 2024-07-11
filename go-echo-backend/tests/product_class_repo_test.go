package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"strings"
	"testing"
)

func TestProductClass_Create(t *testing.T) {
	var app = initApp("local")
	var err = repo.NewProductClassRepo(app.DB).Create(repo.ProductClassCreate{
		ProductID: "2",
		Class:     "women-denim",
		Conf:      12.423,
	})
	assert.NoError(t, err)
}

func TestProductClass_List(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewProductClassRepo(app.DB).List(repo.ProductClassList{
		ProductID: "1",
	})
	if err != nil {
		return
	}
	log.Println(len(result))
}

func TestProductClass_Upsert(t *testing.T) {
	var app = initApp("local")
	var err = repo.NewProductClassRepo(app.DB).Upsert(repo.ProductClassUpsert{
		ProductID: "2",
		Class:     "women-denim",
		Conf:      15.423,
	})
	assert.NoError(t, err)
}

func TestProductClass_BatchUpsert(t *testing.T) {
	var app = initApp("local")
	var err = repo.NewProductClassRepo(app.DB).BatchUpsert(repo.ProductClassBatchUpsert{
		ProductClasses: []models.ProductClass{
			{
				ProductID: "1",
				Class:     "women-denim",
				Conf:      65.423,
			},
			{
				ProductID: "1",
				Class:     "women-activewear",
				Conf:      31.423,
			},
			{
				ProductID: "2",
				Class:     "women-denim",
				Conf:      27.423,
			},
		},
	})
	assert.NoError(t, err)
}

func TestProductClass_Assign(t *testing.T) {
	var app = initApp("local")
	type Class struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
	}
	var sliceProductCategory []Class
	var err = app.DB.Model(&models.Product{}).
		Select("products.id", "c.slug").
		Joins("inner join categories c ON (c.parent_category_id != '' AND c.id = products.category_id)").
		Where(gorm.Expr("coalesce(category_id, '') != ''")).Find(&sliceProductCategory).Error
	if err != nil {
		return
	}
	var productClasses []models.ProductClass
	for _, pc := range sliceProductCategory {
		className := strings.ReplaceAll(pc.Slug, "-and", "")
		productClasses = append(productClasses, models.ProductClass{
			ProductID: pc.ID,
			Class:     className,
			Conf:      100,
		})
	}
	err = repo.NewProductClassRepo(app.DB).BatchUpsert(repo.ProductClassBatchUpsert{
		ProductClasses: productClasses,
	})
	assert.NoError(t, err)
}

func TestProductClass_ClassifyProduct(t *testing.T) {
	var app = initApp("local")
	var err = repo.NewProductClassRepo(app.DB).ClassifyProduct()
	assert.NoError(t, err)
}
