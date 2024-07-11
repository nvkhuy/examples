package repo

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"golang.org/x/exp/slices"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v74"
	stripePrice "github.com/stripe/stripe-go/v74/price"
	"github.com/thaitanloi365/go-utils/values"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InquiryRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewInquiryRepo(db *db.DB) *InquiryRepo {
	return &InquiryRepo{
		db:     db,
		logger: logger.New("repo/Inquiry"),
	}
}

type PaginateInquiryParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	OrderReferenceID       string                     `json:"order_reference_id" query:"order_reference_id" form:"order_reference_id"`
	Statuses               []enums.InquiryStatus      `json:"statuses" query:"statuses" form:"statuses"`
	BuyerQuotationStatuses []enums.InquiryBuyerStatus `json:"buyer_quotation_statuses" query:"buyer_quotation_statuses" form:"buyer_quotation_statuses"`
	UserID                 string                     `json:"user_id" query:"user_id" form:"user_id"`
	TeamID                 string                     `json:"team_id" query:"team_id" form:"team_id"`

	AssigneeID  string   `json:"assignee_id" query:"assignee_id" form:"assignee_id"`
	AssigneeIDs []string `json:"assignee_ids" query:"assignee_ids" form:"assignee_ids"`

	OwnerIDs []string `json:"owner_ids" query:"owner_ids" form:"owner_ids"`

	PostedDateFrom         int64 `json:"posted_date_from" query:"posted_date_from"`
	PostedDateTo           int64 `json:"posted_date_to" query:"posted_date_to"`
	IsQueryAll             bool  `json:"-"`
	PotentialOverdue       bool  `json:"-"`
	IncludePurchaseOrder   bool  `json:"-"`
	IncludeUser            bool  `json:"-"`
	IncludeAssignee        bool  `json:"-"`
	IncludeAuditLog        bool  `json:"-"`
	IncludeShippingAddress bool  `json:"-"`
	IncludeCollection      bool  `json:"-"`
}

func (r *InquiryRepo) PaginateInquiry(params PaginateInquiryParams) *query.Pagination {
	var userID = params.GetUserID()
	if params.TeamID != "" && !params.GetRole().IsAdmin() && !params.IsQueryAll {
		if err := r.db.Select("ID").First(&models.BrandTeam{}, "team_id = ? AND user_id = ?", params.TeamID, userID).Error; err != nil {
			return &query.Pagination{
				Records: []*models.PurchaseOrder{},
			}
		}
		userID = params.TeamID
	}

	var builder = queryfunc.NewInquiryBuilder(queryfunc.InquiryBuilderOptions{
		IncludePurchaseOrder:   params.IncludePurchaseOrder,
		IncludeAssignee:        params.IncludeAssignee,
		IncludeUser:            params.IncludeUser,
		IncludeAuditLog:        params.IncludeAuditLog,
		IncludeShippingAddress: params.IncludeShippingAddress,
		IncludeCollection:      params.IncludeCollection,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if !params.IsQueryAll {
				if params.GetRole().IsAdmin() {
					if params.UserID != "" {
						builder.Where("iq.user_id = ?", params.UserID)
					}
				} else {
					builder.Where("iq.user_id = ? AND iq.deleted_at IS NULL", userID)
				}
			}

			if params.PotentialOverdue {
				builder.Where("EXTRACT(EPOCH FROM now()) - iq.created_at >= ?", int(time.Hour.Seconds()*12))
				builder.Where("iq.status = ? AND iq.buyer_quotation_status = ? AND iq.deleted_at IS NULL", enums.InquiryStatusNew, enums.InquiryBuyerStatusNew)
			}

			if params.PostedDateFrom > 0 {
				builder.Where("iq.created_at >= ?", params.PostedDateFrom)
			}

			if params.PostedDateTo > 0 {
				builder.Where("iq.created_at <= ?", params.PostedDateTo)
			}

			if len(params.Statuses) > 0 {
				builder.Where("iq.status IN ?", params.Statuses)
			}

			if len(params.OwnerIDs) > 0 {
				builder.Where("iq.user_id IN ?", params.OwnerIDs)
			}

			if len(params.BuyerQuotationStatuses) > 0 {
				builder.Where("iq.buyer_quotation_status IN ?", params.BuyerQuotationStatuses)
			}

			if params.AssigneeID != "" {
				builder.Where("count_elements(iq.assignee_ids,?) >= 1", pq.StringArray([]string{params.AssigneeID}))
			}

			if len(params.AssigneeIDs) > 0 {
				builder.Where("count_elements(iq.assignee_ids,?) > 0", pq.StringArray(params.AssigneeIDs))
			}

			if strings.TrimSpace(params.OrderReferenceID) != "" {
				var q = "%" + params.OrderReferenceID + "%"
				builder.Where("order_reference_id ILIKE @query_po", sql.Named("query_po", q))
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				if strings.HasPrefix(keyword, "IQ") {
					builder.Where("iq.reference_id ILIKE ?", q)
				} else {
					builder.Where("(iq.id ILIKE @keyword OR iq.title ILIKE @keyword OR iq.sku_note ILIKE @keyword)", sql.Named("keyword", q))

				}
			}
		}).
		OrderBy("iq.updated_at DESC, iq.order_group_id ASC").
		Page(params.Page).
		WithoutCount(params.WithoutCount).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *InquiryRepo) PaginateInquiryForCreateOrder(params PaginateInquiryParams) *query.Pagination {
	var builder = queryfunc.NewInquiryOrderBuilder(queryfunc.InquiryOrderBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("o.id IS NULL")
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *InquiryRepo) AdminCreateInquiry(form models.InquiryAdminCreateForm) (*models.Inquiry, error) {
	var userAdmin models.User
	var err = r.db.Select("ID", "Role").First(&userAdmin, "id = ?", form.GetUserID()).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	if !userAdmin.Role.IsAdmin() {
		err = errs.ErrAdminActionOnly
		return nil, eris.Wrap(err, err.Error())
	}

	var userClient models.User
	err = r.db.Select("ID", "Role", "Name", "Avatar").First(&userClient, "id = ?", form.BuyerId).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	if userClient.Role != enums.RoleClient {
		err = errs.ErrClientActionOnly
		return nil, eris.Wrap(err, err.Error())
	}

	var inquiry models.Inquiry
	err = copier.Copy(&inquiry, &form)
	if err != nil {
		return nil, err
	}

	inquiry.UserID = form.BuyerId
	inquiry.CreatedBy = form.GetUserID()
	if form.AssigneeId != "" {
		var assignee models.User
		err = r.db.Select("ID", "Role").First(&assignee, "id = ?", form.AssigneeId).Error
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
		if !assignee.Role.IsAdmin() {
			err = errs.ErrAdminActionOnly
			return nil, eris.Wrap(err, err.Error())
		}
		inquiry.AssigneeIDs = []string{form.AssigneeId}
	}

	if inquiry.ShippingAddress != nil {
		if err = inquiry.ShippingAddress.CreateOrUpdate(r.db); err == nil {
			inquiry.ShippingAddressID = inquiry.ShippingAddress.ID
		}
	}

	var countryTax *models.SettingTax
	countryTax, err = NewSettingTaxRepo(r.db).GetAffectedSettingTax(models.GetAffectedSettingTaxForm{
		CurrencyCode: enums.VND,
	})
	if err == nil && countryTax != nil {
		inquiry.TaxPercentage = aws.Float64(countryTax.TaxPercentage)
	}

	err = r.db.Create(&inquiry).Error
	if err != nil {
		return nil, err
	}

	inquiry.User = &userClient

	// InquiryAudit
	err = NewInquiryAuditRepo(r.db).CreateInquiryAudit(models.InquiryAuditCreateForm{
		InquiryID:   inquiry.ID,
		ActionType:  enums.AuditActionTypeInquiryCreated,
		UserID:      userClient.ID,
		Description: fmt.Sprintf("%s has been created an inquiry: %s", userClient.Name, inquiry.ReferenceID),
	})
	if err != nil {
		return nil, err
	}

	return &inquiry, nil
}

type GetInquiryByIDParams struct {
	models.JwtClaimsInfo
	queryfunc.InquiryBuilderOptions

	InquiryID string `param:"inquiry_id" validate:"required"`
	UserID    string `json:"-"`
}

func (r *InquiryRepo) GetInquiryByID(params GetInquiryByIDParams) (*models.Inquiry, error) {
	params.InquiryBuilderOptions.Role = params.GetRole()
	params.IncludeUser = true
	var builder = queryfunc.NewInquiryBuilder(params.InquiryBuilderOptions)
	var Inquiry models.Inquiry
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if strings.HasPrefix(params.InquiryID, "IQ-") {
				builder.Where("iq.reference_id = ?", params.InquiryID)
			} else {
				builder.Where("iq.id = ?", params.InquiryID)
			}

			if params.UserID != "" {
				builder.Where("iq.user_id = ?", params.UserID)
			}
			if params.Role.IsBuyer() {
				builder.Where("iq.status != ?", enums.InquiryStatusClosed)
			}
		}).
		FirstFunc(&Inquiry)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &Inquiry, nil
}

func (r *InquiryRepo) UpdateInquiryByID(form models.InquiryUpdateForm) (*models.Inquiry, error) {
	var update models.Inquiry

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	currentInquiry, err := r.GetInquiryByID(GetInquiryByIDParams{InquiryID: form.InquiryID})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	if !form.GetRole().IsAdmin() && currentInquiry.EditTimeout != nil && *currentInquiry.EditTimeout < time.Now().Unix() {
		err = errs.ErrInquiryEditTimeout
		return nil, eris.Wrap(err, err.Error())
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Omit(clause.Associations).
			Set(models.InquiryAuditCurrentInquiryKey, currentInquiry).
			Set(models.InquiryAuditCurrentUserKey, func() interface{} {
				var user models.User
				var err = tx.Select("ID", "Name", "Avatar").First(&user, "id = ?", form.GetUserID()).Error
				if err == nil && user.ID != "" {
					return &user
				}
				return nil
			}()).
			Model(&models.Inquiry{}).
			Where("id = ?", form.InquiryID).
			Updates(&update).Error
		if err != nil {
			return err
		}

		var purchaseOrderUpdates = models.PurchaseOrder{
			Currency: form.Currency,
		}
		err = tx.Model(&models.PurchaseOrder{}).Where("inquiry_id = ?", form.InquiryID).Updates(&purchaseOrderUpdates).Error

		return err
	})
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(currentInquiry, &update, copier.Option{DeepCopy: true, IgnoreEmpty: true})

	return currentInquiry, err
}

func (r *InquiryRepo) UpdateInquiryEditTimeoutByID(form models.InquiryEditTimeoutUpdateForm) (err error) {
	var updates = models.Inquiry{
		EditTimeout:           aws.Int64(form.EditTimeout),
		IsEditTimeoutExtended: aws.Bool(true),
	}
	err = r.db.Model(&models.Inquiry{}).Where("id = ?", form.InquiryID).Updates(&updates).Error
	return err
}

