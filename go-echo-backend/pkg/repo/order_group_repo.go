package repo

import (
	"database/sql"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"

	"gorm.io/gorm"
)

type OrderGroupRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewOrderGroupRepo(db *db.DB) *OrderGroupRepo {
	return &OrderGroupRepo{
		db:     db,
		logger: logger.New("repo/order_group"),
	}
}

func (repo *OrderGroupRepo) CreateOrderGroup(req *models.CreateOrderGroupRequest) (*models.OrderGroup, error) {
	if req.GetRole().IsAdmin() {
		if req.UserID == "" {
			return nil, errs.ErrOrderGroupUserIDEmpty
		}
	} else {
		req.UserID = req.GetUserID()
	}
	var orderGroup = models.OrderGroup{
		Name:   req.Name,
		UserID: req.UserID,
	}
	if err := repo.db.Create(&orderGroup).Error; err != nil {
		return nil, err
	}
	return &orderGroup, nil
}
func (repo *OrderGroupRepo) GetOrderGroupList(params *models.GetOrderGroupListRequest) (*query.Pagination, error) {
	var builder = queryfunc.NewOrderGroupBuilder(queryfunc.OrderGroupBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Page == 0 {
		params.Page = 1
	}
	if params.GetRole() == enums.RoleClient {
		params.UserID = params.GetUserID()
	}
	sampleCompleteStatus := []enums.PoTrackingStatus{
		enums.PoTrackingStatusDeliveryConfirmed,
		enums.PoTrackingStatusDelivering,
	}
	bulkCompleteStatus := []enums.BulkPoTrackingStatus{
		enums.BulkPoTrackingStatusDeliveryConfirmed,
		enums.BulkPoTrackingStatusDelivering,
		enums.BulkPoTrackingStatusDelivered,
	}
	result := query.New(repo.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.UserID != "" {
				builder.Where("user_id = ?", params.UserID)
			}

			if params.OrderGroupStatus == enums.OrderGroupStatusComplete {
				builder.Where("NOT EXISTS (SELECT 1 FROM order_groups og LEFT JOIN inquiries iq ON og.id = iq.order_group_id WHERE og.id = o.id AND iq.status != ?)", enums.InquiryStatusFinished)
				builder.Where("NOT EXISTS (SELECT 1 FROM order_groups og LEFT JOIN purchase_orders po ON og.id = po.order_group_id WHERE og.id = o.id AND po.tracking_status NOT IN ?)", sampleCompleteStatus)
				builder.Where("NOT EXISTS (SELECT 1 FROM order_groups og LEFT JOIN bulk_purchase_orders bpo ON og.id = bpo.order_group_id WHERE og.id = o.id AND bpo.tracking_status NOT IN ?)", bulkCompleteStatus)
			} else if params.OrderGroupStatus == enums.OrderGroupStatusOnGoing {
				builder.Where(`(
					EXISTS (SELECT 1 FROM order_groups og LEFT JOIN inquiries iq ON og.id = iq.order_group_id WHERE og.id = o.id AND (iq.id is null OR iq.status != ?))
					OR EXISTS (SELECT 1 FROM order_groups og LEFT JOIN purchase_orders po ON og.id = po.order_group_id WHERE og.id = o.id AND (po.id is null OR po.tracking_status NOT IN ?))
					OR EXISTS (SELECT 1 FROM order_groups og LEFT JOIN bulk_purchase_orders bpo ON og.id = bpo.order_group_id WHERE og.id = o.id AND (bpo.id is null OR bpo.tracking_status NOT IN ?))
				)`, enums.InquiryStatusFinished, sampleCompleteStatus, bulkCompleteStatus)
			}

			if params.Keyword != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("id ILIKE @keyword", sql.Named("keyword", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result, nil
}
func (repo *OrderGroupRepo) GetOrderGroupDetail(params *models.GetOrderGroupDetailRequest) (*models.OrderGroup, error) {
	var builder = queryfunc.NewOrderGroupBuilder(queryfunc.OrderGroupBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		WithOrderDetail: true,
	})
	var userID string
	if params.GetRole() == enums.RoleClient {
		userID = params.GetUserID()
	}
	var orderGroup models.OrderGroup
	var err = query.New(repo.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("o.id = ?", params.OrderGroupID)
			if userID != "" {
				builder.Where("o.user_id = ?", userID)
			}
		}).
		Limit(1).
		FirstFunc(&orderGroup)

	return &orderGroup, err
}

func (repo *OrderGroupRepo) AssignOrderGroups(req *models.AssignOrderGroupRequest) error {
	var inquires models.Inquiries
	var purchaseOrders models.PurchaseOrders
	var bulkPurchaseOrders models.BulkPurchaseOrders
	var orderRetriever = NewOrderRetriever(repo.db)
	var err error
	switch req.OrderType {
	case enums.OrderGroupTypeRFQ:
		inquires, purchaseOrders, bulkPurchaseOrders, err = orderRetriever.RetrieveByRFQs(req.OrderIDs)

	case enums.OrderGroupTypeSample:
		inquires, purchaseOrders, bulkPurchaseOrders, err = orderRetriever.RetrieveBySamples(req.OrderIDs)

	case enums.OrderGroupTypeBulk:
		inquires, purchaseOrders, bulkPurchaseOrders, err = orderRetriever.RetrieveByBulks(req.OrderIDs)
	}
	if err != nil {
		return err
	}

	var orderGroup models.OrderGroup
	if err := repo.db.Select("ID").First(&orderGroup, "id = ? ", req.OrderGroupID).Error; err != nil {
		if repo.db.IsRecordNotFoundError(err) {
			return errs.ErrOrderGroupNotFound
		}
		return err
	}

	if err := repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Inquiry{}).Where("id IN ?", inquires.IDs()).UpdateColumn("order_group_id", req.OrderGroupID).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.PurchaseOrder{}).Where("id IN ?", purchaseOrders.IDs()).UpdateColumn("order_group_id", req.OrderGroupID).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.BulkPurchaseOrder{}).Where("id IN ?", bulkPurchaseOrders.IDs()).UpdateColumn("order_group_id", req.OrderGroupID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
