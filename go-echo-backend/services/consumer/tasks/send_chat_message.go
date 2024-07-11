package tasks

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
)

type SendChatMessageTask struct {
	ChatMessage *models.ChatMessage `json:"chat_message" validate:"required"`
}

func (task SendChatMessageTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task SendChatMessageTask) TaskName() string {
	return "send_chat_message"
}

// Handler handler
func (task SendChatMessageTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var roomUserIDs []string

	workerInstance.App.DB.Model(&models.ChatRoomUser{}).
		Select("UserID").
		Find(&roomUserIDs, "room_id = ? AND user_id <> ?", task.ChatMessage.ReceiverID, task.ChatMessage.SenderID)

	chatRoom, err := repo.NewChatRoomRepo(workerInstance.App.DB).GetChatRoomInfo(task.ChatMessage.ReceiverID)
	if err != nil {
		return err
	}

	workerInstance.Logger.Debugf("Message receiver_id=%s type=%s total_users=%d err=%+v", task.ChatMessage.ReceiverID, task.ChatMessage.MessageType, len(roomUserIDs), err)
	_, err = SendWSEventTask{
		ChatMessage:           task.ChatMessage,
		ChatRoom:              chatRoom,
		ChannelParticipantIDs: roomUserIDs,
	}.Dispatch(ctx)

	var senderInfo models.User
	if err := workerInstance.App.DB.Select("ID", "Name", "Avatar", "Role", "ContactOwnerIDs").First(&senderInfo, "id = ?", task.ChatMessage.SenderID).Error; err != nil {
		return err
	}

	var roomUsers []*models.User
	if err := workerInstance.App.DB.Select("ID", "Role", "Name").Find(&roomUsers, "id IN ?", roomUserIDs).Error; err != nil {
		return err
	}

	var attachments models.Attachments
	if task.ChatMessage.Attachments != nil {
		attachments = task.ChatMessage.Attachments.GenerateFileURL()
	}
	var pathName = "inquiries"
	if chatRoom.ChatRoomStatus.Stage == enums.ChatRoomStageSample {
		pathName = "samples"
	} else if chatRoom.ChatRoomStatus.Stage == enums.ChatRoomStageBulk {
		pathName = "bulks"
	}
	var recentMessage models.ChatMessage
	if err := workerInstance.App.DB.Select("ID").First(&recentMessage, "id != ? AND receiver_id = ? AND created_at > ? - 3600", task.ChatMessage.ID, task.ChatMessage.ReceiverID, task.ChatMessage.CreatedAt).Error; err != nil {
		if !workerInstance.App.DB.IsRecordNotFoundError(err) {
			return err
		}
	}
	if !senderInfo.Role.IsBuyer() && !senderInfo.Role.IsSeller() {
		// send email to buyer or seller
		var receiveUser, ok = lo.Find(roomUsers, func(user *models.User) bool {
			return user.Role.IsBuyer() || user.Role.IsSeller()
		})
		if ok {
			var baseURL = config.GetInstance().BrandPortalBaseURL
			if receiveUser.Role.IsSeller() {
				baseURL = config.GetInstance().SellerPortalBaseURL
				if pathName == "inquiries" {
					pathName = "rfqs"
				}
			}
			_, err = TrackCustomerIOTask{
				UserID: receiveUser.ID,
				Event:  customerio.EventNewChatMessage,
				Data: map[string]interface{}{
					"room_id":            task.ChatMessage.ReceiverID,
					"message_id":         task.ChatMessage.ID,
					"sender":             senderInfo.GetCustomerIOMetadata(nil),
					"receiver":           receiveUser.GetCustomerIOMetadata(nil),
					"order_type":         chatRoom.ChatRoomStatus.Stage,
					"order_reference_id": chatRoom.ChatRoomStatus.ReferenceID,
					"message":            task.ChatMessage.Message,
					"attachments":        attachments.GenerateFileURL(),
					"is_send_email":      recentMessage.ID == "",
					"chat_url":           fmt.Sprintf("%s/%s/%s?open_chat=true", baseURL, pathName, chatRoom.ChatRoomStatus.StageID),
				},
			}.Dispatch(ctx)
		}

	} else {
		// send email to admins
		for _, user := range roomUsers {
			_, err = TrackCustomerIOTask{
				UserID: user.ID,
				Event:  customerio.EventNewChatMessage,
				Data: map[string]interface{}{
					"room_id":            task.ChatMessage.ReceiverID,
					"message_id":         task.ChatMessage.ID,
					"sender":             senderInfo.GetCustomerIOMetadata(nil),
					"receiver":           user.GetCustomerIOMetadata(nil),
					"order_type":         chatRoom.ChatRoomStatus.Stage,
					"order_reference_id": chatRoom.ChatRoomStatus.ReferenceID,
					"message":            task.ChatMessage.Message,
					"attachments":        attachments,
					"is_send_email":      recentMessage.ID == "",
					"chat_url":           fmt.Sprintf("%s/%s/%s/customer?open_chat=true", config.GetInstance().AdminPortalBaseURL, pathName, chatRoom.ChatRoomStatus.StageID),
				},
			}.Dispatch(ctx)
		}
	}

	return err

}

// Dispatch dispatch event
func (task SendChatMessageTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
