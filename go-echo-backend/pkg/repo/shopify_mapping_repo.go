package repo

// import (
// 	"fmt"
// 	"strconv"

// 	goshopify "github.com/bold-commerce/go-shopify/v3"
// 	"github.com/engineeringinflow/inflow-backend/pkg/adb"
// 	"github.com/engineeringinflow/inflow-backend/pkg/helper"
// 	"github.com/engineeringinflow/inflow-backend/pkg/models"
// 	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
// 	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
// 	"github.com/rotisserie/eris"
// 	"gorm.io/gorm/clause"
// )

// type ShopifyMappingRepo struct {
// 	adb *adb.DB
// }

// func NewShopifyMappingRepo(adb *adb.DB) *ShopifyMappingRepo {
// 	return &ShopifyMappingRepo{
// 		adb: adb,
// 	}
// }

// type UpdateOrCreateProductParams struct {
// 	ShopChannel    *models.ShopChannel `json:"shop_channel"`
// 	ShopifyProduct *goshopify.Product  `json:"shopify_product"`
// 	ShopName       string              `json:"shop_name" query:"shop_name" param:"shop_name" form:"shop_name"`
// }

// func (r *ShopifyMappingRepo) UpdateOrCreateProduct(params UpdateOrCreateProductParams) (*models.Product, error) {
// 	var shopChannel = params.ShopChannel
// 	if shopChannel == nil {
// 		var err = r.adb.First(&shopChannel, "shop_name = ?", params.ShopName).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	var product = models.Product{
// 		Name:            params.ShopifyProduct.Title,
// 		SourceProductID: fmt.Sprintf("%d", params.ShopifyProduct.ID),
// 		Description:     params.ShopifyProduct.BodyHTML,
// 		ShopID:          shopChannel.ShopID,
// 	}
// 	product.ID = helper.GenerateXID()

// 	var attachments models.Attachments
// 	for _, image := range params.ShopifyProduct.Images {
// 		attachments = append(attachments, &models.Attachment{
// 			FileURL:      image.Src,
// 			ThumbnailURL: image.Src,
// 			Metadata: map[string]interface{}{
// 				"width":             image.Width,
// 				"height":            image.Height,
// 				"source_product_id": image.ProductID,
// 				"souce_id":          image.ID,
// 			},
// 		})
// 	}
// 	product.Attachments = &attachments

// 	var productVariants []*models.Variant
// 	for _, variant := range params.ShopifyProduct.Variants {
// 		productVariants = append(productVariants, &models.Variant{
// 			Title:                 variant.Title,
// 			Sku:                   variant.Sku,
// 			Price:                 price.NewFromDecimal(*variant.Price),
// 			ProductID:             product.ID,
// 			SourceProductID:       product.SourceProductID,
// 			SourceVariantID:       fmt.Sprintf("%d", variant.ID),
// 			SourceInventoryItemID: fmt.Sprintf("%d", variant.InventoryItemId),
// 			SourceLocationID:      shopChannel.SourcePrimaryLocationID,
// 			Stock:                 variant.InventoryQuantity,
// 		})
// 	}

// 	var updatedProduct models.Product
// 	var sqlResult = r.adb.Model(&updatedProduct).Clauses(clause.Returning{
// 		Columns: []clause.Column{
// 			{Name: "id"},
// 		},
// 	}).Where("source_product_id = ?", fmt.Sprintf("%d", params.ShopifyProduct.ID)).Updates(&product)
// 	if sqlResult.RowsAffected == 0 {
// 		var err = r.adb.CreateFromPayload(&product).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		product.ID = updatedProduct.ID
// 	}

// 	for _, productVariant := range productVariants {
// 		productVariant.ProductID = product.ID
// 		var sqlResult = r.adb.Model(&models.Variant{}).Where("source_variant_id = ?", productVariant.SourceVariantID).Updates(&productVariant)
// 		if sqlResult.RowsAffected == 0 {
// 			var err = r.adb.CreateFromPayload(&productVariant).Error
// 			if err != nil {
// 				return nil, err
// 			}
// 		}

// 		product.Variants = append(product.Variants, productVariant)
// 	}

// 	return &product, nil
// }

// type UpdatePlatformProductParams struct {
// 	ShopifyClientInfo *shopify.ClientInfo `json:"shopify_client_info"`
// 	ShopifyProduct    *goshopify.Product  `json:"shopify_product"`
// }

// func (r *ShopifyMappingRepo) UpdatePlatformProduct(shopifyClientInfo *shopify.ClientInfo, platformProduct *models.Product) (*goshopify.Product, error) {
// 	if platformProduct.SourceProductID == "" {
// 		return nil, eris.Errorf("Source product id is empty")
// 	}

// 	id, err := strconv.ParseInt(platformProduct.SourceProductID, 10, 64)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var shopifyClient = shopify.GetInstance().NewClient(shopifyClientInfo.ShopName, shopifyClientInfo.Token)

// 	for _, variant := range platformProduct.Variants {
// 		if variant.SourceVariantID == "" {
// 			continue
// 		}

// 		id, err := strconv.ParseInt(variant.SourceVariantID, 10, 64)
// 		if err != nil {
// 			continue
// 		}

// 		var price = variant.Price.Decimal()
// 		_, err = shopifyClient.Variant.Update(goshopify.Variant{
// 			ID:    id,
// 			Title: variant.Title,
// 			Price: &price,
// 		})
// 		if err != nil {
// 			continue
// 		}
// 	}

// 	shopifyProduct, err := shopifyClient.Product.Update(goshopify.Product{
// 		ID:       id,
// 		Title:    platformProduct.Name,
// 		BodyHTML: platformProduct.Description,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return shopifyProduct, err
// }
