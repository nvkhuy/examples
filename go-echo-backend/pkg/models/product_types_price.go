package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
	"strings"
)

const ProductTypePriceSheetID = "1YHyMOTHwp5W7AO9C4TlX3xyenvFTwrKww6fL6crB6ak"
const ProductTypePriceSheetName = "RWD: ITEM + CMPT Price"
const ProductTypePriceReadSheetFrom = "A3"
const ProductTypePriceReadSheetTo = "AE"

var UpdateProductTypePriceColumns = clause.AssignmentColumns([]string{
	"product", "gender", "fabric_type", "feature", "category", "item", "form", "description",
	"knit_material", "knit_composition", "knit_weight",
	"woven_material", "woven_composition", "woven_weight",
	"cut_width",
	"fabric_consumption", "fabric_price_0_to_50", "fabric_price_50_to_100", "fabric_price_100_to_500", "fabric_price_above_500", "fabric_description",
	"cm_price_0_to_50", "cm_price_50_to_100", "cm_price_100_to_500", "cm_price_above_500",
	"total_price_0_to_50", "total_price_50_to_100", "total_price_100_to_500", "total_price_above_500",
})

var UpdateImagesProductTypePriceColumns = clause.AssignmentColumns([]string{"pic_web_url", "pic_catalog_url"})

type ProductTypePrice struct {
	Model
	PICWebURL           string  `json:"pic_web_url"`
	PICCatalogURL       string  `json:"pic_catalog_url"`
	RowId               int     `gorm:"index:idx_row_id,unique" json:"row_id"`
	Product             string  `json:"product" query:"product" form:"product"`
	Gender              string  `json:"gender" query:"gender" form:"gender"`
	FabricType          string  `json:"fabric_type" query:"fabric_type" form:"fabric_type"`
	Feature             string  `json:"feature" query:"feature" form:"feature"`
	Category            string  `json:"category" query:"category" form:"category"`
	Item                string  `json:"item" query:"item" form:"item"`
	Form                string  `json:"form" query:"form" form:"form"`
	Description         string  `json:"description" query:"description" form:"description"`
	KnitMaterial        string  `json:"knit_material" query:"knit_material" form:"knit_material"`
	KnitComposition     string  `json:"knit_composition" query:"knit_composition" form:"knit_composition"`
	KnitWeight          float64 `json:"knit_weight" query:"knit_weight" form:"knit_weight"`
	WovenMaterial       string  `json:"woven_material" query:"woven_material" form:"woven_material"`
	WovenComposition    string  `json:"woven_composition" query:"woven_composition" form:"woven_composition"`
	WovenWeight         float64 `json:"woven_weight" query:"woven_weight" form:"woven_weight"`
	CutWidth            float64 `json:"cut_width" query:"cut_width" form:"cut_width"`
	FabricDescription   string  `json:"fabric_description" query:"fabric_description" form:"fabric_description"`
	FabricConsumption   float64 `json:"fabric_consumption" query:"fabric_consumption" form:"fabric_consumption"`
	FabricPrice0To50    float64 `gorm:"column:fabric_price_0_to_50" json:"fabric_price_0_to_50"`
	FabricPrice50To100  float64 `gorm:"column:fabric_price_50_to_100" json:"fabric_price_50_to_100"`
	FabricPrice100To500 float64 `gorm:"column:fabric_price_100_to_500" json:"fabric_price_100_to_500"`
	FabricPriceAbove500 float64 `gorm:"column:fabric_price_above_500" json:"fabric_price_above_500"`
	CMPrice0To50        float64 `gorm:"column:cm_price_0_to_50" json:"cm_price_0_to_50"`
	CMPrice50To100      float64 `gorm:"column:cm_price_50_to_100" json:"cm_price_50_to_100"`
	CMPrice100To500     float64 `gorm:"column:cm_price_100_to_500" json:"cm_price_100_to_500"`
	CMPriceAbove500     float64 `gorm:"column:cm_price_above_500" json:"cm_price_above_500"`
	TotalPrice0To50     float64 `gorm:"column:total_price_0_to_50" json:"total_price_0_to_50"`
	TotalPrice50To100   float64 `gorm:"column:total_price_50_to_100" json:"total_price_50_to_100"`
	TotalPrice100To500  float64 `gorm:"column:total_price_100_to_500" json:"total_price_100_to_500"`
	TotalPriceAbove500  float64 `gorm:"column:total_price_above_500" json:"total_price_above_500"`
}

type ProductTypePriceSlice []*ProductTypePrice

func (p ProductTypePriceSlice) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *ProductTypePriceSlice) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, p)
	}
	return errors.New(fmt.Sprint("Failed to unmarshal JSON from DB", src))
}

