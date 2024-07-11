package enums

type CardType string

var (
	CardTypeOrderDesign     CardType = "order_design"
	CardTypeOrderTechpack   CardType = "order_techpack"
	CardTypeOrderSizeSpec   CardType = "order_size_spec"
	CardTypeOrderZipper     CardType = "order_zipper"
	CardTypeOrderThread     CardType = "order_thread"
	CardTypeOrderLabel      CardType = "order_label"
	CardTypeOrderButton     CardType = "order_button"
	CardTypeSkuFabricDetail CardType = "sku_fabric_detail"

	CardTypeProtofitInProgress CardType = "protofit_inprogress"
	CardTypeProtofitShipping   CardType = "protofit_shipping"
	CardTypeProtofitFeedback   CardType = "protofit_feedback"

	CardTypePreProductionInProgress CardType = "pre_production_inprogress"
	CardTypePreProductionShipping   CardType = "pre_production_shipping"
)
