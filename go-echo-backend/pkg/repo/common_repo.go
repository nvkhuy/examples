package repo

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/runner"
	"github.com/engineeringinflow/inflow-backend/pkg/zalo"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/thaitanloi365/go-utils/values"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/rotisserie/eris"
	"github.com/rs/xid"
)

type CommonRepo struct {
	db *db.DB
}

func NewCommonRepo(db *db.DB) *CommonRepo {
	return &CommonRepo{
		db: db,
	}
}

type GetS3SignaturesParams struct {
	models.JwtClaimsInfo

	Forms []*models.S3SignatureForm
}

func (r *CommonRepo) GetS3Signatures(params GetS3SignaturesParams) []s3.Signature {
	var records []s3.Signature
	for _, record := range params.Forms {
		records = append(records, r.GetS3Signature(GetS3SignatureParams{
			JwtClaimsInfo: params.JwtClaimsInfo,
			Form:          record,
		}))
	}

	return records
}

type GetS3SignatureParams struct {
	models.JwtClaimsInfo

	Form *models.S3SignatureForm
}

func (r *CommonRepo) GetS3Signature(params GetS3SignatureParams) s3.Signature {
	var id = params.GetUserID()
	if id == "" {
		id = "anonymous"
	}

	var extension = params.Form.ContentType.GetExtension()
	var key = fmt.Sprintf("uploads/user/%s/%s/%s%s", id, params.Form.Resource, xid.New().String(), extension)
	if params.Form.ContentType.IsMedia() {
		key = fmt.Sprintf("uploads/media/%s_%s_%s%s", id, params.Form.Resource, xid.New().String(), extension)
	}

	var maxSize = 100 * 1024 * 1024 // 100MB
	if params.GetRole().IsAdmin() {
		maxSize = 10000 * 1024 * 1024 // 10000MB
	}

	var result = s3.New(r.db.Configuration).GenerateSignature(s3.GenerateSignatureParams{
		ACL:           "private",
		Bucket:        r.db.Configuration.AWSS3StorageBucket,
		ContentType:   string(params.Form.ContentType),
		ExpiryMinutes: 30,
		Key:           key,
		MaxFileSize:   maxSize,
	})

	return result
}

func (r *CommonRepo) CheckExists(form models.CheckExistsForm) (*models.CheckExistsResponse, error) {
	var resp models.CheckExistsResponse
	var user = models.User{
		Email: form.Email,
	}

	if form.Email != "" {
		var err = r.db.Select("ID").Unscoped().First(&models.User{}, "email = ?", user.Email).Error
		if err != nil {
			if r.db.IsRecordNotFoundError(err) {

			} else {
				return nil, err
			}
		} else {
			resp.IsExists = true
		}

	} else {
		return nil, errs.ErrPhoneOrEmailOrUserNameRequired
	}

	return &resp, nil
}

type GetAttachmentParams struct {
	models.JwtClaimsInfo

	ThumbnailSize string `json:"thumbnail_size" query:"thumbnail_size" form:"thumbnail_size"`
	FileKey       string `json:"file_key" query:"file_key" form:"file_key" param:"file_key"`
	NoCache       bool   `json:"no_cache" query:"no_cache" form:"no_cache" param:"no_cache"`
}

func (r *CommonRepo) GetAttachment(s3Client *s3.Client, params GetAttachmentParams) string {
	var key = strings.TrimPrefix(params.FileKey, "/")
	var storageURL = fmt.Sprintf("https://%s/%s", r.db.Configuration.StorageURL, key)

	if params.ThumbnailSize != "" && helper.IsImageExt(params.FileKey) {
		var keyWithSize = fmt.Sprintf("%s/%s", params.ThumbnailSize, key)
		if _, err := s3Client.CheckFile(r.db.Configuration.AWSS3CdnBucket, keyWithSize); err == nil && !params.NoCache {
			return fmt.Sprintf("https://%s/%s", r.db.Configuration.CDNURL, keyWithSize)
		}

		token, err := r.db.Configuration.GetMediaToken(keyWithSize)
		if err != nil {
			return storageURL
		}
		if params.NoCache {
			return fmt.Sprintf("%s?token=%s&size=%s&key=%s&no_cache=%t", r.db.Configuration.LambdaAPIResizeURL, token, params.ThumbnailSize, params.FileKey, params.NoCache)
		}
		return fmt.Sprintf("%s?token=%s&size=%s&key=%s", r.db.Configuration.LambdaAPIResizeURL, token, params.ThumbnailSize, params.FileKey)

	}

	if params.ThumbnailSize != "" && helper.IsVideoExt(params.FileKey) {
		return r.GetThumbnailAttachment(s3Client, params)
	}

	return storageURL

}

