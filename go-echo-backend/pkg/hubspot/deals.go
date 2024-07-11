package hubspot

import (
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/google/go-querystring/query"
)

// https://app.hubspot.com/property-settings/24143466/properties?type=0-3&search=deal%20stage&action=edit&property=dealstage
type DealStage string

var (
	DealStagePending                      DealStage = "109551430"
	DealStageSample                       DealStage = "appointmentscheduled"
	DealStageFashionDesigningPPS          DealStage = "145245518"
	DealStageTeckpackCreatingDesigningPPS DealStage = "145245523"
	DealStageClosedWonSample              DealStage = "150684720"
	DealStageClosedLostProspecting        DealStage = "150684721"
)

// https://app.hubspot.com/property-settings/24143466/properties?type=0-3&search=pi&action=edit&property=pipeline
type Pipeline string

var (
	PipelineManufacturing Pipeline = "default"
	PipelinePreProduction Pipeline = "appointmentscheduled"
	PipelineProspecting   Pipeline = "79482640"
)

type GetDealsParams struct {
	Properties []string   `json:"properties"`
	Inputs     []*InputID `json:"inputs"`
}

type Deal struct {
	ID         string          `json:"id,omitempty"`
	Properties *DealProperties `json:"properties,omitempty"`
	CreatedAt  time.Time       `json:"createdAt,omitempty"`
	UpdatedAt  time.Time       `json:"updatedAt,omitempty"`
	Archived   bool            `json:"archived,omitempty"`
}

type DealPropertiesForm struct {
	Amount           string     `json:"amount,omitempty" mapstructure:"amount,omitempty"`
	Closedate        *time.Time `json:"closedate,omitempty" mapstructure:"closedate,omitempty"`
	Dealname         string     `json:"dealname,omitempty" mapstructure:"dealname,omitempty"`
	Pipeline         Pipeline   `json:"pipeline,omitempty" mapstructure:"pipeline,omitempty"`
	Dealstage        DealStage  `json:"dealstage,omitempty" mapstructure:"dealstage,omitempty"`
	HubspotOwnerID   string     `json:"hubspot_owner_id,omitempty" mapstructure:"hubspot_owner_id,omitempty"`
	HubspotContactID string     `json:"hubspot_contact_id,omitempty" mapstructure:"-"`
}

type DealProperties struct {
	Amount                          string     `json:"amount,omitempty"`
	AmountInHomeCurrency            string     `json:"amount_in_home_currency,omitempty"`
	Createdate                      *time.Time `json:"createdate,omitempty"`
	DaysToClose                     string     `json:"days_to_close,omitempty"`
	Dealname                        string     `json:"dealname,omitempty"`
	Dealstage                       string     `json:"dealstage,omitempty"`
	HsAllOwnerIds                   string     `json:"hs_all_owner_ids,omitempty"`
	HsClosedAmount                  string     `json:"hs_closed_amount,omitempty"`
	HsClosedAmountInHomeCurrency    string     `json:"hs_closed_amount_in_home_currency,omitempty"`
	HsCreatedate                    *time.Time `json:"hs_createdate,omitempty"`
	HsDaysToCloseRaw                string     `json:"hs_days_to_close_raw,omitempty"`
	HsDealStageProbabilityShadow    string     `json:"hs_deal_stage_probability_shadow,omitempty"`
	HsForecastAmount                string     `json:"hs_forecast_amount,omitempty"`
	HsIsActiveSharedDeal            string     `json:"hs_is_active_shared_deal,omitempty"`
	HsIsClosed                      string     `json:"hs_is_closed,omitempty"`
	HsIsClosedWon                   string     `json:"hs_is_closed_won,omitempty"`
	HsIsDealSplit                   string     `json:"hs_is_deal_split,omitempty"`
	HsIsOpenCount                   string     `json:"hs_is_open_count,omitempty"`
	HsLastmodifieddate              *time.Time `json:"hs_lastmodifieddate,omitempty"`
	HsObjectID                      string     `json:"hs_object_id,omitempty"`
	HsObjectSource                  string     `json:"hs_object_source,omitempty"`
	HsObjectSourceID                string     `json:"hs_object_source_id,omitempty"`
	HsObjectSourceLabel             string     `json:"hs_object_source_label,omitempty"`
	HsProjectedAmount               string     `json:"hs_projected_amount,omitempty"`
	HsProjectedAmountInHomeCurrency string     `json:"hs_projected_amount_in_home_currency,omitempty"`
	HsUserIdsOfAllOwners            string     `json:"hs_user_ids_of_all_owners,omitempty"`
	HubspotOwnerAssigneddate        *time.Time `json:"hubspot_owner_assigneddate,omitempty"`
	HubspotOwnerID                  string     `json:"hubspot_owner_id,omitempty"`
	Pipeline                        string     `json:"pipeline,omitempty"`
}

