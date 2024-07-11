package repo

import (
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"gorm.io/gorm/clause"
)

type UserNotificationRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewUserNotificationRepo(db *db.DB) *UserNotificationRepo {
	return &UserNotificationRepo{
		db:     db,
		logger: logger.New("repo/UserNotification"),
	}
}

type PaginateUserNotificationsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	UserID string `json:"user_id" query:"user_id" form:"user_id"`
}

func (r *UserNotificationRepo) PaginateUserNotifications(params PaginateUserNotificationsParams) *query.Pagination {
	var builder = queryfunc.NewUserNotificationBuilder(queryfunc.UserNotificationBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("un.deleted_at IS NULL AND un.seen_at IS NULL")

			if params.GetRole().IsAdmin() {
				if params.UserID != "" {
					builder.Where("un.user_id = ?", params.UserID)
				}
			} else {
				builder.Where("un.user_id = ?", params.GetUserID())
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *UserNotificationRepo) CreateUserNotification(params models.UserNotificationForm) (*models.UserNotification, error) {
	var form models.UserNotification
	err := copier.Copy(&form, &params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Create(&form).Error
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	return &form, nil
}

func (r *UserNotificationRepo) GetUserNotificationByID(id string) (*models.UserNotification, error) {
	var notification models.UserNotification
	var builder = queryfunc.NewUserNotificationBuilder(queryfunc.UserNotificationBuilderOptions{})
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cn.deleted_at IS NULL AND cn.id = ?", id)

		}).
		FirstFunc(&notification)

	return &notification, err
}

type MarkSeenNotificationParams struct {
	models.JwtClaimsInfo

	NotificationID string `json:"notification_id" query:"notification_id" param:"notification_id" validate:"required"`
}

func (r *UserNotificationRepo) MarkSeen(params MarkSeenNotificationParams) error {

	var err = r.db.Model(&models.UserNotification{}).
		Where("id = ?", params.NotificationID).
		Where("user_id = ?", params.GetUserID()).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

func (r *UserNotificationRepo) MarkSeenAll(claims models.JwtClaimsInfo) error {
	var err = r.db.Model(&models.UserNotification{}).
		Where("seen_at IS NULL").
		Where("user_id = ?", claims.GetUserID()).
		Update("seen_at", time.Now().Unix()).Error

	return err
}
