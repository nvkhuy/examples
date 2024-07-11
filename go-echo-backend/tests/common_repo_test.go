package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/engineeringinflow/inflow-backend/pkg/seeder"
	"github.com/engineeringinflow/inflow-backend/services/backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCommon_GetQRCode(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/common/qrcode?content=SG240621IF,SG240621IF110,31C4090611-239-M:11", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNtYjFtMHJiMmhqZjU3anJhZGpnIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6ImNsaWVudCIsImlzcyI6ImNtYjFtMHJiMmhqZjU3anJhZGswIiwic3ViIjoiY2xpZW50In0.KLStPVNn7wx86X4iJfdhYbnkUx-MOppWFo23pZ0jkUM")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestCommon_CreateQRCodes(t *testing.T) {
	var app = initApp("dev")
	app.Config.EFSPath = "./"
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = []byte(`{"contents":["SG240621IF,SG240621IF110,31C4090611-239-M:11"]}`)

	var req = httptest.NewRequest(echo.POST, "/api/v1/common/qrcode", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNtYjFtMHJiMmhqZjU3anJhZGpnIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6ImNsaWVudCIsImlzcyI6ImNtYjFtMHJiMmhqZjU3anJhZGswIiwic3ViIjoiY2xpZW50In0.KLStPVNn7wx86X4iJfdhYbnkUx-MOppWFo23pZ0jkUM")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestCommon_DownloadQRCode(t *testing.T) {
	var app = initApp("dev")
	app.Config.EFSPath = "./"
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var body = []byte(`{"contents":["SG240621IF,SG240621IF110,31C4090611-239-M:11"]}`)

	var req = httptest.NewRequest(echo.GET, "/api/v1/common/qrcode/download?content=SG240621IF,SG240621IF110,31C4090611-239-M:11", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNtYjFtMHJiMmhqZjU3anJhZGpnIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6ImNsaWVudCIsImlzcyI6ImNtYjFtMHJiMmhqZjU3anJhZGswIiwic3ViIjoiY2xpZW50In0.KLStPVNn7wx86X4iJfdhYbnkUx-MOppWFo23pZ0jkUM")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestCommon_Migration(t *testing.T) {
	var app = initApp()
	initMigration(app)
}

func TestCommon_SeedAccounts(t *testing.T) {
	var app = initApp()
	seeder.New(app.DB).SeedAccounts()
}

func TestCommon_GetAttachment(t *testing.T) {
	var app = initApp("prod")
	var s3Client = s3.New(app.Config)
	var resp = repo.NewCommonRepo(app.DB).GetAttachment(s3Client, repo.GetAttachmentParams{
		FileKey:       "uploads/media/cl5lm4tf3m50cehr8770_blogs_cmp7m0jb2hjaaej7j1tg.png",
		ThumbnailSize: "720w",
	})

	fmt.Println(time.Now().Unix())
	fmt.Println(resp)
}

func TestCommon_GetThumbnailAttachment(t *testing.T) {
	var app = initApp("dev")
	var s3Client = s3.New(app.Config)
	var resp = repo.NewCommonRepo(app.DB).GetThumbnailAttachment(s3Client, repo.GetAttachmentParams{
		FileKey: "uploads/media/cg5anr2llkm6ctpvq8k0_fabric_cmqj1kcuqje7m7gi1hgg.mp4",
	})

	fmt.Println(time.Now().Unix())
	fmt.Println(resp)
}

func TestCommon_GetBlurAttachment(t *testing.T) {
	var app = initApp("dev")
	var s3Client = s3.New(app.Config)
	var resp = repo.NewCommonRepo(app.DB).GetBlurAttachment(s3Client, repo.GetAttachmentParams{
		FileKey:       "uploads/media/cmj2e5rb2hj2bhd9qsd0_rfq_attachments_cnbbilrqcpdo5vcmdt7g.png",
		ThumbnailSize: "720w",
	})

	fmt.Println(time.Now().Unix())
	fmt.Println(resp)
}

func TestCommon_GetAttachmentAPI(t *testing.T) {
	var app = initApp("dev")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var req = httptest.NewRequest(echo.GET, "/api/v1/common/attachments/uploads/media/ci83c8djtqd6mfut3fmg_rfq_attachments_cmjr4frb2hjfcfm7in10.png?thumbnail_size=1280w", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNtYjFtMHJiMmhqZjU3anJhZGpnIiwidHoiOiJBc2lhL1NhaWdvbiIsImF1ZCI6ImNsaWVudCIsImlzcyI6ImNtYjFtMHJiMmhqZjU3anJhZGswIiwic3ViIjoiY2xpZW50In0.KLStPVNn7wx86X4iJfdhYbnkUx-MOppWFo23pZ0jkUM")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestCommon_UploadSignatures(t *testing.T) {
	var app = initApp("prod")
	var router = routes.NewRouter(app)
	router.SetupRoutes()

	var data = []byte(`{"records":[{"content_type":"application/pdf","resource":"supplier_onboarding"}]}`)

	var req = httptest.NewRequest(echo.POST, "/api/v1/common/upload/signatures", bytes.NewBuffer(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("signature", "sBcSSP2ccE8YtD6x0umQxVKIV/3Hhgz7jGymvJ7/KG7vaXNMzXpSvbHroqhuqu06J1Pc55L4SK/j9ubkEZEHDWK/DZU7BlRqWf89CdfeemQm1dMTidIRlJ6uGUNKN8PsUc+e3duwfyH5iGLpsEGxqKYyEDCNa5b6rR2019L4DMQ9WzBXTC0IhAHkiYTQxbI8Cdqaif1Lc2+buzJ8Y9COH73+AIaWojbz/DKLWNIqk0LzT5QNGMOQ6S01U6H9A1V5IWpA5555IIQwkyQZDdWm6qrRqIXIDO1R8PDbqWJkqUg/hTzKunSSmFpclib34mZb8asfz9NCtab8i3IZtEX9jA==")

	var rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	helper.PrintJSONBytes(rec.Body.Bytes())

}

func TestCommon_GenerateSitemap(t *testing.T) {
	var app = initApp("local")
	var resp = repo.NewCommonRepo(app.DB).GetShareLink(repo.GetShareLinkParams{
		LinkID: "YnBjL2Nta2RmM2piMmhqZGZ1MGp0Y2hn",
	})
	// https://dev-t.joininflow.io/YnBjL2Nta2RmM2piMmhqZGZ1MGp0Y2hn

	helper.PrintJSON(resp)
}
func TestCommon_CreateShareLink(t *testing.T) {
	var app = initApp("local")
	resp, err := repo.NewCommonRepo(app.DB).CreateShareLink(repo.CreateShareLinkParams{
		ReferenceID: "BPO-IKVD-82786",
		Action:      "bulk_preview_checkout",
	})
	assert.NoError(t, err)

	helper.PrintJSON(resp)
}