func (c *Client) GetDeals(params GetDealsParams) (*Pagination[Deal], error) {
	v, _ := query.Values(&params)

	var pagination Pagination[Deal]
	var apiErr ApiError

	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetQueryParamsFromValues(v).
		SetResult(&pagination).
		SetError(&apiErr).
		Get("https://api.hubapi.com/crm/v3/objects/deals")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}

	return &pagination, err
}

func (c *Client) CreateDeal(params *DealPropertiesForm) (*Deal, error) {
	var payload = map[string]interface{}{
		"properties": helper.StructToMap(params, "mapstructure"),
		"associations": []map[string]interface{}{
			{
				"to": map[string]interface{}{
					"id": params.HubspotContactID,
				},
				"types": []map[string]interface{}{
					{
						"associationCategory": "HUBSPOT_DEFINED",
						"associationTypeId":   3,
					},
				},
			},
		},
	}

	var resp Deal
	var apiErr ApiError
	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetBody(&payload).
		SetError(&apiErr).
		SetResult(&resp).
		Post("https://api.hubapi.com/crm/v3/objects/deals")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}
	return &resp, nil
}

func (c *Client) UpdateDeal(dealId string, params *DealPropertiesForm) (*Deal, error) {
	var payload = map[string]interface{}{
		"properties": params,
	}

	var resp Deal
	var apiErr ApiError
	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetBody(&payload).
		SetError(&apiErr).
		SetResult(&resp).
		SetPathParam("dealId", dealId).
		Patch("https://api.hubapi.com/crm/v3/objects/deals/{dealId}")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}
	return &resp, nil
}

func (c *Client) GetDeal(dealId string) (*Deal, error) {
	var resp Deal
	var apiErr ApiError
	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetError(&apiErr).
		SetResult(&resp).
		SetPathParam("dealId", dealId).
		Get("https://api.hubapi.com/crm/v3/objects/deals/{dealId}")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}
	return &resp, nil
}

type AssignDealToContactResponse struct {
	ID           string             `json:"id"`
	Properties   *ContactProperties `json:"properties"`
	CreatedAt    *time.Time         `json:"createdAt"`
	UpdatedAt    *time.Time         `json:"updatedAt"`
	Archived     bool               `json:"archived"`
	Associations struct {
		Contacts struct {
			Results []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"results"`
		} `json:"contacts"`
	} `json:"associations"`
}

func (c *Client) AssignDealToContact(dealId string, contactId string) (*AssignDealToContactResponse, error) {
	var resp AssignDealToContactResponse
	var apiErr ApiError
	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetError(&apiErr).
		SetResult(&resp).
		SetPathParams(map[string]string{
			"dealId":            dealId,
			"toObjectId":        contactId,
			"toObjectType":      "contacts",
			"associationTypeId": "3", //deal_to_contact
		}).
		Put("https://api.hubapi.com/crm/v3/objects/deals/{dealId}/associations/{toObjectType}/{toObjectId}/{associationTypeId}")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}
	return &resp, nil
}
