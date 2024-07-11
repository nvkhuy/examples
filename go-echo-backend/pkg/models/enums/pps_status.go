package enums

type PpsStatus string

var (
	PpsStatusNone     PpsStatus = "none"
	PpsStatusWaiting  PpsStatus = "waiting_for_approval"
	PpsStatusApproved PpsStatus = "approved"
	PpsStatusRejected PpsStatus = "rejected"
)

func (p PpsStatus) String() string {
	return string(p)
}
