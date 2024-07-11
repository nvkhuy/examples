package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/rsa"
	"github.com/stretchr/testify/assert"
)

func TestRSA_Validate(t *testing.T) {
	var cfg = initConfig("prod")
	var err = rsa.New(cfg).Validate("8dsJUmDn36N8FygoUbSJKQETqU5/0YC2SR0b/QUx/Qw1nKBlooF5ibZcXXku17NDdapxRnLIxW9b+DoGdguO2zLNv/OSW6n3iPvHwLjQqE9ZdakgjzKTEaLnkiwQExVJn15lgW6P7V12pyJFwuzIPIszN6Fr42ZaE/h20A40vD32RaUh9ole+0Id+ISiQ8xdCDWWPeDXvPUyc4TMY1RUFsLaGqN4cWHZI1fSsWGrYxFQoQ9MMI56i+9wzOdIc+A63ayINkcZJfLkthUEN80RsxhIPSkypDxXVd/iHEyTOcstxEyikmbBh0SKsACm/GUzKoEgFFA4PsX3VpIWD6auUQ==")

	assert.NoError(t, err)

}

func TestRSA_GenRsaKey(t *testing.T) {
	rsa.GenRsaKey(2048)
}
