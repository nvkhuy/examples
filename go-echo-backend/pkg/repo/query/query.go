package query

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QueryBuilder interface {
	GetRawSQL() string
	GetCountRawSQL() string
	GetGroupByRawSQL() string
	GetOrderByRawSQL() string
	GetHavingRawSQL() string
	GetPaginationFunc() ExecFunc
	IsWrapJSON() bool
	IsWrapSelect() bool
	GetClauses() []clause.Expression
}

type ExecFunc = func(db *db.DB, rawSQL *db.DB) (interface{}, error)

type WhereFunc = func(builder *Builder)

type Pagination struct {
	HasNext            bool        `json:"has_next,omitempty"`
	HasPrev            bool        `json:"has_prev,omitempty"`
	PerPage            int         `json:"per_page,omitempty"`
	NextPage           int         `json:"next_page,omitempty"`
	Page               int         `json:"current_page,omitempty"`
	PrevPage           int         `json:"prev_page,omitempty"`
	Offset             int         `json:"offset,omitempty"`
	Records            interface{} `json:"records,omitempty"`
	TotalRecord        int         `json:"total_record,omitempty"`
	TotalPage          int         `json:"total_page,omitempty"`
	Metadata           interface{} `json:"metadata,omitempty"`
	TotalCurrentRecord int         `json:"total_current_record,omitempty"`
}

type Builder struct {
	db               *db.DB
	RawSQLString     string
	countRawSQL      string
	limit            int
	page             int
	hasWhere         bool
	whereValues      []interface{}
	namedWhereValues map[string]interface{}
	orderBy          string
	groupBy          string
	having           string
	wrapJSON         bool
	wrapSelect       bool
	withoutCount     bool
	qf               QueryBuilder
	clauses          []clause.Expression
}

func New(db *db.DB, qf QueryBuilder) *Builder {
	var builder = &Builder{
		db:               db,
		RawSQLString:     qf.GetRawSQL(),
		whereValues:      []interface{}{},
		namedWhereValues: map[string]interface{}{},
		hasWhere:         false,
		orderBy:          qf.GetOrderByRawSQL(),
		groupBy:          qf.GetGroupByRawSQL(),
		wrapJSON:         qf.IsWrapJSON(),
		wrapSelect:       qf.IsWrapSelect(),
		countRawSQL:      qf.GetCountRawSQL(),
		having:           qf.GetHavingRawSQL(),
		qf:               qf,
		clauses:          qf.GetClauses(),
	}

	return builder
}

func (b *Builder) WithWrapJSON(isWrapJSON bool) *Builder {
	b.wrapJSON = isWrapJSON
	return b
}

func (b *Builder) WithWrapSelect(isWrapSelect bool) *Builder {
	b.wrapSelect = isWrapSelect
	return b
}

func (b *Builder) Where(query interface{}, args ...interface{}) *Builder {
	switch value := query.(type) {
	case map[string]interface{}:
		for key, v := range value {
			b.namedWhereValues[key] = v
		}
	case map[string]string:
		for key, v := range value {
			b.namedWhereValues[key] = v
		}
	case sql.NamedArg:
		b.namedWhereValues[value.Name] = value.Value
	default:
		if len(args) > 0 {
			b.whereValues = append(b.whereValues, args...)
		}

		if b.hasWhere {
			b.RawSQLString = fmt.Sprintf("%s AND %v", b.RawSQLString, query)
			if b.countRawSQL != "" {
				b.countRawSQL = fmt.Sprintf("%s AND %v", b.countRawSQL, query)
			}
		} else {
			b.RawSQLString = fmt.Sprintf("%s WHERE %v", b.RawSQLString, query)
			if b.countRawSQL != "" {
				b.countRawSQL = fmt.Sprintf("%s WHERE %v", b.countRawSQL, query)
			}
			b.hasWhere = true

		}

	}

	return b
}

func (b *Builder) OrderBy(orderBy ...string) *Builder {
	if len(orderBy) > 0 {
		b.orderBy = strings.Join(orderBy, ",")
	}
	return b
}

func (b *Builder) GroupBy(groupBy string) *Builder {
	b.groupBy = groupBy
	return b
}

func (b *Builder) WhereFunc(f WhereFunc) *Builder {
	f(b)
	return b
}

func (b *Builder) WithClauses(clauses ...clause.Expression) *Builder {
	b.clauses = clauses
	return b
}

func (b *Builder) Limit(limit int) *Builder {
	b.limit = limit

	return b
}

func (b *Builder) Page(page int) *Builder {
	b.page = page
	return b
}

func (b *Builder) WithoutCount(withoutCount bool) *Builder {
	b.withoutCount = withoutCount
	return b
}

