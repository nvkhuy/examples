package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type Service string

var (
	ServiceBackend  Service = "backend"
	ServiceConsumer Service = "consumer"
	ServiceCrawler  Service = "crawler"
	ServiceChat     Service = "chat"
)

var instance *Configuration

type CommonInfo struct {
	Name           string `json:"name"`
	DocURL         string `json:"doc_url"`
	StatsURL       string `json:"stats_url"`
	ConfigURL      string `json:"config_url"`
	HealthCheckURL string `json:"health_check_url"`
	BuildInfo
}

type BuildInfo struct {
	BuildServiceName string `json:"build_service_name"`
	BuildEnv         string `json:"env"`
	BuildDate        string `json:"build_date"`
	BuildVersion     string `json:"build_version"`
	BuildNumber      string `json:"build_number"`
	BuildGitSummary  string `json:"build_git_summary"`
	BuildGitSHA1     string `json:"build_git_sha1"`
	BuildGitBranch   string `json:"build_git_branch"`
	BuildGitCommit   string `json:"build_git_commit"`
}

type Configuration struct {
	BuildInfo

	AppName string `mapstructure:"APP_NAME" json:"app_name"`

	EFSPath string `mapstructure:"EFS_PATH" json:"efs_path"`

	// Server
	ServerPort              string `mapstructure:"SERVER_PORT" json:"server_port"`
	ServerName              string `mapstructure:"SERVER_NAME" json:"server_name"`
	ServerBaseURL           string `mapstructure:"SERVER_BASE_URL" json:"server_base_url"`
	ServerWebCrawlerBaseURL string `mapstructure:"SERVER_WEB_CRAWLER_BASE_URL" json:"server_web_crawler_base_url"`

	// Consumer server
	ConsumerServerName    string `mapstructure:"CONSUMER_SERVER_NAME" json:"consumer_server_name"`
	ConsumerServerBaseURL string `mapstructure:"CONSUMER_SERVER_BASE_URL" json:"consumer_server_base_url"`

	AdminUserName     string `mapstructure:"ADMIN_USER_NAME" json:"admin_user_name"`
	AdminUserPassword string `mapstructure:"ADMIN_USER_PASSWORD" json:"admin_user_password"`

	// Database
	DBName        string `mapstructure:"DB_NAME" json:"db_name"`
	DBUser        string `mapstructure:"DB_USER" json:"db_user"`
	DBPassword    string `mapstructure:"DB_PASSWORD" json:"db_password"`
	DBReplicaHost string `mapstructure:"DB_REPLICA_HOST" json:"db_replica_host"`
	DBPort        string `mapstructure:"DB_PORT" json:"db_port"`
	DBHost        string `mapstructure:"DB_HOST" json:"db_host"`
	DBURI         string `mapstructure:"DB_URI" json:"db_uri"`
	DBSSLMode     string `mapstructure:"DB_SSL_MODE" json:"db_sslmode"`

	// Analytic Database
	ADBName        string `mapstructure:"ADB_NAME" json:"adb_name"`
	ADBUser        string `mapstructure:"ADB_USER" json:"adb_user"`
	ADBPassword    string `mapstructure:"ADB_PASSWORD" json:"adb_password"`
	ADBPort        string `mapstructure:"ADB_PORT" json:"adb_port"`
	ADBHost        string `mapstructure:"ADB_HOST" json:"adb_host"`
	ADBReplicaHost string `mapstructure:"ADB_REPLICA_HOST" json:"adb_replica_host"`
	ADBURI         string `mapstructure:"ADB_URI" json:"adb_uri"`
	ADBSSLMode     string `mapstructure:"ADB_SSL_MODE" json:"adb_sslmode"`

	// Mailer
	SendgridAPIKey     string `mapstructure:"SENDGRID_API_KEY" json:"sendgrid_api_key"`
	SenderMail         string `mapstructure:"SENDER_MAIL" json:"sender_mail"`
	SenderName         string `mapstructure:"SENDER_NAME" json:"sender_name"`
	BCCAddresses       string `mapstructure:"BCC_ADDRESSES" json:"bcc_addresses"`
	AdminMailTo        string `mapstructure:"ADMIN_MAIL_TO" json:"admin_mail_to"`
	AdminNameTo        string `mapstructure:"ADMIN_NAME_TO" json:"admin_name_to"`
	AdminContactMailTo string `mapstructure:"ADMIN_CONTACT_MAIL_TO" json:"admin_contact_mail_to"`
	AdminContactNameTo string `mapstructure:"ADMIN_CONTACT_NAME_TO" json:"admin_contact_name_to"`
	SupportEmail       string `mapstructure:"SUPPORT_EMAIL" json:"support_email"`

	// Jwt
	JWTSecret string        `mapstructure:"JWT_SECRET" json:"jwt_secret"`
	JWTExpiry time.Duration `mapstructure:"JWT_EXPIRY" json:"jwt_expiry"`

	JWTAssetSecret string        `mapstructure:"JWT_ASSET_SECRET" json:"jwt_asset_secret"`
	JWTAssetExpiry time.Duration `mapstructure:"JWT_ASSET_EXPIRY" json:"jwt_asset_expiry"`

	JWTResetPasswordSecret string        `mapstructure:"JWT_RESET_PASSWORD_SECRET" json:"jwt_reset_password_secret"`
	JWTResetPasswordExpiry time.Duration `mapstructure:"JWT_RESET_PASSWORD_EXPIRY" json:"jwt_reset_password_expiry"`

	JWTEmailVerificationSecret string        `mapstructure:"JWT_EMAIL_VERIFICATION_SECRET" json:"jwt_email_verification_secret"`
	JWTEmailVerificationExpiry time.Duration `mapstructure:"JWT_EMAIL_VERIFICATION_EXPIRY" json:"jwt_email_verification_expiry"`

	// AWS
	AWSAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID" json:"aws_access_key_id"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY" json:"aws_secret_access_key"`

	AWSS3Region                   string `mapstructure:"AWS_REGION" json:"aws_region"`
	AWSProfile                    string `mapstructure:"TF_VAR_profile" json:"aws_profile"`
	AWSS3ACL                      string `mapstructure:"AWS_S3_ACL" json:"aws_s3_acl"`
	AWSS3MaxFileSize              int    `mapstructure:"AWS_S3_MAX_FILE" json:"aws_s3_max_file_size"`
	AWSS3SignatureExpiryInMinutes int    `mapstructure:"AWS_S3_SIGNATURE_EXPIRY_IN_MINUTES" json:"aws_s3_signature_expiry_in_minutes"`
	AWSS3StorageBucket            string `mapstructure:"AWS_S3_STORAGE_BUCKET" json:"aws_s3_storage_bucket"`
	AWSS3TrendingBucket           string `mapstructure:"AWS_S3_TRENDING_BUCKET" json:"aws_s3_trending_bucket"`
	AWSS3CdnBucket                string `mapstructure:"AWS_S3_CDN_BUCKET" json:"aws_s3_cdn_bucket"`
	AWSS3WebScraperBucket         string `mapstructure:"AWS_S3_WEBSCRAPER_BUCKET" json:"aws_3_webscraper_bucket"`

	// Redis
	RedisNamespace string   `mapstructure:"REDIS_NAMESPACE" json:"redis_namespace"`
	RedisAddress   []string `mapstructure:"REDIS_ADDRESS" json:"redis_address"`
	RedisPassword  string   `mapstructure:"REDIS_PASSWORD" json:"redis_password"`
	RedisDBMode    int      `mapstructure:"REDIS_DB_MODE" json:"redis_db_mode"`
	// Chat server
	ChatServerName    string `mapstructure:"CHAT_SERVER_NAME" json:"chat_server_name"`
	ChatServerBaseURL string `mapstructure:"CHAT_SERVER_BASE_URL" json:"chat_server_base_url"`

	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID" json:"google_client_id"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET" json:"google_client_secret"`

	WebAppBaseURL       string `mapstructure:"WEB_APP_BASE_URL" json:"web_app_base_url"`
	BrandPortalBaseURL  string `mapstructure:"BRAND_PORTAL_BASE_URL" json:"brand_portal_base_url"`
	AdminPortalBaseURL  string `mapstructure:"ADMIN_PORTAL_BASE_URL" json:"admin_portal_base_url"`
	SellerPortalBaseURL string `mapstructure:"SELLER_PORTAL_BASE_URL" json:"seller_portal_base_url"`
	ShareBaseURL        string `mapstructure:"SHARE_BASE_URL" json:"share_base_url"`

	StripeDashboardURL     string `mapstructure:"STRIPE_DASHBOARD_URL" json:"stripe_dashboard_url"`
	StripeSecretKey        string `mapstructure:"STRIPE_SECRET_KEY" json:"stripe_secret_key"`
	StripePublicKey        string `mapstructure:"STRIPE_PUBLIC_KEY" json:"stripe_public_key"`
	StripeWebhookSecretKey string `mapstructure:"STRIPE_WEBHOOK_SECRET_KEY" json:"stripe_webhook_secret_key"`

	ShopifyApiKey    string `mapstructure:"SHOPIFY_API_KEY" json:"shopify_api_key"`
	ShopifyApiSecret string `mapstructure:"SHOPIFY_API_SECRET" json:"shopify_api_secret"`

	CustomerIOSiteID         string `mapstructure:"CUSTOMERIO_SITE_ID" json:"customerio_side_id"`
	CustomerIOApiAppKey      string `mapstructure:"CUSTOMERIO_API_APP_KEY" json:"customerio_api_app_key"`
	CustomerIOApiTrackingKey string `mapstructure:"CUSTOMERIO_API_TRACKING_KEY" json:"customerio_api_tracking_key"`

	LambdaAPIResizeURL string `mapstructure:"LAMBDA_API_RESIZE_URL" json:"lambda_api_resize_url"`
	LambdaAPIPDFURL    string `mapstructure:"LAMBDA_API_PDF_URL" json:"lambda_api_pdf_url"`
	LambdaAPIFfmpegURL string `mapstructure:"LAMBDA_API_FFMPEG_URL" json:"lambda_api_ffmpef_url"`
	LambdaAPIBlurURL   string `mapstructure:"LAMBDA_API_BLUR_URL" json:"lambda_api_blur_url"`

	SuperAdminUserID string `mapstructure:"SUPER_ADMIN_USER_ID" json:"super_admin_user_id"`

	StorageURL   string `mapstructure:"TF_VAR_storage_domain" json:"storage_domain"`
	CDNURL       string `mapstructure:"TF_VAR_cdn_domain" json:"cdn_domain"`
	CDNStaticURL string `mapstructure:"CDN_STATIC_DOMAIN" json:"cdn_static_domain"`

	ResetPasswordResendInterval time.Duration `mapstructure:"RESET_PASSWORD_RESEND_INTERVAL" json:"reset_password_resend_interval"`

	CasbinModelConfURL    string `mapstructure:"CASBIN_MODEL_CONF_URL" json:"casbin_model_conf_url"`
	CasbinPolicyCSVURL    string `mapstructure:"CASBIN_POLICY_CSV_URL" json:"casbin_policy_csv_url"`
	GoogleClientSecretURL string `mapstructure:"GOOGLE_CLIENT_SECRET_URL" json:"google_client_secret_url"`
	QRCodeLogoURL         string `mapstructure:"QR_CODE_LOGO_URL" json:"qr_code_logo_url"`

	InflowSaleGroupEmail        string `mapstructure:"INFLOW_SALE_GROUP_EMAIL" json:"inflow_sale_group_email"`
	InflowMerchandiseGroupEmail string `mapstructure:"INFLOW_MERCHANDISE_GROUP_EMAIL" json:"inflow_merchandise_group_email"`
	OpenAIKey                   string `mapstructure:"OPEN_AI_KEY" json:"open_ai_key"`
	OpenAIAssistant             string `mapstructure:"OPEN_AI_ASSISTANT" json:"open_ai_assistant"`

	PayosClientID    string `mapstructure:"PAYOS_CLIENT_ID" json:"payos_client_id"`
	PayosApiKey      string `mapstructure:"PAYOS_API_KEY" json:"payos_api_key"`
	PayosChecksumKey string `mapstructure:"PAYOS_CHECKSUM_KEY" json:"payos_checksum_key"`

	RsaPublicPemFile  string `mapstructure:"RSA_PUBLIC_PEM_FILE" json:"rsa_public_pem_file"`
	RsaPrivatePemFile string `mapstructure:"RSA_PRIVATE_PEM_FILE" json:"rsa_private_pem_file"`
	RsaSecret         string `mapstructure:"RSA_SECRET" json:"rsa_secret"`

	MediaJWTSecret string `mapstructure:"TF_VAR_media_jwt_secret" json:"media_jwt_secret"`
	MediaBaseURL   string `mapstructure:"MEDIA_BASE_URL" json:"media_base_url"`
	MediaResizeURL string `mapstructure:"MEDIA_RESIZE_URL" json:"media_resize_url"`
	SentryDsn      string `mapstructure:"SENTRY_DSN" json:"sentry_dsn"`

	HubspotAccessToken string `mapstructure:"HUBSPOT_ACCESS_TOKEN" json:"hubspot_access_token"`

	CheckoutJwtSecret string `mapstructure:"CHECKOUT_JWT_SECRET" json:"checkout_jwt_secret"`
	CheckoutJwtExpiry string `mapstructure:"CHECKOUT_JWT_EXPIRY" json:"checkout_jwt_expiry"`

	RunningTest bool `json:"-"`
}

