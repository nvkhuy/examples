package enums

type SettingInquiry string

var (
	SettingInquiryEditTimeoutType SettingInquiry = "rfq_edit_timeout"
)

func (p SettingInquiry) String() string {
	return string(p)
}

func (p SettingInquiry) IsValid() bool {
	switch p {
	case SettingInquiryEditTimeoutType:
		return true
	}
	return false
}
