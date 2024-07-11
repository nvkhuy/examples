package tests

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestParsePhoneNumber(t *testing.T) {
	var phones = []string{
		"+84 231 2341 122",
		"+84 231 2341 121",
	}
	for _, phone := range phones {
		p, err := helper.ParsePhoneNumber(phone, "")
		if err != nil {
			assert.NoError(t, err)
		}
		fmt.Println("phoness", phone, p.GetNationalNumber(), p.GetNumberOfLeadingZeros(), p.GetItalianLeadingZero(), p.GetCountryCode(), p.GetCountryCodeSource(), p.GetExtension(), p.GetPreferredDomesticCarrierCode())
	}

}

func TestToExcel(t *testing.T) {
	data, err := helper.ToExcel([][]interface{}{
		{
			"country_code", "level1_en", "level2_en", "level3_en", "level1_vi", "level2_vi", "level3_vi",
		},
		{
			"VN", "Quang Ngai", "Duc Pho Town", "Nguyen Nghiem Ward", "Quảng Ngãi", "Thị xã Đức Phổ", "Phường Nguyễn Nghiêm",
		},
		{
			"VN", "Quang Ngai", "Duc Pho Town", "Pho An Ward", "Quảng Ngãi", "Thị xã Đức Phổ", func(xlsx *excelize.File, sheetName, cell string) error {
				xlsx.SetCellValue(sheetName, cell, "teasdfasdf111")
				return nil
			},
		},
	})

	assert.NoError(t, err)

	ioutil.WriteFile("test.xlsx", data, 0664)
	helper.PrintJSON(data)
}

func TestExt(t *testing.T) {
	var ext = path.Ext("https://d30oz70zacuj76.cloudfront.net/images/products/IMG_2023-12-22_024412_46.png")

	fmt.Println("*** ext", ext)
}

func TestStructToMap(t *testing.T) {
	var params = hubspot.DealPropertiesForm{
		Amount: "100",
	}
	var result = helper.StructToMap(&params)

	helper.PrintJSON(result)
}

func TestGenerateXID(t *testing.T) {
	var id = helper.GenerateXID()

	fmt.Println("id", id, time.Now().Unix())
}

func TestGenerateQRCode(t *testing.T) {
	var contents = []string{
		"MY240621IF,MY240621IF001,31I4090612-239-26:14",
		"MY240621IF,MY240621IF002,31I4090612-239-28:12,31I4090612-239-32:4",
		"MY240621IF,MY240621IF003,31I4090612-239-30:6,31I4090612-239-24:6",
		"MY240621IF,MY240621IF004,31C4090611-239-S:13",
		"MY240621IF,MY240621IF005,31C4090611-239-M:11,31C4090611-239-XL:4",
		"MY240621IF,MY240621IF006,31C4090611-239-XS:9,31C4090611-239-L:5",
	}

	var folderName = "QR_Code_Label_611-612-613_LBMY"
	var newpath = filepath.Join(".", folderName)
	var err = os.MkdirAll(newpath, os.ModePerm)
	assert.NoError(t, err)

	for index, link := range contents {
		var fileName = contents[index]
		buf, err := helper.GenerateQRCode(helper.GenerateQRCodeOptions{
			Content: link,
		})
		if err != nil {
			fmt.Println("Generate qr error", link, err)
			continue
		}

		var file = fmt.Sprintf("%s/%s.jpeg", newpath, fileName)
		err = ioutil.WriteFile(file, buf.Bytes(), 0664)
		fmt.Printf("Generate qr success %d/%d %s \n", index, len(contents), file)
	}

}
