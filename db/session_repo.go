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

func (sr *SessionRepo) SaveToken(ctx context.Context, token string, expirationTime time.Time, userID int) error {
	_, err := sr.db.Exec(`
	INSERT INTO sessions (token, expiration_time, user_id)
	VALUES ($1, $2, $3)
	`, token, expirationTime, userID)

	return err
}

func (sr *SessionRepo) GetUserID(ctx context.Context, token string) (int, error) {
	var userID int
	err := sr.db.QueryRowContext(ctx, `
	SELECT user_id
	FROM sessions
	WHERE token = $1 AND expiration_time > NOW()
	LIMIT 1;
	`, token).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, util.ErrNotFound
		}
		return 0, err
	}

	return userID, nil
}
