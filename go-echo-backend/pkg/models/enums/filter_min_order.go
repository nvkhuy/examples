package enums

type FilterMinOrder int

var (
	FilterRating10  FilterMinOrder = 10
	FilterRating100 FilterMinOrder = 100
	FilterRating500 FilterMinOrder = 500
)

func (p FilterMinOrder) DisplayName() string {
	var name string

	switch p {
	case FilterRating10:
		name = "From 10 items"
	case FilterRating100:
		name = "From 100 items"
	case FilterRating500:
		name = "From 500 items"
	}

	return name
}
