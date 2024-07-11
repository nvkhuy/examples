package enums

type OrderSource string

var (
	OrderSourceAddToCart OrderSource = "add_to_cart"
	OrderSourceInquiry   OrderSource = "inquiry"
)

func (p OrderSource) String() string {
	return string(p)
}
