package db

import (
	"context"
	"database/sql"
	"time"
)

type AuthAuditRepo struct {
	db *sql.DB
}

type AuthAuditEvent struct {
	Datetime time.Time
	Event    EventType
}

type EventType string

const (
	Login           EventType = "login"
	InvalidPassword EventType = "invalid_password"
	Block           EventType = "block"
)

func (aar *AuthAuditRepo) WriteEvent(ctx context.Context, userID int, event EventType) error {
	_, err := aar.db.ExecContext(ctx, `
	INSERT INTO auth_audit (
		user_id,
		event
	) VALUES (
		$1, $2
	);
    `, userID, event)
	if err != nil {
		return err
	}

	return nil
}

func (aar *AuthAuditRepo) GetAuthAuditByUserID(ctx context.Context, userID int) ([]AuthAuditEvent, error) {
	rows, err := aar.db.QueryContext(ctx, `
	SELECT event, event_time
	FROM auth_audit
	WHERE user_id = $1
	ORDER BY event_time
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []AuthAuditEvent
	for rows.Next() {
		var item AuthAuditEvent
		if err := rows.Scan(
			&item.Event,
			&item.Datetime,
		); err != nil {
			return nil, err
		}
		events = append(events, item)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (aar *AuthAuditRepo) ClearAuthAuditByUserID(ctx context.Context, userID int) error {
	_, err := aar.db.ExecContext(ctx, `
	DELETE FROM auth_audit
	WHERE user_id = $1
    `, userID)
	if err != nil {
		return err
	}

	return nil
}
