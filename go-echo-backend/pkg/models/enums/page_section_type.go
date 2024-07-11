package enums

type PageSectionType string

var (
	// For catalog page
	PageSectionTypeCategoryMain      PageSectionType = "category_main"
	PageSectionTypeCatalogDrop       PageSectionType = "catalog_drop"
	PageSectionTypeCatalogCollection PageSectionType = "catalog_collection"
	PageSectionTypeCatalogCloset     PageSectionType = "catalog_closet"

	// For home page
	PageSectionHomeTop        PageSectionType = "home_top"
	PageSectionHomeShop       PageSectionType = "home_shop"
	PageSectionHomeCollection PageSectionType = "home_collection"
	PageSectionHomeClient     PageSectionType = "home_client"

	// For shop custom
	PageSectionShopFeaturedCollection PageSectionType = "shop_featured_collection"
	PageSectionShopPortfolio          PageSectionType = "shop_portfolio"
	PageSectionShopIntro              PageSectionType = "shop_intro"
)

func (p PageSectionType) String() string {
	return string(p)
}

func (p PageSectionType) DisplayName() string {
	var name = string(p)

	switch p {
	case PageSectionTypeCategoryMain:
		name = "Category Main"
	case PageSectionTypeCatalogCollection:
		name = "Catalog Collection"
	case PageSectionTypeCatalogCloset:
		name = "Catalog Closet"

	case PageSectionHomeTop:
		name = "Home Top"
	case PageSectionHomeShop:
		name = "Home Shop"
	case PageSectionHomeCollection:
		name = "Home Collection"
	case PageSectionHomeClient:
		name = "Top brand"
	}

	return name
}
