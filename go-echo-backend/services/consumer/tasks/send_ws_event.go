package tasks

import (
	"context"
	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/ws"
	"github.com/hibiken/asynq"

	"github.com/rotisserie/eris"
)

type SendWSEventTask struct {
	ChatMessage           *models.ChatMessage `json:"chat_message" validate:"required"`
	ChatRoom              *models.ChatRoom    `json:"chat_room" validate:"required"`
	ChannelParticipantIDs []string            `json:"channel_participant_ids"`
}

func (task SendWSEventTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task SendWSEventTask) TaskName() string {
	return "send_ws_event"
}

func (task SendWSEventTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

// Handler handler
func (task SendWSEventTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	err = ws.GetInstance().BroadcastToUsers(&ws.BroadcastChatMessage{
		Type:           enums.ChatMessageWsTypeReceiveMsg.String(),
		Message:        task.ChatMessage,
		ChatRoom:       task.ChatRoom,
		ParticipantIDs: task.ChannelParticipantIDs,
	})

	if err != nil {
		workerInstance.Logger.Debugf("Send ws event type=%s receiver=%s sender=%s error=%+v", task.ChatMessage.MessageType, task.ChatMessage.ReceiverID, task.ChatMessage.SenderID, err)
		return eris.Wrapf(err, "Send ws event type=%s receiver=%s sender=%s", task.ChatMessage.MessageType, task.ChatMessage.ReceiverID, task.ChatMessage.SenderID)
	}

	return nil
}