func (r *CommonRepo) GetThumbnailAttachment(s3Client *s3.Client, params GetAttachmentParams) string {
	var key = strings.TrimPrefix(params.FileKey, "/")
	var storageURL = fmt.Sprintf("https://%s/%s", r.db.Configuration.StorageURL, key)

	if helper.IsVideoExt(params.FileKey) {
		token, err := r.db.Configuration.GetMediaToken(key)
		if err != nil {
			return storageURL
		}

		if params.NoCache {
			return fmt.Sprintf("%s?token=%s&file_key=%s&no_cache=true&type=thumbnail", r.db.Configuration.LambdaAPIFfmpegURL, token, params.FileKey)
		}

		return fmt.Sprintf("%s?token=%s&file_key=%s&type=thumbnail", r.db.Configuration.LambdaAPIFfmpegURL, token, params.FileKey)

	}

	if params.ThumbnailSize != "" && helper.IsImageExt(params.FileKey) {
		return r.GetAttachment(s3Client, params)
	}

	return storageURL

}

func (r *CommonRepo) GetBlurAttachment(s3Client *s3.Client, params GetAttachmentParams) string {
	var key = strings.TrimPrefix(params.FileKey, "/")
	var storageURL = fmt.Sprintf("https://%s/%s", r.db.Configuration.StorageURL, key)

	if params.ThumbnailSize != "" && helper.IsImageExt(params.FileKey) {
		var keyWithSize = fmt.Sprintf("blur/%s/%s", params.ThumbnailSize, params.FileKey)
		var webpKey = strings.ReplaceAll(params.FileKey, path.Ext(params.FileKey), ".webp")
		var webpKeyWithSize = fmt.Sprintf("blur/%s/%s", params.ThumbnailSize, webpKey)
		if _, err := s3Client.CheckFile(r.db.Configuration.AWSS3CdnBucket, webpKeyWithSize); err == nil && !params.NoCache {
			return fmt.Sprintf("https://%s/%s", r.db.Configuration.CDNURL, webpKeyWithSize)
		}

		token, err := r.db.Configuration.GetMediaToken(keyWithSize)
		if err != nil {
			return storageURL
		}
		if params.NoCache {
			return fmt.Sprintf("%s?token=%s&size=%s&key=%s&no_cache=%t", r.db.Configuration.LambdaAPIBlurURL, token, params.ThumbnailSize, params.FileKey, params.NoCache)
		}
		return fmt.Sprintf("%s?token=%s&size=%s&key=%s", r.db.Configuration.LambdaAPIBlurURL, token, params.ThumbnailSize, params.FileKey)

	}

	if params.ThumbnailSize != "" && helper.IsVideoExt(params.FileKey) {
		return r.GetThumbnailAttachment(s3Client, params)
	}

	return storageURL

}

type GetDownloadLinkParams struct {
	FileKey string `json:"file_key" query:"file_key" form:"file_key" param:"file_key"`
}

func (r *CommonRepo) GetDownloadLink(params GetDownloadLinkParams) string {
	var key = strings.TrimPrefix(params.FileKey, "/")

	return fmt.Sprintf("https://%s/%s", r.db.Configuration.StorageURL, key)

}

