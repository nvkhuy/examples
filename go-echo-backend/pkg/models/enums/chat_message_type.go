package enums

type ChatMessageType string

var (
	ChatMessageTypeUser  ChatMessageType = "user"
	ChatMessageTypeGroup ChatMessageType = "group"
)

func (p ChatMessageType) String() string {
	return string(p)
}

func (p ChatMessageType) DisplayName() string {
	var name = string(p)

	switch p {
	case ChatMessageTypeUser:
		name = "User"
	case ChatMessageTypeGroup:
		name = "Group"
	}

	return name
}
