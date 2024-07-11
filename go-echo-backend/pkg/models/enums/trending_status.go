package enums

type TrendingStatus string

var (
	TrendingStatusNew           TrendingStatus = "new"
	TrendingStatusPendingReview TrendingStatus = "pending_review"
	TrendingStatusPublished     TrendingStatus = "published"
	TrendingStatusInactive      TrendingStatus = "inactive"
	TrendingStatusDraft         TrendingStatus = "draft"
)

func (p TrendingStatus) String() string {
	return string(p)
}

func (p TrendingStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case TrendingStatusNew:
		name = "New"
	case TrendingStatusPendingReview:
		name = "Pending Review"
	case TrendingStatusPublished:
		name = "Published"
	case TrendingStatusInactive:
		name = "Inactive"
	case TrendingStatusDraft:
		name = "Draft"
	}

	return name
}
