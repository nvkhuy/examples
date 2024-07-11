package models

type UserSysNotification struct {
	Model

	SysNotificationID string                    `json:"sys_notification_id,omitempty"`
	Message           string                    `json:"message,omitempty"`
	SeenAt            *int64                    `json:"seen_at,omitempty"`
	UserID            string                    `json:"user_id,omitempty"`
	Metadata          *UserNotificationMetadata `json:"metadata,omitempty"`
}
