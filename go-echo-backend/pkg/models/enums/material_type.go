package enums

type RWDMaterial string

const KnitMaterial RWDMaterial = "knit"
const WovenMaterial RWDMaterial = "woven"

func (r RWDMaterial) IsInvalid() bool {
	switch r {
	case KnitMaterial, WovenMaterial:
		return true
	}
	return false
}

func (r RWDMaterial) String() string {
	return string(r)
}

func (r RWDMaterial) Pointer() *RWDMaterial {
	return &r
}
