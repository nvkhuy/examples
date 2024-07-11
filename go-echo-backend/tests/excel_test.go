package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/engineeringinflow/inflow-backend/pkg/excel"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestExcel_ParseBulkPurchaseOrders(t *testing.T) {
	data, err := ioutil.ReadFile("./Order Summary_Denim March_Inflow.xlsx")
	assert.NoError(t, err)

	results, err := excel.ParseBulkPurchaseOrders(bytes.NewBuffer(data))
	assert.NoError(t, err)

	helper.PrintJSON(results)
}

func TestExcel_ParseBOM(t *testing.T) {
	data, err := ioutil.ReadFile("./MAY24_DENIM_BOM_IF.xlsx")
	assert.NoError(t, err)

	_, err = excel.ParseBOM(bytes.NewBuffer(data))
	assert.NoError(t, err)

}

func TestExcel_ParseBulkPurchaseOrdersV2(t *testing.T) {
	data, err := os.ReadFile("./files/Order Summary_Inflow.xlsx")
	assert.NoError(t, err)
	result, err := excel.ParseBulkPurchaseOrdersV2(bytes.NewBuffer(data))
	fmt.Printf("========currency: %s\n", result.Currency)
	fmt.Printf("========1st payment: %v\n", result.FirstPaymentPercentage)
	fmt.Printf("========shipping fee: %d\n", result.ShippingFee.ToInt64())
	fmt.Printf("========tax percentage: %v\n", result.TaxPercentage)

	fmt.Printf("========bulks-length: %d\n", len(result.Bulks))
	fmt.Printf("========po-items-length: %d\n", len(result.Bulks[0].PurchaseOrderItems))
	fmt.Printf("========bulks-items-length: %d\n", len(result.Bulks[1].Items))
	for _, item := range result.Bulks[3].PurchaseOrderItems {
		fmt.Print(item.Size, "\t")
		fmt.Print(item.ColorName, "\t")
		fmt.Print(item.Style, "\t")
		fmt.Print(item.Qty, "\t")
		fmt.Println()
	}
	fmt.Println()

	for _, item := range result.Bulks[0].Items {
		fmt.Print(item.ColorName, "\t")
	}
	fmt.Println()
	for _, item := range result.Bulks[0].Items {
		fmt.Print(item.Style, "\t")
	}
	fmt.Println()
	for _, item := range result.Bulks[0].Items {
		fmt.Print(item.Size, "\t")
	}
	fmt.Println()
	for _, item := range result.Bulks[0].Items {
		fmt.Print(item.Qty, "\t")
	}
	assert.Error(t, err)

}
