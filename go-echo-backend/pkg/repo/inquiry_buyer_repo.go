package repo

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/clause"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/thaitanloi365/go-utils"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
)

type InquiryBuyerRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewInquiryBuyerRepo(db *db.DB) *InquiryBuyerRepo {
	return &InquiryBuyerRepo{
		db:     db,
		logger: logger.New("repo/InquiryBuyer"),
	}
}

type PaginateInquiryBuyerParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	InquiryID       string                `json:"inquiry_id" query:"inquiry_id" form:"inquiry_id" param:"inquiry_id"`
	UserID          string                `json:"user_id" query:"user_id" form:"user_id" param:"user_id"`
	TeamID          string                `json:"team_id" query:"team_id" form:"team_id" param:"team_id"`
	Statuses        []enums.InquiryStatus `json:"statuses" query:"statuses" form:"statuses" param:"statuses"`
	ExcludeStatuses []enums.InquiryStatus `json:"exclude_statuses" query:"exclude_statuses" form:"statuses" param:"exclude_statuses"`

	DateFrom int64 `json:"date_from" query:"date_from" form:"date_from"`
	DateTo   int64 `json:"date_to" query:"date_to" form:"date_to"`

	IncludeInquiry    bool `json:"-"`
	IncludeCollection bool `json:"-"`
	IsQueryAll        bool `json:"-"`
}

func (r *InquiryBuyerRepo) PaginateInquiry(params PaginateInquiryBuyerParams) *query.Pagination {
	var userID = params.GetUserID()
	if params.TeamID != "" && !params.GetRole().IsAdmin() && !params.IsQueryAll {
		if err := r.db.Select("ID").First(&models.BrandTeam{}, "team_id = ? AND user_id = ?", params.TeamID, userID).Error; err != nil {
			return &query.Pagination{
				Records: []*models.PurchaseOrder{},
			}
		}
		userID = params.TeamID
	}

	var builder = queryfunc.NewInquiryBuyerBuilder(queryfunc.InquiryBuyerBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeInquiry:    params.IncludeInquiry,
		IncludeCollection: params.IncludeCollection,
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

			if params.UserID != "" {
				builder.Where("iq.user_id = ?", params.UserID)
			}

			if params.DateFrom > 0 {
				builder.Where("iq.created_at >= ?", params.DateFrom)
			}

			if params.DateTo > 0 {
				builder.Where("iq.created_at <= ?", params.DateTo)
			}

			if params.Statuses != nil {
				builder.Where("iq.buyer_quotation_status IN ?", params.Statuses)
			}

			if params.ExcludeStatuses != nil {
				builder.Where("iq.status NOT IN ?", params.ExcludeStatuses)
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"

				if strings.HasPrefix(keyword, "IQ-") {
					builder.Where("iq.reference_id ILIKE @keyword", sql.Named("keyword", q))
				} else if strings.HasPrefix(keyword, "COL-") {
					builder.Where("iq.order_group_id ILIKE @keyword", sql.Named("keyword", q))
				} else {
					builder.Where("(iq.title ILIKE @keyword OR iq.sku_note ILIKE @keyword)", sql.Named("keyword", q))
				}

			}

		}).
		OrderBy("iq.order_group_id ASC, iq.updated_at DESC").
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *InquiryBuyerRepo) CreateInquiry(form models.InquiryCreateForm) (*models.Inquiry, error) {
	managerID := NewUserRepo(r.db).GetTeamManagerID(form.GetUserID())
	var err = r.ValidateTeamMemberCreateAction(form.GetUserID(), managerID)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.Select("ID", "Name", "Avatar", "ContactOwnerIDs").First(&user, "id = ?", managerID).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var inquiry models.Inquiry
	err = copier.Copy(&inquiry, &form)
	if err != nil {
		return nil, err
	}

	inquiry.AssigneeIDs = user.ContactOwnerIDs
	inquiry.UserID = managerID

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

	var settingIQ *models.SettingInquiry
	settingIQ, err = NewSettingInquiryRepo(r.db).Get(GetSettingInquiryParams{
		Type: enums.SettingInquiryEditTimeoutType,
	})
	if err == nil && settingIQ != nil {
		inquiry.EditTimeout = aws.Int64(time.Now().Unix() + settingIQ.EditTimeout)
	}

	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		inquiry.ID = helper.GenerateXID()
		err = tx.Create(&inquiry).Error
		if err != nil {
			return err
		}

		var purchaseOrder = models.PurchaseOrder{
			UserID:              inquiry.UserID,
			ProductName:         inquiry.Title,
			InquiryID:           inquiry.ID,
			Attachments:         inquiry.Attachments,
			Document:            inquiry.Document,
			Design:              inquiry.Design,
			FabricAttachments:   inquiry.FabricAttachments,
			TechpackAttachments: inquiry.TechpackAttachments,
			Currency:            inquiry.Currency,
			AssigneeIDs:         inquiry.AssigneeIDs,
			OrderGroupID:        inquiry.OrderGroupID,
			ShippingAddressID:   inquiry.ShippingAddressID,
		}
		err = tx.Create(&purchaseOrder).Error
		return err
	})
	if err != nil {
		return nil, err
	}

	inquiry.User = &user

	return &inquiry, nil
}

