package queryfunc

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type InventoryBuilderOptions struct {
	ForRole enums.Role

	IncludeShop bool

	Comment string
}

type InventoryAlias struct {
	*models.Variant

	Product *models.Product `gorm:"embedded;embeddedPrefix:p__"`
}

func NewInventoryBuilder(options InventoryBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* InventoryBuilder - %s */ v.*,
	p.id AS p__id,
	p.name AS p__name,
	p.short_description AS p__short_description,
	p.description AS p__description,
	p.shop_id AS p__shop_id,
	p.category_id AS p__category_id,
	p.sku AS p__sku,
	p.price AS p__price,
	p.sold_quantity AS p__sold_quantity,
	p.product_type AS p__product_type,
	p.ready_to_ship AS p__ready_to_ship,
	p.daily_deal AS p__daily_deal,
	p.rating_count AS p__rating_count,
	p.rating_star AS p__rating_star,
	p.trade_unit AS p__trade_unit,
	p.min_order AS p__min_order,
	p.attachments AS p__attachments

	FROM variants v
	JOIN products p ON p.id = v.product_id
	`

	return NewBuilder(fmt.Sprintf(rawSQL, options.Comment)).
		WithOrderBy("v.stock ASC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Variant, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var copy InventoryAlias
				err = db.ScanRows(rows, &copy)
				if err != nil {
					continue
				}

				copy.Variant.Product = copy.Product
				records = append(records, copy.Variant)
			}

			return &records, nil
		})
}
