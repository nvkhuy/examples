package models

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
)

func (form UpdatePasswordForm) Validate(db *db.DB) (User, error) {
	var user User
	var claims JwtClaims

	var err = claims.ValidateToken(db.Configuration.JWTSecret, form.TokenResetPassword)
	if err != nil {
		return user, errs.ErrTokenInvalid
	}

	err = db.Select("TokenResetPassword").First(&user, "id = ?", claims.ID).Error
	if err != nil {
		return user, err
	}

	if user.TokenResetPassword == nil || *user.TokenResetPassword != form.TokenResetPassword {
		return user, errs.ErrTokenInvalid
	}

	return user, nil
}

func (l *ForgotPasswordResponse) TransformMessage() *ForgotPasswordResponse {
	l.Message = fmt.Sprintf("Please wait in %d seconds", l.NextInSeconds)
	return l
}
func (l *ResendVerificationEmailResponse) TransformMessage() *ResendVerificationEmailResponse {
	l.Message = fmt.Sprintf("Please wait in %d seconds", l.NextInSeconds)
	return l
}
