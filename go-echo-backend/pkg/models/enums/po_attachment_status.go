package enums

type PoAttachmentStatus string

var (
	PoAttachmentStatusRejected           PoAttachmentStatus = "rejected"
	PoAttachmentStatusWaitingForApproval PoAttachmentStatus = "waiting_for_approval"
	PoAttachmentStatusApproved           PoAttachmentStatus = "approved"
)

func (p PoAttachmentStatus) String() string {
	return string(p)
}

func (p PoAttachmentStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case PoAttachmentStatusRejected:
		name = "Rejected"
	case PoAttachmentStatusWaitingForApproval:
		name = "Waiting for approval"
	case PoAttachmentStatusApproved:
		name = "Approved"
	}

	return name
}
