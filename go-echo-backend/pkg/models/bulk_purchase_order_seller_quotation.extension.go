package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/price"

func (m BulkPurchaseOrderSellerQuotation) GetQuotedPrice() (p price.Price) {
	if m.FabricCost != nil {
		p = p.AddPtr(m.FabricCost)
	}

	if m.DecorationCost != nil {
		p = p.AddPtr(m.DecorationCost)
	}

	if m.MakingCost != nil {
		p = p.AddPtr(m.MakingCost)
	}

	if m.OtherCost != nil {
		p = p.AddPtr(m.OtherCost)
	}

	return
}

func (bsqs BulkPurchaseOrderSellerQuotations) IDs() []string {
	var bsqIDs = make([]string, 0, len(bsqs))
	for _, bsq := range bsqs {
		bsqIDs = append(bsqIDs, bsq.ID)
	}
	return bsqIDs
}

func (bsqs BulkPurchaseOrderSellerQuotations) BulkIDs() []string {
	bulkIDs := make([]string, 0, len(bsqs))
	for _, bsq := range bsqs {
		bulkIDs = append(bulkIDs, bsq.BulkPurchaseOrderID)
	}
	return bulkIDs
}
