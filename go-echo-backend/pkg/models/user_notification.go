package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type UserNotification struct {
	Model

	Message          string                     `json:"message,omitempty"`
	SeenAt           *int64                     `json:"seen_at,omitempty"`
	UserID           string                     `json:"user_id,omitempty"`
	NotificationType enums.UserNotificationType `json:"notification_type,omitempty"`
	Metadata         *UserNotificationMetadata  `json:"metadata,omitempty"`
}

type UserNotificationForm struct {
	Message          string                     `json:"message,omitempty"`
	UserID           string                     `json:"user_id,omitempty"`
	NotificationType enums.UserNotificationType `json:"notification_type,omitempty"`
	Metadata         *UserNotificationMetadata  `json:"metadata,omitempty"`
}