func ProductTypePriceFromSlice(id int, args []interface{}) (p *ProductTypePrice) {
	p = &ProductTypePrice{}
	p.RowId = id
	for i, v := range args {
		switch i {
		case 2:
			p.Product = convertInterfaceToString(v)
		case 3:
			p.Gender = convertInterfaceToString(v)
		case 4:
			p.FabricType = convertInterfaceToString(v)
		case 5:
			p.Feature = convertInterfaceToString(v)
		case 6:
			p.Category = convertInterfaceToString(v)
		case 7:
			p.Item = convertInterfaceToString(v)
		case 8:
			p.Form = convertInterfaceToString(v)
		case 9:
			p.Description = convertInterfaceToString(v)
		case 10:
			p.KnitMaterial = convertInterfaceToString(v)
		case 11:
			p.KnitComposition = convertInterfaceToString(v)
		case 12:
			p.KnitWeight = convertInterfaceToFloat64(v)
		case 13:
			p.WovenMaterial = convertInterfaceToString(v)
		case 14:
			p.WovenComposition = convertInterfaceToString(v)
		case 15:
			p.WovenWeight = convertInterfaceToFloat64(v)
		case 16:
			p.CutWidth = convertInterfaceToFloat64(v)
		case 17:
			p.FabricDescription = convertInterfaceToString(v)
		case 18:
			p.FabricConsumption = convertInterfaceToFloat64(v)
		case 19:
			p.FabricPrice0To50 = convertInterfaceToFloat64(v)
		case 20:
			p.FabricPrice50To100 = convertInterfaceToFloat64(v)
		case 21:
			p.FabricPrice100To500 = convertInterfaceToFloat64(v)
		case 22:
			p.FabricPriceAbove500 = convertInterfaceToFloat64(v)
		case 23:
			p.CMPrice0To50 = convertInterfaceToFloat64(v)
		case 24:
			p.CMPrice50To100 = convertInterfaceToFloat64(v)
		case 25:
			p.CMPrice100To500 = convertInterfaceToFloat64(v)
		case 26:
			p.CMPriceAbove500 = convertInterfaceToFloat64(v)
		case 27:
			p.TotalPrice0To50 = convertInterfaceToFloat64(v)
		case 28:
			p.TotalPrice50To100 = convertInterfaceToFloat64(v)
		case 29:
			p.TotalPrice100To500 = convertInterfaceToFloat64(v)
		case 30:
			p.TotalPriceAbove500 = convertInterfaceToFloat64(v)
		}
	}
	return p
}

func convertInterfaceToString(v interface{}) (vs string) {
	var ok bool
	vs, ok = v.(string)
	if ok == true {
		vs = strings.TrimSpace(vs)
	}
	return
}

func convertInterfaceToFloat64(v interface{}) (vf float64) {
	if vs, ok := v.(string); ok {
		vs = strings.Replace(vs, ".", "", -1)
		vs = strings.Replace(vs, ",", ".", -1)
		pattern := `[0-9]+(\.[0-9]+)?`
		regex := regexp.MustCompile(pattern)
		match := regex.FindString(vs)
		if match != "" {
			vf, _ = strconv.ParseFloat(match, 64)
			return
		}
	}
	return
}

type PaginateProductTypesPriceParams struct {
	PaginationParams
	JwtClaimsInfo
	Product     *string `json:"product" query:"product" param:"product"`
	Gender      *string `json:"gender" query:"gender" param:"gender"`
	FabricType  *string `json:"fabric_type" query:"fabric_type" param:"fabric_type"`
	Feature     *string `json:"feature" query:"feature" param:"feature"`
	Category    *string `json:"category" query:"category" param:"category"`
	Item        *string `json:"item" query:"item" param:"item"`
	Form        *string `json:"form" query:"form" param:"form"`
	Description *string `json:"description" query:"description" param:"description"`

	KnitMaterial    *string  `json:"knit_material" query:"knit_material" param:"knit_material"`
	KnitComposition *string  `json:"knit_composition" query:"knit_composition" param:"knit_composition"`
	KnitWeight      *float64 `json:"knit_weight" query:"knit_weight" param:"knit_weight"`

	WovenMaterial    *string  `json:"woven_material" query:"woven_material" param:"woven_material"`
	WovenComposition *string  `json:"woven_composition" query:"woven_composition" param:"woven_composition"`
	WovenWeight      *float64 `json:"woven_weight" query:"woven_weight" param:"woven_weight"`

	CutWidth          *float64 `json:"cut_width" query:"cut_width" param:"cut_width"`
	FabricDescription *string  `json:"fabric_description" query:"fabric_description" param:"fabric_description"`
	FabricConsumption *float64 `json:"fabric_consumption" query:"fabric_consumption" param:"fabric_consumption"`
}

