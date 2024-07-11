package repo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"

	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChatRoomRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewChatRoomRepo(db *db.DB) *ChatRoomRepo {
	return &ChatRoomRepo{
		db:     db,
		logger: logger.New("repo/chat_room"),
	}
}
func (r *ChatRoomRepo) GetChatRoomList(params *models.GetChatRoomListRequest) (*query.Pagination, error) {
	var builder = queryfunc.NewChatRoomBuilder(queryfunc.ChatRoomBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		ReferenceIDKeyWord: params.Keyword,
	})
	if params.Limit == 0 {
		params.Limit = 12
	}
	if params.Page == 0 {
		params.Page = 1
	}
	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.GetRole() == enums.RoleClient {
				builder.Where(map[string]string{
					"userID": params.GetUserID(),
				})
			}
			if params.Role != "" {
				builder.Where("EXISTS (SELECT 1 FROM chat_room_users cru JOIN users u on cru.user_id = u.id where room_id = cc.id and u.role = ?)", params.Role)
			}
			if params.Keyword != "" {
				builder.Where("cc.stage_id IN (SELECT id FROM union_data)")
			} else {
				builder.Where("cc.latest_message_json IS NOT NULL")
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		WithoutCount(true).
		PagingFunc()

	return result, nil
}
func (r *ChatRoomRepo) GetChatRoomInfo(roomID string) (*models.ChatRoom, error) {
	var builder = queryfunc.NewChatRoomBuilder(queryfunc.ChatRoomBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{},
	})
	var chatRoom models.ChatRoom
	err := query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cc.id = ?", roomID)
		}).
		WithoutCount(true).
		FirstFunc(&chatRoom)
	if err != nil {
		return nil, err
	}
	return &chatRoom, nil
}

func (r *ChatRoomRepo) CreateChatRoom(req *models.CreateChatRoomRequest) (room *models.ChatRoom, err error) {
	if err := validateCreateChatRoom(req); err != nil {
		return nil, err
	}
	var buyer, seller models.User
	if req.BuyerID != "" {
		if err := r.db.Select("ID", "ContactOwnerIDs").First(&buyer, "id = ?", req.BuyerID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrUserNotFound
			}
			return nil, err
		}
	}
	if req.SellerID != "" {
		if err := r.db.Select("ID").First(&seller, "id = ?", req.SellerID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrUserNotFound
			}
			return nil, err
		}
	}
	whereMap := map[string]interface{}{
		"inquiry_id":             "",
		"purchase_order_id":      "",
		"bulk_purchase_order_id": "",
		"buyer_id":               req.BuyerID,
		"seller_id":              req.SellerID,
	}

	var chatRoomUsers []*models.ChatRoomUser
	switch {
	case req.InquiryID != "":
		inquiry := &models.Inquiry{}
		if err := r.db.Select("ID", "AssigneeIDs", "UserID").First(inquiry, "id = ?", req.InquiryID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrInquiryNotFound
			}
			return nil, err
		}
		whereMap["inquiry_id"] = req.InquiryID
		for _, assignee := range inquiry.AssigneeIDs {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: assignee})
		}
		if req.BuyerID != "" {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: inquiry.UserID})
		} else {
			var iqSeller models.InquirySeller
			if err := r.db.Select("UserID").First(&iqSeller, "inquiry_id = ?", inquiry.ID).Error; err != nil {
				if r.db.IsRecordNotFoundError(err) {
					return nil, errs.ErrInquirySellerNotFound
				}
				return nil, err
			}
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: iqSeller.UserID})
		}

	case req.PurchaseOrderID != "":
		purchaseOrder := &models.PurchaseOrder{}
		if err := r.db.Select("ID", "inquiry_id", "AssigneeIDs", "UserID").First(purchaseOrder, "id = ?", req.PurchaseOrderID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrPONotFound
			}
			return nil, err
		}

		if purchaseOrder.InquiryID != "" {
			req.InquiryID = purchaseOrder.InquiryID
		}
		whereMap["purchase_order_id"] = req.PurchaseOrderID
		whereMap["inquiry_id"] = req.InquiryID
		for _, assignee := range purchaseOrder.AssigneeIDs {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: assignee})
		}
		if req.BuyerID != "" {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: purchaseOrder.UserID})
		} else {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: purchaseOrder.SampleMakerID})
		}
	case req.BulkPurchaseOrderID != "":
		bulkPurchaseOrder := &models.BulkPurchaseOrder{}
		if err := r.db.Select("ID", "inquiry_id", "purchase_order_id", "AssigneeIDs", "UserID").First(bulkPurchaseOrder, "id = ?", req.BulkPurchaseOrderID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrBulkPoNotFound
			}
			return nil, err
		}

		if bulkPurchaseOrder.InquiryID != "" {
			req.InquiryID = bulkPurchaseOrder.InquiryID
		}
		if bulkPurchaseOrder.PurchaseOrderID != "" {
			req.PurchaseOrderID = bulkPurchaseOrder.PurchaseOrderID
		}
		whereMap["bulk_purchase_order_id"] = req.BulkPurchaseOrderID
		whereMap["inquiry_id"] = req.InquiryID
		whereMap["purchase_order_id"] = req.PurchaseOrderID
		for _, assignee := range bulkPurchaseOrder.AssigneeIDs {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: assignee})
		}
		if req.BuyerID != "" {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: bulkPurchaseOrder.UserID})
		} else {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: bulkPurchaseOrder.SellerID})
		}
	}
	var chatRoom models.ChatRoom
	err = r.db.Where(whereMap).First(&chatRoom).Error
	if err == nil {
		return &chatRoom, nil
	}
	if !r.db.IsRecordNotFoundError(err) {
		return nil, err
	}

	if err := copier.Copy(&chatRoom, req); err != nil {
		return nil, eris.Wrap(err, "copy attribute error")
	}
	chatRoom.HostID = req.GetUserID()

	if req.BuyerID != "" {
		for _, userID := range buyer.ContactOwnerIDs {
			_, _, ok := lo.FindIndexOf(chatRoomUsers, func(item *models.ChatRoomUser) bool {
				return item.UserID == userID
			})
			if !ok {
				chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{UserID: userID})
			}
		}
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&chatRoom).Error; err != nil {
			return err
		}
		for _, chatRoomUser := range chatRoomUsers {
			chatRoomUser.RoomID = chatRoom.ID
		}
		if err := tx.Create(&chatRoomUsers).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &chatRoom, nil
}