type SendInquiryToSellerParams struct {
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`

	Sellers []*models.SellerRequestQuotationInfo `json:"sellers" validate:"required"`
}

func (r *InquiryRepo) SendInquiryToSeller(form SendInquiryToSellerParams) ([]*models.InquirySeller, error) {
	var inquiry models.Inquiry
	var err = r.db.Select("ID", "Currency").First(&inquiry, "id = ?", form.InquiryID).Error
	if err != nil {
		return nil, err
	}

	var records = lo.Map(form.Sellers, func(seller *models.SellerRequestQuotationInfo, index int) *models.InquirySeller {
		var record = &models.InquirySeller{
			InquiryID:                   inquiry.ID,
			Currency:                    inquiry.Currency,
			UserID:                      seller.SellerID,
			Status:                      enums.InquirySellerStatusNew,
			VarianceAmount:              seller.VarianceAmount,
			VariancePercentage:          seller.VariancePercentage,
			AdminSentAt:                 values.Int64(time.Now().Unix()),
			OfferPrice:                  seller.OfferPrice,
			OfferRemark:                 seller.OfferRemark,
			OrderType:                   seller.OrderType,
			ExpectedStartProductionDate: seller.ExpectedStartProductionDate,
		}

		return record
	})
	err = r.db.Clauses(clause.OnConflict{
		DoNothing: true,
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "inquiry_id"},
			{Name: "purchase_order_id"},
		},
	}).
		Create(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

type UpdateInquiryCostingParams struct {
	models.JwtClaimsInfo

	InquiryID  string      `param:"inquiry_id" validate:"required"`
	CMCost     price.Price `json:"cm_cost" param:"cm_cost"`
	FabricCost price.Price `json:"fabric_cost" param:"fabric_cost"`
	BasicCost  price.Price `json:"basic_cost" param:"basic_cost"`

	EstTACost    price.Price `json:"est_ta_cost" param:"est_ta_cost"`
	EstOtherCost price.Price `json:"est_other_cost" param:"est_other_cost"`
	EstTotalCost price.Price `json:"est_total_cost" param:"est_total_cost"`
}

func (r *InquiryRepo) UpdateInquiryCosting(params UpdateInquiryCostingParams) (*models.Inquiry, error) {
	var inquiry models.Inquiry
	var err = r.db.Select("ID", "UserID", "CMCost", "FabricCost", "BasicCost", "EstTACost", "EstOtherCost", "EstTotalCost").First(&inquiry, "id = ?", params.InquiryID).Error
	if err != nil {
		return nil, err
	}

	var admin models.User
	err = r.db.Select("ID", "Name", "Email").First(&admin, "id = ?", params.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	var updates models.Inquiry
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	var audit = models.InquiryAudit{
		InquiryID:   inquiry.ID,
		UserID:      admin.ID,
		ActionType:  enums.AuditActionTypeInquiryAdminUpdateCosting,
		Description: fmt.Sprintf("%s has updated costing", admin.Name),
		Metadata: &models.InquiryAuditMetadata{
			Before: map[string]interface{}{
				"cm_cost":        inquiry.CMCost,
				"fabric_cost":    inquiry.FabricCost,
				"basic_cost":     inquiry.BasicCost,
				"est_ta_cost":    inquiry.EstTACost,
				"est_other_cost": inquiry.EstOtherCost,
				"est_total_cost": inquiry.EstTotalCost,
			},
			After: map[string]interface{}{
				"cm_cost":        updates.CMCost,
				"fabric_cost":    updates.FabricCost,
				"basic_cost":     updates.BasicCost,
				"est_ta_cost":    updates.EstTACost,
				"est_other_cost": updates.EstOtherCost,
				"est_total_cost": updates.EstTotalCost,
				"refunded_by": map[string]interface{}{
					"id":    admin.ID,
					"name":  admin.Name,
					"email": admin.Email,
				},
			},
		},
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Clauses(clause.OnConflict{DoNothing: true}).Where("id = ?", params.InquiryID).Updates(&updates).Error
		if err != nil {
			return err
		}

		err = tx.Create(&audit).Error
		return err
	})
	if err != nil {
		return nil, err
	}

	return &inquiry, nil
}

type SendInquiryToBuyerResponse struct {
	Inquiry *models.Inquiry
	Admin   *models.User
}

func (r *InquiryRepo) SendInquiryToBuyer(form models.SendInquiryToBuyerForm) (*SendInquiryToBuyerResponse, error) {
	var admin models.User
	var err = r.db.Select("ID", "Name", "Email").First(&admin, "id = ?", form.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	inquiry, err := NewInquiryRepo(r.db).GetInquiryByID(GetInquiryByIDParams{
		InquiryID:     form.InquiryID,
		JwtClaimsInfo: form.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	validStatus := []enums.InquirySkuStatus{enums.InquirySkuStatusNew, enums.InquirySkuStatusWaitingForApproval, enums.InquirySkuStatusRejected}
	if ok := slices.Contains(validStatus, inquiry.BuyerQuotationStatus); !ok {
		return nil, errs.ErrInquiryInvalidToSendQuotationToBuyer
	}

	var user models.User
	err = r.db.Select("ID", "Name", "Email").First(&user, "id = ?", inquiry.UserID).Error
	if err != nil {
		return nil, err
	}

	var inquiryUpdates = &models.Inquiry{
		QuotationAt:          values.Int64(time.Now().Unix()),
		AdminQuotations:      form.Quotations,
		Status:               enums.InquiryStatusQuoteInProcess,
		BuyerQuotationStatus: enums.InquirySkuStatusWaitingForApproval,
		ShippingFee:          form.ShippingFee,
		TaxPercentage:        form.TaxPercentage,
		ProductWeight:        form.ProductWeight,
		AssigneeIDs:          lo.Union(append(inquiry.AssigneeIDs, form.GetUserID())),
	}

	err = r.db.Clauses(clause.OnConflict{UpdateAll: true}).Where("id = ?", inquiry.ID).Updates(inquiryUpdates).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	inquiry.TaxPercentage = form.TaxPercentage
	inquiry.User = &user
	inquiry.QuotationAt = inquiryUpdates.QuotationAt
	inquiry.AdminQuotations = inquiryUpdates.AdminQuotations
	inquiry.BuyerQuotationStatus = inquiryUpdates.BuyerQuotationStatus
	inquiry.ShippingFee = inquiryUpdates.ShippingFee
	inquiry.ProductWeight = inquiryUpdates.ProductWeight

	_, _ = r.InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
		InquiryID:     inquiry.ID,
		Inquiry:       inquiry,
		JwtClaimsInfo: form.JwtClaimsInfo,
	})

	var resp = SendInquiryToBuyerResponse{
		Inquiry: inquiry,
		Admin:   &admin,
	}

	return &resp, err
}

type AdminSubmitQuotationResponse struct {
	Inquiry *models.Inquiry
	Admin   *models.User
}

func (r *InquiryRepo) AdminSubmitQuotation(form models.SendInquiryToBuyerForm) (*AdminSubmitQuotationResponse, error) {
	var admin models.User
	var err = r.db.Select("ID", "Name", "Email").First(&admin, "id = ?", form.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	inquiry, err := NewInquiryRepo(r.db).GetInquiryByID(GetInquiryByIDParams{
		InquiryID:     form.InquiryID,
		JwtClaimsInfo: form.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	validStatus := []enums.InquirySkuStatus{enums.InquirySkuStatusNew, enums.InquirySkuStatusWaitingForApproval, enums.InquirySkuStatusRejected}
	if ok := slices.Contains(validStatus, inquiry.BuyerQuotationStatus); !ok {
		return nil, errs.ErrInquiryInvalidToSendQuotationToBuyer
	}

	var user models.User
	err = r.db.Select("ID", "Name", "Email").First(&user, "id = ?", inquiry.UserID).Error
	if err != nil {
		return nil, err
	}

	var inquiryUpdates = &models.Inquiry{
		InternalQuotationCreatedAt:  values.Int64(time.Now().Unix()),
		InternalQuotationCreatedBy:  admin.ID,
		InternalQuotationApprovedBy: "",
		InternalQuotationApprovedAt: nil,
		AdminQuotations:             form.Quotations,
		Status:                      enums.InquiryStatusQuoteInProcess,
		ShippingFee:                 form.ShippingFee.ToPtr(),
		TaxPercentage:               form.TaxPercentage,
		ProductWeight:               form.ProductWeight,
		AssigneeIDs:                 lo.Union(append(inquiry.AssigneeIDs, form.GetUserID())),
	}

	err = r.db.Clauses(clause.OnConflict{UpdateAll: true}).Where("id = ?", inquiry.ID).Updates(inquiryUpdates).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var resp = AdminSubmitQuotationResponse{
		Inquiry: inquiry,
		Admin:   &admin,
	}

	return &resp, err
}

func (r *InquiryRepo) AdminSubmitMultipleInquiryQuotations(req *models.SubmitMultipleInquiryQuotationRequest) ([]*models.Inquiry, error) {
	var inquiryIDs = make([]string, 0, len(req.Quotations))
	for _, quotation := range req.Quotations {
		inquiryIDs = append(inquiryIDs, quotation.InquiryID)
	}
	var inquiries models.Inquiries
	if err := r.db.Select("ID", "BuyerQuotationStatus", "AssigneeIDs", "UserID").Find(&inquiries, "id IN ?", inquiryIDs).Error; err != nil {
		return nil, err
	}
	var dbInquiryIDs = inquiries.IDs()
	for _, id := range inquiryIDs {
		if !helper.StringContains(dbInquiryIDs, id) {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "inquiry_id:%s", id)
		}
	}

	var users models.Users
	if err := r.db.Find(&users, "id IN ?", inquiries.UserIDs()).Error; err != nil {
		return nil, err
	}
	var mapUserIDToUser = make(map[string]*models.User, len(users))
	for _, user := range users {
		mapUserIDToUser[user.ID] = user
	}

	validStatus := []enums.InquirySkuStatus{enums.InquirySkuStatusNew, enums.InquirySkuStatusWaitingForApproval, enums.InquirySkuStatusRejected}
	for idx, iq := range inquiries {
		if ok := slices.Contains(validStatus, iq.BuyerQuotationStatus); !ok {
			return nil, eris.Wrapf(errs.ErrInquiryInvalidToSendQuotationToBuyer, "index:%d", idx)
		}
	}
	var inquiriesToUpdate []*models.Inquiry

	for _, quotation := range req.Quotations {
		inquiry, _ := lo.Find(inquiries, func(item *models.Inquiry) bool {
			return item.ID == quotation.InquiryID
		})
		var inquiryUpdate = &models.Inquiry{
			Model:                       models.Model{ID: inquiry.ID},
			InternalQuotationCreatedAt:  values.Int64(time.Now().Unix()),
			InternalQuotationCreatedBy:  req.GetUserID(),
			InternalQuotationApprovedBy: "",
			InternalQuotationApprovedAt: nil,
			AdminQuotations:             quotation.Quotations,
			Status:                      enums.InquiryStatusQuoteInProcess,
			BuyerQuotationStatus:        enums.InquirySkuStatusWaitingForApproval,
			ShippingFee:                 quotation.ShippingFee.ToPtr(),
			TaxPercentage:               quotation.TaxPercentage,
			ProductWeight:               quotation.ProductWeight,
			AssigneeIDs:                 lo.Union(append(inquiry.AssigneeIDs, req.GetUserID())),
		}
		// inject additional attributes to return
		user := mapUserIDToUser[inquiry.UserID]
		inquiryUpdate.User = user
		inquiriesToUpdate = append(inquiriesToUpdate, inquiryUpdate)
	}
	if err := r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"internal_quotation_created_at",
				"internal_quotation_created_by",
				"internal_quotation_approved_at",
				"internal_quotation_approved_by",
				"admin_quotations",
				"status",
				"buyer_quotation_status",
				"shipping_fee",
				"tax_percentage",
				"product_weight",
				"assignee_ids",
			})}).Create(inquiriesToUpdate).Error; err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return inquiriesToUpdate, nil
}

type ApproveInquiryRequestParams struct {
	models.JwtClaimsInfo

	InquiryID       string `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`
	InquirySellerID string `json:"inquiry_seller_id" param:"inquiry_seller_id" query:"inquiry_request_id" validate:"required"`
}

func (r *InquiryRepo) AdminInternalApproveQuotation(form models.AdminInternalApproveQuotationForm) (*SendInquiryToBuyerResponse, error) {
	var admin models.User
	var err = r.db.Select("ID", "Name", "Email").First(&admin, "id = ?", form.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	inquiry, err := NewInquiryRepo(r.db).GetInquiryByID(GetInquiryByIDParams{
		InquiryID:     form.InquiryID,
		JwtClaimsInfo: form.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	validStatus := []enums.InquirySkuStatus{enums.InquirySkuStatusNew, enums.InquirySkuStatusWaitingForApproval, enums.InquirySkuStatusRejected}
	if ok := slices.Contains(validStatus, inquiry.BuyerQuotationStatus); !ok {
		return nil, errs.ErrInquiryInvalidToSendQuotationToBuyer
	}
	if inquiry.InternalQuotationCreatedAt == nil || inquiry.InternalQuotationApprovedBy != "" {
		return nil, errs.ErrInquiryInvalidToSendQuotationToBuyer
	}

	var user models.User
	err = r.db.Select("ID", "Name", "Email").First(&user, "id = ?", inquiry.UserID).Error
	if err != nil {
		return nil, err
	}

	var inquiryUpdates = &models.Inquiry{
		QuotationAt:                 values.Int64(time.Now().Unix()),
		InternalQuotationApprovedBy: admin.ID,
		InternalQuotationApprovedAt: values.Int64(time.Now().Unix()),
		Status:                      enums.InquiryStatusQuoteInProcess,
		BuyerQuotationStatus:        enums.InquirySkuStatusWaitingForApproval,
	}

	err = r.db.Clauses(clause.OnConflict{UpdateAll: true}).Where("id = ?", inquiry.ID).Updates(inquiryUpdates).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	inquiry.User = &user
	inquiry.QuotationAt = inquiryUpdates.QuotationAt
	inquiry.AdminQuotations = inquiryUpdates.AdminQuotations
	inquiry.BuyerQuotationStatus = inquiryUpdates.BuyerQuotationStatus
	inquiry.ShippingFee = inquiryUpdates.ShippingFee
	inquiry.ProductWeight = inquiryUpdates.ProductWeight

	_, _ = r.InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
		InquiryID:     inquiry.ID,
		Inquiry:       inquiry,
		JwtClaimsInfo: form.JwtClaimsInfo,
	})

	var resp = SendInquiryToBuyerResponse{
		Inquiry: inquiry,
		Admin:   &admin,
	}

	err = NewInquiryAuditRepo(r.db).CreateInquiryAudit(models.InquiryAuditCreateForm{
		InquiryID:   inquiry.ID,
		ActionType:  enums.AuditActionTypeInquiryAdminSendBuyerQuotation,
		UserID:      admin.ID,
		Description: fmt.Sprintf("Admin %s has sent quotation to buyer %s", admin.Name, inquiry.User.Name),
		Metadata: &models.InquiryAuditMetadata{
			After: map[string]interface{}{
				"quotations": inquiry.AdminQuotations,
			},
		},
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	return &resp, err
}

type PaginateInquiryAuditsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	InquiryID   string                  `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" form:"inquiry_id" validate:"required"`
	DateFrom    int64                   `json:"date_from" query:"date_from" form:"date_from"`
	DateTo      int64                   `json:"date_to" query:"date_to" form:"date_to"`
	UserID      string                  `json:"user_id" query:"user_id" form:"user_id"`
	ActionTypes []enums.AuditActionType `json:"action_types" query:"action_types" form:"action_types"`
}

func (r *InquiryRepo) PaginateInquiryAudits(params PaginateInquiryAuditsParams) *query.Pagination {
	var builder = queryfunc.NewInquiryAuditBuilder(queryfunc.InquiryAuditBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("ia.inquiry_id = ?", params.InquiryID)

			if params.UserID != "" {
				builder.Where("ia.user_id = ?", params.UserID)
			}
			if params.DateFrom != 0 {
				builder.Where("ia.created_at >= ?", values.Int64(params.DateFrom))
			}
			if params.DateTo != 0 {
				builder.Where("ia.created_at <= ?", values.Int64(params.DateTo))
			}

			if len(params.ActionTypes) > 0 {
				builder.Where("ia.action_type IN ?", params.ActionTypes)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *InquiryRepo) AdminMarkSeen(form models.InquiryMarkSeenForm) error {

	var update = &models.Inquiry{
		NewSeenAt:    values.Int64(time.Now().Unix()),
		UpdateSeenAt: values.Int64(time.Now().Unix()),
	}

	var err = r.db.Model(&models.Inquiry{}).
		Where("id = ?", form.InquiryID).
		Updates(&update).Error

	return err
}

func (r *InquiryRepo) CloneInquiry(params models.InquiryIDParam) (*models.Inquiry, error) {
	var inquiry models.Inquiry
	var err error

	if params.GetRole().IsAdmin() {
		err = r.db.First(&inquiry, "id = ?", params.InquiryID).Error
	} else {
		err = r.db.First(&inquiry, "id = ? AND user_id = ?", params.InquiryID, params.GetUserID()).Error
	}
	if err != nil {
		return nil, err
	}

	inquiry.ID = helper.GenerateXID()
	inquiry.Status = enums.InquiryStatusNew
	inquiry.UpdateSeenAt = nil
	inquiry.NewSeenAt = nil
	inquiry.ReferenceID = helper.GenerateInquiryReferenceID()
	inquiry.DeliveryDate = nil
	inquiry.BuyerQuotationStatus = enums.InquirySkuStatusNew

	err = r.db.Create(&inquiry).Error
	if err != nil {
		return nil, err
	}

	return &inquiry, err
}

func (r *InquiryRepo) CloneInquiryAndQuotation(params models.InquiryIDParam) (*models.Inquiry, error) {
	var inquiry models.Inquiry
	var err error

	if params.GetRole().IsAdmin() {
		err = r.db.First(&inquiry, "id = ?", params.InquiryID).Error
	} else {
		err = r.db.First(&inquiry, "id = ? AND user_id = ?", params.InquiryID, params.GetUserID()).Error
	}
	if err != nil {
		return nil, err
	}

	inquiry.ID = helper.GenerateXID()
	inquiry.Status = enums.InquiryStatusQuoteInProcess
	inquiry.UpdateSeenAt = nil
	inquiry.NewSeenAt = nil
	inquiry.ReferenceID = helper.GenerateInquiryReferenceID()
	inquiry.DeliveryDate = nil
	inquiry.BuyerQuotationStatus = enums.InquirySkuStatusWaitingForApproval
	inquiry.QuotationAt = values.Int64(time.Now().Unix())

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&inquiry).Error
		if err != nil {
			return err
		}
		return err
	})

	return &inquiry, err
}

type PaginateInquiryCollectionParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	UserID string `json:"user_id" query:"user_id" form:"user_id"`
}

func (r *InquiryRepo) PaginateInquiryCollections(params PaginateInquiryCollectionParams) *query.Pagination {
	var builder = queryfunc.NewInquiryCollectionBuilder(queryfunc.InquiryCollectionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.UserID != "" {
				builder.Where("c.user_id = ?", params.UserID)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *InquiryRepo) CreateInquiryCollection(form models.InquiryCollectionUpdateForm) (*models.InquiryCollection, error) {
	var inquiryCollection = models.InquiryCollection{}
	err := copier.Copy(&inquiryCollection, &form)
	if err != nil {
		return nil, err
	}
	err = r.db.Omit(clause.Associations).
		Where(models.InquiryCollection{UserID: form.UserID, Name: form.Name}).
		Create(&inquiryCollection).Error
	if err != nil {
		if duplicated, _ := r.db.IsDuplicateConstraint(err); duplicated {
			return nil, errs.ErrInquiryCollectionTaken
		}
		return nil, eris.Wrap(err, err.Error())
	}
	return &inquiryCollection, nil
}

func (r *InquiryRepo) UpdateCartItems(form models.InquiryCartItemsUpdateForm) ([]*models.OrderCartItem, error) {
	var inquiry models.Inquiry
	var err = r.db.First(&inquiry, "id = ? AND user_id = ?", form.InquiryID, form.GetUserID()).Error
	if err != nil {
		return nil, err
	}
	var purchaseOrder models.PurchaseOrder
	if err := r.db.Select("ID").First(&purchaseOrder, "inquiry_id = ? AND user_id = ?", form.InquiryID, form.GetUserID()).Error; err != nil {
		return nil, err
	}

	var itemsToCreate []*models.OrderCartItem
	for _, item := range form.Items {
		var unitPrice = inquiry.GetSampleUnitPrice()
		var totalPrice = unitPrice.MultipleInt(int64(item.Qty))
		itemsToCreate = append(itemsToCreate, &models.OrderCartItem{
			PurchaseOrderID: purchaseOrder.ID,
			UnitPrice:       unitPrice,
			TotalPrice:      totalPrice,
			ColorName:       item.ColorName,
			Size:            item.Size,
			Qty:             item.Qty,
			NoteToSupplier:  item.NoteToSupplier,
			Style:           item.Style,
		})
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Delete(&models.OrderCartItem{}, "purchase_order_id = ?", purchaseOrder.ID).Error; err != nil {
			return err
		}

		if len(itemsToCreate) > 0 {
			err = tx.Create(&itemsToCreate).Error
			if err != nil {
				return err
			}
		}

		var purchaseOrder = models.PurchaseOrder{
			UserID:            inquiry.UserID,
			InquiryID:         inquiry.ID,
			Quotations:        inquiry.AdminQuotations,
			ProductWeight:     inquiry.ProductWeight,
			ShippingAddressID: inquiry.ShippingAddressID,
			Currency:          inquiry.Currency,
			OrderGroupID:      inquiry.OrderGroupID,
		}
		purchaseOrder.TaxPercentage = inquiry.TaxPercentage
		purchaseOrder.ShippingFee = inquiry.ShippingFee

		err = tx.Clauses(clause.OnConflict{UpdateAll: true}).Where("inquiry_id = ?", inquiry.ID).FirstOrCreate(&purchaseOrder).Error

		return err
	})
	if err != nil {
		return nil, err
	}

	return itemsToCreate, nil
}

func (r *InquiryRepo) GetCartItems(form models.GetInquiryCartItemsParams) ([]*models.OrderCartItem, error) {
	var items []*models.OrderCartItem
	var inquiry models.Inquiry
	var err error
	err = r.db.Select("ID").First(&inquiry, "id = ? AND user_id = ?", form.InquiryID, form.GetUserID()).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrInquiryNotFound
		}
		return nil, err
	}
	var purchaseOrder models.PurchaseOrder
	if err := r.db.Select("ID").First(&purchaseOrder, "inquiry_id = ? AND user_id = ?", inquiry.ID, form.GetUserID()).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrInquiryNotFound
		}
		return nil, err

	}

	err = r.db.Find(&items, "purchase_order_id = ?", purchaseOrder.ID).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

type PaginateCartsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *InquiryRepo) PaginateCarts(params PaginateCartsParams) *query.Pagination {
	var builder = queryfunc.NewInquiryCartBuilder(queryfunc.InquiryCartBuilderOptions{})

	return query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("po.status IN ?", []string{string(enums.PurchaseOrderStatusPending)})
			builder.Where("iq.user_id = ?", params.GetUserID())
			builder.Where("c.id IS NOT NULL")
		}).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
}

type InquiryCheckoutParams struct {
	models.JwtClaimsInfo

	InquiryID         string            `param:"inquiry_id" validate:"required"`
	CheckoutSessionID string            `param:"checkout_session_id"` // only for multi inquiry in cart
	PaymentType       enums.PaymentType `json:"payment_type" validate:"oneof=bank_transfer card"`
	PaymentMethodID   string            `json:"payment_method_id" validate:"required_if=PaymentType card"`

	TransactionRefID      string             `json:"transaction_ref_id" validate:"required_if=PaymentType bank_transfer"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required_if=PaymentType bank_transfer"`
}

func (r *InquiryRepo) InquiryCheckout(params InquiryCheckoutParams) (*models.PurchaseOrder, error) {
	purchaseOrder, err := r.InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		InquiryID:     params.InquiryID,
		PaymentType:   params.PaymentType,
		UserID:        params.GetUserID(),
	})
	if err != nil {
		return nil, err
	}

	if purchaseOrder.PaymentType == enums.PaymentTypeBankTransfer {
		r.db.Transaction(func(tx *gorm.DB) error {

			// create transaction
			var transaction = models.PaymentTransaction{
				PurchaseOrderID:   purchaseOrder.ID,
				PaidAmount:        purchaseOrder.TotalPrice,
				PaymentType:       params.PaymentType,
				Milestone:         enums.PaymentMilestoneFinalPayment,
				UserID:            purchaseOrder.UserID,
				TransactionRefID:  params.TransactionRefID,
				Status:            enums.PaymentStatusWaitingConfirm,
				PaymentPercentage: values.Float64(100),
				TotalAmount:       purchaseOrder.TotalPrice,
				Currency:          purchaseOrder.Inquiry.Currency,
				Metadata: &models.PaymentTransactionMetadata{
					InquiryID:                purchaseOrder.Inquiry.ID,
					InquiryReferenceID:       purchaseOrder.Inquiry.ReferenceID,
					PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
					PurchaseOrderID:          purchaseOrder.ID,
				},
			}
			if params.TransactionAttachment != nil {
				transaction.Attachments = &models.Attachments{params.TransactionAttachment}
			}
			err = tx.Create(&transaction).Error
			if err != nil {
				return err
			}

			var updates = models.PurchaseOrder{
				Status:                        enums.PurchaseOrderStatusWaitingConfirm,
				PaymentType:                   params.PaymentType,
				TransactionRefID:              params.TransactionRefID,
				TransactionAttachment:         params.TransactionAttachment,
				TransferedAt:                  values.Int64(time.Now().Unix()),
				Currency:                      purchaseOrder.Inquiry.Currency,
				PaymentTransactionReferenceID: transaction.ReferenceID,
			}
			updates.TaxPercentage = purchaseOrder.Inquiry.TaxPercentage
			var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
			if sqlResult.Error != nil {
				return eris.Wrap(sqlResult.Error, sqlResult.Error.Error())
			}

			if sqlResult.RowsAffected == 0 {
				return eris.New("Purchase order not found")
			}

			return err

		})

		return purchaseOrder, err
	}

	var user models.User
	err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", params.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	stripeConfig, err := stripehelper.GetCurrencyConfig(purchaseOrder.Inquiry.Currency)
	if err != nil {
		return nil, err
	}

	var stripeParams = stripehelper.CreatePaymentIntentParams{
		Amount:                  purchaseOrder.TotalPrice.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64(),
		Currency:                purchaseOrder.Inquiry.Currency,
		PaymentMethodID:         params.PaymentMethodID,
		CustomerID:              user.StripeCustomerID,
		IsCaptureMethodManually: false,
		Description:             fmt.Sprintf("Charges for %s/%s", purchaseOrder.ReferenceID, purchaseOrder.Inquiry.ReferenceID),
		PaymentMethodTypes:      []string{"card"},
		Metadata: map[string]string{
			"inquiry_id":                  purchaseOrder.InquiryID,
			"inquiry_reference_id":        purchaseOrder.Inquiry.ReferenceID,
			"purchase_order_id":           purchaseOrder.ID,
			"purchase_order_reference_id": purchaseOrder.ReferenceID,
			"cart_item_ids":               strings.Join(purchaseOrder.CartItemIDs, ","),
			"action_source":               string(stripehelper.ActionSourceInquiryPayment),
		},
	}

	if purchaseOrder.Inquiry.ShippingAddress != nil {
		stripeParams.Shipping = &stripe.ShippingDetailsParams{
			Name: &purchaseOrder.Inquiry.ShippingAddress.Name,
		}
		if purchaseOrder.Inquiry.ShippingAddress.Coordinate != nil {
			stripeParams.Shipping.Address = &stripe.AddressParams{
				State:      &purchaseOrder.Inquiry.ShippingAddress.Coordinate.Level1,
				City:       &purchaseOrder.Inquiry.ShippingAddress.Coordinate.Level2,
				PostalCode: &purchaseOrder.Inquiry.ShippingAddress.Coordinate.PostalCode,
				Country:    stripe.String(purchaseOrder.Inquiry.ShippingAddress.Coordinate.CountryCode.String()),
				Line1:      &purchaseOrder.Inquiry.ShippingAddress.Coordinate.FormattedAddress,
			}
		}

	}

	pi, err := stripehelper.GetInstance().CreatePaymentIntent(stripeParams)
	if err != nil {
		return nil, err
	}
	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		if pi.NextAction != nil {
			intent, err := stripehelper.GetInstance().ConfirmPaymentIntent(stripehelper.ConfirmPaymentIntentParams{
				PaymentIntentID: pi.ID,
				ReturnURL:       fmt.Sprintf("%s/api/v1/callback/stripe/payment_intents/inquiries/%s/purchase_orders/%s/confirm", r.db.Configuration.ServerBaseURL, purchaseOrder.InquiryID, purchaseOrder.ID),
			})
			if err != nil {
				return nil, err
			}

			if intent.Status == stripe.PaymentIntentStatusSucceeded {
				goto PaymentSuccess
			}

			purchaseOrder.PaymentIntentNextAction = intent.NextAction
			purchaseOrder.PaymentIntentClientSecret = intent.ClientSecret
			return purchaseOrder, nil
		} else {
			return nil, eris.Errorf("Payment error with status %s", pi.Status)
		}
	}

PaymentSuccess:
	var updates = models.PurchaseOrder{
		PaymentIntentID:       pi.ID,
		Status:                enums.PurchaseOrderStatusPaid,
		PaymentType:           params.PaymentType,
		MarkAsPaidAt:          values.Int64(time.Now().Unix()),
		TransactionRefID:      params.TransactionRefID,
		TransactionAttachment: params.TransactionAttachment,
		TransferedAt:          values.Int64(time.Now().Unix()),
		Currency:              purchaseOrder.Inquiry.Currency,
	}

	if len(purchaseOrder.Quotations) > 0 {
		sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeSample
		})
		if sampleQuotation != nil {
			updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
			updates.StartDate = updates.MarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	}

	updates.TaxPercentage = purchaseOrder.Inquiry.TaxPercentage
	err = r.db.Transaction(func(tx *gorm.DB) error {

		// create transaction
		var transaction = models.PaymentTransaction{
			PurchaseOrderID:   purchaseOrder.ID,
			PaidAmount:        purchaseOrder.TotalPrice,
			PaymentType:       purchaseOrder.PaymentType,
			UserID:            purchaseOrder.UserID,
			TransactionRefID:  purchaseOrder.TransactionRefID,
			PaymentIntentID:   purchaseOrder.PaymentIntentID,
			Status:            enums.PaymentStatusPaid,
			TotalAmount:       purchaseOrder.TotalPrice,
			Milestone:         enums.PaymentMilestoneFinalPayment,
			PaymentPercentage: values.Float64(100),
			MarkAsPaidAt:      values.Int64(time.Now().Unix()),
			Currency:          purchaseOrder.Inquiry.Currency,
			Metadata: &models.PaymentTransactionMetadata{
				InquiryID:                purchaseOrder.Inquiry.ID,
				InquiryReferenceID:       purchaseOrder.Inquiry.ReferenceID,
				PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
				PurchaseOrderID:          purchaseOrder.ID,
			},
		}
		if params.TransactionAttachment != nil {
			transaction.Attachments = &models.Attachments{params.TransactionAttachment}
		}

		err = tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		updates.PaymentTransactionReferenceID = transaction.ReferenceID
		var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return eris.New("Purchase order is not found")
		}

		return tx.Model(&models.Inquiry{}).Where("id = ?", purchaseOrder.InquiryID).UpdateColumn("Status", enums.InquiryStatusFinished).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	purchaseOrder.PaymentIntentID = updates.PaymentIntentID
	purchaseOrder.Status = updates.Status
	purchaseOrder.PaymentType = updates.PaymentType
	purchaseOrder.MarkAsPaidAt = updates.MarkAsPaidAt
	purchaseOrder.TaxPercentage = updates.TaxPercentage
	purchaseOrder.TransactionRefID = updates.TransactionRefID
	purchaseOrder.TransactionAttachment = updates.TransactionAttachment
	purchaseOrder.TransferedAt = updates.TransferedAt

	return purchaseOrder, err

}

