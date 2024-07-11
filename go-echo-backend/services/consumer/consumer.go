package consumer

import (
	"context"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/worker"
	"github.com/engineeringinflow/inflow-backend/pkg/ws"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/rotisserie/eris"
	"github.com/thaitanloi365/go-utils/values"
)

var instance *worker.Worker

// New consumer
func New(app *app.App, IsConsumer bool) *worker.Worker {
	instance = worker.New(&worker.Config{
		Namespace:    "inflow_consumer",
		DefaultQueue: "inflow_consumer_queue",
		Logger:       logger.New("consumer/queue"),
		App:          app,
		Size:         1024,
		IsConsumer:   IsConsumer,
	})

	tasks.Register(instance, IsConsumer)

	return instance
}

func WSMessageHandler(message ws.Message) error {
	if message.UserID == "" {
		return eris.New("User ID is required")
	}
	switch message.Type {
	case ws.MessageTypePing:
		tasks.UserPingTask{
			UserID:       message.UserID,
			IsOffline:    values.Bool(false),
			LastOnlineAt: values.Int64(time.Now().Unix()),
		}.Dispatch(context.Background())

	case ws.MessageTypeTyping:
		tasks.ChatTypingTask{
			RoomID:       message.Data.ChatRoomID,
			TypingUserID: message.UserID,
		}.Dispatch(context.Background())
	case ws.MessageTypeCancelTyping:
		tasks.CancelChatTypingTask{
			RoomID:             message.Data.ChatRoomID,
			CancelTypingUserID: message.UserID,
		}.Dispatch(context.Background())

	case ws.MessageTypeSeenMessage:

	}
	return nil
}