func (r *InquiryBuyerRepo) CreateMultipleInquiries(req *models.CreateMultipleInquiriesRequest) ([]*models.Inquiry, error) {
	managerID := NewUserRepo(r.db).GetTeamManagerID(req.GetUserID())
	var err = r.ValidateTeamMemberCreateAction(req.GetUserID(), managerID)
	if err != nil {
		return nil, err
	}
	if err := r.validateCreateMultipleInquiries(req); err != nil {
		return nil, err
	}
	var user models.User
	err = r.db.Select("ID", "Name", "Avatar", "ContactOwnerIDs").First(&user, "id = ?", managerID).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var countryTax *models.SettingTax
	countryTax, err = NewSettingTaxRepo(r.db).GetAffectedSettingTax(models.GetAffectedSettingTaxForm{
		CurrencyCode: enums.VND,
	})
	if err != nil && !r.db.IsRecordNotFoundError(err) {
		return nil, eris.Wrap(err, err.Error())
	}

	var settingIQ *models.SettingInquiry
	settingIQ, err = NewSettingInquiryRepo(r.db).Get(GetSettingInquiryParams{
		Type: enums.SettingInquiryEditTimeoutType,
	})
	if err != nil && !r.db.IsRecordNotFoundError(err) {
		return nil, eris.Wrap(err, err.Error())
	}

	var inquiriesToCreate []*models.Inquiry
	var purchaseOrdersToCreate []*models.PurchaseOrder
	var shippingAddressesToCreate []*models.Address
	var coordinatesToCreate []*models.Coordinate

	mapAddressID := make(map[string]struct{})
	mapCoordinateID := make(map[string]struct{})

	for _, iqReq := range req.Inquiries {
		var inquiry models.Inquiry
		err = copier.Copy(&inquiry, &iqReq)
		if err != nil {
			return nil, err
		}
		inquiry.ID = helper.GenerateXID()
		inquiry.AssigneeIDs = user.ContactOwnerIDs
		inquiry.UserID = managerID
		if countryTax != nil {
			inquiry.TaxPercentage = aws.Float64(countryTax.TaxPercentage)
		}
		if settingIQ != nil {
			inquiry.EditTimeout = aws.Int64(time.Now().Unix() + settingIQ.EditTimeout)
		}
		inquiriesToCreate = append(inquiriesToCreate, &inquiry)

		assignInquiryShippingAddress(&inquiry, iqReq.ShippingAddress)
		if _, ok := mapAddressID[iqReq.ShippingAddress.ID]; !ok {
			shippingAddressesToCreate = append(shippingAddressesToCreate, iqReq.ShippingAddress)
			mapAddressID[iqReq.ShippingAddress.ID] = struct{}{}
		}
		if _, ok := mapCoordinateID[iqReq.ShippingAddress.Coordinate.ID]; !ok {
			coordinatesToCreate = append(coordinatesToCreate, iqReq.ShippingAddress.Coordinate)
			mapCoordinateID[iqReq.ShippingAddress.Coordinate.ID] = struct{}{}
		}

		var purchaseOrder = models.PurchaseOrder{
			UserID:              inquiry.UserID,
			InquiryID:           inquiry.ID,
			Attachments:         inquiry.Attachments,
			FabricAttachments:   inquiry.FabricAttachments,
			TechpackAttachments: inquiry.TechpackAttachments,
			Currency:            inquiry.Currency,
			AssigneeIDs:         inquiry.AssigneeIDs,
			OrderGroupID:        inquiry.OrderGroupID,
			ShippingAddressID:   inquiry.ShippingAddressID,
		}
		purchaseOrdersToCreate = append(purchaseOrdersToCreate, &purchaseOrder)
	}
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&coordinatesToCreate).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&shippingAddressesToCreate).Error; err != nil {
			return err
		}

		if err := tx.Create(&inquiriesToCreate).Error; err != nil {
			return err
		}
		if err := tx.Create(&purchaseOrdersToCreate).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	for _, iq := range inquiriesToCreate {
		iq.User = &user
	}

	return inquiriesToCreate, nil
}

