package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataAnalyticGrowingTagRepo_Paginate(t *testing.T) {
	var app = initApp("dev")
	result, err := repo.NewAnalyticGrowingTagRepo(app.AnalyticDB).Paginate(repo.PaginateAnalyticGrowingTagsParams{
		PaginationParams: models.PaginationParams{
			Limit: 4,
		},
	})
	assert.NoError(t, err)
	helper.PrintJSON(result)
}
