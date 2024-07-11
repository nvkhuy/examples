package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
)

type BrandTeam struct {
	ID        string `gorm:"unique" json:"id"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	UserID string `gorm:"primaryKey" json:"user_id,omitempty"` // User can in one team only
	TeamID string `gorm:"primaryKey" json:"team_id,omitempty"` // team_id is user_id of user who start this brand in inflow

	Role    enums.BrandTeamRole `json:"role,omitempty"`
	Actions pq.StringArray      `gorm:"type:varchar(200)[]" json:"actions,omitempty"`

	TeamName string `gorm:"-" json:"team_name,omitempty"`
}
