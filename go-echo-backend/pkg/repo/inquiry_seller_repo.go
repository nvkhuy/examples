package repo

import (
	"database/sql"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InquirySellerRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewInquirySellerRepo(db *db.DB) *InquirySellerRepo {
	return &InquirySellerRepo{
		db:     db,
		logger: logger.New("repo/InquirySeller"),
	}
}

type PaginateInquirySellerParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" query:"inquiry_id" form:"inquiry_id" param:"inquiry_id"`

	UserID string `json:"user_id" query:"user_id" form:"user_id"`

	DateFrom         int64  `json:"date_from" query:"date_from" form:"date_from"`
	DateTo           int64  `json:"date_to" query:"date_to" form:"date_to"`
	OrderReferenceID string `json:"order_reference_id" query:"order_reference_id" form:"order_reference_id"`

	Statuses []enums.InquirySellerStatus `json:"statuses" query:"statuses" form:"statuses"`

	IncludeInquiry            bool `json:"-"`
	IncludeUnseenCommentCount bool `json:"-"`
}

func (r *InquirySellerRepo) PaginateInquirySellerRequest(params PaginateInquirySellerParams) *query.Pagination {
	var builder = queryfunc.NewInquirySellerRequestBuilder(queryfunc.InquirySellerRequestBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeInquiry:            params.IncludeInquiry,
		IncludeUnseenCommentCount: params.IncludeUnseenCommentCount,
		CurrentUserID:             params.GetUserID(),
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.InquiryID != "" {
				builder.Where("rq.inquiry_id = ?", params.InquiryID)
			}

			if params.UserID != "" {
				builder.Where("rq.user_id = ?", params.UserID)
			}
			if params.DateFrom > 0 {
				builder.Where("rq.created_at >= ?", params.DateFrom)
			}
			if params.DateTo > 0 {
				builder.Where("rq.created_at <= ?", params.DateTo)
			}
			if len(params.Statuses) > 0 {
				builder.Where("rq.status IN ?", params.Statuses)
			}

			if strings.TrimSpace(params.OrderReferenceID) != "" {
				var q = "%" + params.OrderReferenceID + "%"
				builder.Where("rq.order_reference_id ILIKE @query_po", sql.Named("query_po", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type PaginateMatchingSellersParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	SellerID      string
	ProductGroups []string `json:"product_groups" query:"product_groups" form:"product_groups"`
	ProductTypes  []string `json:"product_types" query:"product_types" form:"product_types"`
	FabricTypes   []string `json:"fabric_types" query:"fabric_types" form:"fabric_types"`
}

// Current only get matching manufacturer seller
func (r *InquirySellerRepo) PaginateMatchingSellers(params PaginateMatchingSellersParams) *query.Pagination {
	var builder = queryfunc.NewInquirySellerMatchingBuilder(queryfunc.InquirySellerMatchingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: enums.RoleSeller,
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("role = ?", enums.RoleSeller)
			builder.Where("u.supplier_type = ?", enums.SupplierTypeManufacturer)

			if len(params.FabricTypes) > 0 {
				builder.Where("count_elements(bu.excepted_fabric_types,?) = 0", pq.StringArray(params.FabricTypes))
			}
			if len(params.ProductTypes) > 0 {
				builder.Where("count_elements(bu.product_types,?) > 0", pq.StringArray(params.ProductTypes))
			}
			if len(params.ProductGroups) > 0 {
				builder.Where("count_elements(bu.product_groups,?) > 0", pq.StringArray(params.ProductGroups))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *InquirySellerRepo) GetInquirySellerRequestByID(requestID string, options queryfunc.InquirySellerRequestBuilderOptions) (*models.InquirySeller, error) {
	var builder = queryfunc.NewInquirySellerRequestBuilder(options)
	var inquiryRequest models.InquirySeller
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("rq.id = ?", requestID)
		}).
		FirstFunc(&inquiryRequest)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &inquiryRequest, nil
}

func (r *InquirySellerRepo) FindInquirySeller(params PaginateInquirySellerParams) (*models.InquirySeller, error) {
	var builder = queryfunc.NewInquirySellerRequestBuilder(queryfunc.InquirySellerRequestBuilderOptions{})
	var inquiryRequest models.InquirySeller
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.InquiryID != "" {
				builder.Where("rq.inquiry_id = ?", params.InquiryID)
			}
			if params.UserID != "" {
				builder.Where("rq.user_id = ?", params.UserID)
			}
		}).
		FirstFunc(&inquiryRequest)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &inquiryRequest, nil
}

func (r *InquirySellerRepo) SellerCreateInquiryQuotation(form models.InquirySellerCreateQuatationParams) (*models.InquirySeller, error) {
	var inquirySeller models.InquirySeller
	var err = r.db.Find(&inquirySeller, "id = ?", form.InquirySellerID).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var inquirySellerUpdates models.InquirySeller
	err = copier.Copy(&inquirySellerUpdates, &form)
	if err != nil {
		return nil, err
	}

	inquirySellerUpdates.QuotationAt = values.Int64(time.Now().Unix())
	inquirySellerUpdates.Status = enums.InquirySellerStatusWaitingForApproval

	r.db.Clauses(clause.OnConflict{DoNothing: true}).Where("id = ?", inquirySeller.ID).Updates(inquirySellerUpdates)

	return r.GetInquirySellerRequestByID(form.InquirySellerID, queryfunc.InquirySellerRequestBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: form.GetRole(),
		},
	})
}

