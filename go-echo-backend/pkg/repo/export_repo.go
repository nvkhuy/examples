package repo

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/golang-jwt/jwt"
)

type ExportRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewExportRepo(db *db.DB) *ExportRepo {
	return &ExportRepo{
		db:     db,
		logger: logger.New("repo/export"),
	}
}

type ExportOrdersCSVResponse struct {
	DownloadURL string `json:"download_url"`
}

type FileInfo struct {
	FilePath    string   `json:"-"`
	File        *os.File `json:"-"`
	FileName    string   `json:"-"`
	DownloadURL string   `json:"download_url"`
	RawData     []byte   `json:"-"`
}

func (r *ExportRepo) WriteFile(filePath string, data []byte) (*FileInfo, error) {
	var dir = filepath.Dir(filePath)

	var err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	var flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(filePath, flag, 0666)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(data)

	var result = &FileInfo{
		FilePath: filePath,
		File:     file,
		FileName: path.Base(filePath),
		RawData:  data,
	}
	return result, err
}

func (r *ExportRepo) GenerateDownloadFileURL(fileName string, data []byte) (*FileInfo, error) {
	var filePath = fmt.Sprintf("%s/files/%s", r.db.Configuration.EFSPath, fileName)
	fileInfo, err := r.WriteFile(filePath, data)
	if err != nil {
		return nil, err
	}

	helper.PrintJSON(fileInfo)

	var exp = r.db.Configuration.JWTAssetExpiry
	var claims = &models.AssetCustomClaims{
		FileName: fileName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(exp).Unix(),
		},
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	assetToken, err := token.SignedString([]byte(r.db.Configuration.JWTAssetSecret))
	if err != nil {
		return nil, err
	}

	link, err := url.Parse(fmt.Sprintf("%s/files/%s", r.db.Configuration.ServerBaseURL, fileName))
	if err != nil {
		return nil, err
	}
	var q = link.Query()
	q.Add("token", assetToken)
	link.RawQuery = q.Encode()

	fileInfo.DownloadURL = link.String()
	fileInfo.RawData = data

	return fileInfo, nil
}
