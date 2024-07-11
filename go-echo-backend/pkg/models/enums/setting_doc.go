package enums

type SettingDoc string

var (
	SettingDocNDAType SettingDoc = "nda"
	SettingDocTNCType SettingDoc = "tnc" // terms and conditions
)

func (p SettingDoc) String() string {
	return string(p)
}

func (p SettingDoc) IsValid() bool {
	switch p {
	case SettingDocNDAType:
		return true
	case SettingDocTNCType:
		return true
	}
	return false
}
