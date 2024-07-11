package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"github.com/thaitanloi365/go-utils/random"

	"gorm.io/gorm"
)

func TestDocumentRepo__AdminCreateDocumentCategoryBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	repo := repo.NewDocumentCategoryRepo(app.DB)

	testCases := []struct {
		name      string
		req       func(uniStr string) models.CreateDocumentCategoryRequest
		setup     func(t *testing.T, uniStr string)
		expect    func(uniStr string) models.DocumentCategory
		expectErr error
	}{
		{
			name: "should create document category correctly",
			req: func(uniStr string) models.CreateDocumentCategoryRequest {
				return models.CreateDocumentCategoryRequest{
					Name: fmt.Sprintf("cate name %s", uniStr),
				}
			},
			expect: func(uniStr string) models.DocumentCategory {
				return models.DocumentCategory{
					Name: fmt.Sprintf("cate name %s", uniStr),
					Slug: fmt.Sprintf("cate-name-%s", strings.ToLower(uniStr)),
				}
			},
		},
		{
			name: "should return error when create document category name with existing name",
			req: func(uniStr string) models.CreateDocumentCategoryRequest {
				return models.CreateDocumentCategoryRequest{
					Name: fmt.Sprintf("same cate name %s", uniStr),
				}
			},
			setup: func(t *testing.T, uniStr string) {
				_, err := repo.CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
					Name: fmt.Sprintf("same cate name %s", uniStr),
				})
				assert.NoError(t, err)
			},
			expectErr: errs.New(http.StatusBadRequest, "Category name is already existed"),
		},
	}
	for _, testCase := range testCases {
		randomStr := random.String(8, random.Alphabet)
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setup != nil {
				testCase.setup(t, randomStr)
			}
			req := testCase.req(randomStr)
			resp, err := repo.CreateDocumentCategory(&req)

			if testCase.expectErr != nil {
				assert.Equal(t, testCase.expectErr, err)
				return
			}
			assert.NoError(t, err)
			expect := testCase.expect(randomStr)
			expect.ID = resp.ID
			expect.CreatedAt = resp.CreatedAt
			expect.UpdatedAt = resp.UpdatedAt
			dbDocument := models.DocumentCategory{}
			err = app.DB.First(&dbDocument, "id = ?", resp.ID).Error

			assert.NoError(t, err)
			assert.Equal(t, expect, dbDocument)
		})
	}
}

func TestDocumentRepo__AdminUpdateDocumentCategoryBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	repo := repo.NewDocumentCategoryRepo(app.DB)

	categoryPayload := &models.CreateDocumentCategoryRequest{
		Name: fmt.Sprintf("exist cate %s", random.String(8, random.Alphabet)),
	}
	createdCategory, err := repo.CreateDocumentCategory(categoryPayload)
	assert.NoError(t, err)

	testCases := []struct {
		name      string
		req       func(uniStr string) models.UpdateDocumentCategoryRequest
		setup     func(t *testing.T, uniStr string)
		expect    func(uniStr string) models.DocumentCategory
		expectErr error
	}{
		{
			name: "should update document category correctly",
			req: func(uniStr string) models.UpdateDocumentCategoryRequest {
				return models.UpdateDocumentCategoryRequest{
					DocumentCategoryID: createdCategory.ID,
					CreateDocumentCategoryRequest: models.CreateDocumentCategoryRequest{
						Name: fmt.Sprintf("cate name %s updated", uniStr),
					},
				}
			},
			expect: func(uniStr string) models.DocumentCategory {
				return models.DocumentCategory{
					Model: createdCategory.Model,
					Name:  fmt.Sprintf("cate name %s updated", uniStr),
					Slug:  fmt.Sprintf("cate-name-%s-updated", strings.ToLower(uniStr)),
				}
			},
		},
		{
			name: "should return error when update document category name to existing one",
			req: func(uniStr string) models.UpdateDocumentCategoryRequest {
				return models.UpdateDocumentCategoryRequest{
					DocumentCategoryID: createdCategory.ID,
					CreateDocumentCategoryRequest: models.CreateDocumentCategoryRequest{
						Name: fmt.Sprintf("same cate name %s", uniStr),
					},
				}
			},
			setup: func(t *testing.T, uniStr string) {
				_, err := repo.CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
					Name: fmt.Sprintf("same cate name %s", uniStr),
				})
				assert.NoError(t, err)
			},
			expectErr: errs.New(http.StatusBadRequest, "Category name is already existed"),
		},
		{
			name: "should return error when document category does not exist",
			req: func(uniStr string) models.UpdateDocumentCategoryRequest {
				return models.UpdateDocumentCategoryRequest{
					DocumentCategoryID: "invalid-id",
					CreateDocumentCategoryRequest: models.CreateDocumentCategoryRequest{
						Name: fmt.Sprintf("cate name %s updated", uniStr),
					},
				}
			},
			expectErr: errs.ErrRecordNotFound,
		},
	}
	for _, testCase := range testCases {
		randomStr := random.String(8, random.Alphabet)
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setup != nil {
				testCase.setup(t, randomStr)
			}
			req := testCase.req(randomStr)
			resp, err := repo.UpdateDocumentCategory(&req)

			if testCase.expectErr != nil {
				assert.Equal(t, testCase.expectErr, err)
				return
			}
			assert.NoError(t, err)
			expect := testCase.expect(randomStr)
			dbDocument := models.DocumentCategory{}
			err = app.DB.First(&dbDocument, "id = ?", resp.ID).Error

			assert.NoError(t, err)
			assert.Equal(t, expect, dbDocument)
		})
	}
}