// New setup config
func New(cfgFile string, buildInfo ...BuildInfo) *Configuration {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".go-base" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	var err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	instance = &Configuration{}

	err = viper.Unmarshal(&instance)
	if err != nil {
		panic(err)
	}

	if len(buildInfo) > 0 {
		instance.BuildInfo = buildInfo[0]

		if instance.BuildInfo.BuildServiceName == "" {
			instance.BuildInfo.BuildServiceName = viper.GetString("SERVICE_NAME")
		}

		if instance.BuildInfo.BuildEnv == "" {
			instance.BuildInfo.BuildEnv = viper.GetString("ENV")
		}
	}

	if port := viper.GetString("SERVER_PORT"); port != "" {
		instance.ServerPort = port
	}
	return instance
}

// GetInstance get instance
func GetInstance() *Configuration {
	if instance == nil {
		panic("instance is getting NULL")
	}
	return instance
}

// GetServerName get app name including env
func (c Configuration) GetServerName(service Service) (name string) {
	switch service {
	case ServiceConsumer:
		name = fmt.Sprintf("[%s] %s", c.BuildEnv, c.ConsumerServerName)
	default:
		name = fmt.Sprintf("[%s] %s", c.BuildEnv, c.ServerName)
	}
	return
}

func (c Configuration) GetDocDescription() string {
	var buildDate = c.BuildDate
	localDate, err := time.Parse(time.RFC3339, c.BuildDate)
	if err == nil {
		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err == nil {
			buildDate = localDate.In(loc).Format(time.RFC822)
		}

	}
	var texts = []string{
		fmt.Sprintf("<strong>This is api docs of %s. Information of this version:</strong>", c.BuildServiceName),
		fmt.Sprintf("- Build Version:  <strong>%s</strong>", c.BuildVersion),
		fmt.Sprintf("- Build Date:     <strong>%s</strong>", buildDate),
		fmt.Sprintf("- Git Branch:     <strong>%s</strong>", c.BuildGitBranch),
		fmt.Sprintf("- Git Summary:    <strong>%s</strong>", c.BuildGitSummary),
		fmt.Sprintf("- Git SHA1:       <strong>%s</strong>", c.BuildGitSHA1),
		fmt.Sprintf("- Git Commit:     <strong>%s</strong>", c.BuildGitCommit),
	}

	return strings.Join(texts, "\n")
}

