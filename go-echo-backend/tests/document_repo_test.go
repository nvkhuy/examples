package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/thaitanloi365/go-utils/random"
	"gorm.io/gorm"
)

func TestDocumentRepo__AdminCreateDocument(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	router := routes.NewRouter(app)
	router.SetupRoutes()
	loginResp, err := adminLogin(app.DB)
	assert.NoError(t, err)
	createdCate, err := repo.NewDocumentCategoryRepo(app.DB).
		CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
			Name: fmt.Sprintf("category %s", random.String(8, random.Alphabet)),
		})
	assert.NoError(t, err)
	testCases := []struct {
		name             string
		req              models.CreateDocumentRequest
		expectedCode     int
		expectedErrorMsg string
	}{
		{
			name: "should create document correctly",
			req: models.CreateDocumentRequest{
				Title:      fmt.Sprintf("title %s", random.String(8, random.Alphabet)),
				Content:    fmt.Sprintf("content-%s", random.String(8, random.Alphabet)),
				CategoryID: createdCate.ID,
				Status:     enums.DocumentStatusNew,
				FeaturedImage: &models.Attachment{
					ContentType: "image/jpeg",
					FileKey:     fmt.Sprintf("key-%s", random.String(8, random.Alphabet)),
					Metadata: map[string]interface{}{
						"name": fmt.Sprintf("name-%s", random.String(8, random.Alphabet)),
						"size": 3000,
					},
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "should validate document title empty",
			req: models.CreateDocumentRequest{
				Content:    fmt.Sprintf("content-%s", random.String(8, random.Alphabet)),
				CategoryID: createdCate.ID,
				Status:     enums.DocumentStatusNew,
				FeaturedImage: &models.Attachment{
					ContentType: "image/jpeg",
					FileKey:     fmt.Sprintf("key-%s", random.String(8, random.Alphabet)),
					Metadata: map[string]interface{}{
						"name": fmt.Sprintf("name-%s", random.String(8, random.Alphabet)),
						"size": 3000,
					},
				},
			},
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "Title is required",
		},
		{
			name: "should validate document content empty",
			req: models.CreateDocumentRequest{
				Title:      fmt.Sprintf("title %s", random.String(8, random.Alphabet)),
				CategoryID: createdCate.ID,
				Status:     enums.DocumentStatusNew,
				FeaturedImage: &models.Attachment{
					ContentType: "image/jpeg",
					FileKey:     fmt.Sprintf("key-%s", random.String(8, random.Alphabet)),
					Metadata: map[string]interface{}{
						"name": fmt.Sprintf("name-%s", random.String(8, random.Alphabet)),
						"size": 3000,
					},
				},
			},
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "Content is required",
		},
		{
			name: "should validate document categoryId empty",
			req: models.CreateDocumentRequest{
				Content: fmt.Sprintf("content-%s", random.String(8, random.Alphabet)),
				Title:   fmt.Sprintf("title %s", random.String(8, random.Alphabet)),
				Status:  enums.DocumentStatusNew,
				FeaturedImage: &models.Attachment{
					ContentType: "image/jpeg",
					FileKey:     fmt.Sprintf("key-%s", random.String(8, random.Alphabet)),
					Metadata: map[string]interface{}{
						"name": fmt.Sprintf("name-%s", random.String(8, random.Alphabet)),
						"size": 3000,
					},
				},
			},
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "CategoryID is required",
		},
		{
			name: "should validate document feature image empty",
			req: models.CreateDocumentRequest{
				Content:    fmt.Sprintf("content-%s", random.String(8, random.Alphabet)),
				Title:      fmt.Sprintf("title %s", random.String(8, random.Alphabet)),
				CategoryID: createdCate.ID,
				Status:     enums.DocumentStatusNew,
			},
			expectedCode:     http.StatusBadRequest,
			expectedErrorMsg: "FeaturedImage is required",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			jsonReq, err := json.Marshal(testCase.req)
			assert.NoError(t, err)
			var req = httptest.NewRequest(http.MethodPost, "/api/v1/admin/documents", strings.NewReader(string(jsonReq)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginResp.Token))

			var rec = httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedCode, rec.Code)
			if testCase.expectedErrorMsg != "" {
				bodyStruct := struct {
					Message string `json:"message"`
				}{}
				err := json.Unmarshal(rec.Body.Bytes(), &bodyStruct)
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedErrorMsg, bodyStruct.Message)
			}
		})
	}
}
func TestDocumentRepo__AdminCreateDocumentBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	documentRepo := repo.NewDocumentRepo(app.DB)
	loginResp, err := adminLogin(app.DB)
	assert.NoError(t, err)
	createdCate, err := repo.NewDocumentCategoryRepo(app.DB).
		CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
			Name: fmt.Sprintf("category %s", random.String(8, random.Alphabet)),
		})
	assert.NoError(t, err)
	testCases := []struct {
		name   string
		req    func(uniStr string) models.CreateDocumentRequest
		setup  func(t *testing.T, uniStr string)
		expect func(uniStr string) models.Document
	}{
		{
			name: "should create document correctly",
			req: func(uniStr string) models.CreateDocumentRequest {
				return models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					Status:     enums.DocumentStatusInactive,
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": 3000,
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
						Status:  enums.DocumentStatusInactive,
					},
				}
			},
			expect: func(uniStr string) models.Document {
				return models.Document{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					Status:     enums.DocumentStatusInactive,
					CategoryID: createdCate.ID,
					Slug:       fmt.Sprintf("title-%s", strings.ToLower(uniStr)),
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": float64(3000),
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
						Status:  enums.DocumentStatusInactive,
						Slug:    fmt.Sprintf("title-%s", strings.ToLower(uniStr)),
					},
				}
			},
		},
		{
			name: "should set default document status correctly",
			req: func(uniStr string) models.CreateDocumentRequest {
				return models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": 3000,
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
					},
				}
			},
			expect: func(uniStr string) models.Document {
				return models.Document{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					Status:     enums.DocumentStatusNew,
					CategoryID: createdCate.ID,
					Slug:       fmt.Sprintf("title-%s", strings.ToLower(uniStr)),
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": float64(3000),
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
						Status:  enums.DocumentStatusNew,
						Slug:    fmt.Sprintf("title-%s", strings.ToLower(uniStr)),
					},
				}
			},
		},
		{
			name: "should generate document slug postfix correctly",
			setup: func(t *testing.T, uniStr string) {
				_, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("same title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": 3000,
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("same title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
					},
				})
				assert.NoError(t, err)
			},
			req: func(uniStr string) models.CreateDocumentRequest {
				return models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("same title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": 3000,
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("same title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
					},
				}
			},
			expect: func(uniStr string) models.Document {
				return models.Document{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("same title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					Status:     enums.DocumentStatusNew,
					CategoryID: createdCate.ID,
					Slug:       fmt.Sprintf("same-title-%s-%d", strings.ToLower(uniStr), 1),
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": float64(3000),
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("same title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
						Status:  "new",
						Slug:    fmt.Sprintf("same-title-%s-%d", strings.ToLower(uniStr), 1),
					},
				}
			},
		},
	}
	for _, testCase := range testCases {
		randomStr := random.String(8, random.Alphabet)
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setup != nil {
				testCase.setup(t, randomStr)
			}
			req := testCase.req(randomStr)
			expect := testCase.expect(randomStr)
			resp, err := documentRepo.CreateDocument(&req)
			assert.NoError(t, err)
			expect.Model = resp.Model
			dbDocument := models.Document{}
			err = app.DB.First(&dbDocument, "id = ?", resp.ID).Error

			assert.NoError(t, err)
			assert.Equal(t, expect, dbDocument)
		})
	}
}

