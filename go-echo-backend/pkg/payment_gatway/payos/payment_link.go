package payos

import "fmt"

func (c *Client) CreatePaymentLink(req *CreatePaymentLinkRequest) (*APIResponse[CreatePaymentLinkResponse], error) {
	var data = fmt.Sprintf("amount=%d&cancelUrl=%s&description=%s&orderCode=%d&returnUrl=%s", req.Amount, req.CancelURL, req.Description, req.OrderCode, req.ReturnURL)
	c.generateHMAC(data)

	var resp APIResponse[CreatePaymentLinkResponse]

	_, err := c.client.R().
		SetHeader("x-api-key", c.cfg.PayosApiKey).
		SetHeader("x-client-id", c.cfg.PayosClientID).
		SetBody(req).
		SetResult(&resp).
		Post("/v2/payment-requests")

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetPaymentLinkInfo(id string) (*APIResponse[GetPaymentLinkResponse], error) {
	var resp APIResponse[GetPaymentLinkResponse]

	_, err := c.client.R().
		SetHeader("x-api-key", c.cfg.PayosApiKey).
		SetHeader("x-client-id", c.cfg.PayosClientID).
		SetResult(&resp).
		SetPathParam("id", id).
		Get("/v2/payment-requests/{id}")

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CancelPaymentLinkInfo(id string, cancellationReason string) (*APIResponse[GetPaymentLinkResponse], error) {
	var resp APIResponse[GetPaymentLinkResponse]

	_, err := c.client.R().
		SetHeader("x-api-key", c.cfg.PayosApiKey).
		SetHeader("x-client-id", c.cfg.PayosClientID).
		SetBody(map[string]interface{}{
			"cancellationReason": cancellationReason,
		}).
		SetResult(&resp).
		SetPathParam("id", id).
		Get("/v2/payment-requests/{id}/cancel")

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
