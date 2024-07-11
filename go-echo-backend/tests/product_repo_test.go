package tests

import (
	"log"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func TestProductRepo_PaginateProducts(t *testing.T) {
	var app = initApp("dev")

	var resp = repo.NewProductRepo(app.DB).PaginateProducts(repo.PaginateProductParams{
		CategorySlug:    "men",
		SubCategorySlug: "men-and-shirts",
		PaginationParams: models.PaginationParams{
			Limit: 5,
		},
	})

	helper.PrintJSON(resp)
}

func TestProductRepo_GetProductByPageSectionType(t *testing.T) {
	var app = initApp("local")
	sectionType := string(enums.PageSectionTypeCatalogDrop)
	var resp, err = repo.NewProductRepo(app.DB).GetProductByPageSectionType(sectionType, queryfunc.ProductBuilderOptions{})
	if err != nil {
		return
	}

	helper.PrintJSON(resp)
}

func TestProductRepo_GetProduct(t *testing.T) {
	var app = initApp("local")
	var resp, err = repo.NewProductRepo(app.DB).GetProduct(repo.GetProductParams{
		SlugID: "test-luu",
	})
	if err != nil {
		return
	}

	helper.PrintJSON(resp)
}

func TestProductRepo_CreateProduct(t *testing.T) {
	var app = initApp("local")
	var resp, err = repo.NewProductRepo(app.DB).CreateProduct(models.ProductCreateForm{
		Name: "fasdfsadf",
	})
	if err != nil {
		return
	}

	helper.PrintJSON(resp)
}

func TestProductRepo_UpdateProduct(t *testing.T) {
	var app = initApp("local")
	var resp, err = repo.NewProductRepo(app.DB).UpdateProduct(models.ProductUpdateForm{
		ProductID: "cl49tgfskmp3mvc5rdpg",
		Name:      "fasdfsadf",
		Variants: []*models.VariantAttributeUpdateForm{
			{
				Model: models.Model{
					ID: "cl49tjvskmp3mvc5rdr0",
				},
				Title:       "T1",
				ProductName: "Sample Variant 1 - 01",
			},
		},
	})
	if err != nil {
		return
	}

	helper.PrintJSON(resp)
}

func TestProductRepo_GetQRCode(t *testing.T) {
	var app = initApp("local")
	logo := app.Config.QRCodeLogoURL
	result, err := repo.NewProductRepo(app.DB).GetQRCode(models.GetProductQRCodeParams{
		Logo:      logo,
		Bucket:    app.Config.AWSS3StorageBucket,
		ProductId: "cl49tgfskmp3mvc5rdpg",
		URL:       "https://www.joininflow.io/products/cl49tgfskmp3mvc5rdpg",
		Override:  false,
	})
	log.Println(result, err)
}

func TestProductRepo_ExportExcel(t *testing.T) {
	var app = initApp("local")
	resp, err := repo.NewProductRepo(app.DB).ExportExcel(repo.PaginateProductParams{})
	if err != nil {
		return
	}
	log.Println(resp)
}

func TestProductRepo_PaginateRecommendations(t *testing.T) {
	var app = initApp("dev")

	var resp = repo.NewProductRepo(app.DB).PaginateRecommendations(repo.PaginateRecommendProductParams{
		RecommendProductID: "ciqc9v9dbgeapa7oq6qg",
		PaginationParams: models.PaginationParams{
			Limit: 10,
		},
	})

	helper.PrintJSON(resp)
}