type InquiryPreviewCheckoutParams struct {
	models.JwtClaimsInfo

	InquiryID   string                    `param:"inquiry_id" validate:"required"`
	PaymentType enums.PaymentType         `json:"payment_type" validate:"oneof=bank_transfer card"`
	CartItems   []*models.InquiryCartItem `json:"-"`
	Inquiry     *models.Inquiry           `json:"-"`
	UserID      string                    `param:"user_id"`

	UpdatePricing bool `json:"-"`
}

func (r *InquiryRepo) InquiryPreviewCheckout(params InquiryPreviewCheckoutParams) (*models.PurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("inquiry_%s", params.InquiryID), time.Second*30)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var inquiry = params.Inquiry
	if inquiry == nil {
		if params.GetRole().IsAdmin() {
			err = r.db.Select("ID", "ReferenceID", "UserID", "ShippingFee", "ProductWeight", "Currency", "ShippingAddressID", "TaxPercentage", "AssigneeIDs", "AdminQuotations").
				First(&inquiry, "id = ?", params.InquiryID).Error
		} else {
			err = r.db.Select("ID", "ReferenceID", "UserID", "ShippingFee", "ProductWeight", "Currency", "ShippingAddressID", "TaxPercentage", "AssigneeIDs", "AdminQuotations").
				First(&inquiry, "id = ? AND user_id = ?", params.InquiryID, params.UserID).Error
		}

	}

	if err != nil {
		return nil, err
	}

	if inquiry.ShippingAddressID != "" {
		inquiry.ShippingAddress, _ = NewAddressRepo(r.db).GetAddress(GetAddressParams{
			AddressID: inquiry.ShippingAddressID,
		})
	}

	var items = params.CartItems
	if len(items) == 0 {
		items = []*models.InquiryCartItem{}
		err = r.db.Find(&items, "inquiry_id = ?", params.InquiryID).Error
		if err != nil {
			return nil, err
		}
	}

	if len(items) == 0 {
		return nil, eris.Errorf("Empty items")
	}

	var subTotalPrice = price.NewFromFloat(0)
	var cartItemsIDs = lo.Map(items, func(item *models.InquiryCartItem, index int) string {
		subTotalPrice = subTotalPrice.Add(item.TotalPrice)
		return item.ID
	})

	var existingPurchaseOrder models.PurchaseOrder
	err = r.db.First(&existingPurchaseOrder, "user_id = ? AND inquiry_id = ?", inquiry.UserID, inquiry.ID).Error
	if err == nil && existingPurchaseOrder.ID != "" {
		if existingPurchaseOrder.Status == enums.PurchaseOrderStatusPaid && existingPurchaseOrder.SubTotal.GreaterThan(0) {
			existingPurchaseOrder.Inquiry = inquiry
			return &existingPurchaseOrder, err
		}

		existingPurchaseOrder.Currency = inquiry.Currency
		existingPurchaseOrder.TaxPercentage = inquiry.TaxPercentage
		existingPurchaseOrder.ShippingFee = inquiry.ShippingFee
		existingPurchaseOrder.PaymentType = params.PaymentType
		existingPurchaseOrder.CartItemIDs = cartItemsIDs
		existingPurchaseOrder.CartItems = items
		existingPurchaseOrder.SubTotal = subTotalPrice.ToPtr()
		existingPurchaseOrder.Inquiry = inquiry
		existingPurchaseOrder.ProductWeight = inquiry.ProductWeight
		existingPurchaseOrder.ShippingAddressID = inquiry.ShippingAddressID
		existingPurchaseOrder.Quotations = inquiry.AdminQuotations
		err = existingPurchaseOrder.UpdatePrices()
		if err != nil {
			return nil, err
		}

		if !params.GetRole().IsAdmin() || params.UpdatePricing {
			err = r.db.Omit(clause.Associations).Model(&models.PurchaseOrder{}).Where("id = ?", existingPurchaseOrder.ID).Updates(&existingPurchaseOrder).Error
			if err != nil {
				return nil, err
			}
		}

		return &existingPurchaseOrder, err
	}

	var purchaseOrder = models.PurchaseOrder{
		CartItemIDs:       cartItemsIDs,
		CartItems:         items,
		UserID:            inquiry.UserID,
		InquiryID:         inquiry.ID,
		Inquiry:           inquiry,
		PaymentType:       params.PaymentType,
		Currency:          inquiry.Currency,
		ShippingAddressID: inquiry.ShippingAddressID,
		Quotations:        inquiry.AdminQuotations,
	}
	purchaseOrder.ShippingFee = inquiry.ShippingFee
	purchaseOrder.SubTotal = subTotalPrice.ToPtr()
	purchaseOrder.Inquiry = inquiry
	purchaseOrder.TaxPercentage = inquiry.TaxPercentage
	purchaseOrder.ProductWeight = inquiry.ProductWeight
	purchaseOrder.Quotations = inquiry.AdminQuotations
	purchaseOrder.OrderGroupID = inquiry.OrderGroupID
	err = purchaseOrder.UpdatePrices()
	if err != nil {
		return nil, err
	}

	if !params.GetRole().IsAdmin() || params.UpdatePricing {
		// Create purchase order
		err = r.db.Omit(clause.Associations).
			Clauses(clause.OnConflict{UpdateAll: true}).
			Create(&purchaseOrder).Error
		if err != nil {
			if duplicated, _ := r.db.IsDuplicateConstraint(err); duplicated {
				purchaseOrder.ID = existingPurchaseOrder.ID
				purchaseOrder.ReferenceID = existingPurchaseOrder.ReferenceID
				purchaseOrder.Status = existingPurchaseOrder.Status
			} else {
				return nil, err
			}

		}
		if err != nil {
			return nil, err
		}
	}

	return &purchaseOrder, err
}

func (r *InquiryRepo) InquiryMarkAsPaid(params models.InquiryIDParam) (*models.PurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("inquiry_%s", params.InquiryID), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var purchaseOrder = params.PurchaseOrder
	if params.InquiryID != "" && purchaseOrder == nil {
		purchaseOrder = new(models.PurchaseOrder)
		err = r.db.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs").First(purchaseOrder, "inquiry_id = ?", params.InquiryID).Error
		if err != nil {
			return nil, err
		}

	}

	if params.PurchaseOrderID != "" && purchaseOrder == nil {
		purchaseOrder = new(models.PurchaseOrder)
		err = r.db.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs").First(purchaseOrder, "id = ?", params.PurchaseOrderID).Error
		if err != nil {
			return nil, err
		}
	}

	if purchaseOrder.ID == "" {
		return nil, errs.ErrInquiryNotFound
	}

	if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
		return nil, errs.ErrPoIsAlreadyPaid
	}

	var updates = models.PurchaseOrder{
		MarkAsPaidAt: values.Int64(time.Now().Unix()),
		Status:       enums.PurchaseOrderStatusPaid,
	}

	var paymentTransactionUpdates = models.PaymentTransaction{
		MarkAsPaidAt: updates.MarkAsPaidAt,
		Status:       enums.PaymentStatusPaid,
	}

	if len(purchaseOrder.Quotations) > 0 {
		sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeSample
		})
		if sampleQuotation != nil {
			updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
			updates.StartDate = updates.MarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	} else {
		var inquiry models.Inquiry
		r.db.Select("AdminQuotations").First(&inquiry, "id = ?", params.InquiryID)
		sampleQuotation, _ := lo.Find(inquiry.AdminQuotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeSample
		})
		if sampleQuotation != nil {
			updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
			updates.StartDate = updates.MarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var err = tx.Model(&updates).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Where("inquiry_id = ? AND payment_type = ?", params.InquiryID, enums.PaymentTypeBankTransfer).
			Updates(&updates).Error
		if err != nil {
			return err
		}

		return tx.Model(&models.PaymentTransaction{}).
			Where("purchase_order_id = ?", updates.ID).
			Updates(&paymentTransactionUpdates).Error
	})
	return purchaseOrder, err
}

func (r *InquiryRepo) InquiryMarkAsUnpaid(params models.InquiryIDParam) (*models.PurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("inquiry_%s", params.InquiryID), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var purchaseOrder models.PurchaseOrder

	if params.InquiryID != "" {
		err = r.db.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs").First(&purchaseOrder, "inquiry_id = ?", params.InquiryID).Error
		if err != nil {
			return nil, err
		}

	}

	if params.PurchaseOrderID != "" && purchaseOrder.ID == "" {
		err = r.db.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs").First(&purchaseOrder, "id = ?", params.PurchaseOrderID).Error
		if err != nil {
			return nil, err
		}
	}

	if purchaseOrder.ID == "" {
		return nil, errs.ErrInquiryNotFound
	}

	if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
		return nil, errs.ErrPoIsAlreadyPaid
	}

	err = r.db.Transaction(func(tx *gorm.DB) (e error) {
		var updates = models.PurchaseOrder{
			MarkAsUnpaidAt: values.Int64(time.Now().Unix()),
			Status:         enums.PurchaseOrderStatusUnpaid,
		}
		e = tx.Model(&updates).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Where("inquiry_id = ? AND payment_type = ?", params.InquiryID, enums.PaymentTypeBankTransfer).
			Updates(&updates).Error
		if e != nil {
			return
		}

		var paymentTrxUpdates = models.PaymentTransaction{
			Status:         enums.PaymentStatusUnpaid,
			Note:           params.Note,
			MarkAsUnpaidAt: updates.MarkAsUnpaidAt,
		}
		e = tx.Model(&paymentTrxUpdates).
			Where("purchase_order_id = ?", updates.ID).
			Updates(&paymentTrxUpdates).Error
		return
	})
	return &purchaseOrder, err
}

type MultiInquiryParams struct {
	models.JwtClaimsInfo

	CheckoutSessionID string `param:"checkout_session_id" validate:"required"`
	Note              string `json:"note"`
}

func (r *InquiryRepo) MultiInquiryMarkAsPaid(params MultiInquiryParams) ([]*models.PurchaseOrder, error) {
	var purchaseOrders []*models.PurchaseOrder

	if params.CheckoutSessionID != "" {
		var err = r.db.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs", "CheckoutSessionID").Find(&purchaseOrders, "checkout_session_id = ?", params.CheckoutSessionID).Error
		if err != nil {
			return nil, err
		}
	}

	if len(purchaseOrders) == 0 {
		return nil, errs.ErrInquiryNotFound
	}

	var paymentTransactionUpdates = models.PaymentTransaction{
		MarkAsPaidAt: values.Int64(time.Now().Unix()),
		Status:       enums.PaymentStatusPaid,
	}

	var err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, purchaseOrder := range purchaseOrders {
			if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
				return errs.ErrPoIsAlreadyPaid
			}

			var updates = models.PurchaseOrder{
				MarkAsPaidAt: paymentTransactionUpdates.MarkAsPaidAt,
				Status:       enums.PurchaseOrderStatusPaid,
			}

			if len(purchaseOrder.Quotations) > 0 {
				sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
					return item.Type == enums.InquiryTypeSample
				})
				if sampleQuotation != nil {
					updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
					updates.StartDate = updates.MarkAsPaidAt
					updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
				}
			}

			var err = tx.Model(&updates).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
				Where("checkout_session_id = ? AND payment_type = ?", params.CheckoutSessionID, enums.PaymentTypeBankTransfer).
				Updates(&updates).Error
			if err != nil {
				return err
			}
		}

		return tx.Model(&models.PaymentTransaction{}).
			Where("checkout_session_id = ?", params.CheckoutSessionID).
			Updates(&paymentTransactionUpdates).Error
	})

	return purchaseOrders, err
}

func (r *InquiryRepo) MultiInquiryMarkAsUnpaid(params MultiInquiryParams) ([]*models.PurchaseOrder, error) {
	var purchaseOrders []*models.PurchaseOrder

	if params.CheckoutSessionID != "" {
		var err = r.db.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs", "CheckoutSessionID").Find(&purchaseOrders, "checkout_session_id = ?", params.CheckoutSessionID).Error
		if err != nil {
			return nil, err
		}
	}

	if len(purchaseOrders) == 0 {
		return nil, errs.ErrInquiryNotFound
	}

	for _, purchaseOrder := range purchaseOrders {
		if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
			return nil, errs.ErrPoIsAlreadyPaid
		}
	}

	var err = r.db.Transaction(func(tx *gorm.DB) (e error) {
		var updates = models.PurchaseOrder{
			MarkAsUnpaidAt: values.Int64(time.Now().Unix()),
			Status:         enums.PurchaseOrderStatusUnpaid,
		}
		e = tx.Model(&updates).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Where("checkout_session_id = ? AND payment_type = ?", params.CheckoutSessionID, enums.PaymentTypeBankTransfer).
			Updates(&updates).Error
		if e != nil {
			return
		}

		var paymentTrxUpdates = models.PaymentTransaction{
			Status:         enums.PaymentStatusUnpaid,
			Note:           params.Note,
			MarkAsUnpaidAt: updates.MarkAsUnpaidAt,
		}
		e = tx.Model(&paymentTrxUpdates).
			Where("checkout_session_id = ?", updates.ID).
			Updates(&paymentTrxUpdates).Error
		return
	})

	return purchaseOrders, err
}

