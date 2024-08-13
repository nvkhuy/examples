package models

import "leaderboard-app/common"

type UpdateUserPayload struct {
	ID     string `json:"id" form:"id" query:"id" param:"id" validate:"required"`
	Name   string `json:"name"`
	Points int64  `json:"points"`
}

type PaginateUserParams struct {
	common.PaginateParams
}
