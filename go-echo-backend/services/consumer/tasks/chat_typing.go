package tasks

import (
	"context"
	"strings"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/ws"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type ChatTypingTask struct {
	RoomID       string `json:"room_id" validate:"required"`
	TypingUserID string `json:"typing_user_id" validate:"required"`
}

func (task ChatTypingTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task ChatTypingTask) TaskName() string {
	return "chat_typing"
}

// Handler handler
func (task ChatTypingTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var userIDs []string

	workerInstance.App.DB.Model(&models.ChatRoomUser{}).
		Select("UserID").
		Find(&userIDs, "room_id = ? AND user_id <> ?", task.RoomID, task.TypingUserID)

	chatRoom, err := repo.NewChatRoomRepo(workerInstance.App.DB).GetChatRoomInfo(task.RoomID)
	if err != nil {
		return err
	}
	err = ws.GetInstance().BroadcastToUsers(&ws.BroadcastChatMessage{
		Type:           enums.ChatMessageWsTypeTyping.String(),
		Message:        chatRoom.LatestMessage,
		ChatRoom:       chatRoom,
		ParticipantIDs: userIDs,
	})

	if err != nil {
		workerInstance.Logger.Debugf("Send ws event room=%s receiver=%s seen_user=%s error=%+v", task.RoomID, strings.Join(userIDs, ","), task.TypingUserID, err)
		return eris.Wrapf(err, "Send ws event room=%s receiver=%s seen_user=%s", task.RoomID, strings.Join(userIDs, ","), task.TypingUserID)
	}

	return nil

}

// Dispatch dispatch event
func (task ChatTypingTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
