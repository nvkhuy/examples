package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type FabricCollectionBuilderOptions struct {
	QueryBuilderOptions
}

type FabricCollectionAlias struct {
	*models.FabricCollection
}

func NewFabricCollectionBuilder(options FabricCollectionBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ fc.*
	FROM fabric_collections fc
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM fabric_collections fc
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
		WithOrderBy("fc.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.FabricCollection, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var collectionIDs []string
			for rows.Next() {
				var alias FabricCollectionAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				if alias.FabricCollection != nil {
					collectionIDs = append(collectionIDs, alias.FabricCollection.ID)
				}
				records = append(records, alias.FabricCollection)
			}
			if len(collectionIDs) > 0 {
				var fabricInCollections []models.FabricInCollection
				db.Model(&models.FabricInCollection{}).
					Where("fabric_collection_id IN ?", collectionIDs).Find(&fabricInCollections)
				collectionFabricIDs := make(map[string][]string)
				var fabricIds []string
				for _, v := range fabricInCollections {
					collectionFabricIDs[v.FabricCollectionID] = append(collectionFabricIDs[v.FabricCollectionID], v.FabricID)
					fabricIds = append(fabricIds, v.FabricID)
				}
				var fabrics []models.Fabric
				db.Model(&models.Fabric{}).
					Where("id IN ?", fabricIds).Find(&fabrics)
				m1 := make(map[string]models.Fabric)
				for _, v := range fabrics {
					m1[v.ID] = v
				}
				for _, record := range records {
					for _, fabricID := range collectionFabricIDs[record.ID] {
						if _, ok := m1[fabricID]; ok {
							record.Fabrics = append(record.Fabrics, m1[fabricID])
						}
					}
				}
			}

			return records, nil
		})
}
