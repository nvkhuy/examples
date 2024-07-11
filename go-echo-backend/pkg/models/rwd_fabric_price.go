package models

import (
	"encoding/json"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"regexp"
	"strconv"
)

const RWDFabricPriceSheetID = "1YHyMOTHwp5W7AO9C4TlX3xyenvFTwrKww6fL6crB6ak"
const RWDKnitFabricPriceSheetName = "RWD KNIT FABRIC"
const RWDWovenFabricPriceSheetName = "RWD WOVEN FABRIC"
const RWDFabricPriceReadSheetFrom = "A3"
const RWDFabricPriceReadSheetTo = "L"

type RWDFabricPriceSlice []*RWDFabricPrice

type RWDFabricPrice struct {
	Model
	RowId         int               `gorm:"index:idx_row_material_type_id,unique" json:"row_id"`
	MaterialType  enums.RWDMaterial `gorm:"index:idx_row_material_type_id,unique" json:"material_type"`
	Description   string            `json:"description"`
	Supplier      string            `json:"supplier"`
	Material      string            `json:"material"`
	Composition   string            `json:"composition"`
	Weight        float64           `json:"weight"`
	CutWidth      float64           `json:"cut_width"`
	Price0To50    float64           `gorm:"column:price_0_to_50" json:"price_0_to_50"`
	Price50To100  float64           `gorm:"column:price_50_to_100" json:"price_50_to_100"`
	Price100To500 float64           `gorm:"column:price_100_to_500" json:"price_100_to_500"`
	PriceAbove500 float64           `gorm:"column:price_above_500" json:"price_above_500"`
	Remark        string            `json:"remark"`
}

func RWDFabricPriceFromSlice(id int, materialType enums.RWDMaterial, args []interface{}) (p *RWDFabricPrice) {
	p = &RWDFabricPrice{}
	p.RowId = id
	p.MaterialType = materialType
	for i, v := range args {
		switch i {
		case 0:
			p.Description = convertInterfaceToString(v)
		case 1:
		case 2:
			p.Supplier = convertInterfaceToString(v)
		case 3:
			p.Material = convertInterfaceToString(v)
		case 4:
			p.Composition = convertInterfaceToString(v)
		case 5:
			p.Weight = convertInterfaceToFloat64(v)
		case 6:
			p.CutWidth = convertInterfaceToFloat64(v)
		case 7:
			p.Price0To50 = convertInterfaceToFloat64(v)
		case 8:
			p.Price50To100 = convertInterfaceToFloat64(v)
		case 9:
			p.Price100To500 = convertInterfaceToFloat64(v)
		case 10:
			p.PriceAbove500 = convertInterfaceToFloat64(v)
		case 11:
			p.Remark = convertInterfaceToString(v)
		}
	}
	return
}

type FetchRWDFabricPriceParams struct {
	JwtClaimsInfo
	MaterialType  enums.RWDMaterial `json:"material_type" query:"material_type" param:"material_type" validate:"omitempty,oneof=knit woven"`
	SpreadsheetId string            `json:"spreadsheet_id" query:"spreadsheet_id" param:"spreadsheet_id"`
	SheetName     string            `json:"sheet_name" query:"sheet_name" param:"sheet_name"`
	From          string            `json:"from" query:"from" param:"from"`
	To            string            `json:"to" query:"to" param:"to"`
	FromNum       int
	ToNum         int
}

func (f *FetchRWDFabricPriceParams) Fetch() *FetchRWDFabricPriceParams {
	if f.SpreadsheetId == "" {
		f.SpreadsheetId = RWDFabricPriceSheetID
	}
	if f.SheetName == "" {
		switch f.MaterialType {
		case enums.KnitMaterial:
			f.SheetName = RWDKnitFabricPriceSheetName
		case enums.WovenMaterial:
			f.SheetName = RWDWovenFabricPriceSheetName
		}
	}
	if f.From == "" {
		f.From = RWDFabricPriceReadSheetFrom
	}
	if f.To == "" {
		f.To = RWDFabricPriceReadSheetTo
	}
	f.FromNum = f.extractInt(f.From)
	f.ToNum = f.extractInt(f.To)
	return f
}

func (f *FetchRWDFabricPriceParams) extractInt(s string) int {
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

type PaginateRWDFabricPriceParams struct {
	PaginationParams
	JwtClaimsInfo
	FabricType  *string  `json:"fabric_type" query:"fabric_type" param:"fabric_type"`
	Material    *string  `json:"material" query:"material" param:"material"`
	Composition *string  `json:"composition" query:"composition" param:"composition"`
	Weight      *float64 `json:"weight" query:"weight" param:"weight"`
	CutWidth    *float64 `json:"cut_width" query:"cut_width" param:"cut_width"`
}

func RWDFabricPriceSliceFromInterface(v interface{}) (slice RWDFabricPriceSlice) {
	var b []byte
	b, _ = json.Marshal(v)
	_ = json.Unmarshal(b, &slice)
	return
}

func (p *PaginateRWDFabricPriceParams) ToVineSlice() (res []string) {
	res = append(res, "material_type")
	if p.FabricType == nil {
		return
	}

	res = append(res, "material")
	if p.Material == nil {
		return
	}

	res = append(res, "composition")
	if p.Composition == nil {
		return
	}

	res = append(res, "weight")
	if p.Weight == nil {
		return
	}

	res = append(res, "cut_width")
	return
}
