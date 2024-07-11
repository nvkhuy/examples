package enums

type OrderGroupStatus string

var (
	OrderGroupStatusComplete OrderGroupStatus = "complete"
	OrderGroupStatusOnGoing  OrderGroupStatus = "on_going"
)

func (status OrderGroupStatus) String() string {
	return string(status)
}

func (status OrderGroupStatus) DisplayName() string {
	var name = string(status)

	switch status {
	case OrderGroupStatusComplete:
		name = "Complete"
	case OrderGroupStatusOnGoing:
		name = "On going"
	}

	return name
}
