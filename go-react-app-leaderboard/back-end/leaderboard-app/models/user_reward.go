package models

type UserRewardSlice []*UserReward
type UserReward struct {
	Reward string `json:"reward,omitempty"`
}
