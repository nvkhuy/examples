package enums

type LabelKind string

var (
	LabelKindHangTag     LabelKind = "hang_tag"
	LabelKindSizeTag     LabelKind = "size_tag"
	LabelKindCareLabel   LabelKind = "care_label"
	LabelKindScarfLabel  LabelKind = "scarf_label"
	LabelKindShoeSticker LabelKind = "shoe_sticker"
)

func (p LabelKind) String() string {
	return string(p)
}
