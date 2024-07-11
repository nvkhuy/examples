package tests

import (
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestInquiryBuyer_ValidateAction(t *testing.T) {
	var app = initApp("local")
	err := repo.NewInquiryBuyerRepo(app.DB).ValidateTeamMemberUpdateAction(
		"cl9jm6fskmpavu4b7dkg",
		"ciqd78ctq0u8gav5nu50",
		enums.BrandMemberActionApproveRFQ)

	assert.Nil(t, err)
}

func TestInquiryBuyer_BuyerApproveInquiryQuotation(t *testing.T) {
	var app = initApp()
	result, err := repo.NewInquiryBuyerRepo(app.DB).ApproveInquiryQuotation(repo.BuyerApproveInquiryQuotationParams{
		InquiryID: "",
		ApproveRejectMeta: &models.InquiryApproveRejectMeta{
			Comment: "This is find",
		},
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}

func TestInquiryBuyer_BuyerPaginateInquiry(t *testing.T) {
	var app = initApp("dev")
	claims := models.NewJwtClaimsInfo().SetUserID("cjvso05ooc2b8f45a1mg")
	result := repo.NewInquiryBuyerRepo(app.DB).PaginateInquiry(repo.PaginateInquiryBuyerParams{
		PaginationParams:  models.PaginationParams{Page: 1, Limit: 12},
		JwtClaimsInfo:     *claims,
		IncludeCollection: true,
	})

	assert.NotNil(t, result)

	helper.PrintJSON(result)
}

func TestInquiryBuyer_UpdateAttachments(t *testing.T) {
	var app = initApp("local")
	claims := models.NewJwtClaimsInfo().SetUserID("cinqo1cav87smdts4js0")
	result, err := repo.NewInquiryBuyerRepo(app.DB).UpdateAttachments(repo.UpdateAttachmentsParams{
		JwtClaimsInfo: *claims,
		InquiryID:     "ciqe3hgpfgj7pe1llu90",
		Attachments: &models.Attachments{
			&models.Attachment{
				FileKey: "uploads/media/att_1.png",
			},
		},
		//FabricAttachments: &models.Attachments{
		//	&models.Attachment{
		//		FileKey: "uploads/media/fabric_1.png",
		//	},
		//},
	})
	assert.NoError(t, err)

	helper.PrintJSON(result)
}
