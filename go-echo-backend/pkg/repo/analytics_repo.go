package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
)

type AnalyticsRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewAnalyticsRepo(db *db.DB) *AnalyticsRepo {
	return &AnalyticsRepo{
		db:     db,
		logger: logger.New("repo/Analytics"),
	}
}

type PaginatePotentialOverdueInquiriesParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *AnalyticsRepo) PaginatePotentialOverdueInquiries(params PaginatePotentialOverdueInquiriesParams) *query.Pagination {
	return NewInquiryRepo(r.db).PaginateInquiry(PaginateInquiryParams{
		PaginationParams: params.PaginationParams,
		JwtClaimsInfo:    params.JwtClaimsInfo,
		PotentialOverdue: true,
		IncludeUser:      true,
		IncludeAssignee:  true,
	})
}

func (r *AnalyticsRepo) PaginateInquiriesTimeline(params PaginatePotentialOverdueInquiriesParams) *query.Pagination {
	return NewInquiryRepo(r.db).PaginateInquiry(PaginateInquiryParams{
		PaginationParams: params.PaginationParams,
		JwtClaimsInfo:    params.JwtClaimsInfo,
		PotentialOverdue: true,
		IncludeUser:      true,
		IncludeAssignee:  true,
		IncludeAuditLog:  true,
	})
}
