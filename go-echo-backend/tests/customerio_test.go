package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/stretchr/testify/assert"
)

func TestCustomerio_GetActivities(t *testing.T) {
	var cfg = initConfig("prod")
	_, err := customerio.New(cfg).GetActivities("cjip00ifleak1ibumcf0", customerio.GetActivitiesParams{
		Type: "event",
		Name: "track_activity",
	})
	assert.NoError(t, err)
	// helper.PrintJSON(result)
}
