package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"

	"github.com/hibiken/asynq"
	"github.com/samber/lo"
)

type RemindUnseenMessageTask struct{}

func (task RemindUnseenMessageTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task RemindUnseenMessageTask) TaskName() string {
	return "remind_unseen_message"
}

func (task RemindUnseenMessageTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	unseenRooms, err := repo.NewChatRoomRepo(workerInstance.App.DB).GetUnseenChatRoomsOver6Hours()
	if err != nil {
		return err
	}

	for _, room := range unseenRooms {
		senderInfo, ok := lo.Find(room.RoomUsers, func(user *models.User) bool {
			if room.LatestMessage == nil {
				return false
			}
			return user.ID == room.LatestMessage.SenderID
		})
		if !ok {
			continue
		}

		var attachments models.Attachments
		if room.LatestMessage.Attachments != nil {
			attachments = room.LatestMessage.Attachments.GenerateFileURL()
		}
		var pathName = "inquiries"
		if room.ChatRoomStatus.Stage == enums.ChatRoomStageSample {
			pathName = "samples"
		} else if room.ChatRoomStatus.Stage == enums.ChatRoomStageBulk {
			pathName = "bulks"
		}
		if !senderInfo.Role.IsBuyer() && !senderInfo.Role.IsSeller() {
			receiver, ok := lo.Find(room.RoomUsers, func(user *models.User) bool {
				return user.Role.IsBuyer() || user.Role.IsSeller()
			})
			if !ok {
				continue
			}
			var baseURL = config.GetInstance().BrandPortalBaseURL
			if receiver.Role.IsSeller() {
				baseURL = config.GetInstance().SellerPortalBaseURL
				if pathName == "inquiries" {
					pathName = "rfqs"
				}
			}
			_, err := TrackCustomerIOTask{
				UserID: receiver.ID,
				Event:  customerio.EventRemindUnseenChatMessage,
				Data: map[string]interface{}{
					"order_type":         room.ChatRoomStatus.Stage,
					"order_reference_id": room.ChatRoomStatus.ReferenceID,
					"sender":             senderInfo,
					"receiver":           receiver.GetCustomerIOMetadata(nil),
					"room_id":            room.ID,
					"message_id":         room.LatestMessage.ID,
					"message":            room.LatestMessage.Message,
					"attachments":        attachments,
					"chat_url":           fmt.Sprintf("%s/%s/%s?open_chat=true", baseURL, pathName, room.ChatRoomStatus.StageID),
				},
			}.Dispatch(ctx)
			if err != nil {
				return err
			}
		} else {
			var admins = lo.Filter(room.RoomUsers, func(user *models.User, index int) bool {
				return !user.Role.IsBuyer() && !user.Role.IsSeller()
			})
			for _, admin := range admins {
				_, err := TrackCustomerIOTask{
					UserID: admin.ID,
					Event:  customerio.EventRemindUnseenChatMessage,
					Data: map[string]interface{}{
						"order_type":         room.ChatRoomStatus.Stage,
						"order_reference_id": room.ChatRoomStatus.ReferenceID,
						"sender":             senderInfo,
						"receiver":           admin.GetCustomerIOMetadata(nil),
						"room_id":            room.ID,
						"message_id":         room.LatestMessage.ID,
						"message":            room.LatestMessage.Message,
						"attachments":        attachments,
						"chat_url":           fmt.Sprintf("%s/%s/%s/customer?open_chat=true", config.GetInstance().AdminPortalBaseURL, pathName, room.ChatRoomStatus.StageID),
					},
				}.Dispatch(ctx)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (task RemindUnseenMessageTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
