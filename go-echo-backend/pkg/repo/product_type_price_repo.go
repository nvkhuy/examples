package repo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/jinzhu/copier"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm/clause"
)

type ProductTypesPriceRepo struct {
	db       *db.DB
	logger   *logger.Logger
	sheetAPI *sheets.Service
}

type imageData struct {
	CurrentValueObj struct {
		Type     string `json:"type"`
		Display  string `json:"display"`
		ImageURL string `json:"imageUrl"`
	} `json:"currentValueObj"`
}

func NewProductTypesPriceRepo(db *db.DB) *ProductTypesPriceRepo {
	return &ProductTypesPriceRepo{
		db:     db,
		logger: logger.New("repo/ProductTypesPrice"),
	}
}

func (r *ProductTypesPriceRepo) WithSheetAPI(api *sheets.Service) *ProductTypesPriceRepo {
	r.sheetAPI = api
	return r
}

func (r *ProductTypesPriceRepo) FetchProductTypesPrice(params *models.FetchProductTypesPriceParams) (result models.ProductTypePriceSlice, err error) {
	if r.sheetAPI == nil {
		err = errors.New("empty sheet API")
		return
	}
	params = params.Fetch()
	readRange := fmt.Sprintf("%s!%s:%s", params.SheetName, params.From, params.To)
	resp, err := r.sheetAPI.Spreadsheets.Values.Get(params.SpreadsheetId, readRange).Do()
	if err != nil {
		return
	}
	for id, row := range resp.Values {
		pp := models.ProductTypePriceFromSlice(id+params.FromNum, row)
		result = append(result, pp)
	}
	if len(result) == 0 {
		return
	}
	err = r.db.Clauses(clause.OnConflict{ // Upsert
		Columns:   []clause.Column{{Name: "row_id"}},
		DoUpdates: models.UpdateProductTypePriceColumns,
	}).Model(&result).Create(&result).Error
	return
}

func (r *ProductTypesPriceRepo) PatchSheetImageURL(params *models.PatchSheetImageURLParams) (err error) {
	var m sync.Map
	var wg sync.WaitGroup
	params = params.Fetch()

	workQueue := make(chan int, params.Concurrency)

	worker := func(rowId int) {
		defer wg.Done()
		webUrl, _ := r.GetSheetImage(rowId, 1, params.Token, params.Cookie)
		catalogURL, _ := r.GetSheetImage(rowId, 2, params.Token, params.Cookie)
		m.Store(rowId, models.ProductTypePrice{
			RowId:         rowId,
			PICWebURL:     webUrl,
			PICCatalogURL: catalogURL,
		})
		<-workQueue
	}

	for i := params.From; i <= params.To; i++ {
		workQueue <- 1
		wg.Add(1)
		go worker(i)
	}
	wg.Wait()

	var prices models.ProductTypePriceSlice
	m.Range(func(k, v interface{}) bool {
		prices = append(prices, models.ProductTypesPriceFromInterface(v))
		return true
	})
	if len(prices) == 0 {
		return
	}
	err = r.db.Clauses(clause.OnConflict{ // Upsert
		Columns:   []clause.Column{{Name: "row_id"}},
		DoUpdates: models.UpdateImagesProductTypePriceColumns,
	}).Model(&prices).Create(&prices).Error
	return err
}

func (r *ProductTypesPriceRepo) GetSheetImage(row, col int, token, cookie string) (imgURL string, err error) {
	if token == "" || cookie == "" {
		err = errors.New("empty token")
		return
	}
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/1YHyMOTHwp5W7AO9C4TlX3xyenvFTwrKww6fL6crB6ak/blame?token=%s", token)
	method := "POST"

	payload := &bytes.Buffer{}
	where := fmt.Sprintf("[30710966,\"[[[\\\"630508040\\\",%d,%d],[[\\\"630508040\\\",%d,%d,%d,%d]]]]\"]", row-1, col-1, row-1, row, col-1, col)
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("selection", where)
	_ = writer.WriteField("clientRevision", "900000000")
	_ = writer.WriteField("includeDiffs", "true")
	err = writer.Close()
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return
	}
	req.Header.Add("authority", "docs.google.com")
	req.Header.Add("accept", "*/*")
	req.Header.Add("cookie", cookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	body = body[len(")]}'"):]
	var img imageData
	if err = json.Unmarshal(body, &img); err != nil {
		return
	}
	imgURL, err = r.SaveImageToS3(row, col, img.CurrentValueObj.ImageURL)
	return
}