func (r *CommonRepo) GenerateSitemap() models.SitemapResponse {
	var resp models.SitemapResponse

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var productRoutes []string
		r.db.Raw("SELECT '/product/' || slug FROM products").Find(&productRoutes)
		if len(productRoutes) > 0 {
			resp.Data = append(resp.Data, productRoutes...)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var blogRoutes []string
		r.db.Raw("SELECT '/blogs/' || slug FROM posts").Find(&blogRoutes)
		if len(blogRoutes) > 0 {
			resp.Data = append(resp.Data, blogRoutes...)
		}
	}()

	wg.Wait()
	return resp

}

type GetShareLinkParams struct {
	LinkID string `json:"link_id" param:"link_id" query:"link_id" validate:"required"`
}

func (r *CommonRepo) GetShareLink(params GetShareLinkParams) string {
	rawText, err := base64.StdEncoding.DecodeString(params.LinkID)
	if err != nil {
		return r.db.Configuration.BrandPortalBaseURL
	}

	var uri = string(rawText)
	var parts = strings.Split(uri, "/")
	if len(parts) != 2 {
		return r.db.Configuration.BrandPortalBaseURL
	}
	var prefix = parts[0]
	var id = parts[1]

	switch prefix {
	case "r": // RFQ
		return fmt.Sprintf("%s/inquiries/%s", r.db.Configuration.BrandPortalBaseURL, id)
	case "s": // Sample
		return fmt.Sprintf("%s/samples/%s", r.db.Configuration.BrandPortalBaseURL, id)
	case "b": // Bulk
		return fmt.Sprintf("%s/bulks/%s", r.db.Configuration.BrandPortalBaseURL, id)
	case "bpc": // Bulk Preview Checkout
		return fmt.Sprintf("%s/bulks/%s/checkout_info", r.db.Configuration.BrandPortalBaseURL, id)
	case "fc": // Fabric Collection
		return fmt.Sprintf("%s/fabric_collections/%s", r.db.Configuration.WebAppBaseURL, id)
	case "fb": // Fabrics
		return fmt.Sprintf("%s/fabrics/%s", r.db.Configuration.WebAppBaseURL, id)
	}

	return r.db.Configuration.BrandPortalBaseURL
}

type CreateShareLinkParams struct {
	models.JwtClaimsInfo
	ReferenceID string                `json:"reference_id" param:"reference_id" query:"reference_id" validate:"required"`
	Action      enums.ShareLinkAction `json:"action" param:"action" query:"action" form:"action"`
}

type CreateShareLinkResponse struct {
	Link string `json:"link"`
}

func (r *CommonRepo) CreateShareLink(params CreateShareLinkParams) (*CreateShareLinkResponse, error) {
	var linkID = ""
	if strings.HasPrefix(params.ReferenceID, "IQ") {
		var record models.Inquiry
		var err = r.db.Select("ID").First(&record, "reference_id = ?", params.ReferenceID).Error
		if err != nil {
			return nil, err
		}
		linkID = fmt.Sprintf("r/%s", record.ID)
	} else if strings.HasPrefix(params.ReferenceID, "PO") {
		var record models.PurchaseOrder
		var err = r.db.Select("ID").First(&record, "reference_id = ?", params.ReferenceID).Error
		if err != nil {
			return nil, err
		}
		linkID = fmt.Sprintf("s/%s", record.ID)
	} else if strings.HasPrefix(params.ReferenceID, "BPO") {
		var record models.BulkPurchaseOrder
		var err = r.db.Select("ID").First(&record, "reference_id = ?", params.ReferenceID).Error
		if err != nil {
			return nil, err
		}

		switch params.Action {
		case enums.ShareLinkActionBulkPreviewCheckout:
			linkID = fmt.Sprintf("bpc/%s", record.ID)
		default:
			linkID = fmt.Sprintf("b/%s", record.ID)
		}
	} else if strings.HasPrefix(params.ReferenceID, "FC") {
		var record models.FabricCollection
		var err = r.db.Select("ID").First(&record, "reference_id = ?", params.ReferenceID).Error
		if err != nil {
			return nil, err
		}
		linkID = fmt.Sprintf("fc/%s", record.ID)
	} else if strings.HasPrefix(params.ReferenceID, "FB") {
		var record models.Fabric
		var err = r.db.Select("ID").First(&record, "reference_id = ?", params.ReferenceID).Error
		if err != nil {
			return nil, err
		}
		linkID = fmt.Sprintf("fb/%s", record.ID)
	} else {
		return nil, eris.New("Invalid reference id")
	}

	var encodedLinkID = base64.StdEncoding.EncodeToString([]byte(linkID))

	var resp = CreateShareLinkResponse{
		Link: fmt.Sprintf("%s/%s", r.db.Configuration.ShareBaseURL, encodedLinkID),
	}
	return &resp, nil
}