type InquirySellerStatusCountParams struct {
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" query:"inquiry_id" form:"inquiry_id" param:"inquiry_id"`
}

type InquirySellerStatusCountResponse struct {
	Quoted  int64 `json:"quoted"`
	Pending int64 `json:"pending"`
}

func (r *InquirySellerRepo) StatusCount(params InquirySellerStatusCountParams) (*InquirySellerStatusCountResponse, error) {
	var result InquirySellerStatusCountResponse
	var pendingStatuses = []enums.InquirySkuStatus{
		enums.InquirySkuStatusNew,
		enums.InquirySkuStatusWaitingForQuotation,
		enums.InquirySkuStatusRejected,
	}

	var quotedStatuses = []enums.InquirySkuStatus{
		enums.InquirySkuStatusWaitingForApproval,
	}
	r.db.Raw(`
	SELECT 
	COUNT(1) FILTER (WHERE inquiry_id = @inquiry_id AND status IN @pending_status) AS pending,
	COUNT(1) FILTER (WHERE inquiry_id = @inquiry_id AND status IN @quoted_status) AS quoted
	FROM inquiry_sellers
	WHERE inquiry_id = @inquiry_id
	GROUP BY inquiry_id = @inquiry_id
	`, sql.Named("inquiry_id", params.InquiryID), sql.Named("pending_status", pendingStatuses), sql.Named("quoted_status", quotedStatuses)).
		Scan(&result)

	return &result, nil
}

type InquirySellerDesignCommentMarkSeenParams struct {
	models.JwtClaimsInfo

	InquirySellerID string `json:"inquiry_seller_id" param:"inquiry_seller_id"`
	FileKey         string `json:"file_key" param:"file_key"`
}