func (r *ProductTypesPriceRepo) SaveImageToS3(row, col int, url string) (imgURL string, err error) {
	var s3Client = s3.New(r.db.Configuration)
	if url == "" {
		return
	}

	// download image
	var response *http.Response
	response, err = http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return
	}
	var imgData []byte
	imgData, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	// save s3
	contentType := models.ContentTypeImageJPG
	imgURL = fmt.Sprintf("sheet_images/quote_cost_system-%d-%d%s", row, col, contentType.GetExtension())
	_, _ = s3Client.UploadFile(s3.UploadFileParams{
		Data:        bytes.NewReader(imgData),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(contentType),
		ACL:         "private",
		Key:         imgURL,
	})
	log.Printf("inserted %s", imgURL)
	return
}

func (r *ProductTypesPriceRepo) PaginateProductTypesPrice(params *models.PaginateProductTypesPriceParams) (result *query.Pagination) {
	var builder = queryfunc.NewProductTypesPriceBuilder(queryfunc.ProductTypesPriceBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			queryfunc.ProductTypesPriceBuilderWhereFunc(builder, params)
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()
	return
}

func (r *ProductTypesPriceRepo) PaginateProductTypesPriceVine(params *models.PaginateProductTypesPriceParams) (result *query.Pagination) {
	var builder = queryfunc.NewProductTypesPriceVineBuilder(queryfunc.ProductTypesPriceBuilderOptions{
		VineSlice: params.ToVineSlice(),
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			queryfunc.ProductTypesPriceBuilderWhereFunc(builder, params)
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()
	return
}

func (r *ProductTypesPriceRepo) PaginateProductTypesPriceQuote(params *models.PaginateProductTypesPriceQuoteParams) (result *query.Pagination) {
	var cmSLice models.ProductTypePriceSlice
	var rwdSlice models.RWDFabricPriceSlice
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		queryCM := r.PaginateProductTypesPrice(&models.PaginateProductTypesPriceParams{
			Product:     params.Product,
			Gender:      params.Gender,
			FabricType:  params.FabricType,
			Feature:     params.Feature,
			Category:    params.Category,
			Item:        params.Item,
			Form:        params.Form,
			Description: params.Description,
		})
		if queryCM != nil && queryCM.Records != nil {
			cmSLice = models.ProductTypesPriceSliceFromInterface(queryCM.Records)
		}
	}()

	go func() {
		defer wg.Done()
		queryRWD := NewRWDFabricPriceRepo(r.db).Paginate(&models.PaginateRWDFabricPriceParams{
			FabricType:  params.FabricType,
			Material:    params.Material,
			Composition: params.Composition,
			Weight:      params.Weight,
			CutWidth:    params.CutWidth,
		})
		if queryRWD != nil && queryRWD.Records != nil {
			rwdSlice = models.RWDFabricPriceSliceFromInterface(queryRWD.Records)
		}
	}()

	wg.Wait()

	var prices models.ProductTypePriceSlice
	for _, cm := range cmSLice {
		for _, rwd := range rwdSlice {
			var p models.ProductTypePrice
			if err := copier.Copy(&p, &cm); err == nil {
				p.ID = fmt.Sprintf("%s-%s", p.ID, rwd.ID)
				p.CutWidth = rwd.CutWidth
				materialType := strings.ToLower(strings.TrimSpace(p.FabricType))
				switch enums.RWDMaterial(materialType) {
				case enums.KnitMaterial:
					p.KnitMaterial = rwd.Material
					p.KnitComposition = rwd.Composition
					p.KnitWeight = rwd.Weight
				case enums.WovenMaterial:
					p.WovenMaterial = rwd.Material
					p.WovenComposition = rwd.Composition
					p.WovenWeight = rwd.Weight
				}

				// Fabric Price
				p.FabricPrice0To50 = cm.FabricConsumption * rwd.Price0To50
				p.FabricPrice50To100 = cm.FabricConsumption * rwd.Price50To100
				p.FabricPrice100To500 = cm.FabricConsumption * rwd.Price100To500
				p.FabricPriceAbove500 = cm.FabricConsumption * rwd.PriceAbove500

				// Total Price
				p.TotalPrice0To50 = p.FabricPrice0To50 + cm.CMPrice0To50
				p.TotalPrice50To100 = p.FabricPrice50To100 + cm.CMPrice50To100
				p.TotalPrice100To500 = p.FabricPrice100To500 + cm.CMPrice100To500
				p.TotalPriceAbove500 = p.FabricPriceAbove500 + cm.CMPriceAbove500

				// Description
				p.FabricDescription = rwd.Description

				prices = append(prices, &p)

			}
		}
	}
	result = &query.Pagination{
		PerPage:            0,
		Page:               1,
		PrevPage:           0,
		Offset:             0,
		Records:            prices,
		TotalRecord:        len(prices),
		TotalCurrentRecord: len(prices),
	}
	return
}
