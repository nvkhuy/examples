package models

type UserRankingSlice []*UserRanking
type UserRanking struct {
	User
	UserReward
	Ranking       int64 `json:"ranking,omitempty"`
	RankingChange int64 `json:"ranking_change,omitempty"`
}
