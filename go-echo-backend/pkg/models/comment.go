package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

// Comment's model
type Comment struct {
	Model

	Content     string                  `json:"content,omitempty" validate:"required"`
	TargetType  enums.CommentTargetType `json:"target_type,omitempty" validate:"required"`
	TargetID    string                  `json:"target_id,omitempty" validate:"required"`
	UserID      string                  `gorm:"primaryKey" json:"user_id,omitempty" validate:"required"`
	Attachments *Attachments            `json:"attachments,omitempty"`
	FileKey     string                  `json:"file_key,omitempty"`
	SeenAt      *int64                  `json:"seen_at,omitempty"`

	PurchaseOrder *PurchaseOrder `gorm:"-" json:"purchase_order,omitempty"`
	Inquiry       *Inquiry       `gorm:"-" json:"inquiry,omitempty"`
	Participants  []*User        `gorm:"-" json:"participants,omitempty"`

	User *User `gorm:"-" json:"user,omitempty"`
}

type CommentCreateForm struct {
	JwtClaimsInfo

	Content        string                  `json:"content,omitempty" validate:"required"`
	TargetType     enums.CommentTargetType `json:"target_type,omitempty" validate:"required"`
	TargetID       string                  `json:"target_id,omitempty" validate:"required"`
	UserID         string                  `json:"user_id,omitempty"`
	FileKey        string                  `json:"file_key,omitempty"`
	Attachments    *Attachments            `json:"attachments,omitempty"`
	MentionUserIDs []string                `json:"mention_user_ids"`
	Message        string                  `json:"message"`
}

type ContentCommentCreateForm struct {
	Content     string       `json:"content,omitempty" validate:"required"`
	Attachments *Attachments `json:"attachments,omitempty"`
}

type CommentStatusCountItem struct {
	FileKey     string `json:"file_key"`
	UnseenCount int64  `json:"unseen_count"`
}
