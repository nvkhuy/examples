package zalo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"net/http"
	"net/url"
)

type SendZNSParams struct {
	Phone        string      `json:"phone"`
	Mode         string      `json:"mode"`
	TemplateID   string      `json:"template_id"`
	TemplateData interface{} `json:"template_data"`
	TrackingID   string      `json:"tracking_id"`
}

type SendZNSResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Data    struct {
		SentTime    string `json:"sent_time"`
		SendingMode string `json:"sending_mode"`
		Quota       struct {
			RemainingQuota string `json:"remainingQuota"`
			DailyQuota     string `json:"dailyQuota"`
		} `json:"quota"`
		MsgID string `json:"msg_id"`
	} `json:"data"`
}

func SendZNS(db *db.DB, params SendZNSParams) (resp SendZNSResponse, err error) {
	// Define the request URL
	requestURL := "https://business.openapi.zalo.me/message/template"

	var config models.ZaloConfig
	config, err = GetZaloConfig(db)
	if err != nil {
		return
	}

	// Define the request headers
	headers := map[string]string{
		"Content-Type": "application/json",
		"access_token": config.AccessToken,
	}

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create the request URL with parameters
	u, err := url.Parse(requestURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Perform the HTTP request
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonBody))
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
	return
}
