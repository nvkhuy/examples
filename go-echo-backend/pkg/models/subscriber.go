package models

// User user's model
type Subscriber struct {
	Model

	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Email  string `gorm:"unique;type:citext;default:null" json:"email,omitempty"`
}

type SubscribeByEmailForm struct {
	JwtClaimsInfo

	Email string `json:"email"`
}
