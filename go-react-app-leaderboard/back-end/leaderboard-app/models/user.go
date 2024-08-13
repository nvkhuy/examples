package models

type UserSlice []*User
type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Points int64  `json:"points"`
}
