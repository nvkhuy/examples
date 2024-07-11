package queryfunc

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type ChatRoomBuilderOptions struct {
	QueryBuilderOptions
	ReferenceIDKeyWord string
}

func NewChatRoomBuilder(options ChatRoomBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* ChatRoomBuilder {{.ChatRoom}} */ cc.* 
	FROM (
		SELECT cr.*,

		(SELECT COALESCE(NULLIF(cr.bulk_purchase_order_id,''),
		NULLIF(cr.purchase_order_id ,''),
		NULLIF(cr.inquiry_id ,''))
		) as stage_id, 
	
		json_agg(json_build_object(
			'id',u.id,
			'avatar',u.avatar,
			'name',u.name,
			'role',u.role,
			'last_online_at',u.last_online_at,
			'is_offline',u.is_offline,
			'contact_owner_ids',u.contact_owner_ids
		)) AS users_json,

		(
			SELECT row_to_json(c) FROM (
				SELECT cm.*
				FROM chat_messages cm
				WHERE cm.receiver_id = cr.id
				ORDER BY cm.id DESC
				LIMIT 1
			) c
		) AS latest_message_json

		FROM chat_rooms cr
		JOIN chat_room_users cru ON cru.room_id = cr.id 
		JOIN users u ON u.id = cru.user_id 
		GROUP BY cr.id 
	) cc	
	`

	if options.Role == enums.RoleClient {
		rawSQL = `
	SELECT /* ChatRoomBuilder {{.ChatRoom}} */ cc.* 
	FROM (
		SELECT cr.*, 

		(SELECT COALESCE(NULLIF(cr.bulk_purchase_order_id,''),
		NULLIF(cr.purchase_order_id ,''),
		NULLIF(cr.inquiry_id ,''))
		) as stage_id,

		json_agg(json_build_object(
			'id',u.id,
			'avatar',u.avatar,
			'name',u.name,
			'role',u.role,
			'last_online_at',u.last_online_at,
			'is_offline',u.is_offline
		)) AS users_json,

		(
			SELECT row_to_json(c) FROM (
				SELECT cm.*
				FROM chat_messages cm
				WHERE cm.receiver_id = cr.id
				ORDER BY created_at DESC
				LIMIT 1
			) c
		) AS latest_message_json

		FROM chat_rooms cr
		JOIN chat_room_users cru ON cru.room_id = cr.id 
		JOIN users u ON u.id = cru.user_id
		WHERE EXISTS ( SELECT 1 FROM chat_room_users cru2 WHERE room_id = cr.id AND user_id = @userID )
		GROUP BY cr.id 
	) cc	
	`
	}
	if options.ReferenceIDKeyWord != "" {
		unionStmt := fmt.Sprintf(`WITH union_data AS (
		SELECT u.id FROM (
			SELECT id,reference_id FROM inquiries UNION
			SELECT id,reference_id FROM purchase_orders po WHERE po.status = 'paid' OR from_catalog = true UNION
			SELECT id,reference_id FROM bulk_purchase_orders
		) AS u
		WHERE u.reference_id ILIKE '%%%s%%'
	)`, options.ReferenceIDKeyWord)
		if options.Role == enums.RoleClient {
			unionStmt = fmt.Sprintf(`WITH union_data AS (
			SELECT u.id FROM (
				SELECT id,reference_id FROM inquiries WHERE user_id = @userID UNION
				SELECT id,reference_id FROM purchase_orders po WHERE (po.status = 'paid' OR from_catalog = true) AND user_id = @userID UNION
				SELECT id,reference_id FROM bulk_purchase_orders WHERE user_id = @userID
			) AS u
			WHERE u.reference_id ILIKE '%%%s%%'
		)`, options.ReferenceIDKeyWord)
		}

		rawSQL = fmt.Sprintf("%s %s", unionStmt, rawSQL)
	}
	return NewBuilder(rawSQL).
		WithOptions(options).
		WithOrderBy("cc.latest_message_json ->> 'created_at' DESC NULLS LAST").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.ChatRoom, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.ErrorAny(err)
				return nil, err

			}
			defer rows.Close()

			var purchaseOrderIds []string
			var bulkPurchaseOrderIds []string
			var inquiryIds []string

			for rows.Next() {
				var alias models.ChatRoomAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				if alias.LatestMessageJson != nil {
					if err := json.Unmarshal(alias.LatestMessageJson, &alias.ChatRoom.LatestMessage); err != nil {
						continue
					}
				}
				if alias.UsersJson != nil {
					if err := json.Unmarshal(alias.UsersJson, &alias.ChatRoom.RoomUsers); err != nil {
						continue
					}
				}

				if alias.BulkPurchaseOrderID != "" {
					bulkPurchaseOrderIds = append(bulkPurchaseOrderIds, alias.BulkPurchaseOrderID)
				} else if alias.PurchaseOrderID != "" {
					purchaseOrderIds = append(purchaseOrderIds, alias.PurchaseOrderID)
				} else {
					inquiryIds = append(inquiryIds, alias.InquiryID)
				}
				records = append(records, &alias.ChatRoom)
			}
			var wg sync.WaitGroup
			wg.Add(3)
			go func() {
				defer wg.Done()
				if len(bulkPurchaseOrderIds) > 0 {
					var bulkPurchaseOrders []*models.BulkPurchaseOrder
					if err := db.Select("ID", "ReferenceID").Find(&bulkPurchaseOrders, "id IN ?", bulkPurchaseOrderIds).Error; err != nil {
						return
					}
					for _, chatRoom := range records {
						for _, bulk := range bulkPurchaseOrders {
							if chatRoom.BulkPurchaseOrderID == bulk.ID {
								chatRoom.ChatRoomStatus = &models.ChatRoomStatus{
									Stage:       enums.ChatRoomStageBulk,
									StageID:     bulk.ID,
									ReferenceID: bulk.ReferenceID,
								}
							}
						}
					}
				}
			}()
			go func() {
				defer wg.Done()
				if len(purchaseOrderIds) > 0 {
					var purchaseOrders []*models.PurchaseOrder
					if err := db.Select("ID", "ReferenceID", "Status", "InquiryID").Find(&purchaseOrders, "id IN ?", purchaseOrderIds).Error; err != nil {
						return
					}
					for _, chatRoom := range records {
						for _, po := range purchaseOrders {
							if chatRoom.PurchaseOrderID == po.ID && chatRoom.BulkPurchaseOrderID == "" {
								chatRoom.ChatRoomStatus = &models.ChatRoomStatus{
									Stage:       enums.ChatRoomStageSample,
									StageID:     po.ID,
									ReferenceID: po.ReferenceID,
								}
							}
						}
					}
				}
			}()
			go func() {
				defer wg.Done()
				if len(inquiryIds) > 0 {
					var inquiries []*models.Inquiry
					if err := db.Select("ID", "ReferenceID").Find(&inquiries, "id IN ?", inquiryIds).Error; err != nil {
						return
					}
					for _, chatRoom := range records {
						for _, iq := range inquiries {
							if chatRoom.InquiryID == iq.ID && chatRoom.PurchaseOrderID == "" && chatRoom.BulkPurchaseOrderID == "" {
								chatRoom.ChatRoomStatus = &models.ChatRoomStatus{
									Stage:       enums.ChatRoomStageRFQ,
									StageID:     iq.ID,
									ReferenceID: iq.ReferenceID,
								}
							}
						}
					}
				}
			}()
			wg.Wait()

			return &records, nil
		})
}
