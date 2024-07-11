package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/labstack/echo/v4"
)

// GetConstants Generate constants
// @Tags Common
// @Summary Generate constants
// @Description Generate constants
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/constants [get]
func GetConstants(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var constants = models.Constants{
		Features: []*models.FeatureConstant{
			{
				Name:  "Bulk Online Payment",
				Value: enums.FeatureTypeBulkOnlinePayment,
			},
		},
		InvoiceTypes: []*models.InvoiceTypeConstant{
			{
				Name:  enums.InvoiceTypeInquiry.DisplayName(),
				Value: enums.InvoiceTypeInquiry,
			},
			{
				Name:  enums.InvoiceTypeBulkPODepositPayment.DisplayName(),
				Value: enums.InvoiceTypeBulkPODepositPayment,
			},
			{
				Name:  enums.InvoiceTypeBulkPOFirstPayment.DisplayName(),
				Value: enums.InvoiceTypeBulkPOFirstPayment,
			},
			{
				Name:  enums.InvoiceTypeBulkPOSecondPayment.DisplayName(),
				Value: enums.InvoiceTypeBulkPOSecondPayment,
			},
			{
				Name:  enums.InvoiceTypeBulkPOFinalPayment.DisplayName(),
				Value: enums.InvoiceTypeBulkPOFinalPayment,
			},
		},
		InflowBillingAddresses: []*models.InvoiceParty{
			models.DefaultVendorForOnlinePayment,
			models.DefaultVendorForLocal,
		},
		BrandTeamRoles: []*models.BrandTeamRoleConstant{
			{
				Value: enums.BrandTeamRoleManager,
				Name:  enums.BrandTeamRoleManager.DisplayName(),
			},
			{
				Value: enums.BrandTeamRoleStaff,
				Name:  enums.BrandTeamRoleStaff.DisplayName(),
			},
		},
		Roles: []*models.RoleConstant{
			{
				Value: enums.RoleSuperAdmin,
				Name:  enums.RoleSuperAdmin.DisplayName(),
			},
			{
				Value: enums.RoleClient,
				Name:  enums.RoleClient.DisplayName(),
			},
			{
				Value: enums.RoleLeader,
				Name:  enums.RoleLeader.DisplayName(),
			},
			{
				Value: enums.RoleStaff,
				Name:  enums.RoleStaff.DisplayName(),
			},
			{
				Value: enums.RoleSeller,
				Name:  enums.RoleSeller.DisplayName(),
			},
		},
		Teams: []*models.TeamConstant{
			{
				Value: enums.TeamDesigner,
				Name:  enums.TeamDesigner.DisplayName(),
			},
			{
				Value: enums.TeamCustomerService,
				Name:  enums.TeamCustomerService.DisplayName(),
			},
			{
				Value: enums.TeamSales,
				Name:  enums.TeamSales.DisplayName(),
			},
			{
				Value: enums.TeamOperator,
				Name:  enums.TeamOperator.DisplayName(),
			},
			{
				Value: enums.TeamMarketing,
				Name:  enums.TeamMarketing.DisplayName(),
			},
			{
				Value: enums.TeamDev,
				Name:  enums.TeamDev.DisplayName(),
			},
			{
				Value: enums.TeamQA,
				Name:  enums.TeamQA.DisplayName(),
			},
			{
				Value: enums.Finance,
				Name:  enums.Finance.DisplayName(),
			},
		},
		DefaultImageTypes: []*models.DefaultImageTypeConstant{
			{
				Value: enums.DefaultUserAvatar,
				URL:   enums.DefaultUserAvatar.URL(),
			},
		},
		BrandTypes: []*models.BrandTypeConstant{
			{
				Value: enums.BrandTypeBrand,
				Name:  enums.BrandTypeBrand.DisplayName(),
			},
			{
				Value: enums.BrandTypeIndividual,
				Name:  enums.BrandTypeIndividual.DisplayName(),
			},
		},
		AccountStatuses: []*models.AccountStatusConstant{
			{
				Value: enums.AccountStatusActive,
				Name:  enums.AccountStatusActive.DisplayName(),
			},
			{
				Value: enums.AccountStatusInactive,
				Name:  enums.AccountStatusInactive.DisplayName(),
			},
			{
				Value: enums.AccountStatusPendingReview,
				Name:  enums.AccountStatusPendingReview.DisplayName(),
			},
			{
				Value: enums.AccountStatusRejected,
				Name:  enums.AccountStatusRejected.DisplayName(),
			},
			{
				Value: enums.AccountStatusSuspended,
				Name:  enums.AccountStatusSuspended.DisplayName(),
			},
		},
		PaymentStatuses: []*models.PaymentStatusConstant{
			{
				Value: enums.PaymentStatusPaid,
				Name:  enums.PaymentStatusPaid.DisplayName(),
			},
			{
				Value: enums.PaymentStatusUnpaid,
				Name:  enums.PaymentStatusUnpaid.DisplayName(),
			},
			{
				Value: enums.PaymentStatusWaitingConfirm,
				Name:  enums.PaymentStatusWaitingConfirm.DisplayName(),
			},
		},
		AddressTypes: []*models.AddressTypeConstant{
			{
				Value: enums.AddressTypePrimary,
				Name:  enums.AddressTypePrimary.DisplayName(),
			},
		},
		PostStatues: []*models.PostStatusConstant{
			{
				Value: enums.PostStatusNew,
				Name:  enums.PostStatusNew.DisplayName(),
			},
			{
				Value: enums.PostStatusInactive,
				Name:  enums.PostStatusInactive.DisplayName(),
			},
			{
				Value: enums.PostStatusDraft,
				Name:  enums.PostStatusDraft.DisplayName(),
			},
			{
				Value: enums.PostStatusPendingReview,
				Name:  enums.PostStatusPendingReview.DisplayName(),
			},
			{
				Value: enums.PostStatusPublished,
				Name:  enums.PostStatusPublished.DisplayName(),
			},
		},
		DocumentStatuses: []*models.DocumentStatusConstant{
			{
				Value: enums.DocumentStatusNew,
				Name:  enums.DocumentStatusNew.DisplayName(),
			},
			{
				Value: enums.DocumentStatusInactive,
				Name:  enums.DocumentStatusInactive.DisplayName(),
			},
			{
				Value: enums.DocumentStatusDraft,
				Name:  enums.DocumentStatusDraft.DisplayName(),
			},
			{
				Value: enums.DocumentStatusPendingReview,
				Name:  enums.DocumentStatusPendingReview.DisplayName(),
			},
			{
				Value: enums.DocumentStatusPublished,
				Name:  enums.DocumentStatusPublished.DisplayName(),
			},
		},
		ChatMessageTypes: []*models.ChatMessageTypeConstant{
			{
				Value: enums.ChatMessageTypeUser,
				Name:  enums.ChatMessageTypeUser.DisplayName(),
			},

			{
				Value: enums.ChatMessageTypeGroup,
				Name:  enums.ChatMessageTypeGroup.DisplayName(),
			},
		},
		FilterRatings: []*models.FilterRatingConstant{
			{
				Value: enums.FilterRating1,
				Name:  enums.FilterRating1.DisplayName(),
			},
			{
				Value: enums.FilterRating2,
				Name:  enums.FilterRating2.DisplayName(),
			},
			{
				Value: enums.FilterRating3,
				Name:  enums.FilterRating3.DisplayName(),
			},
			{
				Value: enums.FilterRating4,
				Name:  enums.FilterRating4.DisplayName(),
			},
			{
				Value: enums.FilterRating5,
				Name:  enums.FilterRating5.DisplayName(),
			},
		},
		FilterMinOrders: []*models.FilterMinOrderConstant{
			{
				Value: enums.FilterRating10,
				Name:  enums.FilterRating10.DisplayName(),
			},
			{
				Value: enums.FilterRating100,
				Name:  enums.FilterRating100.DisplayName(),
			},
			{
				Value: enums.FilterRating500,
				Name:  enums.FilterRating500.DisplayName(),
			},
		},

		ProductUnits: []*models.ProductUnitConstant{
			{
				Value: enums.ProductUnitPiece,
				Name:  enums.ProductUnitPiece.DisplayName(),
			},
			{
				Value: enums.ProductUnitPair,
				Name:  enums.ProductUnitPair.DisplayName(),
			},
			{
				Value: enums.ProductUnitBox,
				Name:  enums.ProductUnitBox.DisplayName(),
			},
		},

		ProductTypes: []*models.ProductTypeConstant{
			{
				Value: enums.ProductTypeClothing,
				Name:  enums.ProductTypeClothing.DisplayName(),
			},
			{
				Value: enums.ProductTypeFabric,
				Name:  enums.ProductTypeFabric.DisplayName(),
			},
			{
				Value: enums.ProductTypeGraphic,
				Name:  enums.ProductTypeGraphic.DisplayName(),
			},
		},
		InquiryStatuses: []*models.InquiryStatusConstant{
			{
				Value: enums.InquiryStatusNew,
				Name:  enums.InquiryStatusNew.DisplayName(),
			},
			{
				Value: enums.InquiryStatusQuoteInProcess,
				Name:  enums.InquiryStatusQuoteInProcess.DisplayName(),
			},
			{
				Value: enums.InquiryStatusProduction,
				Name:  enums.InquiryStatusProduction.DisplayName(),
			},
			{
				Value: enums.InquiryStatusFinished,
				Name:  enums.InquiryStatusFinished.DisplayName(),
			},
			{
				Value: enums.InquiryStatusClosed,
				Name:  enums.InquiryStatusClosed.DisplayName(),
			},
			{
				Value: enums.InquiryStatusCanceled,
				Name:  enums.InquiryStatusCanceled.DisplayName(),
			},
		},
		InquiryPriceTypes: []*models.InquiryPriceTypeConstant{
			{
				Value: enums.InquiryPriceTypeFOB,
				Name:  enums.InquiryPriceTypeFOB.DisplayName(),
			},
			{
				Value: enums.InquiryPriceTypeCIF,
				Name:  enums.InquiryPriceTypeCIF.DisplayName(),
			},
			{
				Value: enums.InquiryPriceTypeEXW,
				Name:  enums.InquiryPriceTypeEXW.DisplayName(),
			},
		},
		InquiryMOQs: []*models.InquiryMOQConstant{
			{
				Value: enums.InquiryMOQTypeGT50,
				Name:  ">50",
			},
			{
				Value: enums.InquiryMOQTypeGT100,
				Name:  ">100",
			},
			{
				Value: enums.InquiryMOQTypeGT200,
				Name:  ">200",
			},
			{
				Value: enums.InquiryMOQTypeGT300,
				Name:  ">300",
			},
			{
				Value: enums.InquiryMOQTypeGT500,
				Name:  ">500",
			},
			{
				Value: enums.InquiryMOQTypeGT1000,
				Name:  ">1000",
			}, {
				Value: enums.InquiryMOQTypeGT2000,
				Name:  ">2000",
			},
			{
				Value: enums.InquiryMOQTypeGT3000,
				Name:  ">3000",
			},
			{
				Value: enums.InquiryMOQTypeLT100,
				Name:  "<100",
			},
		},
		Certifications: []*models.CertificationConstant{
			{
				Value: enums.CertificationTypeBSCI,
				Name:  "BSCI",
			},
			{
				Value: enums.CertificationTypeOEKO,
				Name:  "OEKO",
			},
			{
				Value: enums.CertificationTypeWRAP,
				Name:  "WRAP",
			},
		},
		SellerQuotationFilters: []*models.SellerQuotationFilterConstant{
			{
				Value: enums.SellerQuotationFilterNew,
				Name:  enums.SellerQuotationFilterNew.DisplayName(),
			},
			{
				Value: enums.SellerQuotationFilterSent,
				Name:  enums.SellerQuotationFilterSent.DisplayName(),
			},
			{
				Value: enums.SellerQuotationFilterWaitingForApproval,
				Name:  enums.SellerQuotationFilterWaitingForApproval.DisplayName(),
			},
		},
		InquirySkuStatuses: []*models.InquirySkuStatusConstant{
			{
				Value: enums.InquirySkuStatusNew,
				Name:  enums.InquirySkuStatusNew.DisplayName(),
			},
			{
				Value: enums.InquirySkuStatusApproved,
				Name:  enums.InquirySkuStatusApproved.DisplayName(),
			},
			{
				Value: enums.InquirySkuStatusRejected,
				Name:  enums.InquirySkuStatusRejected.DisplayName(),
			},
			{
				Value: enums.InquirySkuStatusWaitingForApproval,
				Name:  enums.InquirySkuStatusWaitingForApproval.DisplayName(),
			},
			{
				Value: enums.InquirySkuStatusWaitingForQuotation,
				Name:  enums.InquirySkuStatusWaitingForQuotation.DisplayName(),
			},
		},
		InquiryBuyerStatuses: []*models.InquiryBuyerStatusConstant{
			{
				Value: enums.InquiryBuyerStatusNew,
				Name:  enums.InquiryBuyerStatusNew.DisplayName(),
			},
			{
				Value: enums.InquiryBuyerStatusApproved,
				Name:  enums.InquiryBuyerStatusApproved.DisplayName(),
			},
			{
				Value: enums.InquiryBuyerStatusRejected,
				Name:  enums.InquiryBuyerStatusRejected.DisplayName(),
			},
			{
				Value: enums.InquiryBuyerStatusWaitingForApproved,
				Name:  enums.InquiryBuyerStatusWaitingForApproved.DisplayName(),
			},
		},
		LabelDimensions: []*models.LabelDimensionConstant{
			{
				Value: enums.LabelDimension4060MM,
				Name:  "40x60mm",
			},
			{
				Value: enums.LabelDimension2713MM,
				Name:  "27x13mm",
			},
			{
				Value: enums.LabelDimension6070MM,
				Name:  "60x70mm",
			},
		},
		BarcodeDimensions: []*models.LabelDimensionConstant{
			{
				Value: enums.LabelDimension2713MM,
				Name:  "27x13mm",
			},
		},
		LabelKinds: []*models.LabelKindConstant{
			{
				Value: enums.LabelKindHangTag,
				Name:  "Hang Tag",
			},
			{
				Value: enums.LabelKindCareLabel,
				Name:  "Care Label",
			},
			{
				Value: enums.LabelKindScarfLabel,
				Name:  "Scarf Label",
			},
		},
		LabelSizes: []*models.LabelSizeConstant{
			{
				Value: enums.LabelSizeXS,
				Name:  enums.LabelSizeXS.String(),
			},
			{
				Value: enums.LabelSizeS,
				Name:  enums.LabelSizeS.String(),
			},
			{
				Value: enums.LabelSizeM,
				Name:  enums.LabelSizeM.String(),
			},
			{
				Value: enums.LabelSizeL,
				Name:  enums.LabelSizeL.String(),
			},
			{
				Value: enums.LabelSizeXL,
				Name:  enums.LabelSizeXL.String(),
			},
		},
		LabelMaterials: []*models.LabelMaterialConstant{
			{
				Value: enums.LabelMaterialC200,
				Name:  enums.LabelMaterialC200.String(),
			},
			{
				Value: enums.LabelMaterialC250,
				Name:  enums.LabelMaterialC250.String(),
			},
			{
				Value: enums.LabelMaterialC300,
				Name:  enums.LabelMaterialC300.String(),
			},
			{
				Value: enums.LabelMaterialB200,
				Name:  enums.LabelMaterialB200.String(),
			},
			{
				Value: enums.LabelMaterialB250,
				Name:  enums.LabelMaterialB250.String(),
			},
			{
				Value: enums.LabelMaterialCotton,
				Name:  enums.LabelMaterialCotton.String(),
			},
			{
				Value: enums.LabelMaterialPoly,
				Name:  enums.LabelMaterialPoly.String(),
			},
		},
		LabelStatues: []*models.LabelStatusConstant{
			{
				Value: enums.LabelStatusNew,
				Name:  enums.LabelStatusNew.String(),
			},
			{
				Value: enums.LabelStatusApproved,
				Name:  enums.LabelStatusApproved.String(),
			},
			{
				Value: enums.LabelStatusRejected,
				Name:  enums.LabelStatusRejected.String(),
			},
			{
				Value: enums.LabelStatusPrinting,
				Name:  enums.LabelStatusPrinting.String(),
			},
			{
				Value: enums.LabelStatusClosed,
				Name:  enums.LabelStatusClosed.String(),
			},
			{
				Value: enums.LabelStatusFinished,
				Name:  enums.LabelStatusFinished.String(),
			},
		},
		LabelAccessories: []*models.LabelAccessoryConstant{
			{
				Value: enums.LabelAccessoryHangString,
				Name:  "Hang String",
			},
			{
				Value: enums.LabelAccessoryShoeBox,
				Name:  "Shoe Box",
			},
			{
				Value: enums.LabelAccessoryShoeDustBag,
				Name:  "Shoe Dust Bag",
			},
			{
				Value: enums.LabelAccessoryZiplock,
				Name:  "Ziplock",
			},
		},
		InquirySkuRejectReasons: []*models.InquirySkuRejectReasonConstant{
			{
				Value: enums.InquirySkuRejectReasonUnreasonablePrice,
				Name:  "Unreasonable price",
			},
			{
				Value: enums.InquirySkuRejectReasonChangeMaterial,
				Name:  "Change material",
			},
			{
				Value: enums.InquirySkuRejectReasonChangeQuanity,
				Name:  "Change quantity",
			},
			{
				Value: enums.InquirySkuRejectReasonChangeSize,
				Name:  "Change size",
			},
			{
				Value: enums.InquirySkuRejectReasonOther,
				Name:  "Other",
			},
		},
		CardTypes: []*models.CardTypeConstant{
			{
				Value: enums.CardTypeOrderDesign,
				Name:  "Design",
			},
			{
				Value: enums.CardTypeOrderSizeSpec,
				Name:  "Size spec",
			},
			{
				Value: enums.CardTypeOrderTechpack,
				Name:  "Techpack",
			},
			{
				Value: enums.CardTypeOrderButton,
				Name:  "Button",
			},
			{
				Value: enums.CardTypeOrderLabel,
				Name:  "Label",
			},
			{
				Value: enums.CardTypeOrderThread,
				Name:  "Thread",
			},
			{
				Value: enums.CardTypeOrderZipper,
				Name:  "Zipper",
			},
			{
				Value: enums.CardTypeSkuFabricDetail,
				Name:  "Fabric detail",
			},
			{
				Value: enums.CardTypeProtofitInProgress,
				Name:  "Protofit",
			},
			{
				Value: enums.CardTypeProtofitShipping,
				Name:  "Protofit",
			},
			{
				Value: enums.CardTypeProtofitFeedback,
				Name:  "Protofit",
			},
			{
				Value: enums.CardTypePreProductionInProgress,
				Name:  "PreProduction",
			},
			{
				Value: enums.CardTypePreProductionShipping,
				Name:  "PreProduction",
			},
		},
		CardAttributes: []*models.CardAttributeConstant{
			{
				Value: enums.CardAttributeFabricName,
			},
			{
				Value: enums.CardAttributeColor,
			},
			{
				Value: enums.CardAttributeSize,
			},
			{
				Value: enums.CardAttributeMaterial,
			},
			{
				Value: enums.CardAttributeThreadCode,
			},
			{
				Value: enums.CardAttributeQuantity,
			},
			{
				Value: enums.CardAttributeWeight,
			},
		},
		FabricRawStatuses: []*models.FabricRawStatusConstant{
			{
				Value: enums.FabricRawStatusArrived,
				Name:  "Arrived",
			},
			{
				Value: enums.FabricRawStatusWeaving,
				Name:  "Weaving",
			},
			{
				Value: enums.FabricRawStatusDying,
				Name:  "Dying",
			},
			{
				Value: enums.FabricRawStatusRelaxing,
				Name:  "Relaxing",
			},
		},
		FabricBulkProductionStatuses: []*models.FabricBulkProductionConstant{
			{
				Value: enums.FabricBulkProductionStatusCut,
				Name:  "Cut",
			},
			{
				Value: enums.FabricBulkProductionStatusSew,
				Name:  "Sew",
			},
			{
				Value: enums.FabricBulkProductionStatusFinish,
				Name:  "Finish",
			},
			{
				Value: enums.FabricBulkProductionStatusQc,
				Name:  "Quality Control",
			},
			{
				Value: enums.FabricBulkProductionStatusPack,
				Name:  "Pack",
			},
		},
		QcReportTypes: []*models.QcReportTypeConstant{
			{
				Value: enums.QcReportTypeMatInspection,
				Name:  "Mat Ins.",
			},
			{
				Value: enums.QcReportTypeInlineInspection,
				Name:  "Inline Ins.",
			},
			{
				Value: enums.QcReportTypeEndlineInspection,
				Name:  "Endline Ins.",
			},
			{
				Value: enums.QcReportTypeAqlInspection,
				Name:  "AQL Ins.",
			},
		},
		QcReportResults: []*models.QcReportResultConstant{
			{
				Value: enums.QcReportStatusFail,
				Name:  "Fail",
			},
			{
				Value: enums.QcReportStatusPass,
				Name:  "Pass",
			},
		},
		DeliveryTypes: []*models.DeliveryTypeConstant{
			{
				Value: enums.DeliveryTypeDelivery,
				Name:  "Delivery",
			},
			{
				Value: enums.DeliveryTypeFlight,
				Name:  "Flight/vessel",
			},
			{
				Value: enums.DeliveryTypePickup,
				Name:  "Pickup",
			},
			{
				Value: enums.DeliveryTypeExFactory,
				Name:  "Ex-factory",
			},
		},
		DeliveryAttributeNames: []*models.DeliveryAttributeNameConstant{
			{
				Value: enums.DeliveryAttributeNameCurrierSite,
				Name:  "Currier Site",
			},
			{
				Value: enums.DeliveryAttributeNameTrackingNo,
				Name:  "Tracking No",
			},
		},
		DeliveryStatuses: []*models.DeliveryStatusConstant{
			{
				Value: enums.DeliveryStatusOnTime,
				Name:  "On time",
			},
			{
				Value: enums.DeliveryStatusDelayed,
				Name:  "Delayed",
			},
		},
		InflowPaymentMethods: []*models.InflowPaymentMethodConstant{
			{
				Value: enums.InflowPaymentMethodCard,
				Name:  "Card",
			},
			{
				Value: enums.InflowPaymentMethodOnlineBanking,
				Name:  "Online Banking",
			},
		},
		Currencies: []*models.CurrencyConstant{
			{
				Value: enums.USD,
				Name:  "USD",
			},
			{
				Value: enums.VND,
				Name:  "VND",
			},
		},
		InquiryTypes: []*models.InquiryTypeConstant{
			{
				Value: enums.InquiryTypeBulk,
				Name:  "Bulk",
			},
			{
				Value: enums.InquiryTypeSample,
				Name:  "Sample",
			},
		},

		InquiryAttributes: []*models.InquiryAttributeConstant{
			{
				Value: enums.InquiryAttributeSizeList,
			},
			{
				Value: enums.InquiryAttributeSizeChart,
			},
			{
				Value: enums.InquiryAttributeStyleNo,
			},
			{
				Value: enums.InquiryAttributeFabricName,
			},
			{
				Value: enums.InquiryAttributeFabricWeight,
			},
			{
				Value: enums.InquiryAttributeColorCount,
			},
			{
				Value: enums.InquiryAttributeProductWeight,
			},
		},
		InquirySizeCharts: []*models.InquirySizeChartConstant{
			{
				Value: enums.InquirySizeChartAsia,
			},
			{
				Value: enums.InquirySizeChartUk,
			},
			{
				Value: enums.InquirySizeChartUs,
			},
		},
		ProductAttributeMetas: []*models.ProductAttributeMetaConstant{
			{
				Value: enums.ProductAttributeSize,
			},
			{
				Value: enums.ProductAttributeColor,
			},
		},
		PoCatalogTrackingStatuses: []*models.PoCatalogTrackingStatusConstant{
			{
				Value: enums.PoCatalogTrackingStatusNew,
				Name:  enums.PoCatalogTrackingStatusNew.DisplayName(),
			},
			{
				Value: enums.PoCatalogTrackingStatusPayment,
				Name:  enums.PoCatalogTrackingStatusPayment.DisplayName(),
			},
			{
				Value: enums.PoCatalogTrackingStatusDispatch,
				Name:  enums.PoCatalogTrackingStatusDispatch.DisplayName(),
			},
			{
				Value: enums.PoCatalogTrackingStatusDelivery,
				Name:  enums.PoCatalogTrackingStatusDelivery.DisplayName(),
			},
			{
				Value: enums.PoCatalogTrackingStatusCompleted,
				Name:  enums.PoCatalogTrackingStatusCompleted.DisplayName(),
			},
		},
		PoTrackingStatuses: []*models.PoTrackingStatusConstant{
			{
				Value: enums.PoTrackingStatusNew,
				Name:  enums.PoTrackingStatusNew.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusWaitingForApproved,
				Name:  enums.PoTrackingStatusWaitingForApproved.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusDesignApproved,
				Name:  enums.PoTrackingStatusDesignApproved.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusDesignRejected,
				Name:  enums.PoTrackingStatusDesignRejected.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusRawMaterial,
				Name:  enums.PoTrackingStatusRawMaterial.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusMaking,
				Name:  enums.PoTrackingStatusMaking.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusSubmit,
				Name:  enums.PoTrackingStatusSubmit.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusDelivering,
				Name:  enums.PoTrackingStatusDelivering.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusDeliveryConfirmed,
				Name:  enums.PoTrackingStatusDeliveryConfirmed.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusCanceled,
				Name:  enums.PoTrackingStatusCanceled.DisplayName(),
			},
			{
				Value: enums.PoTrackingStatusPaymentReceived,
				Name:  enums.PoTrackingStatusPaymentReceived.DisplayName(),
			},
		},
		PoRawMaterialStatuses: []*models.PoRawMaterialStatusConstant{
			// {
			// 	Value: enums.PoRawMaterialStatusDying,
			// 	Name:  enums.PoRawMaterialStatusDying.DisplayName(),
			// },
			// {
			// 	Value: enums.PoRawMaterialStatusProcessing,
			// 	Name:  enums.PoRawMaterialStatusProcessing.DisplayName(),
			// },
			{
				Value: enums.PoRawMaterialStatusWaitingForApprove,
				Name:  enums.PoRawMaterialStatusWaitingForApprove.DisplayName(),
			},
			{
				Value: enums.PoRawMaterialStatusApproved,
				Name:  enums.PoRawMaterialStatusApproved.DisplayName(),
			},
		},
		BulkPoTrackingStatuses: []*models.BulkPoTrackingStatusConstant{
			{
				Value: enums.BulkPoTrackingStatusNew,
				Name:  enums.BulkPoTrackingStatusNew.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusWaitingForSubmitOrder,
				Name:  enums.BulkPoTrackingStatusWaitingForSubmitOrder.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusWaitingForQuotation,
				Name:  enums.BulkPoTrackingStatusWaitingForQuotation.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusFirstPayment,
				Name:  enums.BulkPoTrackingStatusFirstPayment.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusFirstPaymentConfirm,
				Name:  enums.BulkPoTrackingStatusFirstPaymentConfirm.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusRawMaterial,
				Name:  enums.BulkPoTrackingStatusRawMaterial.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusPps,
				Name:  enums.BulkPoTrackingStatusPps.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusProduction,
				Name:  enums.BulkPoTrackingStatusProduction.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusQc,
				Name:  enums.BulkPoTrackingStatusQc.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusSubmit,
				Name:  enums.BulkPoTrackingStatusSubmit.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusFinalPayment,
				Name:  enums.BulkPoTrackingStatusFinalPayment.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusFinalPaymentConfirm,
				Name:  enums.BulkPoTrackingStatusFinalPaymentConfirm.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusDelivering,
				Name:  enums.BulkPoTrackingStatusDelivering.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusDeliveryConfirmed,
				Name:  enums.BulkPoTrackingStatusDeliveryConfirmed.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusFirstPaymentConfirmed,
				Name:  enums.BulkPoTrackingStatusFirstPaymentConfirmed.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusSecondPayment,
				Name:  enums.BulkPoTrackingStatusSecondPayment.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusSecondPaymentConfirm,
				Name:  enums.BulkPoTrackingStatusSecondPaymentConfirm.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusSecondPaymentConfirmed,
				Name:  enums.BulkPoTrackingStatusSecondPaymentConfirmed.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusPps,
				Name:  enums.BulkPoTrackingStatusPps.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusFinalPaymentConfirmed,
				Name:  enums.BulkPoTrackingStatusFinalPaymentConfirmed.DisplayName(),
			},
			{
				Value: enums.BulkPoTrackingStatusDelivered,
				Name:  enums.BulkPoTrackingStatusDelivered.DisplayName(),
			},
		},
		BulkPOSellerTrackingStatuses: []*models.BulkPoSellerTrackingStatusConstant{
			// {
			// 	Value: enums.SellerBulkPoTrackingStatusNew,
			// 	Name:  enums.SellerBulkPoTrackingStatusNew.DisplayName(),
			// },
			{
				Value: enums.SellerBulkPoTrackingStatusPO,
				Name:  enums.SellerBulkPoTrackingStatusPO.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusPORejected,
				Name:  enums.SellerBulkPoTrackingStatusPORejected.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusWaitingFirstPayment,
				Name:  enums.SellerBulkPoTrackingStatusWaitingFirstPayment.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusWaitingForSubmitOrder,
				Name:  enums.SellerBulkPoTrackingStatusWaitingForSubmitOrder.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusWaitingForQuotation,
				Name:  enums.SellerBulkPoTrackingStatusWaitingForQuotation.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusFirstPayment,
				Name:  enums.SellerBulkPoTrackingStatusFirstPayment.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusFirstPaymentConfirm,
				Name:  enums.SellerBulkPoTrackingStatusFirstPaymentConfirm.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusFirstPaymentConfirmed,
				Name:  enums.SellerBulkPoTrackingStatusFirstPaymentConfirmed.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusFirstPaymentSkipped,
				Name:  enums.SellerBulkPoTrackingStatusFirstPaymentSkipped.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusRawMaterial,
				Name:  enums.SellerBulkPoTrackingStatusRawMaterial.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusPps,
				Name:  enums.SellerBulkPoTrackingStatusPps.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusProduction,
				Name:  enums.SellerBulkPoTrackingStatusProduction.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusQc,
				Name:  enums.SellerBulkPoTrackingStatusQc.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusSubmit,
				Name:  enums.SellerBulkPoTrackingStatusSubmit.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusFinalPayment,
				Name:  enums.SellerBulkPoTrackingStatusFinalPayment.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusFinalPaymentConfirm,
				Name:  enums.SellerBulkPoTrackingStatusFinalPaymentConfirm.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusFinalPaymentConfirmed,
				Name:  enums.SellerBulkPoTrackingStatusFinalPaymentConfirmed.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusDelivering,
				Name:  enums.SellerBulkPoTrackingStatusDelivering.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusDeliveryConfirmed,
				Name:  enums.SellerBulkPoTrackingStatusDeliveryConfirmed.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusDelivered,
				Name:  enums.SellerBulkPoTrackingStatusDelivered.DisplayName(),
			},
			{
				Value: enums.SellerBulkPoTrackingStatusInspection,
				Name:  enums.SellerBulkPoTrackingStatusInspection.DisplayName(),
			},
		},
		ShippingMethods: []*models.ShippingMethodConstant{
			{
				Value:       enums.ShippingMethodFOB,
				Name:        enums.ShippingMethodFOB.DisplayName(),
				Description: enums.ShippingMethodFOB.Description(),
			},
			{
				Value:       enums.ShippingMethodEXW,
				Name:        enums.ShippingMethodEXW.DisplayName(),
				Description: enums.ShippingMethodEXW.Description(),
			},
			{
				Value:       enums.ShippingMethodCIF,
				Name:        enums.ShippingMethodCIF.DisplayName(),
				Description: enums.ShippingMethodCIF.Description(),
			},
		},
		SupplierTypes: []*models.SupplierTypeConstant{
			{
				Value: enums.SupplierTypeManufacturer,
				Name:  enums.SupplierTypeManufacturer.DisplayName(),
			},
			{
				Value: enums.SupplierTypeMill,
				Name:  enums.SupplierTypeMill.DisplayName(),
			},
			{
				Value: enums.SupplierTypeAccessory,
				Name:  enums.SupplierTypeAccessory.DisplayName(),
			},
			{
				Value: enums.SupplierTypeService,
				Name:  enums.SupplierTypeService.DisplayName(),
			},
			{
				Value: enums.SupplierTypeProductDesigner,
				Name:  enums.SupplierTypeProductDesigner.DisplayName(),
			},
		},
		// Onboarding
		OBOrderTypes: []*models.OnboardingConstant{
			{
				Value:       enums.OBOrderTypeFOB.String(),
				Name:        enums.OBOrderTypeFOB.DisplayName(),
				Description: enums.OBOrderTypeFOB.Description(),
			},
			{
				Value:       enums.OBOrderTypeCMT.String(),
				Name:        enums.OBOrderTypeCMT.DisplayName(),
				Description: enums.OBOrderTypeCMT.Description(),
			},
			{
				Value:       enums.OBOrderTypeCM.String(),
				Name:        enums.OBOrderTypeCM.DisplayName(),
				Description: enums.OBOrderTypeCM.Description(),
			},
		},
		OBShippingTerms: []*models.OnboardingConstant{
			{
				Value:       enums.OBShippingTermExWork.String(),
				Name:        enums.OBShippingTermExWork.DisplayName(),
				Description: enums.OBShippingTermExWork.Description(),
			},
			{
				Value:       enums.OBShippingTermFOB.String(),
				Name:        enums.OBShippingTermFOB.DisplayName(),
				Description: enums.OBShippingTermFOB.Description(),
			},
			{
				Value:       enums.OBShippingTermCIF.String(),
				Name:        enums.OBShippingTermCIF.DisplayName(),
				Description: enums.OBShippingTermCIF.Description(),
			},
		},
		OBProductGroups: []*models.OnboardingConstant{
			{
				Value: enums.OBProductGroupClothing.String(),
				Name:  enums.OBProductGroupClothing.DisplayName(),
			},
			{
				Value: enums.OBProductGroupShoes.String(),
				Name:  enums.OBProductGroupShoes.DisplayName(),
			},
			{
				Value: enums.OBProductGroupBag.String(),
				Name:  enums.OBProductGroupBag.DisplayName(),
			},
			{
				Value: enums.OBProductGroupCap.String(),
				Name:  enums.OBProductGroupCap.DisplayName(),
			},
		},
		OBMOQTypes: []*models.OnboardingConstant{
			{
				Value: enums.OBMOQTypeLT100.String(),
				Name:  enums.OBMOQTypeLT100.DisplayName(),
			},
			{
				Value: enums.OBMOQTypeLT300.String(),
				Name:  enums.OBMOQTypeLT300.DisplayName(),
			},
			{
				Value: enums.OBMOQTypeGT500.String(),
				Name:  enums.OBMOQTypeGT500.DisplayName(),
			},
		},
		OBLeadTimes: []*models.OnboardingConstant{
			{
				Value: enums.OBLeadTime10.String(),
				Name:  enums.OBLeadTime10.DisplayName(),
			},
			{
				Value: enums.OBLeadTime15.String(),
				Name:  enums.OBLeadTime15.DisplayName(),
			},
		},
		OBFabricTypes:           models.GenerateOBFabricTypeConstants(),
		OBMillFabricTypes:       models.GenerateOBMillFabricTypeConstants(),
		OBFactoryProductTypes:   models.GenerateOBFactoryProductTypeConstants(),
		OBSewingAccessoryTypes:  models.GenerateOBSewingAccessoryTypeConstants(),
		OBPackingAccessoryTypes: models.GenerateOBPackingAccessoryTypeConstants(),
		OBDevelopmentServices: []*models.OnboardingConstant{
			{
				Value: enums.OBDevelopmentServiceDesignSketches.String(),
				Name:  enums.OBDevelopmentServiceDesignSketches.DisplayName(),
			},
			{
				Value: enums.OBDevelopmentServiceFabricMaterialCombination.String(),
				Name:  enums.OBDevelopmentServiceFabricMaterialCombination.DisplayName(),
			},
		},
		OBOutputUnits: []*models.OnboardingConstant{
			{
				Value: enums.OBOutputUnitPiece.String(),
				Name:  enums.OBOutputUnitPiece.DisplayName(),
			},
			{
				Value: enums.OBOutputUnitMet.String(),
				Name:  enums.OBOutputUnitMet.DisplayName(),
			},
			{
				Value: enums.OBOutputUnitTon.String(),
				Name:  enums.OBOutputUnitTon.DisplayName(),
			},
		},
		OBServiceTypes: []*models.OnboardingConstant{
			{
				Value: enums.OBServiceTypeLogisticAndShipping.String(),
				Name:  enums.OBServiceTypeLogisticAndShipping.DisplayName(),
			},
			{
				Value: enums.OBServiceTypeTesting.String(),
				Name:  enums.OBServiceTypeTesting.DisplayName(),
			},
			{
				Value: enums.OBServiceTypeDecoration.String(),
				Name:  enums.OBServiceTypeDecoration.DisplayName(),
			},
			{
				Value: "other",
				Name:  "Other",
			},
		},
		OBDecorationServices: []*models.OnboardingConstant{
			{
				Value: enums.OBDecorationServiceWashing.String(),
				Name:  enums.OBDecorationServiceWashing.DisplayName(),
			},
			{
				Value: enums.OBDecorationServiceDrying.String(),
				Name:  enums.OBDecorationServiceDrying.DisplayName(),
			},
			{
				Value: "other",
				Name:  "Other",
			},
		},
		OBPaymentTerms: []*models.OnboardingConstant{
			{
				Value: enums.OBPaymentTermFullPayment.String(),
				Name:  enums.OBPaymentTermFullPayment.DisplayName(),
			},
			{
				Value: enums.OBPaymentTermPartialPayment.String(),
				Name:  enums.OBPaymentTermPartialPayment.DisplayName(),
			},
			{
				Value: enums.OBPaymentTermOpenForInvoice.String(),
				Name:  enums.OBPaymentTermOpenForInvoice.DisplayName(),
			},
		},
		// End Onboarding

		InquirySellerStatuses: []*models.InquirySellerStatusConstant{
			{
				Value: enums.InquirySellerStatusNew,
				Name:  enums.InquirySellerStatusNew.DisplayName(),
			},
			{
				Value: enums.InquirySellerStatusOfferRejected,
				Name:  enums.InquirySellerStatusOfferRejected.DisplayName(),
			},
			{
				Value: enums.InquirySellerStatusWaitingForQuotation,
				Name:  enums.InquirySellerStatusWaitingForQuotation.DisplayName(),
			},
			{
				Value: enums.InquirySellerStatusWaitingForApproval,
				Name:  enums.InquirySellerStatusWaitingForApproval.DisplayName(),
			},
			{
				Value: enums.InquirySellerStatusApproved,
				Name:  enums.InquirySellerStatusApproved.DisplayName(),
			},
			{
				Value: enums.InquirySellerStatusRejected,
				Name:  enums.InquirySellerStatusRejected.DisplayName(),
			},
		},
		CommentTargetTypes: []*models.CommentTargetTypeConstant{
			{
				Value: enums.CommentTargetTypeInquirySellerDesign,
			},
		},
		InquirySellerQuotationTypes: []*models.InquirySellerQuotationTypeConstant{
			{
				Value: enums.InquirySellerQuotationTypeMCQ,
				Name:  enums.InquirySellerQuotationTypeMCQ.DisplayName(),
			},
			{
				Value: enums.InquirySellerQuotationTypeMOQ,
				Name:  enums.InquirySellerQuotationTypeMOQ.DisplayName(),
			},
		},
		AdsVideoSections: []*models.AdsVideoSectionConstant{
			{
				Value: enums.AdsVideoSectionRFQ,
				Name:  enums.AdsVideoSectionRFQ.DisplayName(),
			},
			{
				Value: enums.AdsVideoSectionSample,
				Name:  enums.AdsVideoSectionSample.DisplayName(),
			},
			{
				Value: enums.AdsVideoSectionBulk,
				Name:  enums.AdsVideoSectionBulk.DisplayName(),
			},
			{
				Value: enums.AdsVideoSectionCatalogue,
				Name:  enums.AdsVideoSectionCatalogue.DisplayName(),
			},
		},
		FabricWeightUnits: []*models.FabricWeightUnitConstant{
			{
				Value: enums.FabricWeightUnitGSM,
				Name:  enums.FabricWeightUnitGSM.DisplayName(),
			},
			{
				Value: enums.FabricWeightUnitOZ,
				Name:  enums.FabricWeightUnitOZ.DisplayName(),
			},
		},
		MaterialTypes: []*models.MaterialTypesConstant{
			{
				Name:  "Trim",
				Value: "trim",
			},
			{
				Name:  "Fabric",
				Value: "fabric",
			},
			{
				Name:  "Zipper",
				Value: "zipper",
			},
			{
				Name:  "Button",
				Value: "button",
			},
		},
		BulkQCReportTypes: []*models.BulkQCReportTypesConstant{
			{
				Name:  "Inline Inspection",
				Value: "inline_inspection",
			},
			{
				Name:  "Final Inspection",
				Value: "final_inspection",
			},
		},
		BulkPOSellerStatuses: []*models.BulkPOSellerStatusConstant{
			{
				Value: enums.BulkPurchaseOrderSellerStatusWaitingForQuotation,
				Name:  enums.BulkPurchaseOrderSellerStatusWaitingForQuotation.DisplayName(),
			},
			{
				Value: enums.BulkPurchaseOrderSellerStatusWaitingForApproval,
				Name:  enums.BulkPurchaseOrderSellerStatusWaitingForApproval.DisplayName(),
			},
			{
				Value: enums.BulkPurchaseOrderSellerStatusApproved,
				Name:  enums.BulkPurchaseOrderSellerStatusApproved.DisplayName(),
			},
			{
				Value: enums.BulkPurchaseOrderSellerStatusRejected,
				Name:  enums.BulkPurchaseOrderSellerStatusRejected.DisplayName(),
			},
		},
	}

	return cc.Success(constants)
}

