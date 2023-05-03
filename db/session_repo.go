package db

import (
	"context"
	"database/sql"
	"time"

	"app/internal/util"
)

type SessionRepo struct {
	db *sql.DB
}

func (sr *SessionRepo) SaveToken(ctx context.Context, token string, expirationTime time.Time) error {
	_, err := sr.db.Exec(`
	INSERT INTO sessions (token, expiration_time)
	VALUES ($1, $2)
	`, token, expirationTime)

	return err
}

func (sr *SessionRepo) CheckToken(ctx context.Context, token string) (bool, error) {
	var isValid bool
	err := sr.db.QueryRowContext(ctx, `
	SELECT EXISTS (
		SELECT 1
		FROM sessions
		WHERE token = $1 AND expiration_time > NOW()
	);
	`, token).Scan(&isValid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, util.ErrNotFound
		}
		return false, err
	}

	return isValid, nil
}
