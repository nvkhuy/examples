package repo

import (
	"errors"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"

	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type ChatRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewChatRepo(db *db.DB) *ChatRepo {
	return &ChatRepo{
		db:     db,
		logger: logger.New("repo/chat"),
	}
}

func (r *ChatRepo) CreateChatMessage(req *models.CreateChatMessageRequest) (*models.ChatMessage, error) {
	if req.Message == "" && req.Attachments == nil {
		return nil, errors.New("request params are empty")
	}
	if err := r.db.Select("ID").First(&models.ChatRoom{}, "id = ?", req.ReceiverID).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrChatRoomNotFound
		}
		return nil, err
	}
	if err := r.db.Select("ID").First(&models.User{}, "id = ?", req.SenderID).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	var message models.ChatMessage
	if err := copier.Copy(&message, req); err != nil {
		return nil, eris.Wrap(err, "copy attribute error")
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&message).Error; err != nil {
			return err
		}
		err := tx.Select("ID").First(&models.ChatRoomUser{}, "room_id = ? and user_id = ?", req.ReceiverID, req.SenderID).Error
		if err != nil {
			if !r.db.IsRecordNotFoundError(err) {
				return err
			}
			return tx.Create(&models.ChatRoomUser{RoomID: req.ReceiverID, UserID: req.SenderID}).Error
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *ChatRepo) GetMessageList(params *models.GetMessageListRequest) (*query.Pagination, error) {
	chatRoom := &models.ChatRoom{}
	err := r.db.Select("ID").First(chatRoom, "id = ?", params.RoomID).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrChatRoomNotFound
		}
		return nil, err
	}
	var builder = queryfunc.NewChatMessageBuilder(queryfunc.ChatMessageBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}
	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("c.receiver_id = ? ", chatRoom.ID)
		}).
		Page(params.Page).
		Limit(params.Limit).
		WithoutCount(true).
		PagingFunc()

	return result, nil
}