func TestDocumentRepo__AdminUpdateDocumentBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	documentRepo := repo.NewDocumentRepo(app.DB)
	loginResp, err := adminLogin(app.DB)
	assert.NoError(t, err)
	createdCate, err := repo.NewDocumentCategoryRepo(app.DB).
		CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
			Name: fmt.Sprintf("category %s", random.String(8, random.Alphabet)),
		})
	assert.NoError(t, err)
	testCases := []struct {
		name      string
		req       func(id, uniStr string) models.UpdateDocumentRequest
		setup     func(t *testing.T, uniStr string)
		expect    func(id, uniStr string) models.Document
		expectErr error
	}{
		{
			name: "should update document correctly",
			req: func(id, uniStr string) models.UpdateDocumentRequest {
				return models.UpdateDocumentRequest{
					DocumentID: id,
					CreateDocumentRequest: models.CreateDocumentRequest{
						Title:      fmt.Sprintf("title %s updated", uniStr),
						UserID:     loginResp.User.ID,
						Content:    fmt.Sprintf("content-%s updated", uniStr),
						Status:     enums.DocumentStatusInactive,
						CategoryID: createdCate.ID,
						FeaturedImage: &models.Attachment{
							ContentType: "image/png",
							FileKey:     fmt.Sprintf("key-%s updated", uniStr),
							Metadata: map[string]interface{}{
								"name": fmt.Sprintf("name-%s updated", uniStr),
								"size": 3000,
							},
						},
						Vi: &models.DocumentContent{
							Title:   fmt.Sprintf("title %s updated", uniStr),
							Content: fmt.Sprintf("content-%s updated", uniStr),
							Status:  enums.DocumentStatusInactive,
						},
					},
				}
			},
			expect: func(id, uniStr string) models.Document {
				return models.Document{
					Model: models.Model{
						ID: id,
					},
					Title:      fmt.Sprintf("title %s updated", uniStr),
					Slug:       fmt.Sprintf("title-%s-updated", strings.ToLower(uniStr)),
					UserID:     loginResp.User.ID,
					Content:    fmt.Sprintf("content-%s updated", uniStr),
					Status:     enums.DocumentStatusInactive,
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/png",
						FileKey:     fmt.Sprintf("key-%s updated", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s updated", uniStr),
							"size": float64(3000),
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("title %s updated", uniStr),
						Slug:    fmt.Sprintf("title-%s-updated", strings.ToLower(uniStr)),
						Content: fmt.Sprintf("content-%s updated", uniStr),
						Status:  enums.DocumentStatusInactive,
					},
				}
			},
		},
		{
			name: "should set default document status correctly",
			req: func(id, uniStr string) models.UpdateDocumentRequest {
				return models.UpdateDocumentRequest{
					DocumentID: id,
					CreateDocumentRequest: models.CreateDocumentRequest{
						Title:      fmt.Sprintf("title %s updated", uniStr),
						UserID:     loginResp.User.ID,
						Content:    fmt.Sprintf("content-%s updated", uniStr),
						CategoryID: createdCate.ID,
						FeaturedImage: &models.Attachment{
							ContentType: "image/png",
							FileKey:     fmt.Sprintf("key-%s updated", uniStr),
							Metadata: map[string]interface{}{
								"name": fmt.Sprintf("name-%s updated", uniStr),
								"size": 3000,
							},
						},
						Vi: &models.DocumentContent{
							Title:   fmt.Sprintf("title %s updated", uniStr),
							Content: fmt.Sprintf("content-%s updated", uniStr),
						},
					},
				}
			},
			expect: func(id, uniStr string) models.Document {
				return models.Document{
					Model: models.Model{
						ID: id,
					},
					Title:      fmt.Sprintf("title %s updated", uniStr),
					Slug:       fmt.Sprintf("title-%s-updated", strings.ToLower(uniStr)),
					UserID:     loginResp.User.ID,
					Content:    fmt.Sprintf("content-%s updated", uniStr),
					Status:     enums.DocumentStatusNew,
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/png",
						FileKey:     fmt.Sprintf("key-%s updated", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s updated", uniStr),
							"size": float64(3000),
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("title %s updated", uniStr),
						Slug:    fmt.Sprintf("title-%s-updated", strings.ToLower(uniStr)),
						Content: fmt.Sprintf("content-%s updated", uniStr),
						Status:  enums.DocumentStatusNew,
					},
				}
			},
		},
		{
			name: "should generate document slug postfix correctly",
			setup: func(t *testing.T, uniStr string) {
				_, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("same title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": 3000,
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("same title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
					},
				})
				assert.NoError(t, err)
			},
			req: func(id, uniStr string) models.UpdateDocumentRequest {
				return models.UpdateDocumentRequest{
					DocumentID: id,
					CreateDocumentRequest: models.CreateDocumentRequest{
						Title:      fmt.Sprintf("same title %s", uniStr),
						UserID:     loginResp.User.ID,
						Content:    fmt.Sprintf("content-%s updated", uniStr),
						Status:     enums.DocumentStatusInactive,
						CategoryID: createdCate.ID,
						FeaturedImage: &models.Attachment{
							ContentType: "image/png",
							FileKey:     fmt.Sprintf("key-%s updated", uniStr),
							Metadata: map[string]interface{}{
								"name": fmt.Sprintf("name-%s updated", uniStr),
								"size": 3000,
							},
						},
						Vi: &models.DocumentContent{
							Title:   fmt.Sprintf("same title %s", uniStr),
							Content: fmt.Sprintf("content-%s updated", uniStr),
							Status:  enums.DocumentStatusInactive,
						},
					},
				}
			},
			expect: func(id, uniStr string) models.Document {
				return models.Document{
					Model: models.Model{
						ID: id,
					},
					Title:      fmt.Sprintf("same title %s", uniStr),
					Slug:       fmt.Sprintf("same-title-%s-%d", strings.ToLower(uniStr), 1),
					UserID:     loginResp.User.ID,
					Content:    fmt.Sprintf("content-%s updated", uniStr),
					Status:     enums.DocumentStatusInactive,
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/png",
						FileKey:     fmt.Sprintf("key-%s updated", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s updated", uniStr),
							"size": float64(3000),
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("same title %s", uniStr),
						Slug:    fmt.Sprintf("same-title-%s-%d", strings.ToLower(uniStr), 1),
						Content: fmt.Sprintf("content-%s updated", uniStr),
						Status:  enums.DocumentStatusInactive,
					},
				}
			},
		},
		{
			name: "should return not found error correctly when document does not exist",
			setup: func(t *testing.T, uniStr string) {
				_, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("same title %s", uniStr),
					Content:    fmt.Sprintf("content-%s", uniStr),
					CategoryID: createdCate.ID,
					FeaturedImage: &models.Attachment{
						ContentType: "image/jpeg",
						FileKey:     fmt.Sprintf("key-%s", uniStr),
						Metadata: map[string]interface{}{
							"name": fmt.Sprintf("name-%s", uniStr),
							"size": 3000,
						},
					},
					Vi: &models.DocumentContent{
						Title:   fmt.Sprintf("same title %s", uniStr),
						Content: fmt.Sprintf("content-%s", uniStr),
					},
				})
				assert.NoError(t, err)
			},
			req: func(id, uniStr string) models.UpdateDocumentRequest {
				return models.UpdateDocumentRequest{
					DocumentID: "invalid id",
					CreateDocumentRequest: models.CreateDocumentRequest{
						Title:      fmt.Sprintf("title %s update", uniStr),
						UserID:     loginResp.User.ID,
						Content:    fmt.Sprintf("content-%s updated", uniStr),
						Status:     enums.DocumentStatusInactive,
						CategoryID: createdCate.ID,
						FeaturedImage: &models.Attachment{
							ContentType: "image/png",
							FileKey:     fmt.Sprintf("key-%s updated", uniStr),
							Metadata: map[string]interface{}{
								"name": fmt.Sprintf("name-%s updated", uniStr),
								"size": 3000,
							},
						},
						Vi: &models.DocumentContent{
							Title:   fmt.Sprintf("title %s updates", uniStr),
							Content: fmt.Sprintf("content-%s updated", uniStr),
							Status:  enums.DocumentStatusInactive,
						},
					},
				}
			},
			expect: func(id, uniStr string) models.Document {
				return models.Document{}
			},
			expectErr: errs.ErrRecordNotFound,
		},
	}
	for _, testCase := range testCases {
		randomStr := random.String(8, random.Alphabet)
		t.Run(testCase.name, func(t *testing.T) {
			resp, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
				UserID:     loginResp.User.ID,
				Title:      fmt.Sprintf("title %s", randomStr),
				Content:    fmt.Sprintf("content-%s", randomStr),
				CategoryID: createdCate.ID,
				Status:     "new",
				FeaturedImage: &models.Attachment{
					ContentType: "image/jpeg",
					FileKey:     fmt.Sprintf("key-%s", randomStr),
					Metadata: map[string]interface{}{
						"name": fmt.Sprintf("name-%s", randomStr),
						"size": 3000,
					},
				},
				Vi: &models.DocumentContent{
					Title:   fmt.Sprintf("title %s", randomStr),
					Content: fmt.Sprintf("content-%s", randomStr),
					Status:  "new",
				},
			})
			assert.NoError(t, err)

			if testCase.setup != nil {
				testCase.setup(t, randomStr)
			}
			req := testCase.req(resp.ID, randomStr)
			expect := testCase.expect(resp.ID, randomStr)
			resp, err = documentRepo.UpdateDocument(&req)

			if testCase.expectErr != nil {
				assert.Equal(t, testCase.expectErr, err)
				return
			}
			assert.NoError(t, err)
			expect.CreatedAt = resp.CreatedAt
			expect.UpdatedAt = resp.UpdatedAt
			dbDocument := models.Document{}
			err = app.DB.First(&dbDocument, "id = ?", resp.ID).Error
			assert.NoError(t, err)
			assert.Equal(t, expect, dbDocument)
		})
	}
}