// GetRegisterConstants Generate register constants
// @Tags Common
// @Summary Generate register constants
// @Description Generate register constants
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/register_constants [get]
func GetRegisterConstants(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var top_product_categories = []*models.RegisterProductCategory{
		{
			Name:  "Dress",
			Value: "dress",
		},
		{
			Name:  "Hoodies",
			Value: "hoodies",
		},
		{
			Name:  "Sweater",
			Value: "sweater",
		},
		{
			Name:  "Skirts",
			Value: "skirts",
		},
		{
			Name:  "Coats",
			Value: "coats",
		},
	}

	var constants = models.RegisterConstants{
		RegisterBusinesses: []*models.RegisterBusinessConstant{
			{
				Value:   enums.FindProduct,
				Name:    enums.FindProduct.DisplayName(),
				IconUrl: enums.FindProduct.IconUrl(),
			},
			{
				Value:   enums.FindManufacturer,
				Name:    enums.FindManufacturer.DisplayName(),
				IconUrl: enums.FindManufacturer.IconUrl(),
			},
			{
				Value:   enums.FindDesigner,
				Name:    enums.FindDesigner.DisplayName(),
				IconUrl: enums.FindDesigner.IconUrl(),
			},
			{
				Value:   enums.FindOther,
				Name:    enums.FindOther.DisplayName(),
				IconUrl: enums.FindOther.IconUrl(),
			},
		},

		RegisterQuantities: []*models.RegisterQuantityConstant{
			{
				Value: enums.Quantity_Range_50,
				Name:  enums.Quantity_Range_50.DisplayName(),
			},
			{
				Value: enums.Quantity_Range_100,
				Name:  enums.Quantity_Range_100.DisplayName(),
			},
			{
				Value: enums.Quantity_Range_1000,
				Name:  enums.Quantity_Range_1000.DisplayName(),
			},
			{
				Value: enums.Quantity_Range_10000,
				Name:  enums.Quantity_Range_10000.DisplayName(),
			},
		},

		RegisterAreas: []*models.RegisterAreaConstant{
			{
				Value:   enums.AreaUS,
				Name:    enums.AreaUS.DisplayName(),
				IconUrl: enums.AreaUS.IconUrl(),
			},
			{
				Value:   enums.AreaEU,
				Name:    enums.AreaEU.DisplayName(),
				IconUrl: enums.AreaEU.IconUrl(),
			},
			{
				Value:   enums.AreaJapan,
				Name:    enums.AreaJapan.DisplayName(),
				IconUrl: enums.AreaJapan.IconUrl(),
			},
			{
				Value:   enums.AreaKorea,
				Name:    enums.AreaKorea.DisplayName(),
				IconUrl: enums.AreaKorea.IconUrl(),
			},
			{
				Value:   enums.AreaSoutheastAsia,
				Name:    enums.AreaSoutheastAsia.DisplayName(),
				IconUrl: enums.AreaSoutheastAsia.IconUrl(),
			},
			{
				Value:   enums.AreaChina,
				Name:    enums.AreaChina.DisplayName(),
				IconUrl: enums.AreaChina.IconUrl(),
			},
			{
				Value:   enums.AreaIndia,
				Name:    enums.AreaIndia.DisplayName(),
				IconUrl: enums.AreaIndia.IconUrl(),
			},
		},
		RegisterProductCategories: top_product_categories,
	}

	return cc.Success(constants)
}
