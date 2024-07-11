package enums

type QcReportType string

var (
	QcReportTypeMatInspection     QcReportType = "mat_inspection"
	QcReportTypeInlineInspection  QcReportType = "inline_inspection"
	QcReportTypeEndlineInspection QcReportType = "endline_inspection"
	QcReportTypeAqlInspection     QcReportType = "aql_inspection"
)
