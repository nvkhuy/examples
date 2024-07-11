package zalo

import (
	"encoding/json"
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type OARefreshTokenParams struct {
	SecretKey    string `json:"secret_key"`
	RefreshToken string `json:"refresh_token"`
	AppID        string `json:"app_id"`
	GrantType    string `json:"grant_type"`
}

type OARefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

// OARefreshToken https://developers.zalo.me/docs/api/official-account-api/phu-luc/official-account-access-token-post-4307
func OARefreshToken(db *db.DB, params OARefreshTokenParams) (resp OARefreshTokenResponse, err error) {
	var config models.ZaloConfig
	config, err = GetZaloConfig(db)
	if err != nil {
		return
	}

	params.SecretKey = config.SecretKey
	params.RefreshToken = config.RefreshToken
	params.AppID = config.AppID
	params.GrantType = "refresh_token"

	// Define the request URL
	requestURL := "https://oauth.zaloapp.com/v4/oa/access_token"

	// Define the request headers
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"secret_key":   params.SecretKey,
	}

	// Define the request data
	data := url.Values{}
	data.Set("refresh_token", params.RefreshToken)
	data.Set("app_id", params.AppID)
	data.Set("grant_type", params.GrantType)

	// Create the request URL with parameters
	u, err := url.Parse(requestURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	u.RawQuery = data.Encode()

	// Create the HTTP request
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(""))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Perform the HTTP request
	client := &http.Client{}
	requestResp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer requestResp.Body.Close()

	// Parse the response body into OARefreshTokenResponse struct
	err = json.NewDecoder(requestResp.Body).Decode(&resp)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	err = UpdateToken(db, resp.RefreshToken, resp.AccessToken, resp.ExpiresIn)
	return
}

func GetZaloConfig(db *db.DB) (config models.ZaloConfig, err error) {
	err = db.Model(&models.ZaloConfig{}).Where("id = ?", enums.ZaloConfigKeyRefreshToken).First(&config).Error
	return
}

func UpdateToken(db *db.DB, refreshToken, accessToken, expiresIn string) (err error) {
	expireInt, err := strconv.ParseInt(expiresIn, 10, 64)
	if err != nil {
		return
	}

	var expiredAt = time.Now().Add(time.Duration(expireInt) * time.Second).Unix()

	var updates = models.ZaloConfig{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expireInt,
		ExpiredAt:    expiredAt,
	}
	err = db.Model(&models.ZaloConfig{}).Where("id = ?", enums.ZaloConfigKeyRefreshToken).Updates(&updates).Error
	return
}

// SaveZaloConfig https://developers.zalo.me/tools/explorer
func SaveZaloConfig(db *db.DB, values OARefreshTokenResponse) (err error) {
	expireInt, err := strconv.ParseInt(values.ExpiresIn, 10, 64)
	if err != nil {
		return
	}

	var expiredAt = time.Now().Add(time.Duration(expireInt) * time.Second).Unix()

	var updates = models.ZaloConfig{
		Model: models.Model{
			ID: string(enums.ZaloConfigKeyRefreshToken),
		},
		AppID:        "4212732216110924634",
		SecretKey:    "Gl9J9VWHK8Zk277eZHTJ",
		AccessToken:  values.AccessToken,
		RefreshToken: values.RefreshToken,
		ExpiresIn:    expireInt,
		ExpiredAt:    expiredAt,
	}
	err = db.Model(&models.ZaloConfig{}).Where("id = ?", enums.ZaloConfigKeyRefreshToken).Create(&updates).Error
	return
}

// cron job -> each 5 minute -> if expire_in <= 300 -> refresh token
