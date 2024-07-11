package queryfunc

import (
	"strings"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
)

type Vine struct {
	Attribute string `json:"attribute"`
}

type VineAlias struct {
	*Vine
}

type ProductTypesPriceAlias struct {
	*models.ProductTypePrice
	Prices models.ProductTypePriceSlice `json:"prices"`
}

type ProductTypesPriceBuilderOptions struct {
	QueryBuilderOptions
	VineSlice []string
}

func NewProductTypesPriceBuilder(options ProductTypesPriceBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*

	FROM product_type_prices p
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM product_type_prices p
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
			var records = make(models.ProductTypePriceSlice, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias ProductTypesPriceAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.ProductTypePrice)
			}

			return &records, nil
		})
}

func ProductTypesPriceBuilderWhereFunc(b *query.Builder, p *models.PaginateProductTypesPriceParams) *query.Builder {
	if p.Product != nil {
		b.Where("p.product = ?", p.Product)
	}

	if p.Gender != nil {
		b.Where("p.gender = ?", p.Gender)
	}

	if p.FabricType != nil {
		b.Where("p.fabric_type = ?", p.FabricType)
	}

	if p.Feature != nil {
		b.Where("p.feature = ?", p.Feature)
	}

	if p.Category != nil {
		b.Where("p.category = ?", p.Category)
	}

	if p.Item != nil {
		b.Where("p.item = ?", p.Item)
	}

	if p.Form != nil {
		b.Where("p.form = ?", p.Form)
	}

	if p.Description != nil {
		b.Where("p.description = ?", p.Description)
	}

	if p.KnitMaterial != nil {
		b.Where("p.knit_material = ?", p.KnitMaterial)
	}

	if p.KnitComposition != nil {
		b.Where("p.knit_composition = ?", p.KnitComposition)
	}

	if p.KnitWeight != nil {
		b.Where("p.knit_weight = ?", p.KnitWeight)
	}

	if p.WovenMaterial != nil {
		b.Where("p.woven_material = ?", p.WovenMaterial)
	}

	if p.WovenComposition != nil {
		b.Where("p.woven_composition = ?", p.WovenComposition)
	}

	if p.FabricConsumption != nil {
		b.Where("p.fabric_consumption = ?", p.FabricConsumption)
	}
	return b
}

func NewProductTypesPriceVineBuilder(options ProductTypesPriceBuilderOptions) *Builder {
	fields := strings.Join(options.VineSlice, ", ")
	last := options.VineSlice[len(options.VineSlice)-1]
	var rawSQL = `
	SELECT /* {{Description}} */ _v_ as attribute

	FROM product_type_prices p
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM product_type_prices p
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
				if alias.Vine == nil {
					continue
				}
				records = append(records, alias.Attribute)
			}

			return &records, nil
		})
}
