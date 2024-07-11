package enums

import "strings"

type Domain string

var (
	DomainBuyer   Domain = "buyer" // brands
	DomainSeller  Domain = "seller"
	DomainWebsite Domain = "website"
	DomainAdmin   Domain = "admin"
)

func (d Domain) String() string {
	return string(d)
}

func (d Domain) ToLower() string {
	return strings.TrimSpace(strings.ToLower(string(d)))
}

func (d Domain) DefaultIfInvalid() Domain {
	if d == "" {
		return DomainWebsite
	}

	return d
}

func (d Domain) IsValid() bool {
	switch d {
	case DomainBuyer, DomainSeller, DomainWebsite, DomainAdmin:
		return true
	}
	return false
}
