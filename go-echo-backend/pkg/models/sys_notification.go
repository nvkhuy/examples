package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type SysNotification struct {
	Model
	Name    string                    `json:"name,omitempty"`
	Type    enums.SysNotificationType `json:"type,omitempty"`
	Message string                    `json:"message,omitempty"`
}
