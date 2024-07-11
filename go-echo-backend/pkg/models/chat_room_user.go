package models

type ChatRoomUser struct {
	ID        string `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	RoomID string `gorm:"uniqueIndex:idx_chat_room_user" json:"room_id"`
	UserID string `gorm:"uniqueIndex:idx_chat_room_user" json:"user_id"`
}

type MarkSeenChatRoomMessageRequest struct {
	JwtClaimsInfo
	RoomID string `json:"chat_room_id"  param:"chat_room_id" validate:"required"`
}
