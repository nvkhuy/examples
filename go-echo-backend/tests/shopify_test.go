package tests

import (
	"fmt"
	"testing"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
	"github.com/stretchr/testify/assert"
)

func TestShopify_GetAuthURL(t *testing.T) {
	var cfg = initConfig()

	var url = shopify.New(cfg).GetAuthURL("inflow-test-store-loi", nil)

	fmt.Println(url)
}

func TestShopify_ShopGet(t *testing.T) {
	var cfg = initConfig()

	products, err := shopify.New(cfg).NewClient("inflow-test-store-loi", "shpua_5f34282642afc3470bffc19f94b092ab").Shop.Get(nil)
	assert.NoError(t, err)

	helper.PrintJSON(products)
}

func TestShopify_ProductList(t *testing.T) {
	var cfg = initConfig()

	products, err := shopify.New(cfg).NewClient("inflow-test-store-loi", "shpua_5f34282642afc3470bffc19f94b092ab").Product.List(goshopify.ProductListOptions{})
	assert.NoError(t, err)

	helper.PrintJSON(products)
}

func TestShopify_InventoryItemList(t *testing.T) {
	var cfg = initConfig()

	items, err := shopify.New(cfg).NewClient("inflow-test-store-loi", "shpua_5f34282642afc3470bffc19f94b092ab").InventoryItem.List(goshopify.ListOptions{IDs: []int64{47777875525945}})
	assert.NoError(t, err)

	helper.PrintJSON(items)
}

func TestShopify_InventoryLevelList(t *testing.T) {
	var cfg = initConfig()

	items, err := shopify.New(cfg).NewClient("inflow-test-store-loi", "shpua_5f34282642afc3470bffc19f94b092ab").InventoryLevel.List(goshopify.InventoryLevelListOptions{
		LocationIds: []int64{87929717049},
	})
	assert.NoError(t, err)

	helper.PrintJSON(items)
}

func TestShopify_InventoryLevelSet(t *testing.T) {
	var cfg = initConfig()

	items, err := shopify.New(cfg).NewClient("inflow-test-store-loi", "shpua_5f34282642afc3470bffc19f94b092ab").InventoryLevel.Set(goshopify.InventoryLevel{
		InventoryItemId: 47777875525945,
		LocationId:      87929717049,
		Available:       300,
	})
	assert.NoError(t, err)

	helper.PrintJSON(items)
}

func TestShopify_LocationList(t *testing.T) {
	var cfg = initConfig()

	items, err := shopify.New(cfg).NewClient("inflow-test-store-loi", "shpua_5f34282642afc3470bffc19f94b092ab").Location.List(nil)
	assert.NoError(t, err)

	helper.PrintJSON(items)
}

func TestShopify_WebhookList(t *testing.T) {
	var cfg = initConfig()

	items, err := shopify.New(cfg).NewClient("inflow-test-store-loi", "shpua_5f34282642afc3470bffc19f94b092ab").Webhook.List(nil)
	assert.NoError(t, err)

	helper.PrintJSON(items)
}

func TestShopify_WebhookCreate(t *testing.T) {
	var cfg = initConfig()

	var err = shopify.New(cfg).CreateWebhooks(&shopify.ClientInfo{
		ShopName: "inflow-test-store-loi.myshopify.com",
		Token:    "shpua_5f34282642afc3470bffc19f94b092ab",
	})
	assert.NoError(t, err)

}

func TestShopify_DeleteWebhooks(t *testing.T) {
	var cfg = initConfig()

	var err = shopify.New(cfg).DeleteWebhooks(&shopify.ClientInfo{
		ShopName: "inflow-test-store-loi.myshopify.com",
		Token:    "shpua_5f34282642afc3470bffc19f94b092ab",
	})
	assert.NoError(t, err)

}