func (r *InquiryRepo) InquiryQuotationHistory(params PaginateInquiryAuditsParams) *query.Pagination {
	var builder = queryfunc.NewInquiryAuditBuilder(queryfunc.InquiryAuditBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("ia.inquiry_id = ?", params.InquiryID)
			builder.Where("ia.action_type IN ?", []string{string(enums.AuditActionTypeInquiryBuyerRejectQuotation)})
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *InquiryRepo) InquiryRemoveItems(form models.InquiryRemoveItemsForm) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var purchaseOrder models.PurchaseOrder
		if err := tx.Select("ID", "CartItemIDs").First(&purchaseOrder, "inquiry_id = ? AND user_id = ?", form.InquiryID, form.GetUserID()).Error; err != nil {
			return err
		}
		if err := r.db.Unscoped().Delete(&models.OrderCartItem{}, "id IN ? AND purchase_order_id = ?", form.ItemIDs, purchaseOrder.ID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *InquiryRepo) PaginateInquiryWithPurchaseOrders(params PaginateInquiryParams) *query.Pagination {
	var builder = queryfunc.NewInquiryPurchaseOrderBuilder(queryfunc.InquiryPurchaseOrderBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludePurchaseOrder: true,
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("iq.user_id = ?", params.JwtClaimsInfo.GetUserID())
			builder.Where("po.id IS NOT NULL")
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *InquiryRepo) InquiryAssignPIC(params models.InquiryAssignPICParam) (updates models.Inquiry, err error) {
	var user models.User
	err = r.db.Model(&models.User{}).Where("id IN ? AND role IN ?", params.AssigneeIDs, []enums.Role{
		enums.RoleLeader,
		enums.RoleStaff,
	}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = eris.Wrap(errs.ErrUserNotFound, "cannot get assignee")
		return
	}
	if err != nil {
		return
	}
	var inquiry models.Inquiry
	err = r.db.Select("ID", "AssigneeIDs", "UserID").First(&inquiry, "id = ?", params.InquiryID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.ErrInquiryNotFound
		}
		return
	}
	updates.AssigneeIDs = params.AssigneeIDs

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&updates).Clauses(clause.Returning{}).
			Where("id = ?", params.InquiryID).Updates(&updates).Error

		var chatRoom models.ChatRoom
		err = tx.Select("ID").Where(map[string]interface{}{
			"inquiry_id":             inquiry.ID,
			"purchase_order_id":      "",
			"bulk_purchase_order_id": "",
			"buyer_id":               inquiry.UserID,
		}).First(&chatRoom).Error
		if err != nil && !r.db.IsRecordNotFoundError(err) {
			return err
		}

		if r.db.IsRecordNotFoundError(err) {
			chatRoom.InquiryID = params.InquiryID
			chatRoom.HostID = params.GetUserID()
			if err := tx.Create(&chatRoom).Error; err != nil {
				return err
			}
		}
		var chatRoomUsers = []*models.ChatRoomUser{{RoomID: chatRoom.ID, UserID: inquiry.UserID}}
		for _, userId := range params.AssigneeIDs {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{RoomID: chatRoom.ID, UserID: userId})
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&chatRoomUsers).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

type AdminUnarchiveInquiryParams struct {
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`
}

func (r *InquiryRepo) AdminUnarchiveInquiry(params AdminUnarchiveInquiryParams) error {
	return r.db.Unscoped().Model(&models.Inquiry{}).Where("id = ?", params.InquiryID).UpdateColumn("DeletedAt", nil).Error
}

type AdminArchiveInquiryParams struct {
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`
}

func (r *InquiryRepo) AdminArchiveInquiry(params AdminUnarchiveInquiryParams) (err error) {
	var iq models.Inquiry
	r.db.Select("id").Where("id = ? AND status = ?", params.InquiryID, enums.InquiryStatusNew).First(&iq)
	if iq.ID == "" {
		err = errs.ErrCanArchiveOnNewInquiryOnly
		return
	}
	return r.db.Unscoped().
		Model(&models.Inquiry{}).
		Where("id = ? AND status = ?", params.InquiryID, enums.InquiryStatusNew).
		UpdateColumn("DeletedAt", time.Now().Unix()).Error
}

type AdminDeleteInquiryParams struct {
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`
}

func (r *InquiryRepo) AdminDeleteInquiry(params AdminDeleteInquiryParams) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Unscoped().
			Delete(&models.Inquiry{}, "id = ? AND deleted_at is not null", params.InquiryID)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrInquiryNotFound
		}

		var err = tx.Unscoped().Delete(&models.InquiryAudit{}, "inquiry_id = ?", params.InquiryID).Error
		if err != nil {
			return err
		}

		err = tx.Unscoped().Delete(&models.InquiryCartItem{}, "inquiry_id = ?", params.InquiryID).Error
		if err != nil {
			return err
		}

		err = tx.Unscoped().Delete(&models.InquirySeller{}, "inquiry_id = ?", params.InquiryID).Error
		if err != nil {
			return err
		}

		var purchaseOrder models.PurchaseOrder
		err = tx.Unscoped().
			Model(&purchaseOrder).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Delete(&models.PurchaseOrder{}, "inquiry_id = ?", params.InquiryID).Error
		if err != nil {
			return err
		}
		if purchaseOrder.ID != "" {
			tx.Delete(&models.PaymentTransaction{}, "purchase_order_id = ?", purchaseOrder.ID)
		}

		var bulkPurchaseOrder models.BulkPurchaseOrder
		err = tx.Unscoped().
			Model(&bulkPurchaseOrder).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Delete(&models.BulkPurchaseOrder{}, "inquiry_id = ?", params.InquiryID).Error
		if err != nil {
			return err
		}

		if bulkPurchaseOrder.ID != "" {
			tx.Delete(&models.PaymentTransaction{}, "bulk_purchase_order_id = ?", bulkPurchaseOrder.ID)
		}

		return nil
	})
}

type CreatePaymentLinkParams struct {
	models.JwtClaimsInfo

	InquiryID     string                              `param:"inquiry_id" validate:"required"`
	Items         []*models.InquiryCartItemCreateForm `json:"items"`
	Quotations    models.InquiryQuotationItems        `json:"quotations" param:"quotations" query:"quotations"`
	ProductWeight *float64                            `json:"product_weight" param:"product_weight" query:"product_weight"`
	ShippingFee   price.Price                         `json:"shipping_fee"`
	TaxPercentage *float64                            `json:"tax_percentage" validate:"min=0,max=100"`
}

func (r *InquiryRepo) InquiryCreatePaymentLink(params CreatePaymentLinkParams) (*models.Inquiry, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("inquiry_payment_link_%s", params.InquiryID), time.Second*30)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var inquiry models.Inquiry
	err = r.db.First(&inquiry, "id = ?", params.InquiryID).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	inquiry.AdminQuotations = params.Quotations

	var buyer models.User
	err = r.db.Select("Name", "Email").First(&buyer, "id = ?", inquiry.UserID).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var items []*models.InquiryCartItem
	if len(params.Items) == 0 {
		var err = r.db.Unscoped().Delete(&models.InquiryCartItem{}, "inquiry_id = ?", inquiry.ID).Error
		if err != nil {
			return nil, err
		}
	} else {
		err = copier.Copy(&items, &params.Items)
		if err != nil {
			return nil, err
		}

		items = lo.Map(items, func(item *models.InquiryCartItem, index int) *models.InquiryCartItem {
			item.InquiryID = inquiry.ID
			if item.UnitPrice.ToFloat64() > 0 {
				// DO nothing
			} else {
				item.UnitPrice = inquiry.GetSampleUnitPrice()
			}

			item.TotalPrice = item.UnitPrice.MultipleInt(int64(item.Qty))
			return item
		})

		err = r.db.Transaction(func(tx *gorm.DB) error {
			var err = tx.Unscoped().Delete(&models.InquiryCartItem{}, "inquiry_id = ?", inquiry.ID).Error
			if err != nil {
				return err
			}

			err = tx.Create(&items).Error

			if err != nil {
				return err
			}
			return err
		})
		if err != nil {
			return nil, err
		}

	}

	var inquiryUpdates = &models.Inquiry{
		QuotationAt:     values.Int64(time.Now().Unix()),
		AdminQuotations: params.Quotations,
		Status:          enums.InquiryStatusQuoteInProcess,
		ShippingFee:     params.ShippingFee.ToPtr(),
		TaxPercentage:   params.TaxPercentage,
		ProductWeight:   params.ProductWeight,
		AssigneeIDs:     lo.Union(append(inquiry.AssigneeIDs, params.GetUserID())),
		CartItems:       items,
	}

	err = r.db.Omit(clause.Associations).Model(&models.Inquiry{}).Where("id = ?", inquiry.ID).Updates(inquiryUpdates).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	inquiry.ShippingFee = inquiryUpdates.ShippingFee
	inquiry.TaxPercentage = inquiryUpdates.TaxPercentage
	inquiry.ProductWeight = inquiryUpdates.ProductWeight

	purchaseOrder, err := r.InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
		Inquiry:       &inquiry,
		JwtClaimsInfo: params.JwtClaimsInfo,
		InquiryID:     params.InquiryID,
		CartItems:     items,
		UserID:        inquiry.UserID,
		PaymentType:   enums.PaymentTypeCard,
		UpdatePricing: true,
	})
	if err != nil {
		return nil, err
	}

	stripeConfig, err := stripehelper.GetCurrencyConfig(purchaseOrder.Inquiry.Currency)
	if err != nil {
		return nil, err
	}

	var paymentLineItems []*stripe.PaymentLinkLineItemParams
	for _, cartItem := range items {
		priceItem, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(inquiry.Currency)),
			UnitAmount: stripe.Int64(cartItem.UnitPrice.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String(fmt.Sprintf("%s/%s/%s", inquiry.Title, cartItem.Size, cartItem.ColorName)),
			},
		})
		if err != nil {
			return nil, err
		}

		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(int64(cartItem.Qty)),
			Price:    &priceItem.ID,
		})
	}

	// Append for tax and shipping
	if purchaseOrder.Tax.GreaterThan(0) {
		stripePriceParams := &stripe.PriceParams{
			Currency:   stripe.String(string(inquiry.Currency)),
			UnitAmount: stripe.Int64(purchaseOrder.Tax.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Tax"),
			},
		}
		priceItemTax, err := stripePrice.New(stripePriceParams)
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTax.ID,
		})
	}

	if purchaseOrder.ShippingFee.GreaterThan(0) {
		priceItemShippingFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(inquiry.Currency)),
			UnitAmount: stripe.Int64(purchaseOrder.ShippingFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Shipping Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemShippingFee.ID,
		})
	}

	if purchaseOrder.TransactionFee.GreaterThan(0) {
		priceItemTransactionFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(inquiry.Currency)),
			UnitAmount: stripe.Int64(purchaseOrder.TransactionFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Transaction Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTransactionFee.ID,
		})
	}

	var stripeParams = stripehelper.CreatePaymentLinkParams{
		Currency: purchaseOrder.Inquiry.Currency,
		Metadata: map[string]string{
			"inquiry_id":                  purchaseOrder.InquiryID,
			"inquiry_reference_id":        purchaseOrder.Inquiry.ReferenceID,
			"purchase_order_id":           purchaseOrder.ID,
			"purchase_order_reference_id": purchaseOrder.ReferenceID,
			"cart_item_ids":               strings.Join(purchaseOrder.CartItemIDs, ","),
			"action_source":               string(stripehelper.ActionSourceInquiryPayment),
		},
		RedirectURL: fmt.Sprintf("%s/inquiries/%s", r.db.Configuration.BrandPortalBaseURL, inquiry.ID),
		LineItems:   paymentLineItems,
	}

	pl, err := stripehelper.GetInstance().CreatePaymentLink(stripeParams)
	if err != nil {
		return nil, err
	}
	inquiry.SamplePaymentLink = helper.AddURLQuery(pl.URL,
		map[string]string{
			"client_reference_id": inquiry.ReferenceID,
			"prefilled_email":     buyer.Email,
		},
	)

	var updates = models.Inquiry{
		SamplePaymentLink:   inquiry.SamplePaymentLink,
		SamplePaymentLinkID: pl.ID,
	}
	err = r.db.Model(&models.Inquiry{}).Where("id = ?", inquiry.ID).Updates(&updates).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	purchaseOrder.Inquiry = nil
	inquiry.PurchaseOrder = purchaseOrder

	return &inquiry, nil

}

type CreateBuyerPaymentLinkRequest struct {
	models.JwtClaimsInfo
	BuyerID string                       `json:"buyer_id" validate:"required"`
	Data    []CreateBuyerPaymentLinkItem `param:"data" validate:"required"`
}

type CreateBuyerPaymentLinkItem struct {
	InquiryID     string                       `json:"inquiry_id" param:"inquiry_id" validate:"required"`
	Items         []*models.OrderCartItem      `json:"items" param:"items" validate:"required"`
	Quotations    models.InquiryQuotationItems `json:"quotations" param:"quotations" query:"quotations"`
	ProductWeight *float64                     `json:"product_weight" param:"product_weight" query:"product_weight"`
	ShippingFee   price.Price                  `json:"shipping_fee"`
	TaxPercentage *float64                     `json:"tax_percentage" validate:"min=0,max=100"`
}

type CreateBuyerPaymentLinkResponse struct {
	PaymentLink    string                  `json:"payment_link"`
	PurchaseOrders []*models.PurchaseOrder `json:"purchase_orders"`
}

func (r *InquiryRepo) CreateBuyerPaymentLink(params CreateBuyerPaymentLinkRequest) (*CreateBuyerPaymentLinkResponse, error) {
	var user models.User
	if err := r.db.First(&user, "id = ? ", params.BuyerID).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	var inquiryIDs []string
	for _, data := range params.Data {
		inquiryIDs = append(inquiryIDs, data.InquiryID)
	}
	var inquiries = make(models.Inquiries, 0, len(inquiryIDs))
	if err := r.db.Find(&inquiries, "id IN ? AND user_id = ?", inquiryIDs, params.BuyerID).Error; err != nil {
		return nil, err
	}
	var dbInquiryIDs = inquiries.IDs()
	for _, id := range inquiryIDs {
		if !helper.StringContains(dbInquiryIDs, id) {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "inquiry_id:%s", id)
		}
	}
	var mapInquiryIDToInquiry = make(map[string]*models.Inquiry, len(inquiries))
	validStatus := []enums.InquiryStatus{enums.InquiryStatusNew, enums.InquiryStatusQuoteInProcess}
	var currency = inquiries[0].Currency
	for _, iq := range inquiries {
		if ok := slices.Contains(validStatus, iq.Status); !ok {
			return nil, eris.Wrapf(errs.ErrInquiryInvalidToCreatePaymentLink, "inquiry_id:%s", iq.ID)
		}
		if iq.Currency != currency {
			return nil, errs.ErrOrderCurrencyMismatch
		}
		mapInquiryIDToInquiry[iq.ID] = iq
	}

	var purchaseOrders = make(models.PurchaseOrders, 0, len(inquiryIDs))
	if err := r.db.Find(&purchaseOrders, "inquiry_id IN ? AND user_id = ?", inquiryIDs, params.BuyerID).Error; err != nil {
		return nil, err
	}
	var dbPoInquiryIDs = purchaseOrders.InquiryIDs()
	for _, id := range inquiryIDs {
		if !helper.StringContains(dbPoInquiryIDs, id) {
			return nil, eris.Wrapf(errs.ErrPONotFound, "inquiry_id:%s", id)
		}
	}
	var mapInquiryIDToPurchaseOrder = make(map[string]*models.PurchaseOrder, len(purchaseOrders))
	for _, po := range purchaseOrders {
		if po.Status != enums.PurchaseOrderStatusPending {
			return nil, eris.Wrapf(errs.ErrInquiryInvalidToCreatePaymentLink, "inquiry_id:%s", po.InquiryID)
		}
		if po.Currency != currency {
			return nil, errs.ErrOrderCurrencyMismatch
		}
		mapInquiryIDToPurchaseOrder[po.InquiryID] = po
	}
	stripeConfig, err := stripehelper.GetCurrencyConfig(currency)
	if err != nil {
		return nil, err
	}
	var inquiriesToUpdate = make(models.Inquiries, 0, len(inquiries))
	var purchaseOrdersToUpdate = make(models.PurchaseOrders, 0, len(purchaseOrders))
	var orderCartItemsToCreate models.OrderCartItems
	var lineItems []*stripe.PaymentLinkLineItemParams
	var totalTax = price.NewFromFloat(0)
	var totalTransactionFee = price.NewFromFloat(0)
	var totalShippingFee = price.NewFromFloat(0)
	for _, data := range params.Data {
		inquiry, ok := mapInquiryIDToInquiry[data.InquiryID]
		if !ok {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "inquiry_id:%s", data.InquiryID)
		}
		purchaseOrder, ok := mapInquiryIDToPurchaseOrder[inquiry.ID]
		if !ok {
			return nil, eris.Wrapf(errs.ErrPONotFound, "inquiry_id:%s", inquiry.ID)
		}

		inquiry.ShippingFee = data.ShippingFee.ToPtr()
		inquiry.TaxPercentage = data.TaxPercentage
		inquiry.ProductWeight = data.ProductWeight
		inquiry.AdminQuotations = data.Quotations
		inquiry.Status = enums.InquiryStatusQuoteInProcess
		inquiry.BuyerQuotationStatus = enums.InquirySkuStatusApproved
		inquiriesToUpdate = append(inquiriesToUpdate, inquiry)

		var subTotalPrice = price.NewFromFloat(0)
		for _, item := range data.Items {
			var unitPrice = inquiry.GetSampleUnitPrice()
			var totalPrice = unitPrice.MultipleInt(int64(item.Qty))
			subTotalPrice = subTotalPrice.Add(totalPrice)

			orderCartItemsToCreate = append(orderCartItemsToCreate, &models.OrderCartItem{
				PurchaseOrderID: purchaseOrder.ID,
				UnitPrice:       unitPrice,
				TotalPrice:      totalPrice,
				Size:            item.Size,
				Qty:             item.Qty,
				ColorName:       item.ColorName,
				NoteToSupplier:  item.NoteToSupplier,
				Style:           item.Style,
			})

			priceItem, err := stripePrice.New(&stripe.PriceParams{
				Currency:   stripe.String(string(currency)),
				UnitAmount: stripe.Int64(item.UnitPrice.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
				ProductData: &stripe.PriceProductDataParams{
					Name: stripe.String(fmt.Sprintf("%s - %s - %s (Sample)", purchaseOrder.ProductName, item.Size, item.ColorName)),
				},
			})
			if err != nil {
				return nil, err
			}
			lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
				Price:    &priceItem.ID,
				Quantity: stripe.Int64(int64(item.Qty)),
			})
		}

		purchaseOrder.SubTotal = &subTotalPrice
		purchaseOrder.TaxPercentage = data.TaxPercentage
		purchaseOrder.ShippingFee = data.ShippingFee.ToPtr()
		purchaseOrder.PaymentType = enums.PaymentTypeCard
		purchaseOrder.Quotations = data.Quotations
		if err := purchaseOrder.UpdatePrices(); err != nil {
			return nil, err
		}

		purchaseOrdersToUpdate = append(purchaseOrdersToUpdate, purchaseOrder)

		if purchaseOrder.Tax.GreaterThan(0) {
			totalTax = totalTax.AddPtr(purchaseOrder.Tax)
		}
		if purchaseOrder.ShippingFee.GreaterThan(0) {
			totalShippingFee = totalShippingFee.AddPtr(purchaseOrder.ShippingFee)
		}
		if purchaseOrder.TransactionFee.GreaterThan(0) {
			totalTransactionFee = totalTransactionFee.AddPtr(purchaseOrder.TransactionFee)
		}
	}
	if totalShippingFee.GreaterThan(0) {
		priceItemShippingFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalShippingFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Shipping Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemShippingFee.ID,
		})
	}
	if totalTransactionFee.GreaterThan(0) {
		priceItemTxFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalTransactionFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Transaction Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTxFee.ID,
		})
	}
	if totalTax.GreaterThan(0) {
		priceItemTax, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalTax.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Tax"),
			},
		})
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTax.ID,
		})
	}
	var paymentLink string
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&inquiriesToUpdate).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&purchaseOrdersToUpdate).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&models.OrderCartItem{}, "purchase_order_id IN ?", purchaseOrders.IDs()).Error; err != nil {
			return err
		}
		if err := tx.Create(orderCartItemsToCreate).Error; err != nil {
			return err
		}

		var checkoutSessionID = helper.GenerateCheckoutSessionID()
		var poItemsIDs []string
		for _, item := range orderCartItemsToCreate {
			poItemsIDs = append(poItemsIDs, item.ID)
		}
		var stripeParams = stripehelper.CreatePaymentLinkParams{
			Currency: currency,
			Metadata: map[string]string{
				"purchase_order_cart_item_ids": strings.Join(poItemsIDs, ","),
				"user_id":                      params.BuyerID,
				"checkout_session_id":          checkoutSessionID,
				"action_source":                string(stripehelper.ActionSourceOrderCartPayment),
			},
			RedirectURL: fmt.Sprintf("%s/order-checkout?checkout_session_id=%s", r.db.Configuration.BrandPortalBaseURL, checkoutSessionID),
			LineItems:   lineItems,
		}

		pl, err := stripehelper.GetInstance().CreatePaymentLink(stripeParams)
		if err != nil {
			return err
		}
		paymentLink = helper.AddURLQuery(pl.URL,
			map[string]string{
				"prefilled_email": user.Email,
			},
		)
		return nil
	}); err != nil {
		return nil, err
	}
	return &CreateBuyerPaymentLinkResponse{
		PaymentLink:    paymentLink,
		PurchaseOrders: purchaseOrdersToUpdate,
	}, nil

}

type InquirySyncSampleParams struct {
	models.JwtClaimsInfo

	// Payment
	InquiryID             string             `param:"inquiry_id" validate:"required"`
	TransactionRefID      string             `json:"transaction_ref_id" validate:"required_if=PaymentType bank_transfer"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required_if=PaymentType bank_transfer"`

	// cart items
	Items         []*models.InquiryCartItemCreateForm `json:"items" validate:"required"`
	Quotations    models.InquiryQuotationItems        `json:"quotations" param:"quotations" query:"quotations" validate:"required"`
	ProductWeight *float64                            `json:"product_weight" param:"product_weight" query:"product_weight"`
	ShippingFee   price.Price                         `json:"shipping_fee"`
	TaxPercentage *float64                            `json:"tax_percentage" validate:"min=0,max=100"`
}

