package main

import (
	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/caching"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/geo"
	"github.com/engineeringinflow/inflow-backend/pkg/googlesheet"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/engineeringinflow/inflow-backend/pkg/locker"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/mailer"
	"github.com/engineeringinflow/inflow-backend/pkg/migration"
	"github.com/engineeringinflow/inflow-backend/pkg/oauth"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/engineeringinflow/inflow-backend/pkg/seeder"
	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/engineeringinflow/inflow-backend/services/consumer"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		var config = config.New(cfgFile, config.BuildInfo{
			BuildEnv:         Env,
			BuildServiceName: Name,
			BuildVersion:     Version,
			BuildNumber:      BuildNumber,
			BuildDate:        BuildDate,
			BuildGitBranch:   GitBranch,
			BuildGitSHA1:     GitCommitFull,
			BuildGitCommit:   GitCommitMsgFull,
			BuildGitSummary:  GitSummary,
		})

		var logger = logger.Init(
			logger.WithLogDir("logs/backend"),
			logger.WithDebug(true),
			logger.WithConsole(true),
		)
		defer logger.Sync()

		oauth.New(config)

		locker.New(config)

		var cache = caching.New(config)

		var adb = db.NewAnalytic(config, callback.New(), cache)
		var db = db.New(config, callback.New(), cache)

		migration.New(db).AutoMigrate()

		seeder.New(db).SeedAccounts()

		var mailer = mailer.New(config)

		var s3Client = s3.New(config)

		var stripClient = stripehelper.New(config)

		var geoClient = geo.New(config)

		var sheetAPI = googlesheet.New(config)

		var hubspotClient = hubspot.New(config)

		var app = app.New(config).
			WithCache(cache).
			WithDB(db).
			WithAnalyticDB(adb).
			WithMailer(mailer).
			WithS3Client(s3Client).
			WithHubspotClient(hubspotClient).
			WithStripeHelper(stripClient).
			WithGeoClient(geoClient).
			WithCustomerIOClient(customerio.New(config)).
			WithShopifyClient(shopify.New(config)).
			WithSheetAPI(sheetAPI)

		go func() {
			consumer.New(app, false)
		}()

		var router = routes.NewRouter(app)

		router.ListenAndServe()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
