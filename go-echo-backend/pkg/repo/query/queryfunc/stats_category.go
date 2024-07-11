package queryfunc

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type StatsCategoriesBuilderOptions struct {
	ForRole enums.Role

	Comment string
}

func NewStatsCategoriesBuilder(options StatsCategoriesBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* StatsCategoriesBuilder - %s */ COUNT(1) AS total_records

	FROM categories c
	`

	return NewBuilder(fmt.Sprintf(rawSQL, options.Comment))
}
