package repo

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/rs/xid"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	simpleslug "github.com/gosimple/slug"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"

	"gorm.io/gorm/clause"
)

type PostRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewPostRepo(db *db.DB) *PostRepo {
	return &PostRepo{
		db:     db,
		logger: logger.New("repo/Post"),
	}
}

type PaginatePostParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	Name        string             `json:"name" query:"name" form:"name"`
	CategoryIDs []string           `json:"category_ids" query:"category_ids" form:"category_ids"`
	Statuses    []enums.PostStatus `json:"statuses" query:"statuses" form:"statuses"`
	UserIDs     []string           `json:"user_ids" query:"user_ids" form:"user_ids"`

	Language enums.LanguageCode `json:"-"`
}

func (r *PostRepo) PaginatePost(params PaginatePostParams) *query.Pagination {
	var builder = queryfunc.NewPostBuilder(queryfunc.PostBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.CategoryIDs) > 0 {
				builder.Where("p.category_id IN ?", params.CategoryIDs)
			}

			if len(params.Statuses) > 0 {
				var statues = []enums.PostStatus{enums.PostStatusPublished}
				switch params.Language {
				case enums.LanguageCodeVietnam:
					builder.Where("vi ->> 'status' IN ?", statues)
				case enums.LanguageCodeEnglish:
					builder.Where("p.status IN ?", statues)
				default:
					builder.Where("(p.status IN ? OR vi ->> 'status' IN ?)", statues, statues)
				}
			}
			if len(params.UserIDs) > 0 {
				builder.Where("p.user_id IN ?", params.UserIDs)
			}

			if len(params.Keyword) > 0 {
				q := "%" + params.Keyword + "%"
				builder.Where("(p.title ILIKE @keyword OR unaccent(vi ->> 'title') ILIKE unaccent(@keyword))", sql.Named("keyword", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *PostRepo) CreatePost(form models.PostCreateForm) (*models.Post, error) {
	if strings.TrimSpace(form.Title) == "" {
		err := errors.New("post title is empty")
		return nil, err
	}

	var post models.Post
	err := copier.Copy(&post, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}
	post.Slug = simpleslug.Make(form.Title)
	if post.VI != nil {
		post.VI.Slug = simpleslug.Make(post.VI.Title)
	}
	err = r.db.Omit(clause.Associations).Create(&post).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, eris.Wrap(err, "")
	}

	return &post, nil
}

type GetPostParams struct {
	models.JwtClaimsInfo
	PostID string `json:"post_id,omitempty"`
	Slug   string `param:"slug" query:"slug" form:"slug" validate:"required"`
}

func (r *PostRepo) GetPost(params GetPostParams) (*models.Post, error) {
	var builder = queryfunc.NewPostBuilder(queryfunc.PostBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var Post models.Post
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.Slug != "" {
				builder.Where("p.slug = ? or p.vi->>'slug' = ?", params.Slug, params.Slug)
			}
			if params.PostID != "" {
				builder.Where("p.id = ?", params.PostID)
			}
		}).
		FirstFunc(&Post)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &Post, nil
}

func (r *PostRepo) UpdatePost(form models.PostUpdateForm) (post *models.Post, err error) {
	var update models.Post

	err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	update.ID = form.PostID

	var find models.Post
	err = r.db.First(&find, "id = ?", form.PostID).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Omit(clause.Associations).Clauses(clause.Returning{}).Model(&update).Where("id = ?", form.PostID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	if update.SettingSEO != nil {
		update.SettingSEO.Route = "/blogs/" + update.Slug
		update.SettingSEO.ID = find.SettingSeoID
		if update.SettingSEO.ID != "" {
			if err = r.db.Model(&update.SettingSEO).Clauses(clause.Returning{}).Where("id = ?", update.SettingSeoID).Updates(update.SettingSEO).Error; err != nil {
				return
			}
		} else {
			if err = r.db.Model(&update.SettingSEO).Clauses(clause.Returning{}).Create(update.SettingSEO).Error; err != nil {
				return
			}
		}
		update.SettingSeoID = update.SettingSEO.ID
	}

	if update.VI != nil && update.VI.SettingSEO != nil {
		update.VI.SettingSEO.Route = "/blogs/" + update.VI.Slug
		update.VI.SettingSEO.ID = find.VI.SettingSeoID
		if update.VI.SettingSEO.ID != "" {
			if err = r.db.Model(&update.VI.SettingSEO).Clauses(clause.Returning{}).Where("id = ?", update.VI.SettingSeoID).Updates(update.VI.SettingSEO).Error; err != nil {
				return
			}
		} else {
			if err = r.db.Model(&update.VI.SettingSEO).Clauses(clause.Returning{}).Create(update.VI.SettingSEO).Error; err != nil {
				return
			}
		}
		update.VI.SettingSeoID = update.VI.SettingSEO.ID
	}

	err = r.db.Omit(clause.Associations).Clauses(clause.Returning{}).Model(&update).Where("id = ?", form.PostID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	post = &update
	return
}

type DeletePostParams struct {
	models.JwtClaimsInfo

	PostID string `param:"post_id" query:"post_id" path:"post_id" validate:"required"`
}

func (r *PostRepo) DeletePost(params DeletePostParams) error {

	var err = r.db.Unscoped().Delete(&models.Post{}, "id = ?", params.PostID).Error

	return err
}

func (r *PostRepo) GenerateSlug() (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) (e error) {
		var posts models.PostSlice
		r.db.Find(&posts)
		for _, pd := range posts {
			e = tx.Model(pd).Where("id = ?", pd.ID).Updates(&pd).Error
			if e != nil {
				return
			}
		}
		return
	})
	return
}

func (r *PostRepo) UpdateContentURL(id string, languageCode enums.LanguageCode, contentURL string) (err error) {
	if languageCode == enums.LanguageCodeEnglish {
		err = r.db.Model(&models.Post{}).Where("id = ?", id).Update("content_url", contentURL).Error
	} else if languageCode == enums.LanguageCodeVietnam {
		updateJSON := map[string]interface{}{"content_url": contentURL}
		err = r.db.Model(&models.Post{}).Where("id = ?", id).Update("vi", gorm.Expr("vi || ?", updateJSON)).Error
	}
	return
}

func (r *PostRepo) MoveImageToS3(id string) (err error) {
	var post models.Post
	err = r.db.Select("id", "content").Where("id = ?", id).Find(&post).Error
	if err != nil {
		return
	}
	content := post.Content
	imgURLs := r.extractAll("https://lh7-us.googleusercontent", "width", content)
	var s3Client = s3.New(r.db.Configuration)
	for _, url := range imgURLs {
		var data []byte
		data, err = helper.DownloadImageFromURL(url)
		if err != nil {
			continue
		}
		imageBytes := bytes.NewBuffer(data)
		var uploadParams = s3.UploadFileParams{
			Bucket:      r.db.Configuration.AWSS3StorageBucket,
			Data:        imageBytes,
			ContentType: "image/jpeg",
			ACL:         "private",
			Key:         fmt.Sprintf("/uploads/media/post_%s_image_%s.jpg", id, xid.New().String()),
		}
		_, err = s3Client.UploadFile(uploadParams)
		if err != nil {
			continue
		}
		newImageURL := fmt.Sprintf("https://%s%s", r.db.Configuration.StorageURL, uploadParams.Key)
		content = strings.ReplaceAll(content, url, newImageURL)
	}
	r.db.Model(&models.Post{}).Where("id = ?", id).Update("content", content)
	return
}

func (r *PostRepo) extractAll(startWord, endWord string, content string) (results []string) {
	var start = 0
	for {
		start = strings.Index(content, startWord)
		if start == -1 {
			break
		}
		if start >= 0 {
			sub := content[start:]
			end := strings.Index(sub, endWord)
			if end >= 0 {
				imgURL := strings.ReplaceAll(strings.ReplaceAll(sub[:end], " ", ""), strconv.Itoa(int('"')), "")
				imgURL = imgURL[:len(imgURL)-1]
				results = append(results, imgURL)
				content = content[start+end:]
			} else {
				break
			}
		}
	}
	return
}

type GetPostStatsParams struct {
	models.JwtClaimsInfo
	Status enums.PostStatus `param:"status" query:"status" path:"status" json:"status,omitempty"`
}

func (r *PostRepo) GetPostStats(params GetPostStatsParams) (*models.PostStats, error) {
	var result models.PostStats
	if params.Status == "" {
		params.Status = enums.PostStatusPublished
	}
	r.db.Model(&models.Post{}).Where("status = ?", params.Status).Count(&result.Total)
	return &result, nil
}
