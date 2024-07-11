package models

type DataAnalyticChart struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}
type DataAnalyticNewUser struct {
	Count  int64               `json:"count"`
	Charts []DataAnalyticChart `json:"charts"`
}

type DataAnalyticNewCatalogProduct struct {
	Count  int64               `json:"count"`
	Charts []DataAnalyticChart `json:"charts"`
}

type DataAnalyticInquiry struct {
	Num                int64 `json:"num"`
	AssignedNum        int64 `json:"assigned_num"`
	UnAssignedNum      int64 `json:"un_assigned_num"`
	QuoteOverdue       int64 `json:"quote_overdue"`
	QuotePotentialLate int64 `json:"potential_late"`

	RFQSubmitted      int64 `json:"rfq_submitted"`
	SendQuotation     int64 `json:"send_quotation"`
	QuotationApproved int64 `json:"quotation_approved"`
	WaitingForPayment int64 `json:"waiting_for_payment"`
	PaymentConfirmed  int64 `json:"payment_confirmed"`
}

type DataAnalyticPO struct {
	Num                int64 `json:"num"`
	AssignedNum        int64 `json:"assigned_num"`
	UnAssignedNum      int64 `json:"un_assigned_num"`
	QuoteOverdue       int64 `json:"quote_overdue"`
	QuotePotentialLate int64 `json:"potential_late"`

	SampleOrder int64 `json:"sample_order"`
	Design      int64 `json:"design"`
	RawMaterial int64 `json:"raw_material"`
	Review      int64 `json:"review"`
	Making      int64 `json:"making"`
	Submit      int64 `json:"submit"`
	Delivery    int64 `json:"delivery"`
	Approval    int64 `json:"approval"`
}

type DataAnalyticBulkPO struct {
	Num                int64 `json:"num"`
	AssignedNum        int64 `json:"assigned_num"`
	UnAssignedNum      int64 `json:"un_assigned_num"`
	QuoteOverdue       int64 `json:"quote_overdue"`
	QuotePotentialLate int64 `json:"potential_late"`

	New          int64 `json:"new"`
	Review       int64 `json:"review"`
	FirstPayment int64 `json:"first_payment"`
	Making       int64 `json:"making"`
	Submit       int64 `json:"submit"`
	FinalPayment int64 `json:"final_payment"`
	Delivery     int64 `json:"delivery"`
}

type DataAnalyticPerformance struct {
	// Ops
	InquiryQuoteInTime                int64   `json:"inquiry_quote_in_time"`
	InquiryQuoteInTimeDiffPercentage  float64 `json:"inquiry_quote_in_time_diff_percentage"`
	InquiryQuoteInTimeTotalPercentage float64 `json:"inquiry_quote_in_time_total_percentage"`

	POInLeadTime                int64   `json:"po_in_lead_time"`
	POInLeadTimeDiffPercentage  float64 `json:"po_in_lead_time_diff_percentage"`
	POInLeadTimeTotalPercentage float64 `json:"po_in_lead_time_total_percentage"`

	BulkPOInLeadTime                int64   `json:"bulk_po_in_lead_time"`
	BulkPOInLeadTimeDiffPercentage  float64 `json:"bulk_po_in_lead_time_diff_percentage"`
	BulkPOInLeadTimeTotalPercentage float64 `json:"bulk_po_in_lead_time_total_percentage"`

	// Biz
	InquiryApproved                int64   `json:"inquiry_approved"`
	InquiryApprovedDiffPercentage  float64 `json:"inquiry_approved_diff_percentage"`
	InquiryApprovedTotalPercentage float64 `json:"inquiry_approved_total_percentage"`

	POPaid                int64   `json:"po_paid"`
	POPaidDiffPercentage  float64 `json:"po_paid_diff_percentage"`
	POPaidTotalPercentage float64 `json:"po_paid_total_percentage"`

	BulkPOPaid                int64   `json:"bulk_po_paid"`
	BulkPOPaidDiffPercentage  float64 `json:"bulk_po_paid_diff_percentage"`
	BulkPOPaidTotalPercentage float64 `json:"bulk_po_paid_total_percentage"`
}
