package lib

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gox7/notify/services/authorization/internal/repository"
	"github.com/gox7/notify/services/authorization/models"
)

type (
	SessionsService struct {
		db *repository.Postgres
	}
)

func NewSessions(db *repository.Postgres, model *SessionsService) {
	*model = SessionsService{
		db: db,
	}
}

func (svc *SessionsService) CreateSession(userId int64, token string, client string, expires int64) (int64, error) {
	if userId == 0 {
		return 0, errors.New("userId is zero")
	}
	if client == "" {
		client = "none"
	}

	lastId, err := svc.db.RegisterSession(userId, expires, token, client)
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (svc *SessionsService) SearchSession(token string) (*models.PostgresSession, error) {
	session, err := svc.db.SearchSession(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	if session.ExpiresAt < time.Now().Unix() {
		svc.RemoveSession(session.Token)
		return nil, errors.New("token is expired")
	}

	return session, nil
}

func (svc *SessionsService) RemoveSession(token string) error {
	err := svc.db.RemoveSession(token)
	return err
}
