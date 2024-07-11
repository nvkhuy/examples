package migration

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

var schemas = []interface{}{
	&models.User{},
	&models.Category{},
	&models.Product{},
	&models.ProductAttribute{},
	&models.Variant{},
	&models.ProductReview{},
	&models.Subscriber{},
	&models.FactoryTour{},
	&models.QuantityPriceTier{},
	&models.Collection{},
	&models.CollectionProduct{},
	&models.CollectionProductGroup{},
	&models.Page{},
	&models.PageSection{},
	&models.ChatMessage{},
	&models.ChatRoom{},
	&models.ChatRoomUser{},
	&models.Cart{},
	&models.CartItem{},

	&models.Post{},
	&models.BlogCategory{},

	&models.Inquiry{},
	&models.InquirySeller{},
	&models.InquirySellerSku{},
	&models.InquiryAudit{},
	&models.InquiryCollection{},
	&models.InquiryCartItem{},

	&models.Coordinate{},
	&models.Comment{},
	&models.CardAssignment{},
	&models.CmsNotification{},
	&models.UserNotification{},
	&models.PaymentTransaction{},

	&models.PurchaseOrder{},
	&models.BulkPurchaseOrder{},
	&models.BulkPurchaseOrderItem{},
	&models.BulkPurchaseOrderTracking{},
	&models.PurchaseOrderTracking{},

	&models.SettingTax{},
	&models.SettingSize{},
	&models.SettingBank{},
	&models.SettingSEO{},
	&models.SettingDoc{},

	&models.Address{},
	&models.ProductTypePrice{},
	&models.RWDFabricPrice{},
	&models.SeoTranslation{},
	&models.BusinessProfile{},

	&models.PushToken{},

	&models.UserBank{},
	&models.BrandTeam{},
	&models.AdsVideo{},
	&models.Invoice{},
	&models.Fabric{},
	&models.FabricCollection{},
	&models.FabricInCollection{},

	&models.SysNotification{},
	&models.UserSysNotification{},
	&models.UserDocAgreement{},
	&models.Bom{},

	&models.PurchaseOrderItem{},
	&models.ZaloConfig{},
	&models.TNA{},
	&models.ReleaseNote{},
	&models.Document{},
	&models.DocumentCategory{},
	&models.DocumentTag{},
	&models.TaggedDocument{},
	&models.SettingInquiry{},

	&models.OrderGroup{},
	&models.ProductClass{},

	&models.AsFeaturedIn{},
	&models.OrderCartItem{},
	&models.BulkPurchaseOrderSellerQuotation{},
	&models.Trending{},
	&models.ProductFileUploadInfo{},
}

type Migrator struct {
	db *db.DB
}

func New(db *db.DB) *Migrator {
	return &Migrator{db}
}

// AutoMigrate auto migrate
func (m *Migrator) AutoMigrate() {
	for _, schema := range schemas {
		var err = m.db.SetDebug(!m.db.Configuration.IsProd()).AutoMigrate(schema)
		if err != nil {
			m.db.CustomLogger.Errorf("Auto-migrate %T error: %+v", schema, err)
		}
	}
}