func (r *InquirySellerRepo) InquiryDesignCommentMarkSeen(params InquirySellerDesignCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ? AND file_key = ?", enums.CommentTargetTypeInquirySellerDesign, params.InquirySellerID, params.FileKey).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

type InquirySellerDesignCommenStatusCountParams struct {
	models.JwtClaimsInfo
	InquirySellerID string `json:"inquiry_seller_id" param:"inquiry_seller_id"`
}

type InquirySellerDesignCommenStatusCountItem struct {
	FileKey     string `json:"file_key"`
	UnseenCount int64  `json:"unseen_count"`
}

func (r *InquirySellerRepo) InquiryDesignCommentStatusCount(params InquirySellerDesignCommenStatusCountParams) ([]*InquirySellerDesignCommenStatusCountItem, error) {
	inquirySeller, err := r.GetInquirySellerRequestByID(params.InquirySellerID, queryfunc.InquirySellerRequestBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{},
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	inquiry, err := NewInquiryRepo(r.db).GetInquiryByID(GetInquiryByIDParams{
		InquiryID:     inquirySeller.InquiryID,
		JwtClaimsInfo: params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var items []*InquirySellerDesignCommenStatusCountItem

	if inquiry.TechpackAttachments != nil && len(*inquiry.TechpackAttachments) > 0 {
		for _, fileItem := range *inquiry.TechpackAttachments {
			count := int64(0)
			if fileItem.FileKey != "" {
				r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND file_key = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypePurchaseOrderDesign, params.InquirySellerID, fileItem.FileKey).Count(&count)
				items = append(items, &InquirySellerDesignCommenStatusCountItem{
					FileKey:     fileItem.FileKey,
					UnseenCount: count,
				})
			}

		}
	}

	return items, nil
}

type InquirySellerQuotationCommentMarkSeenParams struct {
	models.JwtClaimsInfo

	InquirySellerID string `json:"inquiry_seller_id" param:"inquiry_seller_id"`
}

func (r *InquirySellerRepo) InquiryQuotationCommentMarkSeen(params InquirySellerQuotationCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ?", enums.CommentTargetTypeInquirySellerRequest, params.InquirySellerID).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

type InquirySellerAllocationSearchSellerParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	InquiryID string `json:"inquiry_id" query:"inquiry_id" form:"inquiry_id" param:"inquiry_id"`
}

func (r *InquirySellerRepo) InquirySellerAllocationSearchSeller(params InquirySellerAllocationSearchSellerParams) *query.Pagination {
	var builder = queryfunc.NewInquirySellerAllocationBuilder(queryfunc.InquirySellerAllocationBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("(iqs.id IS NULL OR iqs.inquiry_id = ?)", params.InquiryID)
			builder.Where("u.role = ?", enums.RoleSeller)

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(u.name ILIKE @keyword OR u.company_name ILIKE @keyword OR u.email ILIKE @keyword)", sql.Named("keyword", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result

}

type AdminInquirySellerApproveQuotationParams struct {
	models.JwtClaimsInfo

	InquirySellerID string `json:"inquiry_seller_id" param:"inquiry_seller_id" validate:"required"`
	Note            string `json:"note"`
}

func (r *InquirySellerRepo) AdminInquirySellerApproveQuotation(params AdminInquirySellerApproveQuotationParams) (*models.InquirySeller, error) {
	var inquirySeller models.InquirySeller
	var err = r.db.Select("ID", "UserID", "InquiryID", "PurchaseOrderID", "FabricCost", "DecorationCost", "MakingCost", "OtherCost").First(&inquirySeller, "id = ?", params.InquirySellerID).Error
	if err != nil {
		return nil, err
	}

	var inquiry models.Inquiry
	err = r.db.First(&inquiry, "id = ?", inquirySeller.InquiryID).Error
	if err != nil {
		return nil, err
	}

	var purchaseOrder *models.PurchaseOrder
	if inquirySeller.PurchaseOrderID == "" {
		var po models.PurchaseOrder
		err = r.db.Model(&models.PurchaseOrder{}).First(&po, "inquiry_id = ?", inquirySeller.InquiryID).Error
		if err != nil {
			if r.db.IsRecordNotFoundError(err) {
				purchaseOrder, err = NewPurchaseOrderRepo(r.db).CreatePurchaseOrder(CreatePurchaseOrderParams{
					UserID:        inquiry.UserID,
					Currency:      inquiry.Currency,
					SizeChart:     enums.InquirySizeChart(inquiry.SizeChart),
					Quotations:    inquiry.AdminQuotations,
					TaxPercentage: inquiry.TaxPercentage,
				})
				if err != nil {
					return nil, err
				}
				goto Next
			}
			return nil, err
		}
		purchaseOrder = &po
	} else {
		purchaseOrder, err = NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
			PurchaseOrderID: inquirySeller.PurchaseOrderID,
		})
		if err != nil {
			return nil, err
		}

	}

Next:
	var totalPrice price.Price
	if inquirySeller.SampleUnitPrice == nil {
		inquirySeller.SampleUnitPrice = inquirySeller.GetSampleUnitPrice().ToPtr()
		r.db.Model(&models.InquirySeller{}).Where("id = ?", inquirySeller.ID).UpdateColumn("SampleUnitPrice", inquirySeller.SampleUnitPrice)
	}

	if len(purchaseOrder.Items) > 0 {
		for _, cartItem := range purchaseOrder.Items {
			totalPrice = totalPrice.Add(inquirySeller.SampleUnitPrice.MultipleInt(cartItem.Quantity))
		}
	} else if len(purchaseOrder.OrderCartItems) > 0 {
		for _, cartItem := range purchaseOrder.OrderCartItems {
			totalPrice = totalPrice.Add(inquirySeller.SampleUnitPrice.MultipleInt(cartItem.Qty))
		}
	} else if len(purchaseOrder.CartItems) > 0 {
		for _, cartItem := range purchaseOrder.CartItems {
			totalPrice = totalPrice.Add(inquirySeller.SampleUnitPrice.MultipleInt(cartItem.Qty))
		}
	}

	var updates = models.InquirySeller{
		Status:          enums.InquirySellerStatusApproved,
		Note:            params.Note,
		PurchaseOrderID: purchaseOrder.ID,
		SampleUnitPrice: inquirySeller.GetSampleUnitPrice().ToPtr(),
	}

	var poUpdates = models.PurchaseOrder{
		SampleMakerID:        inquirySeller.UserID,
		SellerTrackingStatus: enums.SellerPoTrackingStatusNew,
		InquiryID:            inquirySeller.InquiryID,
	}
	poUpdates.SellerTotalPrice = &totalPrice

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Model(&models.InquirySeller{}).Where("id = ?", params.InquirySellerID).Updates(&updates)
		if sqlResult.RowsAffected == 0 {
			return errs.ErrInquirySellerInvalidToApprove
		}

		var err = tx.Model(&models.InquirySeller{}).Where("id <> ? AND purchase_order_id = ?", inquirySeller.ID, purchaseOrder.ID).UpdateColumn("Status", enums.InquirySellerStatusRejected).Error
		if err != nil {
			return err
		}

		sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&poUpdates)
		if sqlResult.RowsAffected == 0 {
			return errs.ErrPONotFound
		}

		return sqlResult.Error
	})

	return &updates, nil
}

