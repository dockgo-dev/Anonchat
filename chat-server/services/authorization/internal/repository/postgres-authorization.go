package repository

import (
	"time"

	"github.com/gox7/notify/services/authorization/models"
)

func (p *Postgres) RegisterUser(login, email, password, client string) (int64, error) {
	var id int64
	err := p.conn.QueryRow("INSERT INTO users (login, email, password, client, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		login, email, password, client, time.Now().Unix()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *Postgres) SearchUserByLogin(login string) (*models.PostgresUser, error) {
	var user models.PostgresUser
	err := p.conn.Get(&user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) SearchUserByEmail(email string) (*models.PostgresUser, error) {
	var user models.PostgresUser
	err := p.conn.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) SearchUserByID(id int64) (*models.PostgresUser, error) {
	var user models.PostgresUser
	err := p.conn.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) RemoveUserByID(id int64) error {
	_, err := p.conn.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