type GetDocsParams struct {
	models.JwtClaimsInfo

	Types []enums.SettingDoc `json:"types" query:"types" param:"types"`
}

func (r *CommonRepo) GetDocs(params GetDocsParams) (records []*models.SettingDoc) {
	query.New(r.db, queryfunc.NewSettingDocBuilder(queryfunc.SettingDocBuilderOptions{})).
		WhereFunc(func(builder *query.Builder) {
			if len(params.Types) > 0 {
				builder.Where("sd.type IN ?", params.Types)
			}
		}).
		FindFunc(&records)

	return
}

type SendZNSParams struct {
	Phone        string      `json:"phone"`
	Mode         string      `json:"mode"`
	TemplateID   string      `json:"template_id"`
	TemplateData interface{} `json:"template_data"`
	TrackingID   string      `json:"tracking_id"`
}

func (r *CommonRepo) SendZNS(params SendZNSParams) (result zalo.SendZNSResponse, err error) {
	var sends zalo.SendZNSParams
	if err = copier.Copy(&sends, &params); err != nil {
		return zalo.SendZNSResponse{}, eris.Wrap(err, err.Error())
	}
	result, err = zalo.SendZNS(r.db, sends)
	return
}

type GetCheckoutInfoParams struct {
	LinkID string `json:"link_id" param:"link_id" query:"link_id" validate:"required"`
	Token  string `json:"token" query:"token" validate:"required"`
}

type GetCheckoutInfoResponse struct {
	ReferenceID    string       `json:"reference_id"`
	ShippingFee    *price.Price `json:"shipping_fee"`
	TransactionFee *price.Price `json:"transaction_fee"`
	SubTotal       *price.Price `json:"sub_total"`
	TaxAmount      *price.Price `json:"tax_amount"`
	TotalAmount    *price.Price `json:"total_amount"`
	Description    string       `json:"description"`
	Address        string       `json:"address"`

	Token     string `json:"token"`
	ExpiredAt *int64 `json:"expired_at"`
}

