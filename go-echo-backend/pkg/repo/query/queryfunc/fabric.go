package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type FabricBuilderOptions struct {
	QueryBuilderOptions
}

type FabricAlias struct {
	*models.Fabric
}

func NewFabricBuilder(options FabricBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ fb.*
	FROM fabrics fb
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM fabrics fb
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
		WithOrderBy("fb.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Fabric, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var (
				fabricIDs []string
			)
			for rows.Next() {
				var alias FabricAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				if alias.Fabric != nil {
					fabricIDs = append(fabricIDs, alias.Fabric.ID)
				}
				records = append(records, alias.Fabric)
			}
			if len(fabricIDs) > 0 {
				var relations []models.FabricInCollection
				db.Model(&models.FabricInCollection{}).Where("fabric_id IN ?", fabricIDs).Find(&relations)
				var collectionIDs []string
				var fabricCollections = make(map[string][]string)
				for _, v := range relations {
					collectionIDs = append(collectionIDs, v.FabricCollectionID)
					fabricCollections[v.FabricID] = append(fabricCollections[v.FabricID], v.FabricCollectionID)
				}
				var collections []models.FabricCollection
				db.Model(&models.FabricCollection{}).Where("id IN ?", collectionIDs).Find(&collections)
				mc := make(map[string]models.FabricCollection) // map collection
				for _, v := range collections {
					mc[v.ID] = v
				}
				for _, record := range records {
					for _, collectionID := range fabricCollections[record.ID] {
						if col, ok := mc[collectionID]; ok {
							record.FabricCollections = append(record.FabricCollections, col)
						}
					}
				}
			}
			return records, nil
		})
}