func TestDocumentRepo__AdminGetDocumentDetailBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	documentRepo := repo.NewDocumentRepo(app.DB)
	loginResp, err := adminLogin(app.DB)
	assert.NoError(t, err)
	createdCate, err := repo.NewDocumentCategoryRepo(app.DB).
		CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
			Name: fmt.Sprintf("category %s", random.String(8, random.Alphabet)),
		})
	assert.NoError(t, err)
	randomStr := random.String(8, random.Alphabet)
	createdDocument, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
		UserID:     loginResp.User.ID,
		Title:      fmt.Sprintf("title %s", randomStr),
		Content:    fmt.Sprintf("content-%s", randomStr),
		CategoryID: createdCate.ID,
		Status:     enums.DocumentStatusNew,
		FeaturedImage: &models.Attachment{
			ContentType: "image/jpeg",
			FileKey:     fmt.Sprintf("key-%s", randomStr),
			Metadata: map[string]interface{}{
				"name": fmt.Sprintf("name-%s", randomStr),
				"size": 3000,
			},
		},
		Vi: &models.DocumentContent{
			Title:   fmt.Sprintf("title %s", randomStr),
			Content: fmt.Sprintf("content-%s", randomStr),
			Status:  enums.DocumentStatusNew,
		},
	})
	assert.NoError(t, err)
	userInfo := &models.User{}
	if err := app.DB.Select("ID", "Name", "Avatar").First(userInfo).Error; err != nil {
		t.Error(err)
	}
	userInfo.ID = ""

	expectDocument := &models.Document{
		Model:      createdDocument.Model,
		Title:      fmt.Sprintf("title %s", randomStr),
		Content:    fmt.Sprintf("content-%s", randomStr),
		CategoryID: createdCate.ID,
		Status:     enums.DocumentStatusNew,
		Slug:       fmt.Sprintf("title-%s", strings.ToLower(randomStr)),
		FeaturedImage: &models.Attachment{
			ContentType: "image/jpeg",
			FileKey:     fmt.Sprintf("key-%s", randomStr),
			Metadata: map[string]interface{}{
				"name": fmt.Sprintf("name-%s", randomStr),
				"size": float64(3000),
			},
		},
		Vi: &models.DocumentContent{
			Title:   fmt.Sprintf("title %s", randomStr),
			Content: fmt.Sprintf("content-%s", randomStr),
			Status:  enums.DocumentStatusNew,
			Slug:    fmt.Sprintf("title-%s", strings.ToLower(randomStr)),
		},
		User:     userInfo,
		Category: createdCate,
	}
	testCases := []struct {
		name      string
		params    models.GetDocumentParams
		expectErr error
	}{
		{
			name: "should get document detail by slug correctly",
			params: models.GetDocumentParams{
				Slug: fmt.Sprintf("title-%s", strings.ToLower(randomStr)),
			},
		},
		{
			name: "should return error correctly when document does not exist",
			params: models.GetDocumentParams{
				Slug: "invalid-slug",
			},
			expectErr: errs.ErrRecordNotFound,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dbDocument, err := documentRepo.GetDocument(&testCase.params)
			if testCase.expectErr != nil {
				assert.Equal(t, testCase.expectErr, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, expectDocument, dbDocument)
		})
	}
}

