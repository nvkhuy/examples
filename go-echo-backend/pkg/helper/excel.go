package helper

import (
	"encoding/csv"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

type HandlerFunc func(xlsx *excelize.File, sheetName, cell string) error

func ToExcel(slices [][]interface{}) ([]byte, error) {
	xlsx := excelize.NewFile()
	sheetName := "Sheet1"
	index, err := xlsx.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	for idx, row := range slices {
		for i, r := range row {
			cell, _ := excelize.CoordinatesToCellName(i+1, idx+1)
			if f, ok := r.(func(xlsx *excelize.File, sheetName, cell string) error); ok {
				if e := f(xlsx, sheetName, cell); e != nil {
					return nil, e
				}
			} else {
				xlsx.SetCellValue(sheetName, cell, r)
			}
		}
	}
	xlsx.SetActiveSheet(index)

	err = AutoFitExcelColumnWidth(xlsx, sheetName)
	if err != nil {
		return nil, err
	}

	xlsxBuffer, err := xlsx.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return xlsxBuffer.Bytes(), nil
}

func BytesFromExcel(csvData []byte) ([]byte, error) {
	r := csv.NewReader(strings.NewReader(string(csvData)))
	content, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	xlsx := excelize.NewFile()
	sheetName := "Sheet1"
	index, err := xlsx.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	for key, _da := range content {
		axis := fmt.Sprintf("A%d", key+1)
		_ = xlsx.SetSheetRow(sheetName, axis, &_da)
	}

	xlsx.SetActiveSheet(index)

	err = AutoFitExcelColumnWidth(xlsx, sheetName)
	if err != nil {
		return nil, err
	}

	xlsxBuffer, err := xlsx.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return xlsxBuffer.Bytes(), nil
}

func AutoFitExcelColumnWidth(xlsx *excelize.File, sheetName string) error {
	cols, err := xlsx.GetCols(sheetName)
	if err != nil {
		return err
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return err
		}

		xlsx.SetColWidth(sheetName, name, name, float64(largestWidth))
	}

	return nil
}
