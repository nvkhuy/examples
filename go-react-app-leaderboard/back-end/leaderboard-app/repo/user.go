package repo

import (
	"fmt"
	"leaderboard-app/models"
	"math/rand"
	"sort"
)

var (
	users       models.UserSlice
	rewards     models.UserRewardSlice
	rankChanges map[string]int64
)

func init() {
	rankChanges = make(map[string]int64)
	rewards = models.UserRewardSlice{
		{Reward: "10000$"}, {Reward: "7000$"}, {Reward: "100$"},
	}
	users = models.UserSlice{
		{
			ID:     "1",
			Name:   "Wade Warren",
			Points: 6100,
		},
		{
			ID:     "2",
			Name:   "Dianne Russell",
			Points: 5200,
		},
		{
			ID:     "3",
			Name:   "Esther Howard",
			Points: 5400,
		},
		{
			ID:     "4",
			Name:   "Robert Fox",
			Points: 4900,
		},
		{
			ID:     "5",
			Name:   "Jan Kowalski",
			Points: 4900,
		},
	}
	for i := 0; i < 100; i++ {
		uid := fmt.Sprintf("%d", i+len(users)+1)
		users = append(users, &models.User{
			ID:     uid,
			Name:   fmt.Sprintf("Player %v", uid),
			Points: 10 + int64(rand.Intn(3000)),
		})
	}
}

// UserRankings can use redis for better performance
func UserRankings() (rankings models.UserRankingSlice, err error) {
	for _, u := range users {
		rankings = append(rankings, &models.UserRanking{
			User:          *u,
			RankingChange: rankChanges[u.ID],
		})
	}
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Points > rankings[j].Points
	})
	for i, r := range rankings {
		r.Ranking = int64(i + 1)
	}
	for i := 0; i < minInt(3, len(rankings)); i++ {
		rankings[i].UserReward = *rewards[i]
	}
	return rankings, err
}

func UpdateUserPoints(id string, points int64) (user *models.User, err error) {
	rankBeforeUpdate := getUserRankByPoint(id)
	for _, u := range users {
		if u.ID == id {
			u.Points = points
			user = u
			break
		}
	}
	rankAfterUpdate := getUserRankByPoint(id)
	if rankAfterUpdate != rankBeforeUpdate {
		rankChanges[id] = rankBeforeUpdate - rankAfterUpdate
	}
	return
}

func getUserRankByPoint(id string) (rank int64) {
	sort.Slice(users, func(i, j int) bool {
		return users[i].Points > users[j].Points
	})
	for i, v := range users {
		if v.ID == id {
			return int64(i + 1)
		}
	}
	return 0
}
func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