func TestDocumentRepo__AdminDeleteDocumentBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	documentRepo := repo.NewDocumentRepo(app.DB)
	loginResp, err := adminLogin(app.DB)
	assert.NoError(t, err)
	createdCate, err := repo.NewDocumentCategoryRepo(app.DB).
		CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
			Name: fmt.Sprintf("category %s", random.String(8, random.Alphabet)),
		})
	assert.NoError(t, err)
	randomStr := random.String(8, random.Alphabet)
	createdDocument, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
		UserID:     loginResp.User.ID,
		Title:      fmt.Sprintf("title %s", randomStr),
		Content:    fmt.Sprintf("content-%s", randomStr),
		CategoryID: createdCate.ID,
		Status:     enums.DocumentStatusNew,
		FeaturedImage: &models.Attachment{
			ContentType: "image/jpeg",
			FileKey:     fmt.Sprintf("key-%s", randomStr),
			Metadata: map[string]interface{}{
				"name": fmt.Sprintf("name-%s", randomStr),
				"size": 3000,
			},
		},
		Vi: &models.DocumentContent{
			Title:   fmt.Sprintf("title %s", randomStr),
			Content: fmt.Sprintf("content-%s", randomStr),
			Status:  enums.DocumentStatusNew,
		},
	})
	assert.NoError(t, err)
	testCases := []struct {
		name      string
		params    models.DeleteDocumentParams
		expectErr error
	}{
		{
			name: "should delete documentcorrectly",
			params: models.DeleteDocumentParams{
				DocumentID: createdDocument.ID,
			},
		},
		{
			name: "should return error when document does not exist",
			params: models.DeleteDocumentParams{
				DocumentID: "invalid-id",
			},
			expectErr: errs.ErrRecordNotFound,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := documentRepo.DeleteDocument(&testCase.params)
			if testCase.expectErr != nil {
				assert.Equal(t, testCase.expectErr, err)
				return
			}
			assert.NoError(t, err)
			dbDocument := models.Document{}
			err = app.DB.First(&dbDocument, "id = ?", createdDocument.ID).Error
			assert.Equal(t, gorm.ErrRecordNotFound, err)
		})
	}
}

