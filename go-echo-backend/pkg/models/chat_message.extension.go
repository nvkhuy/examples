package models

import (
	"encoding/json"
)

func (msg ChatMessage) ToJSONRaw() json.RawMessage {
	bytes, _ := json.Marshal(&msg)

	return bytes
}

func (msg WSMessagePayload) ToJSONRaw() json.RawMessage {
	bytes, _ := json.Marshal(&msg)

	return bytes
}
