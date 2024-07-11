package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/rotisserie/eris"
)

type PushTokenRepo struct {
	db *db.DB
}

func NewPushTokenRepo(db *db.DB) *PushTokenRepo {
	return &PushTokenRepo{
		db: db,
	}
}

type GetPushTokensParams struct {
	models.JwtClaimsInfo

	Tokens []string `json:"tokens" query:"tokens" form:"tokens"`
	UserID string   `json:"user_id" query:"user_id" form:"user_id"`
}

func (r *PushTokenRepo) GetPushTokens(params GetPushTokensParams) (models.PushTokens, error) {
	var devices []*models.PushToken
	var err = query.New(r.db, queryfunc.NewPushTokenBuilder(queryfunc.PushTokenBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			if params.UserID != "" {
				builder.Where("pt.user_id = ?", params.UserID)
			}

			if len(params.Tokens) > 0 {
				builder.Where("pt.token IN ?", params.Tokens)
			}
		}).
		FindFunc(&devices)

	return devices, err
}

type PaginatePushTokensParams struct {
	models.PaginationParams

	models.JwtClaimsInfo

	Tokens []string `json:"tokens" query:"tokens" form:"tokens"`
	UserID string   `json:"user_id" query:"user_id" form:"user_id"`
}

func (r *PushTokenRepo) PaginatePushTokens(params PaginatePushTokensParams) *query.Pagination {
	var result = query.New(r.db, queryfunc.NewPushTokenBuilder(queryfunc.PushTokenBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			if params.UserID != "" {
				builder.Where("pt.user_id = ?", params.UserID)
			}

			if len(params.Tokens) > 0 {
				builder.Where("pt.token IN ?", params.Tokens)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *PushTokenRepo) AddPushToken(userID string, form models.PushTokenCreateForm) (*models.PushToken, error) {
	var pushToken = models.PushToken{
		UserID:   userID,
		Token:    form.Token,
		Platform: form.Platform,
		LastUsed: r.db.NowFunc().Unix(),
	}

	var findDevice models.PushToken
	var err = r.db.Model(&models.PushToken{}).Where("token = ?", pushToken.Token).First(&findDevice).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			var devices []models.PushToken
			err = r.db.Order("last_used desc").Limit(10).Offset(10).Where("user_id = ?", pushToken.UserID).Find(&devices).Error
			if err != nil {
				return nil, err
			}

			for _, d := range devices {
				r.db.Unscoped().Model(models.PushToken{}).Where("token = ?", d.Token).Delete(&models.PushToken{})
			}

			err = r.db.Create(&pushToken).Error
			if err != nil {
				if duplicated, _ := r.db.IsDuplicateConstraint(err); duplicated {
					return &pushToken, nil
				}
				return nil, err
			}

			return &pushToken, nil
		}

		return nil, eris.Wrap(err, err.Error())
	}

	err = r.db.Model(&findDevice).Updates(&pushToken).Error
	if err != nil {
		return nil, err
	}

	return &findDevice, nil
}

func (r *PushTokenRepo) DeletePushToken(token string) error {
	var err = r.db.Model(&models.PushToken{}).Where("token = ?", token).Delete(&models.PushToken{}).Error
	return err
}
