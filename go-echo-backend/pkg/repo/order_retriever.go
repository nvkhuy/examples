package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type OrderRetriever struct {
	db *db.DB
}

func NewOrderRetriever(db *db.DB) *OrderRetriever {
	return &OrderRetriever{db: db}
}

func (r *OrderRetriever) RetrieveByRFQs(ids []string) (models.Inquiries, models.PurchaseOrders, models.BulkPurchaseOrders, error) {
	var inquiries models.Inquiries
	if err := r.db.Select("ID").Find(&inquiries, "id IN ?", ids).Error; err != nil {
		return nil, nil, nil, err
	}
	if len(inquiries) != len(ids) {
		return nil, nil, nil, errs.ErrInquiryNotFound
	}

	var purchaseOrders models.PurchaseOrders
	if err := r.db.Select("ID").Find(&purchaseOrders, "inquiry_id IN ?", inquiries.IDs()).Error; err != nil {
		return nil, nil, nil, err
	}

	var bulkPurchaseOrders models.BulkPurchaseOrders
	if err := r.db.Select("ID").Find(&bulkPurchaseOrders, "purchase_order_id IN ?", purchaseOrders.IDs()).Error; err != nil {
		return nil, nil, nil, err
	}

	return inquiries, purchaseOrders, bulkPurchaseOrders, nil
}

func (r *OrderRetriever) RetrieveBySamples(ids []string) (models.Inquiries, models.PurchaseOrders, models.BulkPurchaseOrders, error) {
	var purchaseOrders models.PurchaseOrders
	if err := r.db.Select("ID", "InquiryID").Find(&purchaseOrders, "id IN ?", ids).Error; err != nil {
		return nil, nil, nil, err
	}
	if len(purchaseOrders) != len(ids) {
		return nil, nil, nil, errs.ErrPONotFound
	}

	var inquiries models.Inquiries
	if err := r.db.Select("ID").Find(&inquiries, "id IN ?", purchaseOrders.InquiryIDs()).Error; err != nil {
		return nil, nil, nil, err
	}

	var bulkPurchaseOrders models.BulkPurchaseOrders
	if err := r.db.Select("ID").Find(&bulkPurchaseOrders, "purchase_order_id IN ?", purchaseOrders.IDs()).Error; err != nil {
		return nil, nil, nil, err
	}

	return inquiries, purchaseOrders, bulkPurchaseOrders, nil

}

func (r *OrderRetriever) RetrieveByBulks(ids []string) (models.Inquiries, models.PurchaseOrders, models.BulkPurchaseOrders, error) {
	var bulkPurchaseOrders models.BulkPurchaseOrders
	if err := r.db.Select("ID", "PurchaseOrderID", "InquiryID").Find(&bulkPurchaseOrders, "id IN ?", ids).Error; err != nil {
		return nil, nil, nil, err
	}
	if len(bulkPurchaseOrders) != len(ids) {
		return nil, nil, nil, errs.ErrBulkPoNotFound
	}

	var purchaseOrders models.PurchaseOrders
	if err := r.db.Select("ID").Find(&purchaseOrders, "id IN ?", bulkPurchaseOrders.PurchaseOrderIDs()).Error; err != nil {
		return nil, nil, nil, err
	}

	var inquiries models.Inquiries
	if err := r.db.Select("ID").Find(&inquiries, "id IN ?", bulkPurchaseOrders.InquiryIDs()).Error; err != nil {
		return nil, nil, nil, err
	}

	return inquiries, purchaseOrders, bulkPurchaseOrders, nil
}
