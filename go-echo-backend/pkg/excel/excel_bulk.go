package excel

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

type BulkPurchaseOrderItem struct {
	SKU        string
	Size       string
	Style      string
	NumSamples int64
	Qty        int64
	UnitPrice  *price.Price
	TotalPrice *price.Price
}

type BulkPurchaseOrder struct {
	ReferenceID            string
	ColIndex               int
	Currency               enums.Currency
	FirstPaymentPercentage float64
	TaxPercentage          float64
	Items                  []*BulkPurchaseOrderItem
}

type ParseBulkPurchaseOrdersResult struct {
	Records []*BulkPurchaseOrder
	Items   []*BulkPurchaseOrderItem
}

func ParseBulkPurchaseOrders(reader io.Reader) (*ParseBulkPurchaseOrdersResult, error) {
	var sheetName = "Summary"
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var result ParseBulkPurchaseOrdersResult

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var datas []*BulkPurchaseOrder
	var styleColIndex = -1
	var skuColIndex = -1
	var sizeColIndex = -1
	var priceColIndex = -1
	var photoshootSamplesColIndex = -1
	var startingRowIndex = -1
	var currency enums.Currency = ""
	var firstPaymentPercentage float64 = 0
	var taxPercentage float64 = 0

	for rowIndex, row := range rows {
		var style = ""
		var sku = ""
		var size = ""
		var numSamples int64 = 0
		var unitPrice = price.NewFromFloat(0)
		var totalQty int64 = 0

		for colIdx, colCell := range row {
			if colCell == "Currency" {
				currency = enums.Currency(row[colIdx+1])
				continue
			}

			if colCell == "1st Payment Percentage" {
				if v, err := strconv.ParseFloat(row[colIdx+1], 64); err == nil {
					firstPaymentPercentage = v
				}
				continue
			}

			if colCell == "Tax Percentage" {
				if v, err := strconv.ParseFloat(row[colIdx+1], 64); err == nil {
					taxPercentage = v
				}
				continue
			}

			if colCell == "M-Style" {
				startingRowIndex = rowIndex
				styleColIndex = colIdx
				continue
			}

			if colCell == "SKU" {
				startingRowIndex = rowIndex
				skuColIndex = colIdx
				continue
			}

			if colCell == "Size" {
				startingRowIndex = rowIndex
				sizeColIndex = colIdx
				continue
			}

			if colCell == "Unit Price" {
				startingRowIndex = rowIndex
				priceColIndex = colIdx
				continue
			}

			if colCell == "Photoshoot Samples" {
				startingRowIndex = rowIndex
				photoshootSamplesColIndex = colIdx
				continue
			}

			if startingRowIndex == -1 {
				continue
			}

			switch colIdx {
			case styleColIndex:
				style = colCell
			case skuColIndex:
				sku = colCell
			case sizeColIndex:
				size = colCell
			case priceColIndex:
				unitPrice = price.NewFromString(colCell)
			case photoshootSamplesColIndex:
				if v, err := strconv.ParseInt(colCell, 10, 64); err == nil {
					numSamples = v
				}

			default:
				if rowIndex == startingRowIndex {
					datas = append(datas, &BulkPurchaseOrder{
						ReferenceID:            strings.TrimSpace(colCell),
						ColIndex:               colIdx,
						Currency:               currency,
						FirstPaymentPercentage: firstPaymentPercentage,
						TaxPercentage:          taxPercentage,
					})
				} else {
					_, index, found := lo.FindIndexOf(datas, func(item *BulkPurchaseOrder) bool {
						return item.ColIndex == colIdx
					})

					if found {
						var item = &BulkPurchaseOrderItem{
							SKU:        sku,
							Size:       size,
							Style:      style,
							NumSamples: numSamples,
							UnitPrice:  &unitPrice,
						}
						if v, err := strconv.ParseInt(colCell, 10, 64); err == nil {
							item.Qty = v
							item.TotalPrice = item.UnitPrice.MultipleInt(item.Qty).ToPtr()
							totalQty += v
						}

						datas[index].Items = append(datas[index].Items, item)

					}
				}

			}

		}

		if startingRowIndex != -1 && startingRowIndex != rowIndex {
			result.Items = append(result.Items, &BulkPurchaseOrderItem{
				SKU:        sku,
				Size:       size,
				Style:      style,
				NumSamples: numSamples,
				UnitPrice:  &unitPrice,
				Qty:        totalQty,
				TotalPrice: unitPrice.MultipleInt(totalQty).ToPtr(),
			})
		}

	}

	result.Records = datas

	return &result, nil
}
func ParseBulkPurchaseOrdersV2(reader io.Reader) (*models.CreateMultipleBulkPurchaseOrdersRequest, error) {
	var sheetName = "Order"
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	var sharedFields = []string{"Currency", "1st Payment Percentage", "Shipping Fee", "Tax Percentage"}
	var mapSharedFields = map[string]string{
		"Currency":               "",
		"1st Payment Percentage": "",
		"Shipping Fee":           "",
		"Tax Percentage":         "",
	}

	var orderItemInfoRequiredFields = []string{"M-Style", "SKU", "Size", "Photoshoot Samples"}
	var mapRequiredFieldCellToIndex = make(map[string]int, len(orderItemInfoRequiredFields))
	var titleHeaderRowIndex = -1
	var bulksPlaceHolder []string
	var mapBulkBreakPoints = make(map[int]struct{})

	for rowIdx, row := range rows {
		for colIdx, colCell := range row {
			if helper.StringContains(sharedFields, colCell) {
				mapSharedFields[colCell] = row[colIdx+1]
			}
			if helper.StringContains(orderItemInfoRequiredFields, colCell) {
				if rowIdx == 0 {
					return nil, errors.New("Invalid template")
				}
				titleHeaderRowIndex = rowIdx
				// validate required header
				for i := colIdx; i < len(orderItemInfoRequiredFields); i++ {
					if !helper.StringContains(orderItemInfoRequiredFields, row[i]) {
						return nil, errors.New("Invalid template")
					}
					if _, ok := mapRequiredFieldCellToIndex[row[i]]; ok {
						return nil, errors.New("Duplicate header")
					}
					mapRequiredFieldCellToIndex[row[i]] = i
				}
				// get bulks placeholder
				for orderColIdx, orderCell := range rows[titleHeaderRowIndex-1] {
					if orderCell != "" {
						bulksPlaceHolder = append(bulksPlaceHolder, orderCell)
						if len(bulksPlaceHolder) > 1 {
							mapBulkBreakPoints[orderColIdx] = struct{}{}
						}
					}
				}
				if len(bulksPlaceHolder) == 0 {
					return nil, errors.New("Invalid template")
				}
				break
			}
		}
		if titleHeaderRowIndex != -1 {
			break
		}
	}
	if titleHeaderRowIndex == -1 {
		return nil, errors.New("Invalid template")
	}
	var styleIdx = mapRequiredFieldCellToIndex["M-Style"]
	var skuIndex = mapRequiredFieldCellToIndex["SKU"]
	var sizeIdx = mapRequiredFieldCellToIndex["Size"]
	var photoSampleIdx = mapRequiredFieldCellToIndex["Photoshoot Samples"]

	var bulksOrderCartItems = make([]models.OrderCartItems, len(bulksPlaceHolder))
	var photoSampleItems models.OrderCartItems

	var style = ""

	for rowIdx, row := range rows {
		if rowIdx <= titleHeaderRowIndex {
			continue
		}
		var items = make(models.OrderCartItems, len(bulksPlaceHolder))
		var itemIndex = 0

		var colorName, size, photoSample, quantity = "", "", 0, 0
	rowLoop:
		for colIdx, colCell := range row {
			switch colIdx {
			case styleIdx:
				if colCell != "" {
					style = colCell
				}
				continue
			case skuIndex:
				if colCell == "" {
					break rowLoop
				}
				colorName = colCell
				continue
			case sizeIdx:
				if colCell == "" {
					break rowLoop
				}
				size = colCell
				continue
			case photoSampleIdx:
				photoSample, err = strconv.Atoi(colCell)
				if err != nil {
					photoSample = 0
				}
				continue
			default:
				cellQuantity, err := strconv.Atoi(colCell)
				if err != nil {
					cellQuantity = 0
				}
				if _, ok := mapBulkBreakPoints[colIdx]; ok {
					itemIndex++
					quantity = cellQuantity
				} else {
					quantity += cellQuantity
				}

				items[itemIndex] = &models.OrderCartItem{
					ColorName: colorName,
					Size:      size,
					Style:     style,
					Qty:       int64(quantity),
				}
			}
		}
		if photoSample != 0 {
			photoSampleItems = append(photoSampleItems, &models.OrderCartItem{
				ColorName: colorName,
				Size:      size,
				Style:     style,
				Qty:       int64(photoSample),
			})
		}
		for i := 0; i < len(bulksPlaceHolder); i++ {
			if items[i] != nil {
				bulksOrderCartItems[i] = append(bulksOrderCartItems[i], items[i])
			}

		}
	}

	var bulks = make([]*models.CreateBulkPurchaseOrderRequest, 0, len(bulksPlaceHolder))
	for idx, clientRefID := range bulksPlaceHolder {
		bulks = append(bulks, &models.CreateBulkPurchaseOrderRequest{
			Items:                          bulksOrderCartItems[idx],
			PurchaseOrderClientReferenceID: clientRefID,
			PurchaseOrderItems:             photoSampleItems,
		})
	}

	var currency = mapSharedFields["Currency"]
	if !helper.StringContains([]string{string(enums.USD), string(enums.VND)}, currency) {
		currency = string(enums.USD)
	}
	firstPaymentPct, err := strconv.ParseFloat(mapSharedFields["1st Payment Percentage"], 64)
	if err != nil {
		firstPaymentPct = 0
	}
	shippingFee := price.NewFromString(mapSharedFields["Shipping Fee"])
	taxPct, err := strconv.ParseFloat(mapSharedFields["Tax Percentage"], 64)
	if err != nil {
		taxPct = 0
	}

	return &models.CreateMultipleBulkPurchaseOrdersRequest{
		Currency:               enums.Currency(currency),
		FirstPaymentPercentage: firstPaymentPct,
		ShippingFee:            shippingFee,
		TaxPercentage:          taxPct,
		Bulks:                  bulks,
	}, nil
}
