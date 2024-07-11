package repo

import (
	"errors"
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm/clause"
)

type RWDFabricPriceRepo struct {
	db       *db.DB
	logger   *logger.Logger
	sheetAPI *sheets.Service
}

func NewRWDFabricPriceRepo(db *db.DB) *RWDFabricPriceRepo {
	return &RWDFabricPriceRepo{
		db:     db,
		logger: logger.New("repo/RWDFabricPrice"),
	}
}

func (r *RWDFabricPriceRepo) WithSheetAPI(api *sheets.Service) *RWDFabricPriceRepo {
	r.sheetAPI = api
	return r
}

func (r *RWDFabricPriceRepo) FetchPrice(params *models.FetchRWDFabricPriceParams) (result models.RWDFabricPriceSlice, err error) {
	if r.sheetAPI == nil {
		err = errors.New("empty sheet API")
		return
	}
	if !params.MaterialType.IsInvalid() {
		err = errors.New("invalid fabric material type")
		return
	}
	params = params.Fetch()
	readRange := fmt.Sprintf("%s!%s:%s", params.SheetName, params.From, params.To)
	resp, err := r.sheetAPI.Spreadsheets.Values.Get(params.SpreadsheetId, readRange).Do()
	if err != nil {
		return
	}
	for id, row := range resp.Values {
		pp := models.RWDFabricPriceFromSlice(id+params.FromNum, params.MaterialType, row)
		result = append(result, pp)
	}
	if len(result) == 0 {
		return
	}
	err = r.db.Clauses(clause.OnConflict{ // Upsert
		Columns:   []clause.Column{{Name: "row_id"}, {Name: "material_type"}},
		UpdateAll: true,
	}).Model(&result).Create(&result).Error
	return
}

func (r *RWDFabricPriceRepo) Paginate(params *models.PaginateRWDFabricPriceParams) (result *query.Pagination) {
	var builder = queryfunc.NewRWDFabricPriceBuilder(queryfunc.RWDFabricPriceBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			queryfunc.RWDFabricPriceBuilderWhereFunc(builder, params)
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()
	return
}

func (r *RWDFabricPriceRepo) Vine(params *models.PaginateRWDFabricPriceParams) (result *query.Pagination) {
	var builder = queryfunc.NewRWDFabricPriceVineBuilder(queryfunc.RWDFabricPriceBuilderOptions{
		VineSlice: params.ToVineSlice(),
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			queryfunc.RWDFabricPriceBuilderWhereFunc(builder, params)
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()
	return
}
