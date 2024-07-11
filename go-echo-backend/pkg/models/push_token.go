package models

type PushTokens []*PushToken

type PushToken struct {
	Token    string `gorm:"primaryKey" json:"token"`
	UserID   string `json:"user_id"`
	Platform string `json:"platform"`
	LastUsed int64  `json:"last_used"`
}

type PushTokenCreateForm struct {
	Token    string `json:"token" validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=ios android web"`
}

type PushTokenDeleteForm struct {
	Token string `json:"token"`
}
