package tests

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"log"
	"testing"
)

func TestComment_PaginateComment(t *testing.T) {
	var app = initApp("dev")
	var params = repo.PaginateCommentsParams{
		OrderByQuery: "c.created_at DESC",
		TargetType:   enums.CommentTargetTypePOInternalNotes,
		TargetID:     "cltaok3b2hjb5isjhot0",
		FileKey:      "",
	}
	result := repo.NewCommentRepo(app.DB).PaginateComment(params)
	log.Println(result)
}
