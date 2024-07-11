package shopify

type ClientInfo struct {
	ShopName string `json:"shop_name"`
	Token    string `json:"token"`
}

type AuthState struct {
	SuccessURL string                 `json:"success_url"`
	ErrorURL   string                 `json:"error_url"`
	UserID     string                 `json:"user_id"`
	ShopName   string                 `json:"shop_name"`
	Metadata   map[string]interface{} `json:"metadata"`
}