func (b *Builder) build() (queryString string, countQuery string) {
	var rawSQLString = b.RawSQLString
	queryString = rawSQLString
	countQuery = b.countRawSQL

	if countQuery == "" {
		countQuery = rawSQLString
	}

	if b.groupBy != "" {
		queryString = fmt.Sprintf("%s GROUP BY %s", queryString, b.groupBy)
		countQuery = fmt.Sprintf("%s GROUP BY %s", countQuery, b.groupBy)
	}

	if b.having != "" {
		queryString = fmt.Sprintf("%s HAVING %s", queryString, b.having)
		countQuery = fmt.Sprintf("%s HAVING %s", countQuery, b.having)
	}

	if b.orderBy != "" {
		queryString = fmt.Sprintf("%s ORDER BY %s", queryString, b.orderBy)
	}

	if b.limit > 0 {
		queryString = fmt.Sprintf("%s LIMIT %d", queryString, b.limit)
	}

	if b.page > 0 {
		var offset = 0
		if b.page > 1 {
			offset = (b.page - 1) * b.limit
		}

		queryString = fmt.Sprintf("%s OFFSET %d", queryString, offset)
	}

	if b.wrapJSON {
		queryString = fmt.Sprintf(`
WITH alias AS (
%s
)
SELECT to_jsonb(row_to_json(alias)) AS alias
FROM alias
		`, queryString)
	}

	if b.wrapSelect {
		queryString = fmt.Sprintf(`
SELECT * FROM (
%s
) alias
		`, queryString)
	}

	return
}

func (b *Builder) GetPagingFunc(f ...ExecFunc) ExecFunc {
	if b.qf != nil {
		return b.qf.GetPaginationFunc()
	}

	if len(f) > 0 {
		return f[0]
	}

	return nil
}

func (b *Builder) PagingFunc(f ...ExecFunc) *Pagination {
	if b.withoutCount {
		return b.PagingInfiniteFunc(f...)
	}

	if b.page < 1 {
		b.page = 1
	}
	var fn = b.GetPagingFunc(f...)
	if fn == nil {
		panic(fmt.Errorf("fn is not implement"))
	}

	var offset = (b.page - 1) * b.limit
	var done = make(chan bool, 1)
	var pagination Pagination
	var count int

	sqlString, countSQLString := b.build()

	var values = b.mergeValues()
	countSQLString = fmt.Sprintf(`
SELECT COUNT(1) 
FROM (
%s
) t
	`, countSQLString)
	var countSQL = b.db.WithGorm(b.db.Clauses(b.clauses...).Raw(countSQLString, values...))
	go b.count(countSQL, done, &count)

	result, err := fn(b.db, b.db.WithGorm(b.db.Clauses(b.clauses...).Raw(sqlString, values...)))
	if err != nil {
		b.db.CustomLogger.ErrorAny(err)
	}
	<-done
	close(done)

	pagination.TotalRecord = count
	pagination.Records = result
	pagination.Page = b.page
	pagination.Offset = offset

	if b.limit > 0 {
		pagination.PerPage = b.limit
		pagination.TotalPage = int(math.Ceil(float64(count) / float64(b.limit)))
	} else {
		pagination.TotalPage = 1
		pagination.PerPage = count
	}

	if b.page > 1 {
		pagination.PrevPage = b.page - 1
	} else {
		pagination.PrevPage = b.page
	}

	if b.page == pagination.TotalPage {
		pagination.NextPage = b.page
	} else {
		pagination.NextPage = b.page + 1
	}

	pagination.HasNext = pagination.TotalPage > pagination.Page
	pagination.HasPrev = pagination.Page > 1

	if !pagination.HasNext {
		pagination.NextPage = pagination.Page
	}

	return &pagination
}

func (b *Builder) PagingInfiniteFunc(f ...ExecFunc) *Pagination {
	if b.page < 1 {
		b.page = 1
	}
	var fn = b.GetPagingFunc(f...)
	if fn == nil {
		panic(fmt.Errorf("fn is not implement"))
	}

	var offset = (b.page - 1) * b.limit
	var pagination Pagination
	var count int

	sqlString, _ := b.build()

	var values = b.mergeValues()

	result, err := fn(b.db, b.db.WithGorm(b.db.Clauses(b.clauses...).Raw(sqlString, values...)))
	if err != nil {
		b.db.CustomLogger.ErrorAny(err)
	}

	pagination.Records = result
	pagination.Page = b.page
	pagination.Offset = offset

	if b.limit > 0 {
		pagination.PerPage = b.limit
	} else {
		pagination.PerPage = count
	}

	if b.page > 1 {
		pagination.PrevPage = b.page - 1
	} else {
		pagination.PrevPage = b.page
	}

	if b.page == pagination.TotalPage {
		pagination.NextPage = b.page
	} else {
		pagination.NextPage = b.page + 1
	}

	pagination.TotalCurrentRecord = b.countResult(result)
	pagination.HasNext = pagination.TotalCurrentRecord == pagination.PerPage
	pagination.HasPrev = pagination.Page > 1

	if !pagination.HasNext {
		pagination.NextPage = pagination.Page
	}

	return &pagination
}