func (r *ChatRoomRepo) MarkSeenChatRoomMessage(params *models.MarkSeenChatRoomMessageRequest) error {
	var err = r.db.Model(&models.ChatMessage{}).Clauses(clause.Returning{}).
		Where("seen_at IS NULL").
		Where("receiver_id = ?", params.RoomID).
		Where("sender_id != ?", params.GetUserID()).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

func (r *ChatRoomRepo) CountUnSeenChatMessage(params *models.CountUnSeenChatMessageRequest) (int, error) {
	countSql := `SELECT COUNT(1) FROM (
		select 1
		FROM chat_rooms cr
	 	JOIN chat_messages cm on cm.receiver_id  = cr.id
		JOIN users u on u.id = cm.sender_id
		WHERE cm.seen_at is null AND (u.role = ? or u.role = ?)
		GROUP BY cr.id
		LIMIT 100
	) as c`
	if params.GetRole() == enums.RoleClient {
		countSql = `SELECT COUNT(1) FROM (
			select 1
			FROM chat_rooms cr
			JOIN chat_messages cm on cm.receiver_id  = cr.id
			WHERE cm.seen_at is null AND cm.sender_id != @userID
			AND exists ( select 1 from chat_room_users cru2 where room_id = cr.id and user_id = @userID )
			GROUP BY cr.id
			LIMIT 100
		) as c`
	}

	var count int
	var arg []interface{} = []interface{}{enums.RoleClient, enums.RoleSeller}
	if params.GetRole().IsBuyer() || params.GetRole().IsSeller() {
		arg = []interface{}{sql.Named("userID", params.GetUserID())}
	}
	if err := r.db.Raw(countSql, arg...).Scan(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
func (r *ChatRoomRepo) CountUnSeenChatMessageOnRoom(params *models.CountUnSeenChatMessageOnRoomRequest) (*models.CountUnSeenChatMessageOnRoomResponse, error) {
	whereMap := map[string]interface{}{
		"inquiry_id":             "",
		"purchase_order_id":      "",
		"bulk_purchase_order_id": "",
		"buyer_id":               params.BuyerID,
		"seller_id":              params.SellerID,
	}
	switch {
	case params.InquiryID != "":
		inquiry := &models.Inquiry{}
		if err := r.db.Select("ID").First(inquiry, "id = ?", params.InquiryID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrInquiryNotFound
			}
			return nil, err
		}
		whereMap["inquiry_id"] = params.InquiryID
	case params.PurchaseOrderID != "":
		purchaseOrder := &models.PurchaseOrder{}
		if err := r.db.Select("ID", "inquiry_id").First(purchaseOrder, "id = ?", params.PurchaseOrderID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrPONotFound
			}
			return nil, err
		}

		whereMap["purchase_order_id"] = params.PurchaseOrderID
		whereMap["inquiry_id"] = purchaseOrder.InquiryID

	case params.BulkPurchaseOrderID != "":
		bulkPurchaseOrder := &models.BulkPurchaseOrder{}
		if err := r.db.Select("ID", "inquiry_id", "purchase_order_id").First(bulkPurchaseOrder, "id = ?", params.BulkPurchaseOrderID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrBulkPoNotFound
			}
			return nil, err
		}

		whereMap["bulk_purchase_order_id"] = params.BulkPurchaseOrderID
		whereMap["inquiry_id"] = bulkPurchaseOrder.InquiryID
		whereMap["purchase_order_id"] = bulkPurchaseOrder.PurchaseOrderID
	}
	var chatRoom models.ChatRoom
	if err := r.db.Select("ID").Where(whereMap).First(&chatRoom).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return &models.CountUnSeenChatMessageOnRoomResponse{}, nil
		}
		return nil, err
	}

	var count int64
	if err := r.db.Model(&models.ChatMessage{}).
		Where("receiver_id = ? AND seen_at IS NULL AND sender_id != ?", chatRoom.ID, params.GetUserID()).
		Count(&count).Error; err != nil {
		return nil, err
	}

	var chatMessage models.ChatMessage
	if err := r.db.First(&chatMessage, "receiver_id = ?", chatRoom.ID).Error; err != nil {
		if !r.db.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	hasHistory := false
	if chatMessage.ID != "" {
		hasHistory = true
	}
	return &models.CountUnSeenChatMessageOnRoomResponse{RoomID: chatRoom.ID, Count: int(count), HasChatHistory: hasHistory}, nil
}