func (r *ChatRepo) GetChatRelevantStage(params *models.GetChatUserRelevantStageRequest) ([]*models.ChatRoomStatus, error) {
	if err := validateGetChatRelevantStage(params); err != nil {
		return nil, err
	}
	var result []*models.ChatRoomStatus
	switch {
	case params.InquiryID != "":
		var inquiry models.Inquiry
		if err := r.db.Select("ID", "ReferenceID").First(&inquiry, "id = ?", params.InquiryID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrInquiryNotFound
			}
			return nil, err
		}
		result = append(result, &models.ChatRoomStatus{
			Stage:       enums.ChatRoomStageRFQ,
			StageID:     inquiry.ID,
			ReferenceID: inquiry.ReferenceID,
		})
		if params.GetRole().IsSeller() {
			var purchaseOrder models.PurchaseOrder
			if err := r.db.Select("ID", "ReferenceID", "SampleMakerID").First(&purchaseOrder, "inquiry_id = ?", inquiry.ID).Error; err != nil {
				if r.db.IsRecordNotFoundError(err) {
					return nil, errs.ErrPONotFound
				}
				return nil, err
			}
			if purchaseOrder.SampleMakerID == params.GetUserID() {
				result = append(result, &models.ChatRoomStatus{
					Stage:       enums.ChatRoomStageSample,
					StageID:     purchaseOrder.ID,
					ReferenceID: purchaseOrder.ReferenceID,
				})
			}
			var bulkPurchaseOrder models.BulkPurchaseOrder
			if err := r.db.Select("ID", "ReferenceID", "SellerID").
				First(&bulkPurchaseOrder, "inquiry_id = ? AND purchase_order_id = ?", inquiry.ID, purchaseOrder.ID).Error; err != nil {
				if !r.db.IsRecordNotFoundError(err) {
					return nil, err
				}
			}
			if bulkPurchaseOrder.ID != "" && bulkPurchaseOrder.SellerID == params.GetUserID() {
				result = append(result, &models.ChatRoomStatus{
					Stage:       enums.ChatRoomStageBulk,
					StageID:     bulkPurchaseOrder.ID,
					ReferenceID: bulkPurchaseOrder.ReferenceID,
				})
			}
		} else {
			var purchaseOrder models.PurchaseOrder
			if err := r.db.Select("ID", "Status", "ReferenceID").First(&purchaseOrder, "inquiry_id = ?", inquiry.ID).Error; err != nil {
				if r.db.IsRecordNotFoundError(err) {
					return nil, errs.ErrPONotFound
				}
				return nil, err
			}
			if purchaseOrder.Status == "paid" {
				result = append(result, &models.ChatRoomStatus{
					Stage:       enums.ChatRoomStageSample,
					StageID:     purchaseOrder.ID,
					ReferenceID: purchaseOrder.ReferenceID,
				})
				var bulkPurchaseOrder models.BulkPurchaseOrder
				if err := r.db.Select("ID", "ReferenceID").
					First(&bulkPurchaseOrder, "inquiry_id = ? AND purchase_order_id = ?", inquiry.ID, purchaseOrder.ID).Error; err != nil {
					if !r.db.IsRecordNotFoundError(err) {
						return nil, err
					}
				}
				if bulkPurchaseOrder.ID != "" {
					result = append(result, &models.ChatRoomStatus{
						Stage:       enums.ChatRoomStageBulk,
						StageID:     bulkPurchaseOrder.ID,
						ReferenceID: bulkPurchaseOrder.ReferenceID,
					})
				}
			}
		}
	case params.PurchaseOrderID != "":
		var purchaseOrder models.PurchaseOrder
		if err := r.db.Select("ID", "InquiryID", "ReferenceID").First(&purchaseOrder, "id = ?", params.PurchaseOrderID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrPONotFound
			}
			return nil, err
		}
		result = append(result, &models.ChatRoomStatus{
			Stage:       enums.ChatRoomStageSample,
			StageID:     purchaseOrder.ID,
			ReferenceID: purchaseOrder.ReferenceID,
		})

		if params.GetRole().IsSeller() {
			if purchaseOrder.InquiryID != "" {
				var inquiry models.Inquiry
				if err := r.db.Select("ID", "ReferenceID").First(&inquiry, "id = ?", purchaseOrder.InquiryID).Error; err != nil {
					if r.db.IsRecordNotFoundError(err) {
						return nil, errs.ErrInquiryNotFound
					}
					return nil, err
				}
				var iqSeller models.InquirySeller
				if err := r.db.Select("ID", "UserID").First(&iqSeller, "inquiry_id = ?", purchaseOrder.InquiryID).Error; err != nil {
					if !r.db.IsRecordNotFoundError(err) {
						return nil, err
					}
				}
				if iqSeller.ID != "" && iqSeller.UserID == params.GetUserID() {
					result = append(result, &models.ChatRoomStatus{
						Stage:       enums.ChatRoomStageRFQ,
						StageID:     inquiry.ID,
						ReferenceID: inquiry.ReferenceID,
					})
				}
			}

			var bulkPurchaseOrder models.BulkPurchaseOrder
			if err := r.db.Select("ID", "ReferenceID", "SellerID").
				First(&bulkPurchaseOrder, "purchase_order_id = ?", purchaseOrder.ID).Error; err != nil {
				if !r.db.IsRecordNotFoundError(err) {
					return nil, errs.ErrPONotFound
				}
				return result, nil
			}
			if bulkPurchaseOrder.ID != "" && bulkPurchaseOrder.SellerID == params.GetUserID() {
				result = append(result, &models.ChatRoomStatus{
					Stage:       enums.ChatRoomStageBulk,
					StageID:     bulkPurchaseOrder.ID,
					ReferenceID: bulkPurchaseOrder.ReferenceID,
				})
			}
		} else {
			if purchaseOrder.InquiryID != "" {
				var inquiry models.Inquiry
				if err := r.db.Select("ID", "ReferenceID").First(&inquiry, "id = ?", purchaseOrder.InquiryID).Error; err != nil {
					if r.db.IsRecordNotFoundError(err) {
						return nil, errs.ErrInquiryNotFound
					}
					return nil, err
				}
				result = append(result, &models.ChatRoomStatus{
					Stage:       enums.ChatRoomStageRFQ,
					StageID:     inquiry.ID,
					ReferenceID: inquiry.ReferenceID,
				})
			}
			var bulkPurchaseOrder models.BulkPurchaseOrder
			if err := r.db.Select("ID", "ReferenceID", "SellerID").
				First(&bulkPurchaseOrder, "purchase_order_id = ?", purchaseOrder.ID).Error; err != nil {
				if !r.db.IsRecordNotFoundError(err) {
					return nil, errs.ErrPONotFound
				}
			}
			if bulkPurchaseOrder.ID != "" {
				result = append(result, &models.ChatRoomStatus{
					Stage:       enums.ChatRoomStageBulk,
					StageID:     bulkPurchaseOrder.ID,
					ReferenceID: bulkPurchaseOrder.ReferenceID,
				})
			}
		}
	case params.BulkPurchaseOrderID != "":
		var bulkPurchaseOrder models.BulkPurchaseOrder
		if err := r.db.Select("ID", "ReferenceID", "InquiryID", "PurchaseOrderID").First(&bulkPurchaseOrder, "id = ?", params.BulkPurchaseOrderID).Error; err != nil {
			if r.db.IsRecordNotFoundError(err) {
				return nil, errs.ErrBulkPoNotFound
			}
			return nil, err
		}
		result = append(result, &models.ChatRoomStatus{
			Stage:       enums.ChatRoomStageBulk,
			StageID:     bulkPurchaseOrder.ID,
			ReferenceID: bulkPurchaseOrder.ReferenceID,
		})

		if params.GetRole().IsSeller() {
			if bulkPurchaseOrder.InquiryID != "" {
				var inquiry models.Inquiry
				if err := r.db.Select("ID", "ReferenceID").First(&inquiry, "id = ?", bulkPurchaseOrder.InquiryID).Error; err != nil {
					if r.db.IsRecordNotFoundError(err) {
						return nil, errs.ErrInquiryNotFound
					}
					return nil, err
				}
				var iqSeller models.InquirySeller
				if err := r.db.Select("ID", "UserID").First(&iqSeller, "inquiry_id = ?", bulkPurchaseOrder.InquiryID).Error; err != nil {
					if !r.db.IsRecordNotFoundError(err) {
						return nil, err
					}
				}
				if iqSeller.ID != "" && iqSeller.UserID == params.GetUserID() {
					result = append(result, &models.ChatRoomStatus{
						Stage:       enums.ChatRoomStageRFQ,
						StageID:     inquiry.ID,
						ReferenceID: inquiry.ReferenceID,
					})
				}
			}
			if bulkPurchaseOrder.PurchaseOrderID != "" {
				var purchaseOrder models.PurchaseOrder
				if err := r.db.Select("ID", "ReferenceID", "SampleMakerID").First(&purchaseOrder, "id = ?", bulkPurchaseOrder.PurchaseOrderID).Error; err != nil {
					if r.db.IsRecordNotFoundError(err) {
						return nil, errs.ErrPONotFound
					}
					return nil, err
				}
				if purchaseOrder.SampleMakerID == params.GetUserID() {
					result = append(result, &models.ChatRoomStatus{
						Stage:       enums.ChatRoomStageSample,
						StageID:     purchaseOrder.ID,
						ReferenceID: purchaseOrder.ReferenceID,
					})
				}
			}
		} else {
			if bulkPurchaseOrder.InquiryID != "" {
				var inquiry models.Inquiry
				if err := r.db.Select("ID", "ReferenceID").First(&inquiry, "id = ?", bulkPurchaseOrder.InquiryID).Error; err != nil {
					if r.db.IsRecordNotFoundError(err) {
						return nil, errs.ErrInquiryNotFound
					}
					return nil, err
				}
				result = append(result, &models.ChatRoomStatus{
					Stage:       enums.ChatRoomStageRFQ,
					StageID:     inquiry.ID,
					ReferenceID: inquiry.ReferenceID,
				})
			}
			if bulkPurchaseOrder.PurchaseOrderID != "" {
				var purchaseOrder models.PurchaseOrder
				if err := r.db.Select("ID", "ReferenceID", "SampleMakerID").First(&purchaseOrder, "id = ?", bulkPurchaseOrder.PurchaseOrderID).Error; err != nil {
					if r.db.IsRecordNotFoundError(err) {
						return nil, errs.ErrPONotFound
					}
				}
				result = append(result, &models.ChatRoomStatus{
					Stage:       enums.ChatRoomStageSample,
					StageID:     purchaseOrder.ID,
					ReferenceID: purchaseOrder.ReferenceID,
				})

			}
		}
	}
	return result, nil
}

func validateGetChatRelevantStage(req *models.GetChatUserRelevantStageRequest) error {
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
		return errors.New("request params are empty")
	} else if count > 1 {
		return errors.New("too many id in request params")
	}

	return nil
}