func (r *InquiryRepo) InquirySyncSample(form InquirySyncSampleParams) (*models.PurchaseOrder, error) {
	var items []*models.InquiryCartItem

	var inquiry models.Inquiry
	var err error

	err = r.db.First(&inquiry, "id = ?", form.InquiryID).Error
	if err != nil {
		return nil, err
	}

	var admin models.User
	err = r.db.Select("ID", "Name").First(&admin, "id = ?", form.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.Select("ID", "Name").First(&user, "id = ?", inquiry.UserID).Error
	if err != nil {
		return nil, err
	}

	inquiry.AdminQuotations = form.Quotations

	if len(form.Items) == 0 {
		var err = r.db.Unscoped().Delete(&models.InquiryCartItem{}, "inquiry_id = ?", form.InquiryID).Error
		if err != nil {
			return nil, err
		}

	} else {
		err = copier.Copy(&items, &form.Items)
		if err != nil {
			return nil, err
		}

		items = lo.Map(items, func(item *models.InquiryCartItem, index int) *models.InquiryCartItem {
			item.InquiryID = form.InquiryID

			if item.UnitPrice.Equal(0) {
				item.UnitPrice = inquiry.GetSampleUnitPrice()
			}
			item.TotalPrice = item.UnitPrice.MultipleInt(int64(item.Qty))
			return item
		})

		err = r.db.Transaction(func(tx *gorm.DB) error {
			var err = tx.Unscoped().Delete(&models.InquiryCartItem{}, "inquiry_id = ?", form.InquiryID).Error
			if err != nil {
				return err
			}

			err = tx.Create(&items).Error

			if err != nil {
				return err
			}

			return err
		})
		if err != nil {
			return nil, err
		}
	}

	var inquiryUpdates = &models.Inquiry{
		QuotationAt:          values.Int64(time.Now().Unix()),
		AdminQuotations:      form.Quotations,
		Status:               enums.InquiryStatusFinished,
		BuyerQuotationStatus: enums.InquirySkuStatusApproved,
		ShippingFee:          form.ShippingFee.ToPtr(),
		TaxPercentage:        form.TaxPercentage,
		ProductWeight:        form.ProductWeight,
		AssigneeIDs:          lo.Union(append(inquiry.AssigneeIDs, form.GetUserID())),
	}

	err = r.db.Model(&models.Inquiry{}).Where("id = ?", inquiry.ID).Updates(inquiryUpdates).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	purchaseOrder, err := r.InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
		JwtClaimsInfo: form.JwtClaimsInfo,
		InquiryID:     form.InquiryID,
		PaymentType:   enums.PaymentTypeBankTransfer,
		UserID:        inquiry.UserID,
		UpdatePricing: true,
		CartItems:     items,
		Inquiry:       &inquiry,
	})
	if err != nil {
		return nil, err
	}

	r.db.Transaction(func(tx *gorm.DB) error {
		var updates = models.PurchaseOrder{
			Status:                enums.PurchaseOrderStatusPaid,
			PaymentType:           enums.PaymentTypeBankTransfer,
			TransactionRefID:      form.TransactionRefID,
			TransactionAttachment: form.TransactionAttachment,
			TransferedAt:          values.Int64(time.Now().Unix()),
			MarkAsPaidAt:          values.Int64(time.Now().Unix()),
		}
		updates.TaxPercentage = purchaseOrder.Inquiry.TaxPercentage

		if len(purchaseOrder.Quotations) > 0 {
			sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
				return item.Type == enums.InquiryTypeSample
			})
			if sampleQuotation != nil {
				updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
				updates.StartDate = updates.MarkAsPaidAt
				updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
			}
		}

		var transaction = models.PaymentTransaction{
			Currency:          inquiry.Currency,
			PurchaseOrderID:   purchaseOrder.ID,
			PaidAmount:        purchaseOrder.TotalPrice,
			PaymentType:       enums.PaymentTypeBankTransfer,
			UserID:            purchaseOrder.UserID,
			Attachments:       &models.Attachments{form.TransactionAttachment},
			TransactionRefID:  form.TransactionRefID,
			Status:            enums.PaymentStatusPaid,
			PaymentPercentage: values.Float64(100),
			TotalAmount:       purchaseOrder.TotalPrice,
			Milestone:         enums.PaymentMilestoneFinalPayment,
			MarkAsPaidAt:      updates.MarkAsPaidAt,
			Metadata: &models.PaymentTransactionMetadata{
				InquiryID:                purchaseOrder.Inquiry.ID,
				InquiryReferenceID:       purchaseOrder.Inquiry.ReferenceID,
				PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
				PurchaseOrderID:          purchaseOrder.ID,
			},
		}

		err = r.db.Create(&transaction).Error
		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		updates.PaymentTransactionReferenceID = transaction.ReferenceID
		err = r.db.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates).Error
		if err != nil {
			return err
		}

		var audits = []models.InquiryAudit{
			{
				InquiryID:   inquiry.ID,
				ActionType:  enums.AuditActionTypeInquiryAdminSendBuyerQuotation,
				UserID:      user.ID,
				Description: fmt.Sprintf("%s has sent quotation to buyer %s", admin.Name, user.Name),
				Metadata: &models.InquiryAuditMetadata{
					After: map[string]interface{}{
						"quotations": form.Quotations,
					},
				},
			},
			{
				InquiryID:   inquiry.ID,
				ActionType:  enums.AuditActionTypeInquiryBuyerApproveQuotation,
				UserID:      user.ID,
				Description: fmt.Sprintf("%s has approved quotation", user.Name),
				Metadata: &models.InquiryAuditMetadata{
					After: map[string]interface{}{
						"quotations":          inquiry.AdminQuotations,
						"approve_reject_meta": inquiry.ApproveRejectMeta,
						"quotation_at":        inquiry.QuotationAt,
					},
				},
			},
			{
				InquiryID:       inquiry.ID,
				ActionType:      enums.AuditActionTypeInquirySamplePoCreated,
				UserID:          inquiry.UserID,
				Description:     fmt.Sprintf("New sample PO %s has been created for inquiry", purchaseOrder.ReferenceID),
				PurchaseOrderID: purchaseOrder.ID,
			},
			{
				InquiryID:   purchaseOrder.InquiryID,
				ActionType:  enums.AuditActionTypeInquiryAdminMarkAsPaid,
				UserID:      admin.ID,
				Description: fmt.Sprintf("Admin %s has confirmed the payment", admin.Name),
			},
		}

		err = tx.Create(&audits).Error
		return err
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return purchaseOrder, err
}

