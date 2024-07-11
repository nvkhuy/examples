package queryfunc

import (
	"fmt"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/gorm/clause"
)

type QueryBuilderOptions struct {
	Role    enums.Role
	Comment string
}

type BuilderFunc = func(builder *Builder)

type Handler = func(db *db.DB, rawSQL *db.DB) (interface{}, error)

type Builder struct {
	rawSQL       string
	countRawSQL  string
	orderBy      string
	groupBy      string
	having       string
	isWrapJSON   bool
	isWrapSelect bool
	withoutCount bool
	paginationFn Handler
	clauses      []clause.Expression
	options      interface{}
}

func NewBuilder(rawSQL string, countRawSQL ...string) *Builder {
	var b = &Builder{
		rawSQL:       rawSQL,
		countRawSQL:  "",
		orderBy:      "",
		groupBy:      "",
		isWrapJSON:   false,
		isWrapSelect: false,
	}
	if len(countRawSQL) > 0 {
		b.countRawSQL = countRawSQL[0]
	}
	return b
}

func (b *Builder) WithOptions(options interface{}, funcs ...template.FuncMap) *Builder {
	b.options = options
	b.rawSQL = b.GetRawSQLTemplateStr(funcs...)
	if b.countRawSQL != "" {
		b.countRawSQL = b.GetCountRawSQLTemplateStr(funcs...)
	}
	return b
}

func (b *Builder) WithClauses(clauses ...clause.Expression) *Builder {
	b.clauses = clauses
	return b
}

func (b *Builder) WithPaginationFunc(fn Handler) *Builder {
	b.paginationFn = fn
	return b
}

func (b *Builder) WithRawSQL(sql string) *Builder {
	b.rawSQL = sql
	return b
}

func (b *Builder) WithWrapJSON(wrap bool) *Builder {
	b.isWrapJSON = wrap
	return b
}

func (b *Builder) WithWrapSelect(wrap bool) *Builder {
	b.isWrapSelect = wrap
	return b
}

func (b *Builder) WithCountRawSQL(sql string) *Builder {
	b.countRawSQL = sql
	return b
}

func (b *Builder) WithOrderBy(sql string) *Builder {
	b.orderBy = sql
	return b
}

func (b *Builder) WithGroupBy(sql string) *Builder {
	b.groupBy = sql
	return b
}

func (b *Builder) WithHaving(sql string) *Builder {
	b.having = sql
	return b
}

func (b *Builder) WithoutCount(withoutCount bool) *Builder {
	b.withoutCount = withoutCount
	return b
}

/* Implement interface */
func (b *Builder) GetRawSQL() string {
	return b.rawSQL
}

func (b *Builder) GetCountRawSQL() string {
	return b.countRawSQL
}

func (b *Builder) GetGroupByRawSQL() string {
	return b.groupBy
}

func (b *Builder) GetOrderByRawSQL() string {
	return b.orderBy
}

func (b *Builder) GetHavingRawSQL() string {
	return b.having
}

func (b *Builder) IsWrapJSON() bool {
	return b.isWrapJSON
}

func (b *Builder) IsWrapSelect() bool {
	return b.isWrapSelect
}

func (b *Builder) GetPaginationFunc() Handler {
	return b.paginationFn
}

func (b *Builder) GetClauses() []clause.Expression {
	return b.clauses
}

func (b *Builder) GetRawSQLTemplateStr(funcs ...template.FuncMap) string {
	return helper.TemplateStr(b.rawSQL, b.options, funcs...)
}

func (b *Builder) GetCountRawSQLTemplateStr(funcs ...template.FuncMap) string {
	return helper.TemplateStr(b.countRawSQL, b.options, funcs...)
}

func GetCaller() string {
	var frame = 17
	return fmt.Sprintf("%s:%s:%s", helper.GetFuncName(frame), helper.GetFuncName(frame+1), helper.GetFuncName(frame+2))
}
