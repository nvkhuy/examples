package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/location"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/lib/pq"
)

type Users []*User

// User user's model
type User struct {
	Model

	// Common
	Role enums.Role `gorm:"default:'client'" json:"role,omitempty"`
	Team enums.Team `json:"team,omitempty"`

	Name        string `gorm:"size:200" json:"name,omitempty"`
	FirstName   string `gorm:"size:100" json:"first_name,omitempty"`
	LastName    string `gorm:"size:100" json:"last_name,omitempty"`
	Email       string `gorm:"unique;type:citext;default:null" json:"email,omitempty"`
	Description string `gorm:"size:200" json:"description,omitempty"`
	PhoneNumber string `gorm:"unique;default:null" json:"phone_number,omitempty" validate:"isPhone"`

	AccountStatus          enums.AccountStatus `gorm:"default:'pending_review'" json:"account_status,omitempty"`
	AccountStatusChangedAt *int64              `json:"account_status_changed_at,omitempty"`

	VerificationToken  *string `gorm:"default:''" json:"-"`
	VerificationSentAt *int64  `json:"verification_sent_at,omitempty"`
	AccountVerified    bool    `gorm:"default:false" json:"account_verified,omitempty"`
	AccountVerifiedAt  *int64  `json:"account_verified_at,omitempty"`

	CountryCode    enums.CountryCode `gorm:"default:'VN';size:100" json:"country_code,omitempty"`
	State          string            `gorm:"default:'';size:100" json:"state,omitempty"`
	ZipCode        string            `gorm:"default:'';size:100" json:"zip_code,omitempty"`
	Timezone       enums.Timezone    `gorm:"default:'';size:100" json:"timezone,omitempty"`
	BillingAddress string            `gorm:"default:'';size:100" json:"billing_address,omitempty"`

	// Photo
	Avatar *Attachment `json:"avatar,omitempty"`

	TokenResetPassword       *string `gorm:"default:'';size:1024" json:"-"`
	TokenIssuer              string  `gorm:"default:'';size:200" json:"-"`
	Password                 string  `gorm:"not null;size:200" json:"-"`
	TokenResetPasswordSentAt *int64  `json:"token_reset_password_sent_at,omitempty"`

	LastLogin    *int64 `gorm:"default:null" json:"last_login,omitempty"`
	LoggedOutAt  *int64 `gorm:"default:null" json:"-"`
	LastActivity *int64 `gorm:"default:null" json:"last_activity,omitempty"`

	IsTest *bool `gorm:"default:false" json:"is_test,omitempty"`

	IsGhost bool `gorm:"-" json:"-"`

	IsFirstLogin *bool  `gorm:"default:true" json:"is_first_login,omitempty"`
	IsOffline    *bool  `gorm:"default:false" json:"is_offline,omitempty"`
	LastOnlineAt *int64 `gorm:"default:null" json:"last_online_at,omitempty"`

	SocialProvider string `gorm:"size:100" json:"social_provider,omitempty"`

	StripeCustomerID string `gorm:"size:100" json:"stripe_customer_id,omitempty"`

	IsNew bool `gorm:"-" json:"-"`

	CoordinateID string      `gorm:"size:100" json:"coordinate_id,omitempty"`
	Coordinate   *Coordinate `gorm:"-" json:"coordinate,omitempty"`

	BusinessProfile *BusinessProfile `gorm:"-" json:"business_profile,omitempty"`

	OnboardingSubmitAt *int64 `json:"onboarding_submit_at,omitempty"`

	PreferredCurrency enums.Currency `gorm:"default:'USD'" json:"preferred_currency,omitempty"`

	ContactOwnerIDs pq.StringArray `gorm:"type:varchar(100)[]" json:"contact_owner_ids,omitempty"`
	ContactOwners   []*User        `gorm:"-" json:"contact_owners,omitempty"`

	RequiresChangePassword bool `json:"requires_change_password,omitempty"`

	BrandTeam      *BrandTeam         `gorm:"-" json:"brand_team,omitempty"`
	ProductClasses UserProductClasses `json:"product_classes,omitempty"`

	BrandInfo
	SupplierInfo
	UserSettings
	// Extra info for API response
	InquirySeller      *InquirySeller `gorm:"-" json:"inquiry_seller,omitempty"`
	ZaloID             string         `gorm:"size:100" json:"zalo_id,omitempty"`
	AssignedInquiryIDs pq.StringArray `gorm:"-" json:"assigned_inquiry_ids,omitempty"`
	AssignedPOIDs      pq.StringArray `gorm:"-" json:"assigned_po_ids,omitempty"`
	AssignedBulkPOIDs  pq.StringArray `gorm:"-" json:"assigned_bulk_po_ids,omitempty"`

	HubspotOwnerID   string `gorm:"size:200" json:"hubspot_owner_id,omitempty"`
	HubspotContactID string `gorm:"size:200" json:"hubspot_contact_id,omitempty"`

	Features Features `json:"features,omitempty"`
}

type UserSettings struct {
	CompletedInquiryTutorialAt *int64 `json:"completed_inquiry_tutorial_at,omitempty"`
}
type BrandInfo struct {
	BrandType        enums.BrandType         `gorm:"size:200" json:"brand_type,omitempty"`
	BrandName        string                  `gorm:"size:200" json:"brand_name,omitempty"`
	BrandAge         int64                   `gorm:"size:200" json:"brand_age,omitempty"`
	BrandWebsite     string                  `gorm:"size:200" json:"brand_website,omitempty"`
	Representative   string                  `gorm:"size:200" json:"representative,omitempty"`
	Bio              string                  `gorm:"size:1000" json:"bio,omitempty"`
	TaxID            string                  `gorm:"size:200" json:"tax_id,omitempty"`
	TotalPcsRequired int64                   `json:"total_pcs_required,omitempty"`
	AnnualRevenue    *price.Price            `json:"annual_revenue,omitempty"`
	MonthlyTurnover  *price.Price            `json:"monthly_turnover,omitempty"`
	Requirement      string                  `gorm:"size:5000" json:"requirement,omitempty"`
	Product          *UserProductCaptureMeta `json:"product,omitempty"`
}

