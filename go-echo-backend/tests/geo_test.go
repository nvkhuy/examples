package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/geo"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/stretchr/testify/assert"
)

func TestGeo_Search(t *testing.T) {
	var cfg = initConfig()
	var input = "Marasi Drive 8, Dubai, United Arab Emirates"
	result, err := geo.New(cfg).SearchPlaceIndexForText(geo.SearchPlaceIndexForTextParams{
		Address:     input,
		CountryCode: "VN",
		Language:    enums.LanguageCodeEnglish,
	})

	assert.NoError(t, err)

	helper.PrintJSON(result)
}
