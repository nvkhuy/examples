package app

import (
	"github.com/casbin/casbin/v2"
	"github.com/engineeringinflow/inflow-backend/pkg/caching"
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/geo"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/engineeringinflow/inflow-backend/pkg/mailer"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"google.golang.org/api/sheets/v4"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
)

type App struct {
	DB               *db.DB
	AnalyticDB       *db.DB
	Config           *config.Configuration
	Cache            *caching.Client
	Mailer           *mailer.Client
	S3Client         *s3.Client
	StripeHelper     *stripehelper.StripeClient
	GeoClient        *geo.Client
	ShopifyClient    *shopify.App
	CustomerIOClient *customerio.Client
	HubspotClient    *hubspot.Client
	Enforcer         *casbin.Enforcer
	SheetAPI         *sheets.Service
}

func New(config *config.Configuration) *App {
	return &App{
		Config: config,
	}
}

func (app *App) WithS3Client(client *s3.Client) *App {
	app.S3Client = client
	return app
}

func (app *App) WithMailer(mailer *mailer.Client) *App {
	app.Mailer = mailer
	return app
}

func (app *App) WithDB(db *db.DB) *App {
	app.DB = db
	return app
}

func (app *App) WithAnalyticDB(db *db.DB) *App {
	app.AnalyticDB = db
	return app
}

func (app *App) WithCache(cache *caching.Client) *App {
	app.Cache = cache
	return app
}

func (app *App) WithStripeHelper(helper *stripehelper.StripeClient) *App {
	app.StripeHelper = helper
	return app
}

func (app *App) WithGeoClient(geo *geo.Client) *App {
	app.GeoClient = geo
	return app
}

func (app *App) WithShopifyClient(shopifyClient *shopify.App) *App {
	app.ShopifyClient = shopifyClient
	return app
}

func (app *App) WithCustomerIOClient(client *customerio.Client) *App {
	app.CustomerIOClient = client
	return app
}

func (app *App) WithEnforcer(e *casbin.Enforcer) *App {
	app.Enforcer = e
	return app
}

func (app *App) WithSheetAPI(s *sheets.Service) *App {
	app.SheetAPI = s
	return app
}

func (app *App) WithHubspotClient(c *hubspot.Client) *App {
	app.HubspotClient = c
	return app
}
