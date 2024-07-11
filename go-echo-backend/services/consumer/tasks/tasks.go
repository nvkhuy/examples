package tasks

import (
	"github.com/engineeringinflow/inflow-backend/pkg/worker"
)

var workerInstance *worker.Worker

// Register register tasks
func Register(w *worker.Worker, IsConsumer bool) {
	workerInstance = w

	if IsConsumer {
		workerInstance.CreateTaskHandler(
			UserPingTask{},

			TrackActivityTask{},
			SendMailTask{},
			SendChatMessageTask{},
			SendWSEventTask{},
			SeenChatRoomTask{},
			ChatTypingTask{},
			CancelChatTypingTask{},
			CreateChatRoomTask{},

			ShopifyCreateWebhooksTask{},
			ShopifyUpdateProductTask{},
			// ShopifySyncProductTask{},
			// ShopifySyncProductsTask{},
			// ShopifySyncChannelProductsTask{},

			HubspotCreateContactTask{},
			HubspotCreateDealTask{},
			HubspotUpdateDealTask{},
			HubspotSyncPOTask{},
			HubspotSyncBulkTask{},
			HubspotSyncInquiryTask{},

			SyncCustomerIOUserTask{},
			DeleteCustomerIOUserTask{},
			TrackCustomerIOTask{},
			CreateInquiryAuditTask{},
			CreateCmsNotificationTask{},
			CreateUserNotificationTask{},
			OnboardUserTask{},
			OnboardSellerTask{},
			ApproveUserTask{},
			SendInquiryToBuyerTask{},
			AssignInquiryPICTask{},
			AssignPurchaseOrderPICTask{},
			AssignBulkPurchaseOrderPICTask{},
			InquiryRemindAdminTask{},
			AddCustomerIOUserDeviceTask{},

			PurchaseOrderBankTransferConfirmedTask{},
			PurchaseOrderBankTransferRejectedTask{},
			BulkPurchaseOrderBankTransferConfirmedTask{},
			BulkPurchaseOrderBankTransferRejectedTask{},
			PoDesignNewCommentTask{},
			GeneratePDFTask{},
			// CreatePOPaymentInvoiceTask{},
			CreateBulkPoDepositPaymentInvoiceTask{},
			// CreateBulkPoFirstPaymentInvoiceTask{},
			// CreateBulkPoSecondPaymentInvoiceTask{},
			// CreateBulkPoFinalPaymentInvoiceTask{},
			CreatePaymentInvoiceTask{},
			NotifyAdminConfirmPaymentTask{},

			CreateBulkPoAttachmentPDFsTask{},
			NewInquiryNotesTask{},
			NewBulkPONotesTask{},
			NewPONotesTask{},
			CreatePOPaymentInvoiceForMultipleItemsTask{},
			CreateSysNotificationTask{},
			RefreshTokenZaloTask{},
			TimeAndActionNotificationTask{},
			TimeAndActionSchedulerTask{},
			PurchaseOrderDesignApproveTask{},
			BulkPurchaseQCApproveTask{},
			BulkPurchaseOrderRawMaterialApproveTask{},
			PurchaseOrderRawMaterialApproveTask{},
			UpdateUserProductClassesTask{},

			SellerApprovePOTask{},
			SellerRejectPOTask{},

			AdminApproveSellerQuotationTask{},
			AdminRejectSellerQuotationTask{},
			AdminApproveSellerBulkPurchaseOrderQuotation{},
			AdminRejectSellerBulkPurchaseOrderQuotation{},
			RemindUnseenMessageTask{},

			GenerateBlurTask{},
			UploadProductFileTask{},
		)

		_, _ = w.ScheduleTask(
			"CRON_TZ=Asia/Saigon 0 08 * * *",
			InquiryRemindAdminTask{},
		)

		_, _ = w.ScheduleTask(
			"CRON_TZ=Asia/Saigon */30 * * * *",
			RefreshTokenZaloTask{},
		)

		_, _ = w.ScheduleTask(
			"CRON_TZ=Asia/Saigon 0 */4 * * *",
			RemindUnseenMessageTask{},
		)
	}

}
