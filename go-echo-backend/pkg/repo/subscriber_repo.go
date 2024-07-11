package repo

import (
	"database/sql"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"

	"github.com/rotisserie/eris"

	"gorm.io/gorm/clause"
)

type SubscriberRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewSubscriberRepo(db *db.DB) *SubscriberRepo {
	return &SubscriberRepo{
		db:     db,
		logger: logger.New("repo/subscriber"),
	}
}

type PaginateSubscribersParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	Email string `json:"email" query:"email" form:"email"`
}

func (r *SubscriberRepo) PaginateSubscribers(params PaginateSubscribersParams) *query.Pagination {
	var builder = queryfunc.NewSubscriberBuilder(queryfunc.SubscriberBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.Email != "" {
				builder.Where("sub.email ILIKE ?", "%"+params.Email)
			}

			if strings.TrimSpace(params.Keyword) != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("email ILIKE @query", sql.Named("query", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type SearchSubscribersParams struct {
	models.PaginationParams

	Roles []string `json:"roles" query:"roles" form:"roles"`

	AccountStatuses []enums.AccountStatus `json:"account_statuses" query:"account_statuses" form:"account_statuses"`

	ForRole enums.Role
}

func (r *SubscriberRepo) CreateSubscriber(form models.SubscribeByEmailForm) (*models.Subscriber, error) {
	var subscriber models.Subscriber
	err := copier.Copy(&subscriber, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Clauses(clause.OnConflict{DoNothing: true}).
		Omit(clause.Associations).
		Create(&subscriber).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, eris.Wrap(err, "")
	}

	return &subscriber, nil
}

func (r *SubscriberRepo) GetSubscriberByID(SubscriberID string, options queryfunc.SubscriberBuilderOptions) (*models.Subscriber, error) {
	var builder = queryfunc.NewSubscriberBuilder(options)
	var Subscriber models.Subscriber
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("sub.id = ?", SubscriberID)
		}).
		FirstFunc(&Subscriber)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &Subscriber, nil
}

func (r *SubscriberRepo) UpdateSubscriberByID(SubscriberID string, form models.SubscribeByEmailForm) (*models.Subscriber, error) {
	var update models.Subscriber

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	err = r.db.Omit(clause.Associations).Model(&models.Subscriber{}).Where("id = ?", SubscriberID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	return r.GetSubscriberByID(SubscriberID, queryfunc.SubscriberBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: form.GetRole(),
		},
	})
}

func (r *SubscriberRepo) ArchiveSubscriberByID(SubscriberID string) error {
	var err = r.db.Delete(&models.Subscriber{}, "id = ?", SubscriberID).Error

	return err
}

func (r *SubscriberRepo) UnarchiveSubscriberByID(SubscriberID string) error {
	var err = r.db.Model(&models.Subscriber{}).Where("id = ?", SubscriberID).UpdateColumn("deleted_at", nil).Error

	return err
}
