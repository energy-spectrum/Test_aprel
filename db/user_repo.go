package db

import (
	"app/internal/util"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepo struct {
	db *sql.DB
}

func (ur *UserRepo) GetUserIDAndBlocked(ctx context.Context, login, password string) (int64, bool, error) {
	var userID int64
	var correctPassword bool
	var blocked bool
	err := ur.db.QueryRow(`
	SELECT
		id,
		blocked,
		password = $2 AS correct_password
	FROM users
	WHERE login = $1
	ORDER BY id DESC
	LIMIT 1;
    `, login, password).Scan(&userID, &blocked, &correctPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, util.ErrNotFound
		}
		return 0, false, err
	}
	if !correctPassword {
		return userID, false, util.ErrInvalidPassword
	}

	return userID, blocked, nil
}

func (ur *UserRepo) IncrementFailedLoginAttempts(ctx context.Context, id int64) (int, error) {
	var failedLoginAttempts int
	err := ur.db.QueryRow(`
	UPDATE users
	SET failed_login_attempts = failed_login_attempts + 1
	WHERE id = $1
	RETURNING failed_login_attempts;
	`, id).Scan(&failedLoginAttempts)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}

	return failedLoginAttempts, nil
}


func (ur *UserRepo) Block(ctx context.Context, id int64) error {
	_, err := ur.db.ExecContext(ctx,`
	UPDATE users
	SET blocked = true
	WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("failed to block user by id: %d: %v", id, err)
	}

	return nil
}
