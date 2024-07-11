package queryfunc

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type StatsBuyersBuilderOptions struct {
	ForRole enums.Role

	Comment string
}

func NewStatsBuyersBuilder(options StatsBuyersBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* StatsBuyersBuilder - %s */ COUNT(1) AS total_records

	FROM users u
	`

	return NewBuilder(fmt.Sprintf(rawSQL, options.Comment))
}
