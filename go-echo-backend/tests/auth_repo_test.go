package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Login(t *testing.T) {
	var app = initApp()

	resp, err := repo.NewAuthRepo(app.DB).LoginEmail(models.LoginEmailForm{
		Email:    "admin@inflow.com",
		Password: "1234qwer",
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}

func TestAuth_ForgotPassword(t *testing.T) {
	var app = initApp("local")

	resp, err := repo.NewAuthRepo(app.DB).ForgotPassword(models.ForgotPasswordForm{
		Email:       "invite+01@gmail.com",
		RedirectURL: "https://dev-brand.joininflow.io/forgot-password",
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}

func TestAuth_Register(t *testing.T) {
	var app = initApp()

	resp, err := repo.NewAuthRepo(app.DB).Register(models.RegisterForm{
		Email: "freelance.01@inflow.com",
		// Password: "1234qwer",
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}
func TestAuth_ResetPassword(t *testing.T) {
	var app = initApp("local")
	resp, err := repo.NewAuthRepo(app.DB).ResetPassword(models.ResetPasswordForm{
		NewPassword:        "@Inflow#4687!",
		TokenResetPassword: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNraHBzNWZza21wNDRhMmNvMHEwIiwidHoiOiIiLCJjaWQiOiIiLCJjdHlwZSI6IiIsImF1ZCI6ImNsaWVudCIsImV4cCI6MTY5NjgzNjY1N30.EVE3uZTxB1kQ3QXVU9-9JcYAw7nhIgm1RJYFOeD-HGQ",
	})
	assert.Nil(t, err)

	helper.PrintJSON(resp)
}
func TestAuth_Admin_ResetPassword(t *testing.T) {
	var app = initApp("local")
	resp, err := repo.NewAuthRepo(app.DB).ResetPassword(models.ResetPasswordForm{
		NewPassword:        "@Inflow#4687!",
		TokenResetPassword: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNraG9xYzdza21wM2R1Zm5rZnJnIiwidHoiOiIiLCJjaWQiOiIiLCJjdHlwZSI6IiIsImF1ZCI6InN0YWZmIiwiZXhwIjoxNjk2ODMyMzIwfQ.0axyBdP3Oj9etvP0D81WvHxv94IubS68rcMzj2ADT9M",
		IsAdminPortal:      true,
	})
	assert.Nil(t, err)

	helper.PrintJSON(resp)
}

func TestAuth_SellerRegister(t *testing.T) {
	var app = initApp("local")

	resp, err := repo.NewAuthRepo(app.DB).Register(models.RegisterForm{
		Email:       "seller.02@inflow.com",
		PhoneNumber: "0232323234",
		IsSeller:    true,
	})

	assert.NoError(t, err)

	helper.PrintJSON(resp)
}

func TestAuth_VerifyEmail(t *testing.T) {
	var app = initApp("dev")

	var err = repo.NewAuthRepo(app.DB).VerifyEmail(models.VerifyEmailForm{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNrdTdhbmlqc3VmcmNtaDlrMjBnIiwidHoiOiIiLCJjaWQiOiIiLCJjdHlwZSI6IiIsImF1ZCI6ImNsaWVudCIsImV4cCI6MTY5ODU0NzQyMiwiaXNzIjoiY2t1N2FuaWpzdWZyY21oOWsyMTAifQ.XCvjpfQ36r_RB-k-HkK0YO9f8L9gpnOCqOoZyQolAgo",
	})

	assert.NoError(t, err)

}
