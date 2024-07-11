package models

type ReleaseNote struct {
	Model
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate int64  `json:"release_date"`
}
