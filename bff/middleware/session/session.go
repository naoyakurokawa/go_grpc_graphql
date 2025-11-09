package session

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

const CookieName = "bff_session"

type Manager struct {
	mu       sync.RWMutex
	sessions map[string]uint64
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]uint64),
	}
}

func (m *Manager) CreateSession(userID uint64) string {
	sessionID := uuid.NewString()
	m.mu.Lock()
	m.sessions[sessionID] = userID
	m.mu.Unlock()
	return sessionID
}

func (m *Manager) GetUserID(sessionID string) (uint64, bool) {
	m.mu.RLock()
	userID, ok := m.sessions[sessionID]
	m.mu.RUnlock()
	return userID, ok
}

func (m *Manager) DeleteSession(sessionID string) {
	m.mu.Lock()
	delete(m.sessions, sessionID)
	m.mu.Unlock()
}

type contextKey string

const (
	userIDContextKey    contextKey = "session_user_id"
	sessionIDContextKey contextKey = "session_id"
	echoContextKey      contextKey = "echo_context"
)

func Middleware(sm *Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), echoContextKey, c)

			if cookie, err := c.Cookie(CookieName); err == nil {
				if userID, ok := sm.GetUserID(cookie.Value); ok {
					ctx = context.WithValue(ctx, userIDContextKey, userID)
					ctx = context.WithValue(ctx, sessionIDContextKey, cookie.Value)
				}
			}

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	val := ctx.Value(userIDContextKey)
	if val == nil {
		return 0, false
	}
	userID, ok := val.(uint64)
	return userID, ok
}

func SessionIDFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(sessionIDContextKey)
	if val == nil {
		return "", false
	}
	sessionID, ok := val.(string)
	return sessionID, ok
}

func EchoContextFromContext(ctx context.Context) (echo.Context, bool) {
	val := ctx.Value(echoContextKey)
	if val == nil {
		return nil, false
	}
	c, ok := val.(echo.Context)
	return c, ok
}

func SetSessionCookie(ctx context.Context, sessionID string) error {
	c, ok := EchoContextFromContext(ctx)
	if !ok {
		return echo.ErrInternalServerError
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	}

	// TODO: set Secure=true in production environments.
	c.SetCookie(cookie)
	return nil
}

func ClearSessionCookie(ctx context.Context) error {
	c, ok := EchoContextFromContext(ctx)
	if !ok {
		return echo.ErrInternalServerError
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	c.SetCookie(cookie)
	return nil
}
