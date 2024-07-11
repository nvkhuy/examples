package tests

import (
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/ai"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAI_ClassifyImage(t *testing.T) {
	params := ai.ImageClassifyParams{
		Image:      "https://bananarepublic.gap.com/webcontent/0052/503/816/cn52503816.jpg",
		Size:       640,
		Confidence: 0.5,
		Overlap:    0.1,
	}
	resp, err := ai.ClassifyImage(params)
	assert.NoError(t, err)
	helper.PrintJSON(resp)
}

func TestAI_ClassifyImage_Multi(t *testing.T) {
	params := []ai.ImageClassifyParams{
		{
			Image:      "https://bananarepublic.gap.com/webcontent/0052/503/816/cn52503816.jpg",
			Size:       640,
			Confidence: 0.5,
			Overlap:    0.1,
		},
		{
			Image:      "https://bananarepublic.gap.com/webcontent/0052/503/816/cn52503816.jpg",
			Size:       640,
			Confidence: 0.3,
			Overlap:    0.1,
		},
		{
			Image:      "https://bananarepublic.gap.com/webcontent/0052/503/816/cn52503816.jpg",
			Size:       640,
			Confidence: 0.6,
			Overlap:    0,
		},
	}
	resp, err := ai.ClassifyMultiImage(params)
	assert.NoError(t, err)
	helper.PrintJSON(resp)
}

func TestAI_PredictSeries(t *testing.T) {
	series := []float64{46, 31, 59, 64, 44, 104, 91, 82, 61, 46, 23, 31, 39, 44, 42, 44, 47, 55, 70, 31}
	next := 5
	result := ai.PredictSeries(series, next)
	fmt.Println(result)
}