func (r *CommonRepo) GetCheckoutInfo(params GetCheckoutInfoParams) (*GetCheckoutInfoResponse, error) {
	rawText, err := base64.StdEncoding.DecodeString(params.LinkID)
	if err != nil {
		return nil, errs.ErrInvalidCheckoutLink
	}

	var uri = string(rawText)
	var parts = strings.Split(uri, "/")
	if len(parts) != 2 {
		return nil, errs.ErrInvalidCheckoutLink
	}
	var prefix = parts[0]
	var id = parts[1]
	var result GetCheckoutInfoResponse

	switch prefix {
	case "b1": // bulk first payment
		bulk, err := NewBulkPurchaseOrderRepo(r.db).BulkPurchaseOrderPreviewCheckout(BulkPurchaseOrderPreviewCheckoutParams{
			BulkPurchaseOrderID: id,
			PaymentType:         enums.PaymentTypeBankTransfer,
			Milestone:           enums.PaymentMilestoneFirstPayment,
		})
		if err != nil {
			return nil, err
		}
		result.ReferenceID = bulk.ReferenceID
		result.ShippingFee = bulk.ShippingFee
		result.TransactionFee = bulk.FirstPaymentTransactionFee
		result.SubTotal = bulk.FirstPaymentSubTotal
		result.TaxAmount = bulk.FirstPaymentTax
		result.TotalAmount = bulk.FirstPaymentTotal
		result.Description = fmt.Sprintf("Changes for 1st payment (%s) %s", fmt.Sprintf("%0.f", values.Float64Value(bulk.FirstPaymentPercentage))+"%", bulk.ReferenceID)
		if bulk.ShippingAddress != nil && bulk.ShippingAddress.Coordinate != nil {
			result.Address = bulk.ShippingAddress.Coordinate.FormattedAddress
		}
		return &result, err

	case "b": // bulk final payment
		bulk, err := NewBulkPurchaseOrderRepo(r.db).BulkPurchaseOrderPreviewCheckout(BulkPurchaseOrderPreviewCheckoutParams{
			BulkPurchaseOrderID: id,
			PaymentType:         enums.PaymentTypeBankTransfer,
			Milestone:           enums.PaymentMilestoneFinalPayment,
		})
		if err != nil {
			return nil, err
		}

		result.ReferenceID = bulk.ReferenceID
		result.ShippingFee = bulk.ShippingFee
		result.TransactionFee = bulk.FinalPaymentTransactionFee
		result.SubTotal = bulk.FinalPaymentSubTotal
		result.TaxAmount = bulk.FinalPaymentTax
		result.TotalAmount = bulk.FinalPaymentTotal
		result.Description = fmt.Sprintf("Changes for final payment %s", bulk.ReferenceID)
		if bulk.ShippingAddress != nil && bulk.ShippingAddress.Coordinate != nil {
			result.Address = bulk.ShippingAddress.Coordinate.FormattedAddress
		}
		return &result, err

	case "i": //inquiry
		purchaseOrder, err := NewInquiryRepo(r.db).InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
			InquiryID:   id,
			PaymentType: enums.PaymentTypeBankTransfer,
		})
		if err != nil {
			return nil, err
		}

		result.ReferenceID = purchaseOrder.ReferenceID
		result.ShippingFee = purchaseOrder.ShippingFee
		result.TransactionFee = purchaseOrder.TransactionFee
		result.SubTotal = purchaseOrder.SubTotal
		result.TaxAmount = purchaseOrder.Tax
		result.TotalAmount = purchaseOrder.TotalPrice
		result.Description = fmt.Sprintf("Changes for %s - %s", purchaseOrder.ReferenceID, purchaseOrder.Inquiry.Title)
		if purchaseOrder.ShippingAddress != nil && purchaseOrder.ShippingAddress.Coordinate != nil {
			result.Address = purchaseOrder.ShippingAddress.Coordinate.FormattedAddress
		}

		return &result, err
	}

	return nil, errs.ErrInvalidCheckoutLink
}

func (r *CommonRepo) CreateCheckoutShareLink(params CreateShareLinkParams) (*CreateShareLinkResponse, error) {
	var linkID = ""
	if strings.HasPrefix(params.ReferenceID, "IQ") {
		var record models.Inquiry
		var err = r.db.Select("ID").First(&record, "reference_id = ?", params.ReferenceID).Error
		if err != nil {
			return nil, err
		}
		linkID = fmt.Sprintf("i/%s", record.ID)
	} else if strings.HasPrefix(params.ReferenceID, "BPO") {
		var record models.BulkPurchaseOrder
		var err = r.db.Select("ID", "TrackingStatus").First(&record, "reference_id = ?", params.ReferenceID).Error
		if err != nil {
			return nil, err
		}

		switch record.TrackingStatus {
		case enums.BulkPoTrackingStatusFirstPayment:
			linkID = fmt.Sprintf("b1/%s", record.ID)
		case enums.BulkPoTrackingStatusFinalPayment:
			linkID = fmt.Sprintf("b/%s", record.ID)
		default:
			return r.CreateShareLink(params)
		}
	} else {
		return nil, eris.New("Invalid reference id")
	}

	var encodedLinkID = base64.StdEncoding.EncodeToString([]byte(linkID))
	var jwt = sjwt.New()
	jwt.SetExpiresAt(time.Now().Add(time.Hour))
	var token = jwt.Generate([]byte(r.db.Configuration.CheckoutJwtSecret))

	var resp = CreateShareLinkResponse{
		Link: fmt.Sprintf("%s/checkout-info/%s?token=%s", r.db.Configuration.BrandPortalBaseURL, encodedLinkID, token),
	}
	return &resp, nil
}

