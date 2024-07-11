package enums

type FabricRawStatus string

var (
	FabricRawStatusArrived  FabricRawStatus = "arrived"
	FabricRawStatusWeaving  FabricRawStatus = "weaving"
	FabricRawStatusDying    FabricRawStatus = "dying"
	FabricRawStatusRelaxing FabricRawStatus = "relaxing"
)
