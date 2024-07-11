package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type ProductAlias struct {
	*models.Product
	Category *models.Category    `gorm:"embedded;embeddedPrefix:ct__"`
	Images   *models.Attachments `json:"images"`
}

type ProductBuilderOptions struct {
	QueryBuilderOptions
	GetByCollection bool
	GetDetails      bool
}

func NewProductBuilder(options ProductBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*,
	ct.id AS ct__id,
	ct.name AS ct__name,
	ct.slug AS ct__slug

	FROM products p
	LEFT JOIN categories ct ON p.category_id = ct.id
    LEFT JOIN categories pr ON ct.parent_category_id = pr.id
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM products p
	LEFT JOIN categories ct ON p.category_id = ct.id
    LEFT JOIN categories pr ON ct.parent_category_id = pr.id
	`

	if options.GetByCollection {
		rawSQL = `
		SELECT /* {{Description}} */ p.*

		FROM products p
		JOIN collection_products c ON p.id = c.product_id
		`
		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM products p
		JOIN collection_products c ON p.id = c.product_id
		`
	}

	builder := NewBuilder(rawSQL, countSQL).
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
			var records = make([]*models.Product, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias ProductAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				alias.Product.Category = alias.Category
				records = append(records, alias.Product)
			}

			var (
				productIDs      []string
				fabricIDs       []string
				fabricProductID = make(map[string]string)
			)

			for _, record := range records {
				productIDs = append(productIDs, record.ID)
				for _, fabricID := range record.FabricIDs {
					fabricIDs = append(fabricIDs, fabricID)
					fabricProductID[fabricID] = record.ID
				}
			}

			var wg sync.WaitGroup

			if len(productIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var variants []*models.Variant
					err = db.Find(&variants, "product_id IN ?", productIDs).Error
					for _, record := range records {
						for _, variant := range variants {
							if record.ID == variant.ProductID {
								record.Variants = append(record.Variants, variant)
							}
						}
					}
				}()

			}

			if len(fabricIDs) > 0 {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var fabrics []*models.Fabric
					err = db.Find(&fabrics, "id IN ?", fabricIDs).Error
					for _, record := range records {
						for _, fabric := range fabrics {
							if pID, ok := fabricProductID[fabric.ID]; ok && pID == record.ID {
								record.Fabrics = append(record.Fabrics, fabric)
							}
						}
					}
				}()

			}

			wg.Wait()

			return records, nil
		})
	return builder
}

type ProductRecommendBuilderOptions struct {
	QueryBuilderOptions
	GetByCollection bool
	GetDetails      bool
}

func NewProductRecommendBuilder(options ProductRecommendBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*,
	ct.id AS ct__id,
	ct.name AS ct__name,
	ct.slug AS ct__slug

	FROM products p
	LEFT JOIN categories ct ON p.category_id = ct.id
    LEFT JOIN categories pr ON ct.parent_category_id = pr.id
	LEFT JOIN product_classes pc ON p.id = pc.product_id
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM products p
	LEFT JOIN categories ct ON p.category_id = ct.id
    LEFT JOIN categories pr ON ct.parent_category_id = pr.id
	LEFT JOIN product_classes pc ON p.id = pc.product_id
	`

	if options.GetByCollection {
		rawSQL = `
		SELECT /* {{Description}} */ p.*

		FROM products p
		JOIN collection_products c ON p.id = c.product_id
		`
		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM products p
		JOIN collection_products c ON p.id = c.product_id
		`
	}

	builder := NewBuilder(rawSQL, countSQL).
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
			var records = make([]*models.Product, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias ProductAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				alias.Product.Category = alias.Category
				records = append(records, alias.Product)
			}

			var (
				productIDs      []string
				fabricIDs       []string
				fabricProductID = make(map[string]string)
			)

			for _, record := range records {
				productIDs = append(productIDs, record.ID)
				for _, fabricID := range record.FabricIDs {
					fabricIDs = append(fabricIDs, fabricID)
					fabricProductID[fabricID] = record.ID
				}
			}

			if len(productIDs) > 0 {
				var variants []*models.Variant
				err = db.Find(&variants, "product_id IN ?", productIDs).Error
				for _, record := range records {
					for _, variant := range variants {
						if record.ID == variant.ProductID {
							record.Variants = append(record.Variants, variant)
						}
					}
				}
			}

			if len(fabricIDs) > 0 {
				var fabrics []*models.Fabric
				err = db.Find(&fabrics, "id IN ?", fabricIDs).Error
				for _, record := range records {
					for _, fabric := range fabrics {
						if pID, ok := fabricProductID[fabric.ID]; ok && pID == record.ID {
							record.Fabrics = append(record.Fabrics, fabric)
						}
					}
				}
			}

			if len(productIDs) > 0 {
				var classes []*models.ProductClass
				err = db.Find(&classes).Where("product_id IN ?", productIDs).Error
				for _, record := range records {
					for _, class := range classes {
						if record.ID == class.ProductID {
							record.ProductClasses = append(record.ProductClasses, class)
						}
					}
				}
			}

			return records, nil
		})
	return builder
}
