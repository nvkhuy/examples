package models

type Bom struct {
	Model

	UserID string `json:"user_id"`

	BulkPurchaseOrderID string `gorm:"uniqueIndex:idx_bom" json:"bulk_purchase_order_id"`
	MS                  string `gorm:"uniqueIndex:idx_bom" json:"ms" csv:"ms"`
	HenryID             string `gorm:"uniqueIndex:idx_bom" json:"henry_id" csv:"henry_id"`

	Image                 Attachments `json:"image"`
	FitSizing             string      `json:"fit_sizing" csv:"fit_sizing"`
	ConfirmedColor        string      `json:"confirmed_color" csv:"confirmed_color"`
	Cat1                  string      `json:"cat1" csv:"cat1"`
	Cat2                  string      `json:"cat2" csv:"cat2"`
	MainFabric            string      `json:"main_fabric" csv:"main_fabric"`
	MainFabricComposition string      `json:"main_fabric_composition" csv:"main_fabric_composition"`
	Buttons               Attachments `json:"buttons" csv:"buttons"`
	Rivet                 string      `json:"rivet" csv:"rivet"`
	ZipperTapeColor       string      `json:"zipper_tape_color" csv:"zipper_tape_color"`
	ZipperTeeth           string      `json:"zipper_teeth" csv:"zipper_teeth"`
	Thread                string      `json:"thread" csv:"thread"`
	OtherDetail           string      `json:"other_detail" csv:"other_detail"`
	Artwork               string      `json:"artwork" csv:"artwork"`
	Lining                string      `json:"lining" csv:"lining"`
	MainLabel             Attachments `json:"main_label" csv:"main_label"`
	SizeLabel             Attachments `json:"size_label" csv:"size_label"`
	WovenPatch            Attachments `json:"woven_patch" csv:"woven_patch"`
	MainHangtagString     Attachments `json:"main_hangtag_string" csv:"main_hangtag_string"`
	SustainableHangtag    string      `json:"sustainable_hangtag" csv:"sustainable_hangtag"`
	CareLabel             Attachments `json:"care_label" csv:"care_label"`
	BarcodeSticker        Attachments `json:"barcode_sticker" csv:"barcode_sticker"`
}