func (r *InquiryRepo) GetInquiryRemindAdmin() models.Inquiries {
	var inquiries []*models.Inquiry
	query.New(r.db, queryfunc.NewInquiryRemindBuilder(queryfunc.InquiryRemindBuilderOptions{
		IncludeAssignee: true,
		IncludeUser:     true,
	})).
		WhereFunc(func(builder *query.Builder) {
			var buyerStatues = []enums.InquirySkuStatus{
				enums.InquirySkuStatusRejected,
				enums.InquirySkuStatusNew,
			}
			builder.Where("( iq.status = ? OR (iq.status = ? AND iq.buyer_quotation_status IN ?) ) AND iq.deleted_at IS NULL", enums.InquiryStatusNew, enums.InquiryStatusQuoteInProcess, buyerStatues)

			builder.Where("EXTRACT(EPOCH FROM now()) - iq.created_at > ?", 12*int(time.Hour.Seconds()))
		}).
		FindFunc(&inquiries)

	return inquiries
}

func (r *InquiryRepo) CloseInquiry(form models.InquiryCloseForm) error {
	var statues = []enums.InquiryStatus{
		enums.InquiryStatusNew,
		enums.InquiryStatusQuoteInProcess,
	}

	var updates = models.Inquiry{
		Status:      enums.InquiryStatusClosed,
		CloseReason: form.CloseReason,
	}
	var sqlResult = r.db.Unscoped().Model(&models.Inquiry{}).Where("id = ? AND status IN ?", form.InquiryID, statues).Updates(&updates)
	if sqlResult.Error != nil {
		return sqlResult.Error
	}

	if sqlResult.RowsAffected == 0 {
		return errs.ErrInquiryNotFound
	}
	return nil
}

type InquiryNoteMarkSeenParams struct {
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" param:"inquiry_id"`
}

func (r *InquiryRepo) InquiryNoteMarkSeen(params InquiryNoteMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ?", enums.CommentTargetTypeInquiryInternalNotes, params.InquiryID).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

type InquiryNoteUnreadCountParams struct {
	models.JwtClaimsInfo
	InquiryID string `json:"inquiry_id" param:"inquiry_id"`
}

type InquiryNoteUnreadCountResponse struct {
	TotalCount int64 `json:"total_count"`
}

func (r *InquiryRepo) InquiryNoteUnreadCount(params InquiryNoteUnreadCountParams) InquiryNoteUnreadCountResponse {
	var resp InquiryNoteUnreadCountResponse
	r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeInquiryInternalNotes, params.InquiryID).Count(&resp.TotalCount)

	return resp
}

type ExportInquiriesParams struct {
	models.JwtClaimsInfo

	PaginateInquiryParams
}

func (r *InquiryRepo) ExportExcel(params ExportInquiriesParams) (*models.Attachment, error) {
	params.PaginateInquiryParams.JwtClaimsInfo = params.JwtClaimsInfo
	params.PaginateInquiryParams.IsQueryAll = true
	params.PaginateInquiryParams.IncludeUser = true
	params.PaginateInquiryParams.IncludeAssignee = true
	params.PaginateInquiryParams.WithoutCount = true
	var result = r.PaginateInquiry(params.PaginateInquiryParams)
	if result == nil || result.Records == nil {
		return nil, errors.New("empty response")
	}
	trans, ok := result.Records.([]*models.Inquiry)
	if !ok {
		return nil, eris.Errorf("failed to parse from %T", result.Records)
	}

	fileContent, err := models.Inquiries(trans).ToExcel()
	if err != nil {
		return nil, err
	}

	var contentType = models.ContentTypeXLSX
	url := fmt.Sprintf("uploads/inquiries/export/export_inquiry_user_%s%s", params.GetUserID(), contentType.GetExtension())
	_, err = s3.New(r.db.Configuration).UploadFile(s3.UploadFileParams{
		Data:        bytes.NewReader(fileContent),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(contentType),
		ACL:         "private",
		Key:         url,
	})
	if err != nil {
		return nil, err
	}

	var resp = models.Attachment{
		FileKey:     url,
		ContentType: string(contentType),
	}
	return &resp, err
}

type MultiInquiryPreviewCheckoutParams struct {
	models.JwtClaimsInfo

	PaymentType enums.PaymentType         `json:"payment_type" validate:"oneof=bank_transfer card"`
	CartItems   []*models.InquiryCartItem `json:"-"`
	CartItemIDs []string                  `json:"cart_item_ids" params:"cart_item_ids" validate:"required"`
	UserID      string                    `param:"user_id"`

	UpdatePricing bool `json:"-"`
}

func (r *InquiryRepo) MultiInquiryPreviewCheckout(params MultiInquiryPreviewCheckoutParams) ([]*models.PurchaseOrder, error) {
	var orders []*models.PurchaseOrder

	var cartItems = []*models.InquiryCartItem{}
	var err = r.db.Find(&cartItems, "id IN ? AND (checkout_session_id IS NULL OR checkout_session_id = ?)", params.CartItemIDs, "").Error
	if err != nil {
		return nil, err
	}

	if len(cartItems) <= 0 {
		return nil, errs.ErrInquiryCartInvalidToCheckout
	}

	var inquiryIDs []string
	for _, cartItem := range cartItems {
		if !helper.StringContains(inquiryIDs, cartItem.InquiryID) {
			inquiryIDs = append(inquiryIDs, cartItem.InquiryID)
		}
	}

	for _, inquiryID := range inquiryIDs {
		var items = []*models.InquiryCartItem{}
		for _, cartItem := range cartItems {
			if cartItem.InquiryID == inquiryID {
				items = append(items, cartItem)
			}
		}

		if len(items) > 0 {
			order, err := r.InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
				InquiryID:     inquiryID,
				CartItems:     items,
				JwtClaimsInfo: params.JwtClaimsInfo,
				UserID:        params.UserID,
				PaymentType:   params.PaymentType,
				UpdatePricing: params.UpdatePricing,
			})

			if err != nil {
				return nil, err
			}

			if order != nil {
				orders = append(orders, order)
			}
		}
	}

	err = r.db.Model(&models.InquiryCartItem{}).Where("id IN ? AND inquiry_id IN ?", params.CartItemIDs, inquiryIDs).Updates(&models.InquiryCartItem{WaitingForCheckout: values.Bool(true)}).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&models.InquiryCartItem{}).Where("id NOT IN ? AND inquiry_id IN ?", params.CartItemIDs, inquiryIDs).Updates(&models.InquiryCartItem{WaitingForCheckout: values.Bool(false)}).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

type MultiInquiryCheckoutParams struct {
	models.JwtClaimsInfo

	CartItemIDs     []string          `json:"cart_item_ids" params:"cart_item_ids" validate:"required"`
	PaymentType     enums.PaymentType `json:"payment_type" validate:"oneof=bank_transfer card"`
	PaymentMethodID string            `json:"payment_method_id" validate:"required_if=PaymentType card"`

	TransactionRefID      string             `json:"transaction_ref_id" validate:"required_if=PaymentType bank_transfer"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required_if=PaymentType bank_transfer"`
}

type MultiInquiryCheckoutResponse struct {
	Orders                    []*models.PurchaseOrder         `json:"orders"`
	CheckoutSessionID         string                          `json:"checkout_session_id"`
	PaymentTransaction        *models.PaymentTransaction      `json:"payment_transaction"`
	PaymentIntentNextAction   *stripe.PaymentIntentNextAction `json:"payment_intent_next_action,omitempty"`
	PaymentIntentClientSecret string                          `json:"payment_intent_client_secret,omitempty"`
}

func (r *InquiryRepo) MultiInquiryCheckout(params MultiInquiryCheckoutParams) (*MultiInquiryCheckoutResponse, error) {
	var checkoutSessionID = helper.GenerateCheckoutSessionID()

	var resp = MultiInquiryCheckoutResponse{
		CheckoutSessionID: checkoutSessionID,
	}
	var orders []*models.PurchaseOrder

	var cartItems []*models.InquiryCartItem
	var err = r.db.Find(&cartItems, "COALESCE(checkout_session_id,'') = '' AND id IN ?", params.CartItemIDs).Error
	if err != nil {
		return nil, err
	}

	if len(cartItems) <= 0 {
		return nil, errs.ErrInquiryCartInvalidToCheckout
	}

	var inquiryIDs []string
	var orderIDs []string
	for _, cartItem := range cartItems {
		if !helper.StringContains(inquiryIDs, cartItem.InquiryID) {
			inquiryIDs = append(inquiryIDs, cartItem.InquiryID)
		}
	}

	for _, inquiryID := range inquiryIDs {
		order, err := r.MultiInquiryCheckoutOneOrder(InquiryCheckoutParams{
			JwtClaimsInfo:         params.JwtClaimsInfo,
			InquiryID:             inquiryID,
			PaymentType:           params.PaymentType,
			PaymentMethodID:       params.PaymentMethodID,
			TransactionRefID:      params.TransactionRefID,
			TransactionAttachment: params.TransactionAttachment,
			CheckoutSessionID:     checkoutSessionID,
		})
		if err != nil {
			return nil, err
		}

		if order != nil {
			orders = append(orders, order)
			orderIDs = append(orderIDs, order.ID)
		}
	}

	if params.PaymentType == enums.PaymentTypeBankTransfer {
		r.db.Transaction(func(tx *gorm.DB) error {
			var totalAmount price.Price
			var currency enums.Currency
			var purchaseOrderIDs []string
			var purchaseOrderReferenceIDs []string
			var inquiryIDs []string
			var inquiryReferenceIDs []string
			var transactionRefID = helper.GeneratePaymentTransactionReferenceID()
			for _, purchaseOrder := range orders {
				currency = purchaseOrder.Inquiry.Currency
				totalAmount = totalAmount.AddPtr(purchaseOrder.TotalPrice)

				var updates = models.PurchaseOrder{
					Status:                        enums.PurchaseOrderStatusWaitingConfirm,
					PaymentType:                   params.PaymentType,
					TransactionRefID:              params.TransactionRefID,
					TransactionAttachment:         params.TransactionAttachment,
					TransferedAt:                  values.Int64(time.Now().Unix()),
					Currency:                      purchaseOrder.Inquiry.Currency,
					CheckoutSessionID:             checkoutSessionID,
					PaymentTransactionReferenceID: transactionRefID,
				}
				updates.TaxPercentage = purchaseOrder.Inquiry.TaxPercentage
				var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
				if sqlResult.Error != nil {
					return eris.Wrap(sqlResult.Error, sqlResult.Error.Error())
				}

				if sqlResult.RowsAffected == 0 {
					return eris.New("Purchase order not found")
				}

				purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
				purchaseOrderReferenceIDs = append(purchaseOrderReferenceIDs, purchaseOrder.ReferenceID)

				inquiryIDs = append(inquiryIDs, purchaseOrder.Inquiry.ID)
				inquiryReferenceIDs = append(inquiryReferenceIDs, purchaseOrder.Inquiry.ReferenceID)
			}

			// create transaction
			var transaction = models.PaymentTransaction{
				ReferenceID:       transactionRefID,
				PaidAmount:        totalAmount.ToPtr(),
				PaymentType:       params.PaymentType,
				Milestone:         enums.PaymentMilestoneFinalPayment,
				UserID:            params.GetUserID(),
				TransactionRefID:  params.TransactionRefID,
				Status:            enums.PaymentStatusWaitingConfirm,
				PaymentPercentage: values.Float64(100),
				TotalAmount:       totalAmount.ToPtr(),
				Currency:          currency,
				CheckoutSessionID: checkoutSessionID,
				PurchaseOrderIDs:  purchaseOrderIDs,
				Metadata: &models.PaymentTransactionMetadata{
					InquiryIDs:                inquiryIDs,
					InquiryReferenceIDs:       inquiryReferenceIDs,
					PurchaseOrderReferenceIDs: purchaseOrderReferenceIDs,
					PurchaseOrderIDs:          purchaseOrderIDs,
				},
			}
			if params.TransactionAttachment != nil {
				transaction.Attachments = &models.Attachments{params.TransactionAttachment}
			}
			err = tx.Create(&transaction).Error
			if err != nil {
				return eris.Wrap(err, "Create transaction error")
			}
			resp.PaymentTransaction = &transaction
			return err

		})

		if err != nil {
			return nil, err
		}
	}

	if params.PaymentType == enums.PaymentTypeCard {
		var user models.User
		err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", params.GetUserID()).Error
		if err != nil {
			return nil, err
		}
		var totalAmount price.Price
		var currency enums.Currency
		var finalCardItemIDs []string

		for _, purchaseOrder := range orders {
			currency = purchaseOrder.Inquiry.Currency
			totalAmount = totalAmount.AddPtr(purchaseOrder.TotalPrice)
			finalCardItemIDs = append(finalCardItemIDs, purchaseOrder.CartItemIDs...)
		}

		stripeConfig, err := stripehelper.GetCurrencyConfig(currency)
		if err != nil {
			return nil, err
		}
		var stripeParams = stripehelper.CreatePaymentIntentParams{
			Amount:                  totalAmount.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64(),
			Currency:                currency,
			PaymentMethodID:         params.PaymentMethodID,
			CustomerID:              user.StripeCustomerID,
			IsCaptureMethodManually: false,
			Description:             fmt.Sprintf("Charges for checkout session %s", checkoutSessionID),
			PaymentMethodTypes:      []string{"card"},
			Metadata: map[string]string{
				"cart_item_ids":       strings.Join(finalCardItemIDs, ","),
				"checkout_session_id": checkoutSessionID,
				"action_source":       string(stripehelper.ActionSourceMultiInquiryPayment),
			},
		}

		pi, err := stripehelper.GetInstance().CreatePaymentIntent(stripeParams)
		if err != nil {
			return nil, err
		}
		if pi.Status != stripe.PaymentIntentStatusSucceeded {
			if pi.NextAction != nil {
				intent, err := stripehelper.GetInstance().ConfirmPaymentIntent(stripehelper.ConfirmPaymentIntentParams{
					PaymentIntentID: pi.ID,
					ReturnURL:       fmt.Sprintf("%s/api/v1/callback/stripe/payment_intents/inquiry_carts/%s/confirm?cart_items=%s", r.db.Configuration.ServerBaseURL, checkoutSessionID, strings.Join(params.CartItemIDs, ",")),
				})
				if err != nil {
					return nil, err
				}

				if intent.Status == stripe.PaymentIntentStatusSucceeded {
					goto PaymentSuccess
				}

				resp.PaymentIntentNextAction = intent.NextAction
				resp.PaymentIntentClientSecret = intent.ClientSecret
				return &resp, nil

			} else {
				return nil, eris.Errorf("Payment error with status %s", pi.Status)
			}
		}

	PaymentSuccess:
		err = r.db.Transaction(func(tx *gorm.DB) error {
			var purchaseOrderIDs []string
			var purchaseOrderReferenceIDs []string
			var inquiryIDs []string
			var inquiryReferenceIDs []string
			var transactionRefID = helper.GeneratePaymentTransactionReferenceID()

			for _, purchaseOrder := range orders {
				var updates = models.PurchaseOrder{
					PaymentIntentID:               pi.ID,
					Status:                        enums.PurchaseOrderStatusPaid,
					PaymentType:                   params.PaymentType,
					MarkAsPaidAt:                  values.Int64(time.Now().Unix()),
					TransactionRefID:              params.TransactionRefID,
					TransactionAttachment:         params.TransactionAttachment,
					TransferedAt:                  values.Int64(time.Now().Unix()),
					Currency:                      purchaseOrder.Inquiry.Currency,
					CheckoutSessionID:             checkoutSessionID,
					PaymentTransactionReferenceID: transactionRefID,
				}

				if len(purchaseOrder.Quotations) > 0 {
					sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
						return item.Type == enums.InquiryTypeSample
					})
					if sampleQuotation != nil {
						updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
						updates.StartDate = updates.MarkAsPaidAt
						updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
					}
				}

				updates.TaxPercentage = purchaseOrder.Inquiry.TaxPercentage
				var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
				if sqlResult.Error != nil {
					return sqlResult.Error
				}

				if sqlResult.RowsAffected == 0 {
					return eris.New("Purchase order is not found")
				}

				var sqlResult2 = tx.Model(&models.Inquiry{}).Where("id = ?", purchaseOrder.InquiryID).UpdateColumn("Status", enums.InquiryStatusFinished)
				if sqlResult2.Error != nil {
					return sqlResult2.Error
				}

				purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
				purchaseOrderReferenceIDs = append(purchaseOrderReferenceIDs, purchaseOrder.ReferenceID)

				inquiryIDs = append(inquiryIDs, purchaseOrder.Inquiry.ID)
				inquiryReferenceIDs = append(inquiryReferenceIDs, purchaseOrder.Inquiry.ReferenceID)
			}

			// create transaction
			var transaction = models.PaymentTransaction{
				ReferenceID:       transactionRefID,
				PaidAmount:        totalAmount.ToPtr(),
				PaymentType:       params.PaymentType,
				UserID:            params.GetUserID(),
				TransactionRefID:  params.TransactionRefID,
				PaymentIntentID:   pi.ID,
				Status:            enums.PaymentStatusPaid,
				TotalAmount:       totalAmount.ToPtr(),
				Milestone:         enums.PaymentMilestoneFinalPayment,
				PaymentPercentage: values.Float64(100),
				MarkAsPaidAt:      values.Int64(time.Now().Unix()),
				Currency:          currency,
				CheckoutSessionID: checkoutSessionID,
				PurchaseOrderIDs:  purchaseOrderIDs,
				Metadata: &models.PaymentTransactionMetadata{
					InquiryIDs:                inquiryIDs,
					InquiryReferenceIDs:       inquiryReferenceIDs,
					PurchaseOrderReferenceIDs: purchaseOrderReferenceIDs,
					PurchaseOrderIDs:          purchaseOrderIDs,
				},
			}
			if params.TransactionAttachment != nil {
				transaction.Attachments = &models.Attachments{params.TransactionAttachment}
			}

			err = tx.Create(&transaction).Error
			resp.PaymentTransaction = &transaction
			return err
		})

		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}

	}

	var cartItemIDs []string
	for _, cartItem := range cartItems {
		cartItemIDs = append(cartItemIDs, cartItem.ID)
	}

	err = r.db.Model(&models.InquiryCartItem{}).Where("id IN ?", cartItemIDs).Updates(&models.InquiryCartItem{CheckoutSessionID: checkoutSessionID, WaitingForCheckout: values.Bool(false)}).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&models.PurchaseOrder{}).Where("id IN ?", orderIDs).Updates(&models.PurchaseOrder{CheckoutSessionID: checkoutSessionID}).Error
	if err != nil {
		return nil, err
	}

	resp.Orders = orders
	resp.CheckoutSessionID = checkoutSessionID

	return &resp, nil
}

