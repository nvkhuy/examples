package repo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"gorm.io/gorm/clause"
	"time"
)

type SysNotificationRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewSysNotificationRepo(db *db.DB) *SysNotificationRepo {
	return &SysNotificationRepo{
		db:     db,
		logger: logger.New("repo/SysNotification"),
	}
}

type CreateSysNotificationsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	models.SysNotification
}

func (r *SysNotificationRepo) Create(params CreateSysNotificationsParams) (result *models.SysNotification, err error) {
	if err = r.db.Model(&models.SysNotification{}).Clauses(clause.Returning{}).Create(&params.SysNotification).Error; err != nil {
		return
	}
	result = &params.SysNotification
	return
}

type PaginateSysNotificationsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *SysNotificationRepo) Paginate(params PaginateSysNotificationsParams) *query.Pagination {
	var builder = queryfunc.NewSysNotificationBuilder(queryfunc.SysNotificationBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	subQuery := r.db.Model(&models.UserSysNotification{}).
		Select("sys_notification_id").
		Where("user_id = ? AND seen_at is not NULL", params.GetUserID())

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("sn.id NOT IN (?)", subQuery)
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type MarkSeenSysNotificationParams struct {
	models.JwtClaimsInfo

	NotificationID string `json:"notification_id" query:"notification_id" param:"notification_id" validate:"required"`
}

func (r *SysNotificationRepo) MarkSeen(params MarkSeenSysNotificationParams) (err error) {

	var seenSysNotification = &models.UserSysNotification{
		SysNotificationID: params.NotificationID,
		SeenAt:            aws.Int64(time.Now().Unix()),
		UserID:            params.GetUserID(),
	}
	var seenId string
	r.db.Model(&models.UserSysNotification{}).
		Select("id").
		Where("user_id = ? and sys_notification_id = ? and seen_at is not null", params.GetUserID(), params.NotificationID).
		First(&seenId)
	if seenId != "" {
		return
	}

	err = r.db.Model(&models.UserSysNotification{}).Create(&seenSysNotification).Error
	return err
}

type MarkSeenAllSysNotificationParams struct {
	models.JwtClaimsInfo
}

func (r *SysNotificationRepo) MarkSeenAll(params MarkSeenAllSysNotificationParams) (err error) {
	subQuery := r.db.Model(&models.UserSysNotification{}).
		Select("sys_notification_id").
		Where("user_id = ? AND seen_at is not NULL", params.GetUserID())

	var Ids []string // sys notification ids to mark seen
	err = r.db.Model(&models.SysNotification{}).
		Select("id").
		Where("id NOT IN (?)", subQuery).
		Find(&Ids).Error

	var seens []*models.UserSysNotification
	for _, sysNotifyID := range Ids {
		seens = append(seens, &models.UserSysNotification{
			SysNotificationID: sysNotifyID,
			SeenAt:            aws.Int64(time.Now().Unix()),
			UserID:            params.GetUserID(),
		})
	}
	if seens == nil {
		return
	}
	err = r.db.Model(&models.UserSysNotification{}).Create(seens).Error
	return err
}