func assignInquiryShippingAddress(inquiry *models.Inquiry, address *models.Address) {
	var dataStr = address.Coordinate.Coordinate.ToJsonString()
	address.Coordinate.ID = utils.MD5(dataStr)
	address.CoordinateID = address.Coordinate.ID
	address.ID = address.GenerateID()
	inquiry.ShippingAddressID = address.ID
}

func (r *InquiryBuyerRepo) validateCreateMultipleInquiries(req *models.CreateMultipleInquiriesRequest) error {
	acceptedCurrency := []string{string(enums.USD), string(enums.SGD), string(enums.VND)}
	acceptedFabricWeightUnit := []string{string(enums.FabricWeightUnitGSM), string(enums.FabricWeightUnitOZ)}
	var orderGroupIDs []string

	for idx, inquiry := range req.Inquiries {
		if inquiry.OrderGroupID != "" {
			orderGroupIDs = append(orderGroupIDs, inquiry.OrderGroupID)
		}
		switch {
		case inquiry.Title == "":
			return fmt.Errorf("title is required at inquiry %d", idx)
		case inquiry.Quantity < 100:
			return fmt.Errorf("minimum order quantity is 100 at inquiry %d", idx)
		case inquiry.SizeList == "":
			return fmt.Errorf("size list is required at inquiry %d", idx)
		case inquiry.Attachments == nil || len(*inquiry.Attachments) == 0:
			return fmt.Errorf("attachments is required at inquiry %d", idx)
		case !helper.StringContains(acceptedCurrency, string(inquiry.Currency)):
			return fmt.Errorf("currency is invalid at inquiry %d", idx)
		case inquiry.FabricWeightUnit != "" && !helper.StringContains(acceptedFabricWeightUnit, string(inquiry.FabricWeightUnit)):
			return fmt.Errorf("fabric weight unit is invalid at inquiry %d", idx)
		case inquiry.ShippingAddress == nil:
			return fmt.Errorf("shipping address is required at inquiry %d", idx)
		case inquiry.ShippingAddress.Name == "":
			return fmt.Errorf("shipping address name is required at inquiry %d", idx)
		case inquiry.ShippingAddress.PhoneNumber == "":
			return fmt.Errorf("shipping address phone number is required at inquiry %d", idx)
		case inquiry.ShippingAddress.Coordinate == nil:
			return fmt.Errorf("shipping address coordinate is required at inquiry %d", idx)
		}
	}
	if len(orderGroupIDs) > 0 {
		var dbOrderGroups models.OrderGroups
		if err := r.db.Select("ID").Find(&dbOrderGroups, "id IN ?", orderGroupIDs).Error; err != nil {
			return eris.Wrap(err, err.Error())
		}
		dbOrderGroupIDs := dbOrderGroups.IDs()

		for idx, id := range orderGroupIDs {
			if !helper.StringContains(dbOrderGroupIDs, id) {
				return fmt.Errorf("Order group not found at inquiry %d", idx)
			}
		}
	}

	return nil
}

type BuyerApproveInquiryQuotationParams struct {
	models.JwtClaimsInfo

	InquiryID         string                           `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`
	ApproveRejectMeta *models.InquiryApproveRejectMeta `json:"approve_reject_meta" param:"approve_reject_meta" query:"approve_reject_meta"`
}

func (r *InquiryBuyerRepo) ValidateTeamMemberCreateAction(staffId, managerId string) (err error) {
	if staffId == managerId {
		return
	}
	var member models.BrandTeam
	if err = r.db.First(&member, "user_id = ? and team_id = ?", staffId, managerId).Error; err != nil {
		return
	}

	for _, act := range member.Actions {
		if enums.BrandMemberAction(act) == enums.BrandMemberActionCreateRFQ {
			return
		}
	}
	err = errs.ErrInvalidActionBrandTeamMember
	return
}

func (r *InquiryBuyerRepo) ValidateTeamMemberUpdateAction(userId, inquiryId string, action enums.BrandMemberAction) (err error) {
	var iq models.Inquiry
	if err = r.db.Select("id", "user_id").First(&iq, "id = ?", inquiryId).Error; err != nil {
		return
	}
	if userId == iq.UserID {
		return
	}

	var member models.BrandTeam
	if err = r.db.First(&member, "user_id = ? and team_id = ?", userId, iq.UserID).Error; err != nil {
		return
	}

	for _, act := range member.Actions {
		if enums.BrandMemberAction(act) == action {
			return
		}
	}
	err = errs.ErrInvalidActionBrandTeamMember
	return
}

