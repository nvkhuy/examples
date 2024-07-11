package enums

type AddressType string

var (
	AddressTypePrimary AddressType = "primary"
)

func (p AddressType) String() string {
	return string(p)
}

func (p AddressType) DisplayName() string {
	var name = string(p)

	switch p {
	case AddressTypePrimary:
		name = "Primary"
	}

	return name
}