// GetServerCommonInfo common info
func (c Configuration) GetServerCommonInfo(service Service) (info CommonInfo) {
	var baseURL = c.ServerBaseURL
	info = CommonInfo{
		Name:      c.GetServerName(service),
		BuildInfo: c.BuildInfo,
	}

	switch service {
	case ServiceConsumer:
		baseURL = c.ConsumerServerBaseURL
	case ServiceChat:
		baseURL = c.ChatServerBaseURL
	}

	urlInfo, _ := url.Parse(baseURL)

	if _, e := strconv.ParseInt(urlInfo.Port(), 10, 64); e == nil && urlInfo.Port() != c.ServerPort {
		baseURL = strings.ReplaceAll(baseURL, urlInfo.Port(), c.ServerPort)
	}

	info.HealthCheckURL = fmt.Sprintf("%s/health_check", baseURL)
	info.DocURL = fmt.Sprintf("%s/api/docs/index.html", baseURL)
	info.ConfigURL = fmt.Sprintf("%s/api/app/config", baseURL)

	return
}

func (c Configuration) GetDatabaseURI() string {
	var dbHost = c.DBHost
	var dbPort = c.DBPort
	var dbName = c.DBName
	var dbUser = c.DBUser
	var dbPassword = c.DBPassword
	var dbSSLMode = c.DBSSLMode

	var dbURI = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbPort, dbUser, dbName, dbPassword, dbSSLMode)

	return dbURI
}

