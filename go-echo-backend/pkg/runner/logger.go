package runner

import "github.com/engineeringinflow/inflow-backend/pkg/logger"

type Logger struct {
	*logger.Logger
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.WithSkipCaller(1).Debugf(format, args...)
}
