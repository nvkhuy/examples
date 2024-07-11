package oauth

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConfig *oauth2.Config

func New(config *config.Configuration) *oauth2.Config {
	oauthConfig = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/api/v1/oauth/google/callback", config.ServerBaseURL),
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return oauthConfig
}

func GetInstance() *oauth2.Config {
	if oauthConfig == nil {
		panic("Please init oauth package first")
	}
	return oauthConfig
}
