package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/middlewares"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/worker"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/thaitanloi365/melody"
)

var instance *WS

// WS ws
type WS struct {
	melody             *melody.Melody
	app                *app.App
	messageHandlerFunc func(msg Message) error
	Logger             *logger.Logger
}

// New init
func New(e *echo.Echo, app *app.App, messageHandlerFunc func(msg Message) error) *WS {
	instance = &WS{
		melody:             melody.New(),
		Logger:             logger.New("chat/queue"),
		app:                app,
		messageHandlerFunc: messageHandlerFunc,
	}

	e.GET("/ws", instance.handleRequest, middlewares.IsAuthorizedWithQueryToken(app.Config.JWTSecret))

	instance.melody.HandleMessage(instance.handleMessage)

	instance.melody.HandleConnect(instance.handleConnect)

	instance.melody.HandleDisconnect(instance.handleDisconnect)

	return instance
}

// GetInstance get instance
func GetInstance() *WS {
	if instance == nil {
		panic("Must be call New() first")
	}

	return instance
}

func (ws *WS) GetAllSessions() (result []*Session) {
	sessions, err := ws.melody.Sessions()
	if err != nil {
		return nil
	}

	for _, session := range sessions {
		if sessionID, found := session.Get(UserSessionIDContextKey); found {
			if str, ok := sessionID.(string); ok {
				var sess = SessionID(str).GetSession()
				if sess != nil {

					if v, found := session.Get(UserConnectedAtKey); found {
						sess.ConnectedAt = v
					}

					sess.RemoteAddress = session.RemoteAddr().String()

					result = append(result, sess)
				}
			}
		}
	}
	return
}

func (ws *WS) GetUserSessions(userIDs ...string) (result []*melody.Session) {
	if len(userIDs) == 0 {
		return
	}

	sessions, err := ws.melody.Sessions()
	if err != nil {
		return nil
	}

	for _, participantID := range userIDs {
		for _, session := range sessions {
			if sessionID, found := session.Get(UserSessionIDContextKey); found {
				if str, ok := sessionID.(string); ok {
					var sess = SessionID(str).GetSession()
					if sess != nil && participantID == sess.UserID {
						session.Set("user_id", sess.UserID)
						result = append(result, session)
					}
				}
			}
		}
	}

	return
}

func (ws *WS) BroadcastToUser(msg *Message) error {
	var sessions = ws.GetUserSessions(msg.UserID)

	ws.Logger.Debugf("Broadcast to team receiver=%v type=%s total_sessions=%d", msg.UserID, msg.Type, len(sessions))

	if len(sessions) == 0 {
		return nil
	}

	for _, sess := range sessions {
		if v, found := sess.Get("user_id"); found {
			if str, ok := v.(string); ok && str != "" {
				var err = ws.melody.BroadcastMultiple(msg.ToJSONRaw(), []*melody.Session{sess})
				if err != nil {
					ws.Logger.Debugf("Broadcast to team receiver=%s type=%s total_sessions=%d error=%+v", msg.UserID, msg.Type, len(sessions), err)
				}
			}
		}

	}

	return nil
}

func (ws *WS) BroadcastToUsers(msg *BroadcastChatMessage) error {
	var sessions = ws.GetUserSessions(msg.ParticipantIDs...)

	if len(sessions) == 0 {
		return nil
	}

	for _, sess := range sessions {
		var err = ws.melody.BroadcastMultiple(models.WSMessagePayload{Type: msg.Type, Message: msg.Message, ChatRoom: msg.ChatRoom}.ToJSONRaw(), []*melody.Session{sess})
		if err != nil {
			ws.Logger.Debugf("Broadcast to team receiver=%s type=%s total_sessions=%d error=%+v", msg.Message.ReceiverID, msg.Message.MessageType, len(sessions), err)
		}
	}

	return nil
}

func (ws *WS) handleRequest(c echo.Context) error {
	var token = c.QueryParam("token")
	var claims models.JwtClaims
	var err = claims.ValidateToken(ws.app.Config.JWTSecret, token)

	if err != nil {
		ws.Logger.ErrorAny(err)
		return err
	}

	var session = Session{
		UserID:        claims.ID,
		Role:          claims.Audience,
		RemoteAddress: c.RealIP(),
	}
	if session.RemoteAddress == "" {
		session.RemoteAddress = c.Request().RemoteAddr
	}
	var sessionID = session.GetID()
	c.Request().Header.Set(UserSessionIDContextKey, sessionID.String())

	err = ws.melody.HandleRequest(c.Response(), c.Request())
	if err != nil {
		ws.Logger.ErrorAny(err, fmt.Sprintf("User %s", sessionID))
		return err
	}

	return nil
}

func (ws *WS) handleConnect(s *melody.Session) {
	var sessionID = SessionID(s.Request.Header.Get(UserSessionIDContextKey))
	ws.Logger.Debugf("User connected: %s", sessionID)
	if sessionID != "" {
		if sess := sessionID.GetSession(); sess != nil {
			s.Set(UserConnectedAtKey, sess.ConnectedAt)
			s.Set(UserSessionIDContextKey, sessionID.String())

			worker.GetInstance().Client.EnqueueContext(
				s.Request.Context(),
				asynq.NewTask("user_ping", helper.ToJson(map[string]interface{}{
					"user_id":        sess.UserID,
					"is_offline":     false,
					"last_online_at": time.Now().Unix(),
				})),
				worker.QueueLow,
				asynq.MaxRetry(1),
				asynq.Retention(time.Hour*4),
			)

		}

	}

}

func (ws *WS) handleDisconnect(s *melody.Session) {
	sessionID, found := s.Get(UserSessionIDContextKey)

	if !found {
		ws.Logger.Debugf("User %s not found in sessions", sessionID)
		return
	}

	if str, ok := sessionID.(string); ok {
		if sess := SessionID(str).GetSession(); sess != nil {
			s.Set(UserConnectedAtKey, sess.ConnectedAt)
			s.Set(UserSessionIDContextKey, sessionID)

			_, _ = worker.GetInstance().Client.EnqueueContext(
				s.Request.Context(),
				asynq.NewTask("user_ping", helper.ToJson(map[string]interface{}{
					"user_id":    sess.UserID,
					"is_offline": true,
				})),
				worker.QueueLow,
				asynq.MaxRetry(1),
				asynq.Retention(time.Hour*4),
			)

		}
	}

	ws.Logger.Debugf("%s is disconnected", sessionID)

}

func (ws *WS) handleMessage(s *melody.Session, msg []byte) {
	sessionID, found := s.Get(UserSessionIDContextKey)
	if !found {
		ws.Logger.Debugf("User %s not found in sessions", sessionID)
		return
	}
	str, ok := sessionID.(string)
	if !ok {
		ws.Logger.Debugf("%s is disconnected", sessionID)
		return
	}
	if sess := SessionID(str).GetSession(); sess != nil {
		ws.Logger.Debugf("handle message msg=%s", msg)
		var message Message
		var err = json.Unmarshal(msg, &message)
		if err != nil {
			ws.Logger.Debugf("handle message msg=%s err=%+v", msg, err)
			return
		}
		message.UserID = sess.UserID
		message.Role = enums.Role(sess.Role)
		err = ws.messageHandlerFunc(message)
		if err != nil {
			ws.Logger.Debugf("handle message msg=%s err=%+v", msg, err)
			return
		}
	}

}
