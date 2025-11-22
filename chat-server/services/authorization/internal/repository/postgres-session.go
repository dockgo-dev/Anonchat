package repository

import (
	"time"

	"github.com/gox7/notify/services/authorization/models"
)

func (p *Postgres) RegisterSession(userId, expires int64, token, client string) (int64, error) {
	var id int64
	err := p.conn.QueryRow("INSERT INTO sessions (user_id, token, client, expires_at, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		userId, token, client, expires, time.Now().Unix()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *Postgres) SearchSession(token string) (*models.PostgresSession, error) {
	var session models.PostgresSession
	err := p.conn.Get(&session, "SELECT * FROM sessions WHERE token = $1", token)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (p *Postgres) RemoveSession(token string) error {
	_, err := p.conn.Exec("DELETE FROM sessions WHERE token = $1", token)
	return err
}
