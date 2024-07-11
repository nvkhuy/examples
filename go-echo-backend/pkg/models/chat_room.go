package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/datatypes"
)

type ChatRoom struct {
	Model
	HostID              string `json:"host_id"`
	InquiryID           string `gorm:"uniqueIndex:idx_chat_room;not null;default:''" json:"inquiry_id"`
	PurchaseOrderID     string `gorm:"uniqueIndex:idx_chat_room;not null;default:''" json:"purchase_order_id"`
	BulkPurchaseOrderID string `gorm:"uniqueIndex:idx_chat_room;not null;default:''" json:"bulk_purchase_order_id"`
	BuyerID             string `gorm:"uniqueIndex:idx_chat_room;not null;default:''" json:"buyer_id"`
	SellerID            string `gorm:"uniqueIndex:idx_chat_room;not null;default:''" json:"seller_id"`
	// Status              enums.ChatRoomStatus `json:"status"`

	ChatRoomStatus *ChatRoomStatus `gorm:"-" json:"chat_room_status,omitempty"`
	LatestMessage  *ChatMessage    `gorm:"-" json:"latest_message,omitempty"`
	RoomUsers      []*User         `gorm:"-" json:"room_users,omitempty"`
}
type ChatRoomStatus struct {
	Stage       enums.ChatRoomStage `json:"stage"`
	StageID     string              `json:"stage_id"`
	ReferenceID string              `json:"reference_id"`
}

type ChatRoomAlias struct {
	ChatRoom

	LatestMessageJson datatypes.JSON `gorm:"column:latest_message_json"`
	UsersJson         datatypes.JSON `gorm:"column:users_json"`
}
type CreateChatRoomRequest struct {
	JwtClaimsInfo
	InquiryID           string `json:"inquiry_id"`
	PurchaseOrderID     string `json:"purchase_order_id"`
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id"`
	BuyerID             string `json:"buyer_id"`
	SellerID            string `json:"seller_id"`
}

type GetChatRoomListRequest struct {
	JwtClaimsInfo
	PaginationParams
	Role enums.Role `json:"role" validate:"required,oneof=client seller" query:"role"`
}
type CountUnSeenChatMessageRequest struct {
	JwtClaimsInfo
	PaginationParams
}

type CountUnSeenChatMessageResponse struct {
	Count int `json:"count"`
}

type CountUnSeenChatMessageOnRoomRequest struct {
	JwtClaimsInfo
	InquiryID           string `json:"inquiry_id" query:"inquiry_id"`
	PurchaseOrderID     string `json:"purchase_order_id" query:"purchase_order_id"`
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id"`
	BuyerID             string `json:"buyer_id" query:"buyer_id"`
	SellerID            string `json:"seller_id" query:"seller_id"`
}

type CountUnSeenChatMessageOnRoomResponse struct {
	RoomID         string `json:"room_id"`
	Count          int    `json:"count"`
	HasChatHistory bool   `json:"has_chat_history"`
}
