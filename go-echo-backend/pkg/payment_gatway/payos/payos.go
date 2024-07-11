package payos

import (
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	cfg    *config.Configuration
	client *resty.Client
}

func New(cfg *config.Configuration) *Client {
	return &Client{
		cfg:    cfg,
		client: resty.New().SetBaseURL("https://api-merchant.payos.vn"),
	}
}

func (c *Client) ConfirmWebhook(webhookUrl string) (*APIResponse[any], error) {
	var resp APIResponse[any]

	_, err := c.client.R().
		SetHeader("x-api-key", c.cfg.PayosApiKey).
		SetHeader("x-client-id", c.cfg.PayosClientID).
		SetBody(map[string]interface{}{
			"webhookUrl": webhookUrl,
		}).
		SetResult(&resp).
		Post("/confirm-webhook")

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
