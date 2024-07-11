package models

import (
	"time"

	"gorm.io/gorm"
)

func (item *InquiryCartItem) BeforeUpdate(tx *gorm.DB) (err error) {
	item.inquiryUpdateAt(tx)
	return
}

func (item *InquiryCartItem) BeforeCreate(tx *gorm.DB) (err error) {
	item.inquiryUpdateAt(tx)
	return
}

func (item *InquiryCartItem) inquiryUpdateAt(db *gorm.DB) {
	if !(item != nil && item.InquiryID != "") {
		return
	}
	db.Model(&Inquiry{}).Where("id = ?", item.InquiryID).Update("updated_at", time.Now().Unix())
	return
}

func (items InquiryCartItems) IDs() []string {
	var IDs []string
	for _, item := range items {
		IDs = append(IDs, item.ID)
	}
	return IDs
}

func (items InquiryCartItems) InquiryIDs() []string {
	var inquiryIDs []string
	for _, item := range items {
		if item.InquiryID != "" {
			inquiryIDs = append(inquiryIDs, item.InquiryID)
		}
	}
	return inquiryIDs
}
