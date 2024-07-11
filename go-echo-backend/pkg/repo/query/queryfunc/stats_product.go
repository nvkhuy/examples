package queryfunc

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type StatsProductsBuilderOptions struct {
	ForRole enums.Role

	Comment string
}

func NewStatsProductsBuilder(options StatsProductsBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* StatsProductsBuilder - %s */ COUNT(1) AS total_records

	FROM products p
	`

	return NewBuilder(fmt.Sprintf(rawSQL, options.Comment))
}