func TestDocumentRepo__AdminGetDocumentCategoryBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	repo := repo.NewDocumentCategoryRepo(app.DB)
	randomStr := random.String(8, random.Alphabet)
	categoryPayload := &models.CreateDocumentCategoryRequest{
		Name: fmt.Sprintf("exist cate %s", randomStr),
	}
	createdCategory, err := repo.CreateDocumentCategory(categoryPayload)
	assert.NoError(t, err)

	testCases := []struct {
		name      string
		req       models.GetDocumentCategoryParams
		expect    *models.DocumentCategory
		expectErr error
	}{
		{
			name: "should get document category correctly",
			req: models.GetDocumentCategoryParams{
				DocumentCategoryID: createdCategory.ID,
			},
			expect: &models.DocumentCategory{
				Model: createdCategory.Model,
				Name:  fmt.Sprintf("exist cate %s", randomStr),
				Slug:  fmt.Sprintf("exist-cate-%s", strings.ToLower(randomStr)),
			},
		},
		{
			name: "should return error when document category does not exist",
			req: models.GetDocumentCategoryParams{
				DocumentCategoryID: "invalid-id",
			},
			expectErr: errs.ErrRecordNotFound,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp, err := repo.GetDocumentCategory(&testCase.req)

			if testCase.expectErr != nil {
				assert.Equal(t, testCase.expectErr, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, testCase.expect, resp)
		})
	}
}

func TestDocumentRepo__AdminDeleteDocumentCategoryBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	cateRepo := repo.NewDocumentCategoryRepo(app.DB)

	testCases := []struct {
		name         string
		req          func(cateID string) models.GetDocumentCategoryParams
		expectErr    error
		withDocument bool
	}{
		{
			name: "should delete document category correctly",
			req: func(cateID string) models.GetDocumentCategoryParams {
				return models.GetDocumentCategoryParams{
					DocumentCategoryID: cateID,
				}
			},
		},
		{
			name: "should set document category_id to blank when delete reference category correctly",
			req: func(cateID string) models.GetDocumentCategoryParams {
				return models.GetDocumentCategoryParams{
					DocumentCategoryID: cateID,
				}
			},
			withDocument: true,
		},
		{
			name: "should return error when document category does not exist",
			req: func(cateID string) models.GetDocumentCategoryParams {
				return models.GetDocumentCategoryParams{
					DocumentCategoryID: "invalid id",
				}
			},
			expectErr: errs.ErrRecordNotFound,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			randomStr := random.String(8, random.Alphabet)
			createdCategory, err := cateRepo.CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
				Name: fmt.Sprintf("existing cate %s", randomStr),
			})
			assert.NoError(t, err)

			createdDocumentID := ""
			if testCase.withDocument {
				loginResp, err := adminLogin(app.DB)
				assert.NoError(t, err)
				randomStr := random.String(8, random.Alphabet)
				createdDocument, err := repo.NewDocumentRepo(app.DB).CreateDocument(&models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("title %s", randomStr),
					Content:    fmt.Sprintf("content-%s", randomStr),
					CategoryID: createdCategory.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", randomStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", randomStr),
							"size": 3000,
						},
					},
				})
				assert.NoError(t, err)
				createdDocumentID = createdDocument.ID
			}
			req := testCase.req(createdCategory.ID)
			err = cateRepo.DeleteDocumentCategory(&req)
			if testCase.expectErr != nil {
				assert.Equal(t, testCase.expectErr, err)
				return
			}
			assert.NoError(t, err)
			dbDocument := models.DocumentCategory{}
			err = app.DB.First(&dbDocument, "id = ?", createdCategory.ID).Error
			assert.Equal(t, gorm.ErrRecordNotFound, err)
			if testCase.withDocument {
				dbDocument := &models.Document{Model: models.Model{ID: createdDocumentID}}
				err := app.DB.Select("category_id").First(dbDocument).Error
				assert.NoError(t, err)
				assert.Equal(t, dbDocument.CategoryID, "")
			}
		})
	}
}

func TestDocumentRepo__AdminGetDocumentCategoryListBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	repo := repo.NewDocumentCategoryRepo(app.DB)

	_, err := repo.CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
		Name: fmt.Sprintf("exist cate %s", random.String(8, random.Alphabet)),
	})
	assert.NoError(t, err)

	t.Run("should get document category list correctly", func(t *testing.T) {
		list := repo.GetDocumentCategoryList(&models.GetDocumentCategoryListParams{})
		assert.Greater(t, len(list), 0)
	})
}
