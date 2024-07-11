package enums

type PostStatus string

var (
	PostStatusNew           PostStatus = "new"
	PostStatusPendingReview PostStatus = "pending_review"
	PostStatusPublished     PostStatus = "published"
	PostStatusInactive      PostStatus = "inactive"
	PostStatusDraft         PostStatus = "draft"
)

func (p PostStatus) String() string {
	return string(p)
}

func (p PostStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case PostStatusNew:
		name = "New"
	case PostStatusPendingReview:
		name = "Pending Review"
	case PostStatusPublished:
		name = "Published"
	case PostStatusInactive:
		name = "Inactive"
	case PostStatusDraft:
		name = "Draft"
	}

	return name
}
