package shopify

import (
	"encoding/json"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

type AuthorizeShopResponse struct {
	Token string
	Shop  *goshopify.Shop
}

type AuthorizeShopParams struct {
	Shop      string `json:"shop" query:"shop" param:"shop" validate:"required"`
	Code      string `json:"code" query:"code" param:"code" validate:"required"`
	HMAC      string `json:"hmac" query:"hmac" param:"hmac"`
	Host      string `json:"host" query:"host" param:"host"`
	Timestamp int    `json:"timestamp" query:"timestamp" param:"timestamp"`
}

func (a *App) AuthorizeShop(params AuthorizeShopParams) (*AuthorizeShopResponse, error) {
	token, err := a.GetAccessToken(params.Shop, params.Code)
	if err != nil {
		return nil, err
	}

	shop, err := a.NewClient(params.Shop, token).Shop.Get(nil)
	if err != nil {
		return nil, err
	}

	var resp = AuthorizeShopResponse{
		Token: token,
		Shop:  shop,
	}

	return &resp, nil
}

func (a *App) GetAuthURL(shopName string, state *AuthState) string {
	var stateStr = ""
	if state != nil {
		data, _ := json.Marshal(state)
		stateStr = string(data)
	}

	return a.AuthorizeUrl(shopName, stateStr)

}
