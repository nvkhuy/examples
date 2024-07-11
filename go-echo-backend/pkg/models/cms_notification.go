package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type CmsNotification struct {
	Model

	Message string `json:"message,omitempty"`
	SeenAt  *int64 `json:"seen_at,omitempty"`

	NotificationType enums.CmsNotificationType `json:"notification_type,omitempty"`
	Metadata         *NotificationMetadata     `json:"metadata,omitempty"`
}

type CmsNotificationForm struct {
	Message          string                    `json:"message,omitempty"`
	NotificationType enums.CmsNotificationType `json:"notification_type,omitempty"`
	Metadata         *NotificationMetadata     `json:"metadata,omitempty"`
}