type InquirySellerApproveOfferParams struct {
	models.JwtClaimsInfo

	InquirySellerID string `json:"inquiry_seller_id" param:"inquiry_seller_id" validate:"required"`
	Note            string `json:"note"`
}

func (r *InquirySellerRepo) InquirySellerApproveOffer(params InquirySellerApproveOfferParams) error {
	var updates = models.InquirySeller{
		Status: enums.InquirySellerStatusWaitingForQuotation,
		Note:   params.Note,
	}

	var sqlResult = r.db.Model(&models.InquirySeller{}).
		Where("id = ? AND user_id = ?", params.InquirySellerID, params.GetUserID()).
		Updates(&updates)
	if sqlResult.RowsAffected == 0 {
		return errs.ErrInquirySellerInvalidToApprove
	}
	return sqlResult.Error
}

type AdminInquirySellerRejectQuotationParams struct {
	models.JwtClaimsInfo

	InquirySellerID string       `json:"inquiry_seller_id" param:"inquiry_seller_id" validate:"required"`
	RejectReason    string       `json:"reject_reason"`
	ExpectedPrice   *price.Price `json:"expected_price"`
}

func (r *InquirySellerRepo) AdminInquirySellerRejectQuotation(params AdminInquirySellerRejectQuotationParams) error {
	var updates = models.InquirySeller{
		Status:            enums.InquirySellerStatusRejected,
		AdminRejectReason: params.RejectReason,
		ExpectedPrice:     params.ExpectedPrice,
	}

	var sqlResult = r.db.Model(&models.InquirySeller{}).Where("id = ?", params.InquirySellerID).Updates(&updates)
	if sqlResult.RowsAffected == 0 {
		return errs.ErrInquirySellerInvalidToReject
	}
	return sqlResult.Error
}

type InquirySellerRejectOfferParams struct {
	models.JwtClaimsInfo

	InquirySellerID string       `json:"inquiry_seller_id" param:"inquiry_seller_id" validate:"required"`
	RejectReason    string       `json:"reject_reason"`
	ExpectedPrice   *price.Price `json:"expected_price"`
}