func (r *InquiryRepo) MultiInquiryCheckoutOneOrder(params InquiryCheckoutParams) (*models.PurchaseOrder, error) {
	var cartItems = []*models.InquiryCartItem{}
	var err = r.db.Find(&cartItems, "inquiry_id = ? AND (checkout_session_id IS NULL OR checkout_session_id = ?)", params.InquiryID, "").Error
	if err != nil {
		return nil, err
	}

	purchaseOrder, err := r.InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		InquiryID:     params.InquiryID,
		PaymentType:   params.PaymentType,
		CartItems:     cartItems,
		UserID:        params.GetUserID(),
	})
	if err != nil {
		return nil, err
	}

	return purchaseOrder, err

}

type MultiInquiryCheckoutInfoParams struct {
	models.JwtClaimsInfo
	CheckoutSessionID string `json:"checkout_session_id" params:"checkout_session_id" query:"checkout_session_id" validate:"required"`
}

type MultiInquiryCheckoutInfoResponse struct {
	Orders             []*models.PurchaseOrder    `json:"orders"`
	PaymentTransaction *models.PaymentTransaction `json:"payment_transaction"`
}

func (r *InquiryRepo) MultiInquiryCheckoutInfo(params MultiInquiryCheckoutInfoParams) (*MultiInquiryCheckoutInfoResponse, error) {
	var orders []*models.PurchaseOrder

	var builder = queryfunc.NewPurchaseOrderBuilder(queryfunc.PurchaseOrderBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IsConsistentRead:          true,
		IncludeCartItems:          true,
		IncludePaymentTransaction: true,
	})

	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("po.checkout_session_id = ?", params.CheckoutSessionID)
		}).
		FindFunc(&orders)

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	for _, order := range orders {
		if order.Inquiry != nil && order.Inquiry.ShippingAddressID != "" {
			order.Inquiry.ShippingAddress, _ = NewAddressRepo(r.db).GetAddress(GetAddressParams{
				AddressID: order.Inquiry.ShippingAddressID,
			})
		}
	}

	var resp = MultiInquiryCheckoutInfoResponse{
		Orders: orders,
	}

	return &resp, nil
}

type UpdateInquiryLogsParams struct {
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" param:"inquiry_id" validate:"required"`
	LogID     string `json:"log_id" param:"log_id" validate:"required"`

	Notes       string             `json:"notes" validate:"required"`
	Attachments models.Attachments `json:"attachments" validate:"required"`
}

func (r *InquiryRepo) UpdateInquiryLogs(params UpdateInquiryLogsParams) (*models.InquiryAudit, error) {
	var inquiryAudit = models.InquiryAudit{
		Notes:       params.Notes,
		Attachments: params.Attachments,
	}

	var err = r.db.Model(&models.InquiryAudit{}).
		Where("inquiry_id = ? AND id = ?", params.InquiryID, params.LogID).
		Updates(&inquiryAudit).Error

	return &inquiryAudit, err
}

type DeleteInquiryLogsParams struct {
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" param:"inquiry_id" validate:"required"`
	LogID     string `json:"log_id" param:"log_id" validate:"required"`
}

func (r *InquiryRepo) DeleteInquiryLogs(params DeleteInquiryLogsParams) (*models.InquiryAudit, error) {
	var inquiryAudit = models.InquiryAudit{
		Notes:       "",
		Attachments: nil,
	}

	var err = r.db.Model(&models.InquiryAudit{}).
		Select("Notes", "Attachments").
		Where("inquiry_id = ? AND id = ?", params.InquiryID, params.LogID).
		Updates(&inquiryAudit).Error

	return &inquiryAudit, err
}

type BuyerConfirmInquiryParams struct {
	models.JwtClaimsInfo

	InquiryID string                              `param:"inquiry_id" json:"inquiry_id" validate:"required"`
	Items     []*models.InquiryCartItemCreateForm `json:"items" validate:"required"`
}

func (r *InquiryRepo) BuyerConfirmInquiry(params BuyerConfirmInquiryParams) (*models.Inquiry, error) {
	var inquiry models.Inquiry
	var err = r.db.Select("ID", "UserID", "AdminQuotations").First(&inquiry, "id = ?", params.InquiryID).Error
	if err != nil {
		return nil, err
	}
	var samplePrice = inquiry.GetSampleUnitPrice()
	if samplePrice.GreaterThan(0) {
		return nil, errs.ErrInquirySampleOrderIsNotPaid
	}

	var purchaseOrder models.PurchaseOrder
	err = r.db.Select("ID", "UserID", "ReferenceID", "Status").First(&purchaseOrder, "inquiry_id = ?", params.InquiryID).Error
	if err != nil {
		return nil, err
	}

	var admins []*models.User
	err = r.db.Select("ID", "Name").First(&admins, "id IN ?", []string(inquiry.AssigneeIDs)).Error
	if err != nil {
		return nil, err
	}

	var itemIDs []string
	var items = lo.Map(params.Items, func(item *models.InquiryCartItemCreateForm, index int) *models.InquiryCartItem {
		var cartItem models.InquiryCartItem
		var err = copier.Copy(&cartItem, item)
		if err != nil {
			return nil
		}

		cartItem.ID = helper.GenerateXID()
		cartItem.InquiryID = params.InquiryID
		cartItem.UnitPrice = inquiry.GetSampleUnitPrice()
		cartItem.TotalPrice = item.UnitPrice.MultipleInt(int64(item.Qty))

		itemIDs = append(itemIDs, cartItem.ID)
		return &cartItem
	})

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&items).Error
		if err != nil {
			return err
		}

		var updates = models.PurchaseOrder{
			Status:       enums.PurchaseOrderStatusPaid,
			MarkAsPaidAt: values.Int64(time.Now().Unix()),
			TransferedAt: values.Int64(time.Now().Unix()),
			Currency:     inquiry.Currency,
			CartItemIDs:  itemIDs,
		}
		updates.TaxPercentage = inquiry.TaxPercentage
		var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
		if sqlResult.Error != nil {
			return eris.Wrap(sqlResult.Error, sqlResult.Error.Error())
		}

		if sqlResult.RowsAffected == 0 {
			return eris.New("Purchase order not found")
		}

		sqlResult = tx.Model(&models.Inquiry{}).Where("id = ?", params.InquiryID).UpdateColumn("Status", enums.InquiryStatusFinished)
		if sqlResult.Error != nil {
			return eris.Wrap(sqlResult.Error, sqlResult.Error.Error())
		}

		if sqlResult.RowsAffected == 0 {
			return eris.New("Inquiry not found")
		}

		var audits = []models.InquiryAudit{
			{
				InquiryID:       inquiry.ID,
				ActionType:      enums.AuditActionTypeInquirySamplePoCreated,
				UserID:          inquiry.UserID,
				Description:     fmt.Sprintf("New sample PO %s has been created for inquiry", purchaseOrder.ReferenceID),
				PurchaseOrderID: purchaseOrder.ID,
			},
		}

		if len(admins) > 0 {
			audits = append(audits, models.InquiryAudit{
				InquiryID:   inquiry.ID,
				ActionType:  enums.AuditActionTypeInquiryAdminMarkAsPaid,
				UserID:      admins[0].ID,
				Description: fmt.Sprintf("Admin %s has confirmed the payment", admins[0].Name),
			})
		}

		err = tx.Create(&audits).Error

		return err

	})
	if err != nil {
		return nil, err
	}

	inquiry.PurchaseOrder = &purchaseOrder

	return &inquiry, err

}

func (r *InquiryRepo) ApproveMultipleInquiryQuotations(req *models.ApproveMultipleInquiryQuotationsRequest) ([]*models.Inquiry, error) {
	var inquiries = make(models.Inquiries, 0, len(req.InquiryIDs))
	if err := r.db.Find(&inquiries, "id IN ? ", req.InquiryIDs).Error; err != nil {
		return nil, err
	}
	var dbInquiryIDs = inquiries.IDs()
	for _, id := range req.InquiryIDs {
		if !helper.StringContains(dbInquiryIDs, id) {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "inquiry_id:%s", id)
		}
	}
	for _, iq := range inquiries {
		if iq.BuyerQuotationStatus != enums.InquirySkuStatusWaitingForApproval {
			return nil, eris.Wrapf(errs.ErrInquiryInvalidToApproveQuotation, "inquiry_id:%s", iq.ID)
		}
	}

	var users models.Users
	if err := r.db.Find(&users, "id IN ?", inquiries.UserIDs()).Error; err != nil {
		return nil, err
	}
	var mapUserIDToUser = make(map[string]*models.User, len(users))
	for _, user := range users {
		mapUserIDToUser[user.ID] = user
	}

	var inquiriesToUpdate = make(models.Inquiries, 0, len(inquiries))
	for _, iq := range inquiries {
		iq.BuyerQuotationStatus = enums.InquirySkuStatusApproved
		iq.ApproveRejectMeta = &models.InquiryApproveRejectMeta{}
		iq.QuotationApprovedAt = values.Int64(time.Now().Unix())

		// inject additional attributes to return
		user := mapUserIDToUser[iq.UserID]
		iq.User = user

		inquiriesToUpdate = append(inquiriesToUpdate, iq)
	}

	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true}).Create(&inquiries).Error; err != nil {
		return nil, err
	}

	return inquiriesToUpdate, nil
}