func TestDocumentRepo__AdminGetDocumentListBiz(t *testing.T) {
	t.Parallel()
	app := initApp("local")
	loginResp, err := adminLogin(app.DB)
	assert.NoError(t, err)
	documentRepo := repo.NewDocumentRepo(app.DB)
	createdCate, err := repo.NewDocumentCategoryRepo(app.DB).
		CreateDocumentCategory(&models.CreateDocumentCategoryRequest{
			Name: fmt.Sprintf("category %s", random.String(8, random.Alphabet)),
		})
	assert.NoError(t, err)
	testCases := []struct {
		name   string
		params func(uniStr string) models.GetDocumentListParams
		setup  func(uniStr string)
	}{
		{
			name: "should get document list correctly",
			params: func(randomStr string) models.GetDocumentListParams {
				return models.GetDocumentListParams{
					PaginationParams: models.PaginationParams{
						Page:  1,
						Limit: 10,
					},
				}
			},
			setup: func(randomStr string) {
				_, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("title %s", randomStr),
					Content:    fmt.Sprintf("content-%s", randomStr),
					CategoryID: createdCate.ID,
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
			},
		},
		{
			name: "should filter document list by status correctly",
			params: func(randomStr string) models.GetDocumentListParams {
				return models.GetDocumentListParams{
					PaginationParams: models.PaginationParams{
						Page:  1,
						Limit: 10,
					},
					Statuses: []enums.DocumentStatus{enums.DocumentStatusInactive},
				}
			},
			setup: func(randomStr string) {
				_, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("title %s", randomStr),
					Content:    fmt.Sprintf("content-%s", randomStr),
					Status:     enums.DocumentStatusInactive,
					CategoryID: createdCate.ID,
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
			},
		},
		{
			name: "should filter document list by name correctly",
			params: func(randomStr string) models.GetDocumentListParams {
				return models.GetDocumentListParams{
					PaginationParams: models.PaginationParams{
						Page:    1,
						Limit:   10,
						Keyword: randomStr,
					},
				}
			},
			setup: func(randomStr string) {
				_, err := documentRepo.CreateDocument(&models.CreateDocumentRequest{
					UserID:     loginResp.User.ID,
					Title:      fmt.Sprintf("title %s", randomStr),
					Content:    fmt.Sprintf("content-%s", randomStr),
					Status:     enums.DocumentStatusInactive,
					CategoryID: createdCate.ID,
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
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			randomStr := random.String(8, random.Alphabet)
			if testCase.setup != nil {
				testCase.setup(randomStr)
			}
			params := testCase.params(randomStr)
			resp := documentRepo.GetDocumentList(&params)
			documentList, ok := resp.Records.(*[]*models.Document)
			assert.True(t, ok)
			assert.NotEmpty(t, *documentList)
		})
	}
}
