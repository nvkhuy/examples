package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
)

func TestCategoryRepo_GetCategoryTree(t *testing.T) {
	var app = initApp("dev")

	var resp = repo.NewCategoryRepo(app.DB).GetCategoryTree(repo.PaginateCategoriesParams{})

	helper.PrintJSON(resp)
}

func TestCategoryRepo_GetChildCategoryIDs(t *testing.T) {
	var app = initApp()

	var resp = repo.NewCategoryRepo(app.DB).GetChildCategoryIDs("cg8lfou701noh4tunu00")

	helper.PrintJSON(resp)
}

func TestCategoryRepo_GetCategories(t *testing.T) {
	var app = initApp("local")

	var resp = repo.NewCategoryRepo(app.DB).GetCategories(repo.PaginateCategoriesParams{})

	helper.PrintJSON(resp)
}

func TestCategoryRepo_UpdateCategoryByID(t *testing.T) {
	var app = initApp("local")

	cate, err := repo.NewCategoryRepo(app.DB).UpdateCategoryByID(models.CategoryUpdateForm{
		CategoryID:    "cgsgj6brgh1n1gtlhqc0",
		TopProductIDs: []string{"cl8poi96ilsfgpaljn4g"},
	})
	if err != nil {
		return
	}
	helper.PrintJSON(cate)
}

func TestCategoryRepo_GenerateSlug(t *testing.T) {
	var app = initApp("local")
	err := repo.NewCategoryRepo(app.DB).GenerateSlug()
	assert.NoError(t, err)
}

func TestCategoryRepo_CreateCategory(t *testing.T) {
	var app = initApp("local")
	cate, err := repo.NewCategoryRepo(app.DB).CreateCategory(models.CategoryCreateForm{
		Name:             "New Cate 01",
		ParentCategoryID: "cg8lfl6701noh4tuntvg",
		CategoryType:     enums.CategoryOfProduct,
		TopProductIDs:    []string{"cl5ltse28clrj668m8ug"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, cate)
}
