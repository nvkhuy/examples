package models

type Notification struct {
	Model

	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
	SeenAt  *int64 `json:"seen_at,omitempty"`

	Metadata *NotificationMetadata `json:"metadata,omitempty"`
}
