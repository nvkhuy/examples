package repo

import (
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"

	"github.com/rotisserie/eris"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewCommentRepo(db *db.DB) *CommentRepo {
	return &CommentRepo{
		db:     db,
		logger: logger.New("repo/Comment"),
	}
}

type PaginateCommentsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	OrderByQuery string

	TargetType enums.CommentTargetType
	TargetID   string

	FileKey string `json:"file_key" query:"file_key" params:"file_key"`
}

func (r *CommentRepo) PaginateComment(params PaginateCommentsParams) *query.Pagination {
	var builder = queryfunc.NewCommentBuilder(queryfunc.CommentBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var orderBy = "c.created_at ASC"
	if params.OrderByQuery != "" {
		orderBy = params.OrderByQuery
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.TargetType != "" && params.TargetID != "" {
				builder.Where("c.target_type = ? and c.target_id = ?", params.TargetType, params.TargetID)
			}

			if params.FileKey != "" {
				builder.Where("c.file_key = ?", params.FileKey)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		OrderBy(orderBy).
		PagingFunc()

	return result
}

func (r *CommentRepo) CreateComment(form models.CommentCreateForm) (*models.Comment, error) {
	var comment models.Comment
	err := copier.Copy(&comment, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	comment.UserID = form.GetUserID()

	err = r.db.Omit(clause.Associations).Create(&comment).Error
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return &comment, nil
}

type InquirySellerCreateCommentParams struct {
	models.JwtClaimsInfo

	InquirySellerID string                      `json:"inquiry_seller_id" query:"inquiry_seller_id" param:"inquiry_seller_id" validate:"required"`
	Comments        []*models.CommentCreateForm `json:"comments" query:"comments" params:"comments" validate:"required"`
}

func (r *CommentRepo) CreateInquirySellerComment(form InquirySellerCreateCommentParams) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := copier.Copy(&comments, &form.Comments)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range comments {
			item.UserID = form.GetUserID()
			item.TargetID = form.InquirySellerID
		}

		err = tx.Create(&comments).Error

		return err
	})
	if err != nil {
		return nil, err
	}

	return comments, nil
}

type PaginateInquirySellerCommentsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	InquirySellerID string `json:"inquiry_seller_id" query:"inquiry_seller_id" param:"inquiry_seller_id" validate:"required"`
	FileKey         string `json:"file_key" query:"file_key" param:"file_key" validate:"required"`
}

func (r *CommentRepo) PaginateInquirySelerComments(params PaginateInquirySellerCommentsParams) *query.Pagination {
	inquirySeller, err := NewInquirySellerRepo(r.db).
		GetInquirySellerRequestByID(params.InquirySellerID, queryfunc.InquirySellerRequestBuilderOptions{})
	if err != nil {
		return nil
	}

	var builder = queryfunc.NewCommentBuilder(queryfunc.CommentBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("c.target_id = ?", inquirySeller.InquiryID)
			builder.Where("c.target_type = ?", enums.CommentTargetTypeInquirySellerDesign)
			if params.FileKey != "" {
				builder.Where("c.file_key = ?", params.FileKey)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type DeleteCommentParams struct {
	models.JwtClaimsInfo

	CommentID string `param:"comment_id" query:"comment_id" path:"comment_id" validate:"required"`
}

func (r *CommentRepo) DeleteComment(params DeleteCommentParams) error {

	var err = r.db.Unscoped().Delete(&models.Comment{}, "id = ?", params.CommentID).Error

	return err
}

type MarkSeenParams struct {
	models.JwtClaimsInfo

	TargetID   string                  `json:"target_id" param:"target_id"`
	TargetType enums.CommentTargetType `json:"target_type" param:"target_type"`
}

func (r *CommentRepo) MarkSeen(params MarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id <> ?", params.GetUserID()).
		Where("target_type = ?", params.TargetType).
		Where("target_id = ?", params.TargetID).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

type GetUnreadCountParams struct {
	models.JwtClaimsInfo
	TargetID   string                  `json:"target_id" param:"target_id"`
	TargetType enums.CommentTargetType `json:"target_type" param:"target_type"`
}

type GetUnreadCountResponse struct {
	TotalCount int64 `json:"total_count"`
}

func (r *CommentRepo) GetUnreadCount(params GetUnreadCountParams) GetUnreadCountResponse {
	var resp GetUnreadCountResponse
	r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id <> ?", params.GetUserID()).
		Where("target_type = ?", params.TargetType).
		Where("target_id = ?", params.TargetID).
		Count(&resp.TotalCount)

	return resp
}
