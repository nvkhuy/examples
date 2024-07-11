package hubspot

import (
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/validation"
	"github.com/go-resty/resty/v2"
)

// https://developers.hubspot.com/docs/api/overview

type Client struct {
	config    *config.Configuration
	client    *resty.Client
	validator *validation.Validator
}

func New(config *config.Configuration) *Client {
	return &Client{
		config:    config,
		client:    resty.New().SetDebug(true),
		validator: validation.RegisterValidation(),
	}
}
