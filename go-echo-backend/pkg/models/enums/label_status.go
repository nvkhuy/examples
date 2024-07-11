package enums

type LabelStatus string

var (
	LabelStatusNew      LabelStatus = "new"
	LabelStatusPrinting LabelStatus = "printing"
	LabelStatusFinished LabelStatus = "finished"
	LabelStatusClosed   LabelStatus = "closed"
	LabelStatusRejected LabelStatus = "rejected"
	LabelStatusApproved LabelStatus = "approved"
)

func (p LabelStatus) String() string {
	return string(p)
}
