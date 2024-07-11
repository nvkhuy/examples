package hubspot

import (
	"time"

	"github.com/google/go-querystring/query"
)

type GetOwnersParams struct {
	Email string `json:"email,omitempty" url:"email,omitempty"`
	After string `json:"after,omitempty" url:"after,omitempty"`
	Limit int    `json:"limit,omitempty" url:"limit,omitempty"`
}

type Owner struct {
	ID         string           `json:"id,omitempty"`
	Properties *OwnerProperties `json:"properties,omitempty"`
	CreatedAt  time.Time        `json:"createdAt,omitempty"`
	UpdatedAt  time.Time        `json:"updatedAt,omitempty"`
	Archived   bool             `json:"archived,omitempty"`
	Firstname  string           `json:"firstname,omitempty"`
	Lastname   string           `json:"lastname,omitempty"`
	UserId     int              `json:"userId"`
}

type OwnerPropertiesForm struct {
	Createdate       *time.Time `json:"createdate,omitempty"`
	Email            string     `json:"email,omitempty"`
	Firstname        string     `json:"firstname,omitempty"`
	HsObjectID       string     `json:"hs_object_id,omitempty"`
	Lastmodifieddate *time.Time `json:"lastmodifieddate,omitempty"`
	Lastname         string     `json:"lastname,omitempty"`

	Phone          string `json:"phone,omitempty"`
	Company        string `json:"company,omitempty"`
	Website        string `json:"website,omitempty"`
	Lifecyclestage string `json:"lifecyclestage,omitempty"`
}

type OwnerProperties struct {
	Amount                          string    `json:"amount"`
	AmountInHomeCurrency            string    `json:"amount_in_home_currency"`
	Createdate                      time.Time `json:"createdate"`
	DaysToClose                     string    `json:"days_to_close"`
	Dealname                        string    `json:"dealname"`
	Dealstage                       string    `json:"dealstage"`
	HsAllOwnerIds                   string    `json:"hs_all_owner_ids"`
	HsClosedAmount                  string    `json:"hs_closed_amount"`
	HsClosedAmountInHomeCurrency    string    `json:"hs_closed_amount_in_home_currency"`
	HsCreatedate                    time.Time `json:"hs_createdate"`
	HsDaysToCloseRaw                string    `json:"hs_days_to_close_raw"`
	HsDealStageProbabilityShadow    string    `json:"hs_deal_stage_probability_shadow"`
	HsForecastAmount                string    `json:"hs_forecast_amount"`
	HsIsActiveSharedDeal            string    `json:"hs_is_active_shared_deal"`
	HsIsClosed                      string    `json:"hs_is_closed"`
	HsIsClosedWon                   string    `json:"hs_is_closed_won"`
	HsIsDealSplit                   string    `json:"hs_is_deal_split"`
	HsIsOpenCount                   string    `json:"hs_is_open_count"`
	HsLastmodifieddate              time.Time `json:"hs_lastmodifieddate"`
	HsObjectID                      string    `json:"hs_object_id"`
	HsObjectSource                  string    `json:"hs_object_source"`
	HsObjectSourceID                string    `json:"hs_object_source_id"`
	HsObjectSourceLabel             string    `json:"hs_object_source_label"`
	HsProjectedAmount               string    `json:"hs_projected_amount"`
	HsProjectedAmountInHomeCurrency string    `json:"hs_projected_amount_in_home_currency"`
	HsUserIdsOfAllOwners            string    `json:"hs_user_ids_of_all_owners"`
	HubspotOwnerAssigneddate        time.Time `json:"hubspot_owner_assigneddate"`
	HubspotOwnerID                  string    `json:"hubspot_owner_id"`
	Pipeline                        string    `json:"pipeline"`
}

func (c *Client) GetOwners(params GetOwnersParams) (*Pagination[Owner], error) {
	v, _ := query.Values(&params)

	var pagination Pagination[Owner]
	var apiErr ApiError

	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetQueryParamsFromValues(v).
		SetResult(&pagination).
		Get("https://api.hubapi.com/crm/v3/owners")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}

	return &pagination, err
}

func (c *Client) CreateOwner(params OwnerPropertiesForm) (*Owner, error) {
	var payload = map[string]interface{}{
		"properties": params,
	}

	var apiErr ApiError
	var owner Owner
	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetBody(&payload).
		SetResult(&owner).
		SetError(&apiErr).
		Post("https://api.hubapi.com/crm/v3/owners")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}
	return nil, nil
}
