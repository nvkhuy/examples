package models

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/samber/lo"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (inquiry *Inquiry) BeforeCreate(tx *gorm.DB) error {
	if inquiry.ReferenceID == "" {
		var id = helper.GenerateInquiryReferenceID()
		tx.Statement.SetColumn("ReferenceID", id)
		tx.Statement.AddClauseIfNotExists(clause.OnConflict{
			Columns: []clause.Column{{Name: "reference_id"}},
			DoUpdates: clause.Assignments(func() map[string]interface{} {
				var id = helper.GenerateInquiryReferenceID()
				return map[string]interface{}{"reference_id": id}
			}()),
		})
	}

	return nil
}

func (inquiry Inquiry) GetSampleUnitPrice() (unitPrice price.Price) {
	for _, record := range inquiry.AdminQuotations {
		if record.Type == enums.InquiryTypeSample {
			unitPrice = record.Price
		}
	}
	return
}
func (inquiry Inquiry) GetQuotedPrice() (price price.Price) {
	if inquiry.AdminQuotations != nil && inquiry.Quantity != nil {
		sort.Slice(inquiry.AdminQuotations[:], func(i, j int) bool {
			return *inquiry.AdminQuotations[i].Quantity < *inquiry.AdminQuotations[j].Quantity
		})

		for _, record := range inquiry.AdminQuotations {
			if record.Type == enums.InquiryTypeSample {
				price = record.Price
			}

			if record.Type == enums.InquiryTypeBulk && record.Quantity != nil && *inquiry.Quantity > *record.Quantity {
				price = record.Price
			}

		}
	}

	return
}

func (inquiry *Inquiry) GetCustomerIOMetadata(extras map[string]interface{}) map[string]interface{} {
	var cfg = config.GetInstance()
	var result = map[string]interface{}{
		"brand_inquiry_url": fmt.Sprintf("%s/inquiries/%s", cfg.BrandPortalBaseURL, inquiry.ID),
		"admin_inquiry_url": fmt.Sprintf("%s/inquiries/%s/customer/assets", cfg.AdminPortalBaseURL, inquiry.ID),
		"id":                inquiry.ID,
		"reference_id":      inquiry.ReferenceID,
	}
	if inquiry.Currency != "" {
		result["currency"] = inquiry.Currency
		result["currency_customerio_code"] = inquiry.Currency.GetCustomerIOCode()
	}

	if inquiry.BuyerQuotationStatus != "" {
		result["buyer_quotation_status"] = inquiry.BuyerQuotationStatus
	}

	if inquiry.BuyerQuotationStatus != "" {
		result["buyer_quotation_status"] = inquiry.BuyerQuotationStatus
	}

	if inquiry.Status != "" {
		result["status"] = inquiry.Status
	}

	if inquiry.Title != "" {
		result["title"] = inquiry.Title
	}

	if inquiry.ColorList != "" {
		result["color_list"] = inquiry.ColorList
	}

	if inquiry.SizeChart != "" {
		result["size_chart"] = inquiry.SizeChart
	}

	if inquiry.Category != nil {
		result["category"] = *inquiry.Category
	}

	if inquiry.Composition != "" {
		result["composition"] = inquiry.Composition
	}

	if inquiry.Quantity != nil {
		result["quantity"] = *inquiry.Quantity
	}

	if inquiry.FabricName != "" {
		result["fabric_name"] = inquiry.FabricName
	}

	if inquiry.FabricWeight != nil {
		result["fabric_weight"] = *inquiry.FabricWeight
	}

	if inquiry.SkuNote != "" {
		result["sku_note"] = inquiry.SkuNote
	}

	if inquiry.Requirement != "" {
		result["requirement"] = inquiry.Requirement
	}

	if inquiry.Requirements != nil {
		result["requirements"] = *inquiry.Requirements
	}

	if inquiry.ApproveRejectMeta != nil {
		result["approve_reject_meta"] = inquiry.ApproveRejectMeta
	}

	if inquiry.Attachments != nil {
		result["attachments"] = inquiry.Attachments.GenerateFileURL()
	}

	if inquiry.FabricAttachments != nil {
		result["fabric_attachments"] = inquiry.FabricAttachments.GenerateFileURL()
	}

	if inquiry.TechpackAttachments != nil {
		result["techpack_attachments"] = inquiry.TechpackAttachments.GenerateFileURL()
	}

	if inquiry.Design != nil {
		result["design"] = inquiry.Design.GenerateFileURL()
	}

	if inquiry.Document != nil {
		result["document"] = inquiry.Document.GenerateFileURL()
	}

	if inquiry.User != nil {
		result["user"] = inquiry.User.GetCustomerIOMetadata(nil)
	}

	if inquiry.Assignees != nil {
		result["assignees"] = inquiry.Assignees.GetCustomerIOMetadata(nil)
	}

	for k, v := range extras {
		result[k] = v
	}

	return result
}

func (inquiry *Inquiry) GetShortCustomerIOMetadata() map[string]interface{} {
	var cfg = config.GetInstance()
	var result = map[string]interface{}{
		"brand_inquiry_url": fmt.Sprintf("%s/inquiries/%s", cfg.BrandPortalBaseURL, inquiry.ID),
		"admin_inquiry_url": fmt.Sprintf("%s/inquiries/%s/assets", cfg.AdminPortalBaseURL, inquiry.ID),
		"id":                inquiry.ID,
		"reference_id":      inquiry.ReferenceID,
		"created_at":        inquiry.CreatedAt,
	}

	if len(inquiry.Assignees) > 0 {
		result["assignee_names"] = lo.Map(inquiry.Assignees, func(item *User, index int) string {
			return item.Name
		})
	}

	if inquiry.User != nil {
		result["user_name"] = inquiry.User.Name
	}

	return result
}

func (r Inquiries) GetCustomerIOMetadata() (list []map[string]interface{}) {
	for _, record := range r {
		list = append(list, record.GetShortCustomerIOMetadata())
	}

	return
}

func (records Inquiries) ToExcel() ([]byte, error) {
	var data = [][]interface{}{
		{"Reference ID", "User", "Expected Price", "Quantity", "Product", "Status", "Buyer Status", "Assignee", "Posted Date"},
	}
	for _, record := range records {
		data = append(data, []interface{}{
			record.ReferenceID,
			func() string {
				if record.User != nil {
					return record.User.Name
				}
				return ""
			}(),
			func() string {
				if record.ExpectedPrice != nil {
					return record.ExpectedPrice.FormatMoney(record.Currency)
				}
				return ""
			}(),
			values.Int64Value(record.Quantity),
			record.Title,
			record.Status.DisplayName(),
			record.BuyerQuotationStatus.DisplayName(),
			func() interface{} {
				if len(record.Assignees) > 0 {
					var names = lo.Map(record.Assignees, func(item *User, index int) string {
						return item.Name
					})

					return strings.Join(names, ",")
				}
				return nil
			}(),
			time.Unix(record.CreatedAt, 0).In(helper.DefaultTimezone.GetLocation()).Format(`Mon. Jan 2 2006 3:04 PM MST-0700`),
		})
	}

	return helper.ToExcel(data)
}

func (iqs Inquiries) IDs() []string {
	var iqIDs = make([]string, 0, len(iqs))
	for _, iq := range iqs {
		iqIDs = append(iqIDs, iq.ID)
	}
	return iqIDs
}

func (iqs Inquiries) UserIDs() []string {
	userIDs := make([]string, 0, len(iqs))
	for _, iq := range iqs {
		userIDs = append(userIDs, iq.UserID)
	}
	return userIDs
}
