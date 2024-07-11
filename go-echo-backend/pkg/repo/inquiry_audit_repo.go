package repo

import (
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

type InquiryAuditRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewInquiryAuditRepo(db *db.DB) *InquiryAuditRepo {
	return &InquiryAuditRepo{
		db:     db,
		logger: logger.New("repo/InquiryAudit"),
	}
}

func (r *InquiryAuditRepo) CreateInquiryAudit(params models.InquiryAuditCreateForm) error {
	var form models.InquiryAudit
	err := copier.Copy(&form, &params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Create(&form).Error
	if err != nil {
		return eris.Wrap(err, "")
	}
	return err
}

type GetQuotationLogParams struct {
	models.JwtClaimsInfo

	InquiryID  string                `json:"inquiry_id" query:"inquiry_id" form:"inquiry_id" param:"inquiry_id"`
	UserID     string                `json:"user_id" query:"user_id" form:"user_id" param:"user_id"`
	ActionType enums.AuditActionType `json:"action_type" query:"action_type" form:"action_type" param:"action_type"`
}

func (r *InquiryAuditRepo) GetLastQuotationLog(params GetQuotationLogParams) (*models.InquiryAudit, error) {
	var builder = queryfunc.NewInquiryAuditBuilder(queryfunc.InquiryAuditBuilderOptions{})
	var lastLog models.InquiryAudit
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("ia.inquiry_id = ?", params.InquiryID)
			builder.Where("ia.action_type = ?", params.ActionType)
		}).
		OrderBy("ia.created_at DESC").
		FirstFunc(&lastLog)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &lastLog, nil
}
