package enums

type FabricBulkProductionStatus string

var (
	FabricBulkProductionStatusCut    FabricBulkProductionStatus = "cut"
	FabricBulkProductionStatusSew    FabricBulkProductionStatus = "sew"
	FabricBulkProductionStatusFinish FabricBulkProductionStatus = "finish"
	FabricBulkProductionStatusQc     FabricBulkProductionStatus = "qc"
	FabricBulkProductionStatusPack   FabricBulkProductionStatus = "pack"
)
