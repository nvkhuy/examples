package repo

import (
	"errors"
	"fmt"
	"sync"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/rotisserie/eris"
	"google.golang.org/api/sheets/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

type SeoTranslationRepo struct {
	db       *db.DB
	logger   *logger.Logger
	sheetAPI *sheets.Service
}

func NewSeoTranslationRepo(db *db.DB) *SeoTranslationRepo {
	return &SeoTranslationRepo{
		db:     db,
		logger: logger.New("repo/SeoTranslation"),
	}
}

func (r *SeoTranslationRepo) WithSheetAPI(api *sheets.Service) *SeoTranslationRepo {
	r.sheetAPI = api
	return r
}

func (r *SeoTranslationRepo) FetchSeoTranslation(params *models.FetchSeoTranslationParams) (updates models.SeoTranslationSlice, err error) {
	if !params.Domain.IsValid() {
		err = errors.New("invalid domain")
		return
	}
	if r.sheetAPI == nil {
		err = errors.New("empty sheet API")
		return
	}
	params = params.Fetch()
	var m sync.Map
	load := func(domain enums.Domain, lang enums.LanguageCode, read string, q chan error) {
		resp, e := r.sheetAPI.Spreadsheets.Values.Get(params.SpreadsheetId, read).Do()
		if e != nil {
			q <- e
			return
		}
		for _, row := range resp.Values {
			update := models.SeoTranslationFromSliceInterface(domain, lang, row)
			if v, ok := m.Load(update.Keyword); ok {
				exits := models.SeoTranslationFromInterface(v)
				if update.EN == "" {
					update.EN = exits.EN
				}
				if update.VI == "" {
					update.VI = exits.VI
				}
			}
			m.Store(update.Keyword, update)
		}
		q <- nil
	}
	q := make(chan error)
	readEN := fmt.Sprintf("%s!%s:%s", params.SheetENName, params.From, params.To)
	readVI := fmt.Sprintf("%s!%s:%s", params.SheetVIName, params.From, params.To)
	go load(params.Domain, enums.LanguageCodeEnglish, readEN, q)
	go load(params.Domain, enums.LanguageCodeVietnam, readVI, q)
	for i := 0; i < 2; i++ {
		if e := <-q; e != nil {
			return nil, eris.Wrap(e, e.Error())
		}
	}
	m.Range(func(_, v any) bool {
		updates = append(updates, models.SeoTranslationFromInterface(v))
		return true
	})
	err = r.db.Clauses(clause.OnConflict{ // Upsert
		Columns:   []clause.Column{{Name: "keyword"}, {Name: "domain"}},
		UpdateAll: true,
	}).Create(&updates).Error
	return
}

func (r *SeoTranslationRepo) GetSEOTranslation(params models.GetSEOTranslationForm) (result map[string]datatypes.JSONMap, err error) {
	if !params.Domain.IsValid() {
		err = errors.New("invalid domain")
		return
	}
	type data struct {
		Lang enums.LanguageCode
		Data *datatypes.JSONMap
		Err  error
	}
	load := func(domain enums.Domain, lang enums.LanguageCode, ch chan data) {
		var _result *datatypes.JSONMap
		q := fmt.Sprintf("SELECT json_object_agg(keyword, %v) FROM seo_translations WHERE domain = '%s'", lang, domain)
		if err = r.db.Raw(q).Scan(&_result).Error; err != nil {
			ch <- data{
				Err: err,
			}
			return
		}
		ch <- data{
			Lang: lang,
			Data: _result,
		}
		return
	}
	ch := make(chan data)
	go load(params.Domain, enums.LanguageCodeEnglish, ch)
	go load(params.Domain, enums.LanguageCodeVietnam, ch)
	result = make(map[string]datatypes.JSONMap)
	for i := 0; i < 2; i++ {
		v := <-ch
		if v.Err != nil {
			err = eris.Wrap(err, err.Error())
			return
		}
		if v.Data == nil {
			continue
		}
		result[v.Lang.String()] = *v.Data
	}
	return
}