func (r *InquirySellerRepo) InquirySellerRejectOffer(params InquirySellerRejectOfferParams) error {
	var updates = models.InquirySeller{
		Status:            enums.InquirySellerStatusOfferRejected,
		AdminRejectReason: params.RejectReason,
		ExpectedPrice:     params.ExpectedPrice,
	}

	var sqlResult = r.db.Model(&models.InquirySeller{}).
		Where("id = ? AND user_id = ?", params.InquirySellerID, params.GetUserID()).
		Updates(&updates)
	if sqlResult.RowsAffected == 0 {
		return errs.ErrInquirySellerInvalidToReject
	}
	return sqlResult.Error
}

func (r *InquirySellerRepo) SubmitMultipleInquirySellerQuotations(req *models.SubmitMultipleInquirySellerQuotationRequest) ([]*models.InquirySeller, error) {
	var inquirySellerIDs = make([]string, 0, len(req.Quotations))
	for _, quotation := range req.Quotations {
		inquirySellerIDs = append(inquirySellerIDs, quotation.InquirySellerID)
	}
	var inquirySellers models.InquirySellers
	if err := r.db.Find(&inquirySellers, "id IN ? AND user_id = ?", inquirySellerIDs, req.GetUserID()).Error; err != nil {
		return nil, err
	}
	var dbInquirySellerIDs = inquirySellers.IDs()
	for _, id := range inquirySellerIDs {
		if !helper.StringContains(dbInquirySellerIDs, id) {
			return nil, eris.Wrapf(errs.ErrInquirySellerNotFound, "inquiry_seller_id:%s", id)
		}
	}

	var validStatus = []string{enums.InquirySellerStatusWaitingForQuotation.String(), enums.InquirySellerStatusWaitingForApproval.String(), enums.InquirySellerStatusRejected.String()}
	for idx, iq := range inquirySellers {
		if !helper.StringContains(validStatus, string(iq.Status)) {
			return nil, eris.Wrapf(errs.ErrInquirySellerInvalidToSendQuotationToBuyer, "index:%d", idx)
		}
	}
	var inquiries = make(models.Inquiries, 0, len(inquirySellers))
	if err := r.db.Select("ID", "AssigneeIDs").Find(&inquiries, "id IN ?", inquirySellers.InquiryIDs()).Error; err != nil {
		return nil, err
	}
	var mapInquiryIDtoInquiry = make(map[string]*models.Inquiry, len(inquiries))
	for _, iq := range inquiries {
		mapInquiryIDtoInquiry[iq.ID] = iq
	}

	var inquirySellersToUpdate = make(models.InquirySellers, 0, len(req.Quotations))

	for _, quotation := range req.Quotations {
		inquirySeller, _ := lo.Find(inquirySellers, func(item *models.InquirySeller) bool {
			return item.ID == quotation.InquirySellerID
		})
		inquiry, ok := mapInquiryIDtoInquiry[inquirySeller.InquiryID]
		if !ok {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "inquiry_seller_id:%s", inquirySeller.ID)
		}

		inquirySeller.FabricCost = quotation.FabricCost
		inquirySeller.DecorationCost = quotation.DecorationCost
		inquirySeller.MakingCost = quotation.MakingCost
		inquirySeller.OtherCost = quotation.OtherCost
		inquirySeller.SampleUnitPrice = quotation.SampleUnitPrice
		inquirySeller.SampleLeadTime = quotation.SampleLeadTime
		inquirySeller.SellerRemark = quotation.SellerRemark
		inquirySeller.StartProductionDate = quotation.StartProductionDate
		inquirySeller.CapacityPerDay = quotation.CapacityPerDay
		inquirySeller.BulkQuotations = quotation.BulkQuotations
		inquirySeller.QuotationAt = values.Int64(time.Now().Unix())
		inquirySeller.Status = enums.InquirySellerStatusWaitingForApproval

		inquirySeller.Inquiry = inquiry

		inquirySellersToUpdate = append(inquirySellersToUpdate, inquirySeller)
	}
	if err := r.db.
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, UpdateAll: true}).
		Create(inquirySellersToUpdate).Error; err != nil {

		return nil, eris.Wrap(err, err.Error())
	}

	return inquirySellersToUpdate, nil
}
