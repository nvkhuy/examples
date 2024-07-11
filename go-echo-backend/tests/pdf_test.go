package tests

import (
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"os"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/pdf"
	"github.com/stretchr/testify/assert"
)

func TestPDF_GetPDF(t *testing.T) {
	logger.Init()
	var cfg = initConfig()

	var params = pdf.GetPDFParams{
		URL:           fmt.Sprintf("%s/bulks/%s/order-closing", cfg.AdminPortalBaseURL, "cmsf8lnv30js1m62796g"),
		Selector:      "#order-closing-ready-to-print",
		Landscape:     true,
		DisableLocker: true,
	}

	pdfRawData, err := pdf.New(cfg).GetPDF(params)
	assert.NoError(t, err)

	helper.PrintJSON(params)

	_ = os.WriteFile("test.pdf", pdfRawData, 0777)
}
