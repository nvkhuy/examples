package enums

type BrandTeamRole string

var (
	BrandTeamRoleManager BrandTeamRole = "manager"
	BrandTeamRoleStaff   BrandTeamRole = "staff"
)

func (role BrandTeamRole) String() string {
	return string(role)
}

func (role BrandTeamRole) DisplayName() string {
	switch role {
	case BrandTeamRoleManager:
		return "Manager"

	case BrandTeamRoleStaff:
		return "Staff"

	}

	return string(role)
}

func (role BrandTeamRole) IsManager() bool {
	return role == BrandTeamRoleManager
}
