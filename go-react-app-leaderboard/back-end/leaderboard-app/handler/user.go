package handler

import (
	"leaderboard-app/models"
	"leaderboard-app/repo"
)

func UserRankings() (rankings models.UserRankingSlice, err error) {
	rankings, err = repo.UserRankings()
	return
}

func UpdateUser(payload models.UpdateUserPayload) (user *models.User, err error) {
	user, err = repo.UpdateUserPoints(payload.ID, payload.Points)
	return
}
