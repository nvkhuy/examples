package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"gorm.io/gorm"
)

func MigrateChatRooms(role string) {
	var cfg = config.New("../deployment/config/local/env.json")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var buyerSql = `
		update chat_rooms cr set buyer_id = (select cru.user_id from chat_room_users cru join users u on cru.user_id = u.id where cru.room_id = cr.id and u.role = 'client')
		where cr.participant_role = 'admin-buyer' and cr.buyer_id = ''
	`
	var sellerSql = `
		update chat_rooms cr set seller_id = (select cru.user_id from chat_room_users cru join users u on cru.user_id = u.id where cru.room_id = cr.id and u.role = 'seller')
		where cr.participant_role = 'admin-seller' and cr.seller_id = ''
	`
	var sql = buyerSql
	if role == "seller" {
		sql = sellerSql
	}
	if err := app.DB.Exec(sql).Error; err != nil {
		fmt.Printf("Migrate chat room error, err=%+v\n", err)
		return
	}
	fmt.Println("Migrate chat room successfully")
}

func MigrateChatRoomsUpdateIndex() {
	var cfg = config.New("../deployment/config/local/env.json")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var dropIdxSql = `DROP INDEX IF EXISTS idx_chat_room`
	var createIdxSql = `CREATE UNIQUE INDEX IF NOT EXISTS idx_chat_room ON chat_rooms (inquiry_id,purchase_order_id,bulk_purchase_order_id,buyer_id,seller_id)`

	if err := app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(dropIdxSql).Error; err != nil {
			return err
		}
		if err := tx.Exec(createIdxSql).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		fmt.Printf("Create chat room index error, err=%+v\n", err)
		return
	}
	fmt.Println("Update chat room index successfully")
}
