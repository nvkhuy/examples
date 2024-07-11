package shopify

import (
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/samber/lo"
)

var instance *App

type App struct {
	*goshopify.App
	config *config.Configuration
	logger *logger.Logger
}

func New(config *config.Configuration) *App {
	var app = goshopify.App{
		ApiKey:      config.ShopifyApiKey,
		ApiSecret:   config.ShopifyApiSecret,
		RedirectUrl: fmt.Sprintf("%s/api/v1/oauth/shopify/callback", config.ServerBaseURL),
		Scope:       "read_products,write_products,read_product_listings,read_inventory,write_inventory",
	}

	instance = &App{
		App:    &app,
		config: config,
		logger: logger.New("shopify"),
	}

	return instance
}

func GetInstance() *App {
	return instance
}

func (a *App) GetClient(clientInfo *ClientInfo) *goshopify.Client {
	return a.NewClient(clientInfo.ShopName, clientInfo.Token)
}

func (a *App) CreateWebhooks(clientInfo *ClientInfo) error {
	var client = a.NewClient(clientInfo.ShopName, clientInfo.Token).Webhook
	webhooks, err := client.List(nil)
	if err != nil {
		return err
	}

	var topics = []string{
		"products/create",
		"products/delete",
		"products/update",
		"shop/update",
		"inventory_items/create",
		"inventory_items/delete",
		"inventory_items/update",
	}

	for _, topic := range topics {
		var webhook = goshopify.Webhook{
			Topic:   topic,
			Address: fmt.Sprintf("%s/api/v1/webhook/shopify/%s/%s", a.config.ServerBaseURL, clientInfo.ShopName, topic),
			Format:  "json",
		}

		existingWebhook, found := lo.Find(webhooks, func(i goshopify.Webhook) bool {
			return webhook.Topic == i.Topic
		})

		if found {
			webhook.ID = existingWebhook.ID
			_, err = client.Update(webhook)
			a.logger.Debugf("Updating existing webhook topic=%s err=%+v", existingWebhook.Topic, err)

		} else {
			a.logger.Debugf("Creating a new webhook topic=%s err=%+v", existingWebhook.Topic, err)
			_, err = client.Create(webhook)

		}
	}

	return nil
}

func (a *App) DeleteWebhooks(clientInfo *ClientInfo) error {
	var client = a.NewClient(clientInfo.ShopName, clientInfo.Token).Webhook
	webhooks, err := client.List(nil)
	if err != nil {
		return err
	}

	for _, webhook := range webhooks {
		err = client.Delete(webhook.ID)
		a.logger.Debugf("Deleting existing webhook topic=%s err=%+v", webhook.Topic, err)

	}

	return nil
}
