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

func (sr *SessionRepo) CheckToken(ctx context.Context, token string) (time.Time, error) {
	var expirationTime time.Time
	err := sr.db.QueryRowContext(ctx, `
	SELECT expiration_time
	FROM sessions
	WHERE token = $1 AND expiration_time > NOW()
	LIMIT 1;
	`, token).Scan(&expirationTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, util.ErrNotFound
		}
		return time.Time{}, err
	}

	return expirationTime, nil
}
