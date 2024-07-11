package enums

type Team string

var (
	TeamMarketing       Team = "marketing"
	TeamSales           Team = "sales"
	TeamOperator        Team = "operator"
	TeamDev             Team = "dev"
	TeamCustomerService Team = "customer_service"
	TeamDesigner        Team = "designer"
	TeamQA              Team = "qa"
	Finance             Team = "finance"
)

func (t Team) String() string {
	return string(t)
}

func (t Team) DisplayName() string {
	var name = string(t)

	switch t {
	case TeamMarketing:
		name = "Marketing"
	case TeamSales:
		name = "Sales"
	case TeamOperator:
		name = "Operator"
	case TeamDev:
		name = "Dev"
	case TeamCustomerService:
		name = "Customer Service"
	case TeamDesigner:
		name = "Designer"
	case TeamQA:
		name = "QA"
	case Finance:
		name = "Finance"
	}

	return name
}
