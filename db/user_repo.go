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

func (ur *UserRepo) GetUserIDAndBlocked(ctx context.Context, login, hashedPassword string) (int, bool, error) {
	var userID int
	var correctPassword bool
	var blocked bool
	err := ur.db.QueryRowContext(ctx, `
	SELECT
		id,
		blocked,
		password = $2 AS correct_password
	FROM users
	WHERE login = $1
	ORDER BY id DESC
	LIMIT 1;
    `, login, hashedPassword).Scan(&userID, &blocked, &correctPassword)
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

func (ur *UserRepo) IncrementFailedLoginAttempts(ctx context.Context, id int) (int, error) {
	var failedLoginAttempts int
	err := ur.db.QueryRowContext(ctx, `
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

func (ur *UserRepo) Block(ctx context.Context, id int) error {
	_, err := ur.db.ExecContext(ctx, `
	UPDATE users
	SET blocked = true
	WHERE id = $1;
	`, id)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) Unblock(ctx context.Context, id int) error {
	_, err := ur.db.ExecContext(ctx, `
	UPDATE users
	SET blocked = false
	WHERE id = $1;
	`, id)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) Create(ctx context.Context, login, password string) error {
	_, err := ur.db.ExecContext(ctx, `
	INSERT INTO users (login, password)
	VALUES ($1, $2)
	`, login, password)

	return err
}
