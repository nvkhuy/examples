package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/golang-jwt/jwt"
	"github.com/rotisserie/eris"
	"github.com/rs/xid"
	"github.com/thaitanloi365/go-utils/values"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (user *User) BeforeCreate(tx *gorm.DB) error {
	var err = user.HashPassword()
	if err != nil {
		return eris.Wrap(err, "")
	}

	return nil

}

func (user *User) BeforeSave(tx *gorm.DB) error {
	var name = user.GetFullName()
	if name != "" {
		tx.Statement.SetColumn("name", name)
	}
	// if user.CustomerType == "" {
	// 	tx.Statement.SetColumn("customer_type", enums.CustomerTypeBuyer)
	// }
	return nil
}

func (user *User) IsTestAccount() bool {
	return false
}

func (user *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return eris.Wrap(err, "")
	}
	user.Password = string(hash)
	return nil
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (user *User) GetFullName() string {
	if user.Name != "" {
		return user.Name
	}
	if user.FirstName == "" && user.LastName != "" {
		return user.LastName
	} else if user.LastName == "" && user.FirstName != "" {
		return user.FirstName
	} else if user.LastName != "" && user.FirstName != "" {
		return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	return user.Name
}

func (user *User) GenerateToken(secret string, duration ...time.Duration) (string, error) {
	claims := &JwtClaims{}
	claims.ID = user.ID
	claims.ExpiresAt = 0
	claims.Audience = string(user.Role)
	claims.Subject = user.Role.String()
	if user.Team != "" {
		claims.Subject += fmt.Sprintf(":%s", user.Team)
	}
	claims.Subject = strings.ToLower(claims.Subject)

	if user.TokenIssuer == "" {
		user.TokenIssuer = xid.New().String()
	}

	claims.Issuer = user.TokenIssuer

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var key = []byte(secret)
	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (user *User) GenerateResetPasswordTokenAndUpdate(db *db.DB) (string, error) {
	token, err := user.GenerateResetPasswordToken(db.Configuration.JWTResetPasswordSecret, db.Configuration.JWTResetPasswordExpiry)
	if err != nil {
		return "", err
	}

	var updates = User{
		TokenResetPassword:       user.TokenResetPassword,
		TokenResetPasswordSentAt: user.TokenResetPasswordSentAt,
	}
	err = db.Model(&User{}).Where("id = ?", user.ID).Updates(updates).Error
	if err != nil {
		if db.IsRecordNotFoundError(err) {
			return "", errs.ErrUserNotFound
		}
		return "", err
	}

	return token, err
}

func (user *User) UpdateAccountStatus(db *db.DB) (err error) {
	var updates = User{
		AccountStatus: user.AccountStatus,
	}
	err = db.Model(&User{}).Where("id = ?", user.ID).Updates(updates).Error
	if err != nil && db.IsRecordNotFoundError(err) {
		err = errs.ErrUserNotFound
	}
	return
}

func (user *User) GenerateResetPasswordToken(secret string, expiry time.Duration) (string, error) {
	var jwtClaims JwtClaims

	if user.TokenResetPassword != nil && *user.TokenResetPassword != "" {
		var err = jwtClaims.ValidateToken(secret, *user.TokenResetPassword)
		if err == nil {
			return *user.TokenResetPassword, nil
		}
	}

	jwtClaims.ID = user.ID
	jwtClaims.StandardClaims.Audience = string(user.Role)
	jwtClaims.StandardClaims.Issuer = user.TokenIssuer
	jwtClaims.StandardClaims.ExpiresAt = time.Now().Add(expiry).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	key := []byte(secret)

	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	user.TokenResetPassword = &t
	user.TokenResetPasswordSentAt = values.Int64(time.Now().Unix())

	return t, nil
}

func (user *User) GenerateVerificationToken(db *db.DB) (string, error) {
	var jwtClaims JwtClaims

	if user.TokenIssuer == "" {
		user.TokenIssuer = xid.New().String()
	}

	jwtClaims.ID = user.ID
	jwtClaims.StandardClaims.Audience = string(user.Role)
	jwtClaims.StandardClaims.Issuer = user.TokenIssuer
	jwtClaims.StandardClaims.ExpiresAt = time.Now().Add(db.Configuration.JWTEmailVerificationExpiry).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	key := []byte(db.Configuration.JWTEmailVerificationSecret)

	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	var updates = User{VerificationToken: &t, TokenIssuer: user.TokenIssuer, VerificationSentAt: values.Int64(time.Now().Unix())}

	err = db.Model(&User{}).Where("id = ?", user.ID).Updates(updates).Error
	if err != nil {
		if db.IsRecordNotFoundError(err) {
			return "", errs.ErrUserNotFound
		}
		return "", err
	}

	return t, nil
}

// GetStripeCustomerMetadata get medata for stripe customer
func (user *User) GetStripeCustomerMetadata() map[string]string {
	return map[string]string{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}
}

// CreateStripeCustomer create stripe customer
func (user *User) createStripeCustomer(db *db.DB, paymentMethod *string) bool {
	var isChanged = false

	if user.StripeCustomerID == "" {
		var p = &stripehelper.CreateCustomerParams{
			PaymentMethodID: paymentMethod,
			Email:           user.Email,
			Name:            user.Name,
			Metadata:        user.GetStripeCustomerMetadata(),
		}

		id, err := stripehelper.GetInstance().CreateCustomer(p)
		if err == nil {
			user.StripeCustomerID = id
			isChanged = true
		}

	}

	return isChanged
}

// CreateStripeCustomer create stripe customer
func (user *User) CreateStripeCustomer(db *db.DB, paymentMethod *string, updateLater ...bool) error {
	if user.ID != "" {
		var isChanged = user.createStripeCustomer(db, paymentMethod)
		var isSkip = len(updateLater) > 0 && updateLater[0]
		if isChanged && !isSkip && user.StripeCustomerID != "" {
			var err = db.Model(user).Update("stripe_customer_id", user.StripeCustomerID).Error
			return err
		}

	}
	return nil
}

func (user *User) GetCustomerIOMetadata(extras map[string]interface{}) map[string]interface{} {
	var result = map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	if user.Role != "" {
		result["role"] = user.Role
		result["role_display"] = user.Role.DisplayName()
	}

	if user.Team != "" {
		result["team"] = user.Team
		result["team_display"] = user.Team.DisplayName()
	}

	if user.Timezone != "" {
		result["timezone"] = user.Timezone
	}

	if user.AccountStatus != "" {
		result["account_status"] = user.AccountStatus
	}

	if user.CountryCode != "" {
		result["country_code"] = user.CountryCode
	}

	if user.StripeCustomerID != "" {
		result["stripe_customer_id"] = user.StripeCustomerID
	}

	if user.Avatar != nil {
		result["avatar"] = user.Avatar.GenerateFileURL()

	}
	if user.LastLogin != nil {
		result["login_at"] = *user.LastLogin

	}
	if user.HubspotContactID != "" {
		result["hubspot_contact_id"] = user.HubspotContactID
	}

	if user.HubspotOwnerID != "" {
		result["hubspot_owner_id"] = user.HubspotOwnerID
	}

	for k, v := range extras {
		result[k] = v
	}

	return result

}

func (users Users) GetCustomerIOMetadata(extras map[string]interface{}) (list []map[string]interface{}) {
	for _, user := range users {
		list = append(list, user.GetCustomerIOMetadata(extras))
	}
	return

}
