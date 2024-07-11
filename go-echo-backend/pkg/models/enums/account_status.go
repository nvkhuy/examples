package enums

type AccountStatus string

var (
	AccountStatusPendingReview AccountStatus = "pending_review"
	AccountStatusActive        AccountStatus = "active"
	AccountStatusRejected      AccountStatus = "rejected"
	AccountStatusInactive      AccountStatus = "inactive"
	AccountStatusSuspended     AccountStatus = "suspended"
)

func (p AccountStatus) String() string {
	return string(p)
}

func (p AccountStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case AccountStatusPendingReview:
		name = "Pending"
	case AccountStatusActive:
		name = "Active"
	case AccountStatusRejected:
		name = "Rejected"
	case AccountStatusInactive:
		name = "Inactive"
	case AccountStatusSuspended:
		name = "Suspended"
	}

	return name
}
