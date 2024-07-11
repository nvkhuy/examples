package tests

import (
	"fmt"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/googlesheet"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/locker"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/migration"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
)

func initConfig(env ...string) *config.Configuration {
	var defaultEnv = "dev"
	if len(env) > 0 {
		defaultEnv = env[0]
	}

	logger.Init(logger.WithDebug(true), logger.WithConsole(true), logger.WithJSON(false))

	var cfg = config.New(fmt.Sprintf("../deployment/config/%s/env.json", defaultEnv))
	cfg.GoogleClientSecretURL = "../static/google_client_secret.json"
	cfg.QRCodeLogoURL = "../static/logo.jpg"
	cfg.RsaPrivatePemFile = fmt.Sprintf("../deployment/config/%s/private.pem", defaultEnv)
	cfg.RsaPublicPemFile = fmt.Sprintf("../deployment/config/%s/public.pem", defaultEnv)
	cfg.RunningTest = true
	return cfg
}

func initApp(env ...string) *app.App {
	var defaultEnv = "dev"
	if len(env) > 0 {
		defaultEnv = env[0]
	}

	var cfg = initConfig(defaultEnv)
	locker.New(cfg)

	var s3Client = s3.New(cfg)

	var sheetAPI = googlesheet.New(cfg)

	return app.New(cfg).
		WithStripeHelper(stripehelper.New(cfg)).
		WithDB(db.New(cfg, callback.New(), nil)).
		WithAnalyticDB(db.NewAnalytic(cfg, callback.New(), nil)).
		WithSheetAPI(sheetAPI).WithShopifyClient(shopify.New(cfg)).WithS3Client(s3Client)
}

func initMigration(app *app.App) {
	migration.New(app.DB).AutoMigrate()
}

func TestDB_Migration(t *testing.T) {
	var app = initApp("local")

	migration.New(app.DB).AutoMigrate()
}

func TestDB_MigrationUser(t *testing.T) {
	var app = initApp("prod")

	app.DB.AutoMigrate(&models.User{})
}

func TestDB_DuplicateUser(t *testing.T) {
	var app = initApp("dev")

	var err = app.DB.Create(&models.User{Email: "loithai@joininflow.io"}).Error

	duplicated, pgErr := app.DB.IsDuplicateConstraint(err)

	fmt.Println("duplicated", duplicated)

	fmt.Println("pgErr", pgErr)
}

func TestDB_AutoMigrateBrandTeam(t *testing.T) {
	var app = initApp("prod")

	var err = app.DB.AutoMigrate(&models.BrandTeam{})

	fmt.Println("*** err", err)
}

func TestDB_AutoMigratePurchaseOrderItem(t *testing.T) {
	var app = initApp("dev")

	var err = app.DB.AutoMigrate(&models.PurchaseOrderItem{})

	fmt.Println("*** err", err)
}

func TestDB_AutoMigratePurchaseOrder(t *testing.T) {
	var app = initApp("dev")

	var err = app.DB.AutoMigrate(&models.PurchaseOrder{})

	fmt.Println("*** err", err)
}
