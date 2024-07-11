package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type ChatMessage struct {
	Model

	ReceiverID  string                `gorm:"not null" json:"receiver_id"`
	MessageType enums.ChatMessageType `gorm:"default:'user'" json:"message_type"`

	// MessageType = Group
	ReceiverAsGroup *ChatRoom `gorm:"-" json:"receiver_as_group,omitempty"`
	// MessageType = User
	ReceiverAsTeam *User `gorm:"-" json:"receiver_as_user,omitempty"`

	SenderID string `gorm:"not null" json:"sender_id"`
	Sender   *User  `gorm:"-" json:"sender,omitempty"`

	Message     string       `json:"message,omitempty"`
	Attachments *Attachments `json:"attachments,omitempty"`

	SeenAt *int64 `json:"seen_at,omitempty"`
}

type CreateChatMessageRequest struct {
	JwtClaimsInfo
	ReceiverID  string       `json:"receiver_id" validate:"required"`
	SenderID    string       `json:"sender_id"`
	Message     string       `json:"message"`
	Attachments *Attachments `json:"attachments,omitempty"`
}

type GetMessageListRequest struct {
	PaginationParams
	JwtClaimsInfo
	RoomID string `json:"room_id" query:"room_id"`
}

type GetChatUserRelevantStageRequest struct {
	JwtClaimsInfo
	InquiryID           string `json:"inquiry_id" query:"inquiry_id"`
	PurchaseOrderID     string `json:"purchase_order_id" query:"purchase_order_id"`
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id"`
}

type WSMessagePayload struct {
	Type     string       `json:"type"`
	Message  *ChatMessage `json:"message"`
	ChatRoom *ChatRoom    `json:"chat_room"`
}
