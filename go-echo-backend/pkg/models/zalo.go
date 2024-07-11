package models

type ZaloConfig struct {
	Model
	SecretKey    string `json:"secret_key,omitempty"`
	AppID        string `json:"app_id,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	ExpiredAt    int64  `json:"expired_at"`
}