func (c Configuration) GetDatabaseReplicaURI() string {
	var dbHost = c.DBReplicaHost
	var dbPort = c.DBPort
	var dbName = c.DBName
	var dbUser = c.DBUser
	var dbPassword = c.DBPassword
	var dbSSLMode = c.DBSSLMode

	var dbURI = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbPort, dbUser, dbName, dbPassword, dbSSLMode)

	return dbURI
}

func (c Configuration) GetAnalyticDatabaseURI() string {
	var dbHost = c.ADBHost
	var dbPort = c.ADBPort
	var dbName = c.ADBName
	var dbUser = c.ADBUser
	var dbPassword = c.ADBPassword
	var dbSSLMode = c.ADBSSLMode

	var dbURI = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbPort, dbUser, dbName, dbPassword, dbSSLMode)

	return dbURI
}

func (c Configuration) GetAnalyticDatabaseReplicaURI() string {
	var dbHost = c.ADBReplicaHost
	var dbPort = c.ADBPort
	var dbName = c.ADBName
	var dbUser = c.ADBUser
	var dbPassword = c.ADBPassword
	var dbSSLMode = c.ADBSSLMode

	var dbURI = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbPort, dbUser, dbName, dbPassword, dbSSLMode)

	return dbURI
}

func (c Configuration) IsLocal() bool {
	return c.BuildEnv == "local"
}

func (c Configuration) IsDev() bool {
	return c.BuildEnv == "dev"
}

func (c Configuration) IsProd() bool {
	return c.BuildEnv == "prod"
}

func (c Configuration) IsTest() bool {
	return c.BuildEnv == "test" || c.BuildEnv == ""
}

func (c *Configuration) GetMediaToken(audience string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
		Audience:  audience,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(c.MediaJWTSecret))

}
