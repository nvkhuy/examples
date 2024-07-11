package enums

type QcReportResult string

var (
	QcReportStatusPass QcReportResult = "pass"
	QcReportStatusFail QcReportResult = "fail"
)