type SupplierInfo struct {
	SupplierType                   enums.SupplierType `gorm:"size:200" json:"supplier_type,omitempty"`
	YourDesignation                *string            `gorm:"size:200" json:"your_designation,omitempty"`
	CompanyName                    string             `gorm:"size:200" json:"company_name,omitempty"`
	CompanyEmail                   string             `gorm:"size:200" json:"company_email,omitempty"`
	CompanyWebsite                 string             `gorm:"size:200" json:"company_website,omitempty"`
	BusinessLicenseNumber          string             `gorm:"size:200" json:"business_license_number,omitempty"`
	SupplierProfileAttachments     *Attachments       `json:"supplier_profile_attachments,omitempty"`
	SupplierCertificate            string             `gorm:"size:200" json:"supplier_certificate,omitempty"` // separete by comma
	SupplierCertificateAttachments *Attachments       `json:"supplier_certificate_attachments,omitempty"`
	PaymentTerms                   pq.StringArray     `gorm:"type:varchar(200)[]" json:"payment_terms,omitempty" swaggertype:"array,string"`
	CoverAttachment                *Attachment        `json:"cover_attachment,omitempty"`
	SupplierCatalogAttachments     *Attachments       `json:"supplier_catalog_attachments,omitempty"`
}

type AdminUserUpdateForm struct {
	Role enums.Role `gorm:"default:'worker'" json:"role"`
	Team string     `json:"team"`

	AccountStatus          enums.AccountStatus `gorm:"default:'pending_review'" json:"account_status,omitempty"`
	AccountStatusChangedAt *int64              `json:"account_status_changed_at,omitempty"`

	Features Features `json:"features"`

	UserUpdateForm
}

type UserUpdateForm struct {
	JwtClaimsInfo
	UserID string `param:"user_id"`

	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `gorm:"unique;type:citext;default:null" json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty" param:"phone_number" query:"phone_number" form:"phone_number" validate:"omitempty,isPhone"`

	CountryCode    enums.CountryCode `gorm:"default:'VN'" json:"country_code,omitempty"`
	State          string            `gorm:"default:'';size:100" json:"state,omitempty"`
	ZipCode        string            `gorm:"default:'';size:100" json:"zip_code,omitempty"`
	Timezone       enums.Timezone    `gorm:"default:'';size:100" json:"timezone,omitempty"`
	BillingAddress string            `gorm:"default:'';size:100" json:"billing_address,omitempty"`

	Avatar *Attachment `json:"avatar,omitempty"`

	Coordinate location.Coordinate `json:"coordinate,omitempty"`

	PreferredCurrency enums.Currency `gorm:"default:'USD'" json:"preferred_currency,omitempty"`

	BrandInfo
	SupplierInfo
}
type UserTrackActivityForm struct {
	CountryCode enums.CountryCode `json:"country_code"`
	Timezone    enums.Timezone    `json:"timezone"`
}

type SellerDashboardResponse struct {
	RecommendedProducts interface{} `json:"recommended_products,omitempty"`
	News                []string    `json:"news,omitempty"`
	Messages            []string    `json:"messages,omitempty"`
}

type SellerDashboardRevenueResponse struct {
	MonthlyTrend []MonthlyRevenue `json:"monthly_trend,omitempty"`
	WeeklyTrend  []WeekLyRevenue  `json:"weekly_trend,omitempty"`

	MonthlyIncreaseRate float32 `json:"monthly_increase_rate,omitempty"`
	WeeklyIncreaseRate  float32 `json:"weekly_increase_rate,omitempty"`

	MonthlyRevenue price.Price `json:"monthly_revenue,omitempty"`
	WeeklyRevenue  price.Price `json:"weekly_revenue,omitempty"`
}

type MonthlyRevenue struct {
	Month   string      `json:"month,omitempty"`
	Revenue price.Price `json:"revenue,omitempty"`
}

type WeekLyRevenue struct {
	Week    string      `json:"week,omitempty"`
	Revenue price.Price `json:"revenue,omitempty"`
}

type CreateInvitedUserResponse struct {
	User          User   `json:"-"`
	InvitedByUser User   `json:"-"`
	RedirectURL   string `json:"redirect_url"`
}

type UserRoleStat struct {
	Team  string `json:"team"`
	Count int    `json:"count"`
}

type UserPaymentMethod struct {
	ID        string `json:"id,omitempty"`
	CreatedAt int64  `json:"created_at"`
	Name      string `json:"name,omitempty"`
	Last4     string `json:"last4,omitempty"`
	Brand     string `json:"brand,omitempty"`
	Type      string `json:"type,omitempty"`
	IsDefault *bool  `json:"is_default"`
	ExpMonth  int64  `json:"exp_month"`
	ExpYear   int64  `json:"exp_year"`
}

type AssignContactOwnersForm struct {
	JwtClaimsInfo

	UserID          string   `json:"user_id" param:"user_id" query:"user_id" validate:"required"`
	ContactOwnerIDs []string `json:"contact_owner_ids" validate:"required"`
}

type UserIDParam struct {
	JwtClaimsInfo

	UserID string `json:"user_id" param:"user_id" validate:"required"`

	customerio.GetActivitiesParams
}