func (r *InquiryBuyerRepo) ApproveInquiryQuotation(params BuyerApproveInquiryQuotationParams) (*models.Inquiry, error) {
	if err := r.ValidateTeamMemberUpdateAction(params.GetUserID(), params.InquiryID, enums.BrandMemberActionApproveRFQ); err != nil {
		return nil, err
	}

	var inquiry models.Inquiry
	var err = r.db.First(&inquiry, "id = ? and user_id = ?", params.InquiryID, params.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	if inquiry.BuyerQuotationStatus != enums.InquirySkuStatusWaitingForApproval {
		return nil, errs.ErrInquiryInvalidToApproveQuotation
	}

	var user models.User
	err = r.db.Select("ID", "Name", "Email").First(&user, "id = ?", params.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	var updateColumns = []string{
		"BuyerQuotationStatus", "ApproveRejectMeta", "UpdateSeenAt", "QuotationApprovedAt",
	}
	var update = &models.Inquiry{
		BuyerQuotationStatus: enums.InquirySkuStatusApproved,
		ApproveRejectMeta:    params.ApproveRejectMeta,
		QuotationApprovedAt:  values.Int64(time.Now().Unix()),
		UpdateSeenAt:         nil,
	}

	for _, item := range inquiry.AdminQuotations {
		if item.Type == enums.InquiryTypeSample {
			item.Accepted = values.Bool(true)
		}

		update.AdminQuotations = append(update.AdminQuotations, item)
	}
	inquiry.AdminQuotations = update.AdminQuotations

	updateColumns = append(updateColumns, "AdminQuotations")

	err = r.db.Model(&models.Inquiry{}).
		Select(updateColumns).
		Where("id = ?", inquiry.ID).
		Updates(&update).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	inquiry.User = &user
	inquiry.BuyerQuotationStatus = update.BuyerQuotationStatus
	inquiry.ApproveRejectMeta = update.ApproveRejectMeta
	inquiry.ApproveRejectMeta = update.ApproveRejectMeta
	inquiry.UpdateSeenAt = update.UpdateSeenAt
	inquiry.AdminQuotations = update.AdminQuotations

	return &inquiry, err
}

type BuyerRejectInquiryQuotationParams struct {
	models.JwtClaimsInfo

	InquiryID         string                           `json:"inquiry_id" param:"inquiry_id" query:"inquiry_id" validate:"required"`
	ApproveRejectMeta *models.InquiryApproveRejectMeta `json:"approve_reject_meta" param:"approve_reject_meta" query:"approve_reject_meta"`
}

func (r *InquiryBuyerRepo) RejectInquiryQuotation(params BuyerRejectInquiryQuotationParams) (*models.Inquiry, error) {
	if err := r.ValidateTeamMemberUpdateAction(params.GetUserID(), params.InquiryID, enums.BrandMemberActionRejectRFQ); err != nil {
		return nil, err
	}

	var inquiry models.Inquiry
	var err = r.db.First(&inquiry, "id = ? and user_id = ?", params.InquiryID, params.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	if inquiry.BuyerQuotationStatus != enums.InquirySkuStatusWaitingForApproval {
		return nil, errs.ErrInquiryInvalidToApproveQuotation
	}

	var user models.User
	err = r.db.Select("ID", "Name", "Email").First(&user, "id = ?", params.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	var update = &models.Inquiry{
		BuyerQuotationStatus: enums.InquirySkuStatusRejected,
		ApproveRejectMeta:    params.ApproveRejectMeta,
		UpdateSeenAt:         nil,
	}
	if params.ApproveRejectMeta != nil {
		params.ApproveRejectMeta.CreatedAt = time.Now().Unix()
		var items = inquiry.ApproveRejectMetaItems
		items = append(items, params.ApproveRejectMeta)
		update.ApproveRejectMetaItems = items
	}

	err = r.db.Model(&models.Inquiry{}).
		Select("BuyerQuotationStatus", "ApproveRejectMeta", "ApproveRejectMetaItems", "UpdateSeenAt").
		Where("id = ?", inquiry.ID).
		Updates(&update).Error

	if err != nil {
		return &inquiry, eris.Wrap(err, err.Error())
	}

	inquiry.User = &user
	inquiry.BuyerQuotationStatus = update.BuyerQuotationStatus
	inquiry.ApproveRejectMeta = update.ApproveRejectMeta
	inquiry.UpdateSeenAt = update.UpdateSeenAt
	if len(update.ApproveRejectMetaItems) > 0 {
		inquiry.ApproveRejectMetaItems = update.ApproveRejectMetaItems
	}

	return &inquiry, err
}

type UpdateAttachmentsParams struct {
	models.JwtClaimsInfo
	InquiryID           string              `param:"inquiry_id" json:"inquiry_id" validate:"required"`
	Attachments         *models.Attachments `json:"attachments,omitempty"`
	Document            *models.Attachments `json:"document,omitempty"`
	Design              *models.Attachments `json:"design,omitempty"`
	FabricAttachments   *models.Attachments `json:"fabric_attachments,omitempty"`
	TechpackAttachments *models.Attachments `json:"techpack_attachments,omitempty"`
}

func (r *InquiryBuyerRepo) UpdateAttachments(params UpdateAttachmentsParams) (iq *models.Inquiry, err error) {
	if err = r.ValidateTeamMemberUpdateAction(params.GetUserID(), params.InquiryID, enums.BrandMemberActionUpdateRFQ); err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	currentInquiry, err := NewInquiryRepo(r.db).GetInquiryByID(GetInquiryByIDParams{InquiryID: params.InquiryID})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	if currentInquiry.EditTimeout != nil && *currentInquiry.EditTimeout < time.Now().Unix() {
		err = errs.ErrInquiryEditTimeout
		return nil, eris.Wrap(err, err.Error())
	}

	var update = &models.Inquiry{
		Attachments:         params.Attachments,
		Document:            params.Document,
		Design:              params.Design,
		FabricAttachments:   params.FabricAttachments,
		TechpackAttachments: params.TechpackAttachments,
	}
	var result models.Inquiry
	err = r.db.Model(&result).
		Clauses(clause.Returning{}).
		Where("id = ?", params.InquiryID).
		Updates(&update).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	iq = &result
	return
}

type BuyerInquiryCloseForm struct {
	models.JwtClaimsInfo

	InquiryID   string                           `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`
	CloseReason *models.InquiryApproveRejectMeta `json:"close_reason,omitempty"`
}

func (r *InquiryBuyerRepo) CloseInquiry(form BuyerInquiryCloseForm) error {
	userID := form.GetUserID()
	var statuses = []enums.InquiryStatus{
		enums.InquiryStatusNew,
		enums.InquiryStatusQuoteInProcess,
	}

	var updates = models.Inquiry{
		Status:      enums.InquiryStatusClosed,
		CloseReason: form.CloseReason,
	}
	var sqlResult = r.db.Unscoped().Model(&models.Inquiry{}).Where("user_id = ? AND id = ? AND status IN ?", userID, form.InquiryID, statuses).Updates(&updates)
	if sqlResult.Error != nil {
		return sqlResult.Error
	}

	if sqlResult.RowsAffected == 0 {
		return errs.ErrInquiryNotFound
	}
	return nil
}

type BuyerInquiryCancelForm struct {
	models.JwtClaimsInfo

	InquiryID    string                           `json:"inquiry_id" query:"inquiry_id" param:"inquiry_id" validate:"required"`
	CancelReason *models.InquiryApproveRejectMeta `json:"cancel_reason,omitempty"`
}

func (r *InquiryBuyerRepo) CancelInquiry(form BuyerInquiryCancelForm) error {
	userID := form.GetUserID()
	var statuses = []enums.InquiryStatus{
		enums.InquiryStatusNew,
	}

	var updates = models.Inquiry{
		Status:       enums.InquiryStatusCanceled,
		CancelReason: form.CancelReason,
	}
	var sqlResult = r.db.Unscoped().Model(&models.Inquiry{}).Where("user_id = ? AND id = ? AND status IN ?", userID, form.InquiryID, statuses).Updates(&updates)
	if sqlResult.Error != nil {
		return sqlResult.Error
	}

	if sqlResult.RowsAffected == 0 {
		return errs.ErrInquiryNotFound
	}
	return nil
}

func (r *InquiryBuyerRepo) RejectMultipleInquiryQuotations(req *models.RejectMultipleInquiryQuotationsRequest) ([]*models.Inquiry, error) {
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
		iq.BuyerQuotationStatus = enums.InquirySkuStatusRejected
		iq.ApproveRejectMeta = &models.InquiryApproveRejectMeta{
			Reason:    req.Reason,
			CreatedAt: time.Now().Unix(),
		}
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
