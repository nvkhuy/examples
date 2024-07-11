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

type SeenChatRoomTask struct {
	RoomID     string `json:"room_id" validate:"required"`
	SeenUserID string `json:"seen_user_id" validate:"required"`
}

func (task SeenChatRoomTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task SeenChatRoomTask) TaskName() string {
	return "seen_chat_room"
}

// Handler handler
func (task SeenChatRoomTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var userIDs []string

	workerInstance.App.DB.Model(&models.ChatRoomUser{}).
		Select("UserID").
		Find(&userIDs, "room_id = ? AND user_id <> ?", task.RoomID, task.SeenUserID)

	chatRoom, err := repo.NewChatRoomRepo(workerInstance.App.DB).GetChatRoomInfo(task.RoomID)
	if err != nil {
		return err
	}
	err = ws.GetInstance().BroadcastToUsers(&ws.BroadcastChatMessage{
		Type:           enums.ChatMessageWsTypeSeenRoom.String(),
		Message:        chatRoom.LatestMessage,
		ChatRoom:       chatRoom,
		ParticipantIDs: userIDs,
	})

	if err != nil {
		workerInstance.Logger.Debugf("Send ws event room=%s receiver=%s seen_user=%s error=%+v", task.RoomID, strings.Join(userIDs, ","), task.SeenUserID, err)
		return eris.Wrapf(err, "Send ws event room=%s receiver=%s seen_user=%s", task.RoomID, strings.Join(userIDs, ","), task.SeenUserID)
	}

	return nil

}

// Dispatch dispatch event
func (task SeenChatRoomTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
