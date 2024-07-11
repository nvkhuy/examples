package stripehelper

import (
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

// StripeClient client
type StripeClient struct {
	*client.API
	config *config.Configuration
}

var instance *StripeClient

// New init
func New(cfg *config.Configuration) *StripeClient {
	var client = &client.API{}
	var key = cfg.StripeSecretKey
	var maxRetries int64 = 3
	var stripeConfig = &stripe.BackendConfig{
		MaxNetworkRetries: stripe.Int64(maxRetries),
		LeveledLogger:     logger.New("stripe").Sugar(),
	}

	stripe.Key = key
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    cfg.GetServerName(config.ServiceBackend),
		URL:     cfg.ServerBaseURL,
		Version: cfg.BuildVersion,
	})

	client.Init(key, &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, stripeConfig),
	})

	instance = &StripeClient{
		API:    client,
		config: cfg,
	}
	return instance
}

func GetInstance() *StripeClient {
	if instance == nil {
		panic("Must be call new first")
	}

	return instance
}
