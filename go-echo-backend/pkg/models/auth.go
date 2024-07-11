package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/location"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/lib/pq"
)

type GoogleOauthForm struct {
	Platform enums.Role `json:"platform" query:"platform" param:"platform" validate:"required"`
}
type LoginEmailForm struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`

	IsAdminPortal bool `json:"-"`
	IsSeller      bool `json:"-"`
}

type LoginGoogleForm struct {
	Token string `json:"token" validate:"token"`

	IsAdminPortal bool `json:"-"`
	IsSeller      bool `json:"-"`
}

type RegisterForm struct {
	ID string `json:"id"`

	Email       string      `json:"email" validate:"email"`
	Password    string      `json:"password" validate:"required"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Name        string      `json:"name"`
	Avatar      *Attachment `json:"avatar,omitempty"`
	PhoneNumber string      `json:"phone_number,omitempty" param:"phone_number" query:"phone_number" form:"phone_number" validate:"isPhone"`

	RegisterAreas             enums.RegisterArea     `json:"register_area,omitempty"`
	RegisterBusiness          enums.RegisterBusiness `json:"register_business,omitempty"`
	RegisterQuantity          enums.RegisterQuantity `json:"register_quantity,omitempty"`
	RegisterProductCategories []string               `json:"register_product_categories,omitempty"`
	CountryCode               enums.CountryCode      `gorm:"default:'VN'" json:"country_code,omitempty"`

	BrandRegisterInfo
	SupplierRegisterInfo

	Coordinate location.Coordinate `json:"coordinate,omitempty"`

	BusinessProfile *BusinessProfileCreateForm `gorm:"-" json:"profile"`

	IsSeller bool `json:"-"`
}

type BrandRegisterInfo struct {
	BrandName        string                  `json:"brand_name,omitempty"`
	BrandAge         int64                   `json:"brand_age,omitempty"`
	BrandWebsite     string                  `json:"brand_website,omitempty"`
	TotalPcsRequired int64                   `json:"total_pcs_required,omitempty"`
	AnnualRevenue    price.Price             `json:"annual_revenue,omitempty"`
	MonthlyTurnover  price.Price             `json:"monthly_turnover,omitempty"`
	Requirement      string                  `json:"requirement,omitempty" validate:"max=5000"`
	Product          *UserProductCaptureMeta `json:"product,omitempty"`
}

type SupplierRegisterInfo struct {
	SupplierType                   enums.SupplierType `json:"supplier_type,omitempty"`
	YourDesignation                string             `json:"your_designation,omitempty"`
	CompanyName                    string             `json:"company_name,omitempty"`
	CompanyEmail                   string             `json:"company_email,omitempty"`
	CompanyWebsite                 string             `json:"company_website,omitempty"`
	BusinessLicenseNumber          string             `json:"business_license_number,omitempty"`
	SupplierProfileAttachments     *Attachments       `json:"supplier_profile_attachments,omitempty"`
	SupplierCertificate            string             `json:"supplier_certificate,omitempty"` // separete by comma
	SupplierCertificateAttachments *Attachments       `json:"supplier_certificate_attachments,omitempty"`
	PaymentTerms                   pq.StringArray     `gorm:"type:varchar(200)[]" json:"payment_terms,omitempty" swaggertype:"array,string"`
}

type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type UpdatePasswordForm struct {
	NewPassword        string `json:"new_password" validate:"required"`
	TokenResetPassword string `json:"token_reset_password" validate:"required"`
}

type UpdateMyPasswordForm struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}
type ForgotPasswordForm struct {
	Email       string `json:"email" validate:"required,email"`
	RedirectURL string `json:"redirect_url" validate:"required,startswith=http"`

	IsAdminPortal    bool  `json:"-"`
	IsSeller         bool  `json:"-"`
	User             *User `json:"-"`
	IsFromInvitation bool  `json:"-"`
}

type ForgotPasswordResponse struct {
	Email       string `json:"email"`
	RedirectURL string `json:"redirect_url"`

	Message       string `json:"message"`
	NextInSeconds int    `json:"next_in_seconds"`

	User *User `json:"-"`
}

type ResendVerificationEmailResponse struct {
	Message       string `json:"message"`
	NextInSeconds int    `json:"next_in_seconds"`
}
type ResetPasswordForm struct {
	NewPassword        string `json:"new_password" validate:"required"`
	TokenResetPassword string `json:"token_reset_password" validate:"required"`

	IsAdminPortal bool `json:"-"`
	IsSeller      bool `json:"-"`
}

type VerifyEmailForm struct {
	JwtClaimsInfo

	Token       string `query:"token" validate:"required"`
	RedirectURL string `query:"redirect_url" validate:"required"`
}

type ChangePasswordForm struct {
	JwtClaimsInfo

	UserID string `param:"user_id" validate:"required"`

	NewPassword string `json:"new_password" validate:"required"`
}
