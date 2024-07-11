package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTNARepo_Delete(t *testing.T) {
	var app = initApp("dev")
	var err = repo.NewTNARepo(app.DB).Delete(repo.DeleteTNAsParams{
		ID: "cn6s883kl4pkf4p8bnm1",
	})
	if err != nil {
		return
	}
	assert.NoError(t, err)
}
