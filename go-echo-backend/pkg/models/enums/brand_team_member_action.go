package enums

type BrandMemberAction string
type BrandMemberActions []BrandMemberAction

var (
	BrandMemberActionCreateRFQ  BrandMemberAction = "create_rfq"
	BrandMemberActionUpdateRFQ  BrandMemberAction = "update_rfq"
	BrandMemberActionApproveRFQ BrandMemberAction = "approve_rfq"
	BrandMemberActionRejectRFQ  BrandMemberAction = "reject_rfq"
)

func (b *BrandMemberActions) ToStringSlice() (s []string) {
	for _, v := range *b {
		s = append(s, string(v))
	}
	return
}
