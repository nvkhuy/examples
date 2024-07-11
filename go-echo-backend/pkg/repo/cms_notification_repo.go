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

type CmsNotificationRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewCmsNotificationRepo(db *db.DB) *CmsNotificationRepo {
	return &CmsNotificationRepo{
		db:     db,
		logger: logger.New("repo/CmsNotification"),
	}
}

type PaginateCmsNotificationsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *CmsNotificationRepo) PaginateCmsNotifications(params PaginateCmsNotificationsParams) *query.Pagination {
	var builder = queryfunc.NewCmsNotificationBuilder(queryfunc.CmsNotificationBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cn.deleted_at IS NULL AND cn.seen_at IS NULL")
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *CmsNotificationRepo) CreateCmsNotification(params models.CmsNotificationForm) (*models.CmsNotification, error) {
	var form models.CmsNotification
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

func (r *CmsNotificationRepo) GetCmsNotificationByID(id string) (*models.CmsNotification, error) {
	var notification models.CmsNotification
	var builder = queryfunc.NewCmsNotificationBuilder(queryfunc.CmsNotificationBuilderOptions{})
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cn.deleted_at IS NULL AND cn.id = ?", id)

		}).
		FirstFunc(&notification)

	return &notification, err
}

func (r *CmsNotificationRepo) MarkSeen(notificationID string) error {

	var err = r.db.Model(&models.CmsNotification{}).
		Where("id = ?", notificationID).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

func (r *CmsNotificationRepo) MarkSeenAll() error {

	var err = r.db.Model(&models.CmsNotification{}).
		Where("seen_at IS NULL").
		Update("seen_at", time.Now().Unix()).Error

	return err
}
