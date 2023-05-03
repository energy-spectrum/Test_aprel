package db

import (
	"app/internal/util"
	"context"
	"database/sql"
	"errors"
	"time"
)

type AuthAuditRepo struct {
	db *sql.DB
}

type AuditEvent struct {
	Datatime time.Time
	Event    EventType
}

type EventType string

const (
	Login           EventType = "login"
	InvalidPassword EventType = "invalid_password"
	Block           EventType = "block"
)

func (aar *AuthAuditRepo) WriteEvent(ctx context.Context, userID int64, event EventType) error {
	_, err := aar.db.ExecContext(ctx, `
	INSERT INTO auth_audit (
		user_id,
		event,
		event_time
	) VALUES (
		$1, $2, $3
	);
    `, userID, event, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (aar *AuthAuditRepo) GetEvents(ctx context.Context, userID int64) ([]AuditEvent, error) {
	rows, err := aar.db.QueryContext(ctx, `
	SELECT * 
	FROM auth_audit
	WHERE user_id = $1
	ORDER BY event_time
    `, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()

	events := make([]AuditEvent, 0)
	for rows.Next() {
		var i AuditEvent
		if err := rows.Scan(
			&i.Datatime,
			&i.Event,
		); err != nil {
			return nil, err
		}
		events = append(events, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (aar *AuthAuditRepo) ClearAuditByUserID(ctx context.Context, userID int64) error {
	_, err := aar.db.ExecContext(ctx, `
	DELETE FROM auth_audit
	WHERE user_id = $1
    `, userID)
	if err != nil {
		return err
	}

	return nil
}
