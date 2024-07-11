package customerio

import (
	"github.com/customerio/go-customerio/v3"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/go-resty/resty/v2"
)

var instance *Client

type Client struct {
	Track         *customerio.CustomerIO
	Transactional *customerio.APIClient
	restyClient   *resty.Client
	config        *config.Configuration
}

func New(config *config.Configuration) *Client {
	var track = customerio.NewTrackClient(config.CustomerIOSiteID, config.CustomerIOApiTrackingKey, customerio.WithRegion(customerio.RegionUS))
	var transactional = customerio.NewAPIClient(config.CustomerIOApiAppKey)
	instance = &Client{
		Track:         track,
		Transactional: transactional,
		restyClient:   resty.New().SetDebug(true),
		config:        config,
	}
	return instance
}

func GetInstance() *Client {
	if instance == nil {
		panic("Must call New() first")
	}

	return instance
}