func (b *Builder) FindFunc(dest interface{}, f ...ExecFunc) error {
	sqlString, _ := b.build()

	var rOut = reflect.ValueOf(dest)
	if rOut.Kind() != reflect.Ptr {
		return fmt.Errorf("must be a pointer of %T", dest)
	}

	var fn = b.GetPagingFunc(f...)
	if fn == nil {
		panic(fmt.Errorf("fn is not implement"))
	}

	var values = b.mergeValues()
	result, err := fn(b.db, b.db.WithGorm(b.db.Clauses(b.clauses...).Raw(sqlString, values...)))
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return b.copyResult(rOut, result)
}

func (b *Builder) FirstFunc(dest interface{}, f ...ExecFunc) error {
	b.limit = 1
	sqlString, _ := b.build()

	var rOut = reflect.ValueOf(dest)
	if rOut.Kind() != reflect.Ptr {
		return fmt.Errorf("must be a pointer of %T", dest)
	}

	var fn = b.GetPagingFunc(f...)
	if fn == nil {
		panic(fmt.Errorf("fn is not implement"))
	}

	var values = b.mergeValues()
	result, err := fn(b.db, b.db.WithGorm(b.db.Clauses(b.clauses...).Raw(sqlString, values...)))
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return b.copyResult(rOut, result)
}

func (b *Builder) Scan(dest interface{}) error {
	sqlString, _ := b.build()

	var values = b.mergeValues()
	var result = b.db.Clauses(b.clauses...).Raw(sqlString, values...).Scan(dest)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return sql.ErrNoRows
		}
	}

	return result.Error
}

func (b *Builder) Find(dest interface{}) error {
	sqlString, _ := b.build()

	var values = b.mergeValues()
	var result = b.db.Clauses(b.clauses...).Raw(sqlString, values...).Find(dest)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return sql.ErrNoRows
		}
	}

	return result.Error
}

func (b *Builder) ExplainSQL() string {
	sqlString, _ := b.build()

	var values = b.mergeValues()
	var stmt = b.db.Clauses(b.clauses...).Session(&gorm.Session{DryRun: true}).Raw(sqlString, values...).Statement
	return stmt.Explain(stmt.SQL.String(), stmt.Vars...)

}

func (b *Builder) ScanRow(dest interface{}) error {
	sqlString, _ := b.build()

	var values = b.mergeValues()
	var err = b.db.Clauses(b.clauses...).Raw(sqlString, values).Row().Scan(dest)
	if err != nil {
		b.db.CustomLogger.ErrorAny(err)
		return err
	}

	return nil
}

func (b *Builder) count(countSQL *db.DB, done chan bool, count *int) {
	if countSQL != nil {
		var err = countSQL.Clauses(b.clauses...).Row().Scan(count)
		if err != nil {
			b.db.CustomLogger.DebugAny(err)
		}
	}
	done <- true
}

func (b *Builder) mergeValues() []interface{} {
	var values = []interface{}{}
	values = append(values, b.whereValues...)
	if len(b.namedWhereValues) > 0 {
		values = append(values, b.namedWhereValues)
	}
	return values
}

func (b *Builder) copyResult(rOut reflect.Value, result interface{}) error {
	var rResult = reflect.ValueOf(result)

	if rResult.Kind() != reflect.Ptr {
		rResult = helper.ToPtr(rResult)

	}

	if rResult.Type() != rOut.Type() {
		switch rResult.Kind() {
		case reflect.Array, reflect.Slice:
			if rResult.Len() > 0 {
				var elem = rResult.Index(0).Elem()
				rOut.Elem().Set(elem)
				return nil
			} else {
				return sql.ErrNoRows
			}
		case reflect.Ptr:
			switch rResult.Elem().Kind() {
			case reflect.Array, reflect.Slice:
				if rResult.Elem().Len() > 0 {
					var elem = rResult.Elem().Index(0).Elem()
					rOut.Elem().Set(elem)
					return nil
				} else {
					return sql.ErrNoRows
				}
			}
		}

		return fmt.Errorf("%v is not %v", rResult.Type(), rOut.Type())
	}

	rOut.Elem().Set(rResult.Elem())

	return nil
}

func (b *Builder) countResult(result interface{}) int {
	if result == nil {
		return 0
	}

	var rResult = reflect.ValueOf(result)

	if rResult.Kind() != reflect.Ptr {
		rResult = helper.ToPtr(rResult)

	}

	switch rResult.Kind() {
	case reflect.Array, reflect.Slice:
		return rResult.Len()
	case reflect.Ptr:
		switch rResult.Elem().Kind() {
		case reflect.Array, reflect.Slice:
			return rResult.Elem().Len()
		}
	}

	return 0
}