func (p *PaginateProductTypesPriceParams) ToVineSlice() (res []string) {
	res = append(res, "product")
	if p.Product == nil {
		return
	}

	res = append(res, "gender")
	if p.Gender == nil {
		return
	}

	res = append(res, "fabric_type")
	if p.FabricType == nil {
		return
	}

	res = append(res, "feature")
	if p.Feature == nil {
		return
	}

	res = append(res, "category")
	if p.Category == nil {
		return
	}

	res = append(res, "item")
	if p.Item == nil {
		return
	}

	res = append(res, "form")
	if p.Form == nil {
		return
	}

	res = append(res, "description")
	if p.Description == nil {
		return
	}

	res = append(res, "knit_material")
	if p.KnitMaterial == nil {
		return
	}

	res = append(res, "knit_composition")
	if p.KnitComposition == nil {
		return
	}

	res = append(res, "knit_weight")
	if p.KnitWeight == nil {
		return
	}

	res = append(res, "woven_material")
	if p.WovenMaterial == nil {
		return
	}

	res = append(res, "woven_composition")
	if p.WovenComposition == nil {
		return
	}

	res = append(res, "woven_weight")
	if p.WovenWeight == nil {
		return
	}

	res = append(res, "cut_width")
	if p.CutWidth == nil {
		return
	}

	res = append(res, "fabric_description")
	if p.FabricDescription == nil {
		return
	}

	res = append(res, "fabric_consumption")
	return
}

type FetchProductTypesPriceParams struct {
	JwtClaimsInfo
	SpreadsheetId string `json:"spreadsheet_id" query:"spreadsheet_id" param:"spreadsheet_id"`
	SheetName     string `json:"sheet_name" query:"w" param:"sheet_name"`
	From          string `json:"from" query:"from" param:"from"`
	To            string `json:"to" query:"to" param:"to"`
	FromNum       int
	ToNum         int
}

func (f *FetchProductTypesPriceParams) Fetch() *FetchProductTypesPriceParams {
	if f.SpreadsheetId == "" {
		f.SpreadsheetId = ProductTypePriceSheetID
	}
	if f.SheetName == "" {
		f.SheetName = ProductTypePriceSheetName
	}
	if f.From == "" {
		f.From = ProductTypePriceReadSheetFrom
	}
	if f.To == "" {
		f.To = ProductTypePriceReadSheetTo
	}
	f.FromNum = f.extractInt(f.From)
	f.ToNum = f.extractInt(f.To)
	return f
}

func (f *FetchProductTypesPriceParams) extractInt(s string) int {
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(s)
	if match != "" {
		number, err := strconv.Atoi(match)
		if err == nil {
			return number
		}
	}
	return 0
}

func ProductTypesPriceFromInterface(i interface{}) (p *ProductTypePrice) {
	b, _ := json.Marshal(i)
	_ = json.Unmarshal(b, &p)
	return
}

func ProductTypesPriceSliceFromInterface(i interface{}) (p ProductTypePriceSlice) {
	b, _ := json.Marshal(i)
	_ = json.Unmarshal(b, &p)
	return
}

type PatchSheetImageURLParams struct {
	JwtClaimsInfo
	Concurrency int    `json:"concurrency" form:"concurrency"`
	From        int    `json:"from" form:"from"`
	To          int    `json:"to" form:"to"`
	Token       string `json:"token" form:"token" validate:"required"`
	Cookie      string `json:"cookie" form:"cookie" validate:"required"`
}

func (f *PatchSheetImageURLParams) Fetch() *PatchSheetImageURLParams {
	if f.Concurrency == 0 {
		f.Concurrency = 20
	}
	if f.From <= 0 {
		f.From = 3
	}
	if f.To <= 0 {
		f.To = 649
	}
	return f
}

type PaginateProductTypesPriceQuoteParams struct {
	PaginationParams
	JwtClaimsInfo
	Product     *string `json:"product" query:"product" param:"product"`
	Gender      *string `json:"gender" query:"gender" param:"gender"`
	FabricType  *string `json:"fabric_type" query:"fabric_type" param:"fabric_type"`
	Feature     *string `json:"feature" query:"feature" param:"feature"`
	Category    *string `json:"category" query:"category" param:"category"`
	Item        *string `json:"item" query:"item" param:"item"`
	Form        *string `json:"form" query:"form" param:"form"`
	Description *string `json:"description" query:"description" param:"description"`

	Material    *string  `json:"material" query:"material" param:"material"`
	Composition *string  `json:"composition" query:"composition" param:"composition"`
	Weight      *float64 `json:"weight" query:"weight" param:"weight"`
	CutWidth    *float64 `json:"cut_width" query:"cut_width" param:"cut_width"`
}
