package ws

import (
	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

// EventType event
type EventType string

// All types
var (
	EventTypeChatSendMessage EventType = "chat_send_message"
)

type MessageType string

var (
	MessageTypePing         MessageType = "ping"
	MessageTypeTyping       MessageType = "typing"
	MessageTypeCancelTyping MessageType = "cancel_typing"
	MessageTypeSeenMessage  MessageType = "seen_message"
	MessageTypeChat         MessageType = "chat"
)

type Message struct {
	Type   MessageType  `json:"type"`
	UserID string       `json:"user_id"`
	Role   enums.Role   `json:"role"`
	Data   *MessageData `json:"data"`
}

type MessageData struct {
	ChatRoomID string `json:"chat_room_id,omitempty"`
}

type BroadcastChatMessage struct {
	Type                              string              `json:"type"`
	Message                           *models.ChatMessage `json:"message"`
	ChatRoom                          *models.ChatRoom    `json:"chat_room"`
	ParticipantIDs                    []string            `json:"participant_ids"`
	EnabledNotificationParticipantIDs []string            `json:"enabled_notification_participant_ids"`
}

func (m *Message) ToJSONRaw() json.RawMessage {
	data, _ := json.Marshal(m)

	return data
}
