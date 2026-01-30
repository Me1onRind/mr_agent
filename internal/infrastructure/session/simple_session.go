package session

import (
	"context"

	"github.com/Me1onRind/mr_agent/internal/pkg/strutil"
)

type simpleSession struct {
	data map[string]*Data
}

func newSimpleSession() *simpleSession {
	s := &simpleSession{
		data: make(map[string]*Data),
	}
	return s
}

func (s *simpleSession) GetData(ctx context.Context, sessionId string) (*Data, error) {
	sessionData, ok := s.data[sessionId]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return sessionData, nil
}

func (s *simpleSession) SaveData(ctx context.Context, sessionId string, data *Data) error {
	s.data[sessionId] = data
	return nil
}

func (s *simpleSession) NewSessionId(ctx context.Context) string {
	return strutil.NewUUID()
}