type SubscribeUpdatesParams struct {
	Email          string `json:"email" validate:"email,required"`
	TechProduct    bool   `json:"tech_product"`
	Sustainability bool   `json:"sustainability"`
	FashionTrends  bool   `json:"fashion_trends"`
	ProductionTips bool   `json:"production_tips"`
	SupplyChain    bool   `json:"supply_chain"`
	Others         bool   `json:"others"`
}

func (r *CommonRepo) SubscribeUpdates(params SubscribeUpdatesParams) error {
	return customerio.GetInstance().Track.Track(params.Email, string(customerio.EventNewSubscriber), map[string]interface{}{
		"tech_product":    params.TechProduct,
		"sustainability":  params.Sustainability,
		"fashion_trends":  params.FashionTrends,
		"production_tips": params.ProductionTips,
		"supply_chain":    params.SupplyChain,
		"others":          params.Others,
	})
}

type GenerateQRCodeParams struct {
	Content string `json:"content" param:"content" query:"content" validate:"required"`
}

func (r *CommonRepo) GenerateQRCode(params GenerateQRCodeParams) (*bytes.Buffer, error) {
	buf, err := helper.GenerateQRCode(helper.GenerateQRCodeOptions{
		Content: params.Content,
	})
	return buf, err

}

type CreateQRCodesParams struct {
	Contents   []string `json:"contents" param:"contents" query:"contents" validate:"required"`
	FolderName string   `json:"folder_name" param:"folder_name" query:"folder_name"`
}

func (r *CommonRepo) CreateQRCodes(params CreateQRCodesParams) (string, error) {
	var folder = fmt.Sprintf("%s_qrcodes.zip", time.Now().Format("2006-01-02_15-04-05"))
	if params.FolderName != "" {
		folder = params.FolderName
	}

	if !strings.HasSuffix(folder, ".zip") {
		folder = fmt.Sprintf("%s.zip", folder)
	}

	var filePath = fmt.Sprintf("%s/files/%s", r.db.Configuration.EFSPath, folder)
	fmt.Println("filePath", filePath)
	var dir = filepath.Dir(filePath)

	var err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	var flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	archive, err := os.OpenFile(filePath, flag, 0666)
	if err != nil {
		return "", err
	}

	defer archive.Close()
	zipWriter := zip.NewWriter(archive)

	var runner = runner.New(10)
	defer runner.Release()

	for _, content := range params.Contents {
		content := content
		runner.Submit(func() {
			buf, err := helper.GenerateQRCode(helper.GenerateQRCodeOptions{
				Content: content,
			})
			if err != nil {
				return
			}

			w, err := zipWriter.Create(fmt.Sprintf("%s.png", content))
			if err != nil {
				return
			}

			if _, err := io.Copy(w, buf); err != nil {
				return
			}
		})
	}

	runner.Wait()
	zipWriter.Close()

	var exp = r.db.Configuration.JWTAssetExpiry
	var claims = &models.AssetCustomClaims{
		FileName: folder,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(exp).Unix(),
		},
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	assetToken, err := token.SignedString([]byte(r.db.Configuration.JWTAssetSecret))
	if err != nil {
		return "", err
	}

	link, err := url.Parse(fmt.Sprintf("%s/files/%s", r.db.Configuration.ConsumerServerBaseURL, folder))
	if err != nil {
		return "nil", err
	}
	var q = link.Query()
	q.Add("token", assetToken)
	link.RawQuery = q.Encode()

	return link.String(), err

}
