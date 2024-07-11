package enums

type ShopChannel string

var (
	ShopChannelShopify ShopChannel = "shopify"
)

func (l ShopChannel) String() string {
	return string(l)
}

func (l ShopChannel) DisplayName() string {
	switch l {
	case ShopChannelShopify:
		return "Shopify"
	}

	return string(l)
}
