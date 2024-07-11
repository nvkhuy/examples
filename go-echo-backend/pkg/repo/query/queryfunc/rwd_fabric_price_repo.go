package queryfunc

import (
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
)

type RWDFabricPriceAlias struct {
	*models.RWDFabricPrice
}

type RWDFabricPriceBuilderOptions struct {
	QueryBuilderOptions
	VineSlice []string
}

func NewRWDFabricPriceBuilder(options RWDFabricPriceBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*

	FROM rwd_fabric_prices p
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM rwd_fabric_prices p
	`

	return NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make(models.RWDFabricPriceSlice, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias RWDFabricPriceAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.RWDFabricPrice)
			}

			return &records, nil
		})
}

func RWDFabricPriceBuilderWhereFunc(b *query.Builder, p *models.PaginateRWDFabricPriceParams) *query.Builder {
	if p.FabricType != nil {
		_fabricType := strings.TrimSpace(strings.ToLower(aws.StringValue(p.FabricType)))
		b.Where("p.material_type = ?", _fabricType)
	}

	if p.Material != nil {
		b.Where("p.material = ?", p.Material)
	}

	if p.Composition != nil {
		b.Where("p.composition = ?", p.Composition)
	}

	if p.Weight != nil {
		b.Where("p.weight = ?", p.Weight)
	}

	if p.CutWidth != nil {
		b.Where("p.cut_width = ?", p.CutWidth)
	}
	return b
}

func NewRWDFabricPriceVineBuilder(options RWDFabricPriceBuilderOptions) *Builder {
	fields := strings.Join(options.VineSlice, ", ")
	last := options.VineSlice[len(options.VineSlice)-1]
	var rawSQL = `
	SELECT /* {{Description}} */ _v_ as attribute

	FROM rwd_fabric_prices p
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM rwd_fabric_prices p
	`
	var groupBy = "_v_"

	rawSQL = strings.Replace(rawSQL, "_v_", last, 1)
	groupBy = strings.Replace(groupBy, "_v_", fields, 1)

	return NewBuilder(rawSQL, countSQL).
		WithGroupBy(groupBy).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]interface{}, rawSQL.RowsAffected)
			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias VineAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				records = append(records, alias.Attribute)
			}

			return &records, nil
		})
}