func (r *ChatRoomRepo) GetUnseenChatRoomsOver6Hours() ([]*models.ChatRoom, error) {
	var unseenRooms []*models.ChatRoom
	if err := query.New(r.db, queryfunc.NewChatRoomBuilder(queryfunc.ChatRoomBuilderOptions{})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where(`
				EXISTS (SELECT 1 FROM chat_messages cm WHERE cm.receiver_id = cc.id AND cm.seen_at IS NULL
				AND EXTRACT(EPOCH FROM NOW()) - created_at > 3600*6 AND EXTRACT(EPOCH FROM NOW()) - created_at < 3600*12)`)
		}).WithoutCount(true).FindFunc(&unseenRooms); err != nil {
		return nil, err
	}
	return unseenRooms, nil
}

func validateCreateChatRoom(req *models.CreateChatRoomRequest) error {
	count := 0
	if req.InquiryID != "" {
		count++
	}
	if req.PurchaseOrderID != "" {
		count++
	}
	if req.BulkPurchaseOrderID != "" {
		count++
	}
	if count == 0 {
		return errors.New("order params are empty")
	} else if count > 1 {
		return errors.New("too many order id in request params")
	}
	if req.BuyerID == "" && req.SellerID == "" {
		return errors.New("seller_id and buyer_id are empty")
	}
	if req.BuyerID != "" && req.SellerID != "" {
		return errors.New("seller_id and buyer_id can not have both value")
	}

	return nil
}
