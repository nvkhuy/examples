package enums

type ProjectStatus string

var (
	// Admin need to preview (these statuses set by admin)
	ProjectStatusPending  ProjectStatus = "pending"
	ProjectStatusRejected ProjectStatus = "rejected"
	ProjectStatusApproved ProjectStatus = "approved"

	// After project was approved by admin
	ProjectStatusClosed    ProjectStatus = "closed"    // project is closed by some reasons
	ProjectStatusOpen      ProjectStatus = "open"      // project is open for freelancer to apply
	ProjectStatusRecruited ProjectStatus = "recruited" // hired freelancer
	ProjectStatusCompleted ProjectStatus = "completed"

	// Base on Figma CMS
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusSuspended ProjectStatus = "suspended"
)

func (p ProjectStatus) String() string {
	return string(p)
}

func (p ProjectStatus) DisplayName() string {
	var name = string(p)

	switch p {
	// Admin need to preview (these statuses set by admin)
	case ProjectStatusPending:
		name = "Pending"
	case ProjectStatusRejected:
		name = "Rejected"
	case ProjectStatusApproved:
		name = "Approved"

	// After project was approved by admin
	case ProjectStatusClosed:
		name = "Closed" // project is closed by some reasons
	case ProjectStatusOpen:
		name = "Open" // project is open for freelancer to apply
	case ProjectStatusRecruited:
		name = "Recruited" // hired freelancer
	case ProjectStatusCompleted:
		name = "Completed"

		// Base on Figma CMS
	case ProjectStatusActive:
		name = "Active"
	case ProjectStatusSuspended:
		name = "Suspended"
	}

	return name
}
