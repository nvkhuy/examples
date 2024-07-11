package enums

type FilterRating int

var (
	FilterRating1 FilterRating = 1
	FilterRating2 FilterRating = 2
	FilterRating3 FilterRating = 3
	FilterRating4 FilterRating = 4
	FilterRating5 FilterRating = 5
)

func (p FilterRating) DisplayName() string {
	var name string

	switch p {
	case FilterRating1:
		name = "From 1 star"
	case FilterRating2:
		name = "From 2 star"
	case FilterRating3:
		name = "From 3 star"
	case FilterRating4:
		name = "From 4 star"
	case FilterRating5:
		name = "From 5 star"
	}

	return name
}
