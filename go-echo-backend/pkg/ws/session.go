package ws

import (
	"fmt"
	"strings"
)

// UserSessionIDContextKey session key
const UserSessionIDContextKey = "user_session_id"
const separator = "__"
const UserConnectedAtKey = "user_connected_at"

// Session session
type Session struct {
	UserID        string      `json:"user_id"`
	Role          string      `json:"role"`
	RemoteAddress string      `json:"remote_address"`
	ConnectedAt   interface{} `json:"connected_at"`
}

// SessionID id
type SessionID string

// GetID id
func (session *Session) GetID() SessionID {
	var sessionID = fmt.Sprintf("%s%s%s%s%s", session.UserID, separator, session.Role, separator, session.RemoteAddress)
	return SessionID(sessionID)
}
func (id SessionID) String() string {
	return string(id)
}

// GetSession get session
func (id SessionID) GetSession() *Session {
	var keys = strings.Split(string(id), separator)

	return &Session{
		UserID:        keys[0],
		Role:          keys[1],
		RemoteAddress: keys[2],
	}
}
