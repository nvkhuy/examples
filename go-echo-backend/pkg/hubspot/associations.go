package hubspot

import "github.com/google/go-querystring/query"

type Association struct {
}

type GetAssociationsParams struct {
}

func (c *Client) GetAssociations(params GetAssociationsParams) (*Pagination[Association], error) {
	v, _ := query.Values(&params)

	var pagination Pagination[Association]
	var apiErr ApiError

	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetQueryParamsFromValues(v).
		SetResult(&pagination).
		SetError(&apiErr).
		Get("https://api.hubapi.com/crm/v3/associations/deals/contacts/types")

	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}

	return &pagination, err
}
