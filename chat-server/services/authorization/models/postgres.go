package models

type (
	PostgresUser struct {
		ID        int64  `db:"id"`
		Login     string `db:"login"`
		Email     string `db:"email"`
		Password  string `db:"password"`
		Client    string `db:"client"`
		CreatedAt int64  `db:"created_at"`
	}
	PostgresSession struct {
		ID        int64  `db:"id"`
		UserID    int64  `db:"user_id"`
		Token     string `db:"token"`
		Client    string `db:"client"`
		ExpiresAt int64  `db:"expires_at"`
		CreatedAt int64  `db:"created_at"`
	}
)
