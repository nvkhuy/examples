package enums

type Role string

var (
	RoleSuperAdmin Role = "super_admin"
	RoleClient     Role = "client"
	RoleSeller     Role = "seller"

	RoleLeader Role = "leader"
	RoleStaff  Role = "staff"

	RoleManager Role = "manager"
)

func (role Role) String() string {
	return string(role)
}

func (role Role) DisplayName() string {
	switch role {
	case RoleSuperAdmin:
		return "Super Admin"

	case RoleClient:
		return "Client"

	case RoleSeller:
		return "Seller"

	case RoleLeader:
		return "Leader"

	case RoleStaff:
		return "Staff"
	}

	return string(role)
}

func (role Role) IsAdmin() bool {
	return role == RoleSuperAdmin || role == RoleLeader || role == RoleStaff
}

func (role Role) IsSeller() bool {
	return role == RoleSeller
}

func (role Role) IsBuyer() bool {
	return role == RoleClient
}
