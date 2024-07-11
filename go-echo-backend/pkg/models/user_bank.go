package models

type UserBank struct {
	Model
	UserID      string `gorm:"not null" json:"user_id"`
	CountryCode string `gorm:"type:varchar(10);default:'VN'" json:"country_code"`
	Currency    string `gorm:"type:varchar(10);default:'VND'" json:"currency"`
	Content     string `gorm:"type:text;default:''" json:"content"`
	IsDisabled  *bool  `gorm:"type:bool;default:false" json:"is_disabled"`
}

type UserBankSlice []UserBank

type UserBanksForm struct {
	JwtClaimsInfo
	ID          string `json:"id,omitempty"`
	CountryCode string `json:"country_code"`
	Currency    string `json:"currency"`
	Content     string `json:"content"`
	IsDisabled  *bool  `json:"is_disabled"`
}
type DeleteUserBanksForm struct {
	JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" validate:"required"`
}

type DeleteUserBanksByCountryCodeForm struct {
	JwtClaimsInfo
	CountryCode string `json:"country_code" param:"country_code" query:"country_code" validate:"required"`
	Currency    string `json:"currency" validate:"required"`
}
