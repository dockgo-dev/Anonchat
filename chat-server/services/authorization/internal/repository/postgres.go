package repository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gox7/notify/services/authorization/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type (
	Postgres struct {
		conn   *sqlx.DB
		logger *slog.Logger
	}
)

func NewPostgres(config *models.LocalConfig, logger *slog.Logger, model *Postgres) {
	pcs := fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable",
		config.Postgres.User, config.Postgres.Addr, config.Postgres.DB,
	)

	var connect *sqlx.DB
	var err error
	for range 10 {
		connect, err = sqlx.Connect("postgres", pcs)
		if err == nil {
			logger.Info("postgres connect:", slog.String("addr", config.Postgres.Addr))
			fmt.Println("[+] postgres.connect:", config.Postgres.Addr)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		logger.Error(err.Error(), slog.String("addr", config.Postgres.Addr))
		fmt.Println("[-] postgres.connect:", err.Error())
		return
	}

	err = connect.Ping()
	if err != nil {
		logger.Error(err.Error(), slog.String("addr", config.Postgres.Addr))
		fmt.Println("[-] postgres.ping:", err.Error())
		return
	}

	logger.Info("postgres ping:", slog.String("addr", config.Postgres.Addr))
	fmt.Println("[+] postgres.ping:", config.Postgres.Addr)

	*model = Postgres{
		conn:   connect,
		logger: logger,
	}
}

func (p *Postgres) Migration() {
	_, err := p.conn.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		client VARCHAR(255) DEFAULT 'none',
		created_at INTEGER
    )`)
	if err != nil {
		fmt.Println("[-] create table:", err.Error())
	}

	_, err = p.conn.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		token VARCHAR(255) UNIQUE NOT NULL,
    	client VARCHAR(255) NOT NULL,
    	expires_at INTEGER NOT NULL,
    	created_at INTEGER NOT NULL
	)`)
	if err != nil {
		fmt.Println("[-] create table:", err.Error())
	}
}
