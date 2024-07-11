package enums

type ChatRoomStage string

var (
	ChatRoomStageRFQ    ChatRoomStage = "RFQ"
	ChatRoomStageSample ChatRoomStage = "Sample"
	ChatRoomStageBulk   ChatRoomStage = "Bulk"
)

func (status ChatRoomStage) String() string {
	return string(status)
}

func (status ChatRoomStage) DisplayName() string {
	var name = string(status)

	switch status {
	case ChatRoomStageRFQ:
		name = "RFQ"
	case ChatRoomStageSample:
		name = "Sample"
	case ChatRoomStageBulk:
		name = "Bulk"
	}

	return name
}
