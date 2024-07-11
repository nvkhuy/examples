package mailer

type TemplateID string

var (
	TemplateIDConfirmMail               TemplateID = "d-784ed23a992b45b197b456212383f5ee"
	TemplateIDResetPassword             TemplateID = "d-c9911d392ac44c2e8e2badc278850690"
	TemplateIDStaffInvitation           TemplateID = "d-4704a8f95fe142c88c5f9f48cc152e91"
	TemplateIDWelcomToBoard             TemplateID = "d-1c4f7bcaf7a44eb9b84488cb98acde5e"
	TemplateIDNotifyUserApproved        TemplateID = "d-d7441b369eb144f094f36453c1eb1875"
	TemplateIDBuyerApproveSkuQuotation  TemplateID = "d-d7792d8f547e4b91ab3fec37de712015"
	TemplateIDBuyerRejectSkuQuotation   TemplateID = "d-2de21ece8dce40979cc8652b11a132b1"
	TemplateIDAdminSentQuotationToBuyer TemplateID = "d-eeba8a44f1684dd2b82791909a20673a"

	TemplateIDSellerQuotationApproved   TemplateID = "d-b339a58afccc445ab4e37e7cc67cc90a"
	TemplateIDSellerQuotationRejected   TemplateID = "d-c0b10a17977d4924b7ff73088c2aa839"
	TemplateIDSellerNewQuotationRequest TemplateID = "d-96f2d00e58da472abd0c26d2af3ae4df"
	TemplateIDNewSubscriber             TemplateID = "d-7c4bf68dfd684f2b8720aafcf115fb03"
)
