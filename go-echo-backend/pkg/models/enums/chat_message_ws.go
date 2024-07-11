package enums

type ChatMessageWsType string

var (
	ChatMessageWsTypeReceiveMsg   ChatMessageWsType = "receive_message"
	ChatMessageWsTypeSeenRoom     ChatMessageWsType = "mark_seen"
	ChatMessageWsTypeTyping       ChatMessageWsType = "typing"
	ChatMessageWsTypeCancelTyping ChatMessageWsType = "cancel_typing"
)

func (p ChatMessageWsType) String() string {
	return string(p)
}
