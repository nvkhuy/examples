package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3_GetMediaToken(t *testing.T) {
	var cfg = initConfig()
	resp, err := s3.New(cfg).GetBlurhash(s3.GetBlurhashParams{
		FileKey:       "/uploads/media/ci83c8djtqd6mfut3fmg_rfq_attachments_cmjr4frb2hjfcfm7in10.png",
		ThumbnailSize: "128w",
	})
	assert.NoError(t, err)

	helper.PrintJSON(resp)
}
