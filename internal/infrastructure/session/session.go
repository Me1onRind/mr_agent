package session

import (
	"context"
	"errors"
)

var (
	sessionStore       sessionStoreIface
	ErrStoreNotInit    = errors.New("session store not init")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionIdEmtpy  = errors.New("session_id is empty")
)

type sessionIdKey struct{}

type Data struct {
	SessionId string `json:"session_id"`
}

func InitSessionStore(ctx context.Context) error {
	sessionStore = newSimpleSession()
	return nil
}

func GetData(ctx context.Context) (*Data, error) {
	if sessionStore == nil {
		return nil, ErrStoreNotInit
	}
	sessionId, err := GetSessionId(ctx)
	if err != nil {
		return nil, err
	}
	return sessionStore.GetData(ctx, sessionId)
}

func SaveData(ctx context.Context, data *Data) error {
	if sessionStore == nil {
		return ErrStoreNotInit
	}
	sessionId, err := GetSessionId(ctx)
	if err != nil {
		return err
	}
	return sessionStore.SaveData(ctx, sessionId, data)
}

func NewSession(ctx context.Context, data *Data) (context.Context, error) {
	if sessionStore == nil {
		return ctx, ErrStoreNotInit
	}
	sessionId := sessionStore.NewSessionId(ctx)
	if err := sessionStore.SaveData(ctx, sessionId, data); err != nil {
		return ctx, err
	}
	return WithSessionId(ctx, sessionId), nil
}

func WithSessionId(ctx context.Context, sessionId string) context.Context {
	return context.WithValue(ctx, sessionIdKey{}, sessionId)
}

func GetSessionId(ctx context.Context) (string, error) {
	value := ctx.Value(sessionIdKey{})
	if value == nil {
		return "", ErrSessionIdEmtpy
	}
	return value.(string), nil
}

type sessionStoreIface interface {
	GetData(ctx context.Context, sessionId string) (*Data, error)
	SaveData(ctx context.Context, sessionId string, data *Data) error
	NewSessionId(ctx context.Context) string
}
