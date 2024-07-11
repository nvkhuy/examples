package excel

import (
	"fmt"
	"io"

	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

type ImageData struct {
	Data []byte `json:"data"`
	Ext  string `json:"ext"`
}

type Bom struct {
	MS                    string       `json:"ms" csv:"ms"`
	HenryID               string       `json:"henry_id" csv:"henry_id"`
	ImageLink             string       `json:"image_link" csv:"image_link"`
	Image                 []*ImageData `json:"image"`
	FitSizing             string       `json:"fit_sizing" csv:"fit_sizing"`
	ConfirmedColor        string       `json:"confirmed_color" csv:"confirmed_color"`
	Cat1                  string       `json:"cat1" csv:"cat1"`
	Cat2                  string       `json:"cat2" csv:"cat2"`
	MainFabric            string       `json:"main_fabric" csv:"main_fabric"`
	MainFabricComposition string       `json:"main_fabric_composition" csv:"main_fabric_composition"`
	Buttons               []*ImageData `json:"buttons" csv:"buttons"`
	Rivet                 string       `json:"rivet" csv:"rivet"`
	ZipperTapeColor       string       `json:"zipper_tape_color" csv:"zipper_tape_color"`
	ZipperTeeth           string       `json:"zipper_teeth" csv:"zipper_teeth"`
	Thread                string       `json:"thread" csv:"thread"`
	OtherDetail           string       `json:"other_detail" csv:"other_detail"`
	Artwork               string       `json:"artwork" csv:"artwork"`
	Lining                string       `json:"lining" csv:"lining"`
	MainLabel             []*ImageData `json:"main_label" csv:"main_label"`
	SizeLabel             []*ImageData `json:"size_label" csv:"size_label"`
	WovenPatch            []*ImageData `json:"woven_patch" csv:"woven_patch"`
	MainHangtagString     []*ImageData `json:"main_hangtag_string" csv:"main_hangtag_string"`
	SustainableHangtag    string       `json:"sustainable_hangtag" csv:"sustainable_hangtag"`
	CareLabel             []*ImageData `json:"care_label" csv:"care_label"`
	BarcodeSticker        []*ImageData `json:"barcode_sticker" csv:"barcode_sticker"`
}

func ParseBOM(reader io.Reader) ([]*Bom, error) {
	var sheetName = "Inflow"
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

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var datas []*Bom
	var startingRowIndex = -1

	var msIDIndex = -1
	var henryIDIndex = -1
	var imageLinkIndex = -1
	var imageIndex = -1
	var fitSizingIndex = -1
	var colorIndex = -1
	var cat1Index = -1
	var cat2Index = -1
	var mainFabricIndex = -1
	var mainFabricCompositionIndex = -1
	var buttonsIndex = -1
	var rivetIndex = -1
	var zipperTapeColorIndex = -1
	var zipperTeethIndex = -1
	var threadIndex = -1
	var otherDetailIndex = -1
	var artWorkIndex = -1
	var liningIndex = -1
	var mainLabelIndex = -1
	var sizeLabelIndex = -1
	var wovenPatchIndex = -1
	var mainHangTagIndex = -1
	var sustainableHangtagIndex = -1
	var careLabelIndex = -1
	var barcodeStickerIndex = -1

	var colLens = lo.Map(rows, func(row []string, index int) int {
		return len(row)
	})

	var maxCol = lo.Max(colLens)

	for rowIndex, row := range rows {
		var item Bom
		for colIdx := 0; colIdx < maxCol; colIdx++ {
			var colCell = ""
			if colIdx < len(row) {
				colCell = row[colIdx]
			}

			if colCell == "MS" {
				startingRowIndex = rowIndex
				msIDIndex = colIdx
				continue
			}

			if colCell == "Henry ID" {
				startingRowIndex = rowIndex
				henryIDIndex = colIdx
				continue
			}

			if colCell == "Image Link" {
				startingRowIndex = rowIndex
				imageLinkIndex = colIdx
				continue
			}

			if colCell == "Image" {
				startingRowIndex = rowIndex
				imageIndex = colIdx
				continue
			}

			if colCell == "Fit/ Sizing" {
				startingRowIndex = rowIndex
				fitSizingIndex = colIdx
				continue
			}

			if colCell == "Confirmed Color" {
				startingRowIndex = rowIndex
				colorIndex = colIdx
				continue
			}

			if colCell == "Cat1" {
				startingRowIndex = rowIndex
				cat1Index = colIdx
				continue
			}

			if colCell == "Cat2" {
				startingRowIndex = rowIndex
				cat2Index = colIdx
				continue
			}

			if colCell == "Main Fabric" {
				startingRowIndex = rowIndex
				mainFabricIndex = colIdx
				continue
			}

			if colCell == "Main Fabric Composition" {
				startingRowIndex = rowIndex
				mainFabricCompositionIndex = colIdx
				continue
			}

			if colCell == "Buttons" {
				startingRowIndex = rowIndex
				buttonsIndex = colIdx
				continue
			}

			if colCell == "Rivet" {
				startingRowIndex = rowIndex
				rivetIndex = colIdx
				continue
			}

			if colCell == "Zipper Tape Color" {
				startingRowIndex = rowIndex
				zipperTapeColorIndex = colIdx
				continue
			}

			if colCell == "Zipper Teeth" {
				startingRowIndex = rowIndex
				zipperTeethIndex = colIdx
				continue
			}

			if colCell == "Thread" {
				startingRowIndex = rowIndex
				threadIndex = colIdx
				continue
			}

			if colCell == "Other Detail" {
				startingRowIndex = rowIndex
				otherDetailIndex = colIdx
				continue
			}

			if colCell == "Artwork" {
				startingRowIndex = rowIndex
				artWorkIndex = colIdx
				continue
			}

			if colCell == "Lining" {
				startingRowIndex = rowIndex
				liningIndex = colIdx
				continue
			}

			if colCell == "Main Label" {
				startingRowIndex = rowIndex
				mainLabelIndex = colIdx
				continue
			}

			if colCell == "Size Label" {
				startingRowIndex = rowIndex
				sizeLabelIndex = colIdx
				continue
			}

			if colCell == "Woven Patch" {
				startingRowIndex = rowIndex
				wovenPatchIndex = colIdx
				continue
			}

			if colCell == "Main Hangtag/ String" {
				startingRowIndex = rowIndex
				mainHangTagIndex = colIdx
				continue
			}

			if colCell == "Sustainable Hangtag" {
				startingRowIndex = rowIndex
				sustainableHangtagIndex = colIdx
				continue
			}

			if colCell == "Care Label" {
				startingRowIndex = rowIndex
				careLabelIndex = colIdx
				continue
			}

			if colCell == "Barcode Sticker" {
				startingRowIndex = rowIndex
				barcodeStickerIndex = colIdx
				continue
			}

			if startingRowIndex == -1 {
				continue
			}

			switch colIdx {
			case msIDIndex:
				item.MS = colCell

			case henryIDIndex:
				item.HenryID = colCell

			case imageLinkIndex:
				item.ImageLink = colCell

			case imageIndex:
				cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIndex+1)

				if err == nil {
					pics, _ := f.GetPictures(sheetName, cellName)
					for _, pic := range pics {
						item.Image = append(item.Image, &ImageData{
							Data: pic.File,
							Ext:  pic.Extension,
						})
					}

				}

			case fitSizingIndex:
				item.FitSizing = colCell

			case colorIndex:
				item.ConfirmedColor = colCell

			case cat1Index:
				item.Cat1 = colCell

			case cat2Index:
				item.Cat2 = colCell

			case mainFabricIndex:
				item.MainFabric = colCell

			case mainFabricCompositionIndex:
				item.MainFabricComposition = colCell

			case buttonsIndex:
				cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIndex+1)
				if err == nil {
					pics, _ := f.GetPictures(sheetName, cellName)
					for _, pic := range pics {
						item.Buttons = append(item.Buttons, &ImageData{
							Data: pic.File,
							Ext:  pic.Extension,
						})
					}

				}

			case rivetIndex:
				item.Rivet = colCell

			case zipperTapeColorIndex:
				item.ZipperTapeColor = colCell

			case zipperTeethIndex:
				item.ZipperTeeth = colCell

			case threadIndex:
				item.Thread = colCell

			case otherDetailIndex:
				item.OtherDetail = colCell

			case artWorkIndex:
				item.Artwork = colCell

			case liningIndex:
				item.Lining = colCell

			case mainLabelIndex:
				cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIndex+1)
				if err == nil {
					pics, _ := f.GetPictures(sheetName, cellName)
					for _, pic := range pics {
						item.MainLabel = append(item.MainLabel, &ImageData{
							Data: pic.File,
							Ext:  pic.Extension,
						})
					}
				}
			case sizeLabelIndex:
				cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIndex+1)
				if err == nil {
					pics, _ := f.GetPictures(sheetName, cellName)
					for _, pic := range pics {
						item.SizeLabel = append(item.SizeLabel, &ImageData{
							Data: pic.File,
							Ext:  pic.Extension,
						})
					}
				}
			case wovenPatchIndex:
				cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIndex+1)
				if err == nil {
					pics, _ := f.GetPictures(sheetName, cellName)
					for _, pic := range pics {
						item.WovenPatch = append(item.WovenPatch, &ImageData{
							Data: pic.File,
							Ext:  pic.Extension,
						})
					}

				}

			case mainHangTagIndex:
				cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIndex+1)
				if err == nil {
					pics, _ := f.GetPictures(sheetName, cellName)
					for _, pic := range pics {
						item.MainHangtagString = append(item.MainHangtagString, &ImageData{
							Data: pic.File,
							Ext:  pic.Extension,
						})
					}

				}

			case sustainableHangtagIndex:
				item.SustainableHangtag = colCell

			case careLabelIndex:
				cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIndex+1)
				if err == nil {
					pics, _ := f.GetPictures(sheetName, cellName)
					for _, pic := range pics {
						item.CareLabel = append(item.CareLabel, &ImageData{
							Data: pic.File,
							Ext:  pic.Extension,
						})
					}

				}

			case barcodeStickerIndex:
				cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIndex+1)
				if err == nil {
					pics, _ := f.GetPictures(sheetName, cellName)
					for _, pic := range pics {
						item.BarcodeSticker = append(item.BarcodeSticker, &ImageData{
							Data: pic.File,
							Ext:  pic.Extension,
						})
					}

				}

			}

		}

		if startingRowIndex != -1 && startingRowIndex != rowIndex {
			datas = append(datas, &item)
		}

	}

	return datas, nil
}
