package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/stretchr/testify/assert"
)

func TestHubspot_GetContacts(t *testing.T) {
	var cfg = initConfig()
	result, err := hubspot.New(cfg).GetContacts(hubspot.GetContactsParams{})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestHubspot_SearchContacts(t *testing.T) {
	var cfg = initConfig()
	result, err := hubspot.New(cfg).SearchContactsByEmail([]string{"loithai+1@joininflow.io"})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestHubspot_CreateContact(t *testing.T) {
	var cfg = initConfig()
	result, err := hubspot.New(cfg).CreateContact(&hubspot.ContactPropertiesForm{
		Email:          "loithai+1@joininflow.io",
		Firstname:      "Loi",
		Lastname:       "Thai",
		Phone:          "+84327308788",
		Company:        "Inflow",
		Lifecyclestage: "lead",
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestHubspot_CreateDeal(t *testing.T) {
	var cfg = initConfig()
	result, err := hubspot.New(cfg).CreateDeal(&hubspot.DealPropertiesForm{
		Amount:           "100",
		Dealname:         "Loi test deal 2",
		Pipeline:         hubspot.PipelineManufacturing,
		Dealstage:        hubspot.DealStagePending,
		HubspotOwnerID:   "345074727",
		HubspotContactID: "1313351",
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestHubspot_AssignDealToContact(t *testing.T) {
	var cfg = initConfig()
	result, err := hubspot.New(cfg).AssignDealToContact("17314677910", "1313351")
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestHubspot_GetDeal(t *testing.T) {
	var cfg = initConfig()
	result, err := hubspot.New(cfg).GetDeal("17314677910")
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestHubspot_GetOwners(t *testing.T) {
	var cfg = initConfig()
	result, err := hubspot.New(cfg).GetOwners(hubspot.GetOwnersParams{
		Email: "quockhanh@joininflow.io",
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestHubspot_GetAssociations(t *testing.T) {
	var cfg = initConfig()
	result, err := hubspot.New(cfg).GetAssociations(hubspot.GetAssociationsParams{})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}
