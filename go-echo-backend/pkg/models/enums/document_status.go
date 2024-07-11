package enums

type DocumentStatus string

var (
	DocumentStatusNew           DocumentStatus = "new"
	DocumentStatusPendingReview DocumentStatus = "pending_review"
	DocumentStatusPublished     DocumentStatus = "published"
	DocumentStatusInactive      DocumentStatus = "inactive"
	DocumentStatusDraft         DocumentStatus = "draft"
)

func (p DocumentStatus) String() string {
	return string(p)
}

func (p DocumentStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case DocumentStatusNew:
		name = "New"
	case DocumentStatusPendingReview:
		name = "Pending Review"
	case DocumentStatusPublished:
		name = "Published"
	case DocumentStatusInactive:
		name = "Inactive"
	case DocumentStatusDraft:
		name = "Draft"
	}

	return name
}
