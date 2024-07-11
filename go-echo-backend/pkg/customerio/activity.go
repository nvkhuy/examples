package customerio

import "github.com/google/go-querystring/query"

type ActivityResouces struct {
	Activities []*Activity `json:"activities,omitempty"`
	Next       string      `json:"next,omitempty"`
}

type Activity struct {
	CustomerID          string                 `json:"customer_id,omitempty"`
	CustomerIdentifiers *Identifiers           `json:"customer_identifiers,omitempty"`
	Data                map[string]interface{} `json:"data,omitempty"`
	DeliveryID          string                 `json:"delivery_id,omitempty"`
	DeliveryType        string                 `json:"delivery_type,omitempty"`
	ID                  string                 `json:"id,omitempty"`
	Timestamp           int                    `json:"timestamp,omitempty"`
	Type                string                 `json:"type,omitempty"`
	Url                 string                 `json:"url,omitempty"`
	Name                string                 `json:"name,omitempty"`
}

type GetActivitiesParams struct {
	Type  string `json:"type" query:"type" url:"type"`
	Start string `json:"start" query:"start" url:"start"`
	Name  string `json:"name" query:"name" url:"name"`
	Limit int    `json:"limit" query:"limit" url:"limit"`
}

func (client *Client) GetActivities(id string, params GetActivitiesParams) (*ActivityResouces, error) {
	values, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	var apiErr APIErrors
	var resp ActivityResouces
	_, err = client.restyClient.R().
		SetAuthToken(client.config.CustomerIOApiAppKey).
		SetQueryParamsFromValues(values).
		SetQueryParams(map[string]string{
			"customer_id": id,
			"id_type":     "id",
		}).
		SetResult(&resp).
		SetError(&apiErr).
		Get("https://api.customer.io/v1/activities")

	if len(apiErr.Errors) > 0 {
		return &resp, &apiErr
	}

	return &resp, err
}
