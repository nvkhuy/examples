package models

import (
	"encoding/json"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

// Website
const WebsiteSeoTranslationSheetID = "1zYT1sK-yICGqDeUQ1KbzD_4Ce2bd-WA26JpIq2q96Fg"
const WebsiteSeoTranslationSheetENName = "EN"
const WebsiteSeoTranslationSheetVIName = "VI"
const WebsiteSeoTranslationReadSheetFrom = "A1"
const WebsiteSeoTranslationReadSheetTo = "B"

// Buyer
const BuyerSeoTranslationSheetID = "1f7f2oEIXoB8i2BhEWrN7Nsy9LQFFoUe-X3tRokXARn0"
const BuyerSeoTranslationSheetENName = "EN"
const BuyerSeoTranslationSheetVIName = "VI"
const BuyerSeoTranslationReadSheetFrom = "A1"
const BuyerSeoTranslationReadSheetTo = "B"

// Seller
const SellerSeoTranslationSheetID = "1GVXgCRg-PdOS9t59ynUzPMVv3WYA7-WS4CA5wHFl3gs"
const SellerSeoTranslationSheetENName = "EN"
const SellerSeoTranslationSheetVIName = "VI"
const SellerSeoTranslationReadSheetFrom = "A1"
const SellerSeoTranslationReadSheetTo = "B"

// Admin
const AdminSeoTranslationSheetID = "16WEA6hTD2AFWNJHfSjwykoWD6LLtANyTmlOLDEEav6M"
const AdminSeoTranslationSheetENName = "EN"
const AdminSeoTranslationSheetVIName = "VI"
const AdminSeoTranslationReadSheetFrom = "A1"
const AdminSeoTranslationReadSheetTo = "B"

type SeoTranslation struct {
	Model
	Keyword string       `gorm:"index:idx_domain_keyword,unique" json:"keyword"`
	Domain  enums.Domain `gorm:"index:idx_domain_keyword,unique" json:"domain"`
	VI      string       `json:"vi"`
	EN      string       `json:"en"`
}

type SeoTranslationSlice []*SeoTranslation

type GetSEOTranslationForm struct {
	Domain enums.Domain `json:"domain" query:"domain" param:"domain"`
}

type FetchSeoTranslationParams struct {
	Domain        enums.Domain `json:"domain" query:"domain" param:"domain" validate:"required,oneof=buyer seller website admin"`
	SpreadsheetId string       `json:"spreadsheet_id" query:"spreadsheet_id" param:"spreadsheet_id"`
	SheetENName   string       `json:"sheet_en_name" query:"sheet_en_name" param:"sheet_name"`
	SheetVIName   string       `json:"sheet_vi_name" query:"sheet_vi_name" param:"sheet_vi_name"`
	From          string       `json:"from" query:"from" param:"from"`
	To            string       `json:"to" query:"to" param:"to"`
}

func (f *FetchSeoTranslationParams) Fetch() *FetchSeoTranslationParams {
	if f.SpreadsheetId == "" {
		switch f.Domain {
		case enums.DomainBuyer:
			f.SpreadsheetId = BuyerSeoTranslationSheetID
		case enums.DomainSeller:
			f.SpreadsheetId = SellerSeoTranslationSheetID
		case enums.DomainWebsite:
			f.SpreadsheetId = WebsiteSeoTranslationSheetID
		case enums.DomainAdmin:
			f.SpreadsheetId = AdminSeoTranslationSheetID
		}
	}
	if f.SheetENName == "" {
		switch f.Domain {
		case enums.DomainBuyer:
			f.SheetENName = BuyerSeoTranslationSheetENName
		case enums.DomainSeller:
			f.SheetENName = SellerSeoTranslationSheetENName
		case enums.DomainWebsite:
			f.SheetENName = WebsiteSeoTranslationSheetENName
		case enums.DomainAdmin:
			f.SheetENName = AdminSeoTranslationSheetENName
		}
	}
	if f.SheetVIName == "" {
		switch f.Domain {
		case enums.DomainBuyer:
			f.SheetVIName = BuyerSeoTranslationSheetVIName
		case enums.DomainSeller:
			f.SheetVIName = SellerSeoTranslationSheetVIName
		case enums.DomainWebsite:
			f.SheetVIName = WebsiteSeoTranslationSheetVIName
		case enums.DomainAdmin:
			f.SheetVIName = AdminSeoTranslationSheetVIName
		}
	}
	if f.From == "" {
		switch f.Domain {
		case enums.DomainBuyer:
			f.From = BuyerSeoTranslationReadSheetFrom
		case enums.DomainSeller:
			f.From = SellerSeoTranslationReadSheetFrom
		case enums.DomainWebsite:
			f.From = WebsiteSeoTranslationReadSheetFrom
		case enums.DomainAdmin:
			f.From = AdminSeoTranslationReadSheetFrom
		}
	}
	if f.To == "" {
		switch f.Domain {
		case enums.DomainBuyer:
			f.To = BuyerSeoTranslationReadSheetTo
		case enums.DomainSeller:
			f.To = SellerSeoTranslationReadSheetTo
		case enums.DomainWebsite:
			f.To = WebsiteSeoTranslationReadSheetTo
		case enums.DomainAdmin:
			f.To = AdminSeoTranslationReadSheetTo
		}
	}
	return f
}

func SeoTranslationFromSliceInterface(domain enums.Domain, lang enums.LanguageCode, args []interface{}) (s *SeoTranslation) {
	s = &SeoTranslation{}
	s.Domain = domain
	for i, v := range args {
		switch i {
		case 0:
			s.Keyword = convertInterfaceToString(v)
		case 1:
			switch lang {
			case enums.LanguageCodeEnglish:
				s.EN = convertInterfaceToString(v)
			case enums.LanguageCodeVietnam:
				s.VI = convertInterfaceToString(v)
			}
		}
	}
	return
}

func SeoTranslationFromInterface(i interface{}) (s *SeoTranslation) {
	b, _ := json.Marshal(i)
	_ = json.Unmarshal(b, &s)
	return
}
