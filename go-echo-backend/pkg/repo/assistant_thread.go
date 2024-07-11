package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
)

type AssistantThreadRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewAssistantThreadRepo(db *db.DB) *AssistantThreadRepo {
	return &AssistantThreadRepo{
		db:     db,
		logger: logger.New("repo/AssistantThreadRepo"),
	}
}

type AssistantThreadParams struct{}

func (r *AssistantThreadRepo) Chat(params AssistantThreadParams) (err error) {
	return
}

type CreateAssistantThreadParams struct{}

func (r *AssistantThreadRepo) CreateChat(params CreateAssistantThreadParams) (err error) {
	return
}
