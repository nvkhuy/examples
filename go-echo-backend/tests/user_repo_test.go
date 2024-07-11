package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_GetUsers(t *testing.T) {
	var app = initApp("dev")

	var result = repo.NewUserRepo(app.DB).PaginateUsers(repo.PaginateUsersParams{
		//Roles:                      []enums.Role{enums.RoleClient, enums.RoleLeader},
		JwtClaimsInfo:             *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
		IncludeAssignedInquiryIds: true,
		//UserIDs: []string{
		//	"cl712kdbt6ca83b9kk3g", "cl658pd7qfpgrlqqj93g",
		//},
		//Teams: []string{
		//	"marketing",
		//},
		PaginationParams: models.PaginationParams{
			Page:  1,
			Limit: 10,
		},
	})

	helper.PrintJSON(result)
}

func TestUserRepo_GetMe(t *testing.T) {
	var app = initApp()

	result, err := repo.NewUserRepo(app.DB).GetMe(repo.GetMeParams{
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetUserID("cg5anr2llkm6ctpvq8k0"),
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestUserRepo_GetUserByID(t *testing.T) {
	var app = initApp()

	result, err := repo.NewUserRepo(app.DB).GetUser(repo.GetUserParams{
		UserID:        "ck4p2rdq413sb14f1jag",
		JwtClaimsInfo: models.JwtClaimsInfo{},
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestUserRepo_GetCustomerIOUser(t *testing.T) {
	var app = initApp()

	result, err := repo.NewUserRepo(app.DB).GetCustomerIOUser("cir3f55ocd6fbjvrsv10")
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestUserRepo_CreateInvitedUser(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewUserRepo(app.DB).CreateInvitedUser(repo.InviteUserForm{
		Email:       "test_01@gmail.com",
		RedirectURL: "https://dev-brand.joininflow.io/forgot-password",
		Role:        enums.RoleLeader,
		Team:        enums.TeamDesigner,
	})
	assert.Nil(t, err)
	helper.PrintJSON(result)
}

func TestUserRepo_PaginateUsersRole(t *testing.T) {
	var app = initApp("local")
	result, err := repo.NewUserRepo(app.DB).PaginateUsersRole(repo.PaginateUsersRoleParams{
		PaginationParams: models.PaginationParams{},
		JwtClaimsInfo:    *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin),
	})
	if err != nil {
		return
	}
	helper.PrintJSON(result)
}

func TestUserRepo_UpdateUserPasswordDirectly(t *testing.T) {
	var app = initApp("local")
	var user models.User
	user.Password = "1234qwer"
	var err = user.HashPassword()
	if err != nil {
		return
	}
	err = app.DB.Model(&models.User{}).Where("id = ?", "cl658ovbudgljog72hd0").UpdateColumn("Password", user.Password).Error
	if err != nil {
		return
	}
	helper.PrintJSON(user)
}

func TestUserRepo_TeamClientInvite(t *testing.T) {
	var app = initApp("local")
	claims := models.NewJwtClaimsInfo().SetUserID("cku7anijsufrcmh9k20g")
	result, err := repo.NewUserRepo(app.DB).TeamClientInvite(repo.TeamClientInviteForm{
		JwtClaimsInfo: *claims,
		FirstName:     "Huy",
		LastName:      "Nguyen",
		PhoneNumber:   "+84938294687",
		Email:         "invite+01@gmail.com",
		RedirectURL:   "https://dev-brand.joininflow.io/forgot-password",
	})
	assert.Nil(t, err)
	helper.PrintJSON(result)
}

func TestUserRepo_UpdateBrandTeamMemberActions(t *testing.T) {
	var app = initApp("local")
	claims := models.NewJwtClaimsInfo().SetUserID("cku7anijsufrcmh9k20g")
	err := repo.NewUserRepo(app.DB).UpdateBrandTeamMemberActions(repo.UpdateBrandTeamMemberActionsForm{
		JwtClaimsInfo: *claims,
		MemberId:      "cl9jm6fskmpavu4b7dkg",
		Actions: enums.BrandMemberActions{
			enums.BrandMemberActionApproveRFQ,
		},
	})
	assert.Nil(t, err)
}

func TestUserRepo_DeleteBrandTeamMember(t *testing.T) {
	var app = initApp("local")
	claims := models.NewJwtClaimsInfo().SetUserID("cku7anijsufrcmh9k20g")
	err := repo.NewUserRepo(app.DB).DeleteBrandTeamMember(repo.DeleteBrandTeamMemberForm{
		JwtClaimsInfo: *claims,
		MemberId:      "cl9jcovskmpara83jk4g",
	})
	assert.Nil(t, err)
}

func TestUserRepo_PaginateBrandTeamMembers(t *testing.T) {
	var app = initApp("dev")
	claims := models.NewJwtClaimsInfo().SetUserID("clg6guun84ht3df9a60g")
	var members = repo.NewUserRepo(app.DB).PaginateBrandTeamMembers(repo.PaginateBrandTeamMembersParams{
		JwtClaimsInfo: *claims,
	})

	assert.NotNil(t, members)
}

func TestUserRepo_GetAccessToken(t *testing.T) {
	var app = initApp("dev")

	token, err := repo.NewUserRepo(app.DB).GetAccessToken(repo.GetAccessTokenParams{
		UserID: "cg5anr2llkm6ctpvq8k0",
	})

	assert.NoError(t, err)

	fmt.Println("*** token", token)
}

func TestUserRepo_GetLastShippingAddress(t *testing.T) {
	var app = initApp("prod")

	token, err := repo.NewUserRepo(app.DB).GetLastShippingAddress(*models.NewJwtClaimsInfo().SetUserID("cmin48jb2hj6jvehgqr0"))

	assert.NoError(t, err)

	helper.PrintJSON(token)
}

func TestUserRepo_UpdateUser(t *testing.T) {
	var app = initApp("prod")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = []byte(`{"id":"cle8i3kc0j6g390u4u00","created_at":1700563214,"updated_at":1711597838,"role":"client","name":"Thanakit Kamma","first_name":"Thanakit","last_name":"Kamma","email":"adintziez16@gmail.com","phone_number":"+66 99 426 9529","account_status":"active","account_status_changed_at":1700563485,"verification_sent_at":1700563214,"account_verified":true,"account_verified_at":1700563485,"country_code":"TH","timezone":"Asia/Saigon","token_reset_password_sent_at":1701846886,"last_login":1711597838,"is_test":false,"is_first_login":true,"is_offline":false,"last_online_at":1711597838,"stripe_customer_id":"cus_P2xLYOXu4dIe89","coordinate_id":"49d61200b5132f2d2a141a86002d73fb","coordinate":{"id":"49d61200b5132f2d2a141a86002d73fb","created_at":1699613839,"updated_at":1711597980,"postal_code":"","country_code":"VN"},"preferred_currency":"USD","brand_name":"22thoctoberr","total_pcs_required":100,"annual_revenue":"0","monthly_turnover":"0","requirement":"Jean Denim Pants= high quality, lots of details , cargo\nJean Denim Shorts = high quality, lots of details , cargo","completed_inquiry_tutorial_at":1703146805,"hubspot_contact_id":"591201","features":[{"type":"bulk_online_payment","name":"Bulk Online Payment"}],"feature_ids":["bulk_online_payment"]}`)

	var req = httptest.NewRequest(echo.PUT, "/api/v1/admin/users/cle8i3kc0j6g390u4u00", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNqY3ZtaDRwcmY4dDA5OHJiZnYwIiwiYXVkIjoic3VwZXJfYWRtaW4iLCJpc3MiOiJjbW5ucnByYjJoamFhZWo3aDAxZyIsInN1YiI6InN1cGVyX2FkbWluIn0.mmvbDz-iHMLiRp0DsC8Dn7c3XfGgmGf8g2n3GtnYrts")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestUserRepo_UpdateProductClasses(t *testing.T) {
	var app = initApp("local")
	err := repo.NewUserRepo(app.DB).UpdateProductClasses(repo.UpdateUserProductClassesParams{
		JwtClaimsInfo: models.JwtClaimsInfo{},
		UserID:        "cisb5mmvnd8fddu9n920",
		InquiryID:     "cj3o2qumt2dhmng4kfkg",
	})
	assert.Nil(t, err)
}
