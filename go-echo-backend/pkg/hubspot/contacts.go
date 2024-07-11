package hubspot

import (
	"time"

	"github.com/google/go-querystring/query"
	"github.com/samber/lo"
)

type InputID struct {
	ID string `json:"id"`
}
type GetContactsParams struct {
	Properties []string   `json:"properties"`
	Inputs     []*InputID `json:"inputs"`
}

type GetContactParams struct {
	ContactID string `json:"contactId" validate:"required"`
}

type Contact struct {
	ID         string             `json:"id,omitempty"`
	Properties *ContactProperties `json:"properties,omitempty"`
	CreatedAt  time.Time          `json:"createdAt,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty"`
	Archived   bool               `json:"archived,omitempty"`
}

type ContactPropertiesForm struct {
	Email          string `json:"email,omitempty"`
	Firstname      string `json:"firstname,omitempty"`
	Lastname       string `json:"lastname,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Company        string `json:"company,omitempty"`
	Website        string `json:"website,omitempty"`
	Lifecyclestage string `json:"lifecyclestage,omitempty"`
}

type ContactProperties struct {
	Createdate       *time.Time `json:"createdate,omitempty"`
	Email            string     `json:"email,omitempty"`
	Firstname        string     `json:"firstname,omitempty"`
	HsObjectID       string     `json:"hs_object_id,omitempty"`
	Lastmodifieddate *time.Time `json:"lastmodifieddate,omitempty"`
	Lastname         string     `json:"lastname,omitempty"`
	Phone            string     `json:"phone,omitempty"`
}

func (c *Client) GetContacts(params GetContactsParams) (*Pagination[Contact], error) {
	v, _ := query.Values(&params)

	var pagination Pagination[Contact]
	var apiErr ApiError

	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetQueryParamsFromValues(v).
		SetResult(&pagination).
		SetError(&apiErr).
		Get("https://api.hubapi.com/crm/v3/objects/contacts")

	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}

	return &pagination, err
}

func (c *Client) GetContactByID(contactId string) (*Contact, error) {
	var resp Contact
	var apiErr ApiError

	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetPathParam("contactId", contactId).
		SetResult(&resp).
		SetError(&apiErr).
		Get("https://api.hubapi.com/crm/v3/objects/contacts/{contactId}")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}

	return &resp, err
}

func (c *Client) SearchContactsByEmail(emails []string) (*Results[Contact], error) {
	var resp Results[Contact]
	var apiErr ApiError

	var payload = map[string]interface{}{
		"idProperty": "email",
		"inputs": lo.Map(emails, func(email string, index int) map[string]interface{} {
			return map[string]interface{}{
				"id": email,
			}
		}),
	}
	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetBody(&payload).
		SetResult(&resp).
		SetError(&apiErr).
		Post("https://api.hubapi.com/crm/v3/objects/contacts/batch/read")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}

	return &resp, err
}

func (c *Client) CreateContact(params *ContactPropertiesForm) (*Contact, error) {
	var payload = map[string]interface{}{
		"properties": params,
	}

	var resp Contact
	var apiErr ApiError

	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetBody(&payload).
		SetResult(&resp).
		SetError(&apiErr).
		Post("https://api.hubapi.com/crm/v3/objects/contacts")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}
	return &resp, nil
}

func (c *Client) UpdateContact(params *ContactPropertiesForm) (*Contact, error) {
	var payload = map[string]interface{}{
		"properties": params,
	}

	var resp Contact
	var apiErr ApiError

	result, err := c.client.R().
		SetAuthToken(c.config.HubspotAccessToken).
		SetBody(&payload).
		SetError(&apiErr).
		SetResult(&resp).
		Post("https://api.hubapi.com/crm/v3/objects/contacts")
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, &apiErr
	}
	return &resp, nil
}
